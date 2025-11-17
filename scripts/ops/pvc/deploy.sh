#!/bin/bash

# Build and Deploy Script for PVC Expansion Controller
set -e

# Configuration
REGISTRY="your-registry.com"  # Change this to your registry
IMAGE_NAME="pvc-expansion-controller"
TAG="${TAG:-latest}"
FULL_IMAGE="${REGISTRY}/${IMAGE_NAME}:${TAG}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    print_info "Checking prerequisites..."
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker is required but not installed"
        exit 1
    fi
    
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl is required but not installed"
        exit 1
    fi
    
    if ! kubectl cluster-info &> /dev/null; then
        print_error "Cannot connect to Kubernetes cluster"
        exit 1
    fi
    
    print_info "Prerequisites check passed"
}

# Build the Go application
build_binary() {
    print_info "Building Go binary..."
    
    # Ensure we have go.mod
    if [[ ! -f "go.mod" ]]; then
        go mod init pvc-expansion-controller
        go mod tidy
    fi
    
    # Build locally first to check for compilation errors
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o manager controller.go
    
    if [[ $? -eq 0 ]]; then
        print_info "Go binary built successfully"
        rm -f manager  # Clean up local binary
    else
        print_error "Failed to build Go binary"
        exit 1
    fi
}

# Build Docker image
build_image() {
    print_info "Building Docker image: ${FULL_IMAGE}"
    
    docker build -t "${FULL_IMAGE}" .
    
    if [[ $? -eq 0 ]]; then
        print_info "Docker image built successfully"
    else
        print_error "Failed to build Docker image"
        exit 1
    fi
}

# Push Docker image
push_image() {
    print_info "Pushing Docker image: ${FULL_IMAGE}"
    
    docker push "${FULL_IMAGE}"
    
    if [[ $? -eq 0 ]]; then
        print_info "Docker image pushed successfully"
    else
        print_error "Failed to push Docker image"
        exit 1
    fi
}

# Deploy CRD
deploy_crd() {
    print_info "Deploying CRD..."
    
    kubectl apply -f PVCExpansionRequest.yaml
    
    # Wait for CRD to be established
    print_info "Waiting for CRD to be established..."
    kubectl wait --for=condition=Established crd/pvcexpansionrequests.cloudness.io --timeout=60s
    
    print_info "CRD deployed and established"
}

# Deploy controller
deploy_controller() {
    print_info "Deploying controller..."
    
    # Update image in deployment manifest
    sed "s|your-registry/pvc-expansion-controller:latest|${FULL_IMAGE}|g" controller-deployment.yaml | kubectl apply -f -
    
    # Wait for deployment to be ready
    print_info "Waiting for controller deployment to be ready..."
    kubectl wait --for=condition=Available deployment/pvc-expansion-controller -n pvc-expansion-system --timeout=300s
    
    print_info "Controller deployment is ready"
}

# Main deployment function
main() {
    print_info "Starting PVC Expansion Controller deployment..."
    print_info "Image: ${FULL_IMAGE}"
    
    check_prerequisites
    build_binary
    build_image
    
    # Ask if user wants to push
    read -p "Push image to registry? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        push_image
    else
        print_warn "Skipping image push. Make sure the image is available in your cluster."
    fi
    
    deploy_crd
    deploy_controller
    
    print_info "✅ Deployment completed successfully!"
    print_info ""
    print_info "Check controller status:"
    print_info "kubectl get pods -n pvc-expansion-system"
    print_info ""
    print_info "View controller logs:"
    print_info "kubectl logs -n pvc-expansion-system deployment/pvc-expansion-controller"
    print_info ""
    print_info "Create a test PVCExpansionRequest:"
    print_info "kubectl apply -f - <<EOF"
    print_info "apiVersion: cloudness.io/v1alpha1"
    print_info "kind: PVCExpansionRequest"
    print_info "metadata:"
    print_info "  name: test-expansion"
    print_info "spec:"
    print_info "  statefulSetRef:"
    print_info "    name: my-statefulset"
    print_info "  size: \"100Gi\""
    print_info "EOF"
}

# Parse command line arguments
case "${1:-}" in
    "build")
        check_prerequisites
        build_binary
        build_image
        ;;
    "deploy")
        check_prerequisites
        deploy_crd
        deploy_controller
        ;;
    "")
        main
        ;;
    *)
        echo "Usage: $0 [build|deploy]"
        echo "  build  - Only build the image"
        echo "  deploy - Only deploy to cluster (assumes image exists)"
        echo "  (none) - Full build and deploy"
        exit 1
        ;;
esac
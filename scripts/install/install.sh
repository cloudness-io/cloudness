#!/bin/bash

# Cloudness Platform Prerequisites Installation Script
# This script installs and configures all the dependencies required for Cloudness Platform

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
INSTALL_GATEWAY_API="${INSTALL_GATEWAY_API:-true}"
INSTALL_CERT_MANAGER="${INSTALL_CERT_MANAGER:-true}"
INSTALL_TRAEFIK="${INSTALL_TRAEFIK:-true}"
INSTALL_KUBEBLOCKS="${INSTALL_KUBEBLOCKS:-true}"
INSTALL_HTTP_ROUTE="${INSTALL_HTTP_ROUTE:-true}"
CERT_MANAGER_VERSION="${CERT_MANAGER_VERSION:-v1.18.2}"
GATEWAY_API_VERSION="${GATEWAY_API_VERSION:-v1.3.0}"
KUBEBLOCKS_VERSION="${KUBEBLOCKS_VERSION:-v1.0.1}"
TRAEFIK_VERSION="${TRAEFIK_VERSION:-37.1.1}"
VERBOSE="${VERBOSE:-false}"

# Installation URLs and paths
BASE_URL="https://get.cloudness.in"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Detect if running from curl (no local files available)
CURL_INSTALL=false
if [ ! -f "${SCRIPT_DIR}/cloudness-namespace.yaml" ] || [ ! -f "${SCRIPT_DIR}/cloudness-rbac.yaml" ] || [ ! -f "${SCRIPT_DIR}/traefik-gateway.yaml" ]; then
    CURL_INSTALL=true
fi

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --verbose|-v)
            VERBOSE="true"
            shift
            ;;
        --help|-h)
            echo "Cloudness Platform Prerequisites Installer"
            echo ""
            echo "This script installs all required Kubernetes prerequisites for the Cloudness Platform:"
            echo ""
            echo "INFRASTRUCTURE COMPONENTS:"
            echo "• Gateway API CRDs (v${GATEWAY_API_VERSION}) - Traffic routing and management"
            echo "• cert-manager (v${CERT_MANAGER_VERSION}) - Automated TLS certificate management"
            echo "• Traefik Gateway (v${TRAEFIK_VERSION}) - Ingress controller and load balancer"
            echo "• KubeBlocks (v${KUBEBLOCKS_VERSION}) - Database and application management platform"
            echo ""
            echo "PLATFORM COMPONENTS:"
            echo "• Cloudness namespace - Application deployment namespace"
            echo "• Cloudness RBAC - Service account and cluster role for platform operations"
            echo "• Traefik Gateway - Gateway API Gateway and GatewayClass configuration"
            echo ""
            echo "INSTALLATION METHODS:"
            echo ""
            echo "Local installation (with cloned repository):"
            echo "  ./install-prerequisites.sh [OPTIONS]"
            echo ""
            echo "Remote installation (using curl):"
            echo "  curl -fsSL https://cdm.cloudness.io/cloudness/install-prerequisites.sh | bash"
            echo "  curl -fsSL https://cdm.cloudness.io/cloudness/install-prerequisites.sh | bash -s -- --verbose"
            echo ""
            echo "Options:"
            echo "  -v, --verbose    Show full output from kubectl and helm commands"
            echo "  -h, --help       Show this help message"
            echo ""
            echo "Environment variables:"
            echo "  INSTALL_GATEWAY_API     Install Gateway API CRDs (default: true)"
            echo "  INSTALL_CERT_MANAGER    Install Cert-Manager (default: true)"
            echo "  INSTALL_TRAEFIK         Install Traefik (default: true)"
            echo "  INSTALL_KUBEBLOCKS      Install KubeBlocks (default: true)"
            echo "  INSTALL_HTTP_ROUTE      Install Http Route (default: true)"
            echo "  CERT_MANAGER_VERSION    Cert-Manager version (default: ${CERT_MANAGER_VERSION})"
            echo "  GATEWAY_API_VERSION     Gateway API version (default: ${GATEWAY_API_VERSION})"
            echo "  KUBEBLOCKS_VERSION      KubeBlocks version (default: ${KUBEBLOCKS_VERSION})"
            echo "  TRAEFIK_VERSION         Traefik version (default: ${TRAEFIK_VERSION})"
            echo "  VERBOSE                 Show full output (default: false)"
            echo ""
            echo "Examples:"
            echo "  $0                      # Install with default settings"
            echo "  $0 --verbose            # Install with full output"
            echo "  VERBOSE=true $0         # Install with full output (env var)"
            echo ""
            echo "NETWORK ACCESS:"
            echo "  After installation, Traefik will be available on:"
            echo "  • HTTP: Port 30080 (NodePort)"
            echo "  • HTTPS: Port 30443 (NodePort)"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information."
            exit 1
            ;;
    esac
done

# Function to run commands with optional output
run_command() {
    if [ "$VERBOSE" = "true" ]; then
        "$@"
    else
        "$@" > /dev/null 2>&1
    fi
}

# Function to get YAML content (local file or download)
get_yaml_content() {
    local filename="$1"
    local filepath="${SCRIPT_DIR}/${filename}"
    
    if [ "$CURL_INSTALL" = "true" ]; then
        # Download from remote URL
        curl -fsSL "${BASE_URL}/${filename}"
    elif [ -f "$filepath" ]; then
        # Use local file
        cat "$filepath"
    else
        echo "Error: Cannot find $filename locally and not in curl mode" >&2
        exit 1
    fi
}

# Function to apply YAML (local file or download)
apply_yaml() {
    local filename="$1"
    local description="$2"
    
    print_info "Applying ${description}..."
    
    local apply_result
    if [ "$CURL_INSTALL" = "true" ]; then
        # Download and apply directly
        if [ "$VERBOSE" = "true" ]; then
            apply_result=$(curl -fsSL "${BASE_URL}/${filename}" | kubectl apply -f -)
        else
            apply_result=$(curl -fsSL "${BASE_URL}/${filename}" | kubectl apply -f - 2>/dev/null)
        fi
    else
        # Use local file
        local filepath="${SCRIPT_DIR}/${filename}"
        if [ -f "$filepath" ]; then
            if [ "$VERBOSE" = "true" ]; then
                apply_result=$(kubectl apply -f "$filepath")
            else
                apply_result=$(kubectl apply -f "$filepath" 2>/dev/null)
            fi
        else
            print_error "Local file not found: $filepath"
            exit 1
        fi
    fi
    
    if [ $? -eq 0 ]; then
        if [ "$VERBOSE" = "true" ]; then
            echo "$apply_result"
        fi
        print_status "${description} applied successfully"
    else
        print_error "Failed to apply ${description}"
        exit 1
    fi
}

# Function to add and update Helm repositories
add_helm_repo() {
    local repo_name="$1"
    local repo_url="$2"
    
    if ! run_command helm repo add "$repo_name" "$repo_url"; then
        print_warning "$repo_name repo already exists, continuing..."
    fi
    
    if ! run_command helm repo update; then
        print_error "Failed to update Helm repositories. Please check your network connection."
        exit 1
    fi
}

echo -e "${BLUE}🔧 Installing Cloudness Platform Prerequisites${NC}"
echo "=============================================="

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to print status
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "$1"
}

print_section() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Check prerequisites
echo "Checking prerequisites..."

if ! command_exists kubectl; then
    print_error "kubectl is not installed. Please install kubectl first."
    exit 1
fi

if ! command_exists helm; then
    print_error "helm is not installed. Please install Helm 3.x first."
    exit 1
fi

# Check Kubernetes connection
if ! kubectl cluster-info >/dev/null 2>&1; then
    print_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
    exit 1
fi

# Check if we have sufficient permissions
print_info "Checking cluster permissions..."

# Test cluster-level permissions
if ! kubectl auth can-i create clusterroles >/dev/null 2>&1; then
    print_error "Current kubeconfig does not have permission to create ClusterRoles. Please ensure you have cluster-admin privileges."
    exit 1
fi

if ! kubectl auth can-i create clusterrolebindings >/dev/null 2>&1; then
    print_error "Current kubeconfig does not have permission to create ClusterRoleBindings. Please ensure you have cluster-admin privileges."
    exit 1
fi

if ! kubectl auth can-i create customresourcedefinitions >/dev/null 2>&1; then
    print_error "Current kubeconfig does not have permission to create CustomResourceDefinitions. Please ensure you have cluster-admin privileges."
    exit 1
fi

# Test namespace-level permissions
if ! kubectl auth can-i create namespaces >/dev/null 2>&1; then
    print_error "Current kubeconfig does not have permission to create namespaces. Please ensure you have cluster-admin privileges."
    exit 1
fi

# Test if we can create resources in system namespaces (required for cert-manager and traefik)
if ! kubectl auth can-i create deployments --namespace=cert-manager >/dev/null 2>&1; then
    print_warning "Limited permissions detected for cert-manager namespace. Installation may require elevated privileges."
fi

if ! kubectl auth can-i create deployments --namespace=traefik >/dev/null 2>&1; then
    print_warning "Limited permissions detected for traefik namespace. Installation may require elevated privileges."
fi

print_status "Prerequisites check passed"

# Display installation mode
if [ "$CURL_INSTALL" = "true" ]; then
    print_info "🌐 Remote installation mode detected - downloading resources from ${BASE_URL}"
else
    print_info "📁 Local installation mode detected - using files from ${SCRIPT_DIR}"
fi

echo ""
print_info "Starting Cloudness Platform Prerequisites Installation..."
print_info "This will install Gateway API, cert-manager, Traefik, and Cloudness platform resources"
echo ""

# Install Gateway API CRDs
if [ "$INSTALL_GATEWAY_API" = "true" ]; then
    print_info "Installing Gateway API CRDs ${GATEWAY_API_VERSION}..."
    
    if ! run_command kubectl apply -f "https://github.com/kubernetes-sigs/gateway-api/releases/download/${GATEWAY_API_VERSION}/standard-install.yaml"; then
        print_error "Failed to apply Gateway API CRDs. Please check your network connection and kubectl configuration."
        exit 1
    fi
    
    # Wait for CRDs to be established
    if ! run_command kubectl wait --for condition=established --timeout=60s crd/gatewayclasses.gateway.networking.k8s.io || \
       ! run_command kubectl wait --for condition=established --timeout=60s crd/gateways.gateway.networking.k8s.io || \
       ! run_command kubectl wait --for condition=established --timeout=60s crd/httproutes.gateway.networking.k8s.io; then
        print_error "Gateway API CRDs did not become ready in time. Please check your cluster's health."
        exit 1
    fi
    
    print_status "Gateway API CRDs installed"
else
    print_warning "Skipping Gateway API CRDs installation"
fi

# Install Cert-Manager
if [ "$INSTALL_CERT_MANAGER" = "true" ]; then
    print_info "Installing Cert-Manager ${CERT_MANAGER_VERSION}..."
    
    # Add Jetstack Helm repository
    print_info "Adding Jetstack Helm repository..."
    add_helm_repo "jetstack" "https://charts.jetstack.io"
    
    # Install cert-manager CRDs
    print_info "Installing cert-manager CRDs..."
    if ! run_command kubectl apply -f "https://github.com/cert-manager/cert-manager/releases/download/${CERT_MANAGER_VERSION}/cert-manager.crds.yaml"; then
        print_error "Failed to apply cert-manager CRDs. Please check your network connection."
        exit 1
    fi
    
    # Create cert-manager namespace
    if ! kubectl create namespace cert-manager --dry-run=client -o yaml | run_command kubectl apply -f -; then
        print_error "Failed to create cert-manager namespace."
        exit 1
    fi
    
    # Install cert-manager
    print_info "Deploying cert-manager via Helm (this may take a few minutes)..."
    if ! run_command helm upgrade --install cert-manager jetstack/cert-manager \
        --namespace cert-manager \
        --version "${CERT_MANAGER_VERSION}" \
        --set crds.keep=true \
        --wait \
        --timeout 5m; then
        print_error "Cert-Manager installation failed. Please check your Kubernetes cluster and run the script again."
        exit 1
    fi
    
    print_status "Cert-Manager installed"
else
    print_warning "Skipping Cert-Manager installation"
fi

# Install Traefik
if [ "$INSTALL_TRAEFIK" = "true" ]; then
    print_info "Installing Traefik ${TRAEFIK_VERSION}..."
    
    # Add Traefik Helm repository
    print_info "Adding Traefik Helm repository..."
    add_helm_repo "traefik" "https://traefik.github.io/charts"
    
    # Create traefik namespace
    if ! kubectl create namespace traefik --dry-run=client -o yaml | run_command kubectl apply -f -; then
        print_error "Failed to create traefik namespace."
        exit 1
    fi
    
    # Install Traefik with Gateway API support
    print_info "Deploying Traefik with Gateway API support (this may take a few minutes)..."
    if ! run_command helm upgrade --install traefik traefik/traefik \
        --namespace traefik \
        --version "${TRAEFIK_VERSION}" \
        --set providers.kubernetesGateway.enabled=true \
        --set providers.kubernetesCRD.allowCrossNamespace=true \
        --set gateway.enabled=false \
        --set ports.web.port=8000 \
        --set ports.web.exposedPort=80 \
        --set ports.websecure.port=8443 \
        --set ports.websecure.exposedPort=443 \
        --set ports.websecure.tls.enabled=true \
        --set service.type=LoadBalancer \
        --set securityContext.capabilities.drop="{ALL}" \
        --set securityContext.readOnlyRootFilesystem=true \
        --set securityContext.runAsGroup=65532 \
        --set securityContext.runAsNonRoot=true \
        --set securityContext.runAsUser=65532 \
        --wait \
        --timeout 5m; then
        print_error "Traefik installation failed. Please check your Kubernetes cluster and run the script again."
        exit 1
    fi
    
    print_status "Traefik installed"
    
    # Apply Traefik Gateway configuration
    apply_yaml "traefik-gateway.yaml" "Traefik Gateway and GatewayClass configuration"
    
else
    print_warning "Skipping Traefik installation"
fi

# Install KubeBlocks
if [ "$INSTALL_KUBEBLOCKS" = "true" ]; then
    print_info "Installing KubeBlocks ${KUBEBLOCKS_VERSION}..."
    
    # Install KubeBlocks CRDs
    print_info "Installing KubeBlocks CRDs..."
    if ! run_command kubectl apply --server-side=true --validate=false -f "https://github.com/apecloud/kubeblocks/releases/download/${KUBEBLOCKS_VERSION}/kubeblocks_crds.yaml"; then
        print_error "Failed to apply KubeBlocks CRDs. Please check your network connection."
        exit 1
    fi
    
    # Add KubeBlocks Helm repository
    print_info "Adding KubeBlocks Helm repository..."
    add_helm_repo "kubeblocks" "https://apecloud.github.io/helm-charts"
    
    # Create kubeblocks-system namespace
    if ! kubectl create namespace kubeblocks-system --dry-run=client -o yaml | run_command kubectl apply -f -; then
        print_error "Failed to create kubeblocks-system namespace."
        exit 1
    fi
    
    # Install KubeBlocks
    print_info "Deploying KubeBlocks (this may take a few minutes)..."
    if ! run_command helm upgrade --install kubeblocks kubeblocks/kubeblocks \
        --namespace kubeblocks-system \
        --version "${KUBEBLOCKS_VERSION}" \
        --wait \
        --timeout 10m; then
        print_error "KubeBlocks installation failed. Please check your Kubernetes cluster and run the script again."
        exit 1
    fi
    
    print_status "KubeBlocks installed"
else
    print_warning "Skipping KubeBlocks installation"
fi

# Verify installations
echo ""
print_info "Verifying installations..."

if [ "$INSTALL_GATEWAY_API" = "true" ]; then
    if kubectl get crd gatewayclasses.gateway.networking.k8s.io >/dev/null 2>&1; then
        print_status "Gateway API CRDs are installed"
    else
        print_error "Gateway API CRDs installation failed"
    fi
fi

if [ "$INSTALL_CERT_MANAGER" = "true" ]; then
    print_info "Waiting for cert-manager pods to be ready..."
    if run_command kubectl wait --for=condition=ready pod --all -n cert-manager --timeout=60s >/dev/null 2>&1; then
        print_status "Cert-Manager is running"
    else
        print_warning "Cert-Manager pods may not be fully ready yet. Check with: kubectl get pods -n cert-manager"
    fi
fi

if [ "$INSTALL_TRAEFIK" = "true" ]; then
    print_info "Waiting for traefik pods to be ready..."
    if run_command kubectl wait --for=condition=ready pod --all -n traefik --timeout=60s >/dev/null 2>&1; then
        print_status "Traefik is running"
    else
        print_warning "Traefik pods may not be fully ready yet. Check with: kubectl get pods -n traefik"
    fi
fi

if [ "$INSTALL_KUBEBLOCKS" = "true" ]; then
    print_info "Waiting for KubeBlocks pods to be ready..."
    if run_command kubectl wait --for=condition=ready pod --all -n kubeblocks-system --timeout=120s >/dev/null 2>&1; then
        print_status "KubeBlocks is running"
    else
        print_warning "KubeBlocks pods may not be fully ready yet. Check with: kubectl get pods -n kubeblocks-system"
    fi
fi

# Install Cloudness Platform Resources
install_cloudness_resources() {
    print_section "Installing Cloudness Platform Resources"
    
    # Apply Cloudness namespace
    apply_yaml "cloudness-namespace.yaml" "Cloudness namespace"
    
    # Apply Cloudness RBAC resources
    apply_yaml "cloudness-rbac.yaml" "Cloudness RBAC resources"

    apply_yaml "cloudness-runner-rbac.yaml" "Cloudness Runner RBAC resource"
    
    # Verify RBAC resources
    print_info "Verifying Cloudness RBAC resources..."
    
    # Check all resources silently, only report errors
    if ! run_command kubectl get clusterrole cloudness-runner-role &>/dev/null; then
        print_error "ClusterRole 'cloudness-runner-role' not found"
        exit 1
    fi
    
    if ! run_command kubectl get serviceaccount cloudness-runner-account -n cloudness &>/dev/null; then
        print_error "ServiceAccount 'cloudness-runner-account' not found"
        exit 1
    fi
    
    if ! run_command kubectl get clusterrolebinding cloudness-runner-role-binding &>/dev/null; then
        print_error "ClusterRoleBinding 'cloudness-runner-role-binding' not found"
        exit 1
    fi
    
    if ! run_command kubectl get clusterrole cloudness-app-role &>/dev/null; then
        print_error "ClusterRole 'cloudness-app-role' not found"
        exit 1
    fi
    
    if ! run_command kubectl get serviceaccount cloudness-app-account -n cloudness &>/dev/null; then
        print_error "ServiceAccount 'cloudness-app-account' not found"
        exit 1
    fi
    
    if ! run_command kubectl get clusterrolebinding cloudness-app-role-binding &>/dev/null; then
        print_error "ClusterRoleBinding 'cloudness-app-role-binding' not found"
        exit 1
    fi

    # Apply Cloudness Postgres cluster
    apply_yaml "cloudness-postgres.yaml" "Cloudness Postgres Cluster"

    # Apply Cloudness Postgres cluster
    apply_yaml "cloudness-redis.yaml" "Cloudness Redis Cluster"
    
    # Verify Database resources
    print_info "Verifying Cloudness Database resources..."
    
    # Wait for Postgres cluster to be ready
    if ! run_command kubectl wait --for=jsonpath='{.status.components.postgresql.phase}'=Running cluster/pg-cluster --namespace=cloudness --timeout=120s >/dev/null 2>&1; then
        print_warning "Cloudness Postgres cluster may not be fully ready yet. Check with: kubectl get cluster pg-cluster -n cloudness"
    else
        print_status "Cloudness Postgres cluster is ready"
    fi
    
    # Wait for Redis cluster to be ready
    if ! run_command kubectl wait --for=jsonpath='{.status.components.redis.phase}'=Running cluster/redis-cluster --namespace=cloudness --timeout=120s  >/dev/null 2>&1; then
        print_warning "Cloudness Redis cluster may not be fully ready yet. Check with: kubectl get cluster redis-cluster -n cloudness"
    else
        print_status "Cloudness Redis cluster is ready"
    fi
    if ! run_command kubectl wait --for=jsonpath='{.status.components.redis-sentinel.phase}'=Running cluster/redis-cluster --namespace=cloudness --timeout=120s  >/dev/null 2>&1; then
        print_warning "Cloudness Redis Sentinel cluster may not be fully ready yet. Check with: kubectl get cluster redis-cluster -n cloudness"
    else
        print_status "Cloudness Redis Sentinel is ready"
    fi


    # Verify Database resources
    print_info "Setting up Cloudness Database artifacts (DB, triggers, indexes)..."
    PGPASSWORD=$(kubectl get secrets -n cloudness "pg-cluster-postgresql-account-postgres" -o jsonpath='{.data.password}' | base64 -d)
    export PGPASSWORD

    # Get the primary pod name to exec into it
    PRIMARY_POD=$(kubectl get pods -n cloudness -l app.kubernetes.io/instance=pg-cluster,kubeblocks.io/role=primary -o jsonpath='{.items[0].metadata.name}')
    if [ -z "$PRIMARY_POD" ]; then
        print_error "Could not determine primary pod for Cloudness Postgres cluster. Exiting..."
        exit 1
    fi

    # 1. Ensure the dblink extension is enabled in the 'postgres' database where we are executing the script
    print_info "Ensuring dblink extension is available..."
    kubectl exec -ti -n cloudness "$PRIMARY_POD" -- psql -U postgres -d postgres -c "CREATE EXTENSION IF NOT EXISTS dblink;"

    if [ $? -ne 0 ]; then
        print_error "Failed to enable dblink extension. Check permissions or configuration."
        exit 1
    fi

    kubectl exec -ti -n cloudness "$PRIMARY_POD" -- psql -U postgres -d postgres -c "
        DO \$do\$
        BEGIN
            IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'cloudness') THEN
            -- We can now use dblink_exec because the extension is loaded
            PERFORM dblink_exec('host=localhost user=postgres password=$PGPASSWORD dbname=postgres', 'CREATE DATABASE cloudness');
            END IF;
        END
        \$do\$;
    "
    # Apply Cloudness Postgres cluster
    apply_yaml "cloudness-app.yaml" "Cloudness Application"

    if run_command kubectl wait --for=condition=Available deployment cloudness -n cloudness --timeout=120s >/dev/null 2>&1; then
        print_status "Cloudness Application is running"
    else
        print_warning "Cloudness Application pods may not be fully available yet. Check with: kubectl get pods -n cloudness"
    fi

    # Get the ip address from traefik service
    TRAEFIK_IP=$(kubectl get svc traefik -n traefik -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    if [ -z "$TRAEFIK_IP" ]; then
        print_warning "Could not determine Traefik IP address. Check with: kubectl get svc traefik -n traefik"
    else
        print_status "Traefik IP address: $TRAEFIK_IP"
    fi

    if [ "$INSTALL_HTTP_ROUTE" = "true" ]; then
        # Apply Cloudness HTTPRoute with { { .ServiceDomain.Domain } } as cloudness.${TRAEFIK_IP}.sslip.io
        print_info "Applying Cloudness HTTPRoute configuration..."
        CLOUDNESS_DOMAIN="cloudness.${TRAEFIK_IP}.sslip.io"
        CLOUDNESS_HTTP_ROUTE_YAML=$(get_yaml_content "cloudness-http-route.yaml" | sed "s/{{ .ServiceDomain.Domain }}/${CLOUDNESS_DOMAIN}/g")
        echo "$CLOUDNESS_HTTP_ROUTE_YAML" | kubectl apply -f -
        if [ $? -ne 0 ]; then
            print_error "Failed to apply Cloudness HTTPRoute configuration"
            exit 1
        fi
    fi

    print_status "All Cloudness platform resources are available"
}

# Install Cloudness platform resources
install_cloudness_resources

echo ""
print_status "Prerequisites installation completed!"
echo ""
print_info "Next steps:"
echo "1. Configure your domain and TLS certificates"
echo "2. Run the main installation: ./install-platform.sh"
echo ""
print_info "To check the status of services:"
echo "• Gateway API: kubectl get gatewayclasses"
echo "• Traefik Gateway: kubectl get gateway -n traefik"
echo "• Cert-Manager: kubectl get pods -n cert-manager"
echo "• Traefik: kubectl get pods -n traefik"
echo "• Traefik Service: kubectl get svc -n traefik"
echo "• KubeBlocks: kubectl get pods -n kubeblocks-system"
echo "• Cloudness Namespace: kubectl get namespace cloudness"
echo "• Cloudness RBAC: kubectl get clusterrole cloudness-runner-role"

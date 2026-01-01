#!/bin/sh

# Cloudness Kubernetes Deployment Script
# This script deploys applications to Kubernetes with proper resource management

set -eu

# ==============================================================================
# Configuration & Constants
# ==============================================================================

# Timeout configurations (can be overridden via environment)
readonly ROLLOUT_TIMEOUT_STATELESS="${ROLLOUT_TIMEOUT_STATELESS:-60s}"
readonly ROLLOUT_TIMEOUT_STATEFUL="${ROLLOUT_TIMEOUT_STATEFUL:-120s}"
readonly PVC_RESIZE_TIMEOUT="${PVC_RESIZE_TIMEOUT:-300}"
readonly PVC_RESIZE_POLL_INTERVAL="${PVC_RESIZE_POLL_INTERVAL:-5}"

# Required environment variables with defaults
: "${CLOUDNESS_DEPLOY_APP_IDENTIFIER:=}"
: "${CLOUDNESS_DEPLOY_APP_NAMESPACE:=}"
: "${CLOUDNESS_DEPLOY_FLAG_APP_TYPE:=}"
: "${CLOUDNESS_DEPLOY_FLAG_HAS_VOLUME:=0}"
: "${CLOUDNESS_DEPLOY_FLAG_NEED_REMOUNT:=0}"
: "${CLOUDNESS_DEPLOY_FLAG_HAS_ROUTE:=0}"
: "${CLOUDNESS_DEPLOY_PATH:=}"
: "${VERBOSE:=false}"

# YAML file paths (set from CLOUDNESS_DEPLOY_PATH)
CLOUDNESS_DEPLOY_YAML_COMMON=""
CLOUDNESS_DEPLOY_YAML_VOLUME=""
CLOUDNESS_DEPLOY_YAML_APP=""
CLOUDNESS_DEPLOY_YAML_ROUTE=""

# Track cleanup state
CLEANUP_DONE=false

# ==============================================================================
# Validation Functions
# ==============================================================================

validate_dependencies() {
    missing=""

    if ! command_exists kubectl; then
        missing="$missing kubectl"
    fi

    if ! command_exists yq; then
        missing="$missing yq"
    fi

    if [ -n "$missing" ]; then
        log_error "Missing required dependencies:$missing"
        log_error "Please install them before running this script."
        return 1
    fi

    return 0
}

validate_environment() {
    has_errors=0

    if [ -z "$CLOUDNESS_DEPLOY_APP_IDENTIFIER" ]; then
        log_error "  - CLOUDNESS_DEPLOY_APP_IDENTIFIER is required"
        has_errors=1
    fi

    if [ -z "$CLOUDNESS_DEPLOY_APP_NAMESPACE" ]; then
        log_error "  - CLOUDNESS_DEPLOY_APP_NAMESPACE is required"
        has_errors=1
    fi

    if [ -z "$CLOUDNESS_DEPLOY_FLAG_APP_TYPE" ]; then
        log_error "  - CLOUDNESS_DEPLOY_FLAG_APP_TYPE is required"
        has_errors=1
    elif [ "$CLOUDNESS_DEPLOY_FLAG_APP_TYPE" != "Stateless" ] && [ "$CLOUDNESS_DEPLOY_FLAG_APP_TYPE" != "Stateful" ]; then
        log_error "  - CLOUDNESS_DEPLOY_FLAG_APP_TYPE must be 'Stateless' or 'Stateful'"
        has_errors=1
    fi

    if [ -z "$CLOUDNESS_DEPLOY_PATH" ]; then
        log_error "  - CLOUDNESS_DEPLOY_PATH is required"
        has_errors=1
    fi

    if [ "$has_errors" -eq 1 ]; then
        log_error "Environment validation failed"
        return 1
    fi

    # Set YAML file paths from deploy path (mounted from ConfigMap)
    CLOUDNESS_DEPLOY_YAML_COMMON="$CLOUDNESS_DEPLOY_PATH/common.yaml"
    CLOUDNESS_DEPLOY_YAML_VOLUME="$CLOUDNESS_DEPLOY_PATH/volume.yaml"
    CLOUDNESS_DEPLOY_YAML_APP="$CLOUDNESS_DEPLOY_PATH/app.yaml"
    CLOUDNESS_DEPLOY_YAML_ROUTE="$CLOUDNESS_DEPLOY_PATH/route.yaml"

    return 0
}

# ==============================================================================
# Kubernetes Operations
# ==============================================================================

# Returns the resource type based on app type
get_resource_type() {
    if [ "$CLOUDNESS_DEPLOY_FLAG_APP_TYPE" = "Stateless" ]; then
        echo "deployment"
    else
        echo "statefulset"
    fi
}

# Returns the opposite resource type (for cleanup)
get_opposite_resource_type() {
    if [ "$CLOUDNESS_DEPLOY_FLAG_APP_TYPE" = "Stateless" ]; then
        echo "statefulset"
    else
        echo "deployment"
    fi
}

# Apply Kubernetes configuration from YAML file
kube_apply() {
    yaml_file="$1"
    error_output=""

    # Skip if file doesn't exist or is empty
    if [ ! -f "$yaml_file" ] || [ ! -s "$yaml_file" ]; then
        return 0
    fi

    # Apply the YAML
    if [ "$VERBOSE" = "true" ]; then
        if ! kubectl apply -f "$yaml_file"; then
            log_error "Failed to apply Kubernetes configuration from $yaml_file"
            return 1
        fi
    else
        if ! error_output=$(kubectl apply -f "$yaml_file" 2>&1 >/dev/null); then
            log_error "Failed to apply Kubernetes configuration:"
            log_error "$error_output"
            return 1
        fi
    fi

    return 0
}

# Delete a Kubernetes resource
kube_delete() {
    resource_type="$1"
    resource_name="$2"
    namespace="$3"
    error_output=""

    if [ "$VERBOSE" = "true" ]; then
        if ! kubectl delete "$resource_type/$resource_name" -n "$namespace" --ignore-not-found=true; then
            log_warn "Failed to delete $resource_type/$resource_name"
            return 1
        fi
    else
        if ! error_output=$(kubectl delete "$resource_type/$resource_name" -n "$namespace" --ignore-not-found=true 2>&1 >/dev/null); then
            log_warn "Failed to delete $resource_type/$resource_name: $error_output"
            return 1
        fi
    fi

    return 0
}

# Wait for rollout to complete
kube_rollout_status() {
    resource_type=""
    timeout=""
    error_output=""

    resource_type=$(get_resource_type)

    if [ "$resource_type" = "deployment" ]; then
        timeout="$ROLLOUT_TIMEOUT_STATELESS"
    else
        timeout="$ROLLOUT_TIMEOUT_STATEFUL"
    fi

    if [ "$VERBOSE" = "true" ]; then
        if ! kubectl rollout status "$resource_type/$CLOUDNESS_DEPLOY_APP_IDENTIFIER" \
            -n "$CLOUDNESS_DEPLOY_APP_NAMESPACE" \
            --timeout="$timeout"; then
            log_error "Rollout failed, reverting..."
            kubectl rollout undo "$resource_type/$CLOUDNESS_DEPLOY_APP_IDENTIFIER" \
                -n "$CLOUDNESS_DEPLOY_APP_NAMESPACE" || true
            return 1
        fi
    else
        if ! error_output=$(kubectl rollout status "$resource_type/$CLOUDNESS_DEPLOY_APP_IDENTIFIER" \
            -n "$CLOUDNESS_DEPLOY_APP_NAMESPACE" \
            --timeout="$timeout" 2>&1); then
            log_error "$error_output"
            log_error "Rollout failed, reverting..."
            kubectl rollout undo "$resource_type/$CLOUDNESS_DEPLOY_APP_IDENTIFIER" \
                -n "$CLOUDNESS_DEPLOY_APP_NAMESPACE" 2>/dev/null || true
            return 1
        fi
    fi

    return 0
}

# Parse storage size to numeric GiB value
parse_size_to_gib() {
    size_str="$1"
    echo "$size_str" | sed 's/Gi//'
}

# Wait for PVC to resize
kube_wait_pvc_resize() {
    pvc_name="$1"
    new_size="$2"
    namespace="$CLOUDNESS_DEPLOY_APP_NAMESPACE"
    deadline=""
    current_time=""

    current_time=$(date +%s)
    deadline=$((current_time + PVC_RESIZE_TIMEOUT))

    while true; do
        # Check PVC status
        pvc_status=""
        pvc_status=$(kubectl get pvc "$pvc_name" -n "$namespace" -o jsonpath='{.status.phase}' 2>/dev/null || echo "")

        if [ "$VERBOSE" = "true" ]; then
            log_info "PVC '$pvc_name' status: $pvc_status"
        fi

        # Handle WaitForFirstConsumer
        if [ "$pvc_status" = "Pending" ]; then
            pvc_event=""
            pvc_event=$(kubectl get events -n "$namespace" \
                --field-selector "involvedObject.kind=PersistentVolumeClaim,involvedObject.name=$pvc_name" \
                --sort-by=.lastTimestamp \
                -o jsonpath='{.items[-1:].reason}' 2>/dev/null || echo "")
            if [ "$pvc_event" = "WaitForFirstConsumer" ]; then
                return 0
            fi
        fi

        # Check if resize completed
        current_size=""
        current_size=$(kubectl get pvc "$pvc_name" -n "$namespace" -o jsonpath='{.status.capacity.storage}' 2>/dev/null || echo "0Gi")

        if [ "$VERBOSE" = "true" ]; then
            log_info "PVC '$pvc_name' current size: $current_size, target: $new_size"
        fi
        if [ "$(parse_size_to_gib "$current_size")" -ge "$(parse_size_to_gib "$new_size")" ]; then
            return 0
        fi

        # Check for FileSystemResizePending condition
        resize_pending=""
        resize_pending=$(kubectl get pvc "$pvc_name" -n "$namespace" \
            -o jsonpath='{.status.conditions[?(@.type=="FileSystemResizePending")].status}' 2>/dev/null || echo "")
        if [ "$resize_pending" = "True" ]; then
            log_info "Volume resized. Remounting application to finalize."
            return 0
        fi

        # Check timeout
        current_time=$(date +%s)
        if [ "$current_time" -ge "$deadline" ]; then
            log_error "Timed out after ${PVC_RESIZE_TIMEOUT}s waiting for PVC '$pvc_name' to reach $new_size"
            return 1
        fi

        log_info "Waiting for volume '$pvc_name'..."
        sleep "$PVC_RESIZE_POLL_INTERVAL"
    done
}

# ==============================================================================
# Deployment Functions
# ==============================================================================

deploy_common_artifacts() {

    if ! kube_apply "$CLOUDNESS_DEPLOY_YAML_COMMON"; then
        log_error "Failed to set up prerequisite artifacts"
        return 1
    fi

    log_step "Prerequisite artifacts configured"
    return 0
}

deploy_volume() {
    if [ "$CLOUDNESS_DEPLOY_FLAG_HAS_VOLUME" -ne 1 ]; then
        return 0
    fi

    # Handle remount for volume resize
    if [ "$CLOUDNESS_DEPLOY_FLAG_NEED_REMOUNT" -eq 1 ]; then
        log_info "Volume resize detected, removing statefulset for remount..."
        if ! kube_delete "statefulset" "$CLOUDNESS_DEPLOY_APP_IDENTIFIER" "$CLOUDNESS_DEPLOY_APP_NAMESPACE"; then
            log_error "Failed to remove statefulset for remount"
            return 1
        fi
    fi

    # Apply volume configuration
    if ! kube_apply "$CLOUDNESS_DEPLOY_YAML_VOLUME"; then
        log_error "Failed to apply volume configuration"
        return 1
    fi

    # Wait for each PVC to be ready (read from file)
    pvc_data=""
    pvc_data=$(yq -r 'select(.kind == "PersistentVolumeClaim") | .metadata.name + " " + .spec.resources.requests.storage' "$CLOUDNESS_DEPLOY_YAML_VOLUME" 2>/dev/null || echo "")

    echo "$pvc_data" | while read -r pvc_name new_size; do
        # Skip empty lines
        pvc_name=$(echo "$pvc_name" | xargs)
        new_size=$(echo "$new_size" | xargs)
        if [ -z "$pvc_name" ] || [ -z "$new_size" ] || [ "$pvc_name" = "---" ]; then
            continue
        fi

        if ! kube_wait_pvc_resize "$pvc_name" "$new_size"; then
            log_error "Failed to provision PVC '$pvc_name'"
            return 1
        fi
    done

    log_step "Volumes provisioned"
    return 0
}

deploy_application() {
    if ! kube_apply "$CLOUDNESS_DEPLOY_YAML_APP"; then
        log_error "Failed to deploy application"
        return 1
    fi

    if ! kube_rollout_status; then
        return 1
    fi

    log_step "Application deployed"
    return 0
}

deploy_routes() {
    if [ "$CLOUDNESS_DEPLOY_FLAG_HAS_ROUTE" -ne 1 ]; then
        return 0
    fi

    if ! kube_apply "$CLOUDNESS_DEPLOY_YAML_ROUTE"; then
        log_error "Failed to configure HTTP routes"
        return 1
    fi

    log_step "HTTP routes configured"
    return 0
}

# ==============================================================================
# Lifecycle Management
# ==============================================================================

cleanup() {
    if [ "$CLEANUP_DONE" = "true" ]; then
        return 0
    fi
    CLEANUP_DONE=true

    log_info "Running cleanup..."

    opposite_type=""
    opposite_type=$(get_opposite_resource_type)

    kube_delete "$opposite_type" "$CLOUDNESS_DEPLOY_APP_IDENTIFIER" "$CLOUDNESS_DEPLOY_APP_NAMESPACE" || true
}

on_exit() {
    exit_code=$?
    if [ $exit_code -ne 0 ]; then
        log_error "Deployment failed with exit code $exit_code"
    fi
    cleanup
    exit $exit_code
}

# ==============================================================================
# Main Entrypoint
# ==============================================================================

main() {
    # Set up exit trap
    trap on_exit EXIT
    
    print_section "Deploying application"

    # Validate prerequisites
    if ! validate_dependencies; then
        exit 1
    fi

    if ! validate_environment; then
        exit 1
    fi

    # Execute deployment steps
    if ! deploy_common_artifacts; then
        exit 1
    fi

    if ! deploy_volume; then
        exit 1
    fi

    if ! deploy_application; then
        exit 1
    fi

    if ! deploy_routes; then
        exit 1
    fi

    log_success "Deployment completed successfully!"
}

# Run main function
main "$@"

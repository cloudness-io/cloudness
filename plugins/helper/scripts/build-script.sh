#!/bin/sh

# Cloudness Build Script
# This script handles container image building using Dockerfile or Nixpacks

set -eu

# ==============================================================================
# Configuration
# ==============================================================================

# Build type: "dockerfile" or "nixpacks"
: "${CLOUDNESS_BUILD_TYPE:=}"

# Source paths
: "${CLOUDNESS_BUILD_SOURCE_PATH:=}"
: "${CLOUDNESS_BUILD_DOCKERFILE:=Dockerfile}"

# Image configuration
: "${CLOUDNESS_BUILD_IMAGE:=}"
: "${CLOUDNESS_BUILD_CACHE_IMAGE:=}"

# Registry configuration
: "${CLOUDNESS_IMAGE_REGISTRY:=}"
: "${CLOUDNESS_IMAGE_MIRROR_REGISTRY:=}"
: "${CLOUDNESS_MIRROR_ENABLED:=false}"

# Nixpacks specific
: "${CLOUDNESS_BUILD_CMD:=}"
: "${CLOUDNESS_START_CMD:=}"

# Build args (JSON or space-separated key=value pairs)
: "${CLOUDNESS_BUILD_ARGS:=}"

# ==============================================================================
# Validation
# ==============================================================================

validate_inputs() {
    has_errors=0

    if [ -z "$CLOUDNESS_BUILD_TYPE" ]; then
        log_error "CLOUDNESS_BUILD_TYPE is required (dockerfile or nixpacks)"
        has_errors=1
    elif [ "$CLOUDNESS_BUILD_TYPE" != "dockerfile" ] && [ "$CLOUDNESS_BUILD_TYPE" != "nixpacks" ]; then
        log_error "CLOUDNESS_BUILD_TYPE must be 'dockerfile' or 'nixpacks'"
        has_errors=1
    fi

    if [ -z "$CLOUDNESS_BUILD_SOURCE_PATH" ]; then
        log_error "CLOUDNESS_BUILD_SOURCE_PATH is required"
        has_errors=1
    fi

    if [ -z "$CLOUDNESS_BUILD_IMAGE" ]; then
        log_error "CLOUDNESS_BUILD_IMAGE is required"
        has_errors=1
    fi

    if [ "$has_errors" -eq 1 ]; then
        return 1
    fi

    return 0
}

# ==============================================================================
# BuildKit Configuration
# ==============================================================================

setup_buildkit_config() {

    printf "\n"
    BUILDKITD_CONFIG_PATH="$HOME/.config/buildkit/buildkitd.toml"
    mkdir -p "$(dirname "$BUILDKITD_CONFIG_PATH")"
    > "$BUILDKITD_CONFIG_PATH"

    # Main registry configuration
    MAIN_REGISTRY=$(echo "$CLOUDNESS_IMAGE_REGISTRY" | cut -d'/' -f1)
    
    cat >> "$BUILDKITD_CONFIG_PATH" << EOF
[registry."$MAIN_REGISTRY"]
  http = true
  insecure = true
EOF

    # Mirror registry configuration
    if [ "$CLOUDNESS_MIRROR_ENABLED" = "true" ] && [ -n "$CLOUDNESS_IMAGE_MIRROR_REGISTRY" ]; then
        MIRROR_REGISTRY=$(echo "$CLOUDNESS_IMAGE_MIRROR_REGISTRY" | cut -d'/' -f1)
        
        cat >> "$BUILDKITD_CONFIG_PATH" << EOF

[registry."$MIRROR_REGISTRY"]
  http = true
  insecure = true

[registry."docker.io"]
  mirrors = ["$CLOUDNESS_IMAGE_MIRROR_REGISTRY"]
EOF
    fi
}

# ==============================================================================
# Build Functions
# ==============================================================================

build_with_dockerfile() {
    print_section "Building with Dockerfile"

    log_info "Dockerfile: $CLOUDNESS_BUILD_DOCKERFILE"

    # Build base command
    build_cmd="buildctl-daemonless.sh build \
        --frontend=dockerfile.v0 \
        --local context=$CLOUDNESS_BUILD_SOURCE_PATH \
        --local dockerfile=$CLOUDNESS_BUILD_SOURCE_PATH \
        --opt filename=$CLOUDNESS_BUILD_DOCKERFILE \
        --output type=image,name=$CLOUDNESS_BUILD_IMAGE,push=true"

    # Add cache configuration
    if [ -n "$CLOUDNESS_BUILD_CACHE_IMAGE" ]; then
        build_cmd="$build_cmd \
            --export-cache type=registry,ref=$CLOUDNESS_BUILD_CACHE_IMAGE,mode=max \
            --import-cache type=registry,ref=$CLOUDNESS_BUILD_CACHE_IMAGE,mode=max"
    fi

    # Add build args
    if [ -n "$CLOUDNESS_BUILD_ARGS" ]; then
        for arg in $CLOUDNESS_BUILD_ARGS; do
            build_cmd="$build_cmd --opt build-arg:$arg"
        done
    fi

    # Execute build
    log_info "Starting build..."
    if ! eval "$build_cmd"; then
        log_error "Dockerfile build failed"
        return 1
    fi

    return 0
}

build_with_nixpacks() {
    print_section "Building with Nixpacks"

    # Build nixpacks command
    nixpacks_cmd="nixpacks build $CLOUDNESS_BUILD_SOURCE_PATH -o $CLOUDNESS_BUILD_SOURCE_PATH"
    nixpacks_cmd="$nixpacks_cmd --name $CLOUDNESS_BUILD_IMAGE"

    if [ -n "$CLOUDNESS_BUILD_CMD" ]; then
        log_info "Build command: $CLOUDNESS_BUILD_CMD"
        nixpacks_cmd="$nixpacks_cmd --build-cmd \"$CLOUDNESS_BUILD_CMD\""
    fi

    if [ -n "$CLOUDNESS_START_CMD" ]; then
        log_info "Start command: $CLOUDNESS_START_CMD"
        nixpacks_cmd="$nixpacks_cmd --start-cmd \"$CLOUDNESS_START_CMD\""
    fi

    # Add environment variables from build args
    if [ -n "$CLOUDNESS_BUILD_ARGS" ]; then
        for arg in $CLOUDNESS_BUILD_ARGS; do
            key=$(echo "$arg" | cut -d'=' -f1)
            value=$(echo "$arg" | cut -d'=' -f2-)
            nixpacks_cmd="$nixpacks_cmd --env $key=\"$value\""
        done
    fi

    nixpacks_cmd="$nixpacks_cmd"

    # Generate Dockerfile with Nixpacks
    log_info "Generating Dockerfile with Nixpacks..."
    if ! eval "$nixpacks_cmd"; then
        log_error "Nixpacks generation failed"
        return 1
    fi

    # Build and push with BuildKit
    build_cmd="buildctl-daemonless.sh build \
        --frontend=dockerfile.v0 \
        --local context=$CLOUDNESS_BUILD_SOURCE_PATH \
        --local dockerfile=$CLOUDNESS_BUILD_SOURCE_PATH \
        --opt filename=/.nixpacks/Dockerfile \
        --output type=image,name=$CLOUDNESS_BUILD_IMAGE,push=true"

    # Add cache configuration
    if [ -n "$CLOUDNESS_BUILD_CACHE_IMAGE" ]; then
        build_cmd="$build_cmd \
            --export-cache type=registry,ref=$CLOUDNESS_BUILD_CACHE_IMAGE \
            --import-cache type=registry,ref=$CLOUDNESS_BUILD_CACHE_IMAGE,mode=max"
    fi

    log_info "Building and pushing image..."
    if ! eval "$build_cmd"; then
        log_error "Image build/push failed"
        return 1
    fi

    log_step "Image built and pushed successfully"
    return 0
}

# ==============================================================================
# Main
# ==============================================================================

main() {
    if ! validate_inputs; then
        exit 1
    fi

    setup_buildkit_config

    case "$CLOUDNESS_BUILD_TYPE" in
        dockerfile)
            if ! build_with_dockerfile; then
                exit 1
            fi
            ;;
        nixpacks)
            if ! build_with_nixpacks; then
                exit 1
            fi
            ;;
    esac

    log_success "Build completed successfully!"
}

main "$@"

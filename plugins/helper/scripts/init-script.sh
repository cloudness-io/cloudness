#!/bin/sh

# Cloudness Git Clone Script
# This script handles git repository cloning for the build process

set -eu

# ==============================================================================
# Configuration
# ==============================================================================

: "${CLOUDNESS_GIT_REPO_URL:=}"
: "${CLOUDNESS_GIT_BRANCH:=}"
: "${CLOUDNESS_GIT_COMMIT:=}"
: "${CLOUDNESS_BUILD_PATH:=}"

# Optional netrc credentials
: "${GIT_MACHINE:=}"
: "${GIT_LOGIN:=}"
: "${GIT_PASSWORD:=}"

# ==============================================================================
# Validation
# ==============================================================================

validate_inputs() {
    has_errors=0

    if [ -z "$CLOUDNESS_GIT_REPO_URL" ]; then
        log_error "CLOUDNESS_GIT_REPO_URL is required"
        has_errors=1
    fi

    if [ -z "$CLOUDNESS_GIT_BRANCH" ]; then
        log_error "CLOUDNESS_GIT_BRANCH is required"
        has_errors=1
    fi

    if [ -z "$CLOUDNESS_BUILD_PATH" ]; then
        log_error "CLOUDNESS_BUILD_PATH is required"
        has_errors=1
    fi

    if [ "$has_errors" -eq 1 ]; then
        return 1
    fi

    return 0
}

# ==============================================================================
# Git Operations
# ==============================================================================

setup_netrc() {
    if [ -n "$GIT_MACHINE" ] && [ -n "$GIT_LOGIN" ] && [ -n "$GIT_PASSWORD" ]; then
        log_info "Configuring git credentials..."
        echo "machine $GIT_MACHINE login $GIT_LOGIN password $GIT_PASSWORD" > ~/.netrc
        chmod 600 ~/.netrc
    fi
}

clone_repository() {
    print_section "Cloning Repository"
    
    log_info "Repository: $CLOUDNESS_GIT_REPO_URL"
    log_info "Branch: $CLOUDNESS_GIT_BRANCH"
    
    if ! git clone "$CLOUDNESS_GIT_REPO_URL" --branch "$CLOUDNESS_GIT_BRANCH" "$CLOUDNESS_BUILD_PATH"; then
        log_error "Failed to clone repository"
        return 1
    fi

    # Checkout specific commit if provided
    if [ -n "$CLOUDNESS_GIT_COMMIT" ]; then
        log_info "Checking out commit: $CLOUDNESS_GIT_COMMIT"
        git -C "$CLOUDNESS_BUILD_PATH" config advice.detachedHead false
        if ! git -C "$CLOUDNESS_BUILD_PATH" checkout "$CLOUDNESS_GIT_COMMIT"; then
            log_error "Failed to checkout commit $CLOUDNESS_GIT_COMMIT"
            return 1
        fi
    fi

    return 0
}

# ==============================================================================
# Main
# ==============================================================================

main() {
    if ! validate_inputs; then
        exit 1
    fi

    setup_netrc
    
    if ! clone_repository; then
        exit 1
    fi

    log_success "Repository cloned successfully!"
}

main "$@"

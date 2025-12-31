#!/bin/sh
# Cloudness Shared Utilities
# Source this file in scripts: . /usr/local/lib/cloudness-utils.sh

# ==============================================================================
# Colors
# ==============================================================================
readonly CLOUDNESS_RED='\033[1;31m'
readonly CLOUDNESS_GREEN='\033[1;32m'
readonly CLOUDNESS_YELLOW='\033[1;33m'
readonly CLOUDNESS_BLUE='\033[0;34m'
readonly CLOUDNESS_NC='\033[0m'

# ==============================================================================
# Logging Functions
# ==============================================================================

log_error() {
    printf "${CLOUDNESS_RED}[ERROR]${CLOUDNESS_NC} %s\n" "$*" >&2
}

log_warn() {
    printf "${CLOUDNESS_YELLOW}[WARN]${CLOUDNESS_NC} %s\n" "$*"
}

log_info() {
    printf "%s\n" "$*"
}

log_success() {
    printf "${CLOUDNESS_GREEN}[SUCCESS]${CLOUDNESS_NC} %s\n" "$*"
}

log_step() {
    printf "%s ${CLOUDNESS_GREEN}✔${CLOUDNESS_NC}\n" "$*"
}

log_debug() {
    if [ "${VERBOSE:-false}" = "true" ]; then
        printf "${CLOUDNESS_BLUE}[DEBUG]${CLOUDNESS_NC} %s\n" "$*"
    fi
}

print_section() {
    printf "\n"
    printf "${CLOUDNESS_BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${CLOUDNESS_NC}\n"
    printf "${CLOUDNESS_BLUE}  %s${CLOUDNESS_NC}\n" "$1"
    printf "${CLOUDNESS_BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${CLOUDNESS_NC}\n"
}

# ==============================================================================
# Helper Functions
# ==============================================================================

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Run command with optional verbose output
run_command() {
    if [ "${VERBOSE:-false}" = "true" ]; then
        "$@"
    else
        "$@" > /dev/null 2>&1
    fi
}

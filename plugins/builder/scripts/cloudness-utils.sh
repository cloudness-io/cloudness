#!/bin/sh
# Cloudness Shared Utilities
# Source this file in scripts: . /usr/local/lib/cloudness-utils.sh

# ==============================================================================
# Colors
# ==============================================================================
readonly CLOUDNESS_RED='\033[1;31m'
readonly CLOUDNESS_GREEN='\033[1;32m'
readonly CLOUDNESS_YELLOW='\033[1;33m'
readonly CLOUDNESS_BLUE='\033[38;2;40;153;245m'
readonly CLOUDNESS_NC='\033[0m'

# ==============================================================================
# Logging Functions
# ==============================================================================

log_error() {
    printf "%b\n" "${CLOUDNESS_RED}❌ $*${CLOUDNESS_NC}" >&2
}

log_warn() {
    printf "%b\n" "${CLOUDNESS_YELLOW}⚠️ $*${CLOUDNESS_NC}"
}

log_info() {
    printf "%b\n" "$*"
}

log_success() {
    printf "%b\n" "${CLOUDNESS_GREEN}✔ $*${CLOUDNESS_NC}"
}

log_step() {
    printf "%b\n" "$* ${CLOUDNESS_GREEN}✔${CLOUDNESS_NC}"
}

log_debug() {
    if [ "${VERBOSE:-false}" = "true" ]; then
        printf "${CLOUDNESS_BLUE}[DEBUG]${CLOUDNESS_NC} %s\n" "$*"
    fi
}

print_section() {
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

#!/usr/bin/env bash
set -euo pipefail

# Security leak scanner for jira-mcp
# Scans current files and full git history for secrets and credentials.
#
# Usage:
#   ./scripts/security-scan.sh            # Scan everything (files + git history)
#   ./scripts/security-scan.sh --staged   # Scan only staged changes (for pre-commit)

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
GITLEAKS_VERSION="8.21.2"
GITLEAKS_BIN=""

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*"; }

# Detect OS and architecture
detect_platform() {
    local os arch
    os="$(uname -s)"
    arch="$(uname -m)"

    case "$os" in
        Linux)  os="linux" ;;
        Darwin) os="darwin" ;;
        *)      error "Unsupported OS: $os"; exit 1 ;;
    esac

    case "$arch" in
        x86_64)  arch="x64" ;;
        aarch64|arm64) arch="arm64" ;;
        *)       error "Unsupported architecture: $arch"; exit 1 ;;
    esac

    echo "${os}_${arch}"
}

# Install gitleaks if not available
ensure_gitleaks() {
    if command -v gitleaks &>/dev/null; then
        GITLEAKS_BIN="gitleaks"
        return
    fi

    local cache_dir="$REPO_ROOT/.cache/security"
    local cached_bin="$cache_dir/gitleaks"

    if [[ -x "$cached_bin" ]]; then
        GITLEAKS_BIN="$cached_bin"
        return
    fi

    info "Installing gitleaks v${GITLEAKS_VERSION}..."
    mkdir -p "$cache_dir"

    local platform
    platform="$(detect_platform)"
    local url="https://github.com/gitleaks/gitleaks/releases/download/v${GITLEAKS_VERSION}/gitleaks_${GITLEAKS_VERSION}_${platform}.tar.gz"

    if ! curl -sSfL "$url" | tar -xz -C "$cache_dir" gitleaks 2>/dev/null; then
        error "Failed to download gitleaks. Install manually: https://github.com/gitleaks/gitleaks#installing"
        exit 1
    fi

    chmod +x "$cached_bin"
    GITLEAKS_BIN="$cached_bin"
    info "Installed gitleaks to $cached_bin"
}

run_scan() {
    local mode="${1:-full}"
    local exit_code=0

    cd "$REPO_ROOT"

    info "Running gitleaks scan (mode: $mode)..."
    echo ""

    case "$mode" in
        staged)
            # Pre-commit mode: only scan staged changes
            $GITLEAKS_BIN protect --staged --config .gitleaks.toml --verbose || exit_code=$?
            ;;
        *)
            # Full scan: current files + entire git history
            info "=== Scanning git history ==="
            $GITLEAKS_BIN detect --source . --config .gitleaks.toml --verbose --report-format sarif --report-path /dev/null || exit_code=$?
            ;;
    esac

    echo ""
    if [[ $exit_code -eq 0 ]]; then
        info "No leaks detected."
    else
        error "Leaks detected! Please review and fix before committing."
        echo ""
        echo "If a finding is a false positive, add it to .gitleaks.toml allowlist."
    fi

    return $exit_code
}

main() {
    local mode="full"

    for arg in "$@"; do
        case "$arg" in
            --staged) mode="staged" ;;
            --help|-h)
                echo "Usage: $0 [--staged] [--help]"
                echo ""
                echo "  --staged   Scan only staged git changes (for pre-commit hooks)"
                echo "  (default)  Full scan of files and git history"
                exit 0
                ;;
        esac
    done

    ensure_gitleaks
    run_scan "$mode"
}

main "$@"

#!/usr/bin/env bash
set -euo pipefail

REPO_BASE="https://github.com/Infran/wc26"
APP_NAME="wc26"
VERSION="${WC26_VERSION:-latest}"
INSTALL_DIR="${WC26_INSTALL_DIR:-$HOME/.local/bin}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

info()  { echo -e "${CYAN}$1${NC}"; }
ok()    { echo -e "${GREEN}$1${NC}"; }
warn()  { echo -e "${YELLOW}$1${NC}"; }
err()   { echo -e "${RED}$1${NC}" >&2; }

# --- Detect platform ---
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) err "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# --- Prerequisites ---
if ! command -v go &>/dev/null; then
    warn "Go is not installed. Installing..."
    case "$OS" in
        linux)
            curl -fsSL "https://go.dev/dl/go1.24.0.linux-${ARCH}.tar.gz" | tar -C /usr/local -xz
            export PATH="/usr/local/go/bin:$PATH"
            ;;
        darwin)
            if command -v brew &>/dev/null; then
                brew install go
            else
                curl -fsSL "https://go.dev/dl/go1.24.0.darwin-${ARCH}.tar.gz" | tar -C /usr/local -xz
                export PATH="/usr/local/go/bin:$PATH"
            fi
            ;;
        *)
            err "Unsupported OS: $OS. Install Go manually: https://go.dev/dl/"
            exit 1
            ;;
    esac
    ok "Go installed: $(go version)"
else
    ok "Found: $(go version)"
fi

# --- Create install directory ---
mkdir -p "$INSTALL_DIR"

# --- Build from source ---
BUILD_DIR=$(mktemp -d)
trap 'rm -rf "$BUILD_DIR"' EXIT

# Try local source first
if [ -f "./go.mod" ] && grep -q "$APP_NAME" "./go.mod" 2>/dev/null; then
    info "Using local source at $(pwd)"
    cp -r . "$BUILD_DIR/src"
else
    info "Cloning $REPO_BASE ..."
    if command -v git &>/dev/null; then
        git clone --depth 1 "$REPO_BASE" "$BUILD_DIR/src"
    else
        warn "Git not installed. Trying 'go install' directly..."
        GOFLAGS="-ldflags=-X main.Version=$VERSION" \
            go install "$REPO_BASE/cmd/$APP_NAME@$VERSION"
        if command -v "$APP_NAME" &>/dev/null; then
            ok "✓ wc26 installed via 'go install'"
            # Post-install
            "$APP_NAME" config init 2>/dev/null || true
            echo ""
            ok "✓ wc26 installed successfully!"
            echo "  Run '$APP_NAME --help' to get started"
            exit 0
        fi
        err "Cannot clone or install. Install Git or clone manually."
        exit 1
    fi
fi

cd "$BUILD_DIR/src"
info "Building $APP_NAME..."
LDFLAGS="-X main.Version=$VERSION"
go build -ldflags "$LDFLAGS" -o "$INSTALL_DIR/$APP_NAME" "./cmd/$APP_NAME/"

# --- Add to PATH if needed ---
case ":$PATH:" in
    *:$INSTALL_DIR:*) ;;
    *)
        case "$SHELL" in
            *zsh*) echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$HOME/.zshrc" ;;
            *bash*) echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$HOME/.bashrc" ;;
            *) warn "Add $INSTALL_DIR to your PATH manually" ;;
        esac
        export PATH="$INSTALL_DIR:$PATH"
        ;;
esac

# --- Post-install ---
"$INSTALL_DIR/$APP_NAME" config init 2>/dev/null || true

echo ""
ok "✓ wc26 installed successfully!"
echo "  Binary: $INSTALL_DIR/$APP_NAME"
echo "  Run '$APP_NAME --help' to get started"
echo "  Run '$APP_NAME auth login <email> <password>' to authenticate"
echo "  Run '$APP_NAME health' to check API status"

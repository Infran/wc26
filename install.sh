#!/usr/bin/env bash
set -euo pipefail

REPO_BASE="https://github.com/Infran/wc26"
APP_NAME="wc26"
API_RELEASES="https://api.github.com/repos/Infran/wc26/releases"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

info()  { echo -e "${CYAN}$1${NC}"; }
ok()    { echo -e "${GREEN}$1${NC}"; }
warn()  { echo -e "${YELLOW}$1${NC}"; }
err()   { echo -e "${RED}$1${NC}" >&2; }

# --- Resolve version ---
if [ -n "${WC26_VERSION:-}" ]; then
    TAG="$WC26_VERSION"
else
    TAG=$(curl -sfL "$API_RELEASES/latest" | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": "\(.*\)",/\1/')
    if [ -z "$TAG" ]; then
        warn "Could not fetch latest version. Specify WC26_VERSION env var."
        TAG="latest"
    fi
fi
VERSION="${TAG#v}"

info "Installing $APP_NAME $TAG ..."

# --- Detect platform ---
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) err "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# --- Determine binary name ---
case "$OS" in
    linux) BINARY="${APP_NAME}_linux_${ARCH}" ;;
    darwin) BINARY="${APP_NAME}_darwin_${ARCH}" ;;
    *) err "Unsupported OS: $OS"; exit 1 ;;
esac

# --- Determine install directory ---
INSTALL_DIR="${WC26_INSTALL_DIR:-$HOME/.local/bin}"
mkdir -p "$INSTALL_DIR"
TARGET="$INSTALL_DIR/$APP_NAME"

DOWNLOADED=false
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# --- Try pre-built binary ---
BINARY_URL="$REPO_BASE/releases/download/$TAG/$BINARY"
CHECKSUMS_URL="$REPO_BASE/releases/download/$TAG/wc26-checksums.txt"
BIN_PATH="$TMP_DIR/$BINARY"

if command -v curl &>/dev/null; then
    info "Downloading binary..."
    HTTP_CODE=$(curl -sfL "$BINARY_URL" -o "$BIN_PATH" -w "%{http_code}" 2>/dev/null || echo "000")

    if [ "$HTTP_CODE" = "200" ]; then
        info "Downloading checksums..."
        CHECKSUMS_PATH="$TMP_DIR/wc26-checksums.txt"
        if curl -sfL "$CHECKSUMS_URL" -o "$CHECKSUMS_PATH" 2>/dev/null; then
            EXPECTED_HASH=$(grep "$BINARY" "$CHECKSUMS_PATH" | awk '{print $1}')
            if [ -n "$EXPECTED_HASH" ]; then
                ACTUAL_HASH=$(sha256sum "$BIN_PATH" | awk '{print $1}')
                if [ "$ACTUAL_HASH" != "$EXPECTED_HASH" ]; then
                    err "SHA256 mismatch!"
                    err "  Expected: $EXPECTED_HASH"
                    err "  Got:      $ACTUAL_HASH"
                    exit 1
                fi
                ok "SHA256 verified."
            else
                warn "Binary not found in checksums, skipping verification."
            fi
        else
            warn "Could not download checksums, skipping verification."
        fi

        chmod +x "$BIN_PATH"
        # Smoke test
        VERSION_OUT=$("$BIN_PATH" --version 2>&1)
        if [ $? -ne 0 ]; then
            err "Binary smoke test failed: $VERSION_OUT"
            exit 1
        fi

        cp "$BIN_PATH" "$TARGET"
        ok "Binary installed to: $TARGET"
        DOWNLOADED=true
    else
        warn "Pre-built binary not found for $OS/$ARCH (HTTP $HTTP_CODE)."
    fi
else
    warn "curl not available, skipping binary download."
fi

# --- Fallback: build from source ---
if [ "$DOWNLOADED" = false ]; then
    if ! command -v go &>/dev/null; then
        err "Go is required to build from source. Install Go from https://go.dev/dl/"
        exit 1
    fi

    warn "Building from source..."
    SRC_DIR="$TMP_DIR/src"

    if [ -f "./go.mod" ] && grep -q "$APP_NAME" "./go.mod" 2>/dev/null; then
        info "Using local source at $(pwd)"
        SRC_DIR="$PWD"
    elif command -v git &>/dev/null; then
        git clone --depth 1 "$REPO_BASE.git" "$SRC_DIR"
        cd "$SRC_DIR"
    else
        warn "Git not installed. Trying 'go install' directly..."
        GOFLAGS="-ldflags=-X main.Version=$VERSION" go install "$REPO_BASE/cmd/$APP_NAME@$VERSION"
        if command -v "$APP_NAME" &>/dev/null; then
            DOWNLOADED=true
            TARGET="$(command -v "$APP_NAME")"
        else
            err "Cannot install. Install Git or clone manually."
            exit 1
        fi
    fi

    if [ "$DOWNLOADED" = false ]; then
        cd "$SRC_DIR"
        LDFLAGS="-X main.Version=$VERSION"
        go build -ldflags "$LDFLAGS" -o "$TARGET" "./cmd/$APP_NAME/"
        ok "Built from source: $TARGET"
    fi
fi

# --- Add to PATH if needed ---
case ":$PATH:" in
    *:$INSTALL_DIR:*) ;;
    *)
        case "${SHELL##*/}" in
            zsh) echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$HOME/.zshrc" ;;
            bash) echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$HOME/.bashrc" ;;
            fish) echo "fish_add_path $INSTALL_DIR" >> "$HOME/.config/fish/config.fish" ;;
            *) warn "Add $INSTALL_DIR to your PATH manually" ;;
        esac
        export PATH="$INSTALL_DIR:$PATH"
        ;;
esac

# --- Post-install ---
"$TARGET" config init 2>/dev/null || true

echo ""
ok "✓ wc26 $TAG installed!"
echo "  Run '$APP_NAME --help' to get started"
echo "  Run '$APP_NAME auth login <email> <password>' to authenticate"
echo "  Run '$APP_NAME update' to upgrade later"

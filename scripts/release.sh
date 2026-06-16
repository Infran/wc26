#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 1 ]; then
  echo "Usage: $0 <tag> [release notes...]"
  echo "  e.g. $0 v0.3.0 \"Added update command, auth token\""
  exit 1
fi

TAG="$1"
shift
NOTES="${*:-"Release $TAG"}"
APP="wc26"
REPO="Infran/wc26"

echo "=== Building $APP $TAG ==="

mkdir -p dist

LDFLAGS="-X main.Version=${TAG#v}"

build() {
  local GOOS="$1" GOARCH="$2" suffix="$3"
  local name="${APP}_${GOOS}_${GOARCH}${suffix}"
  echo "  → $name"
  GOOS="$GOOS" GOARCH="$GOARCH" go build \
    -ldflags "$LDFLAGS" \
    -o "dist/$name" \
    "./cmd/$APP/"
}

build windows amd64 .exe
build darwin  amd64 ""
build darwin  arm64 ""
build linux   amd64 ""
build linux   arm64 ""

# --- Include install scripts ---
echo ""
echo "=== Including install scripts ==="
cp install.ps1 dist/
cp install.sh dist/
chmod +x dist/install.sh

# --- Generate checksums ---
echo ""
echo "=== Generating checksums ==="
cd dist
sha256sum * > wc26-checksums.txt
cat wc26-checksums.txt
cd ..

# --- Create GitHub release ---
echo ""
echo "=== Creating GitHub release ==="
gh release create "$TAG" dist/* \
  --repo "$REPO" \
  --title "$TAG" \
  --notes "$NOTES"

echo ""
echo "✓ Release $TAG published: https://github.com/$REPO/releases/tag/$TAG"

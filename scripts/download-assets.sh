#!/bin/bash

# Script to download CDN assets for local development
# This allows building with -tags localassets to serve assets locally

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
VENDOR_DIR="$ROOT_DIR/assets/vendor"

echo "Creating vendor directory at $VENDOR_DIR..."
mkdir -p "$VENDOR_DIR"

echo "Downloading CDN assets for local serving..."

# Download TailwindCSS
echo "Downloading TailwindCSS..."
curl -L -o "$VENDOR_DIR/tailwindcss.js" \
  "https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4/lib/cdn.min.js"

# Download Shoelace
echo "Downloading Shoelace..."
curl -L -o "$VENDOR_DIR/shoelace-autoloader.js" \
  "https://cdn.jsdelivr.net/npm/@shoelace-style/shoelace@2.19.1/cdn/shoelace-autoloader.js"
curl -L -o "$VENDOR_DIR/shoelace-light.css" \
  "https://cdn.jsdelivr.net/npm/@shoelace-style/shoelace@2.19.1/cdn/themes/light.css"

# Download HTMX
echo "Downloading HTMX..."
curl -L -o "$VENDOR_DIR/htmx.js" \
  "https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js"

# Download ApexCharts
echo "Downloading ApexCharts..."
curl -L -o "$VENDOR_DIR/apexcharts.min.js" \
  "https://cdn.jsdelivr.net/npm/apexcharts@4.3.0/dist/apexcharts.min.js"

# Download D3
echo "Downloading D3..."
curl -L -o "$VENDOR_DIR/d3.js" \
  "https://cdn.jsdelivr.net/npm/d3@7/dist/d3.min.js"

# Download Helia
echo "Downloading Helia..."
curl -L -o "$VENDOR_DIR/helia.min.js" \
  "https://cdn.jsdelivr.net/npm/helia/dist/index.min.js"

# Download Dexie
echo "Downloading Dexie..."
curl -L -o "$VENDOR_DIR/dexie.js" \
  "https://unpkg.com/dexie@latest/dist/dexie.js"

echo "✅ All assets downloaded successfully!"
echo ""
echo "To build with local assets, use:"
echo "  go build -tags localassets"
echo ""
echo "To use fully local vendor files (no CDN), set:"
echo "  export UI_USE_LOCAL_VENDOR=true"
echo ""
echo "To specify a custom assets path, set:"
echo "  export UI_ASSETS_PATH=/path/to/ui/module"
# Build Tags for Asset Generation

This module supports different asset serving strategies through Go build tags.

## Default Mode (Embedded Assets)

By default, the UI module embeds all assets using `go:embed` directives. This creates a self-contained binary with all JavaScript and CSS assets included.

```bash
# Normal build (assets are embedded)
go build ./...
```

## Local Assets Mode

When building with the `localassets` tag, the module will serve assets from the filesystem instead of embedding them. This is useful for:

1. **Development** - Changes to assets don't require rebuilding
2. **External modules** - Other modules can access and modify assets
3. **Custom deployments** - Assets can be served from a CDN or custom location

### Building with Local Assets

```bash
# Build with local asset serving
go build -tags localassets ./...

# Run with local assets
go run -tags localassets main.go
```

### Environment Variables

When using `localassets` build tag, you can configure asset behavior:

```bash
# Specify custom assets directory
export UI_ASSETS_PATH=/path/to/ui/assets

# Use local vendor files instead of CDN (requires downloading assets first)
export UI_USE_LOCAL_VENDOR=true
```

### Asset Path Resolution

The local assets mode searches for assets in the following order:

1. `UI_ASSETS_PATH` environment variable (if set)
2. Current working directory + `/assets`
3. Parent directories (up to 5 levels) looking for `/assets`
4. Falls back to current directory

## Downloading CDN Assets for Offline Use

To serve all dependencies locally (including CDN resources):

```bash
# Download all CDN assets to assets/vendor/
./scripts/download-assets.sh

# Build with local assets
go build -tags localassets ./...

# Run with fully local assets (no CDN)
UI_USE_LOCAL_VENDOR=true go run -tags localassets main.go
```

## Module Integration Example

Other Go modules can use this UI library with local assets:

```go
// In your module's main.go
package main

import (
    "github.com/go-sonr/ui"
    "net/http"
    "os"
)

func main() {
    // Set assets path if needed
    os.Setenv("UI_ASSETS_PATH", "/path/to/ui/module")
    
    // Create server with UI middleware
    mux := http.NewServeMux()
    ui.ServeAssets(mux)
    
    // Your routes here...
    
    http.ListenAndServe(":8080", mux)
}
```

Then build your module:

```bash
# Build your module with local UI assets
go build -tags localassets ./...
```

## Debugging Asset Loading

To check which assets path is being used:

```go
// When using localassets build tag
if path := ui.GetAssetsPath(); path != "" {
    fmt.Printf("Assets are being served from: %s\n", path)
}
```

## Build Tag Comparison

| Feature | Default (embedded) | `localassets` tag |
|---------|-------------------|-------------------|
| Binary size | Larger (includes assets) | Smaller |
| Asset changes | Requires rebuild | No rebuild needed |
| Distribution | Single binary | Binary + assets directory |
| Performance | Faster (in-memory) | Filesystem I/O |
| CDN fallback | Always uses CDN | Configurable |
| Development | Production-ready | Development-friendly |

## CI/CD Considerations

For production builds, use the default embedded mode:

```yaml
# GitHub Actions example
- name: Build production
  run: go build ./...

# For development/testing with local assets
- name: Build with local assets
  run: go build -tags localassets ./...
```
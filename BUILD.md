# Build Process

This document describes the build process for the UI module, which integrates TailwindCSS, Templ templates, and Go code generation.

## Quick Start

The build process is fully automated through Go:

```bash
# Run complete build (CSS + Templ + Go)
go run build.go

# Or use Make
make all

# Or use go generate
go generate ./...
```

## Build Script (`build.go`)

The `build.go` script orchestrates the entire build process:

### Basic Usage

```bash
# Full build process
go run build.go

# Install dependencies (templ, tailwindcss, templui)
go run build.go -install-deps

# Build with local assets (not embedded)
go run build.go -local-assets

# Development watch mode
go run build.go -watch
```

### Build Options

| Flag | Description |
|------|-------------|
| `-install-deps` | Install required build dependencies |
| `-skip-css` | Skip TailwindCSS generation |
| `-skip-templ` | Skip templ generation |
| `-skip-build` | Skip Go build |
| `-skip-tidy` | Skip go mod tidy |
| `-css-only` | Only generate CSS |
| `-templ-only` | Only generate templ files |
| `-local-assets` | Build with localassets tag |
| `-watch` | Watch mode for development |
| `-v` | Verbose output |

## Build Steps

The complete build process follows these steps:

### 1. TailwindCSS Generation
Converts `assets/css/input.css` to `assets/css/output.css` using TailwindCSS CLI:

```bash
go run build.go -css-only
# or
make css
```

### 2. Templ Generation
Generates Go code from `.templ` files:

```bash
go run build.go -templ-only
# or
make generate
```

### 3. Go Module Build
Builds the Go module with all generated code:

```bash
go run build.go -skip-css -skip-templ
# or
make build
```

## Using go generate

The module supports `go generate` for automated building:

```bash
# Generate everything (CSS + Templ)
go generate ./...

# This runs the directives in generate.go:
# - //go:generate go run build.go -css-only
# - //go:generate templ generate
```

## Makefile Targets

The Makefile provides convenient targets that use `build.go`:

```bash
make all           # Complete build
make generate      # Templ generation only
make css           # CSS generation only
make css-watch     # Watch CSS changes
make build         # Go build only
make build-local   # Build with localassets tag
make install-deps  # Install dependencies
make dev           # Development mode with watchers
```

## Dependencies Installation

### Automatic Installation

```bash
go run build.go -install-deps
# or
make install-deps
```

This installs:
- **templ**: Template engine for Go
- **tailwindcss**: Standalone CLI (no Node.js required)
- **templui**: Component library CLI

### Manual Installation

```bash
# Install templ
go install github.com/a-h/templ/cmd/templ@latest

# Install templui
go install github.com/axzilla/templui/cmd/templui@latest

# Download TailwindCSS CLI (example for Linux x64)
curl -L -o ./bin/tailwindcss \
  https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-x64
chmod +x ./bin/tailwindcss
```

## Development Workflow

### Watch Mode

For development with automatic rebuilding:

```bash
# Watch both CSS and Templ files
go run build.go -watch

# Or use Make for full dev environment
make dev
```

### Adding Components

```bash
# Add a new component
templui add button

# Regenerate after adding
go run build.go
```

### Local Assets Development

When developing with local assets (not embedded):

```bash
# Build with local assets
go run build.go -local-assets

# Set environment for runtime
export UI_ASSETS_PATH=/path/to/ui/module
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Build
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      
      - name: Install dependencies
        run: go run build.go -install-deps
      
      - name: Build
        run: go run build.go
      
      - name: Test
        run: go test ./...
```

### Docker Example

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

# Install dependencies and build
RUN go run build.go -install-deps && \
    go run build.go

# Final image
FROM alpine:latest
COPY --from=builder /app/bin /usr/local/bin
```

## Build Modes

### Production Build (Embedded Assets)

Default mode embeds all assets in the binary:

```bash
go run build.go
# or
go build ./...
```

### Development Build (Local Assets)

Serves assets from filesystem:

```bash
go run build.go -local-assets
# or
go build -tags localassets ./...
```

See [BUILD_TAGS.md](BUILD_TAGS.md) for details on build tags.

## Troubleshooting

### TailwindCSS Not Found

```bash
# Install TailwindCSS CLI
go run build.go -install-deps

# Or manually download to ./bin/
curl -L -o ./bin/tailwindcss [URL for your platform]
chmod +x ./bin/tailwindcss
```

### Templ Not Found

```bash
# Install templ
go install github.com/a-h/templ/cmd/templ@latest

# Verify installation
templ version
```

### Build Fails

```bash
# Clean and rebuild
make clean
go run build.go

# Verbose mode for debugging
go run build.go -v
```

## Performance Tips

1. **Skip unnecessary steps**: Use flags like `-skip-css` if CSS hasn't changed
2. **Use watch mode**: For development, use `-watch` to avoid manual rebuilds
3. **Parallel builds**: The Makefile uses parallel execution where possible
4. **Local assets**: Use `-local-assets` during development to avoid embedding

## Advanced Usage

### Custom TailwindCSS Path

```bash
# Use environment variable
export PATH="./bin:$PATH"
go run build.go
```

### Custom Build Pipeline

Create your own build script:

```go
//go:build ignore

package main

import "os/exec"

func main() {
    // Custom build steps
    exec.Command("tailwindcss", "-i", "input.css", "-o", "output.css").Run()
    exec.Command("templ", "generate").Run()
    exec.Command("go", "build", "./...").Run()
}
```

## Summary

The build process is designed to be:
- **Simple**: Single command builds everything
- **Flexible**: Granular control with flags
- **Fast**: Skip unnecessary steps
- **Portable**: No Node.js or npm required
- **CI-friendly**: Works in automated environments
# `go-sonr/ui`

[![GoDoc](https://godoc.org/github.com/go-sonr/ui?status.svg)](https://godoc.org/github.com/go-sonr/ui)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-sonr/ui)](https://goreportcard.com/report/github.com/go-sonr/ui)

UI Library built for the Sonr ecosystem. It leverages Templ, TailwindCSS 4, and Shoelace Web Components.

## Features

- **40+ UI Components**: Comprehensive set of pre-built components (buttons, cards, forms, charts, etc.)
- **TailwindCSS 4**: Modern CSS framework for rapid UI development
- **Shoelace Web Components**: Customizable and extensible web components
- **Templ**: Type-safe Go HTML template engine
- **HTTP Middleware**: Built-in asset serving for JavaScript and CSS
- **Zero Configuration**: CDN-based dependencies, no local build required
- **Integrated Build Process**: Automated TailwindCSS and Templ generation
- **Flexible Asset Serving**: Embedded or local filesystem modes via build tags
- **go:generate Support**: Standard Go toolchain integration

## Installation

```bash
go get github.com/go-sonr/ui
```

### Build Dependencies (Optional)

If you're contributing or need to regenerate assets:

```bash
# Install build tools
go run github.com/go-sonr/ui/build.go -install-deps

# Or if you have the module cloned
cd $GOPATH/src/github.com/go-sonr/ui
make install-deps
```

## Using as a Module Dependency

When using this module as a dependency in your project, all assets are pre-built and embedded. No additional build steps are required:

```go
import "github.com/go-sonr/ui"
```

The module includes:
- Pre-generated Go code from Templ templates
- Minified CSS from TailwindCSS
- Embedded JavaScript assets
- All components ready to use

## Quick Start

### Standard Library (net/http)

```go
package main

import (
    "net/http"
    ui "github.com/go-sonr/ui"
)

func main() {
    mux := http.NewServeMux()
    
    // Setup asset serving
    ui.ServeAssets(mux)
    
    // Add your routes
    mux.HandleFunc("/", homeHandler)
    
    http.ListenAndServe(":8080", mux)
}
```

### Using Middleware Wrapper

```go
package main

import (
    "net/http"
    ui "github.com/go-sonr/ui"
)

func main() {
    // Wrap your existing handler with UI middleware
    handler := ui.Middleware(yourHandler)
    
    http.ListenAndServe(":8080", handler)
}
```

## Available Components

The library provides 40+ components organized in categories:

- **Forms**: button, input, textarea, checkbox, radio, selectbox, inputotp
- **Layout**: card, accordion, tabs, sidebar, separator, sheet
- **Data**: table, chart, pagination
- **Feedback**: alert, toast, dialog, tooltip, popover
- **Navigation**: breadcrumb, dropdown
- **Display**: avatar, badge, progress, skeleton, rating

## Middleware Functions

### Core Functions

#### `AssetHandler() http.Handler`
Returns a handler for serving embedded JavaScript assets.
- Mount at `/assets/js/` to match component expectations
- Serves minified JavaScript files for components

#### `CSSHandler() http.Handler`
Returns a handler for serving embedded CSS assets.
- Mount at `/assets/css/` to match component expectations
- Serves CSS files including generated Tailwind CSS

#### `Middleware(next http.Handler) http.Handler`
Wraps an existing handler to automatically serve UI assets.
- Intercepts requests to `/assets/js/*` and `/assets/css/*`
- Passes through all other requests

#### `ServeAssets(mux *http.ServeMux)`
Convenience function to setup asset serving on a ServeMux.

#### `WrapHandler(h http.Handler) http.Handler`
Alternative name for `Middleware()` for clarity.

### Security & Performance

#### `SecurityHeaders(next http.Handler) http.Handler`
Adds security headers including:
- Content Security Policy (CSP) for CDN resources
- X-Content-Type-Options
- X-Frame-Options
- X-XSS-Protection
- Referrer-Policy

#### `PreloadMiddleware(next http.Handler) http.Handler`
Adds preload headers for CDN resources to improve initial page load.

#### `CacheControl(maxAge int) func(http.Handler) http.Handler`
Adds cache control headers for static assets.

## Framework Integration

### Chi Router

Build with tag: `go build -tags chi`

```go
import (
    "github.com/go-chi/chi/v5"
    ui "github.com/go-sonr/ui"
)

func main() {
    r := chi.NewRouter()
    
    // Option 1: Configure router
    ui.ChiRouter(r)
    
    // Option 2: Use middleware
    r.Use(ui.ChiMiddleware())
    
    r.Get("/", homeHandler)
    http.ListenAndServe(":8080", r)
}
```

### Gin

Build with tag: `go build -tags gin`

```go
import (
    "github.com/gin-gonic/gin"
    ui "github.com/go-sonr/ui"
)

func main() {
    r := gin.Default()
    
    // Option 1: Configure router
    ui.GinRouter(r)
    
    // Option 2: Use middleware
    r.Use(ui.GinMiddleware())
    
    r.GET("/", homeHandler)
    r.Run(":8080")
}
```

### Echo

```go
import (
    "github.com/labstack/echo/v4"
    ui "github.com/go-sonr/ui"
)

func main() {
    e := echo.New()
    
    // Use standard middleware wrapper
    e.Use(echo.WrapMiddleware(ui.Middleware))
    
    e.GET("/", homeHandler)
    e.Start(":8080")
}
```

### Fiber

```go
import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/adaptor"
    ui "github.com/go-sonr/ui"
)

func main() {
    app := fiber.New()
    
    // Adapt the standard HTTP handler
    app.Use(adaptor.HTTPMiddleware(ui.Middleware))
    
    app.Get("/", homeHandler)
    app.Listen(":8080")
}
```

## Complete Example

```go
package main

import (
    "context"
    "log"
    "net/http"
    
    "github.com/a-h/templ"
    ui "github.com/go-sonr/ui"
    "github.com/go-sonr/ui/layout"
    "github.com/go-sonr/ui/components/button"
    "github.com/go-sonr/ui/components/card"
)

func HomePage() templ.Component {
    return templ.ComponentFunc(func(ctx context.Context, w templ.RenderWriter) error {
        // Use UI layout components
        layout.HTML().Render(ctx, w)
        layout.Head().Render(ctx, w)
        // Head includes all component scripts automatically
        
        layout.Body().Render(ctx, w)
        
        // Your content here
        card.Card(card.Props{
            Class: "max-w-md mx-auto mt-10",
        }).Render(ctx, w)
        
        button.Button(button.Props{
            Variant: button.VariantDefault,
        }).Render(ctx, w)
        
        return nil
    })
}

func main() {
    mux := http.NewServeMux()
    
    // Setup UI assets
    ui.ServeAssets(mux)
    
    // Add routes
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        HomePage().Render(r.Context(), w)
    })
    
    // Apply security and performance middleware
    handler := ui.SecurityHeaders(
        ui.PreloadMiddleware(
            ui.CacheControl(3600)(mux),
        ),
    )
    
    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

## Asset Paths

The middleware serves the following assets:

- `/assets/js/*.min.js` - Minified JavaScript files for interactive components
- `/assets/css/*.css` - CSS files (including generated Tailwind CSS)

Components that include JavaScript:
- avatar, calendar, carousel, chart, code, collapsible, datepicker, dialog
- dropdown, input, inputotp, label, popover, progress, rating, selectbox
- sidebar, slider, tabs, tagsinput, textarea, timepicker, toast

## CDN Resources

The UI components also depend on CDN resources that are loaded directly:
- TailwindCSS 4 (Beta)
- Shoelace Web Components
- HTMX
- ApexCharts
- D3.js
- Helia (IPFS)
- Dexie (IndexedDB)

These are automatically included when using the `layout.Head()` component.

## Development

### Building

The build process integrates TailwindCSS CLI and Templ generation:

```bash
# Complete build (TailwindCSS + Templ + Go)
go run build.go
# or
make all
# or
go generate ./...

# Build with local assets (not embedded)
go build -tags localassets ./...
# or
make build-local

# Development with watch mode
go run build.go -watch
# or
make dev
```

### Build Options

```bash
# Install dependencies only
go run build.go -install-deps

# Generate CSS only
go run build.go -css-only

# Generate Templ files only
go run build.go -templ-only

# Skip specific steps
go run build.go -skip-css -skip-templ

# Verbose output
go run build.go -v
```

### Component Management

```bash
# Add a new component
templui add button

# Add multiple components
templui add button card dialog

# List available components
templui list

# After adding components, regenerate
go generate ./...
# or
make generate
```

### Testing Examples

```bash
# Build all examples
make examples

# Run server example
make examples-run

# Clean example binaries
make examples-clean
```

### CI/CD Integration

For automated builds in CI/CD pipelines:

```yaml
# GitHub Actions
- name: Setup and Build
  run: |
    go run build.go -install-deps
    go generate ./...
    go build ./...
    go test ./...
```

```bash
# Docker
RUN go run build.go -install-deps && \
    go run build.go && \
    go build -o app
```

## Custom Asset Serving

If you need to serve assets from a different path:

```go
// Serve from /static/ instead of /assets/js/
mux.Handle("/static/", http.StripPrefix("/static", ui.StaticHandler()))
```

### Local Assets Mode

For development or custom deployments, you can build with local assets instead of embedding:

```bash
# Build with local assets
go build -tags localassets ./...

# Set custom assets path
export UI_ASSETS_PATH=/path/to/ui/assets

# Use local vendor files (no CDN)
export UI_USE_LOCAL_VENDOR=true
```

See [BUILD_TAGS.md](BUILD_TAGS.md) for detailed information about build tags and [BUILD.md](BUILD.md) for complete build process documentation.

## Development Tips

1. **Always include `layout.Head()`** in your templates - it loads all necessary dependencies
2. **Use `ui.ServeAssets()`** or `ui.Middleware()`** to serve component JavaScript
3. **Apply security headers** in production environments
4. **Enable caching** for better performance with `CacheControl()`
5. **Use preload headers** to improve initial page load times
6. **Run `go generate ./...`** after modifying `.templ` files or CSS
7. **Use `-tags localassets`** for development to avoid rebuilding for asset changes

## Troubleshooting

### Assets not loading (404 errors)
- Ensure you've called `ui.ServeAssets(mux)` or are using `ui.Middleware()`
- Check that assets are mounted at `/assets/js/` and `/assets/css/`
- Verify the embedded files exist in `/assets/` directories

### Components not interactive
- Ensure the component's Script() is included in the Head
- Check browser console for JavaScript errors
- Verify CDN resources are loading (check network tab)

### CSP Errors
- If using custom CSP, ensure CDN domains are whitelisted
- Use the provided `SecurityHeaders()` middleware as a starting point

## License

MIT License - see LICENSE file for details
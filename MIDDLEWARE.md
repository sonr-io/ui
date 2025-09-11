# UI Middleware Documentation

The `go-sonr/ui` library provides HTTP middleware for serving component assets (JavaScript files) that are required by interactive UI components.

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

## Available Middleware Functions

### Core Functions

#### `AssetHandler() http.Handler`
Returns a handler for serving embedded JavaScript assets.
- Mount at `/assets/js/` to match component expectations
- Serves minified JavaScript files for components

#### `Middleware(next http.Handler) http.Handler`
Wraps an existing handler to automatically serve UI assets.
- Intercepts requests to `/assets/js/*`
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
    "github.com/go-sonr/ui/components/button"
    "github.com/go-sonr/ui/components/card"
)

func HomePage() templ.Component {
    return templ.ComponentFunc(func(ctx context.Context, w templ.RenderWriter) error {
        // Use UI layout components
        ui.HTML().Render(ctx, w)
        ui.Head().Render(ctx, w)
        // Head includes all component scripts automatically
        
        ui.Body().Render(ctx, w)
        
        // Your content here
        card.Card(&card.CardProps{
            Class: "max-w-md mx-auto mt-10",
        }).Render(ctx, w)
        
        button.Button(&button.ButtonProps{
            Variant: "default",
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

These are automatically included when using the `ui.Head()` component.

## Custom Asset Serving

If you need to serve assets from a different path:

```go
// Serve from /static/ instead of /assets/js/
mux.Handle("/static/", http.StripPrefix("/static", ui.StaticHandler()))
```

## Development Tips

1. **Always include `ui.Head()`** in your templates - it loads all necessary dependencies
2. **Use `ui.ServeAssets()`** or `ui.Middleware()` to serve component JavaScript
3. **Apply security headers** in production environments
4. **Enable caching** for better performance with `CacheControl()`
5. **Use preload headers** to improve initial page load times

## Troubleshooting

### Assets not loading (404 errors)
- Ensure you've called `ui.ServeAssets(mux)` or are using `ui.Middleware()`
- Check that assets are mounted at `/assets/js/`
- Verify the embedded files exist in `/assets/js/` directory

### Components not interactive
- Ensure the component's Script() is included in the Head
- Check browser console for JavaScript errors
- Verify CDN resources are loading (check network tab)

### CSP Errors
- If using custom CSP, ensure CDN domains are whitelisted
- Use the provided `SecurityHeaders()` middleware as a starting point
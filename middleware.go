//go:build !localassets
// +build !localassets

package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed assets/js/*.min.js assets/css/*.css
var assetsFS embed.FS

// AssetHandler returns an http.Handler that serves embedded JavaScript assets
// from the /assets/js/ path. This handler should be mounted at /assets/js/
// in your HTTP router.
//
// Example with standard library:
//
//	mux := http.NewServeMux()
//	mux.Handle("/assets/js/", ui.AssetHandler())
//
// Example with chi router:
//
//	r := chi.NewRouter()
//	r.Mount("/assets/js/", ui.AssetHandler())
func AssetHandler() http.Handler {
	sub, _ := fs.Sub(assetsFS, "assets/js")
	return http.StripPrefix("/assets/js/", http.FileServer(http.FS(sub)))
}

// CSSHandler returns an http.Handler that serves embedded CSS assets
// from the /assets/css/ path. This handler should be mounted at /assets/css/
// in your HTTP router.
//
// Example:
//
//	mux := http.NewServeMux()
//	mux.Handle("/assets/css/", ui.CSSHandler())
func CSSHandler() http.Handler {
	sub, _ := fs.Sub(assetsFS, "assets/css")
	return http.StripPrefix("/assets/css/", http.FileServer(http.FS(sub)))
}

// Middleware returns an HTTP middleware that serves UI assets and passes through
// other requests. This middleware automatically serves:
// - JavaScript assets at /assets/js/*
// - CSS assets at /assets/css/*
// - External CDN resources are not handled (loaded directly from CDN)
//
// Example:
//
//	mux := http.NewServeMux()
//	mux.Handle("/", ui.Middleware(yourHandler))
func Middleware(next http.Handler) http.Handler {
	jsHandler := AssetHandler()
	cssHandler := CSSHandler()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a JavaScript asset request
		if strings.HasPrefix(r.URL.Path, "/assets/js/") {
			jsHandler.ServeHTTP(w, r)
			return
		}

		// Check if this is a CSS asset request
		if strings.HasPrefix(r.URL.Path, "/assets/css/") {
			cssHandler.ServeHTTP(w, r)
			return
		}

		// Pass through to the next handler
		next.ServeHTTP(w, r)
	})
}

// StaticHandler returns an http.Handler that serves all static assets needed
// by the UI components. This includes JavaScript files and can be extended
// to include other assets like fonts or images if needed.
//
// This handler should be used when you want to handle all UI assets
// from a single mount point.
//
// Example:
//
//	mux := http.NewServeMux()
//	mux.Handle("/static/", http.StripPrefix("/static", ui.StaticHandler()))
func StaticHandler() http.Handler {
	return http.FileServer(http.FS(assetsFS))
}

// MustAssetString returns the content of an asset as a string.
// It panics if the asset cannot be read.
// This is useful for inlining critical assets in templates.
func MustAssetString(name string) string {
	data, err := assetsFS.ReadFile(path.Join("assets/js", name))
	if err != nil {
		panic(err)
	}
	return string(data)
}

// AssetExists checks if an asset exists in the embedded filesystem
func AssetExists(name string) bool {
	_, err := assetsFS.ReadFile(path.Join("assets/js", name))
	return err == nil
}

// ServeAssets is a convenience function that sets up asset serving on the
// provided ServeMux. It configures the standard asset paths used by UI components.
//
// Example:
//
//	mux := http.NewServeMux()
//	ui.ServeAssets(mux)
//	// Now add your other routes
//	mux.HandleFunc("/", homeHandler)
func ServeAssets(mux *http.ServeMux) {
	mux.Handle("/assets/js/", AssetHandler())
	mux.Handle("/assets/css/", CSSHandler())
}

// WrapHandler wraps an existing http.Handler with UI middleware.
// This is useful when you want to add UI asset serving to an existing handler.
//
// Example:
//
//	handler := ui.WrapHandler(yourExistingHandler)
//	http.ListenAndServe(":8080", handler)
func WrapHandler(h http.Handler) http.Handler {
	return Middleware(h)
}

// CDNFallbackMiddleware provides fallback handling for CDN resources.
// If a CDN is unavailable, this middleware can serve local copies.
// Note: This requires you to download and embed the CDN resources.
func CDNFallbackMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For now, we're using CDN directly
		// This can be extended to provide local fallbacks if needed
		next.ServeHTTP(w, r)
	})
}


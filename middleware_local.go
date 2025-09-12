//go:build localassets
// +build localassets

package ui

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// localAssetsFS implements fs.FS for local development
type localAssetsFS struct {
	root string
}

func (l localAssetsFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(l.root, name))
}

func (l localAssetsFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(l.root, name))
}

// assetsFS uses local filesystem when localassets build tag is set
var assetsFS = localAssetsFS{root: getAssetsPath()}

// getAssetsPath determines the assets directory path
// It checks for UI_ASSETS_PATH environment variable first,
// then tries common locations relative to the binary
func getAssetsPath() string {
	// Allow override via environment variable
	if envPath := os.Getenv("UI_ASSETS_PATH"); envPath != "" {
		return envPath
	}

	// Try to find assets relative to current working directory
	cwd, err := os.Getwd()
	if err == nil {
		// Check if we're in the module root
		if _, err := os.Stat(filepath.Join(cwd, "assets")); err == nil {
			return cwd
		}
		
		// Check if we're in a subdirectory of the module
		// Walk up the directory tree to find assets
		dir := cwd
		for i := 0; i < 5; i++ { // Limit depth to prevent infinite loop
			assetsPath := filepath.Join(dir, "assets")
			if _, err := os.Stat(assetsPath); err == nil {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	// Default to current directory
	return "."
}

// AssetHandler returns an http.Handler that serves JavaScript assets
// from the local filesystem when using localassets build tag
func AssetHandler() http.Handler {
	assetsPath := filepath.Join(getAssetsPath(), "assets", "js")
	return http.StripPrefix("/assets/js/", http.FileServer(http.Dir(assetsPath)))
}

// CSSHandler returns an http.Handler that serves CSS assets
// from the local filesystem when using localassets build tag
func CSSHandler() http.Handler {
	assetsPath := filepath.Join(getAssetsPath(), "assets", "css")
	return http.StripPrefix("/assets/css/", http.FileServer(http.Dir(assetsPath)))
}

// Middleware returns an HTTP middleware that serves UI assets from local filesystem
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

// StaticHandler returns an http.Handler that serves all static assets
// from the local filesystem when using localassets build tag
func StaticHandler() http.Handler {
	assetsPath := filepath.Join(getAssetsPath(), "assets")
	return http.FileServer(http.Dir(assetsPath))
}

// MustAssetString returns the content of an asset as a string from local filesystem
func MustAssetString(name string) string {
	data, err := os.ReadFile(filepath.Join(getAssetsPath(), "assets", "js", name))
	if err != nil {
		panic(err)
	}
	return string(data)
}

// AssetExists checks if an asset exists in the local filesystem
func AssetExists(name string) bool {
	_, err := os.Stat(filepath.Join(getAssetsPath(), "assets", "js", name))
	return err == nil
}

// ServeAssets is a convenience function that sets up asset serving on the
// provided ServeMux using local filesystem when localassets build tag is set
func ServeAssets(mux *http.ServeMux) {
	mux.Handle("/assets/js/", AssetHandler())
	mux.Handle("/assets/css/", CSSHandler())
}

// WrapHandler wraps an existing http.Handler with UI middleware
func WrapHandler(h http.Handler) http.Handler {
	return Middleware(h)
}

// CDNFallbackMiddleware provides fallback handling for CDN resources
func CDNFallbackMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// GetAssetsPath returns the current assets path being used
// This is useful for debugging and for modules that need to know where assets are located
func GetAssetsPath() string {
	return getAssetsPath()
}
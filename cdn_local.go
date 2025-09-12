//go:build localassets
// +build localassets

// Package ui provides a set of UI components for Go web applications
package ui

import (
	"net/http"
	"os"
	"strings"
)

// CDNConfig holds configuration for CDN resources or local asset paths
type CDNConfig struct {
	TailwindCSS   string
	Shoelace      string
	ShoelaceTheme string
	HTMX          string
	ApexCharts    string
	D3            string
	Helia         string
	Dexie         string
}

// DefaultCDNConfig returns configuration for local asset serving
// When using localassets build tag, these can point to local copies
// or still use CDN if preferred
func DefaultCDNConfig() CDNConfig {
	// By default, still use CDN even in local mode
	// Users can override by providing local copies
	return CDNConfig{
		TailwindCSS:   "https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4",
		Shoelace:      "https://cdn.jsdelivr.net/npm/@shoelace-style/shoelace@2.19.1/cdn/shoelace-autoloader.js",
		ShoelaceTheme: "https://cdn.jsdelivr.net/npm/@shoelace-style/shoelace@2.19.1/cdn/themes/light.css",
		HTMX:          "https://unpkg.com/htmx.org@2.0.4",
		ApexCharts:    "https://cdn.jsdelivr.net/npm/apexcharts@4.3.0/dist/apexcharts.min.js",
		D3:            "https://cdn.jsdelivr.net/npm/d3@7",
		Helia:         "https://cdn.jsdelivr.net/npm/helia/dist/index.min.js",
		Dexie:         "https://unpkg.com/dexie@latest/dist/dexie.js",
	}
}

// LocalCDNConfig returns configuration for serving all assets locally
// This requires downloading the CDN resources and placing them in assets/vendor/
func LocalCDNConfig() CDNConfig {
	return CDNConfig{
		TailwindCSS:   "/assets/vendor/tailwindcss.js",
		Shoelace:      "/assets/vendor/shoelace-autoloader.js",
		ShoelaceTheme: "/assets/vendor/shoelace-light.css",
		HTMX:          "/assets/vendor/htmx.js",
		ApexCharts:    "/assets/vendor/apexcharts.min.js",
		D3:            "/assets/vendor/d3.js",
		Helia:         "/assets/vendor/helia.min.js",
		Dexie:         "/assets/vendor/dexie.js",
	}
}

// CDNResources represents the CDN resources as HTTP headers
type CDNResources struct {
	config CDNConfig
}

// NewCDNResources creates a new CDN resources handler
// In local mode, it checks for UI_USE_LOCAL_VENDOR env var
func NewCDNResources() *CDNResources {
	// Check if we should use fully local resources
	if strings.ToLower(strings.TrimSpace(strings.ToLower(os.Getenv("UI_USE_LOCAL_VENDOR")))) == "true" {
		return &CDNResources{
			config: LocalCDNConfig(),
		}
	}
	return &CDNResources{
		config: DefaultCDNConfig(),
	}
}

// NewCDNResourcesWithConfig creates a new CDN resources handler with custom config
func NewCDNResourcesWithConfig(config CDNConfig) *CDNResources {
	return &CDNResources{
		config: config,
	}
}

// AddPreloadHeaders adds preload headers for CDN resources to improve performance
func (c *CDNResources) AddPreloadHeaders(w http.ResponseWriter) {
	// Preload critical resources
	w.Header().Add("Link", "<"+c.config.TailwindCSS+">; rel=preload; as=script")
	w.Header().Add("Link", "<"+c.config.Shoelace+">; rel=preload; as=script")
	w.Header().Add("Link", "<"+c.config.HTMX+">; rel=preload; as=script")
	w.Header().Add("Link", "<"+c.config.ShoelaceTheme+">; rel=preload; as=style")
}

// PreloadMiddleware returns middleware that adds preload headers for CDN resources
func PreloadMiddleware(next http.Handler) http.Handler {
	cdn := NewCDNResources()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only add preload headers for HTML responses
		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			cdn.AddPreloadHeaders(w)
		}
		next.ServeHTTP(w, r)
	})
}

// SecurityHeaders adds security headers for serving UI components
// In local mode, CSP is more permissive for development
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// More permissive CSP for local development
		csp := "default-src 'self' http: https:; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval' http: https:; " +
			"style-src 'self' 'unsafe-inline' http: https:; " +
			"font-src 'self' data: http: https:; " +
			"img-src 'self' data: http: https:; " +
			"connect-src 'self' http: https: ws: wss:;"

		w.Header().Set("Content-Security-Policy", csp)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}

// CacheControl adds cache control headers for static assets
// In local mode, disable caching for easier development
func CacheControl(maxAge int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Disable cache in local development mode
			if strings.HasPrefix(r.URL.Path, "/assets/") {
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
			}
			next.ServeHTTP(w, r)
		})
	}
}
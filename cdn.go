// Package ui provides a set of UI components for Go web applications
package ui

import (
	"fmt"
	"net/http"
	"strings"
)

// CDNConfig holds configuration for CDN resources used by UI components
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

// DefaultCDNConfig returns the default CDN configuration
func DefaultCDNConfig() CDNConfig {
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

// CDNResources represents the CDN resources as HTTP headers
type CDNResources struct {
	config CDNConfig
}

// NewCDNResources creates a new CDN resources handler with default config
func NewCDNResources() *CDNResources {
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
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Content Security Policy allowing CDN resources
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval' " +
			"https://cdn.jsdelivr.net " +
			"https://unpkg.com " +
			"https://cdn.tailwindcss.com; " +
			"style-src 'self' 'unsafe-inline' " +
			"https://cdn.jsdelivr.net " +
			"https://cdn.tailwindcss.com; " +
			"font-src 'self' data: " +
			"https://cdn.jsdelivr.net; " +
			"img-src 'self' data: https:; " +
			"connect-src 'self' https: wss:;"

		w.Header().Set("Content-Security-Policy", csp)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}

// CacheControl adds cache control headers for static assets
func CacheControl(maxAge int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set cache headers for static assets
			if strings.HasPrefix(r.URL.Path, "/assets/") {
				w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
			}
			next.ServeHTTP(w, r)
		})
	}
}


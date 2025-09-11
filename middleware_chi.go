//go:build chi
// +build chi

package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ChiRouter configures a chi router with UI asset serving
func ChiRouter(r chi.Router) {
	// Add compression
	r.Use(middleware.Compress(5))

	// Mount asset handler
	r.Mount("/assets/js", AssetHandler())
}

// ChiMiddleware returns a chi-compatible middleware for serving UI assets
func ChiMiddleware() func(next http.Handler) http.Handler {
	return Middleware
}


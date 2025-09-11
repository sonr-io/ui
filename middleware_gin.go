//go:build gin
// +build gin

package ui

import (
	"net/http"
	"strings"
	
	"github.com/gin-gonic/gin"
)

// GinMiddleware returns a gin-compatible middleware for serving UI assets
func GinMiddleware() gin.HandlerFunc {
	assetHandler := AssetHandler()
	
	return func(c *gin.Context) {
		// Check if this is an asset request
		if strings.HasPrefix(c.Request.URL.Path, "/assets/js/") {
			assetHandler.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// GinRouter configures a gin engine with UI asset serving
func GinRouter(r *gin.Engine) {
	// Static file serving
	r.GET("/assets/js/*filepath", func(c *gin.Context) {
		AssetHandler().ServeHTTP(c.Writer, c.Request)
	})
}
package main

import (
	"log"
	"net/http"

	ui "github.com/go-sonr/ui"
)

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Setup UI asset serving
	ui.ServeAssets(mux)

	// Add the home page route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set content type
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Render the page using the templ component
		HomePage().Render(r.Context(), w)
	})

	// Apply security headers and caching
	handler := ui.SecurityHeaders(
		ui.CacheControl(3600)(
			ui.PreloadMiddleware(mux),
		),
	)

	// Start the server
	log.Println("Server starting on http://localhost:8080")
	log.Println("Visit http://localhost:8080 to see the UI components in action")
	log.Fatal(http.ListenAndServe(":8080", handler))
}


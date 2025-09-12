package main

import (
	"log"
	"net/http"

	ui "github.com/go-sonr/ui"
)

func main() {
	// Option 1: Using the middleware wrapper
	handler := ui.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Simple Example</title>
				<script src="/assets/js/button.min.js"></script>
			</head>
			<body>
				<h1>UI Components Example</h1>
				<button class="ui-button">Click Me</button>
			</body>
			</html>
		`))
	}))

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

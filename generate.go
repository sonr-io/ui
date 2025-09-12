// Package ui provides a set of UI components for Go web applications
//
// Code Generation
//
// This package uses go:generate directives to automate the build process.
// To regenerate all assets and templates, run:
//
//	go generate ./...
//
// Or use the Makefile:
//
//	make generate
//
// The generation process includes:
//   - TailwindCSS compilation (CSS generation)
//   - Templ template compilation (Go code generation)
//   - Asset embedding (when not using localassets build tag)
package ui

// Generate CSS from TailwindCSS
//go:generate go run build.go -css-only

// Generate Go code from templ files  
//go:generate templ generate

// Alternative: Run full build process
// Uncomment the line below to run the complete build with go generate
// //go:generate go run build.go
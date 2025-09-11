.PHONY: all generate format lint tidy clean build css css-watch examples examples-run examples-clean templ server tailwind dev help

# Default target
all: format generate css tidy lint build examples

# Generate Go code from templ files
generate:
	@echo "Generating templ files..."
	@templ generate

# Format both Go and templ files
format:
	@echo "Formatting Go files..."
	@go fmt ./...
	@echo "Formatting templ files..."
	@templ fmt .

# Lint Go code
lint:
	@echo "Running Go vet..."
	@go vet ./...
	@echo "Checking module dependencies..."
	@go mod verify

# Tidy Go module dependencies
tidy:
	@echo "Tidying Go module dependencies..."
	@go mod tidy

# Build the module
build:
	@echo "Building module..."
	@go build ./...

# Generate CSS from Tailwind input file
css:
	@echo "Generating CSS with Tailwind..."
	@npx tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --minify 2>/dev/null || \
		(echo "Installing Tailwind CSS..." && \
		npm install -D tailwindcss && \
		npx tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --minify)

# Watch CSS files for changes and regenerate
css-watch:
	@echo "Watching CSS files for changes..."
	@npx tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --watch

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@find . -name "*_templ.go" -type f -delete
	@find . -name "*_templ.txt" -type f -delete
	@rm -f ./assets/css/output.css

# Install or update a component via templui
add:
	@if [ -z "$(COMPONENT)" ]; then \
		echo "Usage: make add COMPONENT=button"; \
		echo "       make add COMPONENT=\"button card dialog\""; \
		exit 1; \
	fi
	@echo "Adding component(s): $(COMPONENT)"
	@templui add $(COMPONENT)

# List available components
list:
	@templui list

# Update templ version
update-templ:
	@echo "Updating templ to latest version..."
	@go get -u github.com/a-h/templ
	@go mod tidy

# Build examples
examples:
	@echo "Building examples..."
	@$(MAKE) -C examples build

# Run server example
examples-run:
	@$(MAKE) -C examples run-server

# Clean examples
examples-clean:
	@$(MAKE) -C examples clean

# Development targets with hot reload
# Run templ generation in watch mode
templ:
	@templ generate --watch --proxy="http://localhost:8090" --open-browser=false

# Run air for Go hot reload
server:
	@air \
	--build.cmd "go build -o tmp/bin/main ./examples/server/main.go" \
	--build.bin "tmp/bin/main" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# Watch Tailwind CSS changes (alias for css-watch)
tailwind:
	@$(MAKE) css-watch

# Start development server with all watchers
dev:
	@echo "Starting development environment..."
	@echo "templ UI: http://localhost:7331"
	@echo "Go server: http://localhost:8090"
	@$(MAKE) -j3 tailwind templ server

# Help target
help:
	@echo "Available targets:"
	@echo "  make all          - Run format, generate, css, tidy, lint, build, and examples"
	@echo "  make generate     - Generate Go code from templ files"
	@echo "  make format       - Format Go and templ files"
	@echo "  make lint         - Run Go linting and verify modules"
	@echo "  make tidy         - Tidy Go module dependencies"
	@echo "  make build        - Build the Go module"
	@echo "  make css          - Generate minified CSS from Tailwind input"
	@echo "  make css-watch    - Watch and regenerate CSS on changes"
	@echo "  make dev          - Start full development environment (templ, CSS, Go server)"
	@echo "  make templ        - Watch and generate templ files"
	@echo "  make server       - Run Go server with hot reload (Air)"
	@echo "  make tailwind     - Watch Tailwind CSS changes (alias for css-watch)"
	@echo "  make examples     - Build all example applications"
	@echo "  make examples-run - Run the server example"
	@echo "  make examples-clean - Clean example binaries"
	@echo "  make clean        - Remove generated files"
	@echo "  make add COMPONENT=name - Add a component via templui"
	@echo "  make list         - List available templui components"
	@echo "  make update-templ - Update templ to latest version"
	@echo "  make help         - Show this help message"
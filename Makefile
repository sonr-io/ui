.PHONY: all generate format lint clean build css css-watch help

# Default target
all: format generate css lint build

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
	@go mod tidy
	@go mod verify

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

# Help target
help:
	@echo "Available targets:"
	@echo "  make all          - Run format, generate, css, lint, and build"
	@echo "  make generate     - Generate Go code from templ files"
	@echo "  make format       - Format Go and templ files"
	@echo "  make lint         - Run Go linting and verify modules"
	@echo "  make build        - Build the Go module"
	@echo "  make css          - Generate minified CSS from Tailwind input"
	@echo "  make css-watch    - Watch and regenerate CSS on changes"
	@echo "  make clean        - Remove generated files"
	@echo "  make add COMPONENT=name - Add a component via templui"
	@echo "  make list         - List available templui components"
	@echo "  make update-templ - Update templ to latest version"
	@echo "  make help         - Show this help message"
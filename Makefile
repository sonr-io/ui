.PHONY: all generate format lint clean build help

# Default target
all: format generate lint build

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

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@find . -name "*_templ.go" -type f -delete
	@find . -name "*_templ.txt" -type f -delete

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
	@echo "  make all          - Run format, generate, lint, and build"
	@echo "  make generate     - Generate Go code from templ files"
	@echo "  make format       - Format Go and templ files"
	@echo "  make lint         - Run Go linting and verify modules"
	@echo "  make build        - Build the Go module"
	@echo "  make clean        - Remove generated files"
	@echo "  make add COMPONENT=name - Add a component via templui"
	@echo "  make list         - List available templui components"
	@echo "  make update-templ - Update templ to latest version"
	@echo "  make help         - Show this help message"
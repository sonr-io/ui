.PHONY: all build run clean server simple help

# Default target
all: build

# Build all examples
build: server simple
	@echo "All examples built successfully"

# Build server example
server:
	@echo "Building server example..."
	@cd server && templ generate && go build -o server .

# Build simple example  
simple:
	@echo "Building simple example..."
	@cd simple && go build -o simple .

# Run server example
run-server: server
	@echo "Running server example on http://localhost:8080"
	@cd server && ./server

# Run simple example
run-simple: simple
	@echo "Running simple example on http://localhost:8081"
	@cd simple && ./simple

# Clean all example binaries
clean:
	@echo "Cleaning example binaries..."
	@rm -f server/server simple/simple
	@rm -f server/*_templ.go simple/*_templ.go
	@echo "Examples cleaned"

# Help
help:
	@echo "Example targets:"
	@echo "  make build        - Build all examples"
	@echo "  make server       - Build server example"
	@echo "  make simple       - Build simple example"
	@echo "  make run-server   - Build and run server example (port 8080)"
	@echo "  make run-simple   - Build and run simple example (port 8081)"
	@echo "  make clean        - Remove all example binaries"
	@echo "  make help         - Show this help message"
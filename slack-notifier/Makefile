.PHONY: help build run test clean deps examples

# Default target
help:
	@echo "Available commands:"
	@echo "  make deps      - Install dependencies"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the main application"
	@echo "  make test      - Run all tests"
	@echo "  make examples  - Run example applications"
	@echo "  make clean     - Clean build artifacts"

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build the application
build:
	go build -o bin/slack-notifier cmd/server/main.go

# Run the main application
run:
	go run cmd/server/main.go

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run example applications
examples:
	@echo "Running basic usage example..."
	go run examples/basic_usage.go
	@echo ""
	@echo "Running event types example..."
	go run examples/event_types.go

# Run basic usage example
example-basic:
	go run examples/basic_usage.go

# Run event types example
example-events:
	go run examples/event_types.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Setup development environment
setup: install-tools deps
	@echo "Development environment setup complete!"
	@echo "Don't forget to set your SLACK_BOT_TOKEN environment variable"

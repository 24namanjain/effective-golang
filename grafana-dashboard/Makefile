.PHONY: help build run test clean deps setup docker-up docker-down fmt lint install-tools generate-load info check-prereqs

# Default target
help:
	@echo "🚀 System Monitor - Available Commands:"
	@echo ""
	@echo "📦 Setup & Dependencies:"
	@echo "  make setup         - Complete development environment setup"
	@echo "  make deps          - Install Go dependencies"
	@echo "  make install-tools - Install development tools"
	@echo ""
	@echo "🔨 Build & Run:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application (local data source)"
	@echo "  make run-grafana   - Run with Grafana data source"
	@echo "  make run-prometheus - Run with Prometheus data source"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  make docker-up         - Start full Grafana stack (with custom exporter)"
	@echo "  make docker-up-simple  - Start simple Grafana stack (without custom exporter)"
	@echo "  make docker-down       - Stop full Grafana stack"
	@echo "  make docker-down-simple - Stop simple Grafana stack"
	@echo "  make docker-build      - Build Docker image for metrics exporter"
	@echo ""
	@echo "🧪 Testing & Quality:"
	@echo "  make test          - Run all tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make fmt           - Format Go code"
	@echo "  make lint          - Lint Go code"
	@echo ""
	@echo "🛠️  Utilities:"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make generate-load - Generate system load for testing"
	@echo "  make info          - Show system information"
	@echo "  make check-prereqs - Check prerequisites"
	@echo ""

# Check prerequisites
check-prereqs:
	@echo "🔍 Checking prerequisites..."
	@command -v go >/dev/null 2>&1 || { echo "❌ Go is not installed. Please install Go 1.21+"; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "⚠️  Docker is not installed. Docker Compose features will not work."; }
	@if command -v docker-compose >/dev/null 2>&1; then \
		echo "✅ Docker Compose found"; \
	elif command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		echo "✅ Docker Compose (v2) found"; \
	else \
		echo "⚠️  Docker Compose is not installed. Docker Compose features will not work."; \
	fi
	@echo "✅ Prerequisites check completed"

# Install development tools
install-tools:
	@echo "🔧 Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/axw/gocov/gocov@latest
	go install github.com/AlekSi/gocov-xml@latest
	@echo "✅ Development tools installed"

# Setup development environment
setup: check-prereqs install-tools deps
	@echo "🎉 Development environment setup complete!"
	@echo ""
	@echo "📋 Next steps:"
	@echo "1. Set your environment variables (see README.md)"
	@echo "2. Run 'make run' to start the application"
	@echo "3. Open http://localhost:8080 in your browser"

# Install dependencies
deps:
	@echo "📦 Installing Go dependencies..."
	go mod tidy
	go mod download
	@echo "✅ Dependencies installed"

# Build the application
build:
	@echo "🔨 Building System Monitor..."
	go build -o bin/system-monitor cmd/monitor/main.go
	@echo "✅ Build complete: bin/system-monitor"

# Run with local data source
run: build
	@echo "🚀 Starting System Monitor (Local Data Source)..."
	@echo "📊 Dashboard will be available at: http://localhost:8080"
	@echo "Press Ctrl+C to stop"
	./bin/system-monitor

# Run with Grafana data source
run-grafana: build
	@echo "🚀 Starting System Monitor (Grafana Data Source)..."
	@echo "📊 Dashboard will be available at: http://localhost:8080"
	@echo "Make sure Grafana is running (make docker-up)"
	@echo "Press Ctrl+C to stop"
	DATA_SOURCE_TYPE=grafana ./bin/system-monitor

# Run with Prometheus data source
run-prometheus: build
	@echo "🚀 Starting System Monitor (Prometheus Data Source)..."
	@echo "📊 Dashboard will be available at: http://localhost:8080"
	@echo "Make sure Prometheus is running (make docker-up)"
	@echo "Press Ctrl+C to stop"
	DATA_SOURCE_TYPE=prometheus ./bin/system-monitor

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "🎨 Formatting Go code..."
	go fmt ./...

# Lint code
lint:
	@echo "🔍 Linting Go code..."
	golangci-lint run

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

# Docker commands
docker-up:
	@echo "🐳 Starting Grafana stack..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose up -d; \
	elif command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		docker compose up -d; \
	else \
		echo "❌ Neither docker-compose nor docker compose found. Please install Docker Compose."; \
		exit 1; \
	fi
	@echo "✅ Grafana stack started"
	@echo "📊 Grafana: http://localhost:3000 (admin/admin123)"
	@echo "📈 Prometheus: http://localhost:9090"
	@echo "📊 Node Exporter: http://localhost:9100"

docker-up-simple:
	@echo "🐳 Starting simple Grafana stack (without custom exporter)..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose -f docker-compose.simple.yml up -d; \
	elif command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		docker compose -f docker-compose.simple.yml up -d; \
	else \
		echo "❌ Neither docker-compose nor docker compose found. Please install Docker Compose."; \
		exit 1; \
	fi
	@echo "✅ Simple Grafana stack started"
	@echo "📊 Grafana: http://localhost:3000 (admin/admin123)"
	@echo "📈 Prometheus: http://localhost:9090"
	@echo "📊 Node Exporter: http://localhost:9100"

docker-down:
	@echo "🐳 Stopping Grafana stack..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose down; \
	elif command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		docker compose down; \
	else \
		echo "❌ Neither docker-compose nor docker compose found."; \
		exit 1; \
	fi
	@echo "✅ Grafana stack stopped"

docker-down-simple:
	@echo "🐳 Stopping simple Grafana stack..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose -f docker-compose.simple.yml down; \
	elif command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		docker compose -f docker-compose.simple.yml down; \
	else \
		echo "❌ Neither docker-compose nor docker compose found."; \
		exit 1; \
	fi
	@echo "✅ Simple Grafana stack stopped"

docker-build:
	@echo "🐳 Building Docker image..."
	docker build -f Dockerfile.metrics -t system-monitor:latest .
	@echo "✅ Docker image built: system-monitor:latest"

# Generate system load for testing
generate-load:
	@echo "🔥 Generating system load for testing..."
	@echo "This will create CPU and memory load for 30 seconds..."
	@echo "Press Ctrl+C to stop early"
	@timeout 30s bash -c 'while true; do dd if=/dev/zero of=/dev/null bs=1M count=1000 2>/dev/null; done' || true
	@echo "✅ Load generation complete"

# Show system information
info:
	@echo "📊 System Information:"
	@echo "Go version: $(shell go version)"
	@echo "Go modules: $(shell go env GOMOD)"
	@echo "Architecture: $(shell uname -m)"
	@echo "OS: $(shell uname -s)"
	@echo "Docker: $(shell docker --version 2>/dev/null || echo 'Not installed')"
	@if command -v docker-compose >/dev/null 2>&1; then \
		echo "Docker Compose: $(shell docker-compose --version)"; \
	elif command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		echo "Docker Compose: $(shell docker compose version)"; \
	else \
		echo "Docker Compose: Not installed"; \
	fi

# Development mode with hot reload (requires air)
dev:
	@echo "🔥 Starting development mode with hot reload..."
	@command -v air >/dev/null 2>&1 || { echo "❌ Air is not installed. Install with: go install github.com/cosmtrek/air@latest"; exit 1; }
	air

# Create .air.toml for hot reload
setup-air:
	@echo "🔥 Setting up Air for hot reload..."
	@command -v air >/dev/null 2>&1 || { echo "❌ Air is not installed. Install with: go install github.com/cosmtrek/air@latest"; exit 1; }
	@if [ ! -f .air.toml ]; then \
		echo "Creating .air.toml configuration..."; \
		air init; \
	fi
	@echo "✅ Air configuration ready"

# Security scan
security-scan:
	@echo "🔒 Running security scan..."
	go list -json -deps ./... | nancy sleuth
	@echo "✅ Security scan complete"

# Performance benchmark
benchmark:
	@echo "⚡ Running performance benchmarks..."
	go test -bench=. -benchmem ./...
	@echo "✅ Benchmarks complete"

# Generate documentation
docs:
	@echo "📚 Generating documentation..."
	@command -v godoc >/dev/null 2>&1 || { echo "❌ godoc is not installed. Install with: go install golang.org/x/tools/cmd/godoc@latest"; exit 1; }
	godoc -http=:6060 &
	@echo "📖 Documentation available at: http://localhost:6060"
	@echo "Press Ctrl+C to stop"

# Release build
release: clean
	@echo "🚀 Building release binaries..."
	GOOS=linux GOARCH=amd64 go build -o bin/system-monitor-linux-amd64 cmd/monitor/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/system-monitor-darwin-amd64 cmd/monitor/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/system-monitor-windows-amd64.exe cmd/monitor/main.go
	@echo "✅ Release binaries created in bin/"

# Show help for environment variables
env-help:
	@echo "🔧 Environment Variables:"
	@echo ""
	@echo "Data Source Configuration:"
	@echo "  DATA_SOURCE_TYPE    - Data source type (local, prometheus, grafana)"
	@echo "  DATA_SOURCE_URL     - Data source URL"
	@echo "  GRAFANA_URL         - Grafana URL (for grafana data source)"
	@echo "  GRAFANA_USERNAME    - Grafana username"
	@echo "  GRAFANA_PASSWORD    - Grafana password"
	@echo "  GRAFANA_API_KEY     - Grafana API key (alternative to username/password)"
	@echo "  PROMETHEUS_URL      - Prometheus URL (for prometheus data source)"
	@echo ""
	@echo "Alert Configuration:"
	@echo "  ALERT_BACKEND_TYPE  - Alert backend type (slack, webhook, email)"
	@echo "  SLACK_BOT_TOKEN     - Slack bot token"
	@echo "  SLACK_CHANNEL       - Slack channel for alerts"
	@echo "  WEBHOOK_URL         - Webhook URL (for webhook alerts)"
	@echo ""
	@echo "Threshold Configuration:"
	@echo "  CPU_THRESHOLD       - CPU usage threshold (percentage)"
	@echo "  MEMORY_THRESHOLD    - Memory usage threshold (percentage)"
	@echo "  LATENCY_THRESHOLD   - Latency threshold (milliseconds)"
	@echo ""
	@echo "Application Configuration:"
	@echo "  DASHBOARD_PORT      - Dashboard port (default: 8080)"
	@echo "  METRICS_INTERVAL    - Metrics collection interval (default: 5s)"
	@echo "  ALERT_COOLDOWN      - Alert cooldown period (default: 5m)"
	@echo "  ENVIRONMENT         - Environment (development, production)"

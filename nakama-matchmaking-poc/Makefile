# Nakama Matchmaking POC Makefile

.PHONY: help build start stop test test-multi clean logs

help: ## Show this help message
	@echo "Nakama Matchmaking POC - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the Go module
	@echo "🔨 Building Go module..."
	docker run --rm -v $(PWD)/modules:/workspace -w /workspace golang:1.21 bash -c "go mod tidy && go build -buildmode=plugin -o main.so main.go"
	@echo "✅ Go module built successfully"

install: ## Install Node.js dependencies
	@echo "📦 Installing Node.js dependencies..."
	npm install
	@echo "✅ Dependencies installed"

start: ## Start Nakama and database
	@echo "🚀 Starting Nakama and database..."
	docker compose up -d
	@echo "✅ Services started"
	@echo "📊 Nakama: http://localhost:7350"
	@echo "🗄️  Database: localhost:26257"

stop: ## Stop Nakama and database
	@echo "🛑 Stopping services..."
	docker compose down
	@echo "✅ Services stopped"

restart: stop start ## Restart all services

test: ## Run single client test
	@echo "🧪 Running single client test..."
	npm test

test-multi: ## Run multi-client test
	@echo "🧪 Running multi-client test..."
	npm run test-multi

logs: ## View Nakama logs
	@echo "📋 Viewing Nakama logs..."
	docker compose logs -f nakama

logs-db: ## View database logs
	@echo "📋 Viewing database logs..."
	docker compose logs -f cockroachdb

status: ## Check service status
	@echo "📊 Service status:"
	docker compose ps

clean: ## Clean up containers and volumes
	@echo "🧹 Cleaning up..."
	docker compose down -v
	docker system prune -f
	@echo "✅ Cleanup complete"

setup: install build start ## Complete setup (install, build, start)
	@echo "🎉 Setup complete! Run 'make test-multi' to test matchmaking"

dev: ## Development mode - start services and show logs
	@echo "🔧 Starting development mode..."
	docker compose up

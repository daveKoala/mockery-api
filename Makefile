.PHONY: build run start stop clean test help

# Default config file
CONFIG ?= config.json

# Binary name
BINARY = mockery-api

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY)..."
	@go build -o $(BINARY)
	@echo "Build complete: ./$(BINARY)"

run: build ## Build and run the server
	@echo "Starting $(BINARY) with $(CONFIG)..."
	@./$(BINARY) -config $(CONFIG)

start: build ## Start the server in the background
	@echo "Starting $(BINARY) in background..."
	@./$(BINARY) -config $(CONFIG) > server.log 2>&1 & echo $$! > server.pid
	@echo "Server started (PID: $$(cat server.pid))"
	@echo "Logs: tail -f server.log"

stop: ## Stop the background server
	@if [ -f server.pid ]; then \
		echo "Stopping server (PID: $$(cat server.pid))..."; \
		kill $$(cat server.pid) 2>/dev/null || echo "Server not running"; \
		rm -f server.pid; \
	else \
		echo "No PID file found. Server may not be running."; \
	fi

status: ## Check if server is running
	@if [ -f server.pid ]; then \
		if ps -p $$(cat server.pid) > /dev/null 2>&1; then \
			echo "Server is running (PID: $$(cat server.pid))"; \
		else \
			echo "Server is not running (stale PID file)"; \
			rm -f server.pid; \
		fi \
	else \
		echo "Server is not running"; \
	fi

logs: ## Tail the server logs
	@if [ -f server.log ]; then \
		tail -f server.log; \
	else \
		echo "No log file found. Server may not have been started with 'make start'"; \
	fi

test: build ## Run basic tests against the server
	@echo "Testing health check..."
	@curl -s http://localhost:3000/_health && echo "" || echo "Health check failed"
	@echo ""
	@echo "Testing public endpoint..."
	@curl -s http://localhost:3000/api/products && echo "" || echo "Public endpoint failed"
	@echo ""
	@echo "Testing auth endpoint (should fail)..."
	@curl -s http://localhost:3000/api/users && echo "" || echo "Expected auth failure"
	@echo ""
	@echo "Testing auth endpoint with header..."
	@curl -s -H "Authorization: Bearer token" http://localhost:3000/api/users && echo "" || echo "Auth endpoint failed"

clean: stop ## Clean build artifacts and logs
	@echo "Cleaning up..."
	@rm -f $(BINARY)
	@rm -f server.log
	@rm -f server.pid
	@echo "Clean complete"

restart: stop start ## Restart the server

dev: ## Run in development mode (auto-reload on config changes - requires fswatch)
	@command -v fswatch >/dev/null 2>&1 || { echo "fswatch not installed. Install with: brew install fswatch"; exit 1; }
	@echo "Watching $(CONFIG) for changes..."
	@make run & echo $$! > server.pid; \
	fswatch -o $(CONFIG) | while read change; do \
		echo "Config changed, restarting..."; \
		make restart; \
	done

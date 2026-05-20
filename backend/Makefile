.PHONY: help dev build test clean docker-up docker-down migrate

help: ## Show this help message
	@echo "RecallFlow Backend - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

dev: ## Run development server
	@echo "Starting RecallFlow API in development mode..."
	go run cmd/api/main.go

build: ## Build the binary
	@echo "Building RecallFlow API..."
	go build -o bin/recallflow-api cmd/api/main.go
	@echo "✅ Build complete: bin/recallflow-api"

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf tmp/

install-deps: ## Install Go dependencies
	go mod download
	go mod tidy

docker-up: ## Start PostgreSQL and Redis with Docker Compose
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

migrate: ## Run database migrations
	@echo "Running database migrations..."
	psql $(DATABASE_URL) < db/schema.sql
	@echo "✅ Migrations complete"

seed: ## Seed database with test data
	@echo "Seeding database..."
	# Add seed script here

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run

.DEFAULT_GOAL := help

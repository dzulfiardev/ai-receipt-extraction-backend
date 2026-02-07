. PHONY: help run build migrate-up migrate-down migrate-drop migrate-version migrate-force migrate-create deps clean

help:  ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+: .*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the API server
	@echo "🚀 Starting API server..."
	@go run cmd/api/main. go

build: ## Build the API server
	@echo "🔨 Building API server..."
	@go build -o bin/api cmd/api/main.go

migrate-up: ## Run all pending migrations
	@echo "🚀 Running migrations..."
	@go run cmd/migrate/main. go -command=up

migrate-down: ## Rollback last migration
	@echo "⏪ Rolling back migration..."
	@go run cmd/migrate/main.go -command=down

migrate-drop: ## Drop all tables (DANGER!)
	@echo "⚠️ Dropping all tables..."
	@go run cmd/migrate/main. go -command=drop

migrate-version: ## Show current migration version
	@go run cmd/migrate/main.go -command=version

migrate-force: ## Force migration to specific version (usage: make migrate-force VERSION=1)
	@go run cmd/migrate/main.go -command=force -version=$(VERSION)
 
migrate-create: ## Create new migration (usage: make migrate-create NAME=add_column)
	@go run cmd/migrate/main.go -command=create -name=$(NAME)

# migrate create manual command optional 
# migrate create -ext sql -dir migrations -seq create_receipts_table

deps: ## Install dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf bin/

docker-up: ## Start Docker containers (Postgres + MinIO)
	@echo "🐳 Starting Docker containers..."
	@docker-compose up -d

docker-down: ## Stop Docker containers
	@echo "🐳 Stopping Docker containers..."
	@docker-compose down

docker-logs: ## Show Docker logs
	@docker-compose logs -f

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

. DEFAULT_GOAL := help
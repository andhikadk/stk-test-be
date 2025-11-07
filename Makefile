.PHONY: help build run test clean migrate seed docker-up docker-down docker-logs docker-reset docker-dev docker-dev-logs docker-dev-down docker-dev-reset install-deps

# Variables
APP_NAME=go-fiber-boilerplate
MAIN_PATH=main.go
BINARY_NAME=./bin/$(APP_NAME)

help: ## Display this help screen
	@echo "Available commands:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-20s %s\n", $$1, $$2}'

install-deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_PATH)

dev: ## Run in development mode with hot reload (requires air)
	@echo "Running in development mode..."
	@air || echo "air not installed. Install with: go install github.com/cosmtrek/air@latest"

test: ## Run unit tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "Clean complete"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run ./... || echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

migrate: ## Run database migrations (AutoMigrate for dev, SQL for prod)
	@echo "Running migrations..."
	@go run $(MAIN_PATH) -migrate

migrate-sql: ## Run SQL migrations from files
	@echo "Running SQL migrations..."
	@go run $(MAIN_PATH) -migrate=sql

migrate-status: ## Show migration status
	@echo "Migration status..."
	@go run $(MAIN_PATH) -status

seed: ## Seed database with sample data
	@echo "Seeding database..."
	@go run $(MAIN_PATH) -seed

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .

docker-up: ## Start Docker containers (docker compose)
	@echo "Starting Docker containers..."
	@docker compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker compose down

docker-logs: ## View Docker logs
	@docker compose logs -f

docker-reset: ## Reset Docker containers and volumes (removes all data)
	@echo "Resetting Docker containers and volumes..."
	@docker compose down -v
	@echo "Containers and volumes removed. Restart with: make docker-up"

docker-dev: ## Start Docker containers with hot reload (development mode)
	@echo "Starting Docker containers with hot reload..."
	@docker compose -f docker-compose.dev.yml up -d

docker-dev-logs: ## View Docker development logs
	@docker compose -f docker-compose.dev.yml logs -f

docker-dev-down: ## Stop Docker development containers
	@echo "Stopping Docker development containers..."
	@docker compose -f docker-compose.dev.yml down

docker-dev-reset: ## Reset Docker development containers and volumes
	@echo "Resetting Docker development containers and volumes..."
	@docker compose -f docker-compose.dev.yml down -v
	@echo "Development containers and volumes removed. Restart with: make docker-dev"

all: clean install-deps build test ## Clean, install, build and test

.DEFAULT_GOAL := help

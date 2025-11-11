.PHONY: help build run test clean sqlc-generate migrate-up migrate-down docker-build docker-run

help: ## Show this help
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/api ./cmd/api

run: ## Run the application
	go run ./cmd/api/main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/

sqlc-generate: ## Generate sqlc code
	sqlc generate

sqlc-verify: ## Verify sqlc configuration
	sqlc verify

migrate-up: ## Run database migrations
	psql $(DATABASE_URL) -f db/migrations/001_init_schema.sql

docker-build: ## Build Docker image
	docker build -t cowatching-backend .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file ../.env cowatching-backend

deps: ## Download dependencies
	go mod download
	go mod tidy

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

dev: ## Run with hot reload (requires air)
	air

.DEFAULT_GOAL := help

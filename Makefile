# Makefile for Transaction Processor

.PHONY: help setup run-api run-worker test lint clean docker-up docker-down aws-setup

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Install dependencies and setup .env
	go mod download
	@if [ ! -f .env ]; then cp .env.example .env; fi
	@echo "Setup complete."

docker-up: ## Start LocalStack containers
	docker-compose up -d
	@echo "Waiting for LocalStack to be ready..."
	@until curl -s http://localhost:4566/_localstack/health | grep "\"sns\": \"running\"" > /dev/null; do sleep 2; done
	@$(MAKE) aws-setup
	@echo "LocalStack is ready."

docker-down: ## Stop LocalStack containers
	docker-compose down

aws-setup: ## Setup AWS resources in LocalStack
	@chmod +x scripts/setup-aws.sh
	./scripts/setup-aws.sh

run-api: ## Run API application
	go run cmd/api/main.go

run-worker: ## Run Worker application
	go run cmd/worker/main.go cmd/worker/worker_base.go

test: ## Run tests
	go test ./... -v -cover

lint: ## Run linter
	golangci-lint run

build: ## Build binaries
	@mkdir -p bin
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go cmd/worker/worker_base.go

clean: ## Remove binaries and temp files
	rm -rf bin/
	go clean

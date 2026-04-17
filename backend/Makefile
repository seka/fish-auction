# Default target
.DEFAULT_GOAL := help

# Server settings
PORT ?= 8080

# Go cache and module cache confinement
# Use project-local .cache directory for reliability in restricted environments
export GOCACHE ?= $(shell pwd)/.cache/go-build
export GOMODCACHE ?= $(shell pwd)/.cache/go-mod
export GOTMPDIR ?= $(shell pwd)/.cache/tmp
export PATH := $(shell pwd)/build/.gopath/bin:$(PATH)

# Ensure cache directories exist
$(shell mkdir -p $(GOCACHE) $(GOMODCACHE) $(GOTMPDIR) 2>/dev/null || true)
$(shell mkdir -p build/.gopath/bin 2>/dev/null || true)

# ============================================================================
# Development tools setup
# ============================================================================

.PHONY: bootstrap
bootstrap: ## Install all development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest
	@echo "All tools installed successfully!"

# ============================================================================
# Development and Build
# ============================================================================

.PHONY: app.run
app.run: ## Run the API server with live reload (Air)
	air -c .air.toml

.PHONY: app.build
app.build: ## Build the API server binary
	@echo "Building binary..."
	go build -o build/main ./cmd/server/main.go

.PHONY: app.clean
app.clean: ## Clean build artifacts
	go clean
	rm -rf build/*

# ============================================================================
# Code Quality
# ============================================================================

.PHONY: app.fmt
app.fmt: ## Format code
	@echo "Running go fmt..."
	go fmt ./...

.PHONY: app.lint
app.lint: ## Run linter
	@echo "Running golangci-lint..."
	golangci-lint run --allow-parallel-runners ./...

.PHONY: app.lint.fix
app.lint.fix: ## Run linter and auto-fix
	@echo "Running go fix and golangci-lint run --fix..."
	go fix ./...
	golangci-lint run --fix --allow-parallel-runners ./...

.PHONY: app.test
app.test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

.PHONY: app.integration-test
app.integration-test: ## Run integration tests with dynamic DB lifecycle
	@TEST_DB_NAME=fish_auction_test_$$(date +%s); \
	TEST_DB_USER=$${POSTGRES_USER:-postgres}; \
	TEST_DB_PASSWORD=$${POSTGRES_PASSWORD:-postgres}; \
	TEST_DB_HOST=$${POSTGRES_HOST:-localhost}; \
	TEST_DB_PORT=$${POSTGRES_PORT:-5432}; \
	echo "Preparing integration test environment (DB: $$TEST_DB_NAME)..."; \
	PGPASSWORD=$$TEST_DB_PASSWORD psql -h $$TEST_DB_HOST -U $$TEST_DB_USER -d postgres -c "CREATE DATABASE $$TEST_DB_NAME;"; \
	echo "Running integration tests..."; \
	export RUN_INTEGRATION_TESTS=true; \
	export POSTGRES_HOST=$$TEST_DB_HOST; \
	export POSTGRES_PORT=$$TEST_DB_PORT; \
	export POSTGRES_USER=$$TEST_DB_USER; \
	export POSTGRES_PASSWORD=$$TEST_DB_PASSWORD; \
	export POSTGRES_DB=$$TEST_DB_NAME; \
	export POSTGRES_SSLMODE=disable; \
	export REDIS_HOST=localhost; \
	export REDIS_PORT=6379; \
	export REDIS_DB=15; \
	export SERVER_HOST=0.0.0.0; \
	export SERVER_PORT=18080; \
	export APP_ENV=test; \
	export SMTP_HOST=localhost; \
	export SMTP_PORT=1025; \
	export SMTP_FROM=test@example.com; \
	export FRONTEND_URL=https://localhost; \
	export ALLOWED_ORIGINS=https://localhost,http://localhost:3000; \
	go test ./cmd/testing/... -run TestServerIntegration -v; \
	EXIT_CODE=$$?; \
	echo "Cleaning up test environment..."; \
	PGPASSWORD=$$TEST_DB_PASSWORD psql -h $$TEST_DB_HOST -U $$TEST_DB_USER -d postgres -c "DROP DATABASE $$TEST_DB_NAME;"; \
	exit $$EXIT_CODE

# Stress test parameters
LOAD_TEST_CONCURRENCY ?= 100
LOAD_TEST_REQUESTS ?= 1000
LOAD_TEST_TARGET_URL ?= http://localhost:8080

.PHONY: app.stress-test
app.stress-test: ## Run stress tests
	@echo "Running stress tests (Concurrency: $(LOAD_TEST_CONCURRENCY), Requests: $(LOAD_TEST_REQUESTS))..."
	LOAD_TEST_CONCURRENCY=$(LOAD_TEST_CONCURRENCY) \
	LOAD_TEST_REQUESTS=$(LOAD_TEST_REQUESTS) \
	LOAD_TEST_TARGET_URL=$(LOAD_TEST_TARGET_URL) \
	go test ./cmd/testing/... -run TestLoadTest -v

# ============================================================================
# Database and Seeding
# ============================================================================

.PHONY: seed
seed: ## Run database seed
	@echo "Running database seed..."
	go run ./cmd/seed/main.go

.PHONY: init-admin
init-admin: ## Initialize admin user (Usage: make init-admin EMAIL=... PASSWORD=...)
	@echo "Initializing admin user..."
	go run ./cmd/init_admin/main.go $(if $(EMAIL),--email $(EMAIL)) $(if $(PASSWORD),--password $(PASSWORD))

# ============================================================================
# Utility targets
# ============================================================================

.PHONY: help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'

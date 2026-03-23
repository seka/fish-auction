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
app.integration-test: ## Run integration tests
	@echo "Running integration tests..."
	go test ./cmd/testing/... -run TestServerIntegration -v

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
init-admin: ## Initialize admin user
	@echo "Initializing admin user..."
	go run ./cmd/init_admin/main.go

# ============================================================================
# Utility targets
# ============================================================================

.PHONY: help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'

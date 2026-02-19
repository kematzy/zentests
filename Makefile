# Makefile for Zen Starter Go Web Project

.PHONY: help test test-verbose fmt vet lint check tidy setup

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

# Testing commands
test: ## Run all tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-coverage: ## Generate test coverage report
	go test -coverprofile=.code-status/coverage.out ./...
	go tool cover -html=.code-status/coverage.out -o .code-status/coverage.html
	@echo "Coverage report generated: coverage.html"


# Code quality commands
fmt: ## Format Go code
	gofmt -w .

vet: ## Run go vet
	go vet ./...

lint: ## Run go linting
	golangci-lint run ./...

check: fmt format-check vet lint test ## Run all checks (format, vet, lint, test)

tidy: ## Tidy Go modules
	go mod tidy

# Setup commands
setup: ## Initial project setup
	go mod tidy
	@echo "Project setup complete. Run 'make dev' to start development server."

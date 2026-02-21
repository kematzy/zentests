# Makefile for zentests

#======================================================================================================================
# VARIABLES
#======================================================================================================================

# Set the date for
DATE := $(shell date +'%Y-%m-%d')
# Read repository from the first line of the `go.mod` file.
REPO := $(shell head -1 go.mod | awk '{print $$2}')

#======================================================================================================================
# MACROS
#======================================================================================================================

define banner
	@echo ""
	@echo -e "\033[2m--\033[1m ${1} \033[0m---------------------------------\033[0m"
	@echo ""
endef

#======================================================================================================================
# DEFAULT
#======================================================================================================================

.DEFAULT_GOAL := help

# non-file targets (commands)
.PHONY: help test test-verbose test-coverage fmt vet lint tidy

## help: Show this help message
help:
	$(call banner,Available commands:)
	@grep -E '^## ' Makefile \
		| sed 's/^## //' \
		| awk -F': ' '{ printf " \033[32m\033[1m%-22s\033[0m \033[2m%s\033[0m\n", $$1, $$2 }'
	@echo -e "\033[0m"


#======================================================================================================================
# CODE QUALITY COMMANDS
#======================================================================================================================

## fmt: Format source code
fmt:
	@gofmt -l -w .

## lint: Run linter
lint:
	@golangci-lint run ./...

## tidy: Run tidy on Go modules
tidy:
	@go mod tidy

## vet: Run `go vet`
vet:
	@go vet ./...

#======================================================================================================================
# TEST COMMANDS
#======================================================================================================================

## test: Run all tests
test:
	@gotest ./...

## test-verbose: Run all tests with verbose colored output
test-verbose:
	@echo ""
	@gotestsum --format testdox

# @gotest -v ./...

## test-coverage: Run tests and generate coverage report
test-coverage:
	@go test -coverprofile=.code-status/coverage.out ./...
	@go tool cover -html=.code-status/coverage.out -o code-coverage.html
	@echo -e "\033[32mCoverage report generated: code-coverage.html \033[0m"


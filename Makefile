# Makefile for zentests

#======================================================================================================================
# VARIABLES
#======================================================================================================================

# Set the current version or default
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
# Get the current commit hash
COMMIT_HASH := $(shell git rev-parse --short HEAD)
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

## modernize: Run tool to identify & replace old Go code with newer standards.
modernize:
	@go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix ./...

## modernize-check: Dry-run with diffs to identify old Go code with newer standards.
modernize-check:
	@go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -diff -fix ./...

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

#======================================================================================================================
# CHANGELOG
#======================================================================================================================

## changelog: Generate changelog using git-cliff
changelog:
	@git-cliff --tag ${VERSION} --output CHANGELOG.md


#======================================================================================================================
# VERSIONS
#======================================================================================================================

## versions: Show current and all versions
versions:
	@echo -e "Version: \033[32m${VERSION}\033[0m \033[02m- #\033[0m\033[32m${COMMIT_HASH}\033[0m"
	@echo "Versions: "
	@git tag -l --sort=version:refname | sed 's/^v/ v/'


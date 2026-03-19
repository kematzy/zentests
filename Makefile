# Makefile for zentests

#======================================================================================================================
# VARIABLES
#======================================================================================================================

# Set the current version or default
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
# Set the minimun coverage percentage
COVERAGE_MIN ?= 90
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
.PHONY: help test test-verbose test-coverage fmt vet lint tidy check release changelog git-push gitea-push github-push docs docs-md modernize modernize-check

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
	@gotestsum --format testname

# @gotest -v ./...

## test-coverage: Run tests and generate coverage report
test-coverage:
	@go test -coverprofile=.code-status/coverage.out ./...
	@go tool cover -html=.code-status/coverage.out -o code-coverage.html
	@echo -e "\033[32mCoverage report generated: code-coverage.html \033[0m"


#======================================================================================================================
# RELEASE
#======================================================================================================================

## check: Run all checks (fmt, vet, lint, check-test-coverage)
check: fmt vet lint check-test-coverage
	@echo "All checks passed."
	@echo ""

## check-test-coverage: Run tests and enforce coverage threshold (default: 90%)
check-test-coverage:
	@go test -coverprofile=.code-status/coverage.out ./...
	@go tool cover -func=coverage.out | tail -1 | awk '{print $$3}' | \
	  tr -d '%' | awk -v threshold=${COVERAGE_MIN} \
	  '{ if ($$1 < threshold) \
	       { print "FAIL: Coverage " $$1 "% is below threshold " threshold "%"; exit 1 } \
	     else \
	       { print "OK: Coverage " $$1 "%" } }'

## release: Release a new tagged version  (e.g. make release VERSION=v1.2.3)
release: check changelog-release
	@[ "${VERSION}" ] || { echo "Usage: make release VERSION=v1.2.3"; exit 1; }

	# Normalize: add 'v' prefix if not already present
	$(eval TAG := $(shell echo "${VERSION}" | grep -q '^v' && echo "${VERSION}" || echo "v${VERSION}"))

	@echo ""
	@echo -e "Starting release: \033[32m${TAG}\033[0m"
	@echo ""
	@git tag -a ${TAG} -m "Release ${TAG}"
	@git remote | xargs -I% git push % ${TAG}
	@echo ""
	@echo -e "\033[32mTag - ${TAG} - pushed to remote(s) & published.\033[0m"
	@echo ""
	@echo "To install:"
	@echo "  go get $(REPO)@${TAG}"
	@echo ""
	@echo "Docs:"
	@echo "  https://pkg.go.dev/$(REPO)@${TAG}"
	@echo "  https://$(REPO)/releases/tag/${TAG}"
	@echo ""

#======================================================================================================================
# CHANGELOG
#======================================================================================================================

## changelog: Generate changelog using git-cliff
changelog:
	@git-cliff --tag ${VERSION} --output CHANGELOG.md

## changelog-release: Generate changelog using git-cliff & Git committ it
changelog-release: changelog
	@git add CHANGELOG.md



#======================================================================================================================
# VERSIONS
#======================================================================================================================

## versions: Show current and all versions
versions:
	@echo -e "Version: \033[32m${VERSION}\033[0m \033[02m- #\033[0m\033[32m${COMMIT_HASH}\033[0m"
	@echo "Versions: "
	@git tag -l --sort=version:refname | sed 's/^v/ v/'


#======================================================================================================================
# GIT
#======================================================================================================================

## git-push: Push code to both repositories
git-push: github-push gitea-push
	@echo "Successfully pushed both remotes."

github-push:
	@echo "Pushing to Github"
	@git push github HEAD

gitea-push:
	@echo "Pushing to Gitea"
	@git push gitea HEAD


#======================================================================================================================
# DEVELOPMENT
#======================================================================================================================

## docs: Open documentation in browser
docs:
	@if command -v pkgsite >/dev/null 2>&1; then \
		echo "Opening docs at http://localhost:8080/pkg/$(REPO)"; \
		xdg-open "http://localhost:8080/pkg/$(REPO)" 2>/dev/null \
			|| open "http://localhost:8080/pkg/$(REPO)" 2>/dev/null \
			|| true; \
		pkgsite -http=:8080; \
	else \
		echo "pkgsite not found. Install with: go install golang.org/x/pkgsite/cmd/pkgsite@latest"; \
		exit 1; \
	fi

## docs-md: Create Markdown API doc with `gomarkdoc` in `docs/API.md`
# See: https://github.com/princjef/gomarkdoc
docs-md:
	@if command -v gomarkdoc >/dev/null 2>&1; then \
		echo "Creating docs/API.md from the code"; \
		gomarkdoc ./... > docs/API.md; \
	else \
		echo "gomarkdoc not found. Install with: go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest"; \
		exit 1; \
	fi


## clean: Clean build artifacts
clean:
	@echo -e "\033[2mCleaning... \033[0m"
	@rm -f coverage.out .code-status/coverage.out
	@echo -e "\033[32mDone! \033[0m"

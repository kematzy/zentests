# Makefile for zentests

#======================================================================================================================
# VARIABLES
#======================================================================================================================

# Set the current version or default
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
# Set the minimum coverage percentage
COVERAGE_MIN ?= 95
# Coverage profile written by coverage-profile and read by the coverage targets
COVERAGE_PROFILE := .code-status/coverage.out
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
.PHONY: help fmt fmt-check vet lint tidy modernize modernize-check check spellcheck
.PHONY: test test-verbose test-dox test-watch coverage-profile test-coverage
.PHONY: release changelog git-push gitea-push github-push docs docs-md

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
	@# -l	list files whose formatting differs from gofmt's
	@# -w	write result to (source) file instead of stdout

## fmt-check: Check source code formatting without modifying files
fmt-check:
	@files=$$(gofmt -l -d -e -s .); \
	if [ -n "$$files" ]; then \
		echo "The following files need formatting:"; \
		printf '%s\n' "$$files"; \
		exit 1; \
	fi
	@# -d	display diffs instead of rewriting files
	@# -e	report all errors (not just the first 10 on different lines)
	@# -l	list files whose formatting differs from gofmt's
	@# -s	simplify code

## lint: Run linter
lint:
	@golangci-lint run --color=always ./...

## tidy: Run tidy on Go modules
tidy:
	@go mod tidy

## vet: Run `go vet`
vet:
	@go vet ./...

## modernize: Identify & replace old Go code with newer standards.
modernize:
	@go fix ./...

## modernize-check: (Dry-run) Identify old Go code with newer standards with diff.
modernize-check:
	@go fix -diff ./... | delta

## check: Check everything before a commit or push
check: check-test-coverage fmt-check
	@echo -e "\nChecked test suite: 			󰄬 (passed)"
	@echo -e "Checking code formatting:		󰄬 (passed)"
	@go vet ./...
	@echo -e "Checked for suspicious code:		󰄬 (passed)"
	@golangci-lint run ./...
	@echo -e "Checked code linting:			󰄬 (passed)\n"
	@cspell --unique --color --no-progress --quiet .
	@echo -e "\nChecked spelling:			󰄬 (passed)"

	@echo -e "\n\033[32m 󰄬   All checks passed \033[0m"

## spellcheck: Check spelling (cspell)
spellcheck:
	@cspell --unique --color --no-progress .

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

## test-dox: Run all tests with testdox colored output
test-dox:
	@echo ""
	@gotestsum --format testdox

## test-watch: Run tests when a file is saved. Reruns failed tests 4 times
test-watch:
	@echo ""
	@gotestsum --watch --format dots

## coverage-profile: Run tests and write the coverage profile
coverage-profile:
	@mkdir -p .code-status
	@echo ""
	@gotest -coverprofile=$(COVERAGE_PROFILE) ./...

## test-coverage: Run tests and generate coverage report
test-coverage:
	@mkdir -p .code-status
	@echo ""
	@gotest -coverprofile=$(COVERAGE_PROFILE) ./...
	@go-better-html-coverage -profile $(COVERAGE_PROFILE) -o code-coverage.html
	@echo -e "\n\033[32mCoverage report generated: code-coverage.html \033[0m\n"
	@#
	@# depends upon: go install github.com/chmouel/go-better-html-coverage@latest


#======================================================================================================================
# RELEASE
#======================================================================================================================

## check-test-coverage: Run tests and enforce coverage threshold (default: 95%)
check-test-coverage: coverage-profile
	@go tool cover -func=$(COVERAGE_PROFILE) | tail -1 | awk '{print $$3}' | \
	  tr -d '%' | awk -v threshold=${COVERAGE_MIN} \
	  '{ if ($$1 < threshold) \
		   { print "FAIL: Coverage " $$1 "% is below threshold " threshold "%"; exit 1 } \
		 else \
		   { print "OK: Coverage " $$1 "%" } }'

## release: Release a new tagged version  (e.g. make release VERSION=v1.2.3)
release: check changelog
	@[ "${VERSION}" ] || { echo "Usage: make release VERSION=v1.2.3"; exit 1; }
	@echo ""
	@echo -e "Starting release: \033[32m${VERSION}\033[0m"
	@echo ""
	@git tag -a ${VERSION} -m "Release ${VERSION}"
	@git remote | xargs -I% git push % ${VERSION}
	@echo ""
	@echo -e "\033[32mTag - ${VERSION} - pushed to remote(s) & published.\033[0m"
	@echo ""
	@echo "To install:"
	@echo "  go get $(REPO)@${VERSION}"
	@echo ""
	@echo "Docs:"
	@echo "  https://pkg.go.dev/$(REPO)@${VERSION}"
	@echo "  https://$(REPO)/releases/tag/${VERSION}"
	@echo ""


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
	@rm -f coverage.out $(COVERAGE_PROFILE)
	@echo -e "\033[32mDone! \033[0m"

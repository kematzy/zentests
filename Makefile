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
# HTML coverage report written by test-coverage (kept out of the repo, see .gitignore)
COVERAGE_REPORT := .code-status/coverage.html
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
.PHONY: examples-check
.PHONY: test test-verbose test-dox test-watch coverage-profile test-coverage
.PHONY: release changelog docs docs-md

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
check: check-test-coverage fmt-check examples-check
	@echo -e "\nChecked test suite: 			󰄬 (passed)"
	@echo -e "Checking code formatting:		󰄬 (passed)"
	@go vet ./...
	@echo -e "Checked for suspicious code:		󰄬 (passed)"
	@golangci-lint run ./...
	@echo -e "Checked code linting:			󰄬 (passed)\n"
	@cspell --unique --color --no-progress --quiet .
	@echo -e "\nChecked spelling:			󰄬 (passed)"

	@echo -e "\n\033[32m 󰄬   All checks passed \033[0m"

## examples-check: Compile and run Go examples embedded in Markdown
examples-check:
	@bash scripts/check-examples.sh

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
test-coverage: coverage-profile
	@go-better-html-coverage -profile $(COVERAGE_PROFILE) -o $(COVERAGE_REPORT)
	@echo -e "\n\033[32mCoverage report generated: $(COVERAGE_REPORT) \033[0m\n"
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
	# Commit the generated CHANGELOG before tagging, but only when it changed,
	# so the tag contains the changelog for its own version.
	@git add CHANGELOG.md
	@git diff --cached --quiet || git commit -m "docs(changelog): update CHANGELOG for ${VERSION}"
	# Create annotated tag
	@git tag -a ${VERSION} -m "Release ${VERSION}"
	# Push tags first, then the current branch, to every configured remote.
	@git remote | xargs -I% git push % ${VERSION}
	@git remote | xargs -I% git push %
	@echo ""
	@echo -e "\033[32mRelease ${VERSION} committed and pushed.\033[0m"
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

## changelog: Add commits since the latest tag to CHANGELOG.md (set VERSION)
changelog:
	@if git rev-parse --quiet --verify "refs/tags/${VERSION}" >/dev/null; then \
		echo "Tag ${VERSION} already exists; choose a new VERSION." >&2; exit 1; \
	fi
	@version=$$(printf '%s' "${VERSION}" | sed 's/^v//'); \
	if grep -Fq "## [$$version]" CHANGELOG.md; then \
		echo "CHANGELOG.md already contains version $$version; refusing to duplicate it." >&2; exit 1; \
	fi
	@git-cliff --unreleased --tag "${VERSION}" --prepend CHANGELOG.md
	@cspell --unique --color --no-progress --quiet CHANGELOG.md


#======================================================================================================================
# VERSIONS
#======================================================================================================================

## versions: Show current and all versions
versions:
	@echo -e "Version: \033[32m${VERSION}\033[0m \033[02m- #\033[0m\033[32m${COMMIT_HASH}\033[0m"
	@echo "Versions: "
	@git tag -l --sort=version:refname | sed 's/^v/ v/'


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
	@rm -f $(COVERAGE_PROFILE) $(COVERAGE_REPORT)
	@echo -e "\033[32mDone! \033[0m"

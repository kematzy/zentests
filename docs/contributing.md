# Zentests - Docs - Contributing

**Contributions are welcome!** Please feel free to submit a Pull Request.

## Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/[yourusername]/zentests.git`
3. Install dependencies: `go mod download`
4. Run tests: `make test`, `make test-verbose`
5. Check coverage: `make test-coverage`

## Guidelines

- **Code Quality**: Maintain the existing code style and documentation standards
- **Tests**: All new features must include comprehensive tests (aim for 100% coverage)
- **Documentation**: Update the `README.md` with any new features or changes
- **Commits**: Use clear, descriptive commit messages
- **Compatibility**: Ensure changes work with the latest stable Go and Fiber versions

## Reporting Issues

Please use GitHub Issues to report bugs or request features. When reporting bugs, include:

- Go version
- Fiber version
- Minimal code example that reproduces the issue
- Expected vs actual behavior

## Branches

| Branch     | Fiber Version | Go Version | Status             |
| ---------- | ------------- | ---------- | ------------------ |
| `master`   | v3.x          | 1.25+      | Active development |
| `fiber-v2` | v2.x          | 1.21+      | Maintenance only   |

The `fiber-v2` branch may receive bug fixes and security updates only. New features are developed on
`master` (Fiber v3).

## Make Commands

```bash
make help                   # Show this help message

# Testing
make test                   # Run all tests
make test-verbose           # Run all tests with verbose colored output
make test-coverage          # Run tests and generate coverage report

# Code Quality
make fmt                    # Format source code
make lint                   # Run linter
make tidy                   # Run tidy on Go modules
make vet                    # Run `go vet`

make modernize              # Run tool to identify & replace old Go code with newer standards.
make modernize-check        # Dry-run with diffs to identify old Go code with newer standards.

make check                  # Run all checks (fmt, vet, lint, check-test-coverage)
make check-test-coverage    # Run tests and enforce coverage threshold (default: 95%)

# Release
make release                # Release a new tagged version (e.g. make release VERSION=v1.2.3)
make release VERSION=v1.2.3 # Create specific version

# Development
make versions               # Show current and all versions

# Documentation
make changelog              # Generate changelog using git-cliff
make docs                   # Open documentation in browser
make docs-md                # Create Markdown API doc with `gomarkdoc` in `docs/API.md`

make clean                  # Clean build artifacts
```

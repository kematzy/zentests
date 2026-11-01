# AGENTS.md - Agentic Coding Guidelines

This file provides guidelines for agentic coding agents operating in this repository.

## Project Overview

`zentests` is a testing utility library for Go Fiber applications. It provides a fluent API for making HTTP requests and asserting responses in tests.

- **Language**: Go 1.25.5
- **Module**: `github.com/kematzy/zentests`
- **Dependencies**: fiber/v3, testify, gorm

---

## Build/Lint/Test Commands

### Running Tests

```bash
go test ./...                      # Run all tests
go test -run TestName ./...        # Run single test by name
go test -run TestName ./file.go   # Run test in specific file
go test -v ./...                  # Run with verbose output
go test -coverprofile=out ./...  # Run with coverage report
gotestsum --format testname       # Run with gotestsum (colored)
```

### Code Quality

```bash
gofmt -l -w .                     # Format code
go vet ./...                      # Run go vet
golangci-lint run ./...           # Run linter
go mod tidy                      # Go mod tidy
make check                       # Full check (fmt, vet, lint, coverage)
```

### Make Targets

```bash
make help          # Show all commands
make test          # Run all tests
make test-verbose  # Run with gotestsum
make test-coverage # Generate coverage report
make test-watch    # Watch mode (auto-run on changes)
make lint          # Run linter
make check         # Full check (fmt, vet, lint, coverage)
make modernize     # Modernize code to newer Go standards
```

---

## Code Style Guidelines

### Formatting

- **Indentation**: Tabs (width 4), configured in `.editorconfig`
- **Line endings**: LF (Unix-style)
- **File endings**: Newline at end of file
- **Trailing whitespace**: Trimmed

### Imports
Standard library first, then third-party. No blank imports.

### Naming
- **Packages**: lowercase (e.g., `zentests`)
- **Types**: PascalCase (e.g., `T`, `Response`, `SuiteContext`)
- **Methods/Functions**: camelCase
- **Constants**: PascalCase (e.g., `statusOK`)
- **Errors**: Prefix with `Err` (e.g., `ErrInvalidInput`)

### Documentation
All exported functions/types must have doc comments with Description, Parameters, Returns, Example sections.

### Code Structure
- Max function length: 100 lines, 40 statements
- Max cyclomatic complexity: 15
- Max cognitive complexity: 20

### Error Handling
- Use `testify/assert` for test assertions
- Wrap errors with context: `fmt.Errorf("context: %w", err)`
- Use `errorlint` for proper error wrapping checks
- Name errors with `Err` prefix

---

## Testing Guidelines

### Test Organization
- Test files: `*_test.go` in same package
- Test names: `TestFunctionName` or `TestType_Method`
- Use table-driven tests where appropriate

### Response Assertions
```go
// Status checks
resp.OK()           // 200
resp.Created()      // 201
resp.NoContent()    // 204
resp.BadRequest()   // 400
resp.Unauthorized() // 401
resp.NotFound()     // 404

// Body checks
resp.IsJSON()            // Content-Type: application/json
resp.Has("path", value)  // JSON path assertion
resp.Body()              // Get raw body as []byte
resp.JSON()              // Parse body as JSON
```

### BDD-Style Testing
```go
zt := zentests.New(t)
zt.Describe("User API", func(ctx *zentests.SuiteContext) {
    ctx.BeforeEach(func(zt *zentests.T) {
        app = setupTestApp()
    })
    ctx.It("creates a user", func(zt *zentests.T) {
        zt.PostJSON(app, "/users", data).Created()
    })
})
```

### GORM Test Helpers
```go
// Setup
db := zentests.SetupTestDBWithModels(t, &User{})
// Create
user := zentests.DBCreate(t, db, &User{Name: "Alice"})
// Query
zentests.DBFind(t, db, &User{}, user.ID)
count := zentests.DBCount(t, db, &User{})
zentests.DBExists(t, db, &User{}, "status = ?", "active")
// Seed
zentests.DBSeed(t, db, []*User{{Name: "A"}, {Name: "B"}})
```

---

## Linter Config

Enabled in `.golangci.yml`:
- errcheck, ineffassign, staticcheck, gosec
- gocyclo (max 15), funlen (100 lines, 40 stmts), gocognit (max 20)
- bodyclose, noctx, rowserrcheck, sqlclosecheck

Exclusions: Test files exclude bodyclose, godoclint, funlen, goconst

---

## Coverage Requirements
- **Minimum coverage**: 90%
- Run: `make check-test-coverage`

---

## Additional Notes
- Targets Go Fiber v3
- Uses testify for assertions
- GORM for database testing support
- HTTP response bodies read lazily (not closed immediately)
- Call Fiber app shutdown in test teardown if needed

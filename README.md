# zentests

Readable Go tests without magic. Inspired by RSpec/Pest

`zentests` provides fluent, chainable HTTP testing for [Fiber](https://gofiber.io/) applications with RSpec-style syntax.

[![Go Report Card](https://goreportcard.com/badge/github.com/kematzy/zentests)](https://goreportcard.com/report/github.com/kematzy/zentests)
[![GoDoc](https://godoc.org/github.com/kematzy/zentests?status.svg)](https://godoc.org/github.com/kematzy/zentests)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Why `zentests`?

- **No magic**: Standard Go testing, no reflection tricks
- **Type safe**: Strict assertions catch bugs early
- **Fiber-native**: Built for Fiber, not adapted from net/http
- **Chainable**: Readable, expressive test code
- **Familiar**: RSpec/Pest-inspired syntax
- **Well-tested**: 98.2% code coverage

### Why only Fiber support?

Simple. I'm using Fiber for my projects and created these tests to make working on my projects faster.

You are welcome to create a version with support for your preferred routing package.

<br>

## Installation

```bash
go get github.com/kematzy/zentests
```

**Requirements:**
- Go 1.25 or higher
- [Fiber](https://github.com/gofiber/fiber) v2

<br>

## Quick Start

```go
package myapp_test

import (
    "testing"
    "github.com/gofiber/fiber/v3"
    "github.com/kematzy/zentests"
)

func TestAPI(t *testing.T) {
    zt := zentests.New(t)
    app := fiber.New()
    
    // Your routes here...
    app.Get("/api/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok"})
    })
    
    // Fluent, chainable assertions
    zt.Get(app, "/api/health").
        OK().
        IsJSON().
        Has("status", "ok")
}

```

## Code Quality & Coverage

**Coverage: 98.2%** of statements

All functions are fully documented with comprehensive examples. See our [test coverage report](./.code-status/coverage.html) for details.

<br>

## Examples

The [`examples/`](./examples/) directory contains practical examples demonstrating different testing styles:

### [flat_style.go](./examples/flat_style.go) - Simple Testing
Demonstrates simple, flat-style testing without suites. Perfect for quick API tests.

```go
// Simple, no suite overhead
zt.Get(app, "/api/health").OK().IsJSON().Has("status", "ok")

zt.PostJSON(app, "/api/users", userData).
    Created().
    Has("id", float64(1)).
    Has("name", "John Doe")
```

### [routes.go](./examples/routes.go) - JSON Testing
Shows comprehensive JSON assertions with dot notation for nested objects and arrays.

```go
zt.Get(app, "/api/users").
    OK().
    IsJSON().
    Has("success", true).
    Has("data.count", 2).
    ArrayLength("data.users", 2).
    Has("data.users.0.name", "John") // dot notation
```

### [suite.go](./examples/suite.go) - BDD Style Testing
Demonstrates BDD-style testing with `Describe`, `It`, and lifecycle hooks (`BeforeEach`, `AfterEach`).

```go
zt.Describe("User API", func(ctx *zentests.SuiteContext) {
    ctx.
        BeforeEach(func(zt *zentests.T) {
            // Fresh app for each test
            app = fiber.New()
            SetupRoutes(app)
        }).
        AfterEach(func(zt *zentests.T) {
            // Cleanup
            app.Shutdown()
        })

    ctx.It("returns user list", func(zt *zentests.T) {
        zt.Get(app, "/api/users").
            OK().
            IsJSON().
            ArrayLength("data.users", 2)
    })
})
```

<br>

## Features

### HTTP Methods

All methods accept `*fiber.App` and `path`, return `*Response` for chaining:

| Method | Description |
|--------|-------------|
| `Get(app, path)` | GET request  |
| `Post(app, path, body)` | POST with raw body  |
| `PostJSON(app, path, data)` | POST with JSON body  |
| `PostForm(app, path, data)` | POST with form data  |
| `Put(app, path, body)` | PUT request  |
| `PutJSON(app, path, data)` | PUT with JSON  |
| `Patch(app, path, body)` | PATCH request  |
| `PatchJSON(app, path, data)` | PATCH with JSON  |
| `Delete(app, path)`	| DELETE request  |
| `DeleteJSON(app, path, data)` | DELETE with JSON body  |


### Status Assertions

```go
resp.OK()              // 200
resp.Created()         // 201
resp.Accepted()        // 202
resp.NoContent()       // 204
resp.BadRequest()      // 400
resp.Unauthorized()    // 401
resp.Forbidden()       // 403
resp.NotFound()        // 404
resp.Unprocessable()   // 422
resp.ServerError()     // 500
resp.Status(code)      // Custom status
```

### Header Assertions

```go
resp.HasHeader("X-Custom", "value")                                     // Exact match
resp.HeaderContains("X-Rate", "limit")                                  // Substring match
resp.HeaderPresent("X-Auth-Token")                                      // Header exists
resp.HeaderNotPresent("X-Auth-Token")                                   // Header does not exists
resp.HeaderHasValues("X-Features", []string{"audio", "video"})          // Header has multiple values
resp.CookieHasValues(map[string]string{"lang": "en", "theme": "dark"})  // Set-Cookie header test

// Content-Type Helpers
resp.HasContentType("text/html; charset=utf-8")                         // Checks Content-Type value

// Content-Type Shortcuts
resp.IsJSON()                                                           // Content-Type: application/json
resp.IsHTML()                                                           // Content-Type: text/html
resp.IsPlainText()                                                      // Content-Type: text/plain
resp.IsCSS()                                                            // Content-Type: text/css
resp.IsJS()                                                             // Content-Type: (application|text)/javascript
resp.IsXML()                                                            // Content-Type: (application|text)/xml
resp.IsImage()                                                          // Content-Type: image/*
resp.IsPNG()                                                            // Content-Type: image/png
resp.IsJPEG()                                                           // Content-Type: image/jpeg
resp.IsGIF()                                                            // Content-Type: image/gif
resp.IsSVG()                                                            // Content-Type: image/svg+xml
resp.IsWebP()                                                           // Content-Type: image/webp

// AJAX Detection:
resp.IsXHR()  // checks for "XMLHttpRequest" in X-Requested-With header
```

### Body Assertions

```go
resp.Contains("substring")             // Body contains text
resp.NotContains("substring")          // Body excludes text
resp.BodyMatches(`regex\d+`)           // Regex match
resp.Equals("exact body")              // Exact match
resp.IsEmpty()                         // Body is empty

```

### JSON Assertions

Dot notation for nested keys: `user.profile.name`

```go
resp.JSON()                            // Parse & cache JSON
resp.HasKey("user.name")               // Key exists
resp.Has("user.name", "John")          // Strict type match
resp.HasString("user.name", "John")    // String value
resp.HasInt("user.age", 30)            // Integer value (handles JSON float64)
resp.HasFloat("user.score", 9.5)       // Float value
resp.HasBool("user.active", true)      // Boolean value
resp.MatchesRegex("user.email", `[\w\.]+@`) // Regex match
resp.IsNull("user.deleted_at")         // Null check
resp.IsNotNull("user.id")              // Not null check
resp.ArrayLength("items", 3)           // Array length

// Bulk assertion
resp.Matches(map[string]any{
    "success": true,
    "data.count": 5,
    "data.user.name": "John",
})

```

### BDD Style with Describe/It

```go
func TestUserAPI(t *testing.T) {
    zt := zentests.New(t)
    
    zt.Describe("User API", func(ctx *zentests.SuiteContext) {
        var app *fiber.App
        
        ctx.BeforeEach(func(zt *zentests.T) {
            app = fiber.New()
            SetupRoutes(app)
        }).AfterEach(func(zt *zentests.T) {
            app.Shutdown()
        })
        
        ctx.It("lists users", func(zt *zentests.T) {
            zt.Get(app, "/users").
                OK().
                IsJSON().
                ArrayLength("data", 2)
        })
        
        ctx.It("creates user", func(zt *zentests.T) {
            zt.PostJSON(app, "/users", map[string]any{
                "name": "John",
            }).
                Created().
                Has("data.name", "John")
        })
    })
}

```

### Type Safety

`zentests` enforces strict type checking:

```go
// JSON numbers are float64 by default
resp.Has("count", 42)           // FAIL: 42 is int, JSON has float64
resp.HasInt("count", 42)        // PASS: handles conversion
resp.Has("count", float64(42))  // PASS: exact type match
resp.Has("count", 42.0)         // PASS: exact type match

```

<br>
<br>

## Development & Building

### Make Commands

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
make check-test-coverage    # Run tests and enforce coverage threshold (default: 90%)

# Release
make release                # Release a new tagged version (e.g. make release VERSION=v1.2.3)
make release VERSION=v1.2.3 # Create specific version

# Development
make versions               # Show current and all versions
make git-push               # Push code to both repositories

# Documentation
make changelog              # Generate changelog using git-cliff
make docs                   # Open documentation in browser
make docs-md                # Create Markdown API doc with `gomarkdoc` in `docs/API.md`

make clean                  # Clean build artifacts
```

<br>

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/[yourusername]/zentests.git`
3. Install dependencies: `go mod download`
4. Run tests: `make test`, `make test-verbose`
5. Check coverage: `make test-coverage`

### Guidelines

- **Code Quality**: Maintain the existing code style and documentation standards
- **Tests**: All new features must include comprehensive tests (aim for 100% coverage)
- **Documentation**: Update the `README.md` with any new features or changes
- **Commits**: Use clear, descriptive commit messages
- **Compatibility**: Ensure changes work with the latest stable Go and Fiber versions

### Reporting Issues

Please use GitHub Issues to report bugs or request features. When reporting bugs, include:
- Go version
- Fiber version
- Minimal code example that reproduces the issue
- Expected vs actual behavior

<br>

## Credits

**Authors**: 
- [Kematzy](https://github.com/kematzy)
- Kimi K2.5 model via [opencode Zen](https://opencode.ai/zen)

**Inspiration**:
- [RSpec](https://rspec.info/) - Ruby testing framework that inspired the BDD syntax
- [Pest](https://pestphp.com/) - PHP testing framework with elegant assertion syntax
- [testify](https://github.com/stretchr/testify) - Go assertion library used internally

**Built with**:
- [Fiber](https://gofiber.io/) - Express inspired web framework for Go
- [testify](https://github.com/stretchr/testify) - A toolkit with common assertions

Special thanks to the Go and Fiber communities for their excellent tools and documentation.

<br>

## License

This code is released under the MIT License.

Copyright (c) 2025 [Kematzy](https://github.com/kematzy)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

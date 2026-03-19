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
- [Fiber](https://github.com/gofiber/fiber) v3.x

> **Fiber v2 Users:** If you're using Fiber v2, use the `fiber-v2` branch:
> ```bash
> go get github.com/kematzy/zentests@fiber-v2
> ```
> See [Branches](#branches) for more information.

<br>

## Branches

| Branch | Fiber Version | Go Version | Status |
|--------|---------------|------------|--------|
| `master`   | v3.x | 1.25+ | Active development |
| `fiber-v2` | v2.x | 1.21+ | Maintenance only |

The `fiber-v2` branch receives bug fixes and security updates only. New features are developed on `master` (Fiber v3).

### Migrating from Fiber v2 to v3

If you're upgrading from Fiber v2 to v3, see the [upgrade guide](#migration-from-fiber-v2).

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
    // Note: Fiber v3 uses fiber.Ctx (interface), not *fiber.Ctx (pointer)
    app.Get("/api/health", func(c fiber.Ctx) error {
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
| `Delete(app, path)` | DELETE request  |
| `DeleteJSON(app, path, data)` | DELETE with JSON body  |

### HTTP Methods with TestConfig (Fiber v3+)

For custom timeout and behavior control, use the `*WithConfig` variants:

| Method | Description |
|--------|-------------|
| `GetWithConfig(app, path, cfg)` | GET with custom TestConfig |
| `PostWithConfig(app, path, body, cfg)` | POST with custom TestConfig |
| `PostJSONWithConfig(app, path, data, cfg)` | POST JSON with custom TestConfig |
| `PostFormWithConfig(app, path, data, cfg)` | POST form with custom TestConfig |
| `PutWithConfig(app, path, body, cfg)` | PUT with custom TestConfig |
| `PutJSONWithConfig(app, path, data, cfg)` | PUT JSON with custom TestConfig |
| `PatchWithConfig(app, path, body, cfg)` | PATCH with custom TestConfig |
| `PatchJSONWithConfig(app, path, data, cfg)` | PATCH JSON with custom TestConfig |
| `DeleteWithConfig(app, path, cfg)` | DELETE with custom TestConfig |
| `DeleteJSONWithConfig(app, path, data, cfg)` | DELETE JSON with custom TestConfig |

#### TestConfig Usage

```go
import (
    "time"
    "github.com/gofiber/fiber/v3"
    "github.com/kematzy/zentests"
)

func TestSlowEndpoint(t *testing.T) {
    zt := zentests.New(t)
    app := fiber.New()
    
    // Custom timeout for slow endpoints
    resp := zt.GetWithConfig(app, "/slow-endpoint", fiber.TestConfig{
        Timeout: 10 * time.Second,
    })
    resp.OK().IsJSON()
    
    // Custom timeout with FailOnTimeout behavior
    resp2 := zt.PostJSONWithConfig(app, "/api/users", userData, fiber.TestConfig{
        Timeout:       5 * time.Second,
        FailOnTimeout: true, // Returns error instead of partial response
    })
    resp2.Created()
}
```

#### TestConfig Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `Timeout` | `time.Duration` | `1s` | Request timeout duration |
| `FailOnTimeout` | `bool` | `true` | Return error on timeout (true) or partial response (false) |


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

## Using with testify/suite

`zentests` integrates seamlessly with `testify/suite` for more complex test scenarios requiring suite-level setup/teardown.

### Basic Structure

```go
package myapp_test

import (
    "testing"
    "github.com/gofiber/fiber/v3"
    "github.com/kematzy/zentests"
    "github.com/stretchr/testify/suite"
)

// Define your suite struct
type UserAPISuite struct {
    suite.Suite
    app *fiber.App
    zt  *zentests.T
}

// SetupSuite runs ONCE before all tests in the suite
func (s *UserAPISuite) SetupSuite() {
    s.app = fiber.New()
    s.zt = zentests.New(s.T())  // Initialize zentests.T

    // Setup routes
    s.app.Get("/users", func(c fiber.Ctx) error {
        return c.JSON(fiber.Map{"users": []string{"John", "Jane"}})
    })

    s.app.Post("/users", func(c fiber.Ctx) error {
        return c.Status(201).JSON(fiber.Map{"id": 1, "name": "John"})
    })
}

// TearDownSuite runs ONCE after all tests in the suite
func (s *UserAPISuite) TearDownSuite() {
    // Cleanup resources (close database, etc.)
}

// SetupTest runs BEFORE EACH test
func (s *UserAPISuite) SetupTest() {
    // Reset state before each test
}

// TearDownTest runs AFTER EACH test
func (s *UserAPISuite) TearDownTest() {
    // Cleanup after each test
}

// Tests use the zentests.T instance
func (s *UserAPISuite) TestListUsers() {
    s.zt.Get(s.app, "/users").
        OK().
        IsJSON().
        ArrayLength("users", 2)
}

func (s *UserAPISuite) TestCreateUser() {
    s.zt.PostJSON(s.app, "/users", map[string]any{
        "name": "John",
    }).
        Created().
        Has("id", float64(1))
}

// Run the suite
func TestUserAPISuite(t *testing.T) {
    suite.Run(t, new(UserAPISuite))
}
```

### Lifecycle Hooks Order

```
SetupSuite()           ← Runs once before all tests
    ├── TestA()
    │       ├── SetupTest()
    │       ├── [test code]
    │       └── TearDownTest()
    │
    └── TestB()
            ├── SetupTest()
            ├── [test code]
            └── TearDownTest()
TearDownSuite()        ← Runs once after all tests
```

### Database Integration

```go
type DatabaseSuite struct {
    suite.Suite
    app *fiber.App
    zt  *zentests.T
    db  *sql.DB  // or *gorm.DB
}

func (s *DatabaseSuite) SetupSuite() {
    // Connect to test database
    s.db, _ = sql.Open("sqlite3", ":memory:")

    s.app = fiber.New()
    s.zt = zentests.New(s.T())

    // Setup routes that use database
    s.app.Get("/users", func(c fiber.Ctx) error {
        var users []User
        s.db.Find(&users)
        return c.JSON(users)
    })
}

func (s *DatabaseSuite) SetupTest() {
    // Clean database before each test
    s.db.Exec("DELETE FROM users")
}

func (s *DatabaseSuite) TearDownSuite() {
    s.db.Close()
}
```

### Multiple Suites with Shared Setup

```go
// Base suite with common setup
type BaseSuite struct {
    suite.Suite
    App *fiber.App
    ZT  *zentests.T
}

func (s *BaseSuite) SetupSuite() {
    s.App = fiber.New()
    s.ZT = zentests.New(s.T())
}

// UserAPISuite inherits BaseSuite
type UserAPISuite struct {
    BaseSuite
}

func (s *UserAPISuite) SetupSuite() {
    s.BaseSuite.SetupSuite()
    // Add user-specific routes
    s.App.Get("/users", s.listUsers)
}

func TestUserAPISuite(t *testing.T) {
    suite.Run(t, new(UserAPISuite))
}
```

### Using TestConfig with Suite

```go
func (s *SlowAPISuite) TestSlowEndpoint() {
    // Use WithConfig for slow endpoints
    s.zt.GetWithConfig(s.app, "/slow-endpoint", fiber.TestConfig{
        Timeout: 10 * time.Second,
    }).OK()
}
```

### Choosing Between zentests.Describe and testify/suite

| Feature | `zentests.Describe` | `testify/suite` |
|---------|---------------------|-----------------|
| Setup | `BeforeEach` | `SetupTest` |
| Teardown | `AfterEach` | `TearDownTest` |
| Suite setup | Manual | `SetupSuite` |
| Suite teardown | Manual | `TearDownSuite` |
| Parallel tests | Limited | Supported |
| Subtests | `It()` | `s.Run()` |
| Assertions | `zentests.T` | `s.T()` + `s.zt` |
| Best for | Quick BDD tests | Complex test setups |

<br>

## Migration from Fiber v2

If you're upgrading from Fiber v2 to v3, here are the key changes affecting `zentests`:

### Import Path Change

```go
// Fiber v2
import "github.com/gofiber/fiber/v2"

// Fiber v3
import "github.com/gofiber/fiber/v3"
```

### Handler Signature Change

Fiber v3 uses `fiber.Ctx` as an interface (not a pointer):

```go
// Fiber v2
app.Get("/route", func(c *fiber.Ctx) error {
    return c.SendString("hello")
})

// Fiber v3
app.Get("/route", func(c fiber.Ctx) error {
    return c.SendString("hello")
})
```

### Body Parser Change

The `BodyParser` method has been replaced with the new `Bind()` API:

```go
// Fiber v2
app.Post("/users", func(c *fiber.Ctx) error {
    var data map[string]any
    c.BodyParser(&data)
    return c.JSON(data)
})

// Fiber v3
app.Post("/users", func(c fiber.Ctx) error {
    var data map[string]any
    if err := c.Bind().Body(&data); err != nil {
        return err
    }
    return c.JSON(data)
})
```

### New TestConfig Support

Fiber v3's `app.Test()` now accepts `TestConfig` instead of timeout duration:

```go
// Fiber v2
resp, err := app.Test(req, 5*time.Second)

// Fiber v3
resp, err := app.Test(req) // Uses default 1s timeout
resp, err := app.Test(req, fiber.TestConfig{
    Timeout: 5 * time.Second,
})
```

`zentests` provides `*WithConfig` methods to expose this functionality:

```go
// Fiber v2 style (still works)
zt.Get(app, "/route")

// Fiber v3 with custom timeout
zt.GetWithConfig(app, "/slow-route", fiber.TestConfig{
    Timeout: 10 * time.Second,
})
```

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
- GLM 5 model via [opencode Go](https://opencode.ai/go)

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


<!-- spellchecker:ignore kimi opencode -->

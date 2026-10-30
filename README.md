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
- **Well-tested**: 97.8% code coverage

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
> See [CONTRIBUTING.md](CONTRIBUTING.md) for branch policy, and the
> [Fiber v2 → v3 migration notes](docs/MIGRATION-FIBER-V2.md) when upgrading.

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

**Coverage: 97.8%** of statements

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
Describe blocks can be nested; hooks of enclosing blocks apply to inner tests. Declare hooks before the `It` calls they should apply to.

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

Full reference for every method, with examples: [docs/API.md](docs/API.md).

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
| `PostMultipart(app, path, fields, files)` | POST with multipart form data and files  |
| `GetQuery(app, path, url.Values)` | GET with encoded query parameters  |

### HTTP Methods with TestConfig (Fiber v3+)

Every request method has a `*WithConfig` twin (`GetWithConfig`, `PostJSONWithConfig`, ...)
that takes a `fiber.TestConfig`:
#### TestConfig Usage

```go
zt.GetWithConfig(app, "/slow-endpoint", fiber.TestConfig{
    Timeout: 10 * time.Second,
}).OK()
```

#### TestConfig Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `Timeout` | `time.Duration` | `1s` | Request timeout duration |
| `FailOnTimeout` | `bool` | `true` | Return error on timeout (true) or partial response (false) |

A zero-value `fiber.TestConfig{}` uses the defaults above. The plain methods (`Get`, `PostJSON`, ...) use it.

Set `ZENTESTS_TIMEOUT=5s` (any Go duration) to raise the default for a whole run, for
example `ZENTESTS_TIMEOUT=5s make test`. An invalid value panics at startup instead of
being silently ignored.


### Default Headers

Headers set on `T` (or on a `Jar`) are sent with every request. A header belonging to a
single request, such as `PostJSON`'s `Content-Type`, wins over a default with the same name.

```go
zt := zentests.New(t).Use(app).SetHeader("Accept", "application/json")
zt.Get("/users").OK()                  // sends Accept: application/json
zt.WithBearer(token)                   // sends Authorization: Bearer <token>

jar := zentests.NewJar(t, app)
jar.SetHeader("X-Api-Key", "key")      // same behavior on the jar
```


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
resp.StatusIn(200, 204) // One of several statuses
resp.IsSuccess()       // Any 2xx
resp.IsClientError()   // Any 4xx
resp.IsServerError()   // Any 5xx
resp.IsRedirect()      // 301, 302, 303, 307 or 308
resp.RedirectsTo("/dashboard") // Redirect + exact Location
```

### Header Assertions

```go
resp.HasHeader("X-Custom", "value")                                     // Exact match
resp.HeaderContains("X-Rate", "limit")                                  // Substring match
resp.HeaderPresent("X-Auth-Token")                                      // Header exists
resp.HeaderNotPresent("X-Auth-Token")                                   // Header does not exist
resp.HeaderHasValues("X-Features", []string{"audio", "video"})          // Header has multiple values
resp.CookieHasValues(map[string]string{"lang": "en", "theme": "dark"})  // Set-Cookie header test
resp.HasCookie("session")                                               // Response set the cookie
resp.CookieNotPresent("session")                                        // Response did not set it
resp.Cookie("session")                                                  // *http.Cookie, for attributes

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
resp.Has("user.name", "John")          // Exact match (any numeric type matches a JSON number)
resp.HasString("user.name", "John")    // String value
resp.HasInt("user.age", 30)            // Integer value (handles JSON float64)
resp.HasFloat("user.score", 9.5)       // Float value
resp.HasBool("user.active", true)      // Boolean value
resp.MatchesRegex("user.email", `[\w\.]+@`) // Regex match
resp.IsNull("user.deleted_at")         // Null check
resp.IsNotNull("user.id")              // Not null check
resp.JSONHasValues(map[string]any{     // Subset check: other keys are ignored
    "user.name": "John",
})
resp.ArrayLength("items", 3)           // Array length
resp.ArrayContains("roles", "admin")   // Array holds the element
resp.ArrayEmpty("roles")               // Array has no elements
resp.Bind(&user)                       // Decode body into a struct

```

### Cookie Jar

`NewJar` sends requests to one app and keeps the cookies it sets, the way a browser does.
Use it for flows that span several requests, such as login, CSRF and session tests. Every
request returns a `*Response`, so the usual assertions compose with jar requests.

```go
jar := zentests.NewJar(t, app)
csrf := jar.Get("/login").OK().BodyString()          // app sets session cookie here
jar.PostForm("/login", url.Values{"csrf": {csrf}, "user": {"jane"}}).Status(302)
jar.PostMultipart("/upload", nil, []zentests.FileUpload{
    {Name: "file", Filename: "a.txt", Content: []byte("hello")},
}).OK()
if jar.Cookie("session") == "" { t.Fatal("no session") }
jar.Get("/settings").OK().Contains("Welcome Jane")
```

- `jar.Get(path)`, `jar.PostForm(path, url.Values)` and `jar.Send(req)` for a caller-built request
- `jar.Cookie(name)` returns the stored value, or `""` when absent
- `resp.Cookie(name)` returns the `*http.Cookie` the response set, so you can check `HttpOnly`,
  `Path`, `Secure` and `SameSite`. It fails the test when the cookie is absent

A cookie is forgotten when a response sets it with `MaxAge < 0` or an empty value. Fiber's
`c.ClearCookie` sends an empty value. The jar is deliberately flat, with no `Path`, `Domain` or
`Expires` matching: every stored cookie is sent to every request on the app.

### BDD Style with Describe/It

```go
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
        zt.Get(app, "/users").OK().IsJSON().ArrayLength("data", 2)
    })
})
```

Nested `Describe` blocks inherit the hooks of enclosing blocks; declare hooks before
the `It` calls they should apply to. A working example lives in
[examples/suite.go](./examples/suite.go).

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

For complex test scenarios, `zentests` works with `testify/suite`: create the `*T` in
`SetupSuite` and use it in your tests.

```go
type UserAPISuite struct {
    suite.Suite
    app *fiber.App
    zt  *zentests.T
}

func (s *UserAPISuite) SetupSuite() {
    s.app = fiber.New()
    s.zt = zentests.New(s.T())
    // declare routes here
}

func TestUserAPISuite(t *testing.T) {
    suite.Run(t, new(UserAPISuite))
}
```

See [docs/SUITE.md](docs/SUITE.md) for lifecycle hooks order, database integration,
shared base suites and TestConfig usage.

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

## Contributing

Contributions are welcome! Development setup, guidelines, the branch policy and the full
`make` command reference live in [CONTRIBUTING.md](CONTRIBUTING.md).

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

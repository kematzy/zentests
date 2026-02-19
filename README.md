# zentests

Readable Go tests without magic. Inspired by RSpec/Pest

`zentests` provides fluent, chainable HTTP testing for [Fiber](https://gofiber.io/) applications with RSpec-style syntax.

## Why `zentests`?

- **No magic**: Standard Go testing, no reflection tricks
- **Type safe**: Strict assertions catch bugs early
- **Fiber-native**: Built for Fiber, not adapted from net/http
- **Chainable**: Readable, expressive test code
- **Familiar**: RSpec/Pest-inspired syntax


### Why only Fiber support?

Simple. I'm using Fiber for my projects and created these tests to make working on my projects faster.

You are welcome to create a version with support for your preferred routing package.


## Code Quality & Coverage

[Coverage: **98.2%** of statements](./.code-status/coverage.html)

## Installation

```bash
go get github.com/kematzy/zentests
```

## Quick Start

```go
package myapp_test

import (
    "testing"
    "github.com/gofiber/fiber/v2"
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

## HTTP Methods

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


## Status Assertions

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

## Header Assertions

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

## Body Assertions

```go
resp.Contains("substring")             // Body contains text
resp.NotContains("substring")          // Body excludes text
resp.BodyMatches(`regex\d+`)           // Regex match
resp.Equals("exact body")              // Exact match
resp.IsEmpty()                         // Body is empty
```

## JSON Assertions

Dot notation for nested keys: `user.profile.name`

```go
resp.JSON()                            // Parse & cache JSON
resp.HasKey("user.name")               // Key exists
resp.Has("user.name", "John")          // Strict type match
resp.HasString("user.name", "John")    // String value
resp.HasInt("user.age", 30)            // Integer value (handles JSON float64)
resp.HasFloat("user.score", 9.5)       // Float value
resp.HasBool("user.active", true)      // Boolean value
resp.MatchesRegex("user.email", `^[\w\.]+@`) // Regex match
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

## BDD Style with Describe/It

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

## Type Safety

`zentests` enforces strict type checking:

```go
// JSON numbers are float64 by default
resp.Has("count", 42)        // FAIL: 42 is int, JSON has float64
resp.HasInt("count", 42)     // PASS: handles conversion
resp.Has("count", 42.0)      // PASS: exact type match
```


## License

MIT

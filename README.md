# zentests

Readable Go tests without magic. Inspired by Pest and RSpec.

`zentests` provides chainable HTTP assertions and GORM database test helpers for
[Fiber](https://gofiber.io/) v3.x apps.

[![Go Reference](https://pkg.go.dev/badge/github.com/kematzy/zentests.svg)](https://pkg.go.dev/github.com/kematzy/zentests)
[![Go Report Card](https://goreportcard.com/badge/github.com/kematzy/zentests)](https://goreportcard.com/report/github.com/kematzy/zentests)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Status:** v0.3.x — API may still change.

## Why `zentests`?

- **No magic**: standard Go testing, no reflection tricks, no DSL or code generation
- **Chainable**: test assertions read easily — `zt.Get(app, "/api/v1/users").OK().IsJSON()`
- **Type safe**: strict assertions catch mistakes early instead of passing quietly
- **Fiber-native**: requests run through the real router and middleware — no listener, no port
- **GORM-focused**: built-in in-memory database setup, seeding and record helpers
- **Well-tested**: coverage gated at 95% in CI on every push

### How it compares

Code speaks louder than words. Below is the same test — status, content type and one JSON field —
for a Fiber app with a `GET /api/v1/healthz` route that returns `{"status": "it works"}`.

**Testing with `net/http` & `encoding/json`**

```go
req := httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
resp, err := app.Test(req)
if err != nil {
	t.Fatalf("request failed: %v", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
	t.Fatalf("status: got %d, want %d", resp.StatusCode, http.StatusOK)
}

if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
	t.Fatalf("content-type: got %q, want application/json", ct)
}

// NOTE! 13 lines abbreviated - see `docs/comparison` link below for full version.
```

**Testing with `testify` & `httptest`**

```go
req := httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
resp, err := app.Test(req)
require.NoError(t, err)
defer resp.Body.Close()

assert.Equal(t, http.StatusOK, resp.StatusCode)
assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

var payload map[string]any
require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))
assert.Equal(t, "it works", payload["status"])
```

**Testing with `zentests`**

```go
zt := zentests.New(t) // t *testing.T

zt.Get(app, "/api/v1/healthz").OK().IsJSON().Has("status", "it works")
```

Yes, just two efficient lines of code. **That's the essence of the `zentests` way!**

More: [docs/comparison](docs/comparison.md)

### Why only Fiber & GORM?

I created these test assertions and helpers for my projects where I use Fiber and GORM.

---

<details open>
  <summary><h2>Contents</h2></summary>

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [Features](#features)
    - [Request Helper Methods](#request-helper-methods)
    - [Status Assertions](#status-assertions)
    - [Headers Assertions](#headers-assertions)
    - [HTML Body Assertions](#html-body-assertions)
    - [JSON Body Assertions](#json-body-assertions)
    - [GORM Database Assertions](#gorm-database-assertions)
    - [Request Cookie Jars](#request-cookie-jars)
    - [Cookie Assertions](#cookie-assertions)
    - [App & Logger Helpers](#app-and-logger-helpers)
4. [Testing Styles](#testing-styles)
5. [Examples, Contributing & Credits](#examples)

</details>

---

## Installation

`zentests` is a test-only dependency: import it from `*_test.go` files and nothing is linked into
your compiled binary.

```bash
go get github.com/kematzy/zentests
```

**Dependencies**: Fiber, `testify`, GORM and their own transitive modules — nothing else.

While the API is still `v0.x`, it is suggested to pin the version explicitly:

```bash
go get github.com/kematzy/zentests@v0.3.0
```

### Requirements

- Go `1.25.5` or newer — the `go` directive in `go.mod`; older toolchains refuse the module.
- Fiber `v3.1` or newer — `github.com/gofiber/fiber/v3`.
- GORM `v1.31+`
    - a `C` compiler and `CGO_ENABLED=1` — only for the database helpers.
    - `gorm.io/driver/sqlite` links `mattn/go-sqlite3`, so the DB helpers need `gcc` (or
      `MSYS2/mingw` on Windows).
    - The HTTP assertions build fine with `CGO_ENABLED=0`.

More: [docs/database](docs/database/README.md)

### Compatibility

- Tested on Go `v1.25`, `v1.26` and `v1.27` (the CI matrix) against Fiber `v3.1+`; master tracks
  Fiber v3.
- Every push runs an equivalent gate in CI: `gofmt`, `go vet`, `goimports`, `golangci-lint`, the
  suite with `-race`, and a **95% coverage gate**. `make check` runs the same locally, plus the
  spell checker.
- Releases are tagged and summarized in `CHANGELOG.md`.

<br />

---

## Quick Start

Getting started with `zentests` is a matter of deciding your overall preferred testing style,
[flat](examples/example-01.md), [BDD](examples/example-02.md) or
[`testify/suite`](examples/example-03.md).

Either way, you need to import `zentests` with `import "github.com/kematzy/zentests"` and then
instantiate it:

### Flat style

```go
import "github.com/kematzy/zentests"

func Test_APIv1_Route_Get_Healthz(t *testing.T) {
	// declared in each test method
	app := fiber.New()            // load your actual Fiber app
	zt := zentests.New(t)

	// perform the tests
	zt.Get(app, "/api/v1/healthz").OK()
}
```

Working code: [examples/01-flat-style](examples/example-01.md) <br />More:
[docs/testing/flat](docs/testing/flat.md)

### `testify/suite` style

```go
// create suite struct
type HealthSuite struct {
	suite.Suite
	app *fiber.App
	zt  *zentests.T
}

// SetupSuite - runs once before the suite.
func (s *HealthSuite) SetupSuite() {
	s.app = fiber.New()           // load your actual Fiber app
	s.zt = zentests.New(s.T())
}

// TestAPIv1RouteGetHealthz - tests the route.
func (s *HealthSuite) TestAPIv1RouteGetHealthz() {
	s.zt.Get(s.app, "/api/v1/healthz").OK()
}

// TestHealthSuite - runs the test suite.
func TestHealthSuite(t *testing.T) { suite.Run(t, new(HealthSuite)) }
```

Working code: [examples/03-testify-suite](examples/example-03.md) <br />More:
[docs/testing/testify-suite](docs/testing/testify-suite.md)

<br />
<br />

---

## Features

`zentests` provides a wide variety of useful test assertions, application and database helpers to
test most common scenarios.

Each area below provides a brief overview with a link to further detailed information.

### Request Helper Methods

Checking your app's routes should be obvious and as simple as `Get(app, "/")`, or
`PostJSON(app, "/login", data)`, or as flexible as `PutJSONWithConfig(app, path, data, cfg)` when a
route requires special configs.

One call and one chain — no request builder, no `httptest`, no body to read or close:

```go
// testing that a GET route returns status 200 and JSON
zt.Get(app, "/api/v1/users").OK().IsJSON()

// testing that a POST route successfully returns an expected token
data := map[string]any{"email": "joe@example.com", "password": "secret"}
zt.PostJSON(app, "/api/v1/login", data).OK().HasKey("token")

// testing that a POST route returns status 200 after submitting a form
zt.PostForm(app, "/api/v1/contact", map[string]string{"subject": "Hi"}).OK()

// testing that a GET route with a query string returns status 200
zt.GetQuery(app, "/api/v1/search", url.Values{"q": {"fiber"}}).OK()
```

`GetQuery` and `PostMultipart` cover query strings and uploads, and every verb method has a
`*WithConfig` twin.

More: [docs/requests](docs/requests/README.md)

---

### Status Assertions

Use simple status names, like `OK()`, `NotFound()` instead of `Status(404)`, or just
`IsServerError()` when the range is more important. `zentests` provides assertions for all common
status codes.

```go
// named statuses read as assertions
zt.Get(app, "/").OK()                             // == Status(200)
zt.Get(app, "/does/not/exist").NotFound()         // == Status(404)
zt.Get(app, "/admin/dashboard").Forbidden()       // == Status(403)
zt.PostJSON(app, "/api/v1/users", data).Created() // == Status(201)

// classes, when the exact code is not the point
zt.Get(app, "/api/v1/search").IsSuccess()         // 2xx range
zt.Get(app, "/api/v1/report").IsServerError()     // 5xx range
zt.Get(app, "/api/v1/users/1").StatusIn(200, 204)
```

More: [docs/assertions/status](docs/assertions/status.md)

---

### Headers Assertions

Save time and improve clarity with assertions like `IsJSON()` instead of
`HeaderContains("Content-Type", "application/json")`, or easily handle complex headers with
`HeaderHasValues("X-Custom", values)`.

```go
// check for the exact header value
zt.Get(app, "/").HasHeader("X-Served-By", "api-1")
// check for presence
zt.Get(app, "/api/v1/users").HeaderPresent("X-Request-Id")
// ensure absent: a one-line security check
zt.Get(app, "/api/v1/users").HeaderNotPresent("X-Powered-By")
// check header for multiple values
zt.Get(app, "/api/v1/users").HeaderHasValues("X-Features", []string{"audio", "video"})
```

Content-type shortcuts: [docs/assertions/content-types](docs/assertions/content-types.md)
<br />More: [docs/assertions/headers](docs/assertions/headers.md)

---

### HTML Body Assertions

Easily assert returned HTML `Contains("some words")` or `NotContains("other words")`, you can even
use regexp with `BodyMatches(regexp)` to assert that content is present.

You can even process the response yourself with the `Body()` or `BodyString()` methods.

```go
zt.Get(app, "/").OK().Contains("Welcome back")
zt.Get(app, "/checkout").OK().NotContains("Your cart is empty")
// use regexp, when the contents might change
zt.Get(app, "/contact").OK().BodyMatches(`Email: \S+@\S+`)
zt.Get(app, "/robots.txt").OK().Equals("User-agent: *")
zt.Get(app, "/blank").OK().IsEmpty()
```

More: [docs/assertions/body-html](docs/assertions/body-html.md)

---

### JSON Body Assertions

Asserting JSON APIs made easy with `HasKey("user.name")`, or find exact matches with
`Has("user.name", "John")` or match just a subset of values with
`JSONHasValues(map[string]any{"user.name": "John"})`.

Dot-notation path support makes testing deep keys easy without structs.

```go
zt.Get(app, "/api/v1/users/1").JSON(). // parse once, cached on the response
	HasKey("user.name").                 // key exists
	HasString("user.name", "John").      // and holds a string
	HasInt("user.age", 42).              // typed: no float64 surprise
	ArrayContains("user.roles", "admin")

// nested objects and arrays, addressed by path
zt.Get(app, "/api/v1/users").Has("users.0.address.city", "Berlin")

// subset check: only the keys you name are compared
zt.Get(app, "/api/v1/users/1").JSONHasValues(map[string]any{
	"user.name":   "John",
	"user.active": true,
})

// or bypass paths entirely and decode into your own type
var user User
zt.Get(app, "/api/v1/users/1").OK().Bind(&user)
```

More: [docs/assertions/json](docs/assertions/body-json.md)

---

### GORM Database Assertions

Conveniently `SetupTestDB()`, `DBMigrate()`, `DBSeed()` and `DBReset()` test databases and assert
you can create `DBCreate()`, read `DBFind()`, update `DBUpdate()` or delete `DBDelete()` records.

A real test database, gone when the test ends: one line opens in-memory SQLite, migrates your models
and closes itself. No fixture files, no cleanup, nothing leaking into the next test.

```go
db := zentests.SetupTestDBWithModels(t, &User{}, &Post{}) // :memory:, AutoMigrate, auto-close

user := zentests.DBCreate(t, db, &User{Name: "Jane", Email: "jane@example.com"})
zentests.DBExists(t, db, &User{}, "email = ?", "jane@example.com")

found := zentests.DBFind(t, db, &User{}, user.ID)
zentests.DBUpdate(t, db, found, found.ID, map[string]any{"name": "Janet"})
zentests.DBCount(t, db, &User{}) // 1
```

Note, uses GORM's SQLite driver so it requires CGO and a C compiler.

More: [docs/database](docs/database/README.md)

---

### Request Cookie Jars

Easily test session flows with `NewJar(t, app)`: cookies set by one response are sent with the next,
and forgotten the moment the app deletes them.

```go
jar := zentests.NewJar(t, app)

// log in once: the session cookie is stored
jar.PostForm("/login", url.Values{"email": {"jane@example.com"}, "password": {"secret"}}).
	OK().
	HasCookie("session")

// later requests are authenticated without ever mentioning the cookie
jar.Get("/api/v1/me").OK().HasString("user.name", "Jane")

// log out, and the same request is refused
jar.Get("/logout").OK()
jar.Get("/api/v1/me").Unauthorized()
```

The jar is deliberately flat: no path, domain or expiry matching, so every stored cookie is sent to
every request. `jar.Cookie("session")` returns the raw value when you need it.

More: [docs/requests/cookie-jars](docs/requests/cookie-jars.md)

---

### Cookie Assertions

Assert that your cookies are present with `HasCookie("session")`, or removed with
`CookieNotPresent("old_session")` and even if they contain multiple values with
`CookieHasValues(values)`.

```go
// one response, several cookie checks
zt.PostJSON(app, "/api/v1/login", data).OK().
	HasCookie("session").
	CookieNotPresent("old_session")

zt.Get(app, "/settings").OK().CookieHasValues(map[string]string{
	"lang":  "en",
	"theme": "dark",
})

// test the cookie attributes
session := zt.PostJSON(app, "/api/v1/login", data).OK().Cookie("session")
assert.True(t, session.HttpOnly)
assert.True(t, session.Secure)
```

Cookies that clear themselves — empty value or `MaxAge < 0` — count as not set for `HasCookie`,
which keeps logout tests honest.

More: [docs/assertions/cookies](docs/assertions/cookies.md)

---

### App and Logger Helpers

Create custom config test apps easily with `NewApp()`, confirm routes are registered with
`HasRoute()`, or even test your `DefaultLogger` on fatal paths with the built-in `MockLogger`
helper.

```go
// declare a custom Fiber app with a custom configs
app := zentests.NewApp(t, fiber.Config{
	Views:             html.New("./views", ".html"), // requires `github.com/gofiber/template/html/v3`
	AppName:           "MyApp",
	PassLocalsToViews: true,
})

// checking route registration without a request
zt.HasRoute("GET", "/api/v1/users", app).NotHasRoute("DELETE", "/api/v1/users", app)

ml := &zentests.MockLogger{} // inject, then assert what was logged
svc := NewService(ml)
svc.Start()
assert.True(ml.WarnCalled) // true
assert.Equal("config not found", ml.WarnMsg) // "config not found"
```

More: [docs/helpers](docs/helpers/README.md)

<br />

---

### And more

The above are just some of the many features `zentests` provides. You can find all the information
in the [docs](docs/README.md) section.

Don't miss:

- [Default headers](docs/requests/headers.md) &mdash; `SetHeader`, `WithBearer`, precedence
- [TestConfig & Timeouts](docs/requests/testconfigs.md) &mdash; `*WithConfig`, `ZENTESTS_TIMEOUT`,
  `FailOnTimeout`
- [Debugging](docs/assertions/debugging.md) &mdash; `Dump`, `Debug`, header redaction
- [Multipart](docs/requests/multipart.md) &mdash; `PostMultipart`, `FileUpload`

<br />

---

## Testing Styles

`zentests` works with three styles of Go test.

- [Flat](docs/testing/flat.md) &mdash; best for independent tests, small files, or parallel
  `t.Parallel()` tests.
- [BDD](docs/testing/bdd.md) &mdash; best for tests grouped by behaviour and readable sub-test
  trees.
- [testify/suite](docs/testing/testify-suite.md) &mdash; best for tests with one app and/or database
  shared by many tests.

Choose your preference. Try the flat style first, then the BDD style, or go straight to the
`testify/suite` style. It really doesn't matter, because the usage is almost identical in all of
them, just instantiate `zt := zentests.New(t)` and begin testing.

More: [docs/testing](docs/testing/README.md)

<br />

---

## Examples

A collection of runnable code examples, with at least one per testing style. Each example contains
the instructions and complete `*_test.go` code so you can test them yourself.

**Available Examples:**

- [Example 1](examples/example-01.md) &mdash; plain `func Test` style tests of status, headers, JSON
  paths and array indices, `PostJSON`, and a 404 response

- [Example 2](examples/example-02.md) &mdash; BDD styled tests with nested blocks, `BeforeEach`
  scope and dotted paths.

- [Example 3](examples/example-03.md) &mdash; uses `testify/suite` styled tests with a shared
  in-memory database, `DBReset` per test and seeding with `DBCreateN`.

These examples are intended to be runnable, while the code in the [docs](docs/README.md) pages is
generally illustrative.

---

## Contributing

**Contributions are welcome!**

Development setup, guidelines and more available in [docs/contributing](docs/contributing.md).

<br />

## Credits

- [Kematzy](https://github.com/kematzy)
- Various AI agents and models

## License

Copyright 2026 Kematzy &mdash; Released under the [MIT License](LICENSE).

# Zentests - Examples - Example 01 - Flat Style

Flat style keeps every test a plain `func TestX(t *testing.T)` — no suite, no shared state, one
helper created at the top of each test. It is the smallest thing that works and the style most Go
developers start with.

## Run it

```bash
mkdir zentests-flat && cd zentests-flat
go mod init example.com/zentests-flat
go get github.com/gofiber/fiber/v3 github.com/kematzy/zentests
# save the code below as api_test.go
go test -v ./...
```

Rename the package to match your project; nothing else in the file is project-specific.

## The code

```go
package api_test

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/kematzy/zentests"
)

// newApp is the app under test: four routes, no database.
func newApp() *fiber.App {
	app := fiber.New()

	app.Get("/api/v1/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "it works"})
	})

	app.Get("/api/v1/users", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"users": []fiber.Map{
			{"id": 1, "name": "Jane", "roles": []string{"admin", "editor"}},
			{"id": 2, "name": "Joe", "roles": []string{"viewer"}},
		}})
	})

	app.Post("/api/v1/users", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": 3, "name": "Sam"})
	})

	app.Get("/api/v1/me", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"authorization": c.Get("Authorization")})
	})

	return app
}

func TestHealthz(t *testing.T) {
	app := newApp()
	zt := zentests.New(t)

	zt.Get(app, "/api/v1/healthz").
		OK().
		IsJSON().
		Has("status", "it works")
}

func TestListUsers(t *testing.T) {
	app := newApp()
	zt := zentests.New(t)

	zt.Get(app, "/api/v1/users").
		OK().
		ArrayLength("users", 2).
		HasString("users.0.name", "Jane").
		HasInt("users.1.id", 2).
		ArrayContains("users.0.roles", "admin")
}

func TestCreateUser(t *testing.T) {
	app := newApp()
	zt := zentests.New(t)

	zt.PostJSON(app, "/api/v1/users", fiber.Map{"name": "Sam"}).
		Created().
		IsJSON().
		HasInt("id", 3).
		HasString("name", "Sam")
}

func TestDefaultHeaderIsSent(t *testing.T) {
	app := newApp()
	zt := zentests.New(t).SetHeader("Authorization", "Bearer secret-token")

	zt.Get(app, "/api/v1/me").
		OK().
		HasString("authorization", "Bearer secret-token")
}

func TestMissingRoute(t *testing.T) {
	app := newApp()
	zt := zentests.New(t)

	zt.Get(app, "/api/v1/nope").NotFound()
}
```

## What it shows

- `zentests.New(t)` once per test — no shared helper, no ordering between tests.
- A request and its assertions in one chain: status, content type and body.
- Dot-notation JSON paths, including array indices (`users.0.name`).
- Typed JSON assertions (`HasString`, `HasInt`) next to the generic `Has`.
- A default header set on `zt` and sent with every request.
- A negative assertion (`NotFound`) with no setup beyond the app.

## Related

- Full reference: [docs/assertions](../docs/assertions/README.md)
- Request methods: [docs/requests/methods](../docs/requests/methods.md)
- Other examples: [Example 02](example-02.md), [Example 03](example-03.md)

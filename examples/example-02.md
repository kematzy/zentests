# Zentests - Examples - Example 02 - BDD Style

BDD style groups related tests under `Describe` blocks, with `BeforeEach`/`AfterEach` hooks and one
`It` per behaviour. It reads like a specification and needs no extra dependency beyond `zentests`.

## Run it

```bash
mkdir zentests-bdd && cd zentests-bdd
go mod init example.com/zentests-bdd
go get github.com/gofiber/fiber/v3 github.com/kematzy/zentests
# save the code below as api_test.go
go test -v ./...
```

Rename the package to match your project. Each `It` becomes a subtest, so `-run` works on names:
`go test -run 'TestUserAPI/User_API/when_listing_users' -v`.

## The code

```go
package api_test

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/kematzy/zentests"
)

func newApp() *fiber.App {
	app := fiber.New()

	app.Get("/api/v1/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "it works"})
	})

	app.Get("/api/v1/users", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"users": []fiber.Map{
			{"id": 1, "name": "Jane"},
			{"id": 2, "name": "Joe"},
		}})
	})

	return app
}

func TestUserAPI(t *testing.T) {
	zt := zentests.New(t)
	var app *fiber.App

	zt.Describe("User API", func(ctx *zentests.SuiteContext) {
		ctx.BeforeEach(func(zt *zentests.T) {
			app = newApp() // fresh app before every It
		})

		ctx.It("reports its health", func(zt *zentests.T) {
			zt.Get(app, "/api/v1/healthz").
				OK().
				Has("status", "it works")
		})

		ctx.Describe("when listing users", func(ctx *zentests.SuiteContext) {
			ctx.BeforeEach(func(zt *zentests.T) {
				// Applies to this block and its nested blocks only.
				zt.SetHeader("Accept", "application/json")
			})

			ctx.It("returns two users", func(zt *zentests.T) {
				zt.Get(app, "/api/v1/users").
					OK().
					IsJSON().
					ArrayLength("users", 2)
			})

			ctx.It("supports dotted paths", func(zt *zentests.T) {
				zt.Get(app, "/api/v1/users").
					OK().
					HasString("users.1.name", "Joe")
			})
		})
	})
}
```

## What it shows

- `Describe`/`It` produce nested subtests, so the test output is a readable tree.
- `BeforeEach` on the outer block rebuilds the app for every `It`; the nested `BeforeEach` adds a
  header for its block only.
- Hooks declared after an `It` do not apply to it — declare hooks first.
- Each `It` receives its own `*zentests.T`, so headers and state never leak between tests.

## Related

- BDD guide: [docs/testing/bdd](../docs/testing/bdd.md)
- Other examples: [Example 01](example-01.md), [Example 03](example-03.md)

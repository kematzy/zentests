# Zentests - Docs - Comparison

The same test, three ways: a health check that must return `200`, `application/json` and
`{"status": "it works"}`.

The README shows abbreviated versions of all three. Below is the complete code. It also shows what
`zentests` removes, what it adds, where its scope ends, and how to migrate an existing test.

## Contents

- [The same test, three ways](#the-same-test-three-ways)
- [What `zentests` removes](#what-zentests-removes)
- [What else it gives you](#what-else-it-gives-you)
- [What `zentests` does not do](#what-zentests-does-not-do)
- [Migrating an existing test](#migrating-an-existing-test)
- [Which should you use](#which-should-you-use)

---

## The same test, three ways

First the Fiber app we are testing, then the same test written three ways.

```go
package api

import "github.com/gofiber/fiber/v3"

func NewAPIApp() *fiber.App {
	app := fiber.New()

	app.Get("/api/v1/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "it works"})
	})

	return app
}
```

The test snippets below assume this app plus the required imports.

### `net/http` and `encoding/json`

```go
func TestRouteAPIv1Healthz_Raw(t *testing.T) {
	app := api.NewAPIApp()  // create app instance

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	var payload struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if payload.Status != "it works" {
		t.Fatalf("status field: got %q, want %q", payload.Status, "it works")
	}
}
```

Summary: **29 lines** of test code.

The test handles request construction, errors, status and header checks, body reading, and JSON
decoding explicitly.

<br />

### `testify` with `httptest`

```go
func TestRouteAPIv1Healthz_Testify(t *testing.T) {
	app := api.NewAPIApp() // create app instance

	req := httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	var payload map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))
	assert.Equal(t, "it works", payload["status"])
}
```

Summary: **13 lines** of test code.

`testify` shortens comparisons, but request construction, error handling, and body decoding remain
explicit. Each additional expectation requires another assertion call.

<br />

### `zentests`

```go
func TestRouteAPIv1Healthz_Zentests(t *testing.T) {
	app := api.NewAPIApp() // create app instance
	zt := zentests.New(t)  // create zentests instance

	zt.Get(app, "/api/v1/healthz").OK().IsJSON().Has("status", "it works")
}
```

Summary: **4 lines** oft test code.

The request and its three expectations fit in one chain. `zentests` is built on `testify`, so it
composes with the assertions you already use rather than replacing them.

---

## What `zentests` removes

The same three assertions from the test above, without the supporting code:

- **Request boilerplate** — `httptest.NewRequest` + `app.Test` + error check collapse into one call:
  `zt.Get(app, path)`, `zt.PostJSON(app, path, data)`.
- **Body lifecycle** — the response body is read once and kept on the `*Response`, so several
  assertions can run against the same body without `io.ReadAll` or `Close`.
- **Decode boilerplate** — `Has("status", "it works")` replaces unmarshal + lookup + compare, and
  reports the key path, expected value and actual value.

Each assertion is a method on the same `*Response`, so adding a fourth check is one more link in the
chain, not another block.

## What else it gives you

Beyond the single request above, the same import also provides:

- **Timeouts** — plain methods use a `1s` timeout by default; `ZENTESTS_TIMEOUT` overrides it for
  the process, and `fiber.TestConfig` customizes a call. A timeout fails the test by default.
- **Multi-step flows** — `NewJar` keeps cookies between requests, so login → authenticated request
  reads like the flow it tests.
- **Database setup** — `SetupTestDBWithModels`, `DBSeed`, `DBReset` replace hand-rolled in-memory
  SQLite wiring.

## What `zentests` does not do

Every library has an edge, and knowing it up front saves an evaluation round-trip:

- **Fiber only.** It does not test `net/http` handlers, `chi`, `echo`, `gin` or the standard mux.
- **GORM only** for the database helpers. No `database/sql`, `sqlx`, `ent` or `sqlc` helpers.
- **Not a replacement for `testify`.** It is built on `testify/assert` and `testify/require`; you
  can and often will use both in one file.

General-purpose mocks, container orchestration, and snapshot testing are **currently** not built in.
Use those tools alongside `zentests`; `MockLogger` is a focused logger test double.

---

## Migrating an existing test

For each test, create a helper and replace request execution plus response handling with one chain.
In a `testify/suite`, create the helper in `SetupTest`, not `SetupSuite`, so it is bound to the
current test's `*testing.T`.

1. **Create the helper once** in `SetupTest`/`SetupSuite` or at the top of the test:

    ```go
    zt := zentests.New(t)
    ```

2. **Replace the request block** with the equivalent method call:

    ```go
    // before
    req := httptest.NewRequest(http.MethodGet, "/users", nil)
    resp, err := app.Test(req)
    require.NoError(t, err)

    zt.Get(app, "/users")	// after
    ```

3. **Replace the assertion block** with a chain, and drop the `defer resp.Body.Close()` — the body
   is read for you:

    ```go
    // before
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode)
    assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
    var got map[string]any
    require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
    assert.Equal(t, "Jane", got["name"])

    zt.Get(app, "/users").OK().IsJSON().Has("name", "Jane") // after
    ```

The response body is read and cached when an assertion needs it; no manual decode or close is
needed.

Keep specialized checks such as `assert.Eventually` and direct struct comparisons as `testify`
assertions. Mixing them with `zentests` is normal.

## Which should you use

| Situation                                           | Recommendation                                                                                |
| --------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| New Fiber v3 app, tests written from scratch        | `zentests` for HTTP and database tests                                                        |
| Existing `testify` + `httptest` suite, working well | Keep it; migrate route by route if the boilerplate hurts                                      |
| Existing `testify/suite` structure                  | Create `zentests.New(s.T())` in `SetupTest`; see [testify-suite.md](testing/testify-suite.md) |
| Non-Fiber framework                                 | Not supported                                                                                 |
| Need mocking, containers, snapshots                 | Use those tools alongside `zentests`, not instead of it                                       |

---

Back to: [README](../README.md) &mdash; Next: [Assertions](assertions/README.md)

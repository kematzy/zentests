# Zentests - Docs - Testing - Flat

Flat style keeps each test a plain `func TestX(t *testing.T)` and creates `zentests.New(t)` inside
it. It is a straightforward fit for independent tests; use `t.Parallel()` only when the app and
other fixtures are safe to share or are created separately per test.

Use it as the default. Choose [BDD](./bdd.md) or [testify/suite](./testify-suite.md) when their
organization better fits the test cases.

## The shape

```go
func TestHealthz(t *testing.T) {
	app := fiber.New()
	app.Get("/api/v1/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "it works"})
	})

	zt := zentests.New(t)
	zt.Get(app, "/api/v1/healthz").OK().IsJSON().Has("status", "it works")
}
```

Build the app and helper in the test, then send a request and assert on its response. Keeping the
setup local makes the test easier to run independently; shared external state still needs its own
isolation.

## Common patterns

### One test, several assertions

Chain them to check several parts of one response. A failed assertion is recorded, and the chain
continues so the test can report other failures too:

```go
zt.Get(app, "/api/v1/users").
	OK().
	IsJSON().
	ArrayLength("users", 2).
	HasString("users.0.name", "Jane")
```

### Default headers

`SetHeader` and `WithBearer` apply to every request the helper sends:

```go
zt := zentests.New(t).WithBearer("secret-token")
zt.Get(app, "/api/v1/me").OK()
```

### Slow endpoints

Wrap the request in `fiber.TestConfig` when the handler legitimately takes longer than the default
timeout (`ZENTESTS_TIMEOUT` raises it for the whole run instead):

```go
zt.GetWithConfig(app, "/slow", fiber.TestConfig{
	Timeout:       10 * time.Second,
	FailOnTimeout: true,
}).OK()
```

### Parallel tests

Each test creates its own app and helper here, so this pattern can use `t.Parallel()`. Avoid sharing
mutable fixtures between parallel tests:

```go
func TestHealthz(t *testing.T) {
	t.Parallel()
	// ...
}
```

## Gotchas

- **Top-level JSON arrays.** JSON assertions address paths from an object root. A response of
  `[{"id":1}]` fails to parse; return `{"users": [...]}` and assert `users.0.id`, or decode with
  `Bind(&slice)`.
- **JSON numbers decode as `float64`.** `Has` normalizes numeric expectations, so `Has("count", 42)`
  matches a JSON number with value 42. Use `HasInt` or `HasFloat` to make your numeric expectation
  explicit. See [assertions](../assertions/README.md).
- **Empty path is not the root.** `ArrayLength("", 2)` does not address the whole body — name the
  key (`users`).
- **Create `zentests.New(t)` inside each test.** Reusing a helper ties assertions to the wrong test
  and can share helper-level state such as default headers.

## Working code

[examples/01-flat-style](../../examples/example-01.md) — the app plus five tests, runnable as-is.

## Related

- [Testing styles overview](./README.md)
- [Assertions](../assertions/README.md)
- [Request methods](../requests/methods.md)

---

Back to: [Testing styles](./README.md) &mdash; Next: [BDD](./bdd.md)

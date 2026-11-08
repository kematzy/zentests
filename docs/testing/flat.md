# Zentests - Docs - Testing - Flat

> **Status:** DRAFT — structure and wording are open for discussion.

Flat style keeps every test a plain `func TestX(t *testing.T)`. The helper is created inside the
test, so tests share no state, run in any order, and work with `t.Parallel()`.

Use it for most tests. Reach for [BDD](./bdd.md) or [testify/suite](./testify-suite.md) only when
their structure earns its keep.

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

Three things happen in order: build the app, create the helper, send a request and assert on the
response. Nothing is shared with the next test, so a failure is always self-contained.

## Common patterns

### One test, several assertions

Chain them — a failure stops at the first broken link and names it:

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
zt.GetWithConfig(app, "/slow", fiber.TestConfig{Timeout: 10 * time.Second}).OK()
```

### Parallel tests

Nothing is global, so `t.Parallel()` is safe:

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
- **Numbers are `float64`.** `Has("count", 42)` fails when the JSON holds `42`; use `HasInt`,
  `HasFloat` or pass `float64(42)`. See [assertions](../assertions/README.md).
- **Empty path is not the root.** `ArrayLength("", 2)` does not address the whole body — name the
  key (`users`).
- **`zentests.New(t)` per test.** Creating it once in `TestMain` or a package variable reintroduces
  shared state and breaks parallelism.

## Working code

[examples/01-flat-style](../../examples/example-01.md) — the app plus five tests, runnable as-is.

## Related

- [Testing styles overview](./README.md)
- [Assertions](../assertions/README.md)
- [Request methods](../requests/methods.md)

---

Back to: [Testing styles](./README.md) &mdash; Next: [BDD](./bdd.md)

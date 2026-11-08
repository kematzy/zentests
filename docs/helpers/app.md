# Zentests - Docs - Helpers - App

> **Status:** DRAFT — structure and wording are open for discussion.

`NewApp` builds a fresh Fiber app and registers its shutdown with `t.Cleanup`, so no test has to end
with a `defer app.Shutdown()`.

## Signature

```go
zentests.NewApp(t, cfg ...fiber.Config) *fiber.App
```

The configuration is optional and variadic, so the zero-config form stays short. Only the first
value is used if several are passed.

## Without configuration

```go
app := zentests.NewApp(t)
app.Get("/api/v1/healthz", healthz)

zt := zentests.New(t)
zt.Get(app, "/api/v1/healthz").OK()
```

Every call returns a new app — nothing is shared between tests, so two tests cannot see each other's
routes.

## With configuration

Any `fiber.Config` can be passed: views engine, app name, error handler, `PassLocalsToViews`, a
custom `ReadTimeout`.

```go
app := zentests.NewApp(t, fiber.Config{
	Views:             html.New("./views", ".html"),
	AppName:           "MyApp",
	PassLocalsToViews: true,
})
```

This is the case `NewApp` exists for: a handler that renders a template or reads `AppName` needs the
same configuration as production, and rebuilding that by hand in every test file is the kind of
setup that drifts.

## Why shutdown is registered

`app.Test` runs requests in-process and never opens a listener, so a test-only app usually has
nothing to release. `NewApp` still registers `Shutdown` because:

- a suite-level app lives for many tests, and leaving it running hides resources it opened;
- configuration may start background work (a views reloader, a prefork listener);
- a leak is silent until it isn't, and one line at creation is cheaper than a `TearDownSuite`.

## Notes

- **One app per test, or one per suite.** Both work; `NewApp` is about the shutdown, not about
  scope.
- **Cleanup runs in reverse order**, so an app created before a database is shut down after it.
- **`NewApp` is not required.** `fiber.New()` works exactly the same with `zentests`; only the
  automatic shutdown is missing.

## Related

- [Routes](../assertions/routes.md) — asserting what was registered on the app
- [testing/testify-suite](../testing/testify-suite.md) — suite-level app setup
- [examples/03-testify-suite](../../examples/example-03.md)

---

Back to: [Helpers](./README.md) &mdash; Next: [Logger](./logger.md)

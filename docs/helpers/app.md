# Zentests - Docs - Helpers - App

`NewApp` builds a fresh Fiber app and registers its shutdown with `t.Cleanup`, so no test has to end
with a `defer app.Shutdown()`.

## Signature

```go
zentests.NewApp(t, cfg ...fiber.Config) *fiber.App
```

| Input  | Type              | Behavior                                           |
| ------ | ----------------- | -------------------------------------------------- |
| `t`    | `*testing.T`      | Registers `app.Shutdown()` with `t.Cleanup`        |
| `cfg`  | `...fiber.Config` | Optional; only the first config is used            |

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

`app.Test` handles requests in-process without opening a listener. `NewApp` still registers
`Shutdown` so Fiber's app-level shutdown lifecycle runs at test cleanup, including any configured
shutdown hooks. It does not start a listening server itself.

## Notes

- **One app per test, or one per suite.** Both work; `NewApp` is about the shutdown, not about
  scope.
- **Cleanup runs in reverse order.** If an app's shutdown hook needs the database, create the DB
  before `NewApp` so the app shuts down first.
- **`NewApp` is not required.** `fiber.New()` works exactly the same with `zentests`; only the
  automatic shutdown is missing.

## Related

- [Routes](../assertions/routes.md) — asserting what was registered on the app
- [testing/testify-suite](../testing/testify-suite.md) — suite-level app setup

---

Back to: [Helpers](./README.md) &mdash; Next: [Logger](./logger.md)

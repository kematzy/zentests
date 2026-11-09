# Zentests - Docs - Assertions - Routes

Route assertions check what the app has _registered_, without sending a request. They catch the case
that produces the most confusing test failures: a route that was never mounted, so the assertion was
really testing a 404.

## Assertions

| Assertion                        | Checks                      |
| -------------------------------- | --------------------------- |
| `HasRoute(method, path, app)`    | The route is registered     |
| `NotHasRoute(method, path, app)` | The route is not registered |

Usage:

```go
zt.HasRoute("GET", "/api/v1/users", app).
	HasRoute("POST", "/api/v1/users", app).
	HasRoute("DELETE", "/api/v1/users/:id", app).
	NotHasRoute("DELETE", "/api/v1/users", app)
```

Both return `*T`, so they chain with each other and with further calls.

## Method and path matching

- **Methods are case-insensitive**: `"get"`, `"Get"` and `"GET"` all match.
- **Paths are exact against the registration**, parameters included. `"/api/v1/users/:id"` matches
  the registered route; `"/api/v1/users/42"` does not, because no request is being made.

## Optional failure message

Both accept an optional message for context. Only the first supplied argument is used, as literal
text; it is not formatted with additional arguments:

```go
zt.NotHasRoute("DELETE", "/api/v1/users", app, "bulk delete must not be exposed")
```

## When to use them

- **After mounting**: assert the routes a router helper is supposed to add, before writing request
  tests against them.
- **Registration checks**: use `NotHasRoute` to confirm a deprecated route is no longer registered.
  This does not prove a URL is inaccessible; a wildcard or parameterized route may still match it.
  Use request tests to verify actual routing and access control.
- **Not** as a substitute for a request test: registration says nothing about status codes, payloads
  or middleware order.

## Gotchas

- **The app must be the same instance** you registered the routes on. A cloned or rebuilt app starts
  empty.
- **Groups and middleware do not hide routes.** A mounted group's routes are visible with their full
  path.
- **Wildcards match literally.** `"/static/*"` must be asserted as `"/static/*"`.

## Related

- [App helper](../helpers/app.md) — `NewApp`
- [Status](status.md) — asserting what the route actually returns
- [examples/01-flat-style](../../examples/example-01.md)

---

Back to: [Assertions](./README.md) &mdash; Next: [Debugging](debugging.md)

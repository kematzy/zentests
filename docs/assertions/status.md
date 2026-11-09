# Zentests - Docs - Assertions - Status

Use simple status names — `OK()`, `NotFound()` — instead of `Status(404)`, or check the class when
the exact code is not the point.

## Assertions

| Assertion            | Checks                           |
| -------------------- | -------------------------------- |
| `OK()`               | Status is 200                    |
| `Created()`          | Status is 201                    |
| `Accepted()`         | Status is 202                    |
| `NoContent()`        | Status is 204                    |
| `BadRequest()`       | Status is 400                    |
| `Unauthorized()`     | Status is 401                    |
| `Forbidden()`        | Status is 403                    |
| `NotFound()`         | Status is 404                    |
| `Unprocessable()`    | Status is 422                    |
| `ServerError()`      | Status is 500                    |
| `Status(code)`       | Status equals `code`             |
| `StatusIn(codes...)` | Status matches one supplied code |

Usage:

```go
zt.Get(app, "/").OK()
zt.PostJSON(app, "/api/v1/users", data).Created()
zt.Get(app, "/does/not/exist").NotFound()
zt.Delete(app, "/api/v1/users/1").NoContent()
```

### Exact code

`Status(code)` checks an exact code. Use `StatusIn` when several codes are legitimate:

```go
zt.Get(app, "/api/v1/users/1").Status(200)
zt.Get(app, "/api/v1/users/1").StatusIn(200, 204)
```

## Status classes

| Assertion         | Checks                     |
| ----------------- | -------------------------- |
| `IsSuccess()`     | Status is in the 2xx range |
| `IsClientError()` | Status is in the 4xx range |
| `IsServerError()` | Status is in the 5xx range |

Usage:

```go
zt.Get(app, "/api/v1/search").IsSuccess()     // 2xx
zt.Get(app, "/api/v1/nope").IsClientError()   // 4xx
zt.Get(app, "/api/v1/report").IsServerError() // 5xx
```

## Redirects

| Assertion               | Checks                                            |
| ----------------------- | ------------------------------------------------- |
| `IsRedirect()`          | Status is 301, 302, 303, 307, or 308              |
| `RedirectsTo(location)` | Redirect status and exact `Location` header match |

Usage:

```go
zt.Get(app, "/old").IsRedirect()              // 301, 302, 303, 307 or 308
zt.Get(app, "/old").RedirectsTo("/new")       // redirect plus an exact Location match
```

`304 Not Modified` is not included. `RedirectsTo(location)` requires the `Location` header to equal
`location` exactly, including any query string.

## Notes

- Every named assertion is a thin wrapper over `Status(n)`, so a failure message always shows the
  expected and actual code.
- Status assertions do not read the body; combine with a body assertion when both matter:
  `zt.PostJSON(...).Created().IsJSON()`.

## Related

- [Headers](./headers.md) — `Location`, `Content-Type` and the rest
- [examples/01-flat-style](../../examples/example-01.md)

---

Back to: [Assertions](./README.md) &mdash; Next: [Headers](./headers.md)

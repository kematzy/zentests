# Zentests - Docs - Assertions - Status

> **Status:** DRAFT — structure and wording are open for discussion.

Use simple status names — `OK()`, `NotFound()` — instead of `Status(404)`, or check the class when
the exact code is not the point.

## Named statuses

| Assertion         | Code |
| ----------------- | ---- |
| `OK()`            | 200  |
| `Created()`       | 201  |
| `NoContent()`     | 204  |
| `BadRequest()`    | 400  |
| `Unauthorized()`  | 401  |
| `Forbidden()`     | 403  |
| `NotFound()`      | 404  |
| `Unprocessable()` | 422  |
| `ServerError()`   | 500  |

```go
zt.Get(app, "/").OK()
zt.PostJSON(app, "/api/v1/users", data).Created()
zt.Get(app, "/does/not/exist").NotFound()
zt.Delete(app, "/api/v1/users/1").NoContent()
```

## Exact code

`Status(200)` is still there when the code itself is the subject of the test, and `StatusIn` accepts
a set when several codes are legitimate:

```go
zt.Get(app, "/api/v1/users/1").Status(200)
zt.Get(app, "/api/v1/users/1").StatusIn(200, 204)
```

## Status classes

```go
zt.Get(app, "/api/v1/search").IsSuccess()     // 2xx
zt.Get(app, "/api/v1/nope").IsClientError()   // 4xx
zt.Get(app, "/api/v1/report").IsServerError() // 5xx
```

## Redirects

```go
zt.Get(app, "/old").IsRedirect()              // 301, 302, 303, 307 or 308
zt.Get(app, "/old").RedirectsTo("/new")       // redirect plus an exact Location match
```

`RedirectsTo` requires the `Location` header to equal the argument exactly, query string included.

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

# Zentests - Docs - Requests - Headers

Headers sent _with_ the request. For assertions about headers the server returns, see
[assertions/headers](../assertions/headers.md).

## Default headers

`SetHeader` stores a default header for every request from that `zt`; `WithBearer` sets the default
`Authorization: Bearer <token>` header.

| Helper                      | Scope | Behavior                                           | Returns |
| --------------------------- | ----- | -------------------------------------------------- | ------- |
| `zt.SetHeader(key, value)`  | `*T`  | Set or replace a default header                    | `*T`    |
| `zt.WithBearer(token)`      | `*T`  | Set default `Authorization: Bearer <token>`        | `*T`    |
| `jar.SetHeader(key, value)` | `*Jar`| Set or replace a default header for jar requests   | —       |
| `jar.WithBearer(token)`     | `*Jar`| Set default bearer authorization for jar requests  | —       |

```go
zt := zentests.New(t).
	SetHeader("Accept", "application/json").
	SetHeader("X-Api-Key", "key-123")

zt.WithBearer("secret-token") // adds Authorization: Bearer secret-token

zt.Get(app, "/api/v1/users").OK() // sends all three
```

`zt` setters return `*T` and can be chained; jar setters return no value. Header names passed to
`SetHeader` are canonicalized, so `"x-api-key"` and `"X-Api-Key"` configure the same default.

## Precedence

A header belonging to a single request wins over a default of the same name:

```go
zt := zentests.New(t).SetHeader("Content-Type", "text/plain")

zt.PostJSON(app, "/api/v1/users", data) // sends application/json, not text/plain
zt.Post(app, "/raw", []byte("hi"))      // sends text/plain
```

That rule lets `PostJSON`, `PostForm` and `PostMultipart` use their required content type even when
a default `Content-Type` is set. Raw-body methods (`Post`, `Put`, `Patch`) do not infer a content type;
set one with `zt.SetHeader` when needed.

## The same thing on a jar

Jar defaults apply to every jar request. A non-empty header already on a request sent with `Send`
wins over the default.

```go
jar := zentests.NewJar(t, app)
jar.SetHeader("X-Api-Key", "key-123")
jar.WithBearer("secret-token")
jar.Get("/api/v1/users").OK()
```

## Notes

- **Defaults live on the helper, not globally.** Two `zt` values in one test can send different
  headers; nothing is shared between tests.
- **Headers are per helper instance.** In a suite, setting headers in `SetupTest` applies to the
  `*T` created there — see [testing/testify-suite](../testing/testify-suite.md).
- **Per-request headers:** request methods do not accept arbitrary header maps. Use `zt.SetHeader`
  for defaults; the package-level `zentests.SetHeader(key, value)` only creates a one-entry map.

## Related

- [Methods](methods.md) — the requests that carry these headers
- [Test config](testconfigs.md) — timeouts and failure behaviour
- [assertions/headers](../assertions/headers.md) — response headers

---

Back to: [Requests](./README.md) &mdash; Next: [Test config](testconfigs.md)

# Zentests - Docs - Requests - Headers

> **Status:** DRAFT — structure and wording are open for discussion.

Headers sent _with_ the request. For assertions about headers the server returns, see
[assertions/headers](../assertions/headers.md).

## Default headers

`SetHeader` stores a header that every request from that `zt` sends. `WithBearer` is shorthand for
`SetHeader("Authorization", "Bearer …")`.

```go
zt := zentests.New(t).
	SetHeader("Accept", "application/json").
	SetHeader("X-Api-Key", "key-123")

zt.WithBearer("secret-token") // adds Authorization: Bearer secret-token

zt.Get(app, "/api/v1/users").OK() // sends all three
```

Both return `*T`, so they chain onto `New` or onto an earlier call. Header names are canonicalised,
so `"x-api-key"` and `"X-Api-Key"` configure the same header.

## Precedence

A header belonging to a single request wins over a default of the same name:

```go
zt := zentests.New(t).SetHeader("Content-Type", "text/plain")

zt.PostJSON(app, "/api/v1/users", data) // sends application/json, not text/plain
zt.Post(app, "/raw", []byte("hi"))      // sends text/plain
```

That rule is what makes `PostJSON`, `PostForm` and `PostMultipart` safe to use with a default
`Content-Type` set for the rest of the file.

## The same thing on a jar

`NewJar` has the same two methods, and applies the same precedence: a header already on the request
wins.

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
- **Use `SetHeader` for one-off maps.** The package-level `zentests.SetHeader(key, value)` returns a
  `map[string]string`, for the raw-body methods that accept a header map.

## Related

- [Methods](methods.md) — the requests that carry these headers
- [Test config](testconfigs.md) — timeouts and failure behaviour
- [assertions/headers](../assertions/headers.md) — response headers

---

Back to: [Requests](./README.md) &mdash; Next: [Test config](testconfigs.md)

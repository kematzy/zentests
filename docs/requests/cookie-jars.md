# Zentests - Docs - Requests - Cookie Jars

> **Status:** DRAFT — structure and wording are open for discussion.

A jar sends requests to one app and keeps the cookies it sets, the way a browser does. Use it for
flows that span several requests: login, CSRF handshakes, session expiry.

## Create once, reuse

```go
jar := zentests.NewJar(t, app)

// log in: the session cookie is stored
jar.PostForm("/login", url.Values{"email": {"jane@example.com"}, "password": {"secret"}}).
	OK().
	HasCookie("session")

// later requests carry it automatically
jar.Get("/api/v1/me").OK().HasString("user.name", "Jane")

// log out, and the same request is refused
jar.Get("/logout").OK()
jar.Get("/api/v1/me").Unauthorized()
```

The jar is created with `t`, so cleanup is registered for you; there is nothing to close.

## Methods

| Method                                   | Sends                                                         |
| ---------------------------------------- | ------------------------------------------------------------- |
| `jar.Get(path)`                          | GET with the stored cookies                                   |
| `jar.PostForm(path, url.Values)`         | POST, form-encoded, with the stored cookies                   |
| `jar.PostMultipart(path, fields, files)` | POST multipart, with the stored cookies                       |
| `jar.Send(req)`                          | A caller-built `*http.Request`, with the stored cookies added |
| `jar.Cookie(name)`                       | The stored value, as a string                                 |

`Send` is the escape hatch for anything the others do not cover — a `PUT`, a JSON body, a custom
method:

```go
body := strings.NewReader(`{"name":"Jane"}`)
req := httptest.NewRequest(http.MethodPatch, "/api/v1/me", body)
req.Header.Set("Content-Type", "application/json")

jar.Send(req).OK().HasString("name", "Jane")
```

## Headers on a jar

```go
jar.SetHeader("X-Api-Key", "key-123")
jar.WithBearer("secret-token")
```

Both apply to every request the jar sends. A header already on the request wins, so a `Send` request
can override them for one call.

## How cookies are stored

The jar is deliberately flat: **no** `Path`, `Domain` or `Expires` matching. Every stored cookie is
sent with every request to the app. That is enough for the flows these tests care about, and it
keeps the behaviour easy to reason about — but it is not a browser.

A cookie is forgotten when a response sets it with `MaxAge < 0` or an empty value, which is how
`ClearCookie` and logout handlers work.

## Notes

- **One jar per app.** It holds a reference to the app it was built with.
- **Requests are sequential by design.** Sharing one jar across parallel tests races on the cookie
  map; create a jar per test.
- **`jar.Cookie(name)` returns `""` for an unknown cookie** — assert with `HasCookie` when absence
  has to fail the test.

## Related

- [assertions/cookies](../assertions/cookies.md) — per-response cookie assertions
- [Methods](methods.md) — the non-jar request surface
- [examples/01-flat-style](../../examples/example-01.md)

---

Back to: [Requests](README.md)

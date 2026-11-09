# Zentests - Docs - Requests - Cookie Jars

A `Jar` stores response cookies by name and sends them with later requests to one Fiber app. Use it
for multi-request flows such as login, CSRF handshakes, and sessions. It is a lightweight test
helper, not a browser cookie jar; see [How cookies are stored](#how-cookies-are-stored).

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

`zentests.NewJar(t, app)` creates an empty jar bound to `app`; it fails the test if `app` is nil.
The jar stores data in memory and owns no external resources, so no teardown or close call is
needed.

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

Both apply to every request the jar sends. A non-empty header already on a `Send` request wins, so
it can override a default for one call.

## How cookies are stored

The jar is a flat map of cookie names to values: it does not match `Path` or `Domain`, enforce
`Secure` or `HttpOnly`, or expire cookies by `Expires` or positive `MaxAge`. Every stored cookie is
sent with every request to the app. Cookies are keyed only by name, so later same-name entries
replace earlier values or remove the entry. Empty values and `MaxAge < 0` remove a cookie, matching
Fiber's `ClearCookie` behavior.

## Notes

- **Each jar is bound to one app.** Create a jar for the app used by that test flow.
- **Requests are sequential by design.** Sharing one jar across parallel tests races on the cookie
  map; create a jar per test.
- **Jar requests use the package defaults.** There are no `*WithConfig` jar methods; set
  `ZENTESTS_TIMEOUT` for a run-wide timeout.
- **`jar.Cookie(name)` returns `""` for an unknown cookie.** Use `assert.Empty` or `assert.NotEmpty`
  on that value to check whether the jar holds a cookie.

## Related

- [assertions/cookies](../assertions/cookies.md) — per-response cookie assertions
- [Methods](methods.md) — the non-jar request surface
- [examples/01-flat-style](../../examples/example-01.md)

---

Back to: [Requests](./README.md)

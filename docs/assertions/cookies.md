# Zentests - Docs - Assertions - Cookies

These assertions read the `Set-Cookie` headers of one response. For carrying cookies _between_
requests, use a [cookie jar](../requests/cookie-jars.md) instead.

## Assertions

| Assertion                                     | Checks                                                   |
| --------------------------------------------- | -------------------------------------------------------- |
| `HasCookie(name)`                             | A non-deletion cookie with that name was set             |
| `CookieNotPresent(name)`                      | No live cookie with that name was set                    |
| `CookieHasValues(expected map[string]string)` | Listed names have expected values; extra cookies ignored |
| `Cookie(name)`                                | Returns the first matching `*http.Cookie`                |

`CookieHasValues` is a subset check: it verifies each listed name/value pair and ignores other
cookies. It checks values, not attributes.

Usage:

```go
zt.PostJSON(app, "/api/v1/login", data).OK().
	HasCookie("session").
	CookieNotPresent("old_session")

zt.Get(app, "/settings").OK().CookieHasValues(map[string]string{
	"lang":  "en",
	"theme": "dark",
})
```

`CookieHasValues` maps **names to values**. It does not check attributes — `Path`, `Domain`,
`HttpOnly` and friends need `Cookie`:

```go
session := zt.PostJSON(app, "/api/v1/login", data).OK().Cookie("session")

assert.True(t, session.HttpOnly)
assert.True(t, session.Secure)
assert.Equal(t, "/", session.Path)
assert.Equal(t, http.SameSiteLaxMode, session.SameSite)
```

`Cookie` fails the test when the cookie was not set, so it is safe to dereference the result.

## Deletion cookies

A cookie that clears itself — empty value or `MaxAge < 0` — counts as **not set** for `HasCookie`
and `CookieNotPresent`. That keeps logout tests honest:

```go
zt.Get(app, "/logout").OK().CookieNotPresent("session") // a deletion cookie passes
```

To confirm a session is really gone, keep using the jar: the deleted cookie must no longer be sent.

## Gotchas

- **A response can set several cookies with the same name** (different paths). `Cookie`,
  `HasCookie`, and `CookieNotPresent` inspect the first match. `CookieHasValues` collapses cookies
  into a name/value map, so the last cookie with a duplicate name wins. These helpers do not compare
  same-name cookies by path.
- **Cookies are per response, not per app.** The session cookie set on login is not sent by the next
  request unless a jar is carrying it.
- **`HttpOnly` cookies are invisible to JavaScript but visible here** — the test reads the raw
  header, which is exactly what makes these assertions useful.

## Related

- [Cookie jars](../requests/cookie-jars.md) — multi-step session flows
- [Headers](./headers.md) — `Set-Cookie` as a raw header
- [examples/bdd](../../examples/example-02.md)

---

Back to: [Assertions](./README.md) &mdash; Next: [Content types](content-types.md)

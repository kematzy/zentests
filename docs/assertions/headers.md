# Zentests - Docs - Assertions - Headers

Header checks that read like the header, without reaching into `resp.Header.Get(...)` by hand.

## Assertions

| Assertion                        | Checks                                      |
| -------------------------------- | ------------------------------------------- |
| `HasHeader(key, value)`          | The first header value equals `value`       |
| `HeaderPresent(key)`             | The first header value is non-empty         |
| `HeaderNotPresent(key)`          | The first header value is missing or empty  |
| `HeaderContains(key, substring)` | The first header value contains `substring` |
| `HeaderHasValues(key, values)`   | All values exactly match, in any order      |

Usage:

```go
zt.Get(app, "/api/v1/users").HasHeader("X-Served-By", "api-1")
zt.Get(app, "/api/v1/users").HeaderPresent("X-Request-Id")
zt.Get(app, "/api/v1/users").HeaderNotPresent("X-Powered-By")
zt.Get(app, "/api/v1/users").HeaderContains("Cache-Control", "max-age")
zt.Get(app, "/api/v1/users").HeaderHasValues("X-Features", []string{"audio", "video"})
```

Header names are case-insensitive. `HasHeader` compares the first value exactly; `HeaderContains`
checks for a substring in that value. `HeaderHasValues` compares the full set of values. For
example, `HasHeader("Cache-Control", "max-age=60")` requires that exact value.

## Multi-value headers

`HeaderHasValues` compares the full set with `assert.ElementsMatch`: no extras, no missing entries,
and order is irrelevant. `HeaderContains` checks only the first value returned by `Header.Get`; it
does not search later values when a header has several separate values.

```go
// exactly these two, in either order
zt.Get(app, "/assets/app.css").HeaderHasValues("Vary", []string{"Accept", "Accept-Encoding"})
```

## Check for absent headers

Use `HeaderNotPresent` to enforce a response-header policy, such as omitting a framework banner.
This checks only the selected header; it does not establish overall security or check the response
body. `HeaderPresent` and `HeaderNotPresent` treat an empty first value as absent.

## Content types are header checks too

`IsJSON()`, `IsHTML()`, `IsPNG()` and friends are `HeaderContains("Content-Type", …)` with a
readable name. They live on their own page: [Content types](content-types.md).

## Related

- [Content types](content-types.md)
- [Default headers](../requests/headers.md) — headers sent _with_ the request
- [examples/01-flat-style](../../examples/example-01.md)

---

Back to: [Assertions](./README.md) &mdash; Next: [Response Body](body-html.md)

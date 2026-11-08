# Zentests - Docs - Assertions - Headers

> **Status:** DRAFT — structure and wording are open for discussion.

Header checks that read like the header, without reaching into `resp.Header.Get(...)` by hand.

## The five assertions

| Assertion                        | Checks                                              |
| -------------------------------- | --------------------------------------------------- |
| `HasHeader(key, value)`          | The header exists and its value is exactly `value`  |
| `HeaderPresent(key)`             | The header exists and is not empty                  |
| `HeaderNotPresent(key)`          | The header is missing or empty                      |
| `HeaderContains(key, substring)` | The header value contains `substring`               |
| `HeaderHasValues(key, values)`   | The header holds exactly those values, in any order |

```go
zt.Get(app, "/api/v1/users").HasHeader("X-Served-By", "api-1")
zt.Get(app, "/api/v1/users").HeaderPresent("X-Request-Id")
zt.Get(app, "/api/v1/users").HeaderNotPresent("X-Powered-By")
zt.Get(app, "/api/v1/users").HeaderContains("Cache-Control", "max-age")
zt.Get(app, "/api/v1/users").HeaderHasValues("X-Features", []string{"audio", "video"})
```

Header names are case-insensitive, as HTTP requires. Values are compared exactly, so `max-age=60`
matches `max-age=60` and nothing else.

## Multi-value headers

`HeaderHasValues` compares the full set with `assert.ElementsMatch`: no extras, no missing entries,
order irrelevant. Use `HeaderContains` when you only care that one value is among others.

```go
// exactly these two, in either order
zt.Get(app, "/assets/app.css").HeaderHasValues("Vary", []string{"Accept", "Accept-Encoding"})
```

## Absence is a security check

`HeaderNotPresent("X-Powered-By")` and `HeaderNotPresent("Server")` are the cheapest way to assert
that a stack-trace or framework banner is not leaking. `HeaderPresent` and `HeaderNotPresent` treat
an empty value as absent.

## Content types are header checks too

`IsJSON()`, `IsHTML()`, `IsPNG()` and friends are `HeaderContains("Content-Type", …)` with a
readable name. They live on their own page: [Content types](content-types.md).

## Related

- [Content types](content-types.md)
- [Default headers](../requests/headers.md) — headers sent _with_ the request
- [examples/01-flat-style](../../examples/example-01.md)

---

Back to: [Assertions](./README.md) &mdash; Next: [HTML body](body-html.md)

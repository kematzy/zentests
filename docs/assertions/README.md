# Zentests - Docs - Assertions

Response assertions are methods on the `*Response` returned by each request. They chain, report
failures against the current test instead of returning errors, and use names that read like
expectations. Route assertions are methods on `*T` and do not send a request.

## Contents

| Page                              | What it answers                                              |
| --------------------------------- | ------------------------------------------------------------ |
| [Status](status.md)               | Named codes, status classes, redirects                       |
| [Headers](./headers.md)           | Header presence, exact values, multi-value headers           |
| [Response Body](body-html.md)     | Substring, regexp, exact body, empty body, raw body          |
| [JSON body](body-json.md)         | Dot-notation paths, typed helpers, arrays, subsets, decoding |
| [Cookies](cookies.md)             | What the server set: names, values, attributes               |
| [Content types](content-types.md) | `IsJSON()`, `IsHTML()`, `IsPNG()` and the rest               |
| [Routes](routes.md)               | `HasRoute`/`NotHasRoute`, without a request                  |
| [Debugging](debugging.md)         | `Debug`, `Dump`, and which headers are redacted              |

## Response assertion behavior

- **Chainable.** Response assertions return the same `*Response`, so one response supports many
  checks in one expression. Value helpers such as `Cookie()` and `Body()` return their values
  instead.
- **Test-failing, not error-returning.** There is nothing to check with `if err != nil`; a failed
  assertion records a failure against the current test.
- **Continue by default.** Most assertions record a failure and let the chain run on, so one run can
  report several problems. `Bind` and `Cookie` stop the test when they fail. A failed JSON parse is
  recorded as a failure but does not stop the test.
- **Read from the same response.** The body is read once and cached, so later assertions reuse the
  bytes instead of re-reading or consuming the response stream.

## Two rules that explain most surprises

1. **JSON numbers decode as `float64`.** `Has` normalizes Go numeric expectations, so both
   `Has("age", 42)` and `Has("age", 42.0)` can match the same JSON number. Use `HasInt` or
   `HasFloat` to make the expected numeric kind explicit; use typed helpers for booleans too.
2. **JSON paths start at an object.** A top-level array cannot be addressed by path; either wrap it
   (`{"users": [...]}`) or decode with `Bind`. See [JSON body](body-json.md).

## Related

- [Requests](../requests/README.md) — how to get the `*Response`
- [Testing styles](../testing/README.md) — where to create `zt`
- [examples/01-flat-style](../../examples/example-01.md) — most of these assertions in a runnable
  file

---

Back to: [Docs](../README.md) &mdash; Next: [Status](status.md)

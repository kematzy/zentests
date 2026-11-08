# Zentests - Docs - Assertions

> **Status:** DRAFT — structure and wording are open for discussion.

Assertions are methods on the `*Response` returned by every request. They chain, they fail the test
instead of returning an error, and they are named after the thing being checked, so a failing line
reads like the expectation.

## Contents

| Page                              | What it answers                                              |
| --------------------------------- | ------------------------------------------------------------ |
| [Status](status.md)               | Named codes, status classes, redirects                       |
| [Headers](./headers.md)           | Header presence, exact values, multi-value headers           |
| [HTML body](body-html.md)         | Substring, regexp, exact body, empty body, raw body          |
| [JSON body](body-json.md)         | Dot-notation paths, typed helpers, arrays, subsets, decoding |
| [Cookies](cookies.md)             | What the server set: names, values, attributes               |
| [Content types](content-types.md) | `IsJSON()`, `IsHTML()`, `IsPNG()` and the rest               |
| [Routes](routes.md)               | `HasRoute`/`NotHasRoute`, without a request                  |
| [Debugging](debugging.md)         | `Debug`, `Dump`, and which headers are redacted              |

## What every assertion shares

- **Chainable.** Each assertion returns the same `*Response`, so one response supports many checks
  in one expression.
- **Test-failing, not error-returning.** There is nothing to check with `if err != nil`; a failed
  assertion records a failure against the current test.
- **Continue by default.** Most assertions record the failure and let the chain run on, so one run
  reports several problems. The few that cannot continue — `Bind`, `Cookie`, and a failed JSON parse
  — stop the test instead.
- **Read from the same response.** The body is read once and cached, so repeated assertions cost
  nothing and cannot consume it.

## Two rules that explain most surprises

1. **Numbers in JSON are `float64`.** `Has("age", 42)` fails where `HasInt("age", 42)` passes. Use
   the typed helpers for numbers and booleans.
2. **JSON paths start at an object.** A top-level array cannot be addressed by path; either wrap it
   (`{"users": [...]}`) or decode with `Bind`. See [JSON body](body-json.md).

## Related

- [Requests](../requests/README.md) — how to get the `*Response`
- [Testing styles](../testing/README.md) — where to create `zt`
- [examples/01-flat-style](../../examples/example-01.md) — most of these assertions in a runnable
  file

---

Back to: [Docs](../README.md) &mdash; Next: [Status](status.md)

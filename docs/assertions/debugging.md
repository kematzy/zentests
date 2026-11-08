# Zentests - Docs - Assertions - Debugging

> **Status:** DRAFT — structure and wording are open for discussion.

When an assertion fails and the message is not enough, `Debug` and `Dump` show the response as the
test sees it — status, headers and body — with credentials redacted.

## `Debug()` prints, `Dump()` returns

```go
// during development: print the response, then keep chaining
zt.Get(app, "/api/v1/users").Debug().OK().IsJSON()

// in a failure message or a log line
t.Logf("response: %s", zt.Get(app, "/api/v1/users").Dump())
```

`Debug()` writes to stdout and returns the `*Response`, so it can sit anywhere in a chain. `Dump()`
returns the same content as a string, for `t.Logf` or your own reporting.

Both show:

```http
Status: 200
Headers: map[Content-Type:[application/json] X-Request-Id:[abc123]]
Body: {"users":[{"id":1,"name":"Jane"}]}
```

`Debug()` additionally prints the parsed JSON, indented, when the body was already parsed.

## Redacted headers

Credential-bearing headers are masked before printing, so a debug run does not leak tokens into CI
logs:

`Authorization`, `Cookie`, `Set-Cookie`, `X-Api-Key`, `X-Auth-Token`, `X-Access-Token`,
`X-Refresh-Token`, `X-Session-Id`, `X-Csrf-Token`.

```go
// prints Authorization: [REDACTED]
zt.Get(app, "/api/v1/me").WithBearer("secret-token").Debug()
```

Everything else — including the body — is printed verbatim, so still avoid `Debug()` on responses
that contain personal data.

## When to use them

- **A failing chain with several assertions**: `Debug()` before the assertions shows which of them
  was the first to disagree.
- **Unexpected 404 or 500**: the body usually explains it faster than reading the handler.
- **Debugging the test, not the app**: `Dump()` in `t.Logf` keeps output quiet until `-v`.

## Notes

- Neither assertion fails a test; they are observation, not assertion. Do not assert on `Dump()`
  output — it is a debugging aid and its format may change.
- `Debug()` writes with `fmt.Printf`, so it bypasses `testing`'s output capture and appears even
  when a test passes. Remove the call once the test is fixed.

## Related

- [Status](status.md) · [Headers](./headers.md) — assert what `Debug` shows
- [Content types](content-types.md) — `IsJSON()` before a JSON assertion

---

Back to: [Assertions](./README.md)

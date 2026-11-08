# Zentests - Docs - Assertions - HTML Body

> **Status:** DRAFT — structure and wording are open for discussion.

Most HTML tests are about a word being there, or not being there. These assertions say exactly that,
and hand you the raw body when you would rather parse it yourself.

## The assertions

| Assertion              | Checks                                  |
| ---------------------- | --------------------------------------- |
| `Contains(text)`       | The body contains `text`                |
| `NotContains(text)`    | The body does not contain `text`        |
| `BodyMatches(pattern)` | The body matches the regular expression |
| `Equals(text)`         | The body is exactly `text`              |
| `IsEmpty()`            | The body has no content                 |

```go
zt.Get(app, "/").OK().Contains("Welcome back")
zt.Get(app, "/checkout").OK().NotContains("Your cart is empty")
zt.Get(app, "/robots.txt").OK().Equals("User-agent: *")
zt.Get(app, "/blank").OK().IsEmpty()
```

`Contains` and `NotContains` are case-sensitive substring checks — good for exact copy, brittle for
prose that gets reworded. When the _shape_ is stable but the wording is not, use a regexp:

```go
zt.Get(app, "/contact").OK().BodyMatches(`Email: \S+@\S+`)
```

The pattern is a Go regular expression, unanchored like `regexp.MatchString`, so add `^` and `$`
when the whole body must match.

## The raw body

```go
body := zt.Get(app, "/feed.xml").OK().Body()        // []byte
text := zt.Get(app, "/feed.xml").OK().BodyString()  // string
```

Both read the cached body; neither consumes it, so further assertions still work.

## Gotchas

- **Regexp metacharacters.** `BodyMatches("@email.com")` matches `@emailXcom` too — escape the dots:
  `@email\.com`.
- **HTML is not parsed.** These are string assertions. To check structure, parse `Body()` with
  `golang.org/x/net/html` in your own helper.
- **`Equals` compares bytes.** Trailing newlines count; prefer `Contains` unless the body is
  generated exactly.

## Related

- [JSON body](body-json.md) — when the response is JSON instead of HTML
- [Content types](content-types.md) — assert `text/html` first if the route can return several
  formats
- [examples/01-flat-style](../../examples/example-01.md)

---

Back to: [Assertions](./README.md) &mdash; Next: [JSON body](body-json.md)

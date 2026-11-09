# Zentests - Docs - Assertions - Content Types

Content-type assertions are header checks with readable names: `IsJSON()` instead of
`HeaderContains("Content-Type", "application/json")`.

## Assertions

| Group                 | Assertions                                                           |
| --------------------- | -------------------------------------------------------------------- |
| Text formats          | `IsJSON()`, `IsHTML()`, `IsXML()`, `IsPlainText()`                   |
| Assets                | `IsCSS()`, `IsJS()`                                                  |
| Images                | `IsImage()`, `IsPNG()`, `IsJPEG()`, `IsGIF()`, `IsSVG()`, `IsWebP()` |
| Other response header | `IsXHR()` — `X-Requested-With` contains `XMLHttpRequest`             |

Usage:

```go
zt.Get(app, "/api/v1/users").OK().IsJSON()
zt.Get(app, "/").OK().IsHTML()
zt.Get(app, "/assets/app.css").OK().IsCSS()
zt.Get(app, "/logo.png").OK().IsPNG().IsImage()
zt.Get(app, "/api/v1/users").IsXHR()
```

The content-type shortcuts use case-sensitive substring checks, so parameters do not matter:
`application/json; charset=utf-8` passes `IsJSON()`. Use the canonical lowercase media-type values.
`IsJSON()` looks for the literal `application/json`; it does not match every JSON-based media type,
such as `application/problem+json`.

## Exact content type

When the whole value matters, compare it:

```go
zt.Get(app, "/page").HasContentType("text/html; charset=utf-8")
```

`HasContentType` is an exact match on the header value — including parameters and spacing. For
"returns HTML in some form", `IsHTML()` is the safer assertion.

## Choosing between them

| You want                            | Use                   |
| ----------------------------------- | --------------------- |
| "This endpoint returns JSON"        | `IsJSON()`            |
| "This endpoint returns HTML"        | `IsHTML()`            |
| "This is one of the image types"    | `IsImage()`           |
| "This is specifically a PNG"        | `IsPNG()`             |
| "The header is exactly this string" | `HasContentType("…")` |

## Gotchas

- **`IsImage()` is a family check.** It passes for any `image/*` content type; use the specific
  assertion when the format is part of the contract.
- **A missing `Content-Type` fails every shortcut**, which is usually the true cause when a handler
  forgets to set one.
- **`IsXHR()` checks the response header.** It looks for `XMLHttpRequest` in the response's
  `X-Requested-With` header; it does not assert that the outgoing request was sent as XHR. This
  header is conventionally used on requests, so do not use this method as a request-header
  assertion.

## Related

- [Headers](./headers.md) — every other header assertion
- [Response Body](body-html.md) · [JSON body](body-json.md) — assertions on the parsed content
- [examples/01-flat-style](../../examples/example-01.md)

---

Back to: [Assertions](./README.md) &mdash; Next: [Routes](routes.md)

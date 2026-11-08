# Zentests - Docs - Assertions - Content Types

> **Status:** DRAFT — structure and wording are open for discussion.

Content-type assertions are header checks with readable names: `IsJSON()` instead of
`HeaderContains("Content-Type", "application/json")`.

## The shortcuts

| Group        | Assertions                                                           |
| ------------ | -------------------------------------------------------------------- |
| Text formats | `IsJSON()`, `IsHTML()`, `IsXML()`, `IsPlainText()`                   |
| Assets       | `IsCSS()`, `IsJS()`                                                  |
| Images       | `IsImage()`, `IsPNG()`, `IsJPEG()`, `IsGIF()`, `IsSVG()`, `IsWebP()` |
| Requests     | `IsXHR()` — the `X-Requested-With: XMLHttpRequest` header            |

```go
zt.Get(app, "/api/v1/users").OK().IsJSON()
zt.Get(app, "/").OK().IsHTML()
zt.Get(app, "/assets/app.css").OK().IsCSS()
zt.Get(app, "/logo.png").OK().IsPNG().IsImage()
zt.Get(app, "/api/v1/users").IsXHR()
```

Each one is a substring check on the `Content-Type` header, so parameters do not matter:
`application/json; charset=utf-8` passes `IsJSON()`. Matching is case-sensitive, as the header
values in practice are.

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
- **`IsXHR()` is not a content type.** It reads `X-Requested-With`; it is grouped here because it
  answers the same question — "what kind of response is this?".

## Related

- [Headers](./headers.md) — every other header assertion
- [HTML body](body-html.md) · [JSON body](body-json.md) — assertions on the parsed content
- [examples/01-flat-style](../../examples/example-01.md)

---

Back to: [Assertions](./README.md) &mdash; Next: [Routes](routes.md)

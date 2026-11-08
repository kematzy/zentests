# Zentests - Docs - Requests - Multipart

> **Status:** DRAFT — structure and wording are open for discussion.

`PostMultipart` builds a `multipart/form-data` body from a field map and a list of files, sets the
`Content-Type` with the boundary, and sends it in one call.

## Signature

```go
zt.PostMultipart(app, path, fields, files) *Response
```

- `fields` — `map[string]string` of ordinary form fields.
- `files` — `[]FileUpload`, where each entry is `{Name, Filename, Content []byte}`. The bytes are in
  hand, so no fixture file has to exist on disk.

## Uploading one file

```go
resp := zt.PostMultipart(app, "/api/v1/avatar",
	map[string]string{"user_id": "42"},
	[]zentests.FileUpload{
		{Name: "avatar", Filename: "jane.png", Content: pngBytes},
	},
)

resp.OK().HasInt("id", 42)
```

## Several files, same field

`FileUpload.Name` is the form field name, and the same name may appear more than once — that is how
HTML multiple-file inputs arrive:

```go
zt.PostMultipart(app, "/api/v1/documents",
	nil, // no plain fields
	[]zentests.FileUpload{
		{Name: "files", Filename: "a.pdf", Content: pdfA},
		{Name: "files", Filename: "b.pdf", Content: pdfB},
	},
).Created()
```

## Reading files from disk

The helper takes bytes, not paths, so a fixture is two lines:

```go
content, err := os.ReadFile("testdata/jane.png")
require.NoError(t, err)

zt.PostMultipart(app, "/api/v1/avatar", nil, []zentests.FileUpload{
	{Name: "avatar", Filename: "jane.png", Content: content},
}).OK()
```

## Notes

- **Content-Type is set for you**, boundary included. A default `Content-Type` on `zt` does not
  override it — per-request headers win.
- **`nil` maps and slices are fine**; a request with fields only, or files only, works.
- **No `*WithConfig` twin.** A multipart request uses the default timeout; raise it run-wide with
  `ZENTESTS_TIMEOUT` for a slow upload path.
- **Field order is not guaranteed**; servers must not depend on it, and neither should tests.

## Related

- [Methods](methods.md) — the rest of the request surface
- [Test config](testconfigs.md) — timeouts
- [assertions/body-json](../assertions/body-json.md) — asserting the upload response

---

Back to: [Requests](./README.md) &mdash; Next: [Cookie jars](cookie-jars.md)

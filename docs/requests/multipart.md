# Zentests - Docs - Requests - Multipart

`PostMultipart` builds a `multipart/form-data` body from a field map and a list of files, sets the
`Content-Type` with the boundary, and sends it in one call.

## Signature

```go
zt.PostMultipart(app, path, fields, files) *Response
```

| Input    | Type                | Meaning                                              |
| -------- | ------------------- | ---------------------------------------------------- |
| `fields` | `map[string]string` | Ordinary form fields; may be `nil`                   |
| `files`  | `[]FileUpload`      | Files built from in-memory byte slices; may be `nil` |

Each `FileUpload` has these fields:

| Field      | Meaning                                        |
| ---------- | ---------------------------------------------- |
| `Name`     | Form field name; can repeat for multiple files |
| `Filename` | Filename presented to the server               |
| `Content`  | File contents as `[]byte`                      |

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
- **Part order is deterministic:** fields are written in sorted key order, followed by files in
  slice order. Avoid making handler behavior depend on multipart part order.

## Related

- [Methods](methods.md) — the rest of the request surface
- [Test config](testconfigs.md) — timeouts
- [assertions/body-json](../assertions/body-json.md) — asserting the upload response

---

Back to: [Requests](./README.md) &mdash; Next: [Cookie jars](cookie-jars.md)

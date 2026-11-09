# Zentests - Docs

Reference pages for each part of `zentests`. Every area has an overview page that states what the
area covers, then one page per topic with examples and edge cases.

The [README](../README.md) stays short on purpose; the detail lives here.

## Start here

| Page                                | What it answers                                                                             |
| ----------------------------------- | ------------------------------------------------------------------------------------------- |
| [Testing styles](testing/README.md) | Flat, BDD or `testify/suite` — and what changes between them                                |
| [Comparison](comparison.md)         | The same test in `net/http`, `testify` and `zentests`; what it removes, what it does not do |
| [Assertions](assertions/README.md)  | Status, headers, bodies, JSON, cookies, content types, routes, debugging                    |
| [Requests](requests/README.md)      | Request methods, default headers, `TestConfig` and timeouts, multipart, cookie jars         |
| [Database](database/README.md)      | In-memory SQLite setup, records, mutations, transactions                                    |
| [Helpers](helpers/README.md)        | `NewApp`, `Logger`, and `MockLogger`                                                       |
| [Contributing](contributing.md)     | Local setup, `make check`, branch policy, release                                           |
| [API reference](API.md)             | Generated signatures and package comments; the pages above are task-oriented guides        |

## Where to begin

- **New to `zentests`** — [README Quick Start](../README.md#quick-start), then
  [Testing styles](testing/README.md).
- **Coming from `testify` + `httptest`** — [Comparison](comparison.md) shows the same test side by
  side and what a migration looks like.
- **Testing routes and responses** — [Assertions](assertions/README.md) and
  [Requests](requests/README.md).
- **Testing with a database** — [Database](database/README.md); the built-in SQLite constructors
  require CGO and a C compiler.

## Examples

Each style has a self-contained file you can copy and paste:

- [examples/01-flat-style](../examples/example-01.md)
- [examples/02-bdd](../examples/example-02.md)
- [examples/03-testify-suite](../examples/example-03.md)

## Conventions

- Code in these pages is illustrative; the files under `examples/` are the ones that run.
- `zt` is the helper returned by `zentests.New(t)`, `app` is the `*fiber.App` under test.
- JSON paths use dot notation, with array indices as numbers: `users.0.name`.

---

Back to: [README](../README.md)

# Zentests - Docs - Assertions - JSON Body

Assertions to check the presence of a JSON path, or compare values or check a selected subset of
JSON paths. Dot notation reaches nested keys without requiring a struct.

## Assertions

| Assertion                                | Checks                                           |
| ---------------------------------------- | ------------------------------------------------ |
| `HasKey(path)`                           | The key exists, whatever its value               |
| `Has(path, want)`                        | The value equals `want`                          |
| `HasString(path, want)`                  | The value equals `want` as a string              |
| `HasInt(path, want)`                     | The numeric value equals the integer expectation |
| `HasFloat(path, want)`                   | The numeric value is within `0.0001` of `want`   |
| `HasBool(path, want)`                    | The value equals `want` as a boolean             |
| `JSONHasValues(expected map[string]any)` | Subset check: only the listed paths are compared |
| `IsNull(path)` / `IsNotNull(path)`       | The value is JSON `null` / is not                |
| `MatchesRegex(path, pattern)`            | The string value matches a regular expression    |
| `ArrayLength(path, n)`                   | The array has `n` elements                       |
| `ArrayContains(path, want)`              | The array holds `want`                           |
| `ArrayEmpty(path)`                       | The array has no elements                        |

## Parse once, assert many times

The `JSON()` helper parses the body and caches the result on the response. Every JSON assertion
calls it implicitly, so the body is decoded at most once per request and a parse failure is reported
once.

```go
zt.Get(app, "/api/v1/users/1").JSON().  // explicit, when you want the parse to be the first check
	HasKey("user.name").
	HasString("user.name", "John").
	HasInt("user.age", 42).
	ArrayContains("user.roles", "admin")
```

## Paths

Paths use dots, and numeric segments address arrays:

```go
zt.Get(app, "/api/v1/users").Has("users.0.address.city", "Berlin")
zt.Get(app, "/api/v1/users/1").HasKey("user.roles.0")
```

A missing key fails with the path in the message, so a typo in the test is obvious.

```go
zt.Get(app, "/api/v1/users/1").JSONHasValues(map[string]any{
	"user.name":   "John",
	"user.active": true,
})

zt.Get(app, "/api/v1/users/1").IsNotNull("user.email")
zt.Get(app, "/api/v1/users/1").MatchesRegex("user.email", `@example\.com$`)
zt.Get(app, "/api/v1/users").ArrayLength("users", 3).ArrayEmpty("warnings")
```

## Helpers

| Helper         | Behavior                                                                    |
| -------------- | --------------------------------------------------------------------------- |
| `JSON()`       | Parses and caches a JSON object for path-based assertions                   |
| `Bind(target)` | Decodes the body into a pointer target with `encoding/json`; fatal on error |

### Decoding into a struct

When paths get in the way, decode the body into your own type. `Bind` reads the raw body, so it
works for top-level arrays and for documents with awkward key names:

```go
var user User
zt.Get(app, "/api/v1/users/1").OK().Bind(&user)

var users []User
zt.Get(app, "/api/v1/users").OK().Bind(&users)
```

`Bind` stops the test when the body does not fit the target.

## Gotchas

- **Numbers are `float64`.** `Has` normalizes numeric expectations, so `Has("age", 42)` and
  `Has("age", 42.0)` both match a JSON number with value 42. `HasInt` makes an integer expectation
  explicit; `HasFloat` compares within a `0.0001` delta. JSON decoding does not preserve whether a
  number was written as `42` or `42.0`.
- **The root must be an object.** A response of `[{"id":1}]` cannot be addressed by path — use
  `Bind`, or have the API return `{"users": [...]}`.
- **The empty path is not the root.** `ArrayLength("", 2)` looks up a key named `""`; name the key
  (`users`).
- **Keys containing dots have no escape syntax.** Dot notation splits every path on `.`; use `Bind`
  when a JSON key itself contains a dot.
- **`JSONMatches` is deprecated.** It is an alias for `JSONHasValues`; use the latter.
- **Type helpers are strict about JSON types.** `HasString("user.age", "42")` fails when the JSON
  holds a number, even though the text looks the same.

## Related

- [Response Body](body-html.md) — for non-JSON responses
- [Content types](content-types.md) — `IsJSON()`
- [examples/01-flat-style](../../examples/example-01.md)

---

Back to: [Assertions](./README.md) &mdash; Next: [Cookies](cookies.md)

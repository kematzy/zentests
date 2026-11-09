# Zentests - Docs - Requests - Methods

Each request method takes `app` and `path`, plus any method-specific input, and returns a chainable
`*Response`. For request headers, see [Headers](./headers.md); for timeouts, see
[Test config](testconfigs.md).

---

## Methods reference

All methods take `app *fiber.App` and `path string`, and return `*Response` for chaining.

| Method                                    | Input                               | Behavior                                    |
| ----------------------------------------- | ----------------------------------- | ------------------------------------------- |
| `Get(app, path)`                          | —                                   | GET                                         |
| `Post(app, path, body)`                   | `body []byte`                       | POST raw body; no content type is inferred  |
| `PostJSON(app, path, data)`               | `data any`                          | POST JSON; sets `Content-Type`              |
| `PostForm(app, path, data)`               | `data map[string]string`            | POST URL-encoded form                       |
| `Put(app, path, body)`                    | `body []byte`                       | PUT raw body; no content type is inferred   |
| `PutJSON(app, path, data)`                | `data any`                          | PUT JSON; sets `Content-Type`               |
| `Patch(app, path, body)`                  | `body []byte`                       | PATCH raw body; no content type is inferred |
| `PatchJSON(app, path, data)`              | `data any`                          | PATCH JSON; sets `Content-Type`             |
| `Delete(app, path)`                       | —                                   | DELETE                                      |
| `DeleteJSON(app, path, data)`             | `data any`                          | DELETE with JSON body; sets `Content-Type`  |
| `PostMultipart(app, path, fields, files)` | `map[string]string`, `[]FileUpload` | POST multipart form data and files          |
| `GetQuery(app, path, query)`              | `url.Values`                        | GET with URL-encoded query parameters       |

`GetQuery` URL-encodes the values and appends them to `path` (`&` when `path` already has a query).
A nil or empty `url.Values` still appends the query separator.

## HTTP Methods with TestConfig

The ten verb methods have a `*WithConfig` twin that takes a `fiber.TestConfig` for a per-call
timeout. `GetQuery` and `PostMultipart` have none — they use the defaults described in
[Test config](testconfigs.md).

| Method                                       | Description                        |
| -------------------------------------------- | ---------------------------------- |
| `GetWithConfig(app, path, cfg)`              | GET with custom TestConfig         |
| `PostWithConfig(app, path, body, cfg)`       | POST with custom TestConfig        |
| `PostJSONWithConfig(app, path, data, cfg)`   | POST JSON with custom TestConfig   |
| `PostFormWithConfig(app, path, data, cfg)`   | POST form with custom TestConfig   |
| `PutWithConfig(app, path, body, cfg)`        | PUT with custom TestConfig         |
| `PutJSONWithConfig(app, path, data, cfg)`    | PUT JSON with custom TestConfig    |
| `PatchWithConfig(app, path, body, cfg)`      | PATCH with custom TestConfig       |
| `PatchJSONWithConfig(app, path, data, cfg)`  | PATCH JSON with custom TestConfig  |
| `DeleteWithConfig(app, path, cfg)`           | DELETE with custom TestConfig      |
| `DeleteJSONWithConfig(app, path, data, cfg)` | DELETE JSON with custom TestConfig |

## TestConfig Usage

```go
import (
    "time"
    "github.com/gofiber/fiber/v3"
    "github.com/kematzy/zentests"
)

func TestSlowEndpoint(t *testing.T) {
    zt := zentests.New(t)
    app := fiber.New()

    // Custom timeout for slow endpoints
    resp := zt.GetWithConfig(app, "/slow-endpoint", fiber.TestConfig{
        Timeout:       10 * time.Second,
        FailOnTimeout: true,
    })
    resp.OK().IsJSON()

    // Custom timeout with FailOnTimeout behavior
    resp2 := zt.PostJSONWithConfig(app, "/api/users", userData, fiber.TestConfig{
        Timeout:       5 * time.Second,
        FailOnTimeout: true, // Stop the test if the request times out
    })
    resp2.Created()
}
```

### TestConfig Options

| Option          | Type            | Package default | Description                                                                  |
| --------------- | --------------- | --------------- | ---------------------------------------------------------------------------- |
| `Timeout`       | `time.Duration` | `1s`            | Maximum request duration                                                     |
| `FailOnTimeout` | `bool`          | `true`          | On timeout, stop the test (`true`) or use Fiber's partial response (`false`) |

Plain methods pass a zero-value config, which `zentests` replaces with these defaults. A non-zero
`*WithConfig` value is passed as-is, not merged with the defaults; set `FailOnTimeout` explicitly if
you want it enabled. See [Test config](testconfigs.md) for details.

---

Back to: [Requests](./README.md) &mdash; Next: [Headers](./headers.md)

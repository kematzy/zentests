# Zentests - Docs - Requests - Methods

> **Status:** DRAFT — structure and wording are open for discussion.

One call per request: the app, the path, and back a `*Response` ready to chain. For headers sent
with every request see [Headers](./headers.md); for timeouts see [Test config](testconfigs.md).

---

## Methods reference

All methods accept `*fiber.App` and `path`, return `*Response` for chaining:

| Method                                    | Description                             |
| ----------------------------------------- | --------------------------------------- |
| `Get(app, path)`                          | GET request                             |
| `Post(app, path, body)`                   | POST with raw body                      |
| `PostJSON(app, path, data)`               | POST with JSON body                     |
| `PostForm(app, path, data)`               | POST with form data                     |
| `Put(app, path, body)`                    | PUT request                             |
| `PutJSON(app, path, data)`                | PUT with JSON                           |
| `Patch(app, path, body)`                  | PATCH request                           |
| `PatchJSON(app, path, data)`              | PATCH with JSON                         |
| `Delete(app, path)`                       | DELETE request                          |
| `DeleteJSON(app, path, data)`             | DELETE with JSON body                   |
| `PostMultipart(app, path, fields, files)` | POST with multipart form data and files |
| `GetQuery(app, path, url.Values)`         | GET with encoded query parameters       |

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
        Timeout: 10 * time.Second,
    })
    resp.OK().IsJSON()

    // Custom timeout with FailOnTimeout behavior
    resp2 := zt.PostJSONWithConfig(app, "/api/users", userData, fiber.TestConfig{
        Timeout:       5 * time.Second,
        FailOnTimeout: true, // Returns error instead of partial response
    })
    resp2.Created()
}
```

### TestConfig Options

| Option          | Type            | Default | Description                                                |
| --------------- | --------------- | ------- | ---------------------------------------------------------- |
| `Timeout`       | `time.Duration` | `1s`    | Request timeout duration                                   |
| `FailOnTimeout` | `bool`          | `true`  | Return error on timeout (true) or partial response (false) |

A zero-value `fiber.TestConfig{}` uses the defaults above, which is what the plain methods (`Get`,
`PostJSON`, …) pass. See [Test config](testconfigs.md) for `ZENTESTS_TIMEOUT` and `FailOnTimeout`.

---

Back to: [Requests](./README.md) &mdash; Next: [Headers](./headers.md)

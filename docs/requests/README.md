# Zentests - Docs - Requests

Each request method takes a Fiber app and path and returns a chainable `*Response`. These pages
cover request methods, default headers, timeouts, multipart uploads, and cookie jars for
multi-request flows.

## Contents

| Page                          | What it answers                                                   |
| ----------------------------- | ----------------------------------------------------------------- |
| [Methods](methods.md)         | Every request method, including the `*WithConfig` twins           |
| [Headers](./headers.md)       | Default headers, bearer tokens, and which header wins             |
| [Test Config](testconfigs.md) | `fiber.TestConfig`, timeouts, `ZENTESTS_TIMEOUT`, `FailOnTimeout` |
| [Multipart](multipart.md)     | `PostMultipart` and `FileUpload`                                  |
| [Cookie Jars](cookie-jars.md) | Keeping cookies across requests: login, CSRF, sessions            |

## What `zt` methods share

- **Signature.** `zt.Method(app, path, ...) *Response` — pass the `*fiber.App` and path, plus any
  method-specific input. Jar methods are separate and are described in
  [Cookie Jars](cookie-jars.md).
- **Return value.** A `*Response`; every assertion is a method on it, and every assertion is
  chainable.
- **Default headers.** `SetHeader` and `WithBearer` apply to requests from that helper. A method's
  own header, such as `PostJSON`'s `Content-Type: application/json`, overrides a same-name default.
- **Timeouts.** Plain methods use a zero-value config, which gets a one-second timeout and
  `FailOnTimeout: true`. Non-zero `*WithConfig` values are passed as-is. `ZENTESTS_TIMEOUT`
  overrides the default timeout for the process.
- **No listener.** Requests run through `app.Test`, so the real router, middleware and error handler
  execute in-process. There is no port to pick and no server to shut down.

## Choosing a method

| You want to                          | Use                                               |
| ------------------------------------ | ------------------------------------------------- |
| Send a JSON body                     | `PostJSON`, `PutJSON`, `PatchJSON`, `DeleteJSON`  |
| Send a raw body you encoded yourself | `Post`, `Put`, `Patch` (takes `[]byte`)           |
| Submit an HTML form                  | `PostForm` (takes `map[string]string`)            |
| Upload files                         | `PostMultipart` (takes fields and `[]FileUpload`) |
| Send query parameters                | `GetQuery` (takes `url.Values`)                   |
| Set a per-call timeout               | any method with a `*WithConfig` twin              |
| Keep cookies between calls           | `NewJar` — see [Cookie Jars](cookie-jars.md)      |

## Related

- [Assertions](../assertions/README.md) — what to do with the `*Response`
- [Testing styles](../testing/README.md) — where to create `zt`
- [examples/01-flat-style](../../examples/example-01.md) — requests and assertions in a runnable
  file

---

Back to: [Docs](../README.md) &mdash; Next: [Methods](methods.md)

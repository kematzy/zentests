# Zentests - Docs - Requests

> **Status:** DRAFT — structure and wording are open for discussion.

Every request is one call on `zt`, taking the app and the path and returning a `*Response` ready to
chain. The pages here cover the methods themselves, the headers sent with every call, timeouts, file
uploads and cookie-keeping request flows.

## Contents

| Page                          | What it answers                                                   |
| ----------------------------- | ----------------------------------------------------------------- |
| [Methods](methods.md)         | Every request method, including the `*WithConfig` twins           |
| [Headers](./headers.md)       | Default headers, bearer tokens, and which header wins             |
| [Test Config](testconfigs.md) | `fiber.TestConfig`, timeouts, `ZENTESTS_TIMEOUT`, `FailOnTimeout` |
| [Multipart](multipart.md)     | `PostMultipart` and `FileUpload`                                  |
| [Cookie Jars](cookie-jars.md) | Keeping cookies across requests: login, CSRF, sessions            |

## What every method shares

- **Signature.** `zt.Method(app, path, ...) *Response` — the `*fiber.App` under test and the path,
  nothing else required.
- **Return value.** A `*Response`; every assertion is a method on it, and every assertion is
  chainable.
- **Default headers.** Anything set with `SetHeader` or `WithBearer` is sent with the request. A
  header belonging to the single request wins over a default of the same name, so `PostJSON`'s
  `Content-Type: application/json` is never overwritten by a default.
- **Timeouts.** A zero-value `fiber.TestConfig{}` means a one-second timeout with
  `FailOnTimeout: true`. `ZENTESTS_TIMEOUT=5s` raises it for a whole run.
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

# Zentests - Docs - Helpers

> **Status:** DRAFT — structure and wording are open for discussion.

Small helpers for the two things every suite ends up writing by hand: an app that shuts itself down,
and a logger you can assert against.

## Contents

| Page                | What it answers                                                |
| ------------------- | -------------------------------------------------------------- |
| [App](app.md)       | `NewApp` — a Fiber app with automatic shutdown                 |
| [Logger](logger.md) | `Logger`, `DefaultLogger`, `MockLogger` — asserting log output |

Route assertions (`HasRoute`, `NotHasRoute`) are helpers too, but they belong with the other
assertions: see [assertions/routes](../assertions/routes.md).

## Related

- [Assertions](../assertions/README.md) — everything else that hangs off `*Response`
- [Requests](../requests/README.md) — sending the requests
- [examples/03-testify-suite](../../examples/example-03.md) — helpers in a suite

---

Back to: [Docs](../README.md) &mdash; Next: [App](app.md)

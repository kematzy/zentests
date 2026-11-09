# Zentests - Docs - Helpers

Helpers for managing a Fiber app's test cleanup and asserting on injected logger calls.

## Contents

| Page                | What it answers                                                |
| ------------------- | -------------------------------------------------------------- |
| [App](app.md)       | `NewApp` — a Fiber app with automatic shutdown                 |
| [Logger](logger.md) | `Logger`, `DefaultLogger`, `MockLogger` — asserting log output |

Route assertions (`HasRoute`, `NotHasRoute`) are helpers too, but they belong with the other
assertions: see [assertions/routes](../assertions/routes.md).

## Related

- [Assertions](../assertions/README.md) — response assertions and route checks
- [Requests](../requests/README.md) — sending the requests

---

Back to: [Docs](../README.md) &mdash; Next: [App](app.md)

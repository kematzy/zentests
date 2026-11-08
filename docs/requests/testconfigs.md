# Zentests - Docs - Requests - Test Config

> **Status:** DRAFT — structure and wording are open for discussion.

Every request runs through Fiber's in-process test transport, which takes a timeout and a policy for
what happens when the handler exceeds it. `zentests` applies sensible defaults and lets you override
them per call or for a whole run.

## Defaults

| Setting         | Default | Meaning                                               |
| --------------- | ------- | ----------------------------------------------------- |
| `Timeout`       | `1s`    | How long the request may take                         |
| `FailOnTimeout` | `true`  | Fail the test instead of returning a partial response |

A zero-value `fiber.TestConfig{}` means "use the defaults", so plain `zt.Get(...)` is already
guarded against a handler that never returns.

## Per call

The ten verb methods have a `*WithConfig` twin:

```go
zt.GetWithConfig(app, "/api/v1/report", fiber.TestConfig{
	Timeout: 10 * time.Second,
}).OK()
```

`GetWithConfig`, `PostWithConfig`, `PostJSONWithConfig`, `PostFormWithConfig`, `PutWithConfig`,
`PutJSONWithConfig`, `PatchWithConfig`, `PatchJSONWithConfig`, `DeleteWithConfig`,
`DeleteJSONWithConfig`.

`GetQuery` and `PostMultipart` have no config twin; they use the defaults below.

## Per run

`ZENTESTS_TIMEOUT` replaces the default timeout for the whole process, which is usually easier than
editing every call when a suite is slow on CI:

```bash
ZENTESTS_TIMEOUT=5s make test
```

Any Go duration works (`500ms`, `2s`, `1m`). An invalid value panics at startup rather than being
silently ignored — a typo would otherwise leave the suite running with the wrong timeout.

## `FailOnTimeout`

With `FailOnTimeout: false`, a request that times out returns whatever the handler had written
rather than failing. That is useful for streaming or long-poll endpoints where a partial response is
the expected outcome:

```go
resp := zt.GetWithConfig(app, "/stream", fiber.TestConfig{
	Timeout:       250 * time.Millisecond,
	FailOnTimeout: false,
})
resp.OK().Contains("first-chunk")
```

## Notes

- **The timeout applies to the whole request**, handler included, not to each write.
- **A per-call config beats the run-wide default**: `*WithConfig` ignores `ZENTESTS_TIMEOUT`.
- **Timeouts fail through `require`**, so a timed-out request stops the current test rather than
  continuing with an empty response.

## Related

- [Methods](methods.md) — which methods have a `*WithConfig` twin
- [Testing styles](../testing/README.md) — using `TestConfig` inside a suite
- [examples/01-flat-style](../../examples/example-01.md)

---

Back to: [Requests](./README.md) &mdash; Next: [Multipart](multipart.md)

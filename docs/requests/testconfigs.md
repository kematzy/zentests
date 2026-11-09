# Zentests - Docs - Requests - Test Config

Every request runs through Fiber's in-process test transport, which takes a timeout and a policy for
what happens when the handler exceeds it. `zentests` applies sensible defaults and lets you override
them per call or for a whole run.

## Defaults

| Setting         | Package default | Meaning                                                                      |
| --------------- | --------------- | ---------------------------------------------------------------------------- |
| `Timeout`       | `1s`            | Maximum duration for the request                                             |
| `FailOnTimeout` | `true`          | On timeout, stop the test (`true`) or use Fiber's partial response (`false`) |

Plain request methods pass a zero-value `fiber.TestConfig{}`, which `zentests` replaces with these
package defaults. A non-zero `*WithConfig` value is passed as-is rather than merged with the
defaults; set every field whose behavior matters. For example, setting only `Timeout` leaves
`FailOnTimeout` at `false`.

## Per call

The ten verb methods have a `*WithConfig` twin:

```go
zt.GetWithConfig(app, "/api/v1/report", fiber.TestConfig{
	Timeout:       10 * time.Second,
	FailOnTimeout: true,
}).OK()
```

All ten `*WithConfig` methods are listed in the
[Methods reference](methods.md#http-methods-with-testconfig).

`GetQuery` and `PostMultipart` have no config twin; they use the defaults below.

## Per run

Set `ZENTESTS_TIMEOUT` before starting the test process to replace the package's default timeout.
This is usually easier than editing every call when a suite is slow on CI:

```bash
ZENTESTS_TIMEOUT=5s make test
```

Use a value accepted by `time.ParseDuration` (`500ms`, `2s`, `1m`). An invalid value panics during
package initialization rather than being silently ignored.

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
- **A non-zero per-call config bypasses the run-wide default**, including `ZENTESTS_TIMEOUT`. A
  zero-value config still uses the package defaults.
- **With `FailOnTimeout: true`, a timeout stops the current test** through `require`; with `false`,
  Fiber may return a partial response.

## Related

- [Methods](methods.md) — which methods have a `*WithConfig` twin
- [Testing styles](../testing/README.md) — using `TestConfig` inside a suite
- [examples/01-flat-style](../../examples/example-01.md)

---

Back to: [Requests](./README.md) &mdash; Next: [Multipart](multipart.md)

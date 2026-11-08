# Zentests - Docs - Helpers - Logger

> **Status:** DRAFT — structure and wording are open for discussion.

A three-method logging interface, a production implementation, and a test double that makes fatal
paths testable — because `Fatalf` on the mock records the call and returns instead of killing the
process.

## The interface

```go
type Logger interface {
	Debugf(format string, v ...any)
	Warnf(format string, v ...any)
	Fatalf(format string, v ...any)
}
```

Declare it in your own code wherever logging must be injectable:

```go
type Service struct {
	log zentests.Logger
}

func NewService(log zentests.Logger) *Service {
	return &Service{log: log}
}
```

## `DefaultLogger`

Writes to stderr with level prefixes and exits on a fatal message.

```go
var logger zentests.Logger = zentests.DefaultLogger{}

logger.Debugf("cache hit for %s", key)   // DEBUG: cache hit for user:42
logger.Warnf("retrying in %d seconds", n) // WARN: retrying in 3 seconds
logger.Fatalf("cannot start: %v", err)    // FATAL: cannot start: … then os.Exit(1)
```

## `MockLogger`

Captures every call, prints nothing, and never exits.

```go
ml := &zentests.MockLogger{}
svc := NewService(ml)

svc.Start() // logs a warning internally

assert.True(t, ml.WarnCalled)
assert.Equal(t, "config not found", ml.WarnMsg)
```

| Field         | Type     | Set by                            |
| ------------- | -------- | --------------------------------- |
| `DebugCalled` | `bool`   | `Debugf`                          |
| `WarnCalled`  | `bool`   | `Warnf`                           |
| `FatalCalled` | `bool`   | `Fatalf`                          |
| `DebugMsg`    | `string` | `Debugf` — last formatted message |
| `WarnMsg`     | `string` | `Warnf` — last formatted message  |
| `FatalMsg`    | `string` | `Fatalf` — last formatted message |

Each `*Msg` field holds the **last** message for that level, not a history. When a test needs the
sequence, inject a small recorder of your own that satisfies the interface.

## Testing a fatal path

```go
ml := &zentests.MockLogger{}
svc := NewService(ml)

svc.MustLoad("missing.json") // would exit in production

assert.True(t, ml.FatalCalled)
assert.Contains(t, ml.FatalMsg, "missing.json")
```

That is the reason the mock exists: `Fatalf` on `DefaultLogger` ends the test process, which no test
can assert on.

## Notes

- **`*MockLogger` is a pointer type**, because the methods record into it.
- **No levels, no structured fields.** Three methods, no configuration — enough for injection, not a
  logging framework.
- **`DefaultLogger` writes to stderr**, so output appears alongside `go test` failures and does not
  mix with stdout captured by the test.

## Related

- [App](app.md) — the other helper
- [examples/03-testify-suite](../../examples/example-03.md) — helpers used together

---

Back to: [Helpers](README.md)

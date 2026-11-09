# Zentests - Docs - Helpers - Logger

A small logging interface, a standard-library implementation, and a test double for checking
logging calls. `MockLogger.Fatalf` records the call and returns, so code that logs fatally can be
tested without exiting the process.

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

Uses the standard library's default logger (stderr by default), prefixes messages with their level,
and exits with status 1 after `Fatalf`.

```go
var logger zentests.Logger = zentests.DefaultLogger{}

logger.Debugf("cache hit for %s", key)    // logs a DEBUG-prefixed message
logger.Warnf("retrying in %d seconds", n) // logs a WARN-prefixed message
logger.Fatalf("cannot start: %v", err)    // logs FATAL, then exits with status 1
```

## `MockLogger`

Records whether each level was called and its last formatted message. It prints nothing and never exits.

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

A direct call to `DefaultLogger.Fatalf` exits the process. Inject `MockLogger` to test code that calls
`Fatalf`; use a subprocess test if you need to verify the real process exit.

## Notes

- **`*MockLogger` is a pointer type**, because the methods record into it.
- **No levels, no structured fields.** Three methods, no configuration — enough for injection, not a
  logging framework.
- **`DefaultLogger` uses the standard logger**, which writes to stderr by default; its output is not
  captured as `testing.T` log output.

## Related

- [App](app.md) — the other helper

---

Back to: [Helpers](./README.md)

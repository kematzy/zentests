package zentests

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/gofiber/fiber/v3"
)

// recordingTB captures assertion failures instead of failing the real test.
// It embeds a real testing.TB so methods it does not override (Helper, Log, ...)
// keep working. FailNow stops the calling goroutine, like the real one does.
type recordingTB struct {
	testing.TB
	errors []string
	fatal  bool
}

// Errorf records a failure message.
func (r *recordingTB) Errorf(format string, args ...any) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}

// FailNow records a fatal failure and exits the calling goroutine.
func (r *recordingTB) FailNow() {
	r.fatal = true
	runtime.Goexit()
}

// runRecorded runs fn against a recordingTB in its own goroutine, so that
// FailNow (runtime.Goexit) only ends fn and not the test. It returns the recorder.
func runRecorded(t *testing.T, fn func(tb testing.TB)) *recordingTB {
	t.Helper()
	rec := &recordingTB{TB: t}
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn(rec)
	}()
	<-done
	return rec
}

// failuresOf runs assertion against the response of GET / served by handler and returns
// the failure messages it recorded. Used to check that assertions report wrong responses.
func failuresOf(t *testing.T, handler fiber.Handler, assertion func(*Response)) []string {
	t.Helper()
	app := fiber.New()
	app.Get("/", handler)

	rec := runRecorded(t, func(tb testing.TB) {
		assertion(execute(tb, app, "GET", "/", nil, nil, fiber.TestConfig{}))
	})
	return rec.errors
}

// sendString returns a handler that replies with body as text/plain.
func sendString(body string) fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.SendString(body)
	}
}

// sendRawJSON returns a handler that replies with body as application/json.
func sendRawJSON(body string) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString(body)
	}
}

package zentests

import (
	"fmt"
	"log"
	"os"
)

// Logger defines the minimal logging interface for use with zentests-compatible packages.
// It allows the host application to inject a logger for debug messages, warnings, and
// fatal errors. Implement this interface in your own code, then inject a [MockLogger]
// in tests to capture and assert on log output without any side effects.
//
// Example — production code:
//
//	type MyService struct {
//	    log zentests.Logger
//	}
//
// Example — test code:
//
//	ml := &zentests.MockLogger{}
//	svc := MyService{log: ml}
//	svc.DoSomething()
//	s.True(ml.WarnCalled)
type Logger interface {
	// Debugf logs a formatted debug message.
	Debugf(format string, v ...any)

	// Warnf logs a formatted warning message (non-fatal).
	Warnf(format string, v ...any)

	// Fatalf logs a formatted fatal error message and terminates the process.
	Fatalf(format string, v ...any)
}

// DefaultLogger is a production-ready [Logger] implementation backed by the standard
// library's log package. All output is written to stderr with level prefixes:
// "DEBUG: ", "WARN: ", or "FATAL: ".
//
// Fatalf calls os.Exit(1) after logging. The exit behaviour can be overridden via the
// unexported exitFn field, which is only used in tests for the DefaultLogger itself.
//
// Example:
//
//	var log zentests.Logger = zentests.DefaultLogger{}
//	log.Warnf("retrying in %d seconds", 5)
type DefaultLogger struct {
	exitFn func(int) // injectable for testing; nil means use os.Exit
}

func (l DefaultLogger) exit(code int) {
	if l.exitFn != nil {
		l.exitFn(code)
		return
	}
	os.Exit(code)
}

// Debugf logs a formatted debug message prefixed with "DEBUG: ".
func (l DefaultLogger) Debugf(format string, v ...any) {
	log.Printf("DEBUG: "+format, v...)
}

// Warnf logs a formatted warning message prefixed with "WARN: ".
func (l DefaultLogger) Warnf(format string, v ...any) {
	log.Printf("WARN: "+format, v...)
}

// Fatalf logs a formatted fatal message prefixed with "FATAL: " then exits with code 1.
// NOTE: uses Printf internally (not log.Fatalf) so the exit is handled by exitFn,
// making it safe to test without actually terminating the process.
func (l DefaultLogger) Fatalf(format string, v ...any) {
	log.Printf("FATAL: "+format, v...)
	l.exit(1)
}

// MockLogger is an exported test double that implements [Logger] by capturing all log
// calls without any side effects. Inject *MockLogger into code-under-test, then inspect
// the exported fields in assertions.
//
// Example:
//
//	func (s *MySuite) Test_Warns_On_Missing_Config() {
//	    ml := &zentests.MockLogger{}
//	    svc := NewMyService(ml)
//	    svc.Start() // triggers Warnf internally
//	    s.True(ml.WarnCalled)
//	    s.Equal("config not found", ml.WarnMsg)
//	}
//
// All three methods (Debugf, Warnf, Fatalf) set their respective *Called and *Msg fields.
// Fatalf does NOT call os.Exit — safe to use in tests.
type MockLogger struct {
	// DebugCalled is true if Debugf was called at least once.
	DebugCalled bool
	// WarnCalled is true if Warnf was called at least once.
	WarnCalled bool
	// FatalCalled is true if Fatalf was called at least once.
	FatalCalled bool

	// DebugMsg holds the last formatted message passed to Debugf.
	DebugMsg string
	// WarnMsg holds the last formatted message passed to Warnf.
	WarnMsg string
	// FatalMsg holds the last formatted message passed to Fatalf.
	FatalMsg string
}

// Debugf records the call and formatted message. Does not produce any output.
func (m *MockLogger) Debugf(format string, v ...any) {
	m.DebugCalled = true
	m.DebugMsg = fmt.Sprintf(format, v...)
}

// Warnf records the call and formatted message. Does not produce any output.
func (m *MockLogger) Warnf(format string, v ...any) {
	m.WarnCalled = true
	m.WarnMsg = fmt.Sprintf(format, v...)
}

// Fatalf records the call and formatted message. Does NOT call os.Exit — safe in tests.
func (m *MockLogger) Fatalf(format string, v ...any) {
	m.FatalCalled = true
	m.FatalMsg = fmt.Sprintf(format, v...)
}

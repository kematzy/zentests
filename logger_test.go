package zentests

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

// --- Logger interface compile-time guard -------------------------------------------------------

type LoggerTestSuite struct {
	suite.Suite
}

func (s *LoggerTestSuite) Test_DefaultLogger_CompileTimeGuard() {
	logger := DefaultLogger{}

	// compile-time guard: DefaultLogger must satisfy Logger
	var _ Logger = logger

	s.NotNil(logger)
}

func (s *LoggerTestSuite) Test_MockLogger_CompileTimeGuard() {
	logger := &MockLogger{}

	// compile-time guard: *MockLogger must satisfy Logger
	var _ Logger = logger

	s.NotNil(logger)
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

// --- MockLogger --------------------------------------------------------------------------------

type MockLoggerTestSuite struct {
	suite.Suite
	logger     Logger
	mockLogger *MockLogger
}

func (s *MockLoggerTestSuite) SetupTest() {
	s.mockLogger = &MockLogger{}
	s.logger = s.mockLogger
}

func (s *MockLoggerTestSuite) Test_MockLogger_InitialState() {
	s.False(s.mockLogger.DebugCalled)
	s.False(s.mockLogger.WarnCalled)
	s.False(s.mockLogger.FatalCalled)
	s.Empty(s.mockLogger.DebugMsg)
	s.Empty(s.mockLogger.WarnMsg)
	s.Empty(s.mockLogger.FatalMsg)
}

// Debugf

func (s *MockLoggerTestSuite) Test_MockLogger_Debugf_WithInt() {
	s.logger.Debugf("debug: %d", 42)
	s.True(s.mockLogger.DebugCalled)
	s.Equal("debug: 42", s.mockLogger.DebugMsg)
}

func (s *MockLoggerTestSuite) Test_MockLogger_Debugf_WithString() {
	s.logger.Debugf("debug: %s", "This works!")
	s.True(s.mockLogger.DebugCalled)
	s.Equal("debug: This works!", s.mockLogger.DebugMsg)
}

func (s *MockLoggerTestSuite) Test_MockLogger_Debugf_WithArray() {
	s.logger.Debugf("debug: %v", []string{"a", "b"})
	s.True(s.mockLogger.DebugCalled)
	s.Equal("debug: [a b]", s.mockLogger.DebugMsg)
}

func (s *MockLoggerTestSuite) Test_MockLogger_Debugf_DoesNotAffectOtherFields() {
	s.logger.Debugf("hello")
	s.False(s.mockLogger.WarnCalled)
	s.False(s.mockLogger.FatalCalled)
}

// Warnf

func (s *MockLoggerTestSuite) Test_MockLogger_Warnf_WithInt() {
	s.logger.Warnf("warn: %d", 42)
	s.True(s.mockLogger.WarnCalled)
	s.Equal("warn: 42", s.mockLogger.WarnMsg)
}

func (s *MockLoggerTestSuite) Test_MockLogger_Warnf_WithString() {
	s.logger.Warnf("warn: %s", "This works!")
	s.True(s.mockLogger.WarnCalled)
	s.Equal("warn: This works!", s.mockLogger.WarnMsg)
}

func (s *MockLoggerTestSuite) Test_MockLogger_Warnf_WithArray() {
	s.logger.Warnf("warn: %v", []string{"a", "b"})
	s.True(s.mockLogger.WarnCalled)
	s.Equal("warn: [a b]", s.mockLogger.WarnMsg)
}

func (s *MockLoggerTestSuite) Test_MockLogger_Warnf_DoesNotAffectOtherFields() {
	s.logger.Warnf("hello")
	s.False(s.mockLogger.DebugCalled)
	s.False(s.mockLogger.FatalCalled)
}

// Fatalf

func (s *MockLoggerTestSuite) Test_MockLogger_Fatalf_WithInt() {
	s.logger.Fatalf("fatal: %d", 42)
	s.True(s.mockLogger.FatalCalled)
	s.Equal("fatal: 42", s.mockLogger.FatalMsg)
}

func (s *MockLoggerTestSuite) Test_MockLogger_Fatalf_WithString() {
	s.logger.Fatalf("fatal: %s", "This works!")
	s.True(s.mockLogger.FatalCalled)
	s.Equal("fatal: This works!", s.mockLogger.FatalMsg)
}

func (s *MockLoggerTestSuite) Test_MockLogger_Fatalf_WithArray() {
	s.logger.Fatalf("fatal: %v", []string{"a", "b"})
	s.True(s.mockLogger.FatalCalled)
	s.Equal("fatal: [a b]", s.mockLogger.FatalMsg)
}

func (s *MockLoggerTestSuite) Test_MockLogger_Fatalf_DoesNotExit() {
	// Fatalf must not call os.Exit — test would be killed if it did.
	s.logger.Fatalf("critical: %s", "should not exit")
	s.True(s.mockLogger.FatalCalled) // still reachable
}

func (s *MockLoggerTestSuite) Test_MockLogger_Fatalf_DoesNotAffectOtherFields() {
	s.logger.Fatalf("hello")
	s.False(s.mockLogger.DebugCalled)
	s.False(s.mockLogger.WarnCalled)
}

// Multiple calls — last message wins

func (s *MockLoggerTestSuite) Test_MockLogger_Debugf_LastMessageWins() {
	s.logger.Debugf("first")
	s.logger.Debugf("second")
	s.Equal("second", s.mockLogger.DebugMsg)
}

func TestMockLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(MockLoggerTestSuite))
}

// --- DefaultLogger -----------------------------------------------------------------------------

type DefaultLoggerTestSuite struct {
	suite.Suite
	logger DefaultLogger
	buf    bytes.Buffer
}

func (s *DefaultLoggerTestSuite) SetupTest() {
	s.buf.Reset()
	log.SetOutput(&s.buf)
	log.SetFlags(0) // remove timestamp noise from assertions
	s.logger = DefaultLogger{}
}

func (s *DefaultLoggerTestSuite) TearDownTest() {
	log.SetOutput(os.Stderr) // restore after each test
	log.SetFlags(log.LstdFlags)
}

func (s *DefaultLoggerTestSuite) Test_Debugf_FormatsCorrectly() {
	s.logger.Debugf("user %s logged in", "mats")
	s.Contains(s.buf.String(), "DEBUG: user mats logged in")
}

func (s *DefaultLoggerTestSuite) Test_Debugf_WithInt() {
	s.logger.Debugf("count: %d", 7)
	s.Contains(s.buf.String(), "DEBUG: count: 7")
}

func (s *DefaultLoggerTestSuite) Test_Warnf_FormatsCorrectly() {
	s.logger.Warnf("token expiring in %d seconds", 30)
	s.Contains(s.buf.String(), "WARN: token expiring in 30 seconds")
}

func (s *DefaultLoggerTestSuite) Test_Warnf_WithString() {
	s.logger.Warnf("deprecated: %s", "use NewFunc instead")
	s.Contains(s.buf.String(), "WARN: deprecated: use NewFunc instead")
}

func (s *DefaultLoggerTestSuite) Test_Fatalf_FormatsCorrectly() {
	var exitCode int
	s.logger = DefaultLogger{
		exitFn: func(code int) { exitCode = code },
	}

	s.logger.Fatalf("critical failure: %s", "db unreachable")

	s.Contains(s.buf.String(), "FATAL: critical failure: db unreachable")
	s.Equal(1, exitCode)
}

func (s *DefaultLoggerTestSuite) Test_Fatalf_ExitsWithCode1() {
	var exitCode int
	s.logger = DefaultLogger{
		exitFn: func(code int) { exitCode = code },
	}

	s.logger.Fatalf("boom")

	s.Equal(1, exitCode)
}

func (s *DefaultLoggerTestSuite) Test_Fatalf_DefaultExitFn_IsNil() {
	// Verify zero-value DefaultLogger has nil exitFn (uses os.Exit in production).
	// We can't call it without exiting, but we can confirm the field is nil.
	logger := DefaultLogger{}
	s.Nil(logger.exitFn)
}

func (s *DefaultLoggerTestSuite) Test_DefaultLogger_SatisfiesLogger_Interface() {
	var _ Logger = DefaultLogger{}
	s.True(true) // compile-time check; reaching here means it passed
}

func TestDefaultLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(DefaultLoggerTestSuite))
}

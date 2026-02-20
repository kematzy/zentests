package zentests

import (
	"testing"
)

// SuiteContext holds hooks and state for a BDD-style describe block.
// Provides lifecycle hooks (BeforeEach, AfterEach) for test setup and teardown.
// Each describe block gets its own SuiteContext to manage test isolation.
//
// Fields:
//   - t: The testing.T instance for running subtests
//   - beforeEach: Function to run before each test case in the suite
//   - afterEach: Function to run after each test case in the suite
type SuiteContext struct {
	t          *testing.T
	beforeEach func(*T)
	afterEach  func(*T)
}

// Describe creates a BDD-style test group with lifecycle hooks.
// Organizes related tests under a descriptive name with optional setup/teardown.
// Similar to describe/it blocks in testing frameworks like RSpec or Mocha.
//
// Parameters:
//   - name: The descriptive name for the test group (e.g., "User API", "Auth Service")
//   - fn: Function that receives the SuiteContext and defines test cases using It()
//
// Example:
//
//	zt.Describe("User API", func(ctx *zentests.SuiteContext) {
//	    ctx.BeforeEach(func(zt *zentests.T) {
//	        // Setup fresh app before each test
//	        app = setupTestApp()
//	    })
//
//	    ctx.It("creates a user", func(zt *zentests.T) {
//	        zt.PostJSON(app, "/users", data).Created()
//	    })
//	})
func (zt *T) Describe(name string, fn func(ctx *SuiteContext)) {
	zt.Run(name, func(t *testing.T) {
		ctx := &SuiteContext{t: t}
		fn(ctx)
	})
}

// It runs a test case within a describe block.
// Individual test cases that execute with BeforeEach/AfterEach hooks.
// Each It() call creates a sub-test that runs independently.
//
// Parameters:
//   - name: The descriptive name for the test case (e.g., "returns user list")
//   - fn: The test function receiving a *T for making requests
//
// Example:
//
//	ctx.It("returns user list", func(zt *zentests.T) {
//	    zt.Get(app, "/users").OK().IsJSON()
//	})
func (ctx *SuiteContext) It(name string, fn func(*T)) {
	ctx.t.Run(name, func(t *testing.T) {
		zt := &T{T: t}

		// Run beforeEach if defined
		if ctx.beforeEach != nil {
			ctx.beforeEach(zt)
		}

		// Run test
		fn(zt)

		// Run afterEach if defined
		if ctx.afterEach != nil {
			ctx.afterEach(zt)
		}
	})
}

// BeforeEach sets the setup hook function to run before each test case.
// Use this for common setup like creating fresh app instances or test data.
// Returns the SuiteContext for method chaining with AfterEach.
//
// Parameters:
//   - fn: Function to execute before each test case, receives *T
//
// Returns:
//   - *SuiteContext: The receiver for method chaining
//
// Example:
//
//	ctx.BeforeEach(func(zt *zentests.T) {
//	    app = fiber.New()
//	    SetupRoutes(app)
//	}).AfterEach(func(zt *zentests.T) {
//	    app.Shutdown()
//	})
func (ctx *SuiteContext) BeforeEach(fn func(*T)) *SuiteContext {
	ctx.beforeEach = fn
	return ctx
}

// AfterEach sets the teardown hook function to run after each test case.
// Use this for cleanup like closing database connections or shutting down apps.
// Returns the SuiteContext for method chaining.
//
// Parameters:
//   - fn: Function to execute after each test case, receives *T
//
// Returns:
//   - *SuiteContext: The receiver for method chaining
//
// Example:
//
//	ctx.AfterEach(func(zt *zentests.T) {
//	    app.Shutdown()
//	    db.Cleanup()
//	})
func (ctx *SuiteContext) AfterEach(fn func(*T)) *SuiteContext {
	ctx.afterEach = fn
	return ctx
}

package zentests

import (
	"testing"
)

// SuiteContext holds hooks and state for a BDD-style describe block.
// Provides lifecycle hooks (BeforeEach, AfterEach) for test setup and teardown.
// Each describe block gets its own SuiteContext to manage test isolation.
//
// Describe blocks can be nested. Hooks of enclosing blocks apply to tests in nested blocks:
// BeforeEach hooks run outermost first, and AfterEach hooks run innermost first.
//
// Hooks are captured when It() is called, so BeforeEach and AfterEach must be declared
// before the It() calls they should apply to.
//
// Fields:
//   - t: The [testing.T] instance for running subtests
//   - parent: The enclosing SuiteContext, or nil for a top-level block
//   - beforeEach: Function to run before each test case in the suite
//   - afterEach: Function to run after each test case in the suite
type SuiteContext struct {
	t          *testing.T
	parent     *SuiteContext
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
		t.Helper()
		fn(&SuiteContext{t: t})
	})
}

// Describe creates a nested describe block inside this one.
// Hooks declared in the enclosing block apply to the tests in the nested block.
//
// Parameters:
//   - name: The descriptive name for the nested group
//   - fn: Function that receives the nested SuiteContext
//
// Example:
//
//	zt.Describe("Users", func(ctx *zentests.SuiteContext) {
//	    ctx.BeforeEach(func(zt *zentests.T) { app = setupTestApp() })
//
//	    ctx.Describe("when logged in", func(ctx *zentests.SuiteContext) {
//	        ctx.It("lists users", func(zt *zentests.T) {
//	            zt.Get(app, "/users").OK()
//	        })
//	    })
//	})
func (ctx *SuiteContext) Describe(name string, fn func(ctx *SuiteContext)) {
	ctx.t.Run(name, func(t *testing.T) {
		t.Helper()
		fn(&SuiteContext{t: t, parent: ctx})
	})
}

// It runs a test case within a describe block.
// Runs the BeforeEach hooks of this block and its enclosing blocks, then fn, then the
// AfterEach hooks. AfterEach hooks are registered with t.Cleanup, so they still run when
// the test stops early (for example through require).
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
	before, after := ctx.hooks()

	ctx.t.Run(name, func(t *testing.T) {
		t.Helper()
		zt := &T{T: t}

		// Register cleanups before setup so they still run if a BeforeEach hook stops the test.
		// Cleanups run last-in first-out, so registering outermost first runs innermost first.
		for _, hook := range after {
			t.Cleanup(func() {
				t.Helper()
				hook(zt)
			})
		}
		for _, hook := range before {
			hook(zt)
		}

		fn(zt)
	})
}

// hooks collects the BeforeEach and AfterEach hooks that apply to this block, including
// those of enclosing blocks. Both slices are ordered outermost first.
func (ctx *SuiteContext) hooks() (before, after []func(*T)) {
	for block := ctx; block != nil; block = block.parent {
		if block.beforeEach != nil {
			before = append([]func(*T){block.beforeEach}, before...)
		}
		if block.afterEach != nil {
			after = append([]func(*T){block.afterEach}, after...)
		}
	}
	return before, after
}

// BeforeEach sets the setup hook function to run before each test case.
// Use this for common setup like creating fresh app instances or test data.
// Must be called before the It() calls it should apply to.
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
// Must be called before the It() calls it should apply to.
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

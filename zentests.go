// Package zentests provides utilities for testing Go Fiber applications.
//
// Basic usage:
//
//	func TestEndpoint(t *testing.T) {
//	    zt := zentests.New(t)
//	    zt.Get(app, "/users").Status(200)
//	}
package zentests

import (
	"testing"

	"github.com/gofiber/fiber/v3"
)

// T wraps [testing.T] with additional context for testing Fiber applications.
// It provides a fluent API for making HTTP requests against a Fiber app in tests.
// Request methods take the Fiber app explicitly; T also carries default headers
// and the testing context for chainable assertions.
//
// Basic usage:
//
//	func TestEndpoint(t *testing.T) {
//	    zt := zentests.New(t)
//	    zt.Get(app, "/users").Status(200)
//	}
//
// Fields:
//   - T: The embedded [testing.T] instance for standard Go test functionality
//   - headers: Default headers sent with every request (see SetHeader)
type T struct {
	*testing.T
	app     *fiber.App
	headers map[string]string
}

// New creates a new zentests context wrapping the provided [testing.T].
// This is the main entry point for using zentests. Pass the Fiber app explicitly to
// request methods such as Get(app, "/").
//
// Parameters:
//   - t: The [testing.T] instance from your test function
//
// Returns:
//   - *T: A new zentests context with no Fiber app configured
//
// Example:
//
//	func TestExample(t *testing.T) {
//	    zt := zentests.New(t)
//	    app := fiber.New()
//	    // declare app routes
//
//	    zt.Get(app, "/").OK()
//	}
func New(t *testing.T) *T {
	return &T{T: t}
}

// Use is retained for source compatibility and checks that app is not nil.
// Request methods still require app to be passed explicitly.
//
// Deprecated: Use does not bind app to request methods. Pass app to each request method.
//
// Parameters:
//   - app: The Fiber application instance; must not be nil
//
// Returns:
//   - *T: The receiver for method chaining
func (zt *T) Use(app *fiber.App) *T {
	zt.Helper()
	if app == nil {
		zt.Fatal("zentests: Use called with a nil *fiber.App")
	}
	zt.app = app
	return zt
}

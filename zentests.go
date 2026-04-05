// Package zentests provides utilities for testing Go Fiber applications.
//
// Basic usage:
//
//	func TestEndpoint(t *testing.T) {
//	    zt := zentests.New(t).Use(app)
//	    zt.Get("/users").ExpectStatus(200)
//	}
package zentests

import (
	"testing"

	"github.com/gofiber/fiber/v3"
)

// T wraps testing.T with additional context for testing Fiber applications.
// It provides a fluent API for making HTTP requests against a Fiber app in tests.
// T maintains a reference to the Fiber app and testing context to enable chainable assertions.
//
// Basic usage:
//
//	func TestEndpoint(t *testing.T) {
//	    zt := zentests.New(t).Use(app)
//	    zt.Get("/users").ExpectStatus(200)
//	}
//
// Fields:
//   - T: The embedded testing.T instance for standard Go test functionality
//   - app: The Fiber application instance for making HTTP requests
type T struct {
	*testing.T
	app *fiber.App
}

// New creates a new zentests context wrapping the provided testing.T.
// This is the main entry point for using zentests. The returned *T has no app set initially
// - call Use() to set the Fiber application for method chaining, or pass the app directly
// to request methods like Get(app, "/").
//
// Parameters:
//   - t: The testing.T instance from your test function
//
// Returns:
//   - *T: A new zentests context with no Fiber app configured
//
// Example:
//
//	func TestExample(t *testing.T) {
//		zt := zentests.New(t)
//		app := fiber.New()
//		// declare app routes
//		zt.Use(app)
//	}
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

// Use sets the Fiber app for making test requests and enables fluent method chaining.
// After calling Use(), subsequent request methods can be called without passing the app parameter.
//
// Parameters:
//   - app: The Fiber application instance to test against
//
// Returns:
//
//   - *T: The receiver for method chaining
//
// Example:
//
//     zt := zentests.New(t).Use(app)
//     zt.Get("/users").OK()

func (zt *T) Use(app *fiber.App) *T {
	if app == nil {
		zt.Fatal("Use called with nil app")
	}
	zt.app = app
	return zt
}

package zentests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

// HasRoute asserts that a route with the given HTTP method and path is registered
// in the Fiber app. The method is case-insensitive. Returns the receiver for chaining.
//
// This is useful in tests that verify route registration without making an actual
// HTTP request — for example, when testing that a router setup function wires up
// all expected endpoints.
//
// Parameters:
//   - method: HTTP method to check (e.g. "GET", "POST"); case-insensitive
//   - path: The exact route path as registered (e.g. "/users/:id")
//   - app: The Fiber application to inspect
//   - msgAndArgs: Optional failure message and format arguments (testify-style)
//
// Returns:
//   - *T: The receiver for method chaining
//
// Example:
//
//	zt := zentests.New(t)
//	zt.HasRoute("GET", "/users", app, "user list route must be registered")
//	zt.HasRoute("POST", "/users", app).
//	    HasRoute("DELETE", "/users/:id", app)
func (zt *T) HasRoute(method, path string, app *fiber.App, msgAndArgs ...any) *T {
	zt.T.Helper()

	found := hasRoute(app, method, path)
	assert.True(zt.T, found, buildRouteMsg(method, path, msgAndArgs...))

	return zt
}

// NotHasRoute asserts that no route with the given HTTP method and path is registered
// in the Fiber app. The complement of HasRoute. Returns the receiver for chaining.
//
// Parameters:
//   - method: HTTP method to check; case-insensitive
//   - path: The route path to verify is absent
//   - app: The Fiber application to inspect
//   - msgAndArgs: Optional failure message and format arguments
//
// Returns:
//   - *T: The receiver for method chaining
//
// Example:
//
//	zt.NotHasRoute("DELETE", "/users", app, "bulk delete must not be exposed").
//	    NotHasRoute("PATCH", "/users", app)
func (zt *T) NotHasRoute(method, path string, app *fiber.App, msgAndArgs ...any) *T {
	zt.T.Helper()

	found := hasRoute(app, method, path)
	assert.False(zt.T, found, buildRouteMsg(method, path, msgAndArgs...))

	return zt
}

// NewApp creates a Fiber application for testing with automatic shutdown registered
// via t.Cleanup, so the app is always stopped when the test finishes regardless of
// pass or fail. The app is created fresh on each call — no shared state between tests.
//
// An optional [fiber.Config] can be passed to customise the app (views engine,
// app name, middleware, error handler, etc.). When no config is supplied the app is
// created with Fiber's defaults, identical to fiber.New().
//
// Parameters:
//   - t: The testing.T instance; marks this as a test helper and registers cleanup.
//   - cfg: Optional Fiber configuration. Only the first value is used if multiple are passed.
//
// Returns:
//   - *fiber.App: A fresh Fiber app ready for route registration.
//
// Example — zero config:
//
//	func (s *MySuite) SetupTest() {
//	    s.app = zentests.NewApp(s.T())
//	    s.app.Get("/ping", handler)
//	}
//
// Example — with views engine:
//
//		engine := html.New("./views", ".html")
//	 // register view helpers
//
//		app := zentests.NewApp(t, fiber.Config{
//		    Views:             engine,
//		    AppName:           "MyApp",
//		    PassLocalsToViews: true,
//		})
func NewApp(t *testing.T, cfg ...fiber.Config) *fiber.App {
	t.Helper()

	var app *fiber.App
	if len(cfg) > 0 {
		app = fiber.New(cfg[0])
	} else {
		app = fiber.New()
	}

	t.Cleanup(func() {
		_ = app.Shutdown()
	})

	return app
}

// hasRoute is the internal implementation used by HasRoute and NotHasRoute.
// It scans the Fiber app's route stack to check whether a route with the given method and path
// is registered. Method comparison is case-insensitive.
//
// In Fiber v3 the stack is indexed by HTTP method integer. Each Route carries a Method field
// (uppercase string) and a Path field (original registered path). Middleware routes registered
// with Use() appear in every method bucket with an empty Method field; they are excluded from
// this check since they are not bound to a specific HTTP verb.
func hasRoute(app *fiber.App, method, path string) bool {
	upper := strings.ToUpper(method)

	for _, routeGroup := range app.Stack() {
		for _, route := range routeGroup {
			if route.Method == upper && route.Path == path {
				return true
			}
		}
	}

	return false
}

// buildRouteMsg constructs a human-readable failure message for route assertions.
// When msgAndArgs is provided its first element is appended to the base label.
// String values are used directly; any other type is formatted with fmt.Sprint.
func buildRouteMsg(method, path string, msgAndArgs ...any) string {
	base := strings.ToUpper(method) + " " + path

	if len(msgAndArgs) == 0 {
		return base
	}

	if msg, ok := msgAndArgs[0].(string); ok {
		return base + " - " + msg
	}

	return base + " - " + fmt.Sprint(msgAndArgs[0])
}

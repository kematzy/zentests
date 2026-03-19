//go:build manual

// Package routes contains example tests demonstrating zentests usage.
// These examples show the flat/functional style of testing without suites.
// For suite-based examples, see suite.go.
// For basic route testing examples, see routes.go.
//
// Main source files:
//   - zentests.go: Core T type and New() constructor
//   - request.go: HTTP methods (Get, PostJSON, etc.)
//   - response.go: Response struct and basic methods
//   - response_http.go: HTTP assertions (OK, HasHeader, etc.)
//   - response_json.go: JSON assertions (Has, HasInt, etc.)
package routes

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/kematzy/zentests"
)

// User represents a user entity for testing.
// Used in the POST endpoint to demonstrate JSON body parsing.
type User struct {
	ID   int
	name string
}

// setupApp creates and configures a Fiber app with test routes.
// Returns a fully configured app ready for testing.
// Demonstrates: GET health check endpoint, POST user creation endpoint.
func setupApp() *fiber.App {
	app := fiber.New()

	// Health check endpoint - returns simple status
	app.Get("/api/health", func(c fiber.Ctx) error {
		return c.JSON(map[string]string{"status": "ok"})
	})

	// User creation endpoint - parses JSON body and returns created user with ID
	app.Post("/api/users", func(c fiber.Ctx) error {
		var user User
		if err := c.Bind().Body(&user); err != nil {
			return err
		}
		user.ID = 1
		return c.Status(fiber.StatusCreated).JSON(user)
	})

	return app
}

// TestAPIFlat demonstrates flat-style testing without suites.
// Shows simple, sequential test execution with method chaining.
//
// Features demonstrated:
//   - Creating zentests context with New()
//   - GET request with OK() and IsJSON() assertions
//   - POST JSON request with Created() assertion
//   - JSON field assertions with Has()
//   - Note: JSON numbers are float64, not int
func TestAPIFlat(t *testing.T) {
	// Initialize zentests context
	zt := zentests.New(t)
	app := setupApp()

	userData := User{
		name: "John Doe",
	}

	// Simple GET test: chain assertions for status, content type, and JSON field
	zt.Get(app, "/api/health").OK().IsJSON().Has("status", "ok")

	// POST JSON test: send data, expect 201 Created, verify response fields
	// Note: int values in JSON are float64 (e.g., float64(1))
	zt.PostJSON(app, "/api/users", userData).
		Created().
		Has("id", float64(1)). // NOTE! int in JSON are float64
		Has("name", "John Doe")
}

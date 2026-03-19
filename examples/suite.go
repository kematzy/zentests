//go:build manual

// Package routes contains example tests demonstrating zentests usage.
// This file demonstrates BDD-style testing with Describe/It blocks and lifecycle hooks.
// For simple flat-style examples, see flat_style.go and routes.go.
//
// Main source files:
//   - zentests.go: Core T type
//   - suite.go: SuiteContext, Describe, It, BeforeEach, AfterEach
//   - request.go: HTTP methods
//   - response_json.go: JSON assertions
package routes

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/kematzy/zentests"
)

// SetupRoutes configures a Fiber app with test routes.
// Separated from test for reusability across multiple test cases.
// Returns the configured app instance.
//
// Routes defined:
//   - GET /simple: Returns simple user data
//   - GET /api/users: Returns list of users
//   - GET /api/users/:id: Returns 404 for missing users
//   - POST /api/users: Creates new user
func SetupRoutes(app *fiber.App) *fiber.App {
	// Simple endpoint returning basic user data
	app.Get("/simple", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":   "John",
			"age":    30,
			"score":  95.5,
			"active": true,
		})
	})

	// Users list endpoint with nested data structure
	app.Get("/api/users", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data": []interface{}{
				fiber.Map{"id": 1, "email": "test@example.com"},
				fiber.Map{"id": 2, "email": "test2@example.com"},
			},
		})
	})

	// Single user endpoint that returns 404 for any ID
	app.Get("/api/users/:id", func(c fiber.Ctx) error {
		c.SendStatus(404)
		return c.JSON(fiber.Map{
			"success": false,
			"error":   "user not found",
			"data":    nil,
		})
	})

	// User creation endpoint
	app.Post("/api/users", func(c *fiber.Ctx) error {
		var data map[string]interface{}
		c.BodyParser(&data)
		return c.JSON(fiber.Map{
			"success": true,
			"data":    data,
		})
	})

	return app
}

// TestUserAPI demonstrates BDD-style testing with lifecycle hooks.
// Uses Describe/It blocks similar to RSpec, Mocha, or Pest.
//
// Features demonstrated:
//   - Describe blocks for test organization
//   - BeforeEach/AfterEach hooks for setup/teardown
//   - Method chaining on hooks (BeforeEach().AfterEach())
//   - Fresh app instance per test for isolation
//   - Multiple It() blocks within one Describe
//   - Various assertion types (OK, Created, NotFound, etc.)
func TestUserAPI(t *testing.T) {
	zt := zentests.New(t)
	var app *fiber.App

	// Describe creates a test group with lifecycle hooks
	zt.Describe("User API", func(ctx *zentests.SuiteContext) {
		// Chain BeforeEach and AfterEach for setup and cleanup
		ctx.
			BeforeEach(func(zt *zentests.T) {
				// Fresh app for each test ensures test isolation
				app = fiber.New()
				SetupRoutes(app)
			}).
			AfterEach(func(zt *zentests.T) {
				// Cleanup resources after each test
				app.Shutdown()
			})

		// Individual test case: verify user list endpoint
		ctx.It("returns user list", func(zt *zentests.T) {
			zt.Get(app, "/api/users").
				OK().
				IsJSON().
				Has("success", true).
				ArrayLength("data.users", 2)
		})

		// Individual test case: verify user creation
		ctx.It("creates a new user", func(zt *zentests.T) {
			zt.PostJSON(app, "/api/users", map[string]any{
				"name":  "John Doe",
				"email": "john@example.com",
			}).
				Created().
				Has("data.name", "John Doe").
				HasString("data.email", "john@example.com")

			// The commented code below shows an alternative flat-style approach
			// for comparison with the BDD style used above
		})

		// Individual test case: verify 404 handling
		ctx.It("returns 404 for missing user", func(zt *zentests.T) {
			zt.Get(app, "/api/users/999").
				NotFound().
				Has("success", false)
		})
	})
}

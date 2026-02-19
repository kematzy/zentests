//go:build manual

// Package routes contains example tests demonstrating zentests usage.
// This file demonstrates basic route testing with various assertion types.
// For flat-style examples, see flat_style.go.
// For suite-based examples with hooks, see suite.go.
//
// Main source files:
//   - zentests.go: Core T type and New() constructor
//   - request.go: HTTP methods (Get, PostJSON, etc.)
//   - response_json.go: JSON assertions with dot notation
package routes

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kematzy/zentests"
)

// TestAPI demonstrates comprehensive JSON response assertions.
// Shows dot notation for nested objects and array access.
//
// Features demonstrated:
//   - Inline route setup within test
//   - Status code assertions (OK())
//   - Content type assertions (IsJSON())
//   - Boolean JSON assertions (Has())
//   - Numeric JSON assertions (HasInt())
//   - Array length assertions (ArrayLength())
//   - Dot notation: "data.users.0.name" accesses nested array items
func TestAPI(t *testing.T) {
	zt := zentests.New(t)
	app := fiber.New()

	// Setup route returning nested JSON structure with array
	app.Get("/api/users", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"users": []fiber.Map{
					{"id": 1, "name": "John"},
					{"id": 2, "name": "Jane"},
				},
				"count": 2,
			},
		})
	})

	// Chain multiple assertions to validate response structure
	zt.Get(app, "/api/users").
		OK().                            // Assert 200 status
		IsJSON().                        // Assert JSON content type
		Has("success", true).            // Assert boolean field
		Has("data.count", 2).            // Assert numeric field (as float64)
		HasInt("data.count", 2).         // Assert numeric field (as int)
		ArrayLength("data.users", 2).    // Assert array has 2 items
		Has("data.users.0.name", "John") // Dot notation for array index access
}

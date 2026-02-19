package zentests

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/stretchr/testify/assert"
)

// JSON parses and caches the response body as JSON.
// Enables JSON-specific assertions on the response. The parsed JSON is cached
// to avoid re-parsing on subsequent JSON assertions. Fails the test if JSON
// parsing fails.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users").JSON().Has("data.0.name", "John")
//
// Implementation:.
func (r *Response) JSON() *Response {
	if r.parsedJSON == nil {
		err := json.Unmarshal(r.Body(), &r.parsedJSON)
		assert.NoError(r.t, err, "JSON parsing failed")
	}
	return r
}

// HasKey asserts that the JSON response contains the specified key path.
// Supports dot notation for nested keys and array indices (e.g., "users.0.name").
//
// Parameters:
//   - path: The JSON key path using dot notation (e.g., "user.email", "items.0.id")
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users/1").HasKey("data.user.name")
//
// Implementation:.
func (r *Response) HasKey(path string) *Response {
	r.JSON()
	_, exists := getNestedValue(r.parsedJSON, path)
	assert.True(r.t, exists, "JSON should have key %q", path)
	return r
}

// Has asserts that the JSON key equals the expected value with strict type checking.
// Supports dot notation for nested keys. Fails if the key doesn't exist or if
// the type doesn't match exactly.
//
// Parameters:
//   - path: The JSON key path using dot notation
//   - expected: The expected value (type must match exactly)
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users/1").Has("data.user.name", "John")
//	zt.Get(app, "/api/count").Has("data.count", float64(42)) // JSON numbers are float64
//
// Implementation:.
func (r *Response) Has(path string, expected any) *Response {
	r.JSON()
	actual, exists := getNestedValue(r.parsedJSON, path)
	assert.True(r.t, exists, "JSON key %q not found", path)

	// Type strictness: compare types first
	assert.IsType(r.t, expected, actual, "type mismatch for key %q", path)
	assert.Equal(r.t, expected, actual, "value mismatch for key %q", path)
	return r
}

// HasInt asserts that the JSON key equals the expected integer value.
// Handles JSON number conversion (JSON numbers are parsed as float64).
//
// Parameters:
//   - path: The JSON key path using dot notation
//   - expected: The expected integer value
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/count").HasInt("data.total", 42)
//
// Implementation:.
func (r *Response) HasInt(path string, expected int) *Response {
	r.JSON()
	actual, exists := getNestedValue(r.parsedJSON, path)
	assert.True(r.t, exists, "JSON key %q not found", path)

	// JSON numbers are float64, convert for comparison
	switch v := actual.(type) {
	case float64:
		assert.Equal(r.t, expected, int(v), "int value mismatch for key %q", path)
	case int:
		assert.Equal(r.t, expected, v, "int value mismatch for key %q", path)
	default:
		assert.Fail(r.t, "type mismatch", "expected numeric value for key %q, got %T", path, actual)
	}
	return r
}

// HasFloat asserts that the JSON key equals the expected float64 value.
// Uses InDelta for comparison to handle floating point precision issues.
//
// Parameters:
//   - path: The JSON key path using dot notation
//   - expected: The expected float64 value
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/score").HasFloat("data.score", 95.5)
//
// Implementation:.
func (r *Response) HasFloat(path string, expected float64) *Response {
	r.JSON()
	actual, exists := getNestedValue(r.parsedJSON, path)
	assert.True(r.t, exists, "JSON key %q not found", path)

	// Handle both float64 and int from JSON
	switch v := actual.(type) {
	case float64:
		assert.InDelta(r.t, expected, v, 0.0001, "float value mismatch for key %q", path)
	case int:
		assert.InDelta(r.t, expected, float64(v), 0.0001, "float value mismatch for key %q", path)
	default:
		assert.Fail(r.t, "type mismatch", "expected numeric value for key %q, got %T", path, actual)
	}
	return r
}

// HasString asserts that the JSON key equals the expected string value.
//
// Parameters:
//   - path: The JSON key path using dot notation
//   - expected: The expected string value
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users/1").HasString("data.user.email", "john@example.com")
//
// Implementation:.
func (r *Response) HasString(path, expected string) *Response {
	r.JSON()
	actual, exists := getNestedValue(r.parsedJSON, path)
	assert.True(r.t, exists, "JSON key %q not found", path)
	assert.Equal(r.t, expected, actual, "string value mismatch for key %q", path)
	return r
}

// HasBool asserts that the JSON key equals the expected boolean value.
//
// Parameters:
//   - path: The JSON key path using dot notation
//   - expected: The expected boolean value
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/status").HasBool("data.active", true)
//
// Implementation:.
func (r *Response) HasBool(path string, expected bool) *Response {
	r.JSON()
	actual, exists := getNestedValue(r.parsedJSON, path)
	assert.True(r.t, exists, "JSON key %q not found", path)
	assert.Equal(r.t, expected, actual, "bool value mismatch for key %q", path)
	return r
}

// MatchesRegex asserts that the JSON key value matches the regex pattern.
// The value at the specified path must be a string.
//
// Parameters:
//   - path: The JSON key path using dot notation
//   - pattern: The regex pattern to match
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users/1").MatchesRegex("data.user.email", `^[\w.-]+@[\w.-]+\.\w+$`)
//
// Implementation:.
func (r *Response) MatchesRegex(path, pattern string) *Response {
	r.JSON()
	actual, exists := getNestedValue(r.parsedJSON, path)
	assert.True(r.t, exists, "JSON key %q not found", path)

	str, ok := actual.(string)
	assert.True(r.t, ok, "expected string value for key %q to match regex, got %T", path, actual)

	matched, err := regexp.MatchString(pattern, str)
	assert.NoError(r.t, err, "invalid regex pattern")
	assert.True(r.t, matched, "value %q should match pattern %q for key %s", str, pattern, path)
	return r
}

// JSONMatches asserts that the entire JSON structure matches the expected map.
// Performs Has() assertions for each key-value pair in the expected map.
//
// Parameters:
//   - expected: Map of key paths to expected values
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users/1").JSONMatches(map[string]interface{}{
//	    "data.user.name": "John",
//	    "data.user.active": true,
//	})
//
// Implementation:.
func (r *Response) JSONMatches(expected map[string]interface{}) *Response {
	r.JSON()
	for path, expectedValue := range expected {
		r.Has(path, expectedValue)
	}
	return r
}

// ArrayLength asserts that the JSON key is an array with the expected length.
//
// Parameters:
//   - path: The JSON key path using dot notation
//   - expected: The expected array length
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users").ArrayLength("data.users", 10)
//
// Implementation:.
func (r *Response) ArrayLength(path string, expected int) *Response {
	r.JSON()
	actual, exists := getNestedValue(r.parsedJSON, path)
	assert.True(r.t, exists, "JSON key %q not found", path)

	arr, ok := actual.([]interface{})
	assert.True(r.t, ok, "expected array for key %q, got %T", path, actual)
	assert.Equal(r.t, expected, len(arr), "array length mismatch for key %q", path)
	return r
}

// IsNull asserts that the JSON key has a null value.
//
// Parameters:
//   - path: The JSON key path using dot notation
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users/999").IsNull("data.user")
//
// Implementation:.
func (r *Response) IsNull(path string) *Response {
	r.JSON()
	actual, exists := getNestedValue(r.parsedJSON, path)
	assert.True(r.t, exists, "JSON key %q not found", path)
	assert.Nil(r.t, actual, "expected null for key %q", path)
	return r
}

// IsNotNull asserts that the JSON key has a non-null value.
//
// Parameters:
//   - path: The JSON key path using dot notation
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users/1").IsNotNull("data.user")
//
// Implementation:.
func (r *Response) IsNotNull(path string) *Response {
	r.JSON()
	actual, exists := getNestedValue(r.parsedJSON, path)
	assert.True(r.t, exists, "JSON key %q not found", path)
	assert.NotNil(r.t, actual, "expected non-null for key %q", path)
	return r
}

// getNestedValue traverses a dot-notation path through JSON data.
// Internal helper function that supports nested objects and array indices.
// Returns the value and a boolean indicating if the path exists.
//
// Parameters:
//   - data: The root JSON object as a map
//   - path: Dot-notation path (e.g., "user.name", "items.0.name")
//
// Returns:
//   - interface{}: The value at the path (nil if not found)
//   - bool: True if the path exists, false otherwise
//
// Example:
//
//	// Internal usage - supports paths like:
//	// "user.name" -> accesses data["user"]["name"]
//	// "items.0.id" -> accesses data["items"][0]["id"]
//
// Implementation:.
func getNestedValue(data map[string]interface{}, path string) (interface{}, bool) {
	parts := strings.Split(path, ".")
	current := interface{}(data)

	for _, part := range parts {
		if current == nil {
			return nil, false
		}

		switch v := current.(type) {
		case map[string]interface{}:
			val, exists := v[part]
			if !exists {
				return nil, false
			}
			current = val
		case []interface{}:
			// Parse array index
			index, err := strconv.Atoi(part)
			if err != nil {
				return nil, false // not a valid index
			}
			if index < 0 || index >= len(v) {
				return nil, false // out of bounds
			}
			current = v[index]
		default:
			// Can't traverse further
			return nil, false
		}
	}

	return current, true
}

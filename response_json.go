package zentests

import (
	"encoding/json"
	"maps"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/stretchr/testify/assert"
)

// JSON parses and caches the response body as JSON.
// Enables JSON-specific assertions on the response. The body is parsed once; later calls
// reuse the result, even when parsing failed. Fails the test if JSON parsing fails.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users").JSON().Has("data.0.name", "John")
func (r *Response) JSON() *Response {
	if r.parsed {
		return r
	}
	r.parsed = true

	err := json.Unmarshal(r.Body(), &r.parsedJSON)
	assert.NoError(r.t, err, "JSON parsing failed")
	return r
}

// lookup parses the JSON body and returns the value at the dot-notation path.
// If the path does not exist, it records a failure and returns false, so callers can
// stop without reporting a second failure for the same key.
//
// Parameters:
//   - path: The JSON key path using dot notation
//
// Returns:
//   - any: The value at path (nil if not found)
//   - bool: True if the path exists
func (r *Response) lookup(path string) (any, bool) {
	r.JSON()
	value, exists := getNestedValue(r.parsedJSON, path)
	assert.True(r.t, exists, "JSON key %q not found", path)
	return value, exists
}

// lookupNumber returns the JSON number at path as a float64. It records a failure and
// returns false if the path does not exist or the value is not a number.
//
// Parameters:
//   - path: The JSON key path using dot notation
//
// Returns:
//   - float64: The numeric value at path
//   - bool: True if the value exists and is a number
func (r *Response) lookupNumber(path string) (float64, bool) {
	actual, exists := r.lookup(path)
	if !exists {
		return 0, false
	}

	number, isNumber := actual.(float64)
	assert.True(r.t, isNumber, "expected numeric value for key %q, got %T", path, actual)
	return number, isNumber
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
func (r *Response) HasKey(path string) *Response {
	r.lookup(path)
	return r
}

// Has asserts that the JSON key equals the expected value.
// Supports dot notation for nested keys. Fails if the key doesn't exist or if the
// value differs. Any Go numeric type (int, float32, ...) is accepted for numbers and
// compared as float64, because JSON numbers always decode to float64. Other types must
// match exactly.
//
// Parameters:
//   - path: The JSON key path using dot notation
//   - expected: The expected value
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users/1").Has("data.user.name", "John")
//	zt.Get(app, "/api/count").Has("data.count", 42) // any numeric type works
func (r *Response) Has(path string, expected any) *Response {
	actual, exists := r.lookup(path)
	if !exists {
		return r
	}

	expected = normalizeNumber(expected)
	if !assert.IsType(r.t, expected, actual, "type mismatch for key %q", path) {
		return r
	}
	assert.Equal(r.t, expected, actual, "value mismatch for key %q", path)
	return r
}

// normalizeNumber converts any Go numeric value to float64, the type JSON numbers decode to.
// Non-numeric values are returned unchanged.
//
// Parameters:
//   - value: The value to normalize
//
// Returns:
//   - any: float64 for numeric values, otherwise value unchanged
func normalizeNumber(value any) any {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return rv.Float()
	default:
		return value
	}
}

// HasInt asserts that the JSON key equals the expected integer value.
// Fails for non-integer numbers, so 42.9 does not match 42.
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
func (r *Response) HasInt(path string, expected int) *Response {
	number, ok := r.lookupNumber(path)
	if ok {
		assert.Equal(r.t, float64(expected), number, "int value mismatch for key %q", path)
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
func (r *Response) HasFloat(path string, expected float64) *Response {
	number, ok := r.lookupNumber(path)
	if ok {
		assert.InDelta(r.t, expected, number, 0.0001, "float value mismatch for key %q", path)
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
func (r *Response) HasString(path, expected string) *Response {
	actual, exists := r.lookup(path)
	if exists {
		assert.Equal(r.t, expected, actual, "string value mismatch for key %q", path)
	}
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
func (r *Response) HasBool(path string, expected bool) *Response {
	actual, exists := r.lookup(path)
	if exists {
		assert.Equal(r.t, expected, actual, "bool value mismatch for key %q", path)
	}
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
func (r *Response) MatchesRegex(path, pattern string) *Response {
	actual, exists := r.lookup(path)
	if !exists {
		return r
	}

	str, isString := actual.(string)
	if !assert.True(r.t, isString, "expected string value for key %q to match regex, got %T", path, actual) {
		return r
	}

	matched, err := regexp.MatchString(pattern, str)
	if !assert.NoError(r.t, err, "invalid regex pattern") {
		return r
	}
	assert.True(r.t, matched, "value %q should match pattern %q for key %s", str, pattern, path)
	return r
}

// JSONHasValues asserts that each listed key path holds its expected value.
// This is a subset check: keys that are not listed are ignored. Each entry is checked
// with Has, so the same type rules apply. Keys are checked in sorted order, so failures
// are reported in a stable order.
//
// Parameters:
//   - expected: Map of key paths to expected values
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users/1").JSONHasValues(map[string]any{
//	    "data.user.name":   "John",
//	    "data.user.active": true,
//	})
func (r *Response) JSONHasValues(expected map[string]any) *Response {
	for _, path := range slices.Sorted(maps.Keys(expected)) {
		r.Has(path, expected[path])
	}
	return r
}

// JSONMatches is an alias for JSONHasValues. It is a subset check: it asserts that each
// listed key path holds its expected value, and ignores any other keys in the body.
// Prefer JSONHasValues in new code.
//
// Parameters:
//   - expected: Map of key paths to expected values
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users/1").JSONMatches(map[string]any{
//	    "data.user.name": "John",
//	})
func (r *Response) JSONMatches(expected map[string]any) *Response {
	return r.JSONHasValues(expected)
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
func (r *Response) ArrayLength(path string, expected int) *Response {
	actual, exists := r.lookup(path)
	if !exists {
		return r
	}

	arr, isArray := actual.([]any)
	if !assert.True(r.t, isArray, "expected array for key %q, got %T", path, actual) {
		return r
	}
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
func (r *Response) IsNull(path string) *Response {
	actual, exists := r.lookup(path)
	if exists {
		assert.Nil(r.t, actual, "expected null for key %q", path)
	}
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
func (r *Response) IsNotNull(path string) *Response {
	actual, exists := r.lookup(path)
	if exists {
		assert.NotNil(r.t, actual, "expected non-null for key %q", path)
	}
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
//   - any: The value at the path (nil if not found)
//   - bool: True if the path exists, false otherwise
//
// Example:
//
//	// Internal usage - supports paths like:
//	// "user.name" -> accesses data["user"]["name"]
//	// "items.0.id" -> accesses data["items"][0]["id"]
func getNestedValue(data map[string]any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	current := any(data)

	for _, part := range parts {
		if current == nil {
			return nil, false
		}

		switch v := current.(type) {
		case map[string]any:
			val, exists := v[part]
			if !exists {
				return nil, false
			}
			current = val
		case []any:
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

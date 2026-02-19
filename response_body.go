package zentests

import (
	"regexp"

	"github.com/stretchr/testify/assert"
)

// Contains asserts that the response body contains the specified substring.
// Case-sensitive substring search in the response body.
//
// Parameters:
//   - substring: The substring to search for in the body
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/page").Contains("Welcome")
//
// Implementation:.
func (r *Response) Contains(substring string) *Response {
	assert.Contains(r.t, r.BodyString(), substring, "body should contain substring")
	return r
}

// NotContains asserts that the response body does not contain the specified substring.
// Useful for ensuring sensitive data is not present in responses.
//
// Parameters:
//   - substring: The substring that should not be in the body
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/public").NotContains("password")
//
// Implementation:.
func (r *Response) NotContains(substring string) *Response {
	assert.NotContains(r.t, r.BodyString(), substring, "body should not contain substring")
	return r
}

// BodyMatches asserts that the response body matches the regex pattern.
// Useful for validating dynamic content or complex text patterns.
//
// Parameters:
//   - pattern: The regex pattern to match against the body
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/users/1").BodyMatches(`User ID: \d+`)
//
// Implementation:.
func (r *Response) BodyMatches(pattern string) *Response {
	matched, err := regexp.MatchString(pattern, r.BodyString())
	assert.NoError(r.t, err, "invalid regex pattern")
	assert.True(r.t, matched, "body should match pattern %q", pattern)
	return r
}

// Equals asserts that the response body equals the exact expected string.
// Performs an exact string comparison of the entire body.
//
// Parameters:
//   - expected: The expected body content
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/ping").Equals("pong")
//
// Implementation:.
func (r *Response) Equals(expected string) *Response {
	assert.Equal(r.t, expected, r.BodyString(), "body mismatch")
	return r
}

// IsEmpty asserts that the response body is empty.
// Checks that the body has zero length.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Delete(app, "/users/1").NoContent().IsEmpty()
//
// Implementation:.
func (r *Response) IsEmpty() *Response {
	assert.Empty(r.t, r.BodyString(), "body should be empty")
	return r
}

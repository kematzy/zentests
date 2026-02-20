package zentests

import (
	"github.com/stretchr/testify/assert"
)

// Status asserts that the HTTP status code equals the expected value.
// This is the base status assertion method used by all other status shortcuts.
//
// Parameters:
//   - expected: The expected HTTP status code (e.g., 200, 404, 500)
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/users").Status(200)
//	zt.Get(app, "/notfound").Status(404)
func (r *Response) Status(expected int) *Response {
	assert.Equal(r.t, expected, r.StatusCode, "status code mismatch")
	return r
}

// OK asserts that the HTTP status code is 200 (OK).
// Convenience method equivalent to Status(200).
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/users").OK()
func (r *Response) OK() *Response {
	return r.Status(200)
}

// Created asserts that the HTTP status code is 201 (Created).
// Convenience method equivalent to Status(201).
// Typically used after POST requests that create new resources.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.PostJSON(app, "/users", data).Created()
func (r *Response) Created() *Response {
	return r.Status(201)
}

// Accepted asserts that the HTTP status code is 202 (Accepted).
// Convenience method equivalent to Status(202).
// Indicates the request has been accepted for processing but not completed.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.PostJSON(app, "/jobs", jobData).Accepted()
func (r *Response) Accepted() *Response {
	return r.Status(202)
}

// NoContent asserts that the HTTP status code is 204 (No Content).
// Convenience method equivalent to Status(204).
// Typically used for DELETE operations or successful operations with no response body.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Delete(app, "/users/1").NoContent()
func (r *Response) NoContent() *Response {
	return r.Status(204)
}

// BadRequest asserts that the HTTP status code is 400 (Bad Request).
// Convenience method equivalent to Status(400).
// Indicates the server cannot process the request due to client error.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.PostJSON(app, "/users", invalidData).BadRequest()
func (r *Response) BadRequest() *Response {
	return r.Status(400)
}

// Unauthorized asserts that the HTTP status code is 401 (Unauthorized).
// Convenience method equivalent to Status(401).
// Indicates authentication is required or has failed.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/protected").Unauthorized()
func (r *Response) Unauthorized() *Response {
	return r.Status(401)
}

// Forbidden asserts that the HTTP status code is 403 (Forbidden).
// Convenience method equivalent to Status(403).
// Indicates the server understood the request but refuses to authorize it.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/admin").Forbidden()
func (r *Response) Forbidden() *Response {
	return r.Status(403)
}

// NotFound asserts that the HTTP status code is 404 (Not Found).
// Convenience method equivalent to Status(404).
// Indicates the requested resource was not found.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/users/999").NotFound()
func (r *Response) NotFound() *Response {
	return r.Status(404)
}

// Unprocessable asserts that the HTTP status code is 422 (Unprocessable Entity).
// Convenience method equivalent to Status(422).
// Indicates the request was well-formed but contains semantic errors.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.PostJSON(app, "/users", invalidUserData).Unprocessable()
func (r *Response) Unprocessable() *Response {
	return r.Status(422)
}

// ServerError asserts that the HTTP status code is 500 (Internal Server Error).
// Convenience method equivalent to Status(500).
// Indicates an unexpected condition prevented the server from fulfilling the request.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/broken").ServerError()
func (r *Response) ServerError() *Response {
	return r.Status(500)
}

// ================================================================================================
// HEADER ASSERTIONS
// ================================================================================================

// HasHeader asserts that a header exists and has the expected value.
// Performs an exact match comparison. Header names are case-insensitive.
//
// Parameters:
//   - key: The header name (e.g., "Content-Type", "X-Custom-Header")
//   - value: The expected header value
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users").HasHeader("Content-Type", "application/json")
func (r *Response) HasHeader(key, value string) *Response {
	actual := r.Header.Get(key)
	assert.Equal(r.t, value, actual, "header %s mismatch", key)
	return r
}

// HeaderContains asserts that a header value contains the expected substring.
// Useful for flexible checks like verifying "application/json" is in Content-Type
// without matching the entire header value including charset.
//
// Parameters:
//   - key: The header name
//   - substring: The substring that must be present in the header value
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users").HeaderContains("Content-Type", "json")
func (r *Response) HeaderContains(key, substring string) *Response {
	actual := r.Header.Get(key)
	assert.Contains(r.t, actual, substring, "header %s should contain %q", key, substring)
	return r
}

// HeaderPresent asserts that a header is present and non-empty.
// Checks for the existence of a header without validating its value.
//
// Parameters:
//   - key: The header name
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/secure").HeaderPresent("X-Auth-Token")
func (r *Response) HeaderPresent(key string) *Response {
	assert.NotEmpty(r.t, r.Header.Get(key), "Expected header %s to exist", key)
	return r
}

// HeaderNotPresent asserts that a header is not present or is empty.
// Useful for security checks to ensure sensitive headers are not exposed.
//
// Parameters:
//   - key: The header name
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/public").HeaderNotPresent("Authorization")
func (r *Response) HeaderNotPresent(key string) *Response {
	assert.Empty(r.t, r.Header.Get(key), "Expected header %s to not exist", key)
	return r
}

// HeaderHasValues asserts that a header has exactly the expected values.
// For headers that can have multiple values (e.g., Set-Cookie, X-Features).
// Order of values does not matter.
//
// Parameters:
//   - key: The header name
//   - values: Slice of expected header values
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/features").HeaderHasValues("X-Features", []string{"audio", "video"})
func (r *Response) HeaderHasValues(key string, values []string) *Response {
	headerValues := r.Header.Values(key)
	assert.ElementsMatch(r.t, values, headerValues)
	return r
}

// CookieHasValues asserts that response cookies match the expected values.
// Parses all Set-Cookie headers and verifies each cookie name-value pair.
//
// Parameters:
//   - expected: Map of cookie names to expected values
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/settings").CookieHasValues(map[string]string{
//	    "lang": "en",
//	    "theme": "dark",
//	})
func (r *Response) CookieHasValues(expected map[string]string) *Response {
	// NOTE! must be called on the response
	cookies := r.resp.Cookies() // Parses "Set-Cookie" headers automatically

	actual := make(map[string]string)
	for _, c := range cookies {
		actual[c.Name] = c.Value
	}

	for name, val := range expected {
		assert.Equal(r.t, val, actual[name], "Cookie %s value mismatch", name)
	}
	return r
}

// ================================================================================================
// CONTENT TYPE SHORTCUTS
// ================================================================================================

// HasContentType asserts that the Content-Type header equals the expected value.
// Performs an exact match on the Content-Type header.
//
// Parameters:
//   - contentType: The expected Content-Type value (e.g., "application/json")
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/page").HasContentType("text/html; charset=utf-8")
func (r *Response) HasContentType(contentType string) *Response {
	return r.HasHeader("Content-Type", contentType)
}

// IsJSON asserts that the Content-Type header contains "application/json".
// Quick check for JSON API responses. Case-sensitive substring match.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users").OK().IsJSON()
func (r *Response) IsJSON() *Response {
	return r.HeaderContains("Content-Type", "application/json")
}

// IsHTML asserts that the Content-Type header contains "text/html".
// Quick check for HTML page responses. Case-sensitive substring match.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/").OK().IsHTML()
func (r *Response) IsHTML() *Response {
	return r.HeaderContains("Content-Type", "text/html")
}

// IsPlainText asserts that the Content-Type header contains "text/plain".
// Quick check for plain text responses. Case-sensitive substring match.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/health").OK().IsPlainText()
func (r *Response) IsPlainText() *Response {
	return r.HeaderContains("Content-Type", "text/plain")
}

// IsCSS asserts that the Content-Type header contains "text/css".
// Quick check for CSS stylesheet responses.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/style.css").OK().IsCSS()
func (r *Response) IsCSS() *Response {
	return r.HeaderContains("Content-Type", "text/css")
}

// IsJS asserts that the Content-Type header contains "javascript".
// Matches both "application/javascript" and "text/javascript".
// Quick check for JavaScript responses.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/app.js").OK().IsJS()
func (r *Response) IsJS() *Response {
	return r.HeaderContains("Content-Type", "javascript")
}

// IsXML asserts that the Content-Type header contains "xml".
// Matches both "application/xml" and "text/xml".
// Quick check for XML responses.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/data.xml").OK().IsXML()
func (r *Response) IsXML() *Response {
	return r.HeaderContains("Content-Type", "xml")
}

// IsXHR asserts that the X-Requested-With header contains "XMLHttpRequest".
// Useful for detecting AJAX requests in responses.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/api/users").OK().IsXHR()
func (r *Response) IsXHR() *Response {
	return r.HeaderContains("X-Requested-With", "XMLHttpRequest")
}

// IsImage asserts that the Content-Type header contains "image/".
// General check for any image type responses.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/photo.jpg").OK().IsImage()
func (r *Response) IsImage() *Response {
	return r.HeaderContains("Content-Type", "image/")
}

// IsPNG asserts that the Content-Type header contains "image/png".
// Quick check for PNG image responses.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/logo.png").OK().IsPNG()
func (r *Response) IsPNG() *Response {
	return r.HeaderContains("Content-Type", "image/png")
}

// IsJPEG asserts that the Content-Type header contains "image/jpeg".
// Quick check for JPEG image responses.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/photo.jpg").OK().IsJPEG()
func (r *Response) IsJPEG() *Response {
	return r.HeaderContains("Content-Type", "image/jpeg")
}

// IsGIF asserts that the Content-Type header contains "image/gif".
// Quick check for GIF image responses.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/animation.gif").OK().IsGIF()
func (r *Response) IsGIF() *Response {
	return r.HeaderContains("Content-Type", "image/gif")
}

// IsSVG asserts that the Content-Type header contains "image/svg+xml".
// Quick check for SVG image responses.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/icon.svg").OK().IsSVG()
func (r *Response) IsSVG() *Response {
	return r.HeaderContains("Content-Type", "image/svg+xml")
}

// IsWebP asserts that the Content-Type header contains "image/webp".
// Quick check for WebP image responses.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/photo.webp").OK().IsWebP()
func (r *Response) IsWebP() *Response {
	return r.HeaderContains("Content-Type", "image/webp")
}

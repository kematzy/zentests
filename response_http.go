package zentests

import (
	"net/http"
	"slices"

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

// IsRedirect asserts that the HTTP status code is one of the redirects that carry a
// Location header: 301, 302, 303, 307 or 308. 304 (Not Modified) is not included.
// Use RedirectsTo to also check where the response points.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/old-url").IsRedirect()
func (r *Response) IsRedirect() *Response {
	r.t.Helper()
	assert.True(r.t, redirectStatuses[r.StatusCode],
		"expected a redirect status (301, 302, 303, 307, 308), got %d", r.StatusCode)
	return r
}

// RedirectsTo asserts that the response is a redirect to the expected location.
// It checks IsRedirect and an exact match on the Location header.
//
// Parameters:
//   - location: The expected Location header value
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.PostForm(app, "/login", data).RedirectsTo("/dashboard")
func (r *Response) RedirectsTo(location string) *Response {
	r.t.Helper()
	r.IsRedirect()
	assert.Equal(r.t, location, r.Header.Get("Location"), "redirect location mismatch")
	return r
}

// redirectStatuses lists the redirect codes whose responses carry a Location header.
var redirectStatuses = map[int]bool{
	http.StatusMovedPermanently:  true,
	http.StatusFound:             true,
	http.StatusSeeOther:          true,
	http.StatusTemporaryRedirect: true,
	http.StatusPermanentRedirect: true,
}

// IsSuccess asserts that the HTTP status code is in the 2xx range.
// Use Status or its shortcuts when the exact code matters.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/health").IsSuccess()
func (r *Response) IsSuccess() *Response {
	r.t.Helper()
	assert.True(r.t, r.StatusCode >= 200 && r.StatusCode < 300, "expected a 2xx status, got %d", r.StatusCode)
	return r
}

// IsClientError asserts that the HTTP status code is in the 4xx range.
// Use Status or its shortcuts when the exact code matters.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/missing").IsClientError()
func (r *Response) IsClientError() *Response {
	r.t.Helper()
	assert.True(r.t, r.StatusCode >= 400 && r.StatusCode < 500, "expected a 4xx status, got %d", r.StatusCode)
	return r
}

// IsServerError asserts that the HTTP status code is in the 5xx range.
// Use Status or its shortcuts when the exact code matters.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/broken").IsServerError()
func (r *Response) IsServerError() *Response {
	r.t.Helper()
	assert.True(r.t, r.StatusCode >= 500 && r.StatusCode < 600, "expected a 5xx status, got %d", r.StatusCode)
	return r
}

// StatusIn asserts that the HTTP status code is one of the expected values.
// Use it when a handler may answer with a small set of codes.
//
// Parameters:
//   - codes: The acceptable HTTP status codes
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Delete(app, "/users/1").StatusIn(200, 202, 204)
func (r *Response) StatusIn(codes ...int) *Response {
	r.t.Helper()
	assert.True(r.t, slices.Contains(codes, r.StatusCode),
		"expected one of the statuses %v, got %d", codes, r.StatusCode)
	return r
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
//
//	zt.Get(app, "/api/users").HasHeader("Content-Type", "application/css") // raises an error
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
//
//	zt.Get(app, "/api/users").HeaderContains("Content-Type", "image") // raises an error
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

	for name, wantValue := range expected {
		gotValue, found := actual[name]
		if assert.True(r.t, found, "Cookie %s was not set", name) {
			assert.Equal(r.t, wantValue, gotValue, "Cookie %s value mismatch", name)
		}
	}
	return r
}

// Cookie returns the named cookie from the response, failing the test when it was not set.
// Use it to check attributes the [Jar] does not keep, such as HttpOnly, Path, Secure and SameSite.
//
// Parameters:
//   - name: The cookie name
//
// Returns:
//   - *[http.Cookie]: The parsed Set-Cookie entry with its attributes
//
// Example:
//
//	cookie := zt.Get(app, "/login").OK().Cookie("session")
//	assert.True(t, cookie.HttpOnly)
func (r *Response) Cookie(name string) *http.Cookie {
	r.t.Helper()
	for _, cookie := range r.resp.Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	r.t.Errorf("cookie %s was not set", name)
	r.t.FailNow()
	return nil
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

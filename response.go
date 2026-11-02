package zentests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Response wraps an HTTP response with fluent assertion methods.
// It provides a chainable API for making assertions on HTTP responses including
// status codes, headers, body content, and JSON data. The response body is lazily
// loaded when first accessed to avoid unnecessary reads.
//
// Response maintains internal state including the parsed JSON body and tracks
// whether the body has been read to prevent multiple reads from the response stream.
//
// Fields:
//   - t: The [testing.TB] instance for making assertions
//   - StatusCode: The HTTP status code from the response
//   - Header: The HTTP headers from the response
//   - resp: The underlying [http.Response]
//   - body: The cached response body bytes (lazy loaded)
//   - bodyRead: Flag indicating if body has been read
//   - parsed: Flag indicating if JSON parsing has been attempted
//   - parsedJSON: The parsed JSON body as a map (cached after first parse)
type Response struct {
	t          testing.TB
	Header     http.Header
	resp       *http.Response
	parsedJSON map[string]any
	body       []byte
	StatusCode int
	bodyRead   bool
	parsed     bool
}

// Body returns the response body as bytes.
// The body is lazily loaded on first call and cached for subsequent accesses.
// This prevents multiple reads from the response stream which would return empty.
//
// Returns:
//   - []byte: The response body as bytes
//
// Example:
//
//	body := resp.Body()
//	fmt.Println(string(body))
func (r *Response) Body() []byte {
	if !r.bodyRead {
		var err error
		r.body, err = io.ReadAll(r.resp.Body)
		assert.NoError(r.t, err)

		err = r.resp.Body.Close()
		assert.NoError(r.t, err)

		r.bodyRead = true
	}
	return r.body
}

// BodyString returns the response body as a string.
// Convenience method that converts Body() bytes to string.
//
// Returns:
//   - string: The response body as a string
//
// Example:
//
//	bodyStr := resp.BodyString()
//	assert.Contains(t, bodyStr, "success")
func (r *Response) BodyString() string {
	return string(r.Body())
}

// redactedValue replaces the values of sensitive headers in Debug and Dump output.
const redactedValue = "[REDACTED]"

// sensitiveHeaders lists canonical header names whose values must not be printed,
// because test logs are often shared or stored in CI.
var sensitiveHeaders = map[string]bool{
	"Authorization":       true,
	"Proxy-Authorization": true,
	"Cookie":              true,
	"Set-Cookie":          true,
	"X-Api-Key":           true,
	"X-Auth-Token":        true,
	"X-Access-Token":      true,
	"X-Refresh-Token":     true,
	"X-Session-Id":        true,
	"X-Csrf-Token":        true,
}

// redactedHeaders returns a copy of the response headers with sensitive values masked.
//
// Returns:
//   - [http.Header]: The headers safe to print
func (r *Response) redactedHeaders() http.Header {
	redacted := make(http.Header, len(r.Header))
	for key, values := range r.Header {
		if sensitiveHeaders[http.CanonicalHeaderKey(key)] {
			redacted[key] = []string{redactedValue}
			continue
		}
		redacted[key] = values
	}
	return redacted
}

// Debug prints detailed response information to stdout for troubleshooting.
// Outputs status code, headers, body, and parsed JSON (if available).
// Sensitive headers such as Authorization and Set-Cookie are redacted.
// This is useful during test development to inspect response details.
// Returns the receiver for method chaining.
//
// Returns:
//   - *Response: The receiver for method chaining
//
// Example:
//
//	zt.Get(app, "/users").Debug().OK()
func (r *Response) Debug() *Response {
	fmt.Printf("\n=== ZENTESTS DEBUG ===\n\n%s\n", r.Dump())

	if r.parsedJSON != nil {
		pretty, err := json.MarshalIndent(r.parsedJSON, "", "  ")
		if assert.NoError(r.t, err, "error marshalIndenting parsed JSON") {
			fmt.Printf("Parsed JSON:\n%s\n", pretty)
		}
	}

	fmt.Printf("\n======================\n\n")
	return r
}

// Dump returns raw response details as a formatted string.
// Useful for custom logging or including response details in test failure messages.
// Unlike Debug(), this returns a string instead of printing to stdout.
// Sensitive headers such as Authorization and Set-Cookie are redacted.
//
// Returns:
//   - string: Formatted string with status, headers, and body
//
// Example:
//
//	details := resp.Dump()
//	t.Logf("Response details: %s", details)
func (r *Response) Dump() string {
	return fmt.Sprintf("Status: %d\nHeaders: %v\nBody: %s",
		r.StatusCode, r.redactedHeaders(), r.BodyString())
}

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
//   - t: The testing.T instance for making assertions
//   - StatusCode: The HTTP status code from the response
//   - Header: The HTTP headers from the response
//   - resp: The underlying http.Response
//   - body: The cached response body bytes (lazy loaded)
//   - bodyRead: Flag indicating if body has been read
//   - parsedJSON: The parsed JSON body as a map (cached after first parse)
type Response struct {
	t          *testing.T
	Header     http.Header
	resp       *http.Response
	parsedJSON map[string]any
	body       []byte
	StatusCode int
	bodyRead   bool
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
	if r.resp == nil {
		return []byte{}
	}
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

// Debug prints detailed response information to stdout for troubleshooting.
// Outputs status code, headers, body, and parsed JSON (if available).
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
	fmt.Println("")
	fmt.Println("=== ZENTESTS DEBUG ===")
	fmt.Println("")
	fmt.Printf("Status: %d\n", r.StatusCode)
	fmt.Printf("Headers: %v\n", sanitizeHeaders(r.Header))
	fmt.Printf("Body: %s\n", r.BodyString())

	if r.parsedJSON != nil {
		var err error
		pretty, err := json.MarshalIndent(r.parsedJSON, "", "  ")
		assert.NoErrorf(r.t, err, "error marshalIndenting parsed JSON: %s", err)
		fmt.Printf("Parsed JSON:\n%s\n", string(pretty))
	}
	fmt.Println("")
	fmt.Println("======================")
	fmt.Println("")
	return r
}

// Dump returns raw response details as a formatted string.
// Useful for custom logging or including response details in test failure messages.
// Unlike Debug(), this returns a string instead of printing to stdout.
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
		r.StatusCode, sanitizeHeaders(r.Header), r.BodyString())
}

var sensitiveHeaders = map[string]bool{
	"Authorization":    true,
	"Cookie":           true,
	"Set-Cookie":       true,
	"X-Api-Key":        true,
	"X-Auth-Token":     true,
	"X-Access-Token":   true,
	"X-Refresh-Token":  true,
	"X-Session-Id":     true,
	"X-Csrf-Token":     true,
	"X-Requested-With": true,
}

func sanitizeHeaders(header http.Header) http.Header {
	if header == nil {
		return nil
	}
	sanitized := make(http.Header)
	for k, v := range header {
		if sensitiveHeaders[k] {
			sanitized[k] = []string{"[REDACTED]"}
		} else {
			sanitized[k] = v
		}
	}
	return sanitized
}

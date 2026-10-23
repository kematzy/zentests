package zentests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

// defaultTestConfig is used whenever a request is sent with a zero-value fiber.TestConfig.
// It matches Fiber's own defaults: a one second timeout that fails the request on expiry.
// Without it, a zero Timeout would make app.Test wait forever for a stuck handler.
var defaultTestConfig = fiber.TestConfig{
	Timeout:       time.Second,
	FailOnTimeout: true,
}

// testConfig returns cfg, or defaultTestConfig when cfg is the zero value.
func testConfig(cfg fiber.TestConfig) fiber.TestConfig {
	if cfg == (fiber.TestConfig{}) {
		return defaultTestConfig
	}
	return cfg
}

// execute performs an HTTP request against the Fiber app and returns a Response.
// This is the single place where requests are sent; every public method delegates here.
// It stops the test (via require) if the request cannot be executed, for example on timeout.
//
// Parameters:
//   - t: The [testing.TB] instance for assertions
//   - app: The Fiber application to test against
//   - method: The HTTP method (GET, POST, PUT, etc.)
//   - path: The request path (e.g., "/api/users")
//   - body: The request body as an [io.Reader] (nil for requests without a body)
//   - headers: Map of headers to set on the request
//   - cfg: TestConfig for timeout/behavior; the zero value uses defaultTestConfig
//
// Returns:
//   - *Response: A Response wrapper containing the HTTP response and assertion methods
//
// Example:
//
//	// Internal usage only - use Get(), Post(), etc. instead
//	resp := execute(t, app, "GET", "/users", nil, nil, fiber.TestConfig{})
func execute(t testing.TB, app *fiber.App, method, path string, body io.Reader, headers map[string]string, cfg fiber.TestConfig) *Response {
	req := httptest.NewRequestWithContext(context.Background(), method, path, body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := app.Test(req, testConfig(cfg)) //nolint:bodyclose // body is stored in Response for lazy reading
	require.NoError(t, err, "request execution failed")

	return &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
	}
}

// executeJSON marshals data to JSON, sets Content-Type to application/json and
// performs the request. It stops the test if marshaling fails.
//
// Parameters:
//   - t: The [testing.TB] instance for assertions
//   - app: The Fiber application to test against
//   - method: The HTTP method
//   - path: The request path
//   - data: Any data structure to be JSON-encoded (nil is sent as JSON null)
//   - cfg: TestConfig for timeout/behavior; the zero value uses defaults
//
// Returns:
//   - *Response: A Response wrapper for making assertions
func executeJSON(t testing.TB, app *fiber.App, method, path string, data any, cfg fiber.TestConfig) *Response {
	jsonBody, err := json.Marshal(data)
	require.NoError(t, err, "JSON marshaling failed")

	return execute(t, app, method, path, bytes.NewReader(jsonBody), SetHeader("Content-Type", "application/json"), cfg)
}

// formBody encodes data as an application/x-www-form-urlencoded request body.
// Keys and values are escaped, so characters such as "&", "=", "+" and spaces
// reach the server unchanged.
func formBody(data map[string]string) io.Reader {
	values := url.Values{}
	for key, value := range data {
		values.Set(key, value)
	}
	return strings.NewReader(values.Encode())
}

// =================================================================================================
// REQUEST METHODS
// Each method without a config delegates to its *WithConfig twin with the zero-value config,
// which means the default timeout applies.
// =================================================================================================

// Get performs a GET request against the Fiber app.
// Returns a *Response for chaining assertions. Use this for testing GET endpoints.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users")
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	zt := zentests.New(t)
//	resp := zt.Get(app, "/users")
//	resp.OK().IsJSON().Has("data.users.0.name", "John")
func (zt *T) Get(app *fiber.App, path string) *Response {
	return zt.GetWithConfig(app, path, fiber.TestConfig{})
}

// Post performs a POST request with a raw body.
// Use PostJSON for JSON payloads or PostForm for form data. This method is suitable
// for sending raw bytes when you need full control over the request body.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users")
//   - body: The raw request body as bytes
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	body := []byte("raw data")
//	zt.Post(app, "/upload", body).OK()
func (zt *T) Post(app *fiber.App, path string, body []byte) *Response {
	return zt.PostWithConfig(app, path, body, fiber.TestConfig{})
}

// PostJSON performs a POST request with JSON-encoded data.
// Automatically marshals the data to JSON and sets Content-Type to application/json.
// Fails the test if JSON marshaling fails.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users")
//   - data: Any data structure to be JSON-encoded (maps, structs, etc.)
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	data := map[string]string{"name": "John", "email": "john@example.com"}
//	zt.PostJSON(app, "/users", data).Created().Has("id", float64(1))
func (zt *T) PostJSON(app *fiber.App, path string, data any) *Response {
	return zt.PostJSONWithConfig(app, path, data, fiber.TestConfig{})
}

// PostForm performs a POST request with form data.
// Automatically encodes the data as application/x-www-form-urlencoded.
// Useful for testing HTML form submissions or traditional form endpoints.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/login")
//   - data: Map of form field names to values
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	data := map[string]string{"username": "john", "password": "secret"}
//	zt.PostForm(app, "/login", data).OK()
func (zt *T) PostForm(app *fiber.App, path string, data map[string]string) *Response {
	return zt.PostFormWithConfig(app, path, data, fiber.TestConfig{})
}

// Put performs a PUT request with a raw body.
// Use PutJSON for JSON payloads. PUT is typically used for full resource updates.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users/1")
//   - body: The raw request body as bytes
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	body := []byte(`{"name": "Updated Name"}`)
//	zt.Put(app, "/users/1", body).OK()
func (zt *T) Put(app *fiber.App, path string, body []byte) *Response {
	return zt.PutWithConfig(app, path, body, fiber.TestConfig{})
}

// PutJSON performs a PUT request with JSON-encoded data.
// Automatically marshals the data to JSON and sets Content-Type to application/json.
// Fails the test if JSON marshaling fails. PUT is typically used for full resource updates.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users/1")
//   - data: Any data structure to be JSON-encoded
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	data := map[string]string{"name": "Updated Name"}
//	zt.PutJSON(app, "/users/1", data).OK()
func (zt *T) PutJSON(app *fiber.App, path string, data any) *Response {
	return zt.PutJSONWithConfig(app, path, data, fiber.TestConfig{})
}

// Patch performs a PATCH request with a raw body.
// Use PatchJSON for JSON payloads. PATCH is typically used for partial resource updates.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users/1")
//   - body: The raw request body as bytes
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	body := []byte(`{"status": "active"}`)
//	zt.Patch(app, "/users/1", body).OK()
func (zt *T) Patch(app *fiber.App, path string, body []byte) *Response {
	return zt.PatchWithConfig(app, path, body, fiber.TestConfig{})
}

// PatchJSON performs a PATCH request with JSON-encoded data.
// Automatically marshals the data to JSON and sets Content-Type to application/json.
// Fails the test if JSON marshaling fails. PATCH is typically used for partial updates.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users/1")
//   - data: Any data structure to be JSON-encoded
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	data := map[string]string{"status": "active"}
//	zt.PatchJSON(app, "/users/1", data).OK()
func (zt *T) PatchJSON(app *fiber.App, path string, data any) *Response {
	return zt.PatchJSONWithConfig(app, path, data, fiber.TestConfig{})
}

// Delete performs a DELETE request.
// Returns a *Response for chaining assertions. Use this for testing DELETE endpoints.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users/1")
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	zt.Delete(app, "/users/1").NoContent()
func (zt *T) Delete(app *fiber.App, path string) *Response {
	return zt.DeleteWithConfig(app, path, fiber.TestConfig{})
}

// DeleteJSON performs a DELETE request with a JSON body.
// Some APIs require a request body for DELETE operations (e.g., bulk delete with IDs).
// Automatically marshals the data to JSON and sets Content-Type to application/json.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users/bulk")
//   - data: Any data structure to be JSON-encoded
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	data := map[string][]int{"ids": {1, 2, 3}}
//	zt.DeleteJSON(app, "/users/bulk", data).OK()
func (zt *T) DeleteJSON(app *fiber.App, path string, data any) *Response {
	return zt.DeleteJSONWithConfig(app, path, data, fiber.TestConfig{})
}

// =================================================================================================
// METHODS WITH TESTCONFIG (Fiber v3+)
// A zero-value TestConfig uses the default timeout (one second, fail on timeout).
// =================================================================================================

// GetWithConfig performs a GET request with custom TestConfig for timeout/behavior control.
// Use this when you need to customize test behavior beyond defaults.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users")
//   - cfg: TestConfig for custom timeout/behavior settings
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	zt.GetWithConfig(app, "/slow-endpoint", fiber.TestConfig{
//	    Timeout: 5 * time.Second,
//	}).OK()
func (zt *T) GetWithConfig(app *fiber.App, path string, cfg fiber.TestConfig) *Response {
	return execute(zt.T, app, "GET", path, nil, nil, cfg)
}

// PostWithConfig performs a POST request with raw body and custom TestConfig.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users")
//   - body: The raw request body as bytes
//   - cfg: TestConfig for custom timeout/behavior settings
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	zt.PostWithConfig(app, "/upload", body, fiber.TestConfig{
//	    Timeout: 10 * time.Second,
//	}).OK()
func (zt *T) PostWithConfig(app *fiber.App, path string, body []byte, cfg fiber.TestConfig) *Response {
	return execute(zt.T, app, "POST", path, bytes.NewReader(body), nil, cfg)
}

// PostJSONWithConfig performs a POST request with JSON data and custom TestConfig.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users")
//   - data: Any data structure to be JSON-encoded
//   - cfg: TestConfig for custom timeout/behavior settings
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	zt.PostJSONWithConfig(app, "/users", data, fiber.TestConfig{
//	    Timeout:       5 * time.Second,
//	    FailOnTimeout: true,
//	}).Created()
func (zt *T) PostJSONWithConfig(app *fiber.App, path string, data any, cfg fiber.TestConfig) *Response {
	return executeJSON(zt.T, app, "POST", path, data, cfg)
}

// PostFormWithConfig performs a POST request with form data and custom TestConfig.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/login")
//   - data: Map of form field names to values
//   - cfg: TestConfig for custom timeout/behavior settings
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	zt.PostFormWithConfig(app, "/login", data, fiber.TestConfig{
//	    Timeout: 5 * time.Second,
//	}).OK()
func (zt *T) PostFormWithConfig(app *fiber.App, path string, data map[string]string, cfg fiber.TestConfig) *Response {
	headers := SetHeader("Content-Type", "application/x-www-form-urlencoded")
	return execute(zt.T, app, "POST", path, formBody(data), headers, cfg)
}

// PutWithConfig performs a PUT request with raw body and custom TestConfig.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users/1")
//   - body: The raw request body as bytes
//   - cfg: TestConfig for custom timeout/behavior settings
//
// Returns:
//   - *Response: A Response wrapper for making assertions
func (zt *T) PutWithConfig(app *fiber.App, path string, body []byte, cfg fiber.TestConfig) *Response {
	return execute(zt.T, app, "PUT", path, bytes.NewReader(body), nil, cfg)
}

// PutJSONWithConfig performs a PUT request with JSON data and custom TestConfig.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users/1")
//   - data: Any data structure to be JSON-encoded
//   - cfg: TestConfig for custom timeout/behavior settings
//
// Returns:
//   - *Response: A Response wrapper for making assertions
func (zt *T) PutJSONWithConfig(app *fiber.App, path string, data any, cfg fiber.TestConfig) *Response {
	return executeJSON(zt.T, app, "PUT", path, data, cfg)
}

// PatchWithConfig performs a PATCH request with raw body and custom TestConfig.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users/1")
//   - body: The raw request body as bytes
//   - cfg: TestConfig for custom timeout/behavior settings
//
// Returns:
//   - *Response: A Response wrapper for making assertions
func (zt *T) PatchWithConfig(app *fiber.App, path string, body []byte, cfg fiber.TestConfig) *Response {
	return execute(zt.T, app, "PATCH", path, bytes.NewReader(body), nil, cfg)
}

// PatchJSONWithConfig performs a PATCH request with JSON data and custom TestConfig.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users/1")
//   - data: Any data structure to be JSON-encoded
//   - cfg: TestConfig for custom timeout/behavior settings
//
// Returns:
//   - *Response: A Response wrapper for making assertions
func (zt *T) PatchJSONWithConfig(app *fiber.App, path string, data any, cfg fiber.TestConfig) *Response {
	return executeJSON(zt.T, app, "PATCH", path, data, cfg)
}

// DeleteWithConfig performs a DELETE request with custom TestConfig.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users/1")
//   - cfg: TestConfig for custom timeout/behavior settings
//
// Returns:
//   - *Response: A Response wrapper for making assertions
func (zt *T) DeleteWithConfig(app *fiber.App, path string, cfg fiber.TestConfig) *Response {
	return execute(zt.T, app, "DELETE", path, nil, nil, cfg)
}

// DeleteJSONWithConfig performs a DELETE request with JSON body and custom TestConfig.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/api/users/bulk")
//   - data: Any data structure to be JSON-encoded
//   - cfg: TestConfig for custom timeout/behavior settings
//
// Returns:
//   - *Response: A Response wrapper for making assertions
func (zt *T) DeleteJSONWithConfig(app *fiber.App, path string, data any, cfg fiber.TestConfig) *Response {
	return executeJSON(zt.T, app, "DELETE", path, data, cfg)
}

// SetHeader creates a header map with a single key-value pair.
// This is a utility function to help construct header maps for custom headers
// when using raw body request methods.
//
// Parameters:
//   - key: The header name
//   - value: The header value
//
// Returns:
//   - map[string]string: A map containing the single header
//
// Example:
//
//	headers := zentests.SetHeader("Authorization", "Bearer token123")
//	// Use with a method that accepts headers parameter
func SetHeader(key, value string) map[string]string {
	return map[string]string{key: value}
}

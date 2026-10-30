package zentests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
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
// It builds the request from its parts and hands it to send, which runs it.
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

	return send(t, app, req, cfg)
}

// send runs a prepared request against the Fiber app and wraps the result in a Response.
// execute and Jar both delegate here, so every request is handled the same way.
// It stops the test (via require) if the request cannot be executed, for example on timeout.
func send(t testing.TB, app *fiber.App, req *http.Request, cfg fiber.TestConfig) *Response {
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
func executeJSON(t testing.TB, app *fiber.App, method, path string, data any, headers map[string]string, cfg fiber.TestConfig) *Response {
	jsonBody, err := json.Marshal(data)
	require.NoError(t, err, "JSON marshaling failed")

	allHeaders := maps.Clone(headers)
	if allHeaders == nil {
		allHeaders = map[string]string{}
	}
	allHeaders["Content-Type"] = "application/json"

	return execute(t, app, method, path, bytes.NewReader(jsonBody), allHeaders, cfg)
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
	return execute(zt.T, app, "GET", path, nil, zt.mergedHeaders(nil), cfg)
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
	return execute(zt.T, app, "POST", path, bytes.NewReader(body), zt.mergedHeaders(nil), cfg)
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
	return executeJSON(zt.T, app, "POST", path, data, zt.mergedHeaders(nil), cfg)
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
	headers := zt.mergedHeaders(SetHeader("Content-Type", "application/x-www-form-urlencoded"))
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
	return execute(zt.T, app, "PUT", path, bytes.NewReader(body), zt.mergedHeaders(nil), cfg)
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
	return executeJSON(zt.T, app, "PUT", path, data, zt.mergedHeaders(nil), cfg)
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
	return execute(zt.T, app, "PATCH", path, bytes.NewReader(body), zt.mergedHeaders(nil), cfg)
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
	return executeJSON(zt.T, app, "PATCH", path, data, zt.mergedHeaders(nil), cfg)
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
	return execute(zt.T, app, "DELETE", path, nil, zt.mergedHeaders(nil), cfg)
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
	return executeJSON(zt.T, app, "DELETE", path, data, zt.mergedHeaders(nil), cfg)
}

// FileUpload describes one file in a PostMultipart request.
type FileUpload struct {
	// Name is the form field name. The same name may appear in several entries.
	Name string
	// Filename is the file name presented to the server.
	Filename string
	// Content is the file's bytes.
	Content []byte
}

// multipartBody encodes fields and files as a multipart/form-data body and returns it
// with the Content-Type header value, which includes the boundary. Fields are written
// first in sorted key order, then files in slice order, so the body is deterministic.
func multipartBody(fields map[string]string, files []FileUpload) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for _, key := range slices.Sorted(maps.Keys(fields)) {
		if err := writer.WriteField(key, fields[key]); err != nil {
			return nil, "", fmt.Errorf("write form field %q: %w", key, err)
		}
	}
	for _, file := range files {
		part, err := writer.CreateFormFile(file.Name, file.Filename)
		if err != nil {
			return nil, "", fmt.Errorf("create form file %q: %w", file.Name, err)
		}
		if _, err := part.Write(file.Content); err != nil {
			return nil, "", fmt.Errorf("write file %q: %w", file.Filename, err)
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("close multipart writer: %w", err)
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}

// PostMultipart performs a POST request with multipart/form-data, typically for file
// uploads. The Content-Type header carries the multipart boundary.
//
// Parameters:
//   - app: The Fiber application to test against
//   - path: The request path (e.g., "/upload")
//   - fields: Plain form fields (nil when there are none)
//   - files: The files to upload (nil when there are none)
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	zt.PostMultipart(app, "/upload", map[string]string{"note": "hi"}, []zentests.FileUpload{
//	    {Name: "file", Filename: "a.txt", Content: []byte("hello")},
//	}).OK()
func (zt *T) PostMultipart(app *fiber.App, path string, fields map[string]string, files []FileUpload) *Response {
	zt.Helper()
	body, contentType, err := multipartBody(fields, files)
	require.NoError(zt.T, err, "building the multipart body failed")

	headers := zt.mergedHeaders(SetHeader("Content-Type", contentType))
	return execute(zt.T, app, "POST", path, bytes.NewReader(body), headers, fiber.TestConfig{})
}

// SetHeader sets a default header that is sent with every request made through zt.
// A header set for a single request, such as PostJSON's Content-Type, takes precedence
// over a default with the same name. Header names are canonicalised, so "x-api-key"
// and "X-Api-Key" configure the same default.
//
// Parameters:
//   - key: The header name (e.g., "Accept")
//   - value: The header value
//
// Returns:
//   - *T: The receiver for method chaining
//
// Example:
//
//	zt := zentests.New(t).Use(app).SetHeader("Accept", "application/json")
func (zt *T) SetHeader(key, value string) *T {
	zt.Helper()
	zt.headers = setDefaultHeader(zt.headers, key, value)
	return zt
}

// WithBearer sets the default Authorization header to "Bearer <token>", which is sent
// with every request made through zt.
//
// Parameters:
//   - token: The bearer token
//
// Returns:
//   - *T: The receiver for method chaining
//
// Example:
//
//	zt := zentests.New(t).Use(app).WithBearer(token)
func (zt *T) WithBearer(token string) *T {
	return zt.SetHeader("Authorization", "Bearer "+token)
}

// mergedHeaders returns the default headers with headers taking precedence on conflict.
// Returns headers unchanged when no defaults are set, so requests without defaults
// allocate nothing.
func (zt *T) mergedHeaders(headers map[string]string) map[string]string {
	if len(zt.headers) == 0 {
		return headers
	}
	merged := make(map[string]string, len(zt.headers)+len(headers))
	maps.Copy(merged, zt.headers)
	maps.Copy(merged, headers)
	return merged
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

// setDefaultHeader stores key-value in headers, creating the map on first use.
// Header names are canonicalised, so two spellings of a name cannot produce two
// defaults for the same header. Shared by T and Jar.
func setDefaultHeader(headers map[string]string, key, value string) map[string]string {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers[http.CanonicalHeaderKey(key)] = value
	return headers
}

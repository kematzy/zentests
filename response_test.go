package zentests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// TestResponse_Body tests the Body() method for lazy loading and caching.
func TestResponse_Body(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)

	assert.NoError(t, err)

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	// First call should read from response
	body1 := r.Body()
	assert.Equal(t, "Hello, World!", string(body1))
	assert.True(t, r.bodyRead)

	// Second call should return cached body
	body2 := r.Body()
	assert.Equal(t, body1, body2)
	assert.Same(t, &body1[0], &body2[0]) // Same underlying array
}

// TestResponse_BodyString tests the BodyString() method.
func TestResponse_BodyString(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Test Body Content")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	bodyStr := r.BodyString()
	assert.Equal(t, "Test Body Content", bodyStr)
}

// TestResponse_Dump tests the Dump() method for formatted output.
func TestResponse_Dump(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/plain")
		c.Status(200)
		return c.SendString("Dump Test Body")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	dump := r.Dump()

	// Verify dump contains expected information
	assert.Contains(t, dump, "Status: 200")
	assert.Contains(t, dump, "Headers:")
	assert.Contains(t, dump, "Content-Type")
	assert.Contains(t, dump, "Dump Test Body")
}

// TestResponse_Dump_WithJSON tests Dump() when JSON has been parsed.
func TestResponse_Dump_WithJSON(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"key": "value"})
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	// Parse JSON first
	r.JSON()

	dump := r.Dump()
	assert.Contains(t, dump, "Status: 200")
	assert.Contains(t, dump, `"key":"value"`)
}

// TestResponse_Body_MultipleCalls tests that Body() handles multiple calls correctly.
func TestResponse_Body_MultipleCalls(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Single Read Test")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	// Multiple calls should all return the same data
	for i := 0; i < 5; i++ {
		body := r.Body()
		assert.Equal(t, "Single Read Test", string(body))
	}
}

// TestResponse_EmptyBody tests handling of empty response body.
func TestResponse_EmptyBody(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(204) // No content
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	body := r.Body()
	assert.Empty(t, body)
	assert.Equal(t, "", r.BodyString())
}

// TestResponse_LargeBody tests handling of larger response bodies.
func TestResponse_LargeBody(t *testing.T) {
	largeContent := bytes.Repeat([]byte("Large content. "), 1000)

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.Send(largeContent)
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	body := r.Body()
	assert.Equal(t, largeContent, body)
	assert.Len(t, body, len(largeContent))
}

// TestResponse_BodyWithParsedJSON tests that Body() works correctly after JSON parsing.
func TestResponse_BodyWithParsedJSON(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "hello"})
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	// Parse JSON first - this should read the body
	r.JSON()

	// Body() should still return the raw JSON
	body := r.Body()
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	assert.NoError(t, err)
	assert.Equal(t, "hello", data["message"])
}

// TestResponse_DumpWithCustomHeaders tests Dump() with custom headers.
func TestResponse_DumpWithCustomHeaders(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		c.Set("X-Custom-Header", "custom-value")
		c.Set("X-Request-ID", "12345")
		return c.SendString("Custom Headers Test")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	dump := r.Dump()
	assert.Contains(t, dump, "X-Custom-Header")
	assert.Contains(t, dump, "X-Request-Id") // Go canonicalizes header names
	assert.Contains(t, dump, "Custom Headers Test")
}

// TestResponse_StructFields tests that Response struct fields are properly accessible.
func TestResponse_StructFields(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		c.Set("X-Test", "test-value")
		return c.SendStatus(201)
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	// Test exported fields
	assert.Equal(t, 201, r.StatusCode)
	assert.Equal(t, "test-value", r.Header.Get("X-Test"))
	assert.False(t, r.bodyRead)
	assert.Nil(t, r.parsedJSON)

	// After reading body
	r.Body()
	assert.True(t, r.bodyRead)
}

// TestResponse_BodyReadFlag tests the bodyRead flag behavior.
func TestResponse_BodyReadFlag(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Flag Test")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	assert.False(t, r.bodyRead)
	r.Body()
	assert.True(t, r.bodyRead)
	r.Body() // Second call should not change flag
	assert.True(t, r.bodyRead)
}

// TestResponse_WithHTTPResponse tests creating Response with raw http.Response.
func TestResponse_WithHTTPResponse(t *testing.T) {
	// Create a raw http.Response
	bodyContent := "Raw HTTP Response"
	resp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"text/plain"}},
		Body:       io.NopCloser(bytes.NewBufferString(bodyContent)),
	}

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	body := r.Body()
	assert.Equal(t, bodyContent, string(body))
}

// TestResponse_DumpFormatting tests the exact formatting of Dump() output.
func TestResponse_DumpFormatting(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.JSON(fiber.Map{"test": "data"})
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	dump := r.Dump()

	// Verify format: Status: N\nHeaders: {...}\nBody: ...
	lines := strings.Split(dump, "\n")
	assert.GreaterOrEqual(t, len(lines), 3)
	assert.True(t, strings.HasPrefix(lines[0], "Status: "))
	assert.True(t, strings.HasPrefix(lines[1], "Headers: "))
	assert.True(t, strings.HasPrefix(lines[2], "Body: "))
}

// TestResponse_BodyAfterBodyString tests that Body() returns same data after BodyString().
func TestResponse_BodyAfterBodyString(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Body After String Test")
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	r := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	// Call BodyString first
	str := r.BodyString()
	assert.Equal(t, "Body After String Test", str)

	// Body() should return same data
	body := r.Body()
	assert.Equal(t, str, string(body))
}

// TestResponse_Debug_WithParsedJSON tests Debug() output when JSON has been parsed.
// This covers the branch in Debug() that prints formatted JSON (lines 96-99).
func TestResponse_Debug_WithParsedJSON(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"test": "debug_json"})
	})

	req := httptest.NewRequest("GET", "/test", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	response := &Response{
		t:          t,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		resp:       resp,
		bodyRead:   false,
	}

	// Parse JSON first to populate parsedJSON
	response.JSON()
	assert.NotNil(t, response.parsedJSON)

	// Capture stdout
	oldStdout := os.Stdout
	pipeReader, pipeWriter, err := os.Pipe()
	assert.NoError(t, err)

	os.Stdout = pipeWriter

	// Call Debug() - should print formatted JSON section
	result := response.Debug()

	// Restore stdout
	pipeWriter.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	_, err = buf.ReadFrom(pipeReader)
	assert.NoErrorf(t, err, "failed to read from stdout buffer:  %s", err)
	output := buf.String()

	// Verify output contains formatted JSON
	assert.Equal(t, response, result)
	assert.Contains(t, output, "=== ZENTESTS DEBUG ===")
	assert.Contains(t, output, "Parsed JSON:")
	assert.Contains(t, output, "test")
	assert.Contains(t, output, "debug_json")
}

package zentests

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	app := fiber.New()

	app.Get("/test", func(c fiber.Ctx) error {
		return c.SendString("GET response")
	})

	app.Post("/test", func(c fiber.Ctx) error {
		// body, _ := io.ReadAll(c.Body())
		body := c.Body()
		// Fix: c.Body() returns []byte, not io.Reader
		return c.SendString("POST: " + string(body))
	})

	app.Post("/json", func(c fiber.Ctx) error {
		var data map[string]any
		if err := c.Bind().Body(&data); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"received": data})
	})

	app.Post("/form", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name": c.FormValue("name"),
			"age":  c.FormValue("age"),
		})
	})

	app.Put("/test", func(c fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Patch("/test", func(c fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Delete("/test", func(c fiber.Ctx) error {
		return c.SendStatus(204)
	})

	return app
}

func TestGet(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	resp := zt.Get(app, "/test")
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "GET response", resp.BodyString())
}

func TestPost(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	resp := zt.Post(app, "/test", []byte("hello"))
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "POST: hello", resp.BodyString())
}

func TestPostJSON(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	data := map[string]any{"name": "John", "age": 30}
	resp := zt.PostJSON(app, "/json", data)

	resp.OK().
		IsJSON().
		Has("received.name", "John").
		HasFloat("received.age", 30)
}

func TestPostForm(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	formData := map[string]string{"name": "Jane", "age": "25"}
	resp := zt.PostForm(app, "/form", formData)

	resp.OK().IsJSON().
		Has("name", "Jane").
		Has("age", "25")
}

func TestPut(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	resp := zt.Put(app, "/test", []byte("data"))
	resp.OK()
}

func TestPutJSON(t *testing.T) {
	app := fiber.New()
	app.Put("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"method": "PUT"})
	})

	zt := New(t)
	resp := zt.PutJSON(app, "/test", map[string]any{"key": "value"})

	resp.OK().IsJSON().Has("method", "PUT")
}

func TestPatch(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	resp := zt.Patch(app, "/test", []byte("patch"))
	resp.OK()
}

func TestPatchJSON(t *testing.T) {
	app := fiber.New()
	app.Patch("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"method": "PATCH"})
	})

	zt := New(t)
	resp := zt.PatchJSON(app, "/test", map[string]any{"key": "value"})

	resp.OK().IsJSON().Has("method", "PATCH")
}

func TestDelete(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	resp := zt.Delete(app, "/test")
	resp.NoContent()
}

func TestDeleteJSON(t *testing.T) {
	app := fiber.New()
	app.Delete("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"method": "DELETE"})
	})

	zt := New(t)
	resp := zt.DeleteJSON(app, "/test", map[string]any{"key": "value"})

	resp.OK().IsJSON().Has("method", "DELETE")
}

func TestSetHeader(t *testing.T) {
	headers := SetHeader("X-Custom", "value")
	assert.Equal(t, "value", headers["X-Custom"])
}

// =================================================================================================
// TESTS FOR *WithConfig METHODS (Fiber v3+)
// =================================================================================================

func TestGetWithConfig(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	// Test with empty config (uses defaults)
	resp := zt.GetWithConfig(app, "/test", fiber.TestConfig{})
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "GET response", resp.BodyString())

	// Test with custom timeout config
	resp2 := zt.GetWithConfig(app, "/test", fiber.TestConfig{
		Timeout: time.Second,
	})
	assert.Equal(t, 200, resp2.StatusCode)
}

func TestPostWithConfig(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	// Test with empty config
	resp := zt.PostWithConfig(app, "/test", []byte("hello"), fiber.TestConfig{})
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "POST: hello", resp.BodyString())

	// Test with custom timeout
	resp2 := zt.PostWithConfig(app, "/test", []byte("world"), fiber.TestConfig{
		Timeout: 2 * time.Second,
	})
	assert.Equal(t, 200, resp2.StatusCode)
	assert.Equal(t, "POST: world", resp2.BodyString())
}

func TestPostJSONWithConfig(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	data := map[string]any{"name": "John", "age": 30}

	// Test with empty config
	resp := zt.PostJSONWithConfig(app, "/json", data, fiber.TestConfig{})

	resp.OK().
		IsJSON().
		Has("received.name", "John").
		HasFloat("received.age", 30)

	// Test with custom config
	data2 := map[string]any{"name": "Jane", "age": 25}
	resp2 := zt.PostJSONWithConfig(app, "/json", data2, fiber.TestConfig{
		Timeout:       time.Second,
		FailOnTimeout: true,
	})

	resp2.OK().
		IsJSON().
		Has("received.name", "Jane").
		HasFloat("received.age", 25)
}

func TestPostFormWithConfig(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	formData := map[string]string{"name": "Jane", "age": "25"}

	// Test with empty config
	resp := zt.PostFormWithConfig(app, "/form", formData, fiber.TestConfig{})

	resp.OK().IsJSON().
		Has("name", "Jane").
		Has("age", "25")

	// Test with custom timeout
	formData2 := map[string]string{"name": "Bob", "age": "30"}
	resp2 := zt.PostFormWithConfig(app, "/form", formData2, fiber.TestConfig{
		Timeout: time.Second,
	})

	resp2.OK().IsJSON().
		Has("name", "Bob").
		Has("age", "30")
}

func TestPutWithConfig(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	// Test with empty config
	resp := zt.PutWithConfig(app, "/test", []byte("data"), fiber.TestConfig{})
	resp.OK()

	// Test with custom timeout
	resp2 := zt.PutWithConfig(app, "/test", []byte("more data"), fiber.TestConfig{
		Timeout: time.Second,
	})
	resp2.OK()
}

func TestPutJSONWithConfig(t *testing.T) {
	app := fiber.New()
	app.Put("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"method": "PUT"})
	})

	zt := New(t)

	// Test with empty config
	resp := zt.PutJSONWithConfig(app, "/test", map[string]any{"key": "value"}, fiber.TestConfig{})

	resp.OK().IsJSON().Has("method", "PUT")

	// Test with custom config
	resp2 := zt.PutJSONWithConfig(app, "/test", map[string]any{"key": "value2"}, fiber.TestConfig{
		Timeout:       time.Second,
		FailOnTimeout: true,
	})

	resp2.OK().IsJSON().Has("method", "PUT")
}

func TestPatchWithConfig(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	// Test with empty config
	resp := zt.PatchWithConfig(app, "/test", []byte("patch"), fiber.TestConfig{})
	resp.OK()

	// Test with custom timeout
	resp2 := zt.PatchWithConfig(app, "/test", []byte("patch2"), fiber.TestConfig{
		Timeout: time.Second,
	})
	resp2.OK()
}

func TestPatchJSONWithConfig(t *testing.T) {
	app := fiber.New()
	app.Patch("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"method": "PATCH"})
	})

	zt := New(t)

	// Test with empty config
	resp := zt.PatchJSONWithConfig(app, "/test", map[string]any{"key": "value"}, fiber.TestConfig{})

	resp.OK().IsJSON().Has("method", "PATCH")

	// Test with custom config
	resp2 := zt.PatchJSONWithConfig(app, "/test", map[string]any{"key": "value2"}, fiber.TestConfig{
		Timeout: time.Second,
	})

	resp2.OK().IsJSON().Has("method", "PATCH")
}

func TestDeleteWithConfig(t *testing.T) {
	app := setupTestApp()
	zt := New(t)

	// Test with empty config
	resp := zt.DeleteWithConfig(app, "/test", fiber.TestConfig{})
	resp.NoContent()

	// Test with custom timeout
	resp2 := zt.DeleteWithConfig(app, "/test", fiber.TestConfig{
		Timeout: time.Second,
	})
	resp2.NoContent()
}

func TestDeleteJSONWithConfig(t *testing.T) {
	app := fiber.New()
	app.Delete("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"method": "DELETE"})
	})

	zt := New(t)

	// Test with empty config
	resp := zt.DeleteJSONWithConfig(app, "/test", map[string]any{"key": "value"}, fiber.TestConfig{})

	resp.OK().IsJSON().Has("method", "DELETE")

	// Test with custom config
	resp2 := zt.DeleteJSONWithConfig(app, "/test", map[string]any{"key": "value2"}, fiber.TestConfig{
		Timeout:       time.Second,
		FailOnTimeout: true,
	})

	resp2.OK().IsJSON().Has("method", "DELETE")
}

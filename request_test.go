package zentests

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("GET response")
	})

	app.Post("/test", func(c *fiber.Ctx) error {
		// body, _ := io.ReadAll(c.Body())
		body := c.Body()
		// Fix: c.Body() returns []byte, not io.Reader
		return c.SendString("POST: " + string(body))
	})

	app.Post("/json", func(c *fiber.Ctx) error {
		var data map[string]any
		c.BodyParser(&data)
		return c.JSON(fiber.Map{"received": data})
	})

	app.Post("/form", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name": c.FormValue("name"),
			"age":  c.FormValue("age"),
		})
	})

	app.Put("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Patch("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Delete("/test", func(c *fiber.Ctx) error {
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
	app.Put("/test", func(c *fiber.Ctx) error {
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
	app.Patch("/test", func(c *fiber.Ctx) error {
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
	app.Delete("/test", func(c *fiber.Ctx) error {
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

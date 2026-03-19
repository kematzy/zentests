package zentests

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func setupBodyApp() *fiber.App {
	app := fiber.New()
	app.Get("/text", func(c fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	app.Get("/html", func(c fiber.Ctx) error {
		return c.SendString("<h1>Title</h1><p>Paragraph</p>")
	})
	app.Get("/empty", func(c fiber.Ctx) error {
		return c.SendStatus(204)
	})
	return app
}

func TestBodyAssertions(t *testing.T) {
	app := setupBodyApp()
	zt := New(t)

	zt.Get(app, "/text").
		Contains("Hello").
		Contains("World").
		NotContains("Goodbye")

	zt.Get(app, "/html").
		Contains("Title").
		BodyMatches(`<h1>\w+</h1>`).
		BodyMatches(`<p>.*</p>`)

	zt.Get(app, "/text").Equals("Hello World")
	zt.Get(app, "/empty").IsEmpty()
}

func TestBodyChaining(t *testing.T) {
	app := setupBodyApp()
	zt := New(t)

	resp := zt.Get(app, "/html")
	resp.Contains("Title").Contains("Paragraph")

	// Test that body is cached
	assert.Equal(t, resp.BodyString(), resp.BodyString())
	assert.Equal(t, resp.Body(), resp.Body())
}

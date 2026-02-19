package zentests

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func setupStatusApp() *fiber.App {
	app := fiber.New()
	app.Get("/200", func(c *fiber.Ctx) error { return c.SendStatus(200) })
	app.Get("/201", func(c *fiber.Ctx) error { return c.SendStatus(201) })
	app.Get("/202", func(c *fiber.Ctx) error { return c.SendStatus(202) })
	app.Get("/204", func(c *fiber.Ctx) error { return c.SendStatus(204) })
	app.Get("/400", func(c *fiber.Ctx) error { return c.SendStatus(400) })
	app.Get("/401", func(c *fiber.Ctx) error { return c.SendStatus(401) })
	app.Get("/403", func(c *fiber.Ctx) error { return c.SendStatus(403) })
	app.Get("/404", func(c *fiber.Ctx) error { return c.SendStatus(404) })
	app.Get("/422", func(c *fiber.Ctx) error { return c.SendStatus(422) })
	app.Get("/500", func(c *fiber.Ctx) error { return c.SendStatus(500) })
	app.Get("/custom", func(c *fiber.Ctx) error { return c.SendStatus(418) })
	return app
}

func TestStatusAssertions(t *testing.T) {
	app := setupStatusApp()
	zt := New(t)

	zt.Get(app, "/200").OK()
	zt.Get(app, "/201").Created()
	zt.Get(app, "/202").Accepted()
	zt.Get(app, "/204").NoContent()
	zt.Get(app, "/400").BadRequest()
	zt.Get(app, "/401").Unauthorized()
	zt.Get(app, "/403").Forbidden()
	zt.Get(app, "/404").NotFound()
	zt.Get(app, "/422").Unprocessable()
	zt.Get(app, "/500").ServerError()
	zt.Get(app, "/custom").Status(418)
}

func TestHeaderAssertions(t *testing.T) {
	app := fiber.New()
	app.Get("/json", func(c *fiber.Ctx) error {
		c.Type("json")
		return c.JSON(fiber.Map{"key": "value"})
	})
	app.Get("/html", func(c *fiber.Ctx) error {
		c.Type("html")
		return c.SendString("<h1>Test</h1>")
	})
	app.Get("/text", func(c *fiber.Ctx) error {
		c.Type("txt")
		return c.SendString("plain text")
	})
	app.Get("/custom-header", func(c *fiber.Ctx) error {
		c.Set("X-Custom", "custom-value")
		return c.SendStatus(200)
	})
	app.Get("/cookies", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{Name: "lang", Value: "en"})
		c.Cookie(&fiber.Cookie{Name: "theme", Value: "dark"})
		return c.SendStatus(200)
	})

	zt := New(t)

	zt.Get(app, "/json").HeaderPresent("Content-Type")
	zt.Get(app, "/json").HeaderNotPresent("Zen-Header")

	zt.Get(app, "/cookies").
		HeaderHasValues(
			"Set-Cookie",
			[]string{
				"lang=en; path=/; SameSite=Lax",
				"theme=dark; path=/; SameSite=Lax",
			},
		)
	zt.Get(app, "/cookies").
		CookieHasValues(map[string]string{
			"lang":  "en",
			"theme": "dark",
		})
	zt.Get(app, "/json").HasContentType("application/json")
	zt.Get(app, "/json").IsJSON()
	zt.Get(app, "/html").IsHTML()
	zt.Get(app, "/text").IsPlainText()
	zt.Get(app, "/custom-header").
		HasHeader("X-Custom", "custom-value").
		HeaderContains("X-Custom", "custom")
}

func TestContentTypeShortcuts(t *testing.T) {
	app := fiber.New()

	// CSS
	app.Get("/style.css", func(c *fiber.Ctx) error {
		c.Type("css")
		return c.SendString("body { color: red; }")
	})

	// JavaScript
	app.Get("/app.js", func(c *fiber.Ctx) error {
		c.Type("js")
		return c.SendString("console.log('hello');")
	})

	// XML
	app.Get("/data.xml", func(c *fiber.Ctx) error {
		c.Type("xml")
		return c.SendString("<?xml version=\"1.0\"?><root><item>value</item></root>")
	})

	// XHR
	app.Get("/xhr", func(c *fiber.Ctx) error {
		c.Set("X-Requested-With", "XMLHttpRequest")
		return c.JSON(fiber.Map{"xhr": true})
	})

	// Images
	app.Get("/photo.png", func(c *fiber.Ctx) error {
		c.Type("png")
		return c.SendStatus(200)
	})
	app.Get("/photo.jpg", func(c *fiber.Ctx) error {
		c.Type("jpg")
		return c.SendStatus(200)
	})
	app.Get("/animation.gif", func(c *fiber.Ctx) error {
		c.Type("gif")
		return c.SendStatus(200)
	})
	app.Get("/icon.svg", func(c *fiber.Ctx) error {
		c.Type("svg")
		return c.SendStatus(200)
	})
	app.Get("/photo.webp", func(c *fiber.Ctx) error {
		c.Type("webp")
		return c.SendStatus(200)
	})
	app.Get("/image", func(c *fiber.Ctx) error {
		c.Type("png") // Any image type should match IsImage()
		return c.SendStatus(200)
	})

	zt := New(t)

	// Test content type assertions
	zt.Get(app, "/style.css").OK().IsCSS()
	zt.Get(app, "/app.js").OK().IsJS()
	zt.Get(app, "/data.xml").OK().IsXML()
	zt.Get(app, "/xhr").OK().IsXHR()

	// Test image assertions
	zt.Get(app, "/photo.png").OK().IsPNG().IsImage()
	zt.Get(app, "/photo.jpg").OK().IsJPEG().IsImage()
	zt.Get(app, "/animation.gif").OK().IsGIF().IsImage()
	zt.Get(app, "/icon.svg").OK().IsSVG().IsImage()
	zt.Get(app, "/photo.webp").OK().IsWebP().IsImage()
	zt.Get(app, "/image").OK().IsImage()
}

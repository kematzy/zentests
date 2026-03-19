package zentests

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/suite"
)

type HTTPRequestsSuite struct {
	suite.Suite
	app *fiber.App
	zt  *T
}

func (s *HTTPRequestsSuite) SetupTest() {
	s.app = fiber.New()
	s.zt = New(s.T())

	s.app.Get("/test", func(c fiber.Ctx) error {
		return c.SendString("GET response")
	})

	s.app.Post("/test", func(c fiber.Ctx) error {
		body := c.Body()
		return c.SendString("POST: " + string(body))
	})

	s.app.Post("/json", func(c fiber.Ctx) error {
		var data map[string]any
		if err := c.Bind().Body(&data); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"received": data})
	})

	s.app.Post("/form", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name": c.FormValue("name"),
			"age":  c.FormValue("age"),
		})
	})

	s.app.Put("/test", func(c fiber.Ctx) error {
		return c.SendStatus(200)
	})

	s.app.Patch("/test", func(c fiber.Ctx) error {
		return c.SendStatus(200)
	})

	s.app.Delete("/test", func(c fiber.Ctx) error {
		return c.SendStatus(204)
	})

	s.app.Put("/json", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"method": "PUT"})
	})

	s.app.Patch("/json", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"method": "PATCH"})
	})

	s.app.Delete("/json", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"method": "DELETE"})
	})
}

func (s *HTTPRequestsSuite) TearDownTest() {
	if s.app != nil {
		_ = s.app.Shutdown()
		s.app = nil
	}
}

// --- GET ---

func (s *HTTPRequestsSuite) Test_HTTP_GET() {
	resp := s.zt.Get(s.app, "/test")

	s.Equal(200, resp.StatusCode)
	s.Equal("GET response", resp.BodyString())
}

// --- POST ---

func (s *HTTPRequestsSuite) Test_HTTP_POST() {
	s.zt.Post(s.app, "/test", []byte("hello")).OK().Contains("POST: hello")

	resp := s.zt.Post(s.app, "/test", []byte("hello"))
	s.Equal(200, resp.StatusCode)
	s.Equal("POST: hello", resp.BodyString())
}

func (s *HTTPRequestsSuite) Test_HTTP_POST_JSON() {
	data := map[string]any{"name": "John", "age": 30}

	s.zt.PostJSON(s.app, "/json", data).OK().JSON().
		Has("received.name", "John").
		HasFloat("received.age", 30)
}

func (s *HTTPRequestsSuite) Test_HTTP_POST_Form() {
	formData := map[string]string{"name": "Jane", "age": "25"}

	s.zt.PostForm(s.app, "/form", formData).OK().JSON().Has("name", "Jane").Has("age", "25")
}

// --- PUT ---

func (s *HTTPRequestsSuite) Test_HTTP_PUT() {
	s.zt.Put(s.app, "/test", []byte("data")).OK()
}

func (s *HTTPRequestsSuite) Test_HTTP_PUT_JSON() {
	s.zt.PutJSON(s.app, "/json", map[string]any{"key": "value"}).OK().JSON().Has("method", "PUT")
}

// --- PATCH ---

func (s *HTTPRequestsSuite) Test_HTTP_PATCH() {
	s.zt.Patch(s.app, "/test", []byte("patch")).OK()
}

func (s *HTTPRequestsSuite) Test_HTTP_PATCH_JSON() {
	s.zt.PatchJSON(s.app, "/json", map[string]any{"key": "value"}).OK().JSON().Has("method", "PATCH")
}

// --- DELETE ---

func (s *HTTPRequestsSuite) Test_HTTP_DELETE() {
	s.zt.Delete(s.app, "/test").NoContent()
}

func (s *HTTPRequestsSuite) Test_HTTP_DELETE_JSON() {
	s.zt.DeleteJSON(s.app, "/json", map[string]any{"key": "value"}).OK().JSON().Has("method", "DELETE")
}

// --- Helpers ---

func (s *HTTPRequestsSuite) Test_HTTP_SetHeader() {
	headers := SetHeader("X-Custom", "value")
	s.Equal("value", headers["X-Custom"])
}

// =================================================================================================
// TESTS FOR *WithConfig METHODS (Fiber v3+)
// =================================================================================================

// --- GET WithConfig ---

func (s *HTTPRequestsSuite) Test_HTTP_GET_WithConfig_Empty() {
	s.zt.GetWithConfig(s.app, "/test", fiber.TestConfig{}).OK().Equals("GET response")
}

func (s *HTTPRequestsSuite) Test_HTTP_GET_WithConfig_Timeout() {
	s.zt.GetWithConfig(s.app, "/test", fiber.TestConfig{Timeout: time.Second}).Status(200)
}

// --- POST WithConfig ---

func (s *HTTPRequestsSuite) Test_HTTP_PostWithConfig() {
	s.zt.PostWithConfig(s.app, "/test", []byte("hello"), fiber.TestConfig{}).Status(200).Equals("POST: hello")
}

func (s *HTTPRequestsSuite) Test_HTTP_PostWithConfig_Timeout() {
	s.zt.PostWithConfig(s.app, "/test", []byte("world"), fiber.TestConfig{Timeout: 2 * time.Second}).
		Status(200).Equals("POST: world")
}

func (s *HTTPRequestsSuite) Test_HTTP_PostJSONWithConfig_Empty() {
	data := map[string]any{"name": "John", "age": 30}

	s.zt.PostJSONWithConfig(s.app, "/json", data, fiber.TestConfig{}).OK().JSON().
		Has("received.name", "John").HasFloat("received.age", 30)
}

func (s *HTTPRequestsSuite) Test_HTTP_PostJSONWithConfig_Timeout() {
	data := map[string]any{"name": "Jane", "age": 25}

	s.zt.PostJSONWithConfig(s.app, "/json", data, fiber.TestConfig{
		Timeout:       time.Second,
		FailOnTimeout: true,
	}).OK().JSON().Has("received.name", "Jane").HasFloat("received.age", 25)
}

func (s *HTTPRequestsSuite) Test_HTTP_PostFormWithConfig_Empty() {
	formData := map[string]string{"name": "Jane", "age": "25"}

	s.zt.PostFormWithConfig(s.app, "/form", formData, fiber.TestConfig{}).OK().JSON().
		Has("name", "Jane").Has("age", "25")
}

func (s *HTTPRequestsSuite) Test_HTTP_PostFormWithConfig_Timeout() {
	formData := map[string]string{"name": "Bob", "age": "30"}
	s.zt.PostFormWithConfig(s.app, "/form", formData, fiber.TestConfig{Timeout: time.Second}).OK().JSON().
		Has("name", "Bob").Has("age", "30")
}

// --- PUT WithConfig ---

func (s *HTTPRequestsSuite) Test_HTTP_PutWithConfig() {
	s.zt.PutWithConfig(s.app, "/test", []byte("data"), fiber.TestConfig{}).OK()
}

func (s *HTTPRequestsSuite) Test_HTTP_PutWithConfig_Timeout() {
	s.zt.PutWithConfig(s.app, "/test", []byte("more data"), fiber.TestConfig{Timeout: time.Second}).OK()
}

func (s *HTTPRequestsSuite) Test_HTTP_PutJSONWithConfig_Empty() {
	s.zt.PutJSONWithConfig(s.app, "/json", map[string]any{"key": "value"}, fiber.TestConfig{}).
		OK().IsJSON().Has("method", "PUT")
}

func (s *HTTPRequestsSuite) Test_HTTP_PutJSONWithConfig_Timeout() {
	s.zt.PutJSONWithConfig(s.app, "/json", map[string]any{"key": "value2"}, fiber.TestConfig{
		Timeout:       time.Second,
		FailOnTimeout: true,
	}).OK().IsJSON().Has("method", "PUT")
}

// --- PATCH WithConfig ---

func (s *HTTPRequestsSuite) Test_HTTP_PATCH_WithConfig_Empty() {
	s.zt.PatchWithConfig(s.app, "/test", []byte("patch"), fiber.TestConfig{}).OK()
}

func (s *HTTPRequestsSuite) Test_HTTP_PATCH_WithConfig_Timeout() {
	s.zt.PatchWithConfig(s.app, "/test", []byte("patch2"), fiber.TestConfig{Timeout: time.Second}).OK()
}

func (s *HTTPRequestsSuite) Test_HTTP_PATCH_JSON_WithConfig_Empty() {
	s.zt.PatchJSONWithConfig(s.app, "/json", map[string]any{"key": "value"}, fiber.TestConfig{}).
		OK().IsJSON().Has("method", "PATCH")
}

func (s *HTTPRequestsSuite) Test_HTTP_PATCH_JSON_WithConfig_Timeout() {
	s.zt.PatchJSONWithConfig(s.app, "/json", map[string]any{"key": "value2"}, fiber.TestConfig{Timeout: time.Second}).
		OK().IsJSON().Has("method", "PATCH")
}

// --- DELETE WithConfig ---

func (s *HTTPRequestsSuite) Test_HTTP_DELETE_WithConfig_Empty() {
	s.zt.DeleteWithConfig(s.app, "/test", fiber.TestConfig{}).NoContent()
}

func (s *HTTPRequestsSuite) Test_HTTP_DELETE_WithConfig_Timeout() {
	s.zt.DeleteWithConfig(s.app, "/test", fiber.TestConfig{Timeout: time.Second}).NoContent()
}

func (s *HTTPRequestsSuite) Test_HTTP_DELETE_JSON_WithConfig_Empty() {
	s.zt.DeleteJSONWithConfig(s.app, "/json", map[string]any{"key": "value"}, fiber.TestConfig{}).
		OK().IsJSON().Has("method", "DELETE")
}

func (s *HTTPRequestsSuite) Test_HTTP_DELETE_JSON_WithConfig_Timeout() {
	s.zt.DeleteJSONWithConfig(s.app, "/json", map[string]any{"key": "value2"}, fiber.TestConfig{
		Timeout:       time.Second,
		FailOnTimeout: true,
	}).OK().IsJSON().Has("method", "DELETE")
}

func TestHTTPRequestsSuite(t *testing.T) {
	suite.Run(t, new(HTTPRequestsSuite))
}

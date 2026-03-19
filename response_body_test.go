package zentests

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/suite"
)

type ResponseBodyAssertionsSuite struct {
	suite.Suite
	app *fiber.App
	zt  *T
}

func (s *ResponseBodyAssertionsSuite) SetupTest() {
	s.app = fiber.New()
	s.zt = New(s.T())

	s.app.Get("/text", func(c fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	s.app.Get("/html", func(c fiber.Ctx) error {
		return c.SendString("<h1>Title</h1><p>Paragraph</p>")
	})
	s.app.Get("/empty", func(c fiber.Ctx) error {
		return c.SendStatus(204)
	})
}

func (s *ResponseBodyAssertionsSuite) TearDownTest() {
	if s.app != nil {
		_ = s.app.Shutdown()
		s.app = nil
	}
}

func (s *ResponseBodyAssertionsSuite) Test_BodyAssertions_Contains() {
	s.zt.Get(s.app, "/text").OK().Contains("Hello World")
}

func (s *ResponseBodyAssertionsSuite) Test_BodyAssertions_NotContains() {
	s.zt.Get(s.app, "/text").OK().NotContains("Hello World!") // Note `!`
}

func (s *ResponseBodyAssertionsSuite) Test_BodyAssertions_BodyMatches() {
	s.zt.Get(s.app, "/html").OK().
		BodyMatches(`<h1>\w+</h1>`).
		BodyMatches(`<p>.*</p>`)
}

func (s *ResponseBodyAssertionsSuite) Test_BodyAssertions_Equals() {
	s.zt.Get(s.app, "/text").OK().Equals("Hello World")
}

func (s *ResponseBodyAssertionsSuite) Test_BodyAssertions_Empty() {
	s.zt.Get(s.app, "/empty").Status(204).IsEmpty()
}

func (s *ResponseBodyAssertionsSuite) Test_BodyAssertions_Contains_Chaining() {
	resp := s.zt.Get(s.app, "/html")
	resp.Contains("Title").Contains("Paragraph")
}

func TestResponseBodyAssertionsSuite(t *testing.T) {
	suite.Run(t, new(ResponseBodyAssertionsSuite))
}

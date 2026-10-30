package zentests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/suite"
)

type JarSuite struct {
	suite.Suite
	app *fiber.App
}

func TestJarSuite(t *testing.T) {
	suite.Run(t, new(JarSuite))
}

func (s *JarSuite) SetupTest() {
	s.app = fiber.New()

	// Echoes the "X-Test" header and the "sid" cookie, comma-separated.
	s.app.Get("/hdr", func(c fiber.Ctx) error {
		return c.SendString(c.Get("X-Test") + "," + c.Cookies("sid"))
	})

	// Echoes the request's Authorization header.
	s.app.Get("/auth", func(c fiber.Ctx) error {
		return c.SendString(c.Get("Authorization"))
	})

	// Echoes the request's Content-Type header.
	s.app.Post("/ctype", func(c fiber.Ctx) error {
		return c.SendString(c.Get("Content-Type"))
	})

	// Echoes the raw Cookie header the request carried, or "" when it carried none.
	// Raw, so a cookie sent with an empty value stays distinct from one not sent at all.
	s.app.Get("/echo", func(c fiber.Ctx) error {
		return c.SendString(c.Get(fiber.HeaderCookie))
	})

	// Sets the "sid" cookie from the query string, e.g. /set?sid=abc.
	s.app.Get("/set", func(c fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{Name: "sid", Value: c.Query("sid")})
		return c.SendString("set")
	})

	// Sets a session cookie with attributes so Response.Cookie can read them.
	s.app.Get("/login", func(c fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{Name: "session", Value: "s1", HTTPOnly: true, Path: "/account"})
		return c.SendString("logged in")
	})

	// Expires the "session" cookie by MaxAge -1, keeping its value, so the jar deletes it on MaxAge alone.
	s.app.Get("/expire", func(c fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{Name: "session", Value: "s1", MaxAge: -1})
		return c.SendString("expired")
	})

	// Expires the "session" cookie with Fiber's ClearCookie, which sends an empty value.
	s.app.Get("/logout", func(c fiber.Ctx) error {
		c.ClearCookie("session")
		return c.SendString("bye")
	})

	// Sets "sid" to an empty value, which the jar treats as a deletion.
	s.app.Get("/blank", func(c fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{Name: "sid", Value: ""})
		return c.SendString("blank")
	})

	// Echoes the "name" form field.
	s.app.Post("/form", func(c fiber.Ctx) error {
		return c.SendString("name=" + c.FormValue("name"))
	})

	// Echoes the "a" and "b" cookies.
	s.app.Get("/pair", func(c fiber.Ctx) error {
		return c.SendString("a=" + c.Cookies("a") + ",b=" + c.Cookies("b"))
	})

	// Sets the "a" and "b" cookies in one response.
	s.app.Get("/set-pair", func(c fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{Name: "a", Value: "1"})
		c.Cookie(&fiber.Cookie{Name: "b", Value: "2"})
		return c.SendString("pair")
	})
}

// Test_Jar_StoresCookieAndSendsItOnNextRequest checks that a cookie set by one response is sent on the next.
func (s *JarSuite) Test_Jar_StoresCookieAndSendsItOnNextRequest() {
	jar := NewJar(s.T(), s.app)

	jar.Get("/set?sid=abc").OK()
	jar.Get("/echo").OK().Equals("sid=abc")
}

// Test_Jar_MaxAgeExpiredCookieIsForgotten checks that a cookie set with MaxAge -1 is removed from the jar.
func (s *JarSuite) Test_Jar_MaxAgeExpiredCookieIsForgotten() {
	jar := NewJar(s.T(), s.app)

	jar.Get("/login").OK()
	s.Equal("s1", jar.Cookie("session"))

	jar.Get("/expire").OK()
	s.Empty(jar.Cookie("session"))
	jar.Get("/echo").OK().Equals("")
}

// Test_Jar_ClearCookieIsForgotten checks that Fiber's ClearCookie removes the cookie from the jar.
// ClearCookie sends an empty value rather than MaxAge -1, so this exercises the empty-value rule.
func (s *JarSuite) Test_Jar_ClearCookieIsForgotten() {
	jar := NewJar(s.T(), s.app)

	jar.Get("/login").OK()
	jar.Get("/logout").OK()
	s.Empty(jar.Cookie("session"))
	jar.Get("/echo").OK().Equals("")
}

// Test_Jar_EmptyValueCookieIsForgotten checks that a Set-Cookie with an empty value removes the cookie.
func (s *JarSuite) Test_Jar_EmptyValueCookieIsForgotten() {
	jar := NewJar(s.T(), s.app)

	jar.Get("/set?sid=abc").OK()
	jar.Get("/blank").OK()

	s.Empty(jar.Cookie("sid"))
	jar.Get("/echo").OK().Equals("")
}

// Test_Jar_Cookie_MissingNameIsEmpty checks that an unknown cookie reads as "" without failing the test.
func (s *JarSuite) Test_Jar_Cookie_MissingNameIsEmpty() {
	jar := NewJar(s.T(), s.app)

	s.Empty(jar.Cookie("missing"))
}

// Test_Jar_Get_ReturnsResponseUsableForAssertions checks that Get returns a Response that chains assertions.
func (s *JarSuite) Test_Jar_Get_ReturnsResponseUsableForAssertions() {
	jar := NewJar(s.T(), s.app)

	jar.Get("/login").OK().IsPlainText().Contains("logged in")
}

// Test_Jar_PostForm_SendsUrlencodedBody checks that PostForm sends an escaped form body the handler can read.
func (s *JarSuite) Test_Jar_PostForm_SendsUrlencodedBody() {
	jar := NewJar(s.T(), s.app)

	jar.PostForm("/form", url.Values{"name": {"Jane & Co = 1+1"}}).OK().Equals("name=Jane & Co = 1+1")
}

// Test_Jar_Send_CarriesStoredCookies checks that a caller-built request is sent with the stored cookies.
func (s *JarSuite) Test_Jar_Send_CarriesStoredCookies() {
	jar := NewJar(s.T(), s.app)
	jar.Get("/set?sid=xyz").OK()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/echo", http.NoBody)
	jar.Send(req).OK().Equals("sid=xyz")
}

// Test_Jar_Send_KeepsCallerCookies checks that Send keeps a cookie the caller set on the request.
func (s *JarSuite) Test_Jar_Send_KeepsCallerCookies() {
	jar := NewJar(s.T(), s.app)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/echo", http.NoBody)
	req.AddCookie(&http.Cookie{Name: "sid", Value: "caller"})
	jar.Send(req).OK().Equals("sid=caller")
}

// Test_Jar_TwoCookiesAreSentIndependently checks that two cookies set together are both sent back.
func (s *JarSuite) Test_Jar_TwoCookiesAreSentIndependently() {
	jar := NewJar(s.T(), s.app)

	jar.Get("/set-pair").OK()
	jar.Get("/pair").OK().Equals("a=1,b=2")
}

// Test_Jar_ResetCookieOverwritesStoredValue checks that a later Set-Cookie replaces the stored value.
func (s *JarSuite) Test_Jar_ResetCookieOverwritesStoredValue() {
	jar := NewJar(s.T(), s.app)

	jar.Get("/set?sid=first").OK()
	jar.Get("/set?sid=second").OK()

	s.Equal("second", jar.Cookie("sid"))
	jar.Get("/echo").OK().Equals("sid=second")
}

// Test_Response_Cookie_ReturnsAttributes checks that Response.Cookie exposes HttpOnly and Path.
func (s *JarSuite) Test_Response_Cookie_ReturnsAttributes() {
	jar := NewJar(s.T(), s.app)

	cookie := jar.Get("/login").OK().Cookie("session")
	s.Equal("s1", cookie.Value)
	s.True(cookie.HttpOnly)
	s.Equal("/account", cookie.Path)
}

// Test_Response_Cookie_MissingFailsTest checks that Response.Cookie fails the test when the cookie was not set.
func (s *JarSuite) Test_Response_Cookie_MissingFailsTest() {
	errs := failuresOf(s.T(), sendString("no cookies"), func(r *Response) {
		r.Cookie("missing")
	})

	s.Require().Len(errs, 1)
	s.Contains(errs[0], "cookie missing was not set")
}

// Test_NewJar_NilAppFailsTest checks that NewJar fails the test when app is nil, like Use(nil).
func (s *JarSuite) Test_NewJar_NilAppFailsTest() {
	rec := runRecorded(s.T(), func(tb testing.TB) {
		NewJar(tb, nil)
	})

	s.True(rec.fatal)
	s.Require().Len(rec.errors, 1)
	s.Contains(rec.errors[0], "nil *fiber.App")
}

// Test_Jar_SetHeader_SentWithEveryRequest checks that a default header reaches every jar
// request, alongside the stored cookies.
func (s *JarSuite) Test_Jar_SetHeader_SentWithEveryRequest() {
	jar := NewJar(s.T(), s.app)
	jar.Get("/set?sid=abc").OK()
	jar.SetHeader("X-Test", "1")

	jar.Get("/hdr").OK().Equals("1,abc")
	jar.Get("/hdr").OK().Equals("1,abc")
}

// Test_Jar_WithBearer_SendsAuthorizationHeader checks that WithBearer sets the
// Authorization header on jar requests.
func (s *JarSuite) Test_Jar_WithBearer_SendsAuthorizationHeader() {
	jar := NewJar(s.T(), s.app)
	jar.WithBearer("tok-1")

	jar.Get("/auth").OK().Equals("Bearer tok-1")
}

// Test_Jar_SetHeader_PerCallHeaderWins checks that a request's own Content-Type is not
// replaced by a default with the same name.
func (s *JarSuite) Test_Jar_SetHeader_PerCallHeaderWins() {
	jar := NewJar(s.T(), s.app)
	jar.SetHeader("Content-Type", "text/plain")

	jar.PostForm("/ctype", url.Values{"name": {"Jane"}}).OK().Equals("application/x-www-form-urlencoded")
}

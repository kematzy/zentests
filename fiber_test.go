package zentests

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/suite"
)

// noop is a minimal Fiber handler used across test helpers to register routes
// without any meaningful logic.
func noop(c fiber.Ctx) error { return nil }

// --- hasRoute (internal helper) ----------------------------------------------------------------

type HasRouteInternalSuite struct {
	suite.Suite
	app *fiber.App
}

func (s *HasRouteInternalSuite) SetupTest() {
	s.app = fiber.New()
	s.app.Get("/users", noop)
	s.app.Post("/users", noop)
	s.app.Put("/users/:id", noop)
	s.app.Patch("/users/:id", noop)
	s.app.Delete("/users/:id", noop)
}

func (s *HasRouteInternalSuite) TearDownTest() {
	if s.app != nil {
		_ = s.app.Shutdown()
		s.app = nil
	}
}

func (s *HasRouteInternalSuite) Test_hasRoute_Returns_True_For_Registered_GET() {
	s.True(hasRoute(s.app, "GET", "/users"))
}

func (s *HasRouteInternalSuite) Test_hasRoute_Returns_True_For_Registered_POST() {
	s.True(hasRoute(s.app, "POST", "/users"))
}

func (s *HasRouteInternalSuite) Test_hasRoute_Returns_True_For_Registered_PUT() {
	s.True(hasRoute(s.app, "PUT", "/users/:id"))
}

func (s *HasRouteInternalSuite) Test_hasRoute_Returns_True_For_Registered_PATCH() {
	s.True(hasRoute(s.app, "PATCH", "/users/:id"))
}

func (s *HasRouteInternalSuite) Test_hasRoute_Returns_True_For_Registered_DELETE() {
	s.True(hasRoute(s.app, "DELETE", "/users/:id"))
}

func (s *HasRouteInternalSuite) Test_hasRoute_Returns_False_For_Unregistered_Path() {
	s.False(hasRoute(s.app, "GET", "/nonexistent"))
}

func (s *HasRouteInternalSuite) Test_hasRoute_Returns_False_For_Wrong_Method() {
	// /users is only GET and POST — PATCH should not match
	s.False(hasRoute(s.app, "PATCH", "/users"))
}

func (s *HasRouteInternalSuite) Test_hasRoute_CaseInsensitive_Method() {
	// lowercase "get" should match the uppercase "GET" stored in the stack
	s.True(hasRoute(s.app, "get", "/users"))
	s.True(hasRoute(s.app, "Get", "/users"))
	s.True(hasRoute(s.app, "gEt", "/users"))
}

func (s *HasRouteInternalSuite) Test_hasRoute_Does_Not_Match_Param_Differently() {
	// /users/:id is registered; /users/42 is a request path, not a route path
	s.False(hasRoute(s.app, "PUT", "/users/42"))
}

func (s *HasRouteInternalSuite) Test_hasRoute_EmptyApp_Returns_False() {
	emptyApp := fiber.New()
	defer emptyApp.Shutdown() //nolint:errcheck
	s.False(hasRoute(emptyApp, "GET", "/"))
}

func (s *HasRouteInternalSuite) Test_hasRoute_Root_Path() {
	s.app.Get("/", noop)
	s.True(hasRoute(s.app, "GET", "/"))
}

func TestHasRouteInternalSuite(t *testing.T) {
	suite.Run(t, new(HasRouteInternalSuite))
}

// --- HasRoute / NotHasRoute (exported methods on *T) ------------------------------------------

type HasRouteAssertionSuite struct {
	suite.Suite
	app *fiber.App
	zt  *T
}

func (s *HasRouteAssertionSuite) SetupTest() {
	s.app = fiber.New()
	s.zt = New(s.T())

	s.app.Get("/ping", noop)
	s.app.Post("/items", noop)
	s.app.Delete("/items/:id", noop)
}

func (s *HasRouteAssertionSuite) TearDownTest() {
	if s.app != nil {
		_ = s.app.Shutdown()
		s.app = nil
	}
}

// --- HasRoute ----------------------------------------------------------------------------------

func (s *HasRouteAssertionSuite) Test_HasRoute_Passes_For_Registered_Route() {
	// Must not fail the test
	s.zt.HasRoute("GET", "/ping", s.app)
}

func (s *HasRouteAssertionSuite) Test_HasRoute_Returns_Self_For_Chaining() {
	result := s.zt.HasRoute("GET", "/ping", s.app)
	s.Same(s.zt, result)
}

func (s *HasRouteAssertionSuite) Test_HasRoute_Chaining_Multiple_Routes() {
	s.zt.
		HasRoute("GET", "/ping", s.app).
		HasRoute("POST", "/items", s.app).
		HasRoute("DELETE", "/items/:id", s.app)
}

func (s *HasRouteAssertionSuite) Test_HasRoute_Accepts_Optional_Message() {
	s.zt.HasRoute("GET", "/ping", s.app, "ping route must be registered")
}

func (s *HasRouteAssertionSuite) Test_HasRoute_CaseInsensitive_Method() {
	s.zt.HasRoute("get", "/ping", s.app)
	s.zt.HasRoute("Get", "/ping", s.app)
}

// --- NotHasRoute -------------------------------------------------------------------------------

func (s *HasRouteAssertionSuite) Test_NotHasRoute_Passes_For_Missing_Route() {
	s.zt.NotHasRoute("GET", "/nonexistent", s.app)
}

func (s *HasRouteAssertionSuite) Test_NotHasRoute_Returns_Self_For_Chaining() {
	result := s.zt.NotHasRoute("PUT", "/ping", s.app)
	s.Same(s.zt, result)
}

func (s *HasRouteAssertionSuite) Test_NotHasRoute_WrongMethod_Passes() {
	// /ping is only GET; PATCH /ping should not exist
	s.zt.NotHasRoute("PATCH", "/ping", s.app)
}

func (s *HasRouteAssertionSuite) Test_NotHasRoute_WrongPath_Passes() {
	s.zt.NotHasRoute("GET", "/items", s.app)
}

func (s *HasRouteAssertionSuite) Test_NotHasRoute_Accepts_Optional_Message() {
	s.zt.NotHasRoute("DELETE", "/ping", s.app, "DELETE /ping must not be exposed")
}

func TestHasRouteAssertionSuite(t *testing.T) {
	suite.Run(t, new(HasRouteAssertionSuite))
}

// --- NewApp ------------------------------------------------------------------------------------

type NewAppSuite struct {
	suite.Suite
}

func (s *NewAppSuite) Test_NewApp_Returns_Non_Nil_App() {
	app := NewApp(s.T())
	s.NotNil(app)
}

func (s *NewAppSuite) Test_NewApp_Each_Call_Returns_Fresh_Instance() {
	app1 := NewApp(s.T())
	app2 := NewApp(s.T())
	s.NotSame(app1, app2)
}

func (s *NewAppSuite) Test_NewApp_App_Accepts_Route_Registration() {
	app := NewApp(s.T())
	app.Get("/hello", noop)
	s.True(hasRoute(app, "GET", "/hello"))
}

func (s *NewAppSuite) Test_NewApp_App_Is_Functional_For_Requests() {
	app := NewApp(s.T())
	app.Get("/ok", func(c fiber.Ctx) error {
		return c.SendString("hello")
	})

	zt := New(s.T())
	zt.Get(app, "/ok").OK().Contains("hello")
}

func (s *NewAppSuite) Test_NewApp_Empty_App_Has_No_Routes() {
	app := NewApp(s.T())
	// No routes registered — every hasRoute check must return false
	s.False(hasRoute(app, "GET", "/"))
	s.False(hasRoute(app, "POST", "/"))
}

// --- NewApp with configs -----------------------------------------------------------------------

func (s *NewAppSuite) Test_NewApp_WithConfig_Returns_Non_Nil_App() {
	app := NewApp(s.T(), fiber.Config{AppName: "TestApp"})
	s.NotNil(app)
}

func (s *NewAppSuite) Test_NewApp_WithConfig_AppliesAppName() {
	app := NewApp(s.T(), fiber.Config{AppName: "ZenTestApp"})
	s.Equal("ZenTestApp", app.Config().AppName)
}

func (s *NewAppSuite) Test_NewApp_WithConfig_AcceptsRoutes() {
	app := NewApp(s.T(), fiber.Config{AppName: "RouteApp"})
	app.Get("/cfg", noop)
	s.True(hasRoute(app, "GET", "/cfg"))
}

func (s *NewAppSuite) Test_NewApp_WithConfig_IsFunctionalForRequests() {
	app := NewApp(s.T(), fiber.Config{AppName: "FunctionalApp"})
	app.Get("/hi", func(c fiber.Ctx) error {
		return c.SendString("configured")
	})

	zt := New(s.T())
	zt.Get(app, "/hi").OK().Contains("configured")
}

func (s *NewAppSuite) Test_NewApp_WithConfig_OnlyFirstConfigUsed() {
	// variadic — only the first config applies; subsequent are ignored
	app := NewApp(s.T(), fiber.Config{AppName: "First"}, fiber.Config{AppName: "Second"})
	s.Equal("First", app.Config().AppName)
}

func TestNewAppSuite(t *testing.T) {
	suite.Run(t, new(NewAppSuite))
}

// --- buildRouteMsg (internal helper) -----------------------------------------------------------

type BuildRouteMsgSuite struct {
	suite.Suite
}

func (s *BuildRouteMsgSuite) Test_NoExtraArgs_Returns_Method_And_Path() {
	msg := buildRouteMsg("GET", "/users")
	s.Equal("GET /users", msg)
}

func (s *BuildRouteMsgSuite) Test_MethodIsUppercased() {
	msg := buildRouteMsg("delete", "/items/:id")
	s.Equal("DELETE /items/:id", msg)
}

func (s *BuildRouteMsgSuite) Test_WithStringArg_AppendsSuffix() {
	msg := buildRouteMsg("POST", "/login", "login route missing")
	s.Equal("POST /login - login route missing", msg)
}

func (s *BuildRouteMsgSuite) Test_WithMultipleArgs_AppendsSuffix() {
	msg := buildRouteMsg("PUT", "/users/:id", "update route missing", "extra")
	s.Equal("PUT /users/:id - update route missing", msg)
}

func (s *BuildRouteMsgSuite) Test_WithMultipleArgs_EmptySuffix() {
	msg := buildRouteMsg("PUT", "/users/:id", "", "extra")
	s.Equal("PUT /users/:id - ", msg)
}

func (s *BuildRouteMsgSuite) Test_WithNonStringArg_FormatsWithSprint() {
	msg := buildRouteMsg("GET", "/count", 42)
	s.Equal("GET /count - 42", msg)
}

func TestBuildRouteMsgSuite(t *testing.T) {
	suite.Run(t, new(BuildRouteMsgSuite))
}

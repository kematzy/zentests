package zentests

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// HTTPRequestsSuite is the test suite struct for HTTP requests.
type HTTPRequestsSuite struct {
	suite.Suite
	app *fiber.App
	zt  *T
}

// SetupTest is called before each test.
func (s *HTTPRequestsSuite) SetupTest() {
	s.app = fiber.New()
	s.zt = New(s.T())

	s.app.Get("/test", func(c fiber.Ctx) error {
		return c.SendString("GET response")
	})

	// Saves an uploaded file's name and content into the given response string.
	s.app.Post("/upload", func(c fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		content, err := readFileHeader(file)
		if err != nil {
			return err
		}
		return c.SendString(c.FormValue("note") + "|" + file.Filename + "|" + string(content))
	})

	// Counts the uploaded files on the "file" field.
	s.app.Post("/upload-multi", func(c fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		return c.SendString(fmt.Sprintf("files=%d", len(form.File["file"])))
	})

	// Echoes the "x", "q" and "lang" query parameters.
	s.app.Get("/search", func(c fiber.Ctx) error {
		return c.SendString(c.Query("x") + "|" + c.Query("q") + "|" + c.Query("lang"))
	})

	// Echoes the request headers the default-header tests assert on.
	s.app.Get("/hdr", func(c fiber.Ctx) error {
		return c.SendString(c.Get("Accept") + "|" + c.Get("Authorization") + "|" + c.Get("X-Api-Key"))
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

// TearDownTest is called after each test.
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

// TestHTTPRequestsSuite runs the test suite.
func TestHTTPRequestsSuite(t *testing.T) {
	suite.Run(t, new(HTTPRequestsSuite))
}

// TestExecute_TransportError_StopsTest checks that a failed app.Test call is
// reported and stops the test before the nil response is dereferenced.
func TestExecute_TransportError_StopsTest(t *testing.T) {
	app := fiber.New()
	app.Get("/slow", func(c fiber.Ctx) error {
		time.Sleep(100 * time.Millisecond)
		return c.SendString("late")
	})
	cfg := fiber.TestConfig{Timeout: 5 * time.Millisecond, FailOnTimeout: true}

	rec := runRecorded(t, func(tb testing.TB) {
		execute(tb, app, "GET", "/slow", nil, nil, cfg)
	})

	assert.True(t, rec.fatal, "execute should stop the test on transport error")
	assert.NotEmpty(t, rec.errors, "transport error should be reported")
}

// TestPostForm_SpecialCharacters_ReachServerUnchanged checks that reserved
// form characters survive encoding instead of splitting or corrupting fields.
func TestPostForm_SpecialCharacters_ReachServerUnchanged(t *testing.T) {
	app := fiber.New()
	app.Post("/form", func(c fiber.Ctx) error {
		return c.SendString(c.FormValue("name"))
	})
	value := "a&b=c d+e"

	zt := New(t)
	zt.PostForm(app, "/form", map[string]string{"name": value}).Equals(value)
}

// TestTestConfig_ZeroValue_UsesDefaults checks that a zero-value config is replaced by
// the default timeout, so a stuck handler cannot block the test forever.
func TestTestConfig_ZeroValue_UsesDefaults(t *testing.T) {
	assert.Equal(t, defaultTestConfig, testConfig(fiber.TestConfig{}))
	assert.Positive(t, testConfig(fiber.TestConfig{}).Timeout)
	assert.True(t, testConfig(fiber.TestConfig{}).FailOnTimeout)
}

// TestTestConfig_Custom_IsKept checks that an explicit config is passed through unchanged.
func TestTestConfig_Custom_IsKept(t *testing.T) {
	custom := fiber.TestConfig{Timeout: 3 * time.Second, FailOnTimeout: false}
	assert.Equal(t, custom, testConfig(custom))
}

// TestSetHeader_SentWithEveryRequest checks that a default header reaches every request.
func (s *HTTPRequestsSuite) Test_SetHeader_SentWithEveryRequest() {
	zt := New(s.T()).SetHeader("Accept", "application/json")

	zt.Get(s.app, "/hdr").Equals("application/json||")
	zt.Get(s.app, "/hdr").Equals("application/json||")
}

// TestSetHeader_CanonicalizesName checks that two spellings of a name configure the same default.
func (s *HTTPRequestsSuite) Test_SetHeader_CanonicalizesName() {
	zt := New(s.T()).SetHeader("x-api-key", "first")

	zt.Get(s.app, "/hdr").Equals("||first")

	zt.SetHeader("X-Api-Key", "second")
	zt.Get(s.app, "/hdr").Equals("||second")
}

// TestWithBearer_SetsAuthorizationHeader checks that WithBearer sends the Authorization header.
func (s *HTTPRequestsSuite) Test_WithBearer_SetsAuthorizationHeader() {
	zt := New(s.T()).WithBearer("tok-1")

	zt.Get(s.app, "/hdr").Equals("|Bearer tok-1|")
}

// TestSetHeader_PerCallHeaderWins checks that a request's own headers take precedence
// over defaults with the same name.
func (s *HTTPRequestsSuite) Test_SetHeader_PerCallHeaderWins() {
	zt := New(s.T()).SetHeader("Content-Type", "text/plain")

	zt.PostJSON(s.app, "/json", map[string]any{"name": "John"}).OK().IsJSON()
}

// readFileHeader reads the full content of an uploaded file, for the upload test handlers.
func readFileHeader(header *multipart.FileHeader) ([]byte, error) {
	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close() //nolint:errcheck // read-only closer
	return io.ReadAll(file)
}

// Test_PostMultipart_UploadsFile checks that fields and file content reach the handler.
func (s *HTTPRequestsSuite) Test_PostMultipart_UploadsFile() {
	s.zt.PostMultipart(s.app, "/upload", map[string]string{"note": "hi"}, []FileUpload{
		{Name: "file", Filename: "a.txt", Content: []byte("hello")},
	}).OK().Equals("hi|a.txt|hello")
}

// Test_PostMultipart_FieldsOnly checks that multipart fields read as form values,
// without any files in the request.
func (s *HTTPRequestsSuite) Test_PostMultipart_FieldsOnly() {
	s.zt.PostMultipart(s.app, "/form", map[string]string{"name": "Jane"}, nil).
		OK().IsJSON().Has("name", "Jane")
}

// Test_PostMultipart_SameFieldTwice checks that a field name can carry two files.
func (s *HTTPRequestsSuite) Test_PostMultipart_SameFieldTwice() {
	s.zt.PostMultipart(s.app, "/upload-multi", nil, []FileUpload{
		{Name: "file", Filename: "a.txt", Content: []byte("a")},
		{Name: "file", Filename: "b.txt", Content: []byte("b")},
	}).OK().Equals("files=2")
}

// TestDefaultTimeout_UnsetEnvIsOneSecond checks that the default timeout stays at one second.
func TestDefaultTimeout_UnsetEnvIsOneSecond(t *testing.T) {
	t.Setenv("ZENTESTS_TIMEOUT", "")

	assert.Equal(t, time.Second, defaultTimeout())
}

// TestDefaultTimeout_SecondsEnvIsUsed checks that a ZENTESTS_TIMEOUT in seconds is used.
func TestDefaultTimeout_SecondsEnvIsUsed(t *testing.T) {
	t.Setenv("ZENTESTS_TIMEOUT", "3s")

	assert.Equal(t, 3*time.Second, defaultTimeout())
}

// TestDefaultTimeout_MillisecondsEnvIsUsed checks that sub-second ZENTESTS_TIMEOUT values are used.
func TestDefaultTimeout_MillisecondsEnvIsUsed(t *testing.T) {
	t.Setenv("ZENTESTS_TIMEOUT", "250ms")

	assert.Equal(t, 250*time.Millisecond, defaultTimeout())
}

// TestDefaultTimeout_InvalidEnvPanics checks that an invalid value panics instead of
// being silently ignored.
func TestDefaultTimeout_InvalidEnvPanics(t *testing.T) {
	t.Setenv("ZENTESTS_TIMEOUT", "not-a-duration")

	assert.Panics(t, func() { defaultTimeout() })
}

// Test_GetQuery_SendsEncodedQuery checks that query values with special characters
// reach the handler unchanged.
func (s *HTTPRequestsSuite) Test_GetQuery_SendsEncodedQuery() {
	query := url.Values{"q": {"go & rust = 1+1"}, "lang": {"en"}}

	s.zt.GetQuery(s.app, "/search", query).OK().Equals("|go & rust = 1+1|en")
}

// Test_GetQuery_AppendsToExistingQuery checks that values append after an existing
// query string in the path.
func (s *HTTPRequestsSuite) Test_GetQuery_AppendsToExistingQuery() {
	s.zt.GetQuery(s.app, "/search?x=1", url.Values{"q": {"z"}}).OK().Equals("1|z|")
}

// Test_GetQuery_EmptyValues checks that an empty query still produces a valid request.
func (s *HTTPRequestsSuite) Test_GetQuery_EmptyValues() {
	s.zt.GetQuery(s.app, "/search", nil).OK().Equals("||")
}

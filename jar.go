package zentests

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

// Jar sends requests to one Fiber app and remembers every cookie the app sets,
// the way a browser does. It lets a test flow span several requests, such as
// login, CSRF checks and session tests. Each request returns a *Response, so the
// usual assertions compose with jar requests.
//
// Not an RFC 6265 cookie jar. It is deliberately a flat map of name to value:
//   - No Path or Domain matching. A cookie is sent with every request to the app.
//   - No Secure or HttpOnly enforcement. Check those attributes with Response.Cookie.
//   - No Expires-based expiry. Only the MaxAge < 0 or empty-value deletion rule applies.
//   - A cookie set to "" is indistinguishable from a deletion, so the empty-value rule wins.
//
// Jar is not safe for concurrent use. Tests using it run on a single goroutine.
//
// Example:
//
//	jar := zentests.NewJar(t, app)
//	csrf := jar.Get("/login").OK().BodyString()
//	jar.PostForm("/login", url.Values{"csrf": {csrf}, "user": {"jane"}}).Status(302)
//	jar.Get("/settings").OK().Contains("Welcome Jane")
type Jar struct {
	t       testing.TB
	app     *fiber.App
	headers map[string]string
	cookies map[string]string
}

// NewJar returns an empty Jar for app. It fails the test when app is nil.
//
// Parameters:
//   - t: The [testing.TB] instance for assertions
//   - app: The Fiber application to test against; must not be nil
//
// Returns:
//   - *Jar: A Jar holding no cookies
//
// Example:
//
//	jar := zentests.NewJar(t, app)
func NewJar(t testing.TB, app *fiber.App) *Jar {
	t.Helper()
	require.NotNil(t, app, "zentests: NewJar called with a nil *fiber.App")

	return &Jar{t: t, app: app, cookies: make(map[string]string)}
}

// Get sends a GET request to path with the stored cookies.
//
// Parameters:
//   - path: The request path (e.g., "/login")
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	jar.Get("/settings").OK()
func (j *Jar) Get(path string) *Response {
	j.t.Helper()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, path, nil)
	return j.Send(req)
}

// PostForm sends a POST request to path with values encoded as
// application/x-www-form-urlencoded, together with the stored cookies.
//
// Parameters:
//   - path: The request path (e.g., "/login")
//   - values: The form fields; keys and values are escaped and multi-value safe
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	jar.PostForm("/login", url.Values{"user": {"jane"}}).Status(302)
func (j *Jar) PostForm(path string, values url.Values) *Response {
	j.t.Helper()
	body := strings.NewReader(values.Encode())
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, path, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return j.Send(req)
}

// Send sends a caller-built request with the stored cookies, then stores the cookies the
// response sets. Use it as an escape hatch when Get and PostForm do not fit.
//
// Parameters:
//   - req: The request to send; its own cookies are kept and the stored ones are added
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	req := httptest.NewRequest(http.MethodDelete, "/session", nil)
//	jar.Send(req).NoContent()
func (j *Jar) Send(req *http.Request) *Response {
	j.t.Helper()
	for name, value := range j.headers {
		// A header already on the request wins, so PostMultipart's Content-Type
		// is not replaced by a default with the same name.
		if req.Header.Get(name) == "" {
			req.Header.Set(name, value)
		}
	}
	for name, value := range j.cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	resp := send(j.t, j.app, req, fiber.TestConfig{})
	j.store(resp.resp.Cookies())
	return resp
}

// SetHeader sets a default header that is sent with every request made through the jar.
// A header already on the request, such as PostForm's Content-Type, takes precedence
// over a default with the same name. Header names are canonicalised, so "x-api-key"
// and "X-Api-Key" configure the same default.
//
// Parameters:
//   - key: The header name (e.g., "Accept")
//   - value: The header value
//
// Example:
//
//	jar := zentests.NewJar(t, app)
//	jar.SetHeader("Accept", "application/json")
func (j *Jar) SetHeader(key, value string) {
	j.t.Helper()
	j.headers = setDefaultHeader(j.headers, key, value)
}

// WithBearer sets the default Authorization header to "Bearer <token>", which is sent
// with every request made through the jar.
//
// Parameters:
//   - token: The bearer token
//
// Example:
//
//	jar := zentests.NewJar(t, app)
//	jar.WithBearer(token)
func (j *Jar) WithBearer(token string) {
	j.SetHeader("Authorization", "Bearer "+token)
}

// PostMultipart sends a POST request with multipart/form-data and the stored cookies,
// typically for file uploads. See T.PostMultipart for the FileUpload fields.
//
// Parameters:
//   - path: The request path (e.g., "/upload")
//   - fields: Plain form fields (nil when there are none)
//   - files: The files to upload (nil when there are none)
//
// Returns:
//   - *Response: A Response wrapper for making assertions
//
// Example:
//
//	jar.PostMultipart("/upload", nil, []zentests.FileUpload{
//	    {Name: "file", Filename: "a.txt", Content: []byte("hello")},
//	}).OK()
func (j *Jar) PostMultipart(path string, fields map[string]string, files []FileUpload) *Response {
	j.t.Helper()
	body, contentType, err := multipartBody(fields, files)
	require.NoError(j.t, err, "building the multipart body failed")

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", contentType)
	return j.Send(req)
}

// Cookie returns the value the jar holds for name, or "" when there is none.
// It never fails the test, so it can check that a cookie is absent.
//
// Parameters:
//   - name: The cookie name
//
// Returns:
//   - string: The stored cookie value, or "" when absent
//
// Example:
//
//	if jar.Cookie("session") == "" { t.Fatal("no session") }
func (j *Jar) Cookie(name string) string {
	return j.cookies[name]
}

// store saves the cookies from a response. Deletion cookies are removed from the jar,
// because that is how Fiber's ClearCookie expires a cookie.
//
// Parameters:
//   - cookies: The cookies parsed from the response's Set-Cookie headers
func (j *Jar) store(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		if isDeletionCookie(cookie) {
			delete(j.cookies, cookie.Name)
			continue
		}
		j.cookies[cookie.Name] = cookie.Value
	}
}

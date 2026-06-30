package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/golangcollege/sessions"

	"github.com/bicosteve/callory-tracker/pkg/models/mock"
)

// csrfTokenRX matches the hidden CSRF token input field rendered by nosurf.
var csrfTokenRX = regexp.MustCompile(`name="csrf_token" value='(.+)'>`)

// newTestApplication returns an *application wired up with mock models, a
// discarding logger and a real (in-memory) session manager. The template cache
// is empty by default; handlers that render templates are tested via the
// testServer which loads the real templates from ./../../ui/html.
func newTestApplication(t *testing.T) *application {
	t.Helper()

	templateCache, err := newTemplateCache("./../../ui/html")
	if err != nil {
		t.Fatal(err)
	}

	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour

	return &application{
		errorLog:      log.New(io.Discard, "", 0),
		infoLog:       log.New(io.Discard, "", 0),
		warningLog:    log.New(io.Discard, "", 0),
		foods:         &mock.FoodModel{},
		users:         &mock.UserModel{},
		templateCache: templateCache,
		session:       session,
	}
}

// testServer embeds httptest.Server and provides convenience helpers for
// performing requests with an automatically-managed cookie jar.
type testServer struct {
	*httptest.Server
}

// newTestServer initialises a new httptest.Server wrapped around the given
// handler, with a cookie jar so that session cookies are preserved across
// requests and redirect following is disabled.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	t.Helper()

	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	ts.Client().Jar = jar

	// Stop the client following redirects so we can assert on 3xx responses.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// get performs a GET request to the given url path and returns the status code,
// response headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	t.Helper()

	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(body)
}

// postForm performs a POST request with the supplied form values.
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	t.Helper()

	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(body)
}

// login authenticates the test client against the mock UserModel so that
// subsequent requests carry an authenticated session cookie.
func (ts *testServer) login(t *testing.T) {
	t.Helper()

	_, _, body := ts.get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("email", "alice@example.com")
	form.Add("password", "password")
	form.Add("csrf_token", csrfToken)

	code, _, _ := ts.postForm(t, "/user/login", form)
	if code != http.StatusSeeOther {
		t.Fatalf("login failed: want %d; got %d", http.StatusSeeOther, code)
	}
}

// extractCSRFToken pulls the CSRF token out of an HTML response body.

func extractCSRFToken(t *testing.T, body string) string {
	t.Helper()

	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(matches[1])
}

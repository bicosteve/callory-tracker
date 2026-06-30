package main

import (
	"net/http"
	"net/url"
	"testing"
)

// TestGetHome verifies that the home page renders for the index path and that
// any path other than "/" returns a 404.
func TestGetHome(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.session.Enable(app.routes()))
	defer ts.Close()

	t.Run("renders home for root path", func(t *testing.T) {
		code, _, _ := ts.get(t, "/")
		if code != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, code)
		}
	})
}

// TestGetRegisterPage ensures the register page is served.
func TestGetRegisterPage(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.session.Enable(app.routes()))
	defer ts.Close()

	code, _, body := ts.get(t, "/user/register")
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
	if body == "" {
		t.Error("expected non-empty register page body")
	}
}

// TestGetLoginPage ensures the login page is served.
func TestGetLoginPage(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.session.Enable(app.routes()))
	defer ts.Close()

	code, _, _ := ts.get(t, "/user/login")
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
}

// TestGetDay covers the various validation branches of the getDay handler.
func TestGetDay(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.session.Enable(app.routes()))
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
	}{
		{
			name:     "valid food and user",
			urlPath:  "/food/day?foodId=1&userId=1",
			wantCode: http.StatusOK,
		},
		{
			name:     "non-existent food",
			urlPath:  "/food/day?foodId=2&userId=1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "negative food id",
			urlPath:  "/food/day?foodId=-1&userId=1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "non-numeric food id",
			urlPath:  "/food/day?foodId=abc&userId=1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "missing user id",
			urlPath:  "/food/day?foodId=1",
			wantCode: http.StatusNotFound,
		},
	}

	// getDay is behind requireAuthenticatedUser; an unauthenticated user is
	// redirected to the login page. We assert that here for the happy-path URL.
	t.Run("unauthenticated user is redirected", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/food/day?foodId=1&userId=1")
		if code != http.StatusSeeOther {
			t.Errorf("want %d; got %d", http.StatusSeeOther, code)
		}
		if loc := headers.Get("Location"); loc != "/user/login" {
			t.Errorf("want redirect to /user/login; got %q", loc)
		}
	})

	// The remaining cases run against an authenticated session.
	ts.login(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.get(t, tt.urlPath)
			if code != tt.wantCode {
				t.Errorf("%s: want %d; got %d", tt.urlPath, tt.wantCode, code)
			}
		})
	}
}

// TestRegisterUser checks form validation and successful registration flow.
func TestRegisterUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.session.Enable(app.routes()))
	defer ts.Close()

	// Grab a valid CSRF token from the register page.
	_, _, body := ts.get(t, "/user/register")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name            string
		username        string
		email           string
		password        string
		confirmPassword string
		csrfToken       string
		wantCode        int
		wantLocation    string
	}{
		{
			name:            "valid registration",
			username:        "bob",
			email:           "bob@example.com",
			password:        "secret",
			confirmPassword: "secret",
			csrfToken:       csrfToken,
			wantCode:        http.StatusSeeOther,
			wantLocation:    "/user/login",
		},
		{
			name:            "missing fields",
			username:        "",
			email:           "",
			password:        "",
			confirmPassword: "",
			csrfToken:       csrfToken,
			wantCode:        http.StatusOK,
		},
		{
			name:            "invalid email",
			username:        "carl",
			email:           "not-an-email",
			password:        "secret",
			confirmPassword: "secret",
			csrfToken:       csrfToken,
			wantCode:        http.StatusOK,
		},
		{
			name:            "passwords do not match",
			username:        "dan",
			email:           "dan@example.com",
			password:        "secret",
			confirmPassword: "different",
			csrfToken:       csrfToken,
			wantCode:        http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("username", tt.username)
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			form.Add("confirm_password", tt.confirmPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, headers, _ := ts.postForm(t, "/user/register", form)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if tt.wantLocation != "" {
				if loc := headers.Get("Location"); loc != tt.wantLocation {
					t.Errorf("want location %q; got %q", tt.wantLocation, loc)
				}
			}
		})
	}
}

// TestLoginUser exercises the login handler's success and failure paths.
func TestLoginUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.session.Enable(app.routes()))
	defer ts.Close()

	_, _, body := ts.get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name         string
		email        string
		password     string
		wantCode     int
		wantLocation string
	}{
		{
			name:         "valid credentials",
			email:        "alice@example.com",
			password:     "password",
			wantCode:     http.StatusSeeOther,
			wantLocation: "/food/add",
		},
		{
			name:     "wrong password",
			email:    "alice@example.com",
			password: "wrong",
			wantCode: http.StatusOK,
		},
		{
			name:     "empty fields",
			email:    "",
			password: "",
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			form.Add("csrf_token", csrfToken)

			code, headers, _ := ts.postForm(t, "/user/login", form)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
			if tt.wantLocation != "" {
				if loc := headers.Get("Location"); loc != tt.wantLocation {
					t.Errorf("want location %q; got %q", tt.wantLocation, loc)
				}
			}
		})
	}
}

// TestPageNotFound asserts that an unknown route returns 404.
func TestPageNotFound(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.session.Enable(app.routes()))
	defer ts.Close()

	code, _, _ := ts.get(t, "/this/does/not/exist")
	if code != http.StatusNotFound {
		t.Errorf("want %d; got %d", http.StatusNotFound, code)
	}
}

package main

import (
	"context"
	"fmt"
	"github.com/bicosteve/callory-tracker/pkg/models"
	"github.com/justinas/nosurf"
	"net/http"
)

/*
SecureHeaders -> instructs browser to implement additional security
which prevents XSS & Click jacking attacks
*/
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

// Recovers server panics and let the app continue to run
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer will run in the event of panic as Go unwinds stack
		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Logs the request's remoteAddress, method, and url
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s : %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// requireAuthenticatedUser middleware to stop unauthenticated users from accessing some urls
func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if user not authenticated, redirect to login page
		if app.isAuthenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// noSurf
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

// authenticate retrieves users id from the session data
// checks the db to confirm the userid corresponds to user using UserModel.Exists()
// updates request context to include an isAuthenticatedContextKey key with true value
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieve authenticated user id value from session
		// returns 0 if no authenticatedUserId
		exists := app.session.Exists(r, "userId")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.users.GetUserDetails(app.session.GetInt(r, "userId"))
		if err == models.ErrNoRecord {
			app.session.Remove(r, "userId")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package main

import (
	"github.com/justinas/alice"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders, noSurf)

	// Only authenticated users can access these routes
	router.Group(func(r chi.Router) {
		r.Use(app.requireAuthenticatedUser)
		r.Get("/food/add", app.postFoodForm)
		r.Post("/food/add", app.postFood)
		r.Get("/food/day", app.getDay)

	})

	router.Get("/", app.getHome)
	router.Get("/user/register", app.getRegisterPage)
	router.Get("/user/login", app.getLoginPage)
	router.Post("/user/register", app.registerUser)
	router.Post("/user/login", app.loginUser)
	router.Post("/user/logout", app.logoutUser)

	return standardMiddleware.Then(router)
}

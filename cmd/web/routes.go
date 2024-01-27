package main

import (
	"github.com/justinas/alice"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Routes
	router.Get("/", app.getHome)
	router.Get("/food/add", app.postFoodForm)
	router.Post("/food/add", app.postFood)

	router.Get("/food/day", app.getDay)

	// User register
	router.Get("/user/register", app.getRegisterPage)
	router.Get("/user/login", app.getLoginPage)

	router.Post("/user/register", app.registerUser)
	router.Post("/user/login", app.loginUser)
	router.Post("/user/logout", app.logoutUser)

	return standardMiddleware.Then(router)
}

package main

import (
	"net/http"

	"github.com/bicosteve/callory-tracker/ui"
	"github.com/justinas/alice"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	// take embedded filesystem and convert it to http.FS type
	var fileServer = http.FS(ui.Files)
	fs := http.FileServer(fileServer)

	// serving static files/css with embed
	router.Handle("/css/*", fs)

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders, noSurf, app.authenticate)

	// Only authenticated users can access these routes
	router.Group(func(r chi.Router) {
		r.Use(app.requireAuthenticatedUser)
		r.Get("/food/add", app.postFoodForm)
		r.Post("/food/add", app.postFood)
		r.Get("/food/day", app.getDay)
		r.Get("/food/edit", app.editFoodForm)
		r.Post("/food/update", app.editFood)

	})

	router.Get("/", app.getHome)
	router.Get("/user/register", app.getRegisterPage)
	router.Get("/user/login", app.getLoginPage)
	router.Post("/user/register", app.registerUser)
	router.Post("/user/login", app.loginUser)
	router.Post("/user/logout", app.logoutUser)

	return standardMiddleware.Then(router)
}

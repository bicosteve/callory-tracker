package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	// Routes
	router.Get("/", app.getHome)
	router.Get("/food/add", app.getFoodPage)
	router.Get("/food/day", app.getDay)
	router.Post("/food/add", app.postFood)

	router.Post("/user/register", app.registerUser)

	return router
}

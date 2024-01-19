package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	// Routes
	router.Get("/", getHome)
	router.Get("/food/add", getFoodPage)
	router.Get("/food/day", getDay)
	router.Post("/food/post", postFood)

	return router
}

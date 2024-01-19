package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes

func routes() http.Handler {
	router := chi.NewRouter()

	// Routes
	router.Get("/", getHome)
	router.Get("/food/add", getAddFoodPage)
	router.Post("/food/add", postFood)
	router.Get("/food/day", getDayConsumption)

	return router
}

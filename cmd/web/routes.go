package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes

func routes() http.Handler {
	router := chi.NewRouter()

	// Routes
	router.Get("/home", getHome)

	return router
}

package main

import (
	"fmt"
	"net/http"
)

// for handlers

func getHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		http.NotFound(w, r)
		return
	}

	fmt.Println("Welcome, home")
}

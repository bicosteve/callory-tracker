package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// for handlers

func getHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		http.NotFound(w, r)
		return
	}

	templateSet, err := template.ParseFiles("./ui/html/home.page.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal error loading home template", 500)
		return
	}

	err = templateSet.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error executing template set", 500)
		return
	}

	fmt.Println("Welcome, home")

}

func getAddFoodPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is get food")
}

func postFood(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	fmt.Println("This is get food")
}

func getDayConsumption(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get day consumption")
}

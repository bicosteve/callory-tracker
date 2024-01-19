package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/bicosteve/callory-tracker/pkg/helpers"
)

var files []string

func getHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	nav, err := helpers.LoadTemplate("./ui/html/partial.nav.html")
	if err != nil {
		log.Fatal(err)
	}

	base, err := helpers.LoadTemplate("./ui/html/layout.base.html")
	if err != nil {
		log.Fatal(err)

	}

	home, err := helpers.LoadTemplate("./ui/html/home.page.html")
	if err != nil {
		log.Fatal(err)

	}

	files = append(files, nav)
	files = append(files, base)
	files = append(files, home)

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal error loading home template", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error executing template set", 500)
		return
	}

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

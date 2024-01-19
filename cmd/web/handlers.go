package main

import (
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
		http.Error(w, "internal error loading home template", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error executing template set", 500)
		return
	}

}

func getFoodPage(w http.ResponseWriter, r *http.Request) {
	nav, _ := helpers.LoadTemplate("./ui/html/partial.nav.html")
	base, _ := helpers.LoadTemplate("./ui/html/layout.base.html")
	add, _ := helpers.LoadTemplate("./ui/html/add_food.page.html")

	files = append(files, nav)
	files = append(files, base)
	files = append(files, add)

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal error parsing templates", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error executing templates", 500)
		return

	}

}

func postFood(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

}

func getDay(w http.ResponseWriter, r *http.Request) {
	nav, _ := helpers.LoadTemplate("./ui/html/partial.nav.html")
	base, _ := helpers.LoadTemplate("./ui/html/layout.base.html")
	day, _ := helpers.LoadTemplate("./ui/html/day.page.html")

	files = append(files, nav)
	files = append(files, base)
	files = append(files, day)

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal error parsing templates", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal error executing templates", 500)
		return
	}

}

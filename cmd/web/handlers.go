package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/bicosteve/callory-tracker/pkg/helpers"
)

var files []string

func (app *application) getHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	nav, err := helpers.LoadTemplate("./ui/html/nav.partial.html")
	if err != nil {
		app.errorLog.Printf(err.Error())
	}

	base, err := helpers.LoadTemplate("./ui/html/layout.base.html")
	if err != nil {
		app.errorLog.Printf(err.Error())

	}

	home, err := helpers.LoadTemplate("./ui/html/home.page.html")
	if err != nil {
		app.errorLog.Printf(err.Error())

	}

	files = append(files, nav)
	files = append(files, base)
	files = append(files, home)

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Printf(err.Error())
		http.Error(w, "internal error loading home template", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error executing template set", http.StatusInternalServerError)
		return
	}

}

func (app *application) getFoodPage(w http.ResponseWriter, r *http.Request) {
	nav, _ := helpers.LoadTemplate("./ui/html/nav.partial.html")
	base, _ := helpers.LoadTemplate("./ui/html/layout.base.html")
	add, _ := helpers.LoadTemplate("./ui/html/add_food.page.html")

	files = append(files, nav)
	files = append(files, base)
	files = append(files, add)

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal error parsing templates", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error executing templates", http.StatusInternalServerError)
		return

	}

}

func (app *application) postFood(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

}

func (app *application) getDay(w http.ResponseWriter, r *http.Request) {
	nav, _ := helpers.LoadTemplate("./ui/html/nav.partial.html")
	base, _ := helpers.LoadTemplate("./ui/html/layout.base.html")
	day, _ := helpers.LoadTemplate("./ui/html/day.page.html")

	files = append(files, nav)
	files = append(files, base)
	files = append(files, day)

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal error parsing templates", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal error executing templates", http.StatusInternalServerError)
		return
	}

}

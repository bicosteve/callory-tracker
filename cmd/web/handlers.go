package main

import (
	"html/template"
	"net/http"

	"github.com/bicosteve/callory-tracker/pkg/helpers"
)

var files []string

const cal = 4

func (app *application) getHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	nav, err := helpers.LoadTemplate("./ui/html/nav.partial.html")
	if err != nil {
		// return error if partial not found
		app.errorLog.Printf(err.Error())
		return
	}

	base, err := helpers.LoadTemplate("./ui/html/layout.base.html")
	if err != nil {
		app.serverError(w, err)
		return

	}

	home, err := helpers.LoadTemplate("./ui/html/home.page.html")
	if err != nil {
		app.serverError(w, err)
		return

	}

	files = append(files, nav)
	files = append(files, base)
	files = append(files, home)

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err) // serverError() helper
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
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		return

	}

}

func (app *application) postFood(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	meal := "lunch"
	name := "pizza"
	protein := 5
	carbohydrates := 20
	fat := 10
	calories := (protein * cal) + (carbohydrates * cal) + (fat * cal)
	userId := 1

	id, err := app.foods.InsertFood(meal, name, protein, carbohydrates, fat, calories, userId)

	if err != nil {
		app.serverError(w, err)
		app.errorLog.Println(err)
		return
	}

	w.Write([]byte("Post foods"))
	app.infoLog.Println(id)

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
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		return
	}

}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	nav, _ := helpers.LoadTemplate("./ui/html/nav.partial.html")
	base, _ := helpers.LoadTemplate("./ui/html/layout.base.html")
	register, _ := helpers.LoadTemplate("./ui/html/register.page.html")

	files = append(files, nav)
	files = append(files, base)
	files = append(files, register)

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	username := "Bix"
	email := "bix@bix.com"
	password := "12345"

	err = app.users.RegisterUser(username, email, password)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		return
	}

}

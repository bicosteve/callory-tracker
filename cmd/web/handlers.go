package main

import (
	"errors"
	"fmt"
	"github.com/bicosteve/callory-tracker/pkg/forms"
	"github.com/bicosteve/callory-tracker/pkg/models"
	"net/http"
	"strconv"
)

const cal = 4

func (app *application) getHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	userID := 1

	foods, err := app.foods.GetFoods(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.renderATemplate(w, r, "home.page.html", &templateData{Foods: foods})

}

// postFoodForm -> renders the add food form
func (app *application) postFoodForm(w http.ResponseWriter, r *http.Request) {
	// passing new empty forms.Form object to the template
	app.renderATemplate(w, r, "add_food.page.html",
		&templateData{Form: forms.NewForm(nil)})
}

func (app *application) postFood(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// using r.ParseForm() method to parse the form
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	/*
		NB: r.PostForm() works only for POST, PATCH, PUT.
		Can also use r.Form() for all http requests. Used for query strings
		/food/add?foo=bar
		r.Form.Get("foo")
	*/

	form := forms.NewForm(r.PostForm)
	form.Required("meal", "name", "protein", "carbohydrate", "fat")
	form.MaxLength("name", 20)
	form.MinValue(1, "protein", "carbohydrate", "fat")

	if !form.Valid() {
		app.renderATemplate(w, r, "add_food.page.html",
			&templateData{Form: form})
		return
	}

	meal := form.Get("meal")
	name := form.Get("name")
	protein, _ := strconv.Atoi(form.Get("protein"))
	carbs, _ := strconv.Atoi(form.Get("carbohydrate"))
	fat, _ := strconv.Atoi(form.Get("fat"))

	calories := (protein * cal) + (carbs * cal) + (fat * cal)
	userId := 1

	id, err := app.foods.InsertFood(meal, name, protein, carbs, fat, calories, userId)
	if err != nil {
		app.serverError(w, err)
		app.errorLog.Println(err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/food/day?foodId=%d&userId=%d", id, userId), http.StatusSeeOther)
}

func (app *application) getDay(w http.ResponseWriter, r *http.Request) {
	foodID, err := strconv.Atoi(r.URL.Query().Get("foodId"))
	userID, err := strconv.Atoi(r.URL.Query().Get("userId"))

	if err != nil || foodID < 1 {
		app.notFound(w)
		return
	}

	if err != nil || userID < 1 {
		app.notFound(w)
		return
	}

	food, err := app.foods.GetFood(foodID, userID)
	if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	}

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.renderATemplate(w, r, "day.page.html", &templateData{Food: food})

}

func (app *application) getRegisterPage(w http.ResponseWriter, r *http.Request) {
	app.renderATemplate(w, r, "register.page.html", nil)
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	username := "Bix"
	email := "bix@bix.com"
	password := "12345"

	err := app.users.RegisterUser(username, email, password)

	_ = err

	app.renderATemplate(w, r, "register.page.html", nil)

}

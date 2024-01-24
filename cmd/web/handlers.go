package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bicosteve/callory-tracker/pkg/models"
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

func (app *application) getFoodPage(w http.ResponseWriter, r *http.Request) {
	app.renderATemplate(w, r, "add_food.page.html", nil)
}

func (app *application) postFood(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	meal := "supper"
	name := "Fish/Ugali"
	protein := 15
	carbohydrates := 10
	fat := 25
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

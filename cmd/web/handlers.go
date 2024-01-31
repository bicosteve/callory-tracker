package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bicosteve/callory-tracker/pkg/forms"
	"github.com/bicosteve/callory-tracker/pkg/models"
)

const cal = 4

func (app *application) getHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	userId := app.session.GetInt(r, "userId")

	foods, err := app.foods.GetFoods(userId)
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
	if r.Method != http.MethodPost {
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
	userId := app.session.GetInt(r, "userId")

	id, err := app.foods.InsertFood(meal, name, protein, carbs, fat, calories, userId)
	if err != nil {
		app.serverError(w, err)
		app.errorLog.Println(err)
		return
	}

	// app.session.Put: adding flash message when meal is created
	app.session.Put(r, "flash", "Meal successfully added!")

	http.Redirect(w, r, fmt.Sprintf("/food/day?foodId=%d&userId=%d", id, userId), http.StatusSeeOther)
}

func (app *application) editFoodForm(w http.ResponseWriter, r *http.Request) {
	foodId, err := strconv.Atoi(r.URL.Query().Get("foodId"))
	if err != nil || foodId < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil || userID < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	food, err := app.foods.GetFood(foodId, userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.renderATemplate(w, r, "edit_food.page.html", &templateData{
		Form: forms.NewForm(nil), Food: food,
	})
}

func (app *application) editFood(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.NewForm(r.PostForm)
	form.Required("meal", "name", "protein", "carbohydrate", "fat")
	form.MaxLength("name", 20)
	form.MinValue(1, "protein", "carbohydrate", "fat")

	if !form.Valid() {
		app.renderATemplate(w, r, "edit_food.page.html",
			&templateData{Form: form})
		return
	}

	foodId, _ := strconv.Atoi(form.Get("foodId"))
	userId := app.session.GetInt(r, "userId")
	meal := form.Get("meal")
	name := form.Get("name")
	protein, _ := strconv.Atoi(form.Get("protein"))
	cabs, _ := strconv.Atoi(form.Get("carbohydrate"))
	fat, _ := strconv.Atoi(form.Get("fat"))
	calories := (protein * cal) + (cabs * cal) + (fat * cal)

	id, err := app.foods.UpdateFood(meal, name, protein, cabs, fat, calories, foodId, userId)
	_ = id

	if err != nil {
		app.serverError(w, err)
		app.errorLog.Println(err)
		return
	}

	app.session.Put(r, "flash", fmt.Sprintf("%s updated successfully", name))

	http.Redirect(w, r, fmt.Sprintf("/food/day?foodId=%d&userId=%d", foodId, userId), http.StatusSeeOther)
}

func (app *application)deleteFood(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		w.Header().Set("Allowed", "POST")
		app.clientError(w,http.StatusMethodNotAllowed)
		return 
	}

	foodId,err := strconv.Atoi(r.FormValue("foodId"))
	if err != nil {
		app.clientError(w,http.StatusBadRequest)
		return 
	}

	userId := app.session.GetInt(r,"userId")

	id, err := app.foods.DeleteFood(foodId,userId)
	if err != nil {
		app.serverError(w,err)
		return 
	}

	app.session.Put(r,"flash",fmt.Sprintf("%d item deleted",id))

	http.Redirect(w,r,"/",http.StatusSeeOther)
}

func (app *application) getDay(w http.ResponseWriter, r *http.Request) {
	foodID, err := strconv.Atoi(r.URL.Query().Get("foodId"))
	if err != nil || foodID < 1 {
		app.notFound(w)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("userId"))
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

	app.renderATemplate(w, r, "day.page.html", &templateData{
		Food: food,
	})
}

func (app *application) getRegisterPage(w http.ResponseWriter, r *http.Request) {
	app.renderATemplate(w, r, "register.page.html",
		&templateData{Form: forms.NewForm(nil)})
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.NewForm(r.PostForm)
	form.Required("username", "email", "password", "confirm_password")
	form.MaxLength("username", 10)
	//form.MinLength("password", 5)
	form.ValidateEmail("email", forms.EmailRegex)
	form.ComparePasswords("password", "confirm_password")

	if !form.Valid() {
		app.renderATemplate(w, r, "register.page.html", &templateData{Form: form})
		return
	}

	username := form.Get("username")
	email := form.Get("email")
	password := form.Get("password")

	err = app.users.RegisterUser(username, email, password)
	if err != nil {
		app.serverError(w, err)
		app.errorLog.Println(err)
		return
	}

	app.session.Put(r, "flash", "Registered successfully")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) getLoginPage(w http.ResponseWriter, r *http.Request) {
	app.renderATemplate(w, r, "login.page.html",
		&templateData{Form: forms.NewForm(nil)})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	form := forms.NewForm(r.PostForm)
	form.Required("email", "password")

	if !form.Valid() {
		app.renderATemplate(w, r, "login.page.html", &templateData{Form: form})
		return
	}

	email := form.Get("email")
	password := form.Get("password")

	userId, err := app.users.LoginUser(email, password)
	if err == models.ErrorInvalidCredentials {
		form.Errors.Add("generic", "Email or password is incorrect")
		app.renderATemplate(w, r, "login.page.html", &templateData{Form: form})
		return
	}

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Successfully logged in")
	app.session.Put(r, "userId", userId)

	http.Redirect(w, r, "/food/add", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Use the RenewToken() method on the current session to change session ID.
	// This is good practise

	// Logout means removing the userId from the session
	app.session.Remove(r, "userId")
	app.session.Put(r, "flash", "You have successfully logout")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

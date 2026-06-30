package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrDuplicateEmail = errors.New("models: user already exists")
var ErrorInvalidCredentials = errors.New("models: incorrect password or email")

type FoodModelInterface interface {
	InsertFood(meal string, name string, protein int, carbohydrate int, fat int, calories int, userId int) (int, error)
	GetFood(foodId, userId int) (*Food, error)
	GetFoodTotal(userId int, createdAt string) (*Food, error)
	GetFoods(userId int) ([]*Food, error)
	UpdateFood(meal string, name string, protein, cabs, fat, calory, foodId, userId int) (int, error)
	DeleteFood(foodId, userId int) (int, error)
}

type UserModelInterface interface {
	RegisterUser(username, email, password string) error
	LoginUser(email, password string) (int, error)
	GetUserDetails(id int) (*User, error)
}

type Food struct {
	ID            int       `json:"id"`
	Meal          string    `json:"meal"`
	Name          string    `json:"name"`
	Protein       int       `json:"protein"`
	Carbohydrates int       `json:"carbohydrates"`
	Fat           int       `json:"fat"`
	Calories      int       `json:"calories"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserID        int       `json:"userId"`
}

type User struct {
	ID              int       `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirm_password"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

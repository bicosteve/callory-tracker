package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrDuplicateEmail = errors.New("models: user already exists")
var ErrorInvalidCredentials = errors.New("models: incorrect password or email")

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

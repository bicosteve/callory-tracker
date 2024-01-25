package configs

import (
	"github.com/bicosteve/callory-tracker/pkg/models/mysql"
	"html/template"
	"log"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Foods         *mysql.FoodModel
	Users         *mysql.UserModel
	templateCache map[string]*template.Template
}

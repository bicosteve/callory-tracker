package configs

import (
	"html/template"
	"log"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	TemplateCache map[string]*template.Template
}

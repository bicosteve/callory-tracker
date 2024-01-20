package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/bicosteve/callory-tracker/pkg/helpers"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
}

func main() {
	/*
		Creating logger for logging info messages
		Creating logger for logging error messages
	*/

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Loading env file
	env, err := helpers.LoadEnv("./env/.env")
	if err != nil {
		errorLog.Fatal(err.Error())
	}

	err = godotenv.Load(env)
	if err != nil {
		errorLog.Fatal(err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":4000"
	}

	htmlDir, err := helpers.LoadTemplates("./ui/html")
	if err != nil {
		errorLog.Fatal(err.Error())
		return
	}

	templatesCache, err := loadTemplateCache(htmlDir)
	if err != nil {
		errorLog.Fatal(err.Error())
		return
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templatesCache,
	}

	serve := &http.Server{
		Addr:     port,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	infoLog.Printf("Server running at %s ... \n", port)
	err = serve.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}

}

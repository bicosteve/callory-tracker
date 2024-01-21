package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/bicosteve/callory-tracker/pkg/db"
	"github.com/bicosteve/callory-tracker/pkg/helpers"
	"github.com/bicosteve/callory-tracker/pkg/models/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
	foods         *mysql.FoodModel
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

	dbUser := os.Getenv("DBUSER")
	dbName := os.Getenv("DBNAME")
	dbHost := os.Getenv("DBHOST")
	dbPassword := os.Getenv("DBPASSWORD")
	dbPort := os.Getenv("DBPORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := db.OpenDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()
	// closes the db connection pool before main func exits

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
		foods:         &mysql.FoodModel{DB: db},
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

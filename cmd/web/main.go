package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"

	"github.com/bicosteve/callory-tracker/pkg/db"
	"github.com/bicosteve/callory-tracker/pkg/helpers"
	"github.com/bicosteve/callory-tracker/pkg/models/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	foods         *mysql.FoodModel
	users         *mysql.UserModel
	templateCache map[string]*template.Template
	session       *sessions.Session
}

type contextKey string

const contextKeyUser = contextKey("user")

func main() {
	// infoLog: logging info messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// errorLog: logging info messages
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	var dbUser string
	var dbName string
	var dbHost string
	var dbPassword string
	var dbPort string
	var secret string
	var dsn string
	var port string

	// Which which environment the application is running and set the configs properly

	if os.Getenv("ENV") == "PROD" {
		dbUser = os.Getenv("DBUSER")
		dbPassword = os.Getenv("DBPASSWORD")
		dbHost = os.Getenv("DBHOST")
		dbName = os.Getenv("DBNAME")
		dbPort = os.Getenv("DBPORT")
		secret = os.Getenv("SESSION")
		port = os.Getenv("PORT")
	} else {
		// Loading env file
		env, err := helpers.LoadEnv(".env")
		if err != nil {
			errorLog.Fatal(err.Error())
		}

		err = godotenv.Load(env)
		if err != nil {
			errorLog.Fatal(err.Error())
		}

		port = os.Getenv("PORT")
		if port == "" {
			port = ":4000"
		}

		dbUser = os.Getenv("DBUSER")
		dbName = os.Getenv("DBNAME")
		dbHost = os.Getenv("DBHOST")
		dbPassword = os.Getenv("DBPASSWORD")
		dbPort = os.Getenv("DBPORT")
		secret = os.Getenv("SESSION")

	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	fmt.Println("Database connection uri", dsn)
	fmt.Println("Port", port)

	db, err := db.OpenDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()
	// closes the db connection pool before main func exits

	templateCache, err := newTemplateCache("./ui/html")
	/*
		NB:
		Having templateCache in the application struct means;
		1. We have an in memory cache of the relevant template.
		2. Handlers have access to this cache via the application struct
	*/
	if err != nil {
		errorLog.Fatal(err.Error())
		return
	}

	session := sessions.New([]byte(secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		foods:         &mysql.FoodModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		templateCache: templateCache,
		session:       session,
	}

	serve := &http.Server{
		Addr:         port,
		Handler:      session.Enable(app.routes()), // wraps handlers with session
		ErrorLog:     errorLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Server running at %s ... \n", port)
	err = serve.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}

}

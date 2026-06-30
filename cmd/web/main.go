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
	"github.com/bicosteve/callory-tracker/pkg/logger"
	"github.com/bicosteve/callory-tracker/pkg/models"
	"github.com/bicosteve/callory-tracker/pkg/models/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	warningLog    *log.Logger
	foods         models.FoodModelInterface
	users         models.UserModelInterface
	templateCache map[string]*template.Template
	session       *sessions.Session
}

type contextKey string

const contextKeyUser = contextKey("user")

func main() {
	// infoLog: logging info messages.
	// Flags used:
	//   log.Ldate      -> the date in the local time zone: 2009/01/23
	//   log.Ltime      -> the time in the local time zone: 01:23:23
	//   log.Lshortfile -> final file name element and line number: handlers.go:42
	// This makes every log line show the log type (INFO), date, time and the
	// affected file with its line number.
	appLogger := logger.Default()
	infoLog := appLogger.Info

	// warningLog: logging non-fatal warnings, prefixed with the log type (WARNING).
	warningLog := appLogger.Warning

	// errorLog: logging error messages with the same date/time and the
	// affected file and line number, prefixed with the log type (ERROR).
	errorLog := appLogger.Error
	var dbUser string
	var dbName string
	var dbHost string
	var dbPassword string
	var dbPort string
	var secret string
	var dsn string
	var port string
	var dbSSLMode string
	var dbCACert string

	// Which which environment the application is running and set the configs properly

	if os.Getenv("ENV") == "prod" {
		dbUser = os.Getenv("DBUSER")
		dbPassword = os.Getenv("DBPASSWORD")
		dbHost = os.Getenv("DBHOST")
		dbName = os.Getenv("DBNAME")
		dbPort = os.Getenv("DBPORT")
		secret = os.Getenv("SESSION")
		port = os.Getenv("PORT")
		dbSSLMode = os.Getenv("DBSSLMODE")
		dbCACert = os.Getenv("DBCACERT")
		if port == "" {
			log.Fatal("Port environment variable not set")
		}
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

		dbUser = os.Getenv("DBUSER")
		dbName = os.Getenv("DBNAME")
		dbHost = os.Getenv("DBHOST")
		dbPassword = os.Getenv("DBPASSWORD")
		dbPort = os.Getenv("DBPORT")
		secret = os.Getenv("SESSION")
		dbSSLMode = os.Getenv("DBSSLMODE")
		dbCACert = os.Getenv("DBCACERT")

	}

	// Managed MySQL providers require an encrypted connection.
	// DBSSLMODE controls how TLS is handled:
	//   - "" or "disable"     -> no TLS (local development)
	//   - "require"/"true"    -> TLS without verifying the server cert
	//   - "verify-ca"         -> TLS verified against the Aiven CA cert (DBCACERT)
	tlsParam := ""
	switch dbSSLMode {
	case "", "disable", "false":
		// no TLS
	case "verify-ca", "verify-full":
		// Register a custom TLS config that trusts the Aiven CA certificate.
		if err := db.RegisterTLSConfig("aiven", dbCACert); err != nil {
			errorLog.Fatal(err)
		}
		tlsParam = "&tls=aiven"
	default:
		// "require"/"true"/"skip-verify": encrypt but skip CA verification.
		tlsParam = "&tls=skip-verify"
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, tlsParam)

	conn, err := db.OpenDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer conn.Close()
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
		warningLog:    warningLog,
		foods:         &mysql.FoodModel{DB: conn},
		users:         &mysql.UserModel{DB: conn},
		templateCache: templateCache,
		session:       session,
	}

	serve := &http.Server{
		Addr:         ":" + port,
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

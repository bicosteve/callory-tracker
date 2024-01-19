package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bicosteve/callory-tracker/pkg/helpers"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Loading env file
	env, err := helpers.LoadEnv("./env/.env")
	if err != nil {
		log.Println(err)
		return
	}

	err = godotenv.Load(env)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")

	app := &application{}

	serve := &http.Server{
		Addr:    port,
		Handler: app.routes(),
	}

	log.Printf("Server running at %s ... \n", port)
	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"log"
	"net/http"
)

const Port = ":4001"

func main() {
	serve := &http.Server{
		Addr:    Port,
		Handler: routes(),
	}

	log.Printf("Server running at %s \n", Port)
	err := serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

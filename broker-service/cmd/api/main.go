package main

import (
	"log"
	"fmt"
	"net/http"
)

type Config struct {
	
}

const port string = "8080"

func main() {
	app := Config{}

	log.Printf("Starting broken service on port %s \n", port)

	// define http server
	srv := &http.Server{
		Addr:  fmt.Sprintf(":%s", port),
		Handler: app.routers(),
	}

	// start http server
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}

}
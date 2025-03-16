package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct{
	Mailer Mail // SMTP Server
}

const port string = "9002"

func main() {
	app := Config{
		Mailer: createMailServer(),
	}

	log.Println("Starting mail service on port: ", port)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: app.routers(),
	}

	err := srv.ListenAndServe() 
	if err != nil {
		log.Panic(err)
	}
}

func createMailServer() Mail {
	port,_ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain: os.Getenv("MAIL_DOMAIN"),
		Host: os.Getenv("MAIL_HOST"),
		Port: port,
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		Encryption: os.Getenv("MAIL_ENCRYPTION"),
		FromName: os.Getenv("FROM_NANE"),
		FromAddr: os.Getenv("FROM_ADDR"),
	}
	return m
}
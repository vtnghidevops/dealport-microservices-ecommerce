package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// import the postgres driver
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const port string = "9000"

type Config struct {
	Database *sql.DB // Connection to the database
	// Models is obj to def in data/models.go that 
	// contain all the models integrated in the database
	Models data.Models 
}

func main() {
	log.Printf("Starting authentication service on %s", port)

	// connect to db
	connect := connectToDB()
	defer connect.Close()
	if connect == nil {
		log.Panic("Could not connect to the database Postgres.")
	}

	// set up config
	app := Config{
		Database: connect,
		Models: data.New(connect), // create a new table in the database => Users
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: app.routers(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func openDB(dsn string) (*sql.DB, error){
	// open the connection to the database
	db, err := sql.Open("pgx", dsn) 
	if err != nil {
		return nil, err
	}

	// check if the connection is working
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	// dsn := "host=localhost port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5"

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Could not connect to the database Postgres. Retrying in 3 seconds.")
			time.Sleep(3 * time.Second)
			continue
		}else {
			log.Println("Connected to the database Postgres.")
			return connection
		}
	}

}
// cmd/api/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"product-service/data"
	"product-service/internal/handler"
	transportHttp "product-service/internal/transport/http"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const port = "8082"

func main() {
	log.Printf("Starting product service on %s", port)

	// Connect to database
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to Postgres")
	}
	defer conn.Close()

	// Set up application config
	handlerConfig := &handler.Config{
		Database: conn,
		Models:   data.New(conn),
	}

	// Create HTTP server
	server := transportHttp.NewServer(handlerConfig)

	// Define HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: server.Routes(),
	}

	// Start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	// For development, provide a fallback if the environment variable is not set
	if dsn == "" {
		// dsn = "host=postgres-products port=5432 user=postgres-products password=password dbname=products sslmode=disable timezone=UTC connect_timeout=5"
		dsn = "host=localhost port=5433 user=postgres-products password=password dbname=products sslmode=disable timezone=UTC connect_timeout=5"
	}

	// Try to connect to the database with retries
	for i := 0; i < 10; i++ {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("PostgreSQL not ready yet...")
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Connected to PostgreSQL!")
		return connection
	}

	log.Println("Could not connect to PostgreSQL after multiple attempts")
	return nil
}

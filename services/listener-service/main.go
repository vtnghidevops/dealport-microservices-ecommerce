package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Printf("Error creating consumer: %v", err)
		os.Exit(1)
	}

	// define topics to listen for
	topics := []string{
		"log.INFO",
		"log.WARNING",
		"log.ERROR",
		"user.registered",
		"user.password_changed",
		"email.send",
		// Add new topics for order events
		"order.created",
		"order.payment_succeeded",
		"order.status_changed",
	}

	// watch the queue and consume events
	err = consumer.Listen(topics)
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// Get RabbitMQ URL from environment variable or use default
	rabbitURL := "amqp://guest:guest@localhost:5672"
	if os.Getenv("RABBITMQ_URL") != "" {
		rabbitURL = os.Getenv("RABBITMQ_URL")
	} else if os.Getenv("RABBITMQ_HOST") != "" {
		// Build URL from individual components if host is specified
		host := os.Getenv("RABBITMQ_HOST")
		user := "guest"
		password := "guest"
		port := "5672"

		if os.Getenv("RABBITMQ_USER") != "" {
			user = os.Getenv("RABBITMQ_USER")
		}
		if os.Getenv("RABBITMQ_PASSWORD") != "" {
			password = os.Getenv("RABBITMQ_PASSWORD")
		}
		if os.Getenv("RABBITMQ_PORT") != "" {
			port = os.Getenv("RABBITMQ_PORT")
		}

		rabbitURL = fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)
	}

	log.Printf("Attempting to connect to RabbitMQ at %s", rabbitURL)

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial(rabbitURL)
		if err != nil {
			log.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Backing off for %v seconds...", backOff.Seconds())
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

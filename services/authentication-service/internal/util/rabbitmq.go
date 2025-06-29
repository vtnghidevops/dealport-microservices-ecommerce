package util

import (
	"fmt"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConnectToRabbitMQ attempts to connect to RabbitMQ with retries
func ConnectToRabbitMQ(rabbitURL string) (*amqp.Connection, error) {
	var counts int64
	var connection *amqp.Connection
	var err error
	var backOff = 1 * time.Second

	// Don't continue until rabbit is ready
	for {
		connection, err = amqp.Dial(rabbitURL)
		if err != nil {
			counts++
			log.Printf("RabbitMQ not yet ready: %v", err)
		} else {
			log.Println("Connected to RabbitMQ!")
			break
		}

		if counts > 5 {
			return nil, fmt.Errorf("cannot connect to RabbitMQ after 5 attempts: %w", err)
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Backing off for %s...", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

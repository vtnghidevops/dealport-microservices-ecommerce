package main

import (
	"listener-service/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Get logger configuration from environment variables
	loggerPrefix := os.Getenv("LOGGER_PREFIX")
	if loggerPrefix == "" {
		loggerPrefix = "listener-service "
	}

	// Create a logger
	logger := log.New(os.Stdout, loggerPrefix, log.LstdFlags)
	logger.Println("Starting listener service...")

	// Connect to RabbitMQ
	rabbitConn, err := connectToRabbitMQ()
	if err != nil {
		logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	// Setup consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		logger.Fatalf("Failed to create consumer: %v", err)
	}

	// Watch the queue for events
	logger.Println("Listening for events...")

	// List all topics we want to listen to
	topics := []string{
		"log.INFO.#",                             // All log events
		"user.registered",                        // User registration
		"user.profile_updated",                   // User profile updated (from user-service)
		"auth.password_reset_requested",          // Password reset requested
		"user.password_changed",                  // Password changed (fixed routing key)
		"auth.password_changed",                  // Password changed (fixed routing key)
		"auth.otp_generated",                     // OTP generation
		"email.send",                             // Email sending
		"order.created",                          // Order creation
		"order.status_changed",                   // Order status changes
		"order.payment_succeeded",                // Payment success
		"order.payment_failed",                   // Payment failure
		"log.INFO.user.login_success",            // User login successful
		"log.INFO.user.login_failed",             // User login failed
		"log.INFO.user.registered",               // User registration
		"log.INFO.user.profile_updated",          // User profile updated
		"log.INFO.user.password_changed",         // User password changed
		"log.INFO.user.password_reset_requested", // User password reset requested
		"log.INFO.user.logout",                   // User logout
		"log.INFO.order.created",                 // Order created log
		"log.INFO.order.status_changed",          // Order status changed log
		"log.INFO.order.payment_succeeded",       // Payment succeeded log
		"log.INFO.order.payment_failed",          // Payment failed log
		"log.INFO.order.cancelled",               // Order cancelled log
		"log.INFO.order.shipped",                 // Order shipped log
		"log.INFO.order.delivered",               // Order delivered log
	}

	logger.Printf("Listener service is watching for the following topics: %v", topics)

	// Listen for messages
	err = consumer.Listen(topics)
	if err != nil {
		logger.Fatalf("Failed to listen for messages: %v", err)
	}
}

// connectToRabbitMQ attempts to connect to RabbitMQ with exponential backoff
func connectToRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var connection *amqp.Connection
	var err error

	// Get RabbitMQ connection details from environment or use defaults
	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	log.Printf("Attempting to connect to RabbitMQ at %s", rabbitURL)

	// Try to connect to RabbitMQ with backoff
	for {
		connection, err = amqp.Dial(rabbitURL)
		if err != nil {
			log.Printf("RabbitMQ not ready yet: %s", err)
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			break
		}

		if counts > 5 {
			log.Println("Could not connect to RabbitMQ after multiple attempts")
			return nil, err
		}

		// Calculate backoff time using exponential backoff
		backOff := time.Second * time.Duration(math.Pow(2, float64(counts)))
		log.Printf("Backing off for %s", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

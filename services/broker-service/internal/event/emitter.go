package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Emitter is responsible for publishing events to RabbitMQ
type Emitter struct {
	connection *amqp.Connection
	logger     *log.Logger
	exchange   string
}

// NewEventEmitter creates a new event emitter with the provided RabbitMQ connection
func NewEventEmitter(conn *amqp.Connection) (*Emitter, error) {
	logger := log.New(os.Stdout, "[BROKER-EVENT] ", log.LstdFlags)

	emitter := &Emitter{
		connection: conn,
		logger:     logger,
		exchange:   "logs_topic",
	}

	err := emitter.setup()
	if err != nil {
		return nil, fmt.Errorf("failed to set up emitter: %w", err)
	}

	return emitter, nil
}

// setup initializes the RabbitMQ exchange
func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer channel.Close()

	return declareExchange(channel)
}

// PublishEvent publishes a standardized event to RabbitMQ
func (e *Emitter) PublishEvent(ctx context.Context, event Event, routingKey string) error {
	// Validate the event
	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	// Update the published timestamp
	event.PublishedAt = time.Now()

	// Convert the event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return e.publishWithRetry(ctx, eventJSON, routingKey, 3)
}

// publishWithRetry attempts to publish a message with retries
func (e *Emitter) publishWithRetry(ctx context.Context, body []byte, routingKey string, maxRetries int) error {
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = e.publish(ctx, body, routingKey)
		if err == nil {
			return nil
		}

		e.logger.Printf("Failed to publish message (attempt %d/%d): %v", attempt, maxRetries, err)

		// Don't sleep on the last attempt
		if attempt < maxRetries {
			// Exponential backoff
			backoff := time.Duration(attempt*attempt) * 100 * time.Millisecond
			time.Sleep(backoff)
		}
	}

	return fmt.Errorf("failed to publish message after %d attempts: %w", maxRetries, err)
}

// publish sends the message to RabbitMQ
func (e *Emitter) publish(ctx context.Context, body []byte, routingKey string) error {
	// Open a channel
	channel, err := e.connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer channel.Close()

	// Set up Dead Letter Exchange
	err = channel.ExchangeDeclare(
		"dead_letter_exchange", // name
		"fanout",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		e.logger.Printf("Warning: Failed to declare dead letter exchange: %v", err)
		// Continue anyway
	}

	// Declare dead letter queue
	_, err = channel.QueueDeclare(
		"dead_letter_queue", // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		e.logger.Printf("Warning: Failed to declare dead letter queue: %v", err)
		// Continue anyway
	}

	// Bind dead letter queue to dead letter exchange
	err = channel.QueueBind(
		"dead_letter_queue",    // queue name
		"",                     // routing key (all messages)
		"dead_letter_exchange", // exchange
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		e.logger.Printf("Warning: Failed to bind dead letter queue: %v", err)
		// Continue anyway
	}

	// Publish the message with a context for timeout
	e.logger.Printf("Publishing message to exchange %s with routing key %s", e.exchange, routingKey)

	err = channel.PublishWithContext(
		ctx,
		e.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
			Timestamp:    time.Now(),
			Headers: amqp.Table{
				"x-retry-count": 0,
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	e.logger.Printf("Successfully published message with routing key %s", routingKey)
	return nil
}

// Push is the legacy method for backward compatibility
// Deprecated: Use PublishEvent instead
func (e *Emitter) Push(event string, routingKey string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()
	e.logger.Println("Pushing to channel:", channel)

	err = channel.Publish(
		"logs_topic",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

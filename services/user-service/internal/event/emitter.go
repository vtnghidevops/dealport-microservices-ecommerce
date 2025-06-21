package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Event represents a standardized message format for all events in the system
type Event struct {
	ID          string      `json:"id"`             // Unique event identifier
	Name        string      `json:"name"`           // Event name/type (e.g., "user.registered")
	Data        interface{} `json:"data,omitempty"` // Event payload data
	DataSchema  string      `json:"data_schema"`    // Schema version for the data
	Source      string      `json:"source"`         // Source service that created the event
	CreatedAt   time.Time   `json:"created_at"`     // Timestamp when event was created
	PublishedAt time.Time   `json:"published_at"`   // Timestamp when event was published
	Version     string      `json:"version"`        // Event schema version
}

// EventEmitter handles publishing events to RabbitMQ
type EventEmitter struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *log.Logger
}

// NewEventEmitter creates a new event emitter
func NewEventEmitter(rabbitURL string, logger *log.Logger) (*EventEmitter, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare the exchange
	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &EventEmitter{
		conn:    conn,
		channel: ch,
		logger:  logger,
	}, nil
}

// EmitUserActivityEvent emits user activity events
func (e *EventEmitter) EmitUserActivityEvent(ctx context.Context, action, userID, message string, metadata map[string]interface{}) error {
	eventData := map[string]interface{}{
		"user_id":   userID,
		"action":    action,
		"message":   message,
		"level":     "INFO",
		"service":   "user-service",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	if metadata != nil {
		eventData["metadata"] = metadata
	}

	event := Event{
		ID:          generateEventID(),
		Name:        fmt.Sprintf("user.%s", action),
		Data:        eventData,
		DataSchema:  "v1",
		Source:      "user-service",
		CreatedAt:   time.Now(),
		PublishedAt: time.Now(),
		Version:     "v1",
	}

	return e.publishEvent(ctx, event, fmt.Sprintf("log.INFO.user.%s", action))
}

// EmitPasswordChangedEvent emits password changed events
// func (e *EventEmitter) EmitPasswordChangedEvent(ctx context.Context, userID string) error {
// 	eventData := map[string]interface{}{
// 		"user_id":    userID,
// 		"changed_at": time.Now().Format(time.RFC3339),
// 	}

// 	event := Event{
// 		ID:          generateEventID(),
// 		Name:        "user.password_changed",
// 		Data:        eventData,
// 		DataSchema:  "v1",
// 		Source:      "user-service",
// 		CreatedAt:   time.Now(),
// 		PublishedAt: time.Now(),
// 		Version:     "v1",
// 	}

// 	return e.publishEvent(ctx, event, "log.INFO.user.password_changed")
// }

// EmitProfileUpdatedEvent emits profile updated events
func (e *EventEmitter) EmitProfileUpdatedEvent(ctx context.Context, userID string, updatedFields map[string]interface{}) error {
	eventData := map[string]interface{}{
		"user_id":        userID,
		"updated_at":     time.Now().Format(time.RFC3339),
		"updated_fields": updatedFields,
	}

	event := Event{
		ID:          generateEventID(),
		Name:        "user.profile_updated",
		Data:        eventData,
		DataSchema:  "v1",
		Source:      "user-service",
		CreatedAt:   time.Now(),
		PublishedAt: time.Now(),
		Version:     "v1",
	}

	return e.publishEvent(ctx, event, "log.INFO.user.profile_updated")
}

// EmitLogoutEvent emits logout events
func (e *EventEmitter) EmitLogoutEvent(ctx context.Context, userID, email string) error {
	eventData := map[string]interface{}{
		"user_id":   userID,
		"email":     email,
		"logout_at": time.Now().Format(time.RFC3339),
		"action":    "logout",
		"level":     "INFO",
		"service":   "user-service",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	event := Event{
		ID:          generateEventID(),
		Name:        "user.logout",
		Data:        eventData,
		DataSchema:  "v1",
		Source:      "user-service",
		CreatedAt:   time.Now(),
		PublishedAt: time.Now(),
		Version:     "v1",
	}

	return e.publishEvent(ctx, event, "log.INFO.user.logout")
}

// publishEvent publishes an event to RabbitMQ
func (e *EventEmitter) publishEvent(ctx context.Context, event Event, routingKey string) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = e.channel.PublishWithContext(
		ctx,
		"logs_topic", // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		e.logger.Printf("Failed to publish event %s: %v", event.Name, err)
		return fmt.Errorf("failed to publish event: %w", err)
	}

	e.logger.Printf("📤 Event published: %s (ID: %s)", event.Name, event.ID)
	return nil
}

// Close closes the connection to RabbitMQ
func (e *EventEmitter) Close() error {
	if e.channel != nil {
		e.channel.Close()
	}
	if e.conn != nil {
		return e.conn.Close()
	}
	return nil
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

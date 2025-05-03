package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// StandardEvent represents a standardized message format for all events
type StandardEvent struct {
	ID          string      `json:"id"`             // Unique event identifier
	Name        string      `json:"name"`           // Event name/type (e.g., "user.registered")
	Data        interface{} `json:"data,omitempty"` // Event payload data
	DataSchema  string      `json:"data_schema"`    // Schema version for the data
	Source      string      `json:"source"`         // Source service that created the event
	CreatedAt   time.Time   `json:"created_at"`     // Timestamp when event was created
	PublishedAt time.Time   `json:"published_at"`   // Timestamp when event was published
	Version     string      `json:"version"`        // Event schema version
}

// UserRegisteredData represents the data for user.registered events
type UserRegisteredData struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	DisplayName  string    `json:"display_name,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	ProfileImage string    `json:"profile_image,omitempty"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Gender       string    `json:"gender,omitempty"`
}

// EmailData represents the data structure for email.send events
type EmailData struct {
	Type      string            `json:"type"`
	From      string            `json:"from,omitempty"`
	FromName  string            `json:"from_name,omitempty"`
	To        string            `json:"to"`
	Subject   string            `json:"subject"`
	Template  string            `json:"template,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
}

// PasswordResetData represents the data for password reset events
type PasswordResetData struct {
	Email     string `json:"email"`
	TokenHash string `json:"token_hash"`
	ExpiresAt string `json:"expires_at"`
}

// OTPGeneratedData represents the data for auth.otp_generated events
type OTPGeneratedData struct {
	Email      string `json:"email"`
	OTP        string `json:"otp"`
	Purpose    string `json:"purpose"`
	ExpiresIn  int    `json:"expires_in"`
	Message    string `json:"message"`
	ActionText string `json:"action_text"`
}

// Emitter is a wrapper around RabbitMQ connection for sending events
type Emitter struct {
	conn   *amqp.Connection
	logger *log.Logger
}

// NewEmitter creates a new event emitter
func NewEmitter(conn *amqp.Connection, logger *log.Logger) (*Emitter, error) {
	if conn == nil {
		return nil, fmt.Errorf("RabbitMQ connection is nil")
	}

	emitter := &Emitter{
		conn:   conn,
		logger: logger,
	}

	// Ensure the exchange exists
	if err := emitter.setupExchange(); err != nil {
		return nil, err
	}

	return emitter, nil
}

// setupExchange declares the exchange
func (e *Emitter) setupExchange() error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

// EmitUserRegistered emits a user.registered event
func (e *Emitter) EmitUserRegistered(user UserRegisteredData) error {
	// Log user data for debugging with more details
	e.logger.Printf("DEBUG EmitUserRegistered: Starting to emit event for user %s (%s)", user.ID, user.Email)
	e.logger.Printf("DEBUG EmitUserRegistered: User detail - ID: %s, Email: %s, Username: %s, FirstName: %s, LastName: %s, Phone: '%s', Role: %s, Status: %s, Active: %t",
		user.ID, user.Email, user.Username, user.FirstName, user.LastName, user.Phone, user.Role, user.Status, user.Active)

	// Tạo custom map để serialize trực tiếp, đảm bảo tất cả trường được bao gồm
	customData := map[string]interface{}{
		"id":           user.ID,
		"email":        user.Email,
		"username":     user.Username,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"display_name": user.DisplayName,
		"phone":        user.Phone,
		"role":         user.Role,
		"status":       user.Status,
		"active":       user.Active,
		"created_at":   user.CreatedAt,
		"updated_at":   user.UpdatedAt,
	}

	event := StandardEvent{
		ID:         uuid.New().String(),
		Name:       "user.registered",
		Data:       customData,
		DataSchema: "v1",
		Source:     "authentication-service",
		CreatedAt:  time.Now(),
		Version:    "1.0",
	}

	// Convert event to JSON for logging
	eventJSON, _ := json.Marshal(event)
	e.logger.Printf("DEBUG EmitUserRegistered: Full event JSON: %s", string(eventJSON))

	// Kiểm tra xem trường phone có trong JSON không
	var dataCheck map[string]interface{}
	json.Unmarshal(eventJSON, &dataCheck)
	if data, ok := dataCheck["data"].(map[string]interface{}); ok {
		e.logger.Printf("DEBUG EmitUserRegistered: Phone field in JSON: '%v'", data["phone"])
	}

	// Specify routing key explicitly
	routingKey := "user.registered"
	e.logger.Printf("DEBUG EmitUserRegistered: Publishing with routing key: %s", routingKey)

	// Get channel
	ch, err := e.conn.Channel()
	if err != nil {
		e.logger.Printf("ERROR EmitUserRegistered: Failed to get channel: %v", err)
		return fmt.Errorf("failed to get RabbitMQ channel: %w", err)
	}
	defer ch.Close()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set PublishedAt time just before publishing
	event.PublishedAt = time.Now()

	// Update JSON with PublishedAt time
	eventJSON, _ = json.Marshal(event)

	// Direct publication for maximum visibility into the process
	err = ch.PublishWithContext(
		ctx,
		"logs_topic", // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventJSON,
		},
	)

	if err != nil {
		e.logger.Printf("ERROR EmitUserRegistered: Failed to publish event: %v", err)
		return fmt.Errorf("failed to publish user.registered event: %w", err)
	}

	e.logger.Printf("SUCCESS EmitUserRegistered: Event published successfully with ID %s and routing key %s", event.ID, routingKey)
	return nil
}

// EmitOTPGenerated emits an auth.otp_generated event
func (e *Emitter) EmitOTPGenerated(email, otp, purpose string, expiresIn int, message, actionText string) error {
	data := OTPGeneratedData{
		Email:      email,
		OTP:        otp,
		Purpose:    purpose,
		ExpiresIn:  expiresIn,
		Message:    message,
		ActionText: actionText,
	}

	event := StandardEvent{
		ID:         uuid.New().String(),
		Name:       "auth.otp_generated",
		Data:       data,
		DataSchema: "v1",
		Source:     "authentication-service",
		CreatedAt:  time.Now(),
		Version:    "1.0",
	}

	return e.publish("email.send", event)
}

// EmitPasswordChanged emits a user.password_changed event
func (e *Emitter) EmitPasswordChanged(email string) error {
	data := map[string]string{
		"email":      email,
		"changed_at": time.Now().Format(time.RFC3339),
	}

	event := StandardEvent{
		ID:         uuid.New().String(),
		Name:       "user.password_changed",
		Data:       data,
		DataSchema: "v1",
		Source:     "authentication-service",
		CreatedAt:  time.Now(),
		Version:    "1.0",
	}

	return e.publish("user.password_changed", event)
}

// EmitPasswordResetRequested emits a auth.password_reset_requested event
func (e *Emitter) EmitPasswordResetRequested(email, tokenHash string, expiresAt time.Time) error {
	data := PasswordResetData{
		Email:     email,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}

	event := StandardEvent{
		ID:         uuid.New().String(),
		Name:       "auth.password_reset_requested",
		Data:       data,
		DataSchema: "v1",
		Source:     "authentication-service",
		CreatedAt:  time.Now(),
		Version:    "1.0",
	}

	return e.publish("auth.password_reset_requested", event)
}

// publish publishes an event to RabbitMQ
func (e *Emitter) publish(routingKey string, event StandardEvent) error {
	ch, err := e.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}
	defer ch.Close()

	// Set the PublishedAt timestamp just before publishing
	event.PublishedAt = time.Now()

	// Convert event to JSON
	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Publish the message
	err = ch.PublishWithContext(
		ctx,
		"logs_topic", // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	e.logger.Printf("Event %s published with routing key %s", event.Name, routingKey)
	return nil
}

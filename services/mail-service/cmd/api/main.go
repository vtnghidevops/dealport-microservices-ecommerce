package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mail-service/internal/config"
	"mail-service/internal/mailer"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config holds configuration for mail service
type Config struct {
	Config config.Config // Service configuration
}

const port string = "9002"

// MailPayload for test endpoint
type MailPayload struct {
	From        string            `json:"from,omitempty"`
	FromName    string            `json:"from_name,omitempty"`
	To          string            `json:"to"`
	Subject     string            `json:"subject"`
	Message     string            `json:"message"`
	Template    string            `json:"template,omitempty"`
	Attachments []string          `json:"attachments,omitempty"`
	Data        map[string]string `json:"data,omitempty"`
}

// EventPayload holds data from RabbitMQ events
type EventPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// JSONResponse for API responses
type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

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

func main() {
	log.Println("Starting mail service...")

	// Try to find .env file in various locations
	envFiles := []string{
		".env",
		"../../.env",
		"../../../.env",
		"../../../../.env", // Project root directory
	}

	for _, envFile := range envFiles {
		absPath, _ := filepath.Abs(envFile)
		err := loadEnvFile(absPath)
		if err == nil {
			log.Printf("Loaded environment from: %s", absPath)
			break // Stop after first successful loading
		}
	}

	// Load configuration
	cfg := config.LoadConfig()

	app := Config{
		Config: cfg,
	}

	// Start HTTP server for Postman testing
	// go app.serveHTTP()

	// Connect to RabbitMQ and start consuming
	go app.setupRabbitMQConsumer()

	log.Println("Mail service started")
	log.Printf("HTTP endpoint available at http://localhost:%s/send (for Postman testing)", port)

	// Keep the application running indefinitely
	select {}
}

// setupRabbitMQConsumer connects to RabbitMQ and starts consuming email.send events
func (app *Config) setupRabbitMQConsumer() {
	// Connect to RabbitMQ with retries - try multiple times with delays
	var rabbitConn *amqp.Connection
	var err error
	var retry bool
	maxRetries := 10
	retryCount := 0

	for {
		rabbitConn, err = connectToRabbitMQ()
		if err != nil {
			retryCount++
			log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v", retryCount, maxRetries, err)

			if retryCount >= maxRetries {
				log.Printf("Cannot connect to RabbitMQ after %d attempts. Will try again in 30 seconds.", maxRetries)
				time.Sleep(30 * time.Second)
				retryCount = 0
				continue
			}

			// Wait before retrying
			delay := time.Duration(math.Pow(2, float64(retryCount))) * time.Second
			log.Printf("Retrying in %v...", delay)
			time.Sleep(delay)
			continue
		}

		// Successfully connected
		defer rabbitConn.Close()
		log.Println("Successfully connected to RabbitMQ, starting consumer...")

		// Start consuming from RabbitMQ
		retry = app.consumeFromRabbitMQ(rabbitConn)

		// If consumeFromRabbitMQ returns true, we should retry connecting
		if retry {
			log.Println("Will attempt to reconnect to RabbitMQ...")
			time.Sleep(5 * time.Second)
			continue
		}

		// If we get here, we're done (likely due to service shutdown)
		break
	}
}

// connectToRabbitMQ connects to RabbitMQ with retries
func connectToRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var connection *amqp.Connection
	var err error
	var backOff = 1 * time.Second

	// Don't continue until rabbit is ready
	// Use localhost instead of container name for local development
	rabbitURL := "amqp://guest:guest@localhost:5672"
	log.Printf("Attempting to connect to RabbitMQ at %s", rabbitURL)

	for {
		connection, err = amqp.Dial(rabbitURL)
		if err != nil {
			counts++
			log.Printf("RabbitMQ not yet ready: %v", err)

			if counts > 5 {
				return nil, fmt.Errorf("cannot connect to RabbitMQ after 5 attempts: %w", err)
			}

			backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
			log.Printf("Backing off for %s...", backOff)
			time.Sleep(backOff)
			continue
		}

		log.Println("Connected to RabbitMQ!")
		break
	}

	return connection, nil
}

// consumeFromRabbitMQ consumes email events from RabbitMQ
// Returns true if we should attempt to reconnect
func (app *Config) consumeFromRabbitMQ(conn *amqp.Connection) bool {
	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open RabbitMQ channel: %v", err)
		return true
	}
	defer ch.Close()

	// Declare exchange
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
		log.Printf("Failed to declare exchange: %v", err)
		return true
	}

	// Declare queue
	q, err := ch.QueueDeclare(
		"mail_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		return true
	}

	// Bind queue to exchange with routing key for email.send events
	err = ch.QueueBind(
		q.Name,       // queue name
		"email.send", // routing key
		"logs_topic", // exchange
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Printf("Failed to bind queue: %v", err)
		return true
	}

	// Set QoS (Quality of Service)
	err = ch.Qos(
		1,     // prefetch count (process one message at a time)
		0,     // prefetch size (no specific size limit)
		false, // global (false means apply to just this channel)
	)
	if err != nil {
		log.Printf("Failed to set QoS: %v", err)
		return true
	}

	// Consume messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Printf("Failed to register consumer: %v", err)
		return true
	}

	// Listen for connection close events
	closeChannel := make(chan *amqp.Error)
	conn.NotifyClose(closeChannel)

	// Forever channel to keep the goroutine running
	forever := make(chan bool)

	// Process messages
	go func() {
		for {
			select {
			case closeErr := <-closeChannel:
				log.Printf("RabbitMQ connection closed: %v", closeErr)
				forever <- true
				return

			case d, ok := <-msgs:
				if !ok {
					log.Println("RabbitMQ channel closed")
					forever <- true
					return
				}

				// Process the message
				// Try to process as a standard event first
				var stdEvent StandardEvent
				stdErr := json.Unmarshal(d.Body, &stdEvent)

				if stdErr == nil && stdEvent.Name != "" {
					// Process as a standard event
					log.Printf("Received standard event: %s (ID: %s)", stdEvent.Name, stdEvent.ID)

					if stdEvent.Name == "auth.otp_generated" {
						// Handle OTP event directly
						err := app.handleOTPGeneratedEvent(stdEvent)
						if err != nil {
							log.Printf("Error processing OTP event: %v", err)
						} else {
							log.Printf("Successfully processed OTP event")
						}
					} else if stdEvent.Name == "email.send" {
						// Process standard email event
						err := app.handleStandardEmailEvent(stdEvent)
						if err != nil {
							log.Printf("Error processing standard email event: %v", err)
						} else {
							log.Printf("Successfully processed standard email event")
						}
					} else {
						log.Printf("Unexpected event name: %s", stdEvent.Name)
					}
				} else {
					// Try legacy format
					var payload EventPayload
					legacyErr := json.Unmarshal(d.Body, &payload)

					if legacyErr == nil {
						log.Printf("Received message from RabbitMQ: %s", payload.Name)

						// Check if it's an email.send event
						if payload.Name == "email.send" {
							// Extract mail information from payload
							var emailData map[string]string
							err := json.Unmarshal([]byte(payload.Data), &emailData)
							if err != nil {
								log.Printf("Error unmarshalling email data: %v", err)
								continue
							}

							// Determine email type and process accordingly
							emailType := emailData["type"]
							var processErr error

							switch emailType {
							case "registration":
								processErr = app.handleRegistrationEmail(payload.Data)
							case "reset_password":
								processErr = app.handlePasswordResetEmail(payload.Data)
							case "password_change":
								processErr = app.handlePasswordChangeEmail(payload.Data)
							case "order_confirmation":
								processErr = app.handleOrderConfirmationEmail(payload.Data)
							default:
								// If no specific type or unknown, send generic email
								processErr = app.handleGenericEmail(payload.Data)
							}

							if processErr != nil {
								log.Printf("Error processing email: %v", processErr)
							} else {
								log.Printf("Successfully sent email of type: %s", emailType)
							}
						} else {
							log.Printf("Unexpected event type: %s - expected email.send", payload.Name)
						}
					} else {
						log.Printf("Failed to parse message in any format: %v", d.Body)
					}
				}
			}
		}
	}()

	log.Printf("Mail service listening for 'email.send' events on RabbitMQ (queue: %s)", q.Name)

	// Wait until we're told to stop
	<-forever
	log.Println("RabbitMQ consumer shutting down")
	return true
}

// handleGenericEmail processes a generic email request
func (app *Config) handleGenericEmail(data string) error {
	// Parse event data
	var emailData map[string]string
	err := json.Unmarshal([]byte(data), &emailData)
	if err != nil {
		return fmt.Errorf("error unmarshalling email data: %w", err)
	}

	// Set default template if not specified
	template := emailData["template"]
	if template == "" {
		template = "mail.html.gohtml"
	}

	msg := mailer.Message{
		From:     emailData["from"],
		FromName: emailData["from_name"],
		To:       emailData["to"],
		Subject:  emailData["subject"],
		Template: template,
		Data:     emailData,
	}

	// Send email
	return app.Config.Mailer.SendSMTPMessage(msg)
}

// serveHTTP starts the HTTP server for testing with Postman
func (app *Config) serveHTTP() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	log.Printf("Starting HTTP server on port %s", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// routes defines HTTP routes for testing with Postman
func (app *Config) routes() http.Handler {
	mux := http.NewServeMux()

	// Health check route
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mail service is alive!"))
	})

	// Test email route (for Postman testing)
	mux.HandleFunc("/send", app.SendMailHandler)

	// OTP email route
	mux.HandleFunc("/send-otp", app.SendOTPMailHandler)

	// Test route for receiving events directly (for manual testing)
	mux.HandleFunc("/handle", app.HandleEvent)

	return mux
}

// HandleEvent handles events from direct HTTP requests (for testing)
func (app *Config) HandleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var eventPayload EventPayload
	err := json.NewDecoder(r.Body).Decode(&eventPayload)
	if err != nil {
		log.Printf("Error decoding event payload: %v", err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("Received test event via HTTP: %s", eventPayload.Name)

	// Handle different types of email events
	switch eventPayload.Name {
	case "mail.registration":
		err = app.handleRegistrationEmail(eventPayload.Data)
	case "mail.reset_password":
		err = app.handlePasswordResetEmail(eventPayload.Data)
	case "mail.password_change":
		err = app.handlePasswordChangeEmail(eventPayload.Data)
	case "mail.order_confirmation":
		err = app.handleOrderConfirmationEmail(eventPayload.Data)
	case "email.send":
		err = app.handleGenericEmail(eventPayload.Data)
	default:
		err = fmt.Errorf("unknown event type: %s", eventPayload.Name)
	}

	if err != nil {
		log.Printf("Error handling event %s: %v", eventPayload.Name, err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully processed test event: %s", eventPayload.Name)

	app.writeJSON(w, http.StatusAccepted, JSONResponse{
		Error:   false,
		Message: "Event processed successfully",
	})
}

// handleRegistrationEmail processes registration email event
func (app *Config) handleRegistrationEmail(data string) error {
	// Parse event data
	var emailData map[string]string
	err := json.Unmarshal([]byte(data), &emailData)
	if err != nil {
		return fmt.Errorf("error unmarshalling registration email data: %w", err)
	}

	msg := mailer.Message{
		From:     "",
		FromName: "",
		To:       emailData["email"],
		Subject:  "Registration Confirmation",
		Template: "register.html.gohtml",
		Data:     emailData,
	}

	// Send email
	return app.Config.Mailer.SendSMTPMessage(msg)
}

// handlePasswordResetEmail processes password reset email event
func (app *Config) handlePasswordResetEmail(data string) error {
	// Parse event data
	var emailData map[string]string
	err := json.Unmarshal([]byte(data), &emailData)
	if err != nil {
		return fmt.Errorf("error unmarshalling password reset email data: %w", err)
	}

	msg := mailer.Message{
		From:     "",
		FromName: "",
		To:       emailData["email"],
		Subject:  "Password Reset Request",
		Template: "reset_password.html.gohtml",
		Data:     emailData,
	}

	// Send email
	return app.Config.Mailer.SendSMTPMessage(msg)
}

// handlePasswordChangeEmail processes password change email event
func (app *Config) handlePasswordChangeEmail(data string) error {
	// Parse event data
	var emailData map[string]string
	err := json.Unmarshal([]byte(data), &emailData)
	if err != nil {
		return fmt.Errorf("error unmarshalling password change email data: %w", err)
	}

	msg := mailer.Message{
		From:     "",
		FromName: "",
		To:       emailData["email"],
		Subject:  "Password Changed",
		Template: "password_change.html.gohtml",
		Data:     emailData,
	}

	// Send email
	return app.Config.Mailer.SendSMTPMessage(msg)
}

// handleOrderConfirmationEmail processes order confirmation email event
func (app *Config) handleOrderConfirmationEmail(data string) error {
	// Parse event data
	var emailData map[string]string
	err := json.Unmarshal([]byte(data), &emailData)
	if err != nil {
		return fmt.Errorf("error unmarshalling order confirmation email data: %w", err)
	}

	msg := mailer.Message{
		From:     "",
		FromName: "",
		To:       emailData["email"],
		Subject:  "Order Confirmation",
		Template: "order_confirmation.html.gohtml",
		Data:     emailData,
	}

	// Send email
	return app.Config.Mailer.SendSMTPMessage(msg)
}

// handleStandardEmailEvent processes a standardized email.send event
func (app *Config) handleStandardEmailEvent(event StandardEvent) error {
	// Extract email data
	var emailData EmailData

	// Handle different data formats
	switch data := event.Data.(type) {
	case EmailData:
		emailData = data
	case map[string]interface{}:
		// Convert map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal email data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &emailData); err != nil {
			return fmt.Errorf("failed to unmarshal email data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &emailData); err != nil {
			return fmt.Errorf("failed to unmarshal email data string: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type for email.send event")
	}

	// Log detailed info about email type being processed
	log.Printf("Processing email of type: %s, to: %s", emailData.Type, emailData.To)

	// Determine email type and process accordingly
	switch emailType := emailData.Type; emailType {
	case "registration":
		// Process directly with EmailData
		return app.handleRegistrationEmailWithData(emailData)

	case "reset_password":
		// Process directly with EmailData
		return app.handlePasswordResetEmailWithData(emailData)

	case "password_change":
		// Process directly with EmailData
		return app.handlePasswordChangeEmailWithData(emailData)

	case "order_confirmation":
		// Process directly with EmailData
		return app.handleOrderConfirmationEmailWithData(emailData)

	default:
		// Send directly using the mail template
		msg := mailer.Message{
			From:     emailData.From,
			FromName: emailData.FromName,
			To:       emailData.To,
			Subject:  emailData.Subject,
			Template: emailData.Template,
			Data:     emailData.Variables,
		}

		// If template not specified, use default
		if msg.Template == "" {
			msg.Template = "mail.html.gohtml"
		}

		return app.Config.Mailer.SendSMTPMessage(msg)
	}
}

// handleOrderConfirmationEmailWithData processes an order confirmation email using EmailData directly
func (app *Config) handleOrderConfirmationEmailWithData(emailData EmailData) error {
	// Validate required fields
	if emailData.To == "" {
		return fmt.Errorf("recipient (to) email is required")
	}

	// Create mail message with all the data passed from EmailData
	msg := mailer.Message{
		From:     emailData.From,
		FromName: emailData.FromName,
		To:       emailData.To,
		Subject:  emailData.Subject,
		Template: "order_confirmation.html.gohtml", // Use proper template
		Data:     emailData.Variables,
	}

	// Override template if not empty in EmailData
	if emailData.Template != "" {
		msg.Template = emailData.Template
	}

	log.Printf("Sending order confirmation email to %s", emailData.To)

	// Send email
	return app.Config.Mailer.SendSMTPMessage(msg)
}

// handleRegistrationEmailWithData processes a registration welcome email using EmailData directly
func (app *Config) handleRegistrationEmailWithData(emailData EmailData) error {
	// Validate required fields
	if emailData.To == "" {
		return fmt.Errorf("recipient (to) email is required")
	}

	// Create mail message
	msg := mailer.Message{
		From:     emailData.From,
		FromName: emailData.FromName,
		To:       emailData.To,
		Subject:  emailData.Subject,
		Template: "welcome.html.gohtml", // Use proper template
		Data:     emailData.Variables,
	}

	// Override template if not empty in EmailData
	if emailData.Template != "" {
		msg.Template = emailData.Template
	}

	log.Printf("Sending registration welcome email to %s", emailData.To)

	// Send email
	return app.Config.Mailer.SendSMTPMessage(msg)
}

// handlePasswordResetEmailWithData processes a password reset email using EmailData directly
func (app *Config) handlePasswordResetEmailWithData(emailData EmailData) error {
	// Validate required fields
	if emailData.To == "" {
		return fmt.Errorf("recipient (to) email is required")
	}

	// Create mail message
	msg := mailer.Message{
		From:     emailData.From,
		FromName: emailData.FromName,
		To:       emailData.To,
		Subject:  emailData.Subject,
		Template: "password_reset.html.gohtml", // Use proper template
		Data:     emailData.Variables,
	}

	// Override template if not empty in EmailData
	if emailData.Template != "" {
		msg.Template = emailData.Template
	}

	log.Printf("Sending password reset email to %s", emailData.To)

	// Send email
	return app.Config.Mailer.SendSMTPMessage(msg)
}

// handlePasswordChangeEmailWithData processes a password change notification email using EmailData directly
func (app *Config) handlePasswordChangeEmailWithData(emailData EmailData) error {
	// Validate required fields
	if emailData.To == "" {
		return fmt.Errorf("recipient (to) email is required")
	}

	// Create mail message
	msg := mailer.Message{
		From:     emailData.From,
		FromName: emailData.FromName,
		To:       emailData.To,
		Subject:  emailData.Subject,
		Template: "password_changed.html.gohtml", // Use proper template
		Data:     emailData.Variables,
	}

	// Override template if not empty in EmailData
	if emailData.Template != "" {
		msg.Template = emailData.Template
	}

	log.Printf("Sending password change notification email to %s", emailData.To)

	// Send email
	return app.Config.Mailer.SendSMTPMessage(msg)
}

// SendMailHandler handles test email sending
func (app *Config) SendMailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var payload MailPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Create mail message
	var data interface{}
	if payload.Data != nil {
		data = payload.Data
	} else {
		data = map[string]string{
			"message": payload.Message,
		}
	}

	// Use default template if not specified
	if payload.Template == "" {
		payload.Template = "mail.html.gohtml"
	}

	msg := mailer.Message{
		From:        payload.From,
		FromName:    payload.FromName,
		To:          payload.To,
		Subject:     payload.Subject,
		Template:    payload.Template,
		Attachments: payload.Attachments,
		Data:        data,
	}

	// Send mail
	if err := app.Config.Mailer.SendSMTPMessage(msg); err != nil {
		log.Printf("Error sending mail: %v", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusAccepted, JSONResponse{
		Error:   false,
		Message: "Mail sent successfully",
	})
}

// SendOTPMailHandler handles sending OTP emails
func (app *Config) SendOTPMailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorJSON(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		Email     string `json:"email"`
		OTP       string `json:"otp"`
		Action    string `json:"action"`
		Message   string `json:"message"`
		ExpiresIn int    `json:"expires_in"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Default expiration time if not provided
	if payload.ExpiresIn == 0 {
		payload.ExpiresIn = 15 // Default expiration: 15 minutes
	}

	// Map payload to template data
	data := map[string]interface{}{
		"email":      payload.Email,
		"otp":        payload.OTP,
		"action":     payload.Action,
		"message":    payload.Message,
		"expires_in": payload.ExpiresIn,
		"year":       time.Now().Year(),
	}

	// Create mail message
	msg := mailer.Message{
		From:     "",
		FromName: "",
		To:       payload.Email,
		Subject:  "Your Verification Code",
		Template: "otp_verification.html.gohtml",
		Data:     data,
	}

	// Send mail
	if err := app.Config.Mailer.SendSMTPMessage(msg); err != nil {
		log.Printf("Error sending OTP mail: %v", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusAccepted, JSONResponse{
		Error:   false,
		Message: "OTP email sent successfully",
	})
}

// writeJSON writes JSON response
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// errorJSON writes error as JSON
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

// loadEnvFile loads environment variables from .env file
func loadEnvFile(filename string) error {
	// Check if file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filename)
	}

	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Parse and set environment variables
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			value = strings.Trim(value, `"'`)
			os.Setenv(key, value)
		}
	}

	return nil
}

// handleOTPGeneratedEvent processes auth.otp_generated events
func (app *Config) handleOTPGeneratedEvent(event StandardEvent) error {
	// Extract OTP data from event
	var otpData struct {
		Email      string `json:"email"`
		OTP        string `json:"otp"`
		Purpose    string `json:"purpose"`
		ExpiresIn  int    `json:"expires_in"`
		Message    string `json:"message"`
		ActionText string `json:"action_text"`
	}

	// Convert the data interface to our expected structure
	switch data := event.Data.(type) {
	case map[string]interface{}:
		// Convert map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal OTP data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &otpData); err != nil {
			return fmt.Errorf("failed to unmarshal OTP data: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type for OTP event")
	}

	// Prepare variables for template
	variables := map[string]interface{}{
		"otp":         otpData.OTP,
		"purpose":     otpData.Purpose,
		"expires_in":  fmt.Sprintf("%d", otpData.ExpiresIn),
		"message":     otpData.Message,
		"action_text": otpData.ActionText,
		"year":        time.Now().Year(),
	}

	// Set subject based on purpose
	subject := "Your Verification Code"
	if otpData.Purpose == "registration" {
		subject = "Complete Your Registration"
	} else if otpData.Purpose == "password_reset" {
		subject = "Password Reset Verification"
	}

	// Create mail message
	msg := mailer.Message{
		To:       otpData.Email,
		Subject:  subject,
		Template: "otp_verification.html.gohtml",
		Data:     variables,
	}

	// Send email
	err := app.Config.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Printf("Error sending OTP email: %v", err)
		return err
	}

	log.Printf("OTP email sent successfully to %s (purpose: %s)", otpData.Email, otpData.Purpose)
	return nil
}

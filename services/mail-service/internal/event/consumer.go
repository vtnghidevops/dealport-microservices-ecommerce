package event

import (
	"encoding/json"
	"fmt"
	"log"

	"mail-service/internal/mailer"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer is an interface for RabbitMQ message consuming
type Consumer struct {
	conn   *amqp.Connection
	mailer mailer.Mail
	logger *log.Logger
}

// StandardEvent represents a standardized message format for all events
type StandardEvent struct {
	ID          string      `json:"id"`             // Unique event identifier
	Name        string      `json:"name"`           // Event name/type (e.g., "user.registered")
	Data        interface{} `json:"data,omitempty"` // Event payload data
	DataSchema  string      `json:"data_schema"`    // Schema version for the data
	Source      string      `json:"source"`         // Source service that created the event
	CreatedAt   string      `json:"created_at"`     // Timestamp when event was created
	PublishedAt string      `json:"published_at"`   // Timestamp when event was published
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

// NewConsumer creates a new consumer for RabbitMQ events
func NewConsumer(conn *amqp.Connection, mailer mailer.Mail) (Consumer, error) {
	consumer := Consumer{
		conn:   conn,
		mailer: mailer,
		logger: log.New(log.Writer(), "[EVENT-CONSUMER] ", log.LstdFlags),
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// setup creates a channel from connection and declares exchange
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// Message payload struct received from queue (legacy format)
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// EmailPayload is the data structure for email events (legacy format)
type EmailPayload struct {
	From     string            `json:"from"`
	To       string            `json:"to"`
	Subject  string            `json:"subject"`
	Template string            `json:"template"`
	Data     map[string]string `json:"data"`
}

// Listen for messages on specified topics
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	// Bind queue with exchange for each topic
	for _, bindingKey := range topics {
		err = ch.QueueBind(
			q.Name,
			bindingKey,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	// Start consuming messages
	messages, err := ch.Consume(
		q.Name,
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	// Create a forever channel to keep the program running
	forever := make(chan bool)

	// Process messages in a goroutine
	go func() {
		for d := range messages {
			// First try to process as StandardEvent
			var stdEvent StandardEvent
			stdErr := json.Unmarshal(d.Body, &stdEvent)

			if stdErr == nil && stdEvent.Name != "" {
				consumer.logger.Printf("Received standard event: %s from %s", stdEvent.Name, stdEvent.Source)
				go consumer.handleStandardEvent(stdEvent)
				continue
			}

			// If not a standard event, try legacy format
			var payload Payload
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				consumer.logger.Printf("Error deserializing message: %v", err)
				continue
			}

			consumer.logger.Printf("Received legacy event: %s", payload.Name)
			go consumer.handleMessage(payload)
		}
	}()

	consumer.logger.Printf("Waiting for messages on topics: %v", topics)
	<-forever

	return nil
}

// handleStandardEvent processes the standardized event format
func (consumer *Consumer) handleStandardEvent(event StandardEvent) {
	// Check routing key based on event name
	if event.Name == "auth.otp_generated" {
		consumer.logger.Printf("Processing email.send event: %s", event.Name)
		consumer.handleOtpGeneratedEvent(event)
		return
	}

	// For "email.send" events
	if event.Name == "email.send" {
		consumer.handleEmailSendEvent(event)
		return
	}
}

// handleEmailSendEvent processes email.send standard events
func (consumer *Consumer) handleEmailSendEvent(event StandardEvent) {
	// Extract EmailData from event
	var emailData EmailData

	// Handle different data formats
	switch data := event.Data.(type) {
	case map[string]interface{}:
		// Convert map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			consumer.logger.Printf("Error marshaling email data: %v", err)
			return
		}
		if err := json.Unmarshal(jsonData, &emailData); err != nil {
			consumer.logger.Printf("Error unmarshaling email data: %v", err)
			return
		}
	default:
		consumer.logger.Printf("Unsupported data type for email event")
		return
	}

	// Determine template based on type if not specified
	if emailData.Template == "" {
		switch emailData.Type {
		case "registration":
			emailData.Template = "welcome.html.gohtml"
		case "reset_password":
			emailData.Template = "password_reset.html.gohtml"
		case "password_change":
			emailData.Template = "password_changed.html.gohtml"
		case "order_confirmation":
			emailData.Template = "order_confirmation.html.gohtml"
		case "otp":
			emailData.Template = "otp_verification.html.gohtml"
		default:
			emailData.Template = "mail.html.gohtml"
		}
	}

	// Create mail message
	msg := mailer.Message{
		From:     emailData.From,
		FromName: emailData.FromName,
		To:       emailData.To,
		Subject:  emailData.Subject,
		Template: emailData.Template,
		Data:     emailData.Variables,
	}

	// Send email
	err := consumer.mailer.SendSMTPMessage(msg)
	if err != nil {
		consumer.logger.Printf("Error sending email: %v", err)
	} else {
		consumer.logger.Printf("Email sent successfully to %s (type: %s)", emailData.To, emailData.Type)
	}
}

// handleOtpGeneratedEvent processes auth.otp_generated events
func (consumer *Consumer) handleOtpGeneratedEvent(event StandardEvent) {
	// Extract OTP data from event
	jsonData, err := json.Marshal(event.Data)
	if err != nil {
		consumer.logger.Printf("Error marshaling OTP data: %v", err)
		return
	}

	var otpData struct {
		Email      string `json:"email"`
		OTP        string `json:"otp"`
		Purpose    string `json:"purpose"`
		ExpiresIn  int    `json:"expires_in"`
		Message    string `json:"message"`
		ActionText string `json:"action_text"`
	}

	if err := json.Unmarshal(jsonData, &otpData); err != nil {
		consumer.logger.Printf("Error unmarshaling OTP data: %v", err)
		return
	}

	// Prepare variables for template
	variables := map[string]string{
		"otp":         otpData.OTP,
		"purpose":     otpData.Purpose,
		"expires_in":  fmt.Sprintf("%d", otpData.ExpiresIn),
		"message":     otpData.Message,
		"action_text": otpData.ActionText,
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
		From:     "", // Will use default
		FromName: "",
		To:       otpData.Email,
		Subject:  subject,
		Template: "otp_verification.html.gohtml",
		Data:     variables,
	}

	// Send email
	err = consumer.mailer.SendSMTPMessage(msg)
	if err != nil {
		consumer.logger.Printf("Error sending OTP email: %v", err)
	} else {
		consumer.logger.Printf("OTP email sent successfully to %s (purpose: %s)", otpData.Email, otpData.Purpose)
	}
}

// handleMessage processes the message based on its type (legacy format)
func (consumer *Consumer) handleMessage(payload Payload) {
	switch payload.Name {
	case "mail.registration":
		consumer.handleRegistrationEmail(payload)
	case "mail.reset_password":
		consumer.handlePasswordResetEmail(payload)
	case "mail.password_change":
		consumer.handlePasswordChangeEmail(payload)
	case "mail.order_confirmation":
		consumer.handleOrderConfirmationEmail(payload)
	default:
		consumer.logger.Printf("Unknown email type: %s", payload.Name)
	}
}

// handleRegistrationEmail processes registration confirmation emails
func (consumer *Consumer) handleRegistrationEmail(payload Payload) {
	var emailData EmailPayload
	err := json.Unmarshal([]byte(payload.Data), &emailData)
	if err != nil {
		consumer.logger.Printf("Error unmarshalling registration email data: %v", err)
		return
	}

	msg := mailer.Message{
		From:     emailData.From,
		To:       emailData.To,
		Subject:  emailData.Subject,
		Template: "register.html.gohtml",
		Data:     emailData.Data,
	}

	err = consumer.mailer.SendSMTPMessage(msg)
	if err != nil {
		consumer.logger.Printf("Error sending registration email: %v", err)
	}
}

// handlePasswordResetEmail processes password reset emails
func (consumer *Consumer) handlePasswordResetEmail(payload Payload) {
	var emailData EmailPayload
	err := json.Unmarshal([]byte(payload.Data), &emailData)
	if err != nil {
		consumer.logger.Printf("Error unmarshalling password reset email data: %v", err)
		return
	}

	msg := mailer.Message{
		From:     emailData.From,
		To:       emailData.To,
		Subject:  emailData.Subject,
		Template: "reset_password.html.gohtml",
		Data:     emailData.Data,
	}

	err = consumer.mailer.SendSMTPMessage(msg)
	if err != nil {
		consumer.logger.Printf("Error sending password reset email: %v", err)
	}
}

// handlePasswordChangeEmail processes password change notification emails
func (consumer *Consumer) handlePasswordChangeEmail(payload Payload) {
	var emailData EmailPayload
	err := json.Unmarshal([]byte(payload.Data), &emailData)
	if err != nil {
		consumer.logger.Printf("Error unmarshalling password change email data: %v", err)
		return
	}

	msg := mailer.Message{
		From:     emailData.From,
		To:       emailData.To,
		Subject:  emailData.Subject,
		Template: "password_change.html.gohtml",
		Data:     emailData.Data,
	}

	err = consumer.mailer.SendSMTPMessage(msg)
	if err != nil {
		consumer.logger.Printf("Error sending password change email: %v", err)
	}
}

// handleOrderConfirmationEmail processes order confirmation emails
func (consumer *Consumer) handleOrderConfirmationEmail(payload Payload) {
	var emailData EmailPayload
	err := json.Unmarshal([]byte(payload.Data), &emailData)
	if err != nil {
		consumer.logger.Printf("Error unmarshalling order confirmation email data: %v", err)
		return
	}

	msg := mailer.Message{
		From:     emailData.From,
		To:       emailData.To,
		Subject:  emailData.Subject,
		Template: "order_confirmation.html.gohtml",
		Data:     emailData.Data,
	}

	err = consumer.mailer.SendSMTPMessage(msg)
	if err != nil {
		consumer.logger.Printf("Error sending order confirmation email: %v", err)
	}
}

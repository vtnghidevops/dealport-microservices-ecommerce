package event

import (
	"encoding/json"
	"fmt"
)

// SchemaValidator provides methods to validate event data against schemas
type SchemaValidator struct {
	validators map[string]func([]byte) error
}

// NewSchemaValidator creates a new validator with standard event schemas
func NewSchemaValidator() *SchemaValidator {
	v := &SchemaValidator{
		validators: make(map[string]func([]byte) error),
	}

	// Register standard validators
	v.RegisterValidator(EventTypeUserRegistered, validateUserRegistered)
	v.RegisterValidator(EventTypeUserPasswordResetRequest, validatePasswordResetRequested)
	v.RegisterValidator(EventTypeUserPasswordChanged, validatePasswordChanged)
	v.RegisterValidator(EventTypeEmailSend, validateEmailSend)

	return v
}

// RegisterValidator registers a custom validator for an event type
func (v *SchemaValidator) RegisterValidator(eventType string, validator func([]byte) error) {
	v.validators[eventType] = validator
}

// Validate checks if an event's data conforms to its schema
func (v *SchemaValidator) Validate(event Event) error {
	// Check if we have a validator for this event type
	validator, exists := v.validators[event.Name]
	if !exists {
		return fmt.Errorf("no validator found for event type: %s", event.Name)
	}

	// Convert data to JSON for validation
	var data []byte
	var err error

	switch d := event.Data.(type) {
	case string:
		data = []byte(d)
	case []byte:
		data = d
	default:
		data, err = json.Marshal(d)
		if err != nil {
			return fmt.Errorf("failed to marshal event data: %w", err)
		}
	}

	// Validate the data
	if err := validator(data); err != nil {
		return fmt.Errorf("event data validation failed: %w", err)
	}

	return nil
}

// Standard validators for each event type

func validateUserRegistered(data []byte) error {
	var userData EventUserRegistered
	if err := json.Unmarshal(data, &userData); err != nil {
		return fmt.Errorf("invalid user.registered data format: %w", err)
	}

	// Validate required fields
	if userData.Email == "" {
		return fmt.Errorf("email is required")
	}

	return nil
}

func validatePasswordResetRequested(data []byte) error {
	var resetData EventPasswordResetRequested
	if err := json.Unmarshal(data, &resetData); err != nil {
		return fmt.Errorf("invalid user.password_reset_requested data format: %w", err)
	}

	// Validate required fields
	if resetData.Email == "" {
		return fmt.Errorf("email is required")
	}
	if resetData.TokenHash == "" {
		return fmt.Errorf("token_hash is required")
	}
	if resetData.ExpiresAt == "" {
		return fmt.Errorf("expires_at is required")
	}

	return nil
}

func validatePasswordChanged(data []byte) error {
	var changedData EventPasswordChanged
	if err := json.Unmarshal(data, &changedData); err != nil {
		return fmt.Errorf("invalid user.password_changed data format: %w", err)
	}

	// Validate required fields
	if changedData.Email == "" {
		return fmt.Errorf("email is required")
	}

	return nil
}

func validateEmailSend(data []byte) error {
	var emailData EventEmailSend
	if err := json.Unmarshal(data, &emailData); err != nil {
		return fmt.Errorf("invalid email.send data format: %w", err)
	}

	// Validate required fields
	if emailData.To == "" {
		return fmt.Errorf("to is required")
	}
	if emailData.Subject == "" {
		return fmt.Errorf("subject is required")
	}
	if emailData.PlainText == "" && emailData.HTMLContent == "" && emailData.Template == "" {
		return fmt.Errorf("at least one of plain_text, html_content, or template is required")
	}

	return nil
}

package event

import (
	"encoding/json"
	"fmt"
	"time"
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

// EventUserRegistered represents the data structure for user.registered events
type EventUserRegistered struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// EventPasswordResetRequested represents the data structure for user.password_reset_requested events
type EventPasswordResetRequested struct {
	Email     string `json:"email"`
	TokenHash string `json:"token_hash"`
	ExpiresAt string `json:"expires_at"`
}

// EventOTPRequested represents the data structure for user.otp_requested events
type EventOTPRequested struct {
	Email     string    `json:"email"`
	OTPCode   string    `json:"otp_code"`
	Purpose   string    `json:"purpose"` // registration, password_reset, etc.
	ExpiresAt time.Time `json:"expires_at"`
}

// EventOTPVerified represents the data structure for user.otp_verified events
type EventOTPVerified struct {
	Email      string `json:"email"`
	Purpose    string `json:"purpose"` // registration, password_reset, etc.
	VerifiedAt string `json:"verified_at"`
}

// EventPasswordChanged represents the data structure for user.password_changed events
type EventPasswordChanged struct {
	Email     string `json:"email"`
	ChangedAt string `json:"changed_at"`
}

// EventEmailSend represents the data structure for email.send events
type EventEmailSend struct {
	From        string            `json:"from"`
	To          string            `json:"to"`
	Subject     string            `json:"subject"`
	PlainText   string            `json:"plain_text"`
	HTMLContent string            `json:"html_content,omitempty"`
	Template    string            `json:"template,omitempty"`
	Variables   map[string]string `json:"variables,omitempty"`
}

// NewEvent creates a new event with proper defaults
func NewEvent(name string, data interface{}, source string) Event {
	now := time.Now()
	return Event{
		ID:          generateEventID(),
		Name:        name,
		Data:        data,
		DataSchema:  "v1",
		Source:      source,
		CreatedAt:   now,
		PublishedAt: now,
		Version:     "v1",
	}
}

// Validate checks if the event has all required fields
func (e Event) Validate() error {
	if e.ID == "" {
		return fmt.Errorf("event ID cannot be empty")
	}
	if e.Name == "" {
		return fmt.Errorf("event name cannot be empty")
	}
	if e.Source == "" {
		return fmt.Errorf("event source cannot be empty")
	}
	return nil
}

// AsJSON converts an event to JSON format
func (e Event) AsJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON parses a JSON string into an Event struct
func EventFromJSON(data []byte) (Event, error) {
	var e Event
	if err := json.Unmarshal(data, &e); err != nil {
		return Event{}, err
	}
	return e, nil
}

// Helper function to generate a unique event ID
// In production, use a proper UUID library
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

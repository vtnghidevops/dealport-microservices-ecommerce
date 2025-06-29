package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"user-service/internal/domain"
	"user-service/internal/service"
)

// UserRegisteredData represents the data for user.registered events
type UserRegisteredData struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
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

// Handler handles events related to user service
type Handler struct {
	userService service.UserService
	logger      *log.Logger
}

// NewEventHandler creates a new event handler
func NewEventHandler(userService service.UserService, logger *log.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}

// HandleEvent handles various events from other services
func (h *Handler) HandleEvent(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the event from request body
	var event StandardEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.logger.Printf("Error decoding event: %v", err)
		http.Error(w, "Invalid event data", http.StatusBadRequest)
		return
	}

	h.logger.Printf("Received event: %s from %s", event.Name, event.Source)

	// Process the event based on its name
	var err error
	switch event.Name {
	case "user.registered":
		err = h.handleUserRegistered(event)
	case "user.password_changed":
		err = h.handlePasswordChanged(event)
	default:
		h.logger.Printf("Unknown event type: %s", event.Name)
		http.Error(w, "Unknown event type", http.StatusBadRequest)
		return
	}

	if err != nil {
		h.logger.Printf("Error handling event: %v", err)
		http.Error(w, fmt.Sprintf("Error processing event: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"processed"}`))
}

// handleUserRegistered processes user.registered events
func (h *Handler) handleUserRegistered(event StandardEvent) error {
	// Extract data
	jsonData, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	var userData UserRegisteredData
	if err := json.Unmarshal(jsonData, &userData); err != nil {
		return fmt.Errorf("failed to parse user data: %w", err)
	}

	// Check if user already exists (to handle potential duplicate events)
	existingUser, err := h.userService.GetUserByID(context.Background(), userData.ID)
	if err == nil && existingUser != nil {
		h.logger.Printf("User %s already exists, skipping creation", userData.ID)
		return nil
	}

	// Create user profile in user service
	createUserReq := &domain.CreateUserRequest{
		Email:     userData.Email,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Username:  userData.Username,
		Role:      "user", // Default role
		// No password as auth is handled by auth service
		DisplayName: fmt.Sprintf("%s %s", userData.FirstName, userData.LastName),
	}

	// Create the user
	user, err := h.userService.CreateUser(context.Background(), createUserReq)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	h.logger.Printf("User %s created successfully", user.ID)
	return nil
}

// handlePasswordChanged processes user.password_changed events
func (h *Handler) handlePasswordChanged(event StandardEvent) error {
	// Extract data
	jsonData, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal password data: %w", err)
	}

	var passwordData map[string]string
	if err := json.Unmarshal(jsonData, &passwordData); err != nil {
		return fmt.Errorf("failed to parse password data: %w", err)
	}

	email := passwordData["email"]

	// Update user's password changed timestamp
	if email == "" {
		return fmt.Errorf("email is required in password changed event")
	}

	// Get user by email
	user, err := h.userService.GetUserByEmail(context.Background(), email)
	if err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}

	// Create update request
	updateUserReq := &domain.UpdateUserRequest{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Role:      user.Role,
		Status:    user.Status,
		Active:    user.Active,
	}

	// Update user
	_, err = h.userService.UpdateUser(context.Background(), updateUserReq)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	h.logger.Printf("Password changed timestamp updated for user %s", user.ID)
	return nil
}

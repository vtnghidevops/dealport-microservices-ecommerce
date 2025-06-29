package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"broker-service/internal/event"
)

// RequestPayload is the format of the JSON request
type RequestPayload struct {
	Action string       `json:"action"`
	Auth   AuthPayload  `json:"auth,omitempty"`
	User   UserPayload  `json:"user,omitempty"`
	Email  EmailPayload `json:"email,omitempty"`
	OTP    OTPPayload   `json:"otp,omitempty"`
}

// AuthPayload is the auth information embedded in RequestPayload
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserPayload contains user data for registration or other user operations
type UserPayload struct {
	ID        string `json:"id,omitempty"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username,omitempty"`
}

// EmailPayload contains email sending data
type EmailPayload struct {
	From        string            `json:"from"`
	To          string            `json:"to"`
	Subject     string            `json:"subject"`
	PlainText   string            `json:"plain_text,omitempty"`
	HTMLContent string            `json:"html_content,omitempty"`
	Template    string            `json:"template,omitempty"`
	Variables   map[string]string `json:"variables,omitempty"`
}

// OTPPayload contains OTP data for requesting and verifying OTPs
type OTPPayload struct {
	Email   string `json:"email"`
	OTPCode string `json:"otp_code,omitempty"`
	Purpose string `json:"purpose"` // registration, password_reset, etc.
}

// Broker handles all incoming requests and routes them appropriately
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// HandleSubmission processes incoming requests based on the action specified
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Initialize the validator
	validator := event.NewSchemaValidator()

	// Create a context with timeout for event operations
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "register":
		app.registerUser(w, ctx, requestPayload.User, validator)
	case "send_email":
		app.sendEmail(w, ctx, requestPayload.Email)
	case "reset_password":
		app.requestPasswordReset(w, ctx, requestPayload.Auth.Email)
	case "request_otp":
		app.requestOTP(w, ctx, requestPayload.OTP)
	case "verify_otp":
		app.verifyOTP(w, ctx, requestPayload.OTP)
	default:
		app.errorJSON(w, errors.New("unknown action"), http.StatusBadRequest)
	}
}

// registerUser handles user registration and publishes a user.registered event
func (app *Config) registerUser(w http.ResponseWriter, ctx context.Context, user UserPayload, validator *event.SchemaValidator) {
	// Create user registration event data
	userData := event.EventUserRegistered{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	// Create the event
	userEvent, err := event.CreateUserRegisteredEvent(ctx, userData, "broker-service")
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Validate the event data
	if err := validator.Validate(userEvent); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Publish the event
	err = app.eventEmitter.PublishEvent(ctx, userEvent, event.EventTypeUserRegistered)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "User registration event published",
		Data:    map[string]string{"event_id": userEvent.ID},
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// sendEmail handles email sending and publishes an email.send event
func (app *Config) sendEmail(w http.ResponseWriter, ctx context.Context, email EmailPayload) {
	// Create email event data
	emailData := event.EventEmailSend{
		From:        email.From,
		To:          email.To,
		Subject:     email.Subject,
		PlainText:   email.PlainText,
		HTMLContent: email.HTMLContent,
		Template:    email.Template,
		Variables:   email.Variables,
	}

	// Create the event
	emailEvent := event.NewEvent(event.EventTypeEmailSend, emailData, "broker-service")

	// Publish the event
	err := app.eventEmitter.PublishEvent(ctx, emailEvent, event.EventTypeEmailSend)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Email send event published",
		Data:    map[string]string{"event_id": emailEvent.ID},
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// requestPasswordReset handles password reset requests and publishes a user.password_reset_requested event
func (app *Config) requestPasswordReset(w http.ResponseWriter, ctx context.Context, email string) {
	// Generate a token hash (in a real app, this would be a secure token)
	tokenHash := "example-token-hash-" + email

	// Set expiry time for the reset token
	expiresAt := time.Now().Add(24 * time.Hour)

	// Create the password reset event
	resetEvent, err := event.CreatePasswordResetRequestedEvent(ctx, email, tokenHash, expiresAt, "broker-service")
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Publish the event
	err = app.eventEmitter.PublishEvent(ctx, resetEvent, event.EventTypeUserPasswordResetRequest)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Password reset requested",
		Data:    map[string]string{"event_id": resetEvent.ID},
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// authenticate is a placeholder for authentication logic
func (app *Config) authenticate(w http.ResponseWriter, auth AuthPayload) {
	// Placeholder for authentication logic
	// In a real application, you would validate credentials and return a token

	resp := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    map[string]string{"authenticated": "true"},
	}

	app.writeJSON(w, http.StatusOK, resp)
}

// requestOTP handles OTP generation requests and publishes a user.otp_requested event
func (app *Config) requestOTP(w http.ResponseWriter, ctx context.Context, otpPayload OTPPayload) {
	// Validate required fields
	if otpPayload.Email == "" {
		app.errorJSON(w, errors.New("email is required"), http.StatusBadRequest)
		return
	}
	if otpPayload.Purpose == "" {
		app.errorJSON(w, errors.New("purpose is required"), http.StatusBadRequest)
		return
	}

	// Generate a random 6-digit OTP code (in a real app, use a proper random generator)
	otpCode := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)

	// Set expiry time for the OTP
	expiresAt := time.Now().Add(15 * time.Minute)

	// Create the OTP requested event
	otpEvent, err := event.CreateOTPRequestedEvent(ctx, otpPayload.Email, otpCode, otpPayload.Purpose, expiresAt, "broker-service")
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Publish the event
	err = app.eventEmitter.PublishEvent(ctx, otpEvent, event.EventTypeUserOTPRequested)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "OTP requested",
		Data:    map[string]string{"event_id": otpEvent.ID},
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// verifyOTP handles OTP verification and publishes a user.otp_verified event
func (app *Config) verifyOTP(w http.ResponseWriter, ctx context.Context, otpPayload OTPPayload) {
	// Validate required fields
	if otpPayload.Email == "" {
		app.errorJSON(w, errors.New("email is required"), http.StatusBadRequest)
		return
	}
	if otpPayload.OTPCode == "" {
		app.errorJSON(w, errors.New("OTP code is required"), http.StatusBadRequest)
		return
	}
	if otpPayload.Purpose == "" {
		app.errorJSON(w, errors.New("purpose is required"), http.StatusBadRequest)
		return
	}

	// In a real application, you would validate the OTP against the stored value
	// Here we just publish an event for the authentication service to handle the validation

	// Create the OTP verified event
	otpEvent, err := event.CreateOTPVerifiedEvent(ctx, otpPayload.Email, otpPayload.Purpose, "broker-service")
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Publish the event
	err = app.eventEmitter.PublishEvent(ctx, otpEvent, event.EventTypeUserOTPVerified)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "OTP verified",
		Data:    map[string]string{"event_id": otpEvent.ID},
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// Helper methods for JSON handling

// readJSON reads JSON from request body into data
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

// writeJSON writes JSON to response
func (app *Config) writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	return err
}

// errorJSON writes an error to the response
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := jsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	app.writeJSON(w, statusCode, payload)
}

// jsonResponse is the standard response format
type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

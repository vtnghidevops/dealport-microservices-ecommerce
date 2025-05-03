package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// MailClient is a client for the mail service
type MailClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewMailClient creates a new mail service client
func NewMailClient(baseURL string) *MailClient {
	return &MailClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendOTPEmail sends an OTP email
func (c *MailClient) SendOTPEmail(email, otp, action, message string, expiresIn int) error {
	// Create request payload
	payload := map[string]interface{}{
		"email":      email,
		"otp":        otp,
		"action":     action,
		"message":    message,
		"expires_in": expiresIn,
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal OTP email payload: %w", err)
	}

	// Create request
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/send-otp", c.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to create request to mail service: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to mail service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusAccepted {
		var errorResp struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return fmt.Errorf("mail service returned status %d", resp.StatusCode)
		}

		return fmt.Errorf("mail service error: %s", errorResp.Message)
	}

	return nil
}

// SendRegistrationOTP sends a registration verification OTP
func (c *MailClient) SendRegistrationOTP(email, otp string, expiresIn int) error {
	message := "We received a request to create an account with this email address."
	return c.SendOTPEmail(
		email,
		otp,
		"account registration",
		message,
		expiresIn,
	)
}

// SendPasswordResetOTP sends a password reset OTP
func (c *MailClient) SendPasswordResetOTP(email, otp string, expiresIn int) error {
	message := "We received a request to reset the password for your account."
	return c.SendOTPEmail(
		email,
		otp,
		"password reset",
		message,
		expiresIn,
	)
}

// SendGenericOTP sends an OTP for a generic purpose
func (c *MailClient) SendGenericOTP(email, otp, purpose string, expiresIn int) error {
	message := fmt.Sprintf("We received a request that requires verification for your account (%s).", purpose)
	return c.SendOTPEmail(
		email,
		otp,
		purpose,
		message,
		expiresIn,
	)
}

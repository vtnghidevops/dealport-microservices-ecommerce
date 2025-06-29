package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Constants for event types
const (
	EventTypeUserRegistered           = "user.registered"
	EventTypeUserPasswordResetRequest = "user.password_reset_requested"
	EventTypeUserPasswordChanged      = "user.password_changed"
	EventTypeEmailSend                = "email.send"
	EventTypeUserOTPRequested         = "user.otp_requested"
	EventTypeUserOTPVerified          = "user.otp_verified"

	// Order events
	EventOrderCreated  = "order.created"
	EventOrderUpdated  = "order.updated"
	EventOrderCanceled = "order.canceled"

	// Cart events
	EventCartCreated   = "cart.created"
	EventCartUpdated   = "cart.updated"
	EventCartAbandoned = "cart.abandoned"

	// Product events
	EventProductCreated = "product.created"
	EventProductUpdated = "product.updated"
	EventProductDeleted = "product.deleted"

	// Inventory events
	EventInventoryUpdated = "inventory.updated"
	EventLowStockAlert    = "inventory.low_stock"
)

// CreateUserRegisteredEvent creates a standardized user.registered event
func CreateUserRegisteredEvent(ctx context.Context, userData EventUserRegistered, source string) (Event, error) {
	return NewEvent(EventTypeUserRegistered, userData, source), nil
}

// CreatePasswordResetRequestedEvent creates a standardized user.password_reset_requested event
func CreatePasswordResetRequestedEvent(ctx context.Context, email, tokenHash string, expiresAt time.Time, source string) (Event, error) {
	data := EventPasswordResetRequested{
		Email:     email,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}
	return NewEvent(EventTypeUserPasswordResetRequest, data, source), nil
}

// CreatePasswordChangedEvent creates a standardized user.password_changed event
func CreatePasswordChangedEvent(ctx context.Context, email string, source string) (Event, error) {
	data := EventPasswordChanged{
		Email:     email,
		ChangedAt: time.Now().Format(time.RFC3339),
	}
	return NewEvent(EventTypeUserPasswordChanged, data, source), nil
}

// CreateEmailSendEvent creates a standardized email.send event
func CreateEmailSendEvent(ctx context.Context, from, to, subject, plainText, htmlContent string, source string) (Event, error) {
	data := EventEmailSend{
		From:        from,
		To:          to,
		Subject:     subject,
		PlainText:   plainText,
		HTMLContent: htmlContent,
	}
	return NewEvent(EventTypeEmailSend, data, source), nil
}

// CreateTemplatedEmailSendEvent creates a standardized email.send event with a template
func CreateTemplatedEmailSendEvent(ctx context.Context, from, to, subject, template string, variables map[string]string, source string) (Event, error) {
	data := EventEmailSend{
		From:      from,
		To:        to,
		Subject:   subject,
		Template:  template,
		Variables: variables,
	}
	return NewEvent(EventTypeEmailSend, data, source), nil
}

// CreateOTPRequestedEvent creates a standardized user.otp_requested event
func CreateOTPRequestedEvent(ctx context.Context, email, otpCode, purpose string, expiresAt time.Time, source string) (Event, error) {
	data := EventOTPRequested{
		Email:     email,
		OTPCode:   otpCode,
		Purpose:   purpose,
		ExpiresAt: expiresAt,
	}
	return NewEvent(EventTypeUserOTPRequested, data, source), nil
}

// CreateOTPVerifiedEvent creates a standardized user.otp_verified event
func CreateOTPVerifiedEvent(ctx context.Context, email, purpose string, source string) (Event, error) {
	data := EventOTPVerified{
		Email:      email,
		Purpose:    purpose,
		VerifiedAt: time.Now().Format(time.RFC3339),
	}
	return NewEvent(EventTypeUserOTPVerified, data, source), nil
}

// ExtractUserRegisteredData extracts the user registered data from an event
func ExtractUserRegisteredData(event Event) (EventUserRegistered, error) {
	var userData EventUserRegistered

	// Handle different data formats
	switch data := event.Data.(type) {
	case EventUserRegistered:
		return data, nil
	case map[string]interface{}:
		// Convert the map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			return userData, fmt.Errorf("failed to marshal data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &userData); err != nil {
			return userData, fmt.Errorf("failed to unmarshal data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &userData); err != nil {
			return userData, fmt.Errorf("failed to unmarshal data string: %w", err)
		}
	case []byte:
		if err := json.Unmarshal(data, &userData); err != nil {
			return userData, fmt.Errorf("failed to unmarshal data bytes: %w", err)
		}
	default:
		return userData, fmt.Errorf("unsupported data type")
	}

	return userData, nil
}

// ExtractEmailSendData extracts the email send data from an event
func ExtractEmailSendData(event Event) (EventEmailSend, error) {
	var emailData EventEmailSend

	// Handle different data formats
	switch data := event.Data.(type) {
	case EventEmailSend:
		return data, nil
	case map[string]interface{}:
		// Convert the map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			return emailData, fmt.Errorf("failed to marshal data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &emailData); err != nil {
			return emailData, fmt.Errorf("failed to unmarshal data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &emailData); err != nil {
			return emailData, fmt.Errorf("failed to unmarshal data string: %w", err)
		}
	case []byte:
		if err := json.Unmarshal(data, &emailData); err != nil {
			return emailData, fmt.Errorf("failed to unmarshal data bytes: %w", err)
		}
	default:
		return emailData, fmt.Errorf("unsupported data type")
	}

	return emailData, nil
}

// ExtractOTPRequestedData extracts the OTP requested data from an event
func ExtractOTPRequestedData(event Event) (EventOTPRequested, error) {
	var otpData EventOTPRequested

	// Handle different data formats
	switch data := event.Data.(type) {
	case EventOTPRequested:
		return data, nil
	case map[string]interface{}:
		// Convert the map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			return otpData, fmt.Errorf("failed to marshal data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &otpData); err != nil {
			return otpData, fmt.Errorf("failed to unmarshal data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &otpData); err != nil {
			return otpData, fmt.Errorf("failed to unmarshal data string: %w", err)
		}
	case []byte:
		if err := json.Unmarshal(data, &otpData); err != nil {
			return otpData, fmt.Errorf("failed to unmarshal data bytes: %w", err)
		}
	default:
		return otpData, fmt.Errorf("unsupported data type")
	}

	return otpData, nil
}

// ExtractOTPVerifiedData extracts the OTP verified data from an event
func ExtractOTPVerifiedData(event Event) (EventOTPVerified, error) {
	var otpData EventOTPVerified

	// Handle different data formats
	switch data := event.Data.(type) {
	case EventOTPVerified:
		return data, nil
	case map[string]interface{}:
		// Convert the map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			return otpData, fmt.Errorf("failed to marshal data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &otpData); err != nil {
			return otpData, fmt.Errorf("failed to unmarshal data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &otpData); err != nil {
			return otpData, fmt.Errorf("failed to unmarshal data string: %w", err)
		}
	case []byte:
		if err := json.Unmarshal(data, &otpData); err != nil {
			return otpData, fmt.Errorf("failed to unmarshal data bytes: %w", err)
		}
	default:
		return otpData, fmt.Errorf("unsupported data type")
	}

	return otpData, nil
}

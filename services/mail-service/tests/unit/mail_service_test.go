package unit_test

import (
	"context"
	"log"
	"mail-service/internal/config"
	"mail-service/internal/grpc"
	"mail-service/internal/mailer"
	"mail-service/proto"
	"mail-service/tests/mocks"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setupTestServer creates a new test server with a mock mailer
func setupTestServer() (*grpc.MailServer, *mocks.MockMailer) {
	// Create a mock mailer
	mockMailer := new(mocks.MockMailer)

	// Create mail struct
	mailStruct := mailer.Mail{
		Domain:      "test.com",
		Host:        "localhost",
		Port:        25,
		Username:    "test",
		Password:    "password",
		Encryption:  "none",
		FromName:    "Test Sender",
		FromAddress: "test@example.com",
	}

	// Create config with the mail struct
	cfg := config.Config{
		Mailer: mailStruct,
	}
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)

	server := grpc.NewMailServer(cfg, logger)

	return server, mockMailer
}

// TestNewMailServer tests the creation of a new mail server
func TestNewMailServer(t *testing.T) {
	// Create mail struct
	mailStruct := mailer.Mail{
		Domain:      "test.com",
		Host:        "localhost",
		Port:        25,
		Username:    "test",
		Password:    "password",
		Encryption:  "none",
		FromName:    "Test Sender",
		FromAddress: "test@example.com",
	}

	// Create config with the mail struct
	cfg := config.Config{
		Mailer: mailStruct,
	}
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)

	// Create the server
	server := grpc.NewMailServer(cfg, logger)

	// Assertions
	assert.NotNil(t, server)
	assert.Equal(t, cfg, server.Config)
}

// TestSendEmail tests the SendEmail method
func TestSendEmail(t *testing.T) {
	// Create mail struct
	mailStruct := mailer.Mail{
		Domain:      "test.com",
		Host:        "localhost",
		Port:        25,
		Username:    "test",
		Password:    "password",
		Encryption:  "none",
		FromName:    "Test Sender",
		FromAddress: "test@example.com",
	}

	// Create config with the mail struct
	cfg := config.Config{
		Mailer: mailStruct,
	}
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)

	server := grpc.NewMailServer(cfg, logger)
	ctx := context.Background()

	t.Run("Basic email request", func(t *testing.T) {
		// Create request
		req := &proto.SendEmailRequest{
			From:     "test@example.com",
			FromName: "Test Sender",
			To:       "recipient@example.com",
			Subject:  "Test Subject",
			Message:  "Test Message",
		}

		// Test the response
		resp, err := server.SendEmail(ctx, req)

		// Basic assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("Email with missing fields", func(t *testing.T) {
		// Create request with missing fields
		req := &proto.SendEmailRequest{
			To:      "recipient@example.com",
			Subject: "Test Subject",
			Message: "Test Message",
		}

		// Test the response
		resp, err := server.SendEmail(ctx, req)

		// Basic assertions - should still process without error
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

// TestSendTemplateEmail tests the SendTemplateEmail method
func TestSendTemplateEmail(t *testing.T) {
	server, _ := setupTestServer()
	ctx := context.Background()

	t.Run("Registration template", func(t *testing.T) {
		// Create request for registration email
		req := &proto.SendTemplateEmailRequest{
			To:      "user@example.com",
			Subject: "Welcome to our platform",
			Variables: map[string]string{
				"email_type": "registration",
				"first_name": "John",
				"last_name":  "Doe",
			},
		}

		// Test response
		resp, err := server.SendTemplateEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("Password reset template", func(t *testing.T) {
		// Create request for password reset
		req := &proto.SendTemplateEmailRequest{
			To:      "user@example.com",
			Subject: "Reset Your Password",
			Variables: map[string]string{
				"email_type": "password_reset",
				"token_hash": "abc123def456",
				"expires_at": "2023-12-31T23:59:59Z",
			},
		}

		// Test response
		resp, err := server.SendTemplateEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("Password change template", func(t *testing.T) {
		// Create request for password change confirmation
		req := &proto.SendTemplateEmailRequest{
			To:      "user@example.com",
			Subject: "Password Changed",
			Variables: map[string]string{
				"email_type": "password_change",
				"changed_at": "2023-12-31T23:59:59Z",
			},
		}

		// Test response
		resp, err := server.SendTemplateEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("Order confirmation template", func(t *testing.T) {
		// Create request for order confirmation
		req := &proto.SendTemplateEmailRequest{
			To:      "customer@example.com",
			Subject: "Order Confirmation",
			Variables: map[string]string{
				"email_type":   "order_confirmation",
				"order_id":     "ord123",
				"order_number": "ORD-123",
				"total":        "99.99",
				"status":       "pending",
			},
		}

		// Test response
		resp, err := server.SendTemplateEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("Order processing template", func(t *testing.T) {
		// Create request for order processing
		req := &proto.SendTemplateEmailRequest{
			To:      "customer@example.com",
			Subject: "Order Processing",
			Variables: map[string]string{
				"email_type":   "order_processing",
				"order_id":     "ord123",
				"order_number": "ORD-123",
			},
		}

		// Test response
		resp, err := server.SendTemplateEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("Order shipped template", func(t *testing.T) {
		// Create request for order shipped
		req := &proto.SendTemplateEmailRequest{
			To:      "customer@example.com",
			Subject: "Order Shipped",
			Variables: map[string]string{
				"email_type":      "order_shipped",
				"order_id":        "ord123",
				"order_number":    "ORD-123",
				"tracking_code":   "TRK123456789",
				"shipping_method": "Express",
			},
		}

		// Test response
		resp, err := server.SendTemplateEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("Order delivered template", func(t *testing.T) {
		// Create request for order delivered
		req := &proto.SendTemplateEmailRequest{
			To:      "customer@example.com",
			Subject: "Order Delivered",
			Variables: map[string]string{
				"email_type":    "order_delivered",
				"order_id":      "ord123",
				"order_number":  "ORD-123",
				"delivery_date": "2023-12-31",
			},
		}

		// Test response
		resp, err := server.SendTemplateEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("Order cancelled template", func(t *testing.T) {
		// Create request for order cancelled
		req := &proto.SendTemplateEmailRequest{
			To:      "customer@example.com",
			Subject: "Order Cancelled",
			Variables: map[string]string{
				"email_type":   "order_cancelled",
				"order_id":     "ord123",
				"order_number": "ORD-123",
				"reason":       "Out of stock",
			},
		}

		// Test response
		resp, err := server.SendTemplateEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("Order status change template", func(t *testing.T) {
		// Create request for order status update
		req := &proto.SendTemplateEmailRequest{
			To:      "customer@example.com",
			Subject: "Order Status Updated",
			Variables: map[string]string{
				"email_type":      "order_status_change",
				"order_id":        "ord123",
				"order_number":    "ORD-123",
				"status":          "ready for pickup",
				"previous_status": "processing",
			},
		}

		// Test response
		resp, err := server.SendTemplateEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("Payment success template", func(t *testing.T) {
		// Create request for payment success
		req := &proto.SendTemplateEmailRequest{
			To:      "customer@example.com",
			Subject: "Payment Successful",
			Variables: map[string]string{
				"email_type":     "payment_success",
				"payment_id":     "pay123",
				"order_id":       "ord123",
				"order_number":   "ORD-123",
				"amount":         "99.99",
				"payment_method": "Credit Card",
			},
		}

		// Test response
		resp, err := server.SendTemplateEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("Custom template", func(t *testing.T) {
		// Create request with custom template
		req := &proto.SendTemplateEmailRequest{
			From:     "noreply@example.com",
			FromName: "Customer Service",
			To:       "customer@example.com",
			Subject:  "Important Information",
			Template: "custom.html.gohtml",
			Variables: map[string]string{
				"name":    "Customer",
				"content": "This is a custom message",
			},
		}

		// Test response
		resp, err := server.SendTemplateEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})
}

// TestSendOTPEmail tests the SendOTPEmail method
func TestSendOTPEmail(t *testing.T) {
	server, _ := setupTestServer()
	ctx := context.Background()

	t.Run("Registration OTP", func(t *testing.T) {
		// Create request for registration OTP
		req := &proto.SendOTPEmailRequest{
			Email:      "user@example.com",
			Otp:        "123456",
			Purpose:    "registration",
			ExpiresIn:  15,
			Message:    "Please use this code to verify your account",
			ActionText: "Verify Account",
		}

		// Test response
		resp, err := server.SendOTPEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("Password reset OTP", func(t *testing.T) {
		// Create request for password reset OTP
		req := &proto.SendOTPEmailRequest{
			Email:      "user@example.com",
			Otp:        "654321",
			Purpose:    "password_reset",
			ExpiresIn:  10,
			Message:    "Please use this code to reset your password",
			ActionText: "Reset Password",
		}

		// Test response
		resp, err := server.SendOTPEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("OTP with no expiration", func(t *testing.T) {
		// Create request with no expiration (should use default)
		req := &proto.SendOTPEmailRequest{
			Email:      "user@example.com",
			Otp:        "987654",
			Purpose:    "login",
			ExpiresIn:  0, // Should use default expiration
			Message:    "Use this code to log in",
			ActionText: "Log In",
		}

		// Test response
		resp, err := server.SendOTPEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("OTP with minimal fields", func(t *testing.T) {
		// Create request with just the required fields
		req := &proto.SendOTPEmailRequest{
			Email: "user@example.com",
			Otp:   "123456",
		}

		// Test response
		resp, err := server.SendOTPEmail(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Message)
	})
}

// TestCreateMessage tests creating and validating a message
func TestCreateMessage(t *testing.T) {
	t.Run("Basic message", func(t *testing.T) {
		// Create a test message
		msg := mailer.Message{
			From:     "test@example.com",
			FromName: "Test Sender",
			To:       "recipient@example.com",
			Subject:  "Test Subject",
			Data: map[string]string{
				"message": "Test Message",
			},
		}

		// Verify the message is constructed correctly
		assert.Equal(t, "test@example.com", msg.From)
		assert.Equal(t, "recipient@example.com", msg.To)
		assert.Equal(t, "Test Subject", msg.Subject)
		assert.NotNil(t, msg.Data)
		assert.Equal(t, "Test Message", msg.Data.(map[string]string)["message"])
	})

	t.Run("Message with template", func(t *testing.T) {
		// Create a test message with template
		msg := mailer.Message{
			From:     "test@example.com",
			FromName: "Test Sender",
			To:       "recipient@example.com",
			Subject:  "Test Subject",
			Template: "test.html.gohtml",
			Data: map[string]interface{}{
				"name": "Test User",
				"url":  "https://example.com",
			},
		}

		// Verify the message is constructed correctly
		assert.Equal(t, "test@example.com", msg.From)
		assert.Equal(t, "recipient@example.com", msg.To)
		assert.Equal(t, "Test Subject", msg.Subject)
		assert.Equal(t, "test.html.gohtml", msg.Template)
		assert.NotNil(t, msg.Data)
		data := msg.Data.(map[string]interface{})
		assert.Equal(t, "Test User", data["name"])
		assert.Equal(t, "https://example.com", data["url"])
	})
}

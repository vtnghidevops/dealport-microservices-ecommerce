package unit_test

import (
	"mail-service/internal/mailer"
	"mail-service/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestSendSMTPMessage tests sending emails through the mailer
func TestSendSMTPMessage(t *testing.T) {
	// Create a mock mailer
	mockMailer := new(mocks.MockMailer)

	// Test data for a generic email
	data := map[string]string{
		"from":      "sender@example.com",
		"from_name": "Sender",
		"to":        "recipient@example.com",
		"subject":   "Test Subject",
		"message":   "Test Message",
		"template":  "mail.html.gohtml",
	}

	// Create the message for expectations
	msg := mailer.Message{
		From:     data["from"],
		FromName: data["from_name"],
		To:       data["to"],
		Subject:  data["subject"],
		Template: data["template"],
		Data:     data,
	}

	// Setup mock expectations
	mockMailer.On("SendSMTPMessage", mock.MatchedBy(func(m mailer.Message) bool {
		return m.From == msg.From &&
			m.FromName == msg.FromName &&
			m.To == msg.To &&
			m.Subject == msg.Subject &&
			m.Template == msg.Template
	})).Return(nil).Once()

	// Test sending the email
	err := mockMailer.SendSMTPMessage(msg)

	// Assertions
	assert.NoError(t, err)
	mockMailer.AssertExpectations(t)
}

// TestRegistrationEmail tests registration email template
func TestRegistrationEmail(t *testing.T) {
	// Create a mock mailer
	mockMailer := new(mocks.MockMailer)

	// Test data
	email := "user@example.com"
	data := map[string]string{
		"first_name": "John",
		"last_name":  "Doe",
		"subject":    "Welcome to our platform",
	}

	// Setup mock expectations for a registration email template
	mockMailer.On("SendSMTPMessage", mock.MatchedBy(func(msg mailer.Message) bool {
		return msg.To == email &&
			msg.Subject == data["subject"] &&
			msg.Template == "register.html.gohtml"
	})).Return(nil).Once()

	// Create a message for the registration template
	msg := mailer.Message{
		To:       email,
		Subject:  data["subject"],
		Template: "register.html.gohtml",
		Data:     data,
	}

	// Send the message
	err := mockMailer.SendSMTPMessage(msg)

	// Assertions
	assert.NoError(t, err)
	mockMailer.AssertExpectations(t)
}

// TestOrderConfirmationEmail tests order confirmation email template
func TestOrderConfirmationEmail(t *testing.T) {
	// Create a mock mailer
	mockMailer := new(mocks.MockMailer)

	// Test data
	email := "customer@example.com"
	data := map[string]string{
		"order_id":     "ord123",
		"order_number": "ORD-123",
		"total":        "99.99",
		"status":       "pending",
	}

	// Setup mock expectations
	mockMailer.On("SendSMTPMessage", mock.MatchedBy(func(msg mailer.Message) bool {
		return msg.To == email &&
			msg.Template == "order_confirmation.html.gohtml" &&
			msg.Subject == "Order Confirmation - #"+data["order_number"]
	})).Return(nil).Once()

	// Create a message for order confirmation
	msg := mailer.Message{
		To:       email,
		Template: "order_confirmation.html.gohtml",
		Subject:  "Order Confirmation - #" + data["order_number"],
		Data:     data,
	}

	// Send the message
	err := mockMailer.SendSMTPMessage(msg)

	// Assertions
	assert.NoError(t, err)
	mockMailer.AssertExpectations(t)
}

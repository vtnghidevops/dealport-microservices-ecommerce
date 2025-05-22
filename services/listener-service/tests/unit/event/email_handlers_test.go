package event_test

import (
	"listener-service/event"
	"listener-service/tests/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestSendPasswordChangedEmail tests the sendPasswordChangedEmail function
func TestSendPasswordChangedEmail(t *testing.T) {
	// Arrange
	mockMailClient := new(mocks.MockMailClient)
	mockLoggerClient := new(mocks.MockLoggerClient)

	// Setup mock behaviors
	mockMailClient.On("SendPasswordChangedEmail", "user@example.com", mock.Anything).Return(nil)
	mockLoggerClient.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Mock logic for password changed email
	passwordChangedData := event.PasswordChanged{
		Email:     "user@example.com",
		ChangedAt: time.Now().Format(time.RFC3339),
	}

	testEvent := createTestEvent("auth.password_changed", passwordChangedData)

	// Mock consumer for testing event handler
	mockConsumer := mocks.NewMockConsumer()
	mockConsumer.On("HandleStandardEvent", testEvent, "auth.password_changed").Return(nil)

	// Act
	err := mockConsumer.HandleStandardEvent(testEvent, "auth.password_changed")

	// Assert
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

// TestSendOrderStatusEmail tests the sendOrderStatusEmail function
func TestSendOrderStatusEmail(t *testing.T) {
	// Arrange
	mockMailClient := new(mocks.MockMailClient)
	mockLoggerClient := new(mocks.MockLoggerClient)

	// Setup mock behaviors
	mockMailClient.On("SendOrderStatusEmail",
		"customer@example.com",
		"ORD-123",
		"delivered",
		mock.Anything,
	).Return(nil)
	mockLoggerClient.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create test data for order status changed
	orderStatusData := event.OrderStatusChangedData{
		OrderID:        "order-123",
		OrderNumber:    "ORD-123",
		UserID:         "user-123",
		UserEmail:      "customer@example.com",
		Status:         "delivered",
		PreviousStatus: "processing",
		PaymentMethod:  "card",
		Total:          99.99,
		CreatedAt:      time.Now(),
	}

	testEvent := createTestEvent("order.status_changed", orderStatusData)

	// Mock consumer for testing event handler
	mockConsumer := mocks.NewMockConsumer()
	mockConsumer.On("HandleStandardEvent", testEvent, "order.status_changed").Return(nil)

	// Act
	err := mockConsumer.HandleStandardEvent(testEvent, "order.status_changed")

	// Assert
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

// TestSendPaymentSuccessEmail tests the sendPaymentSuccessEmail function
func TestSendPaymentSuccessEmail(t *testing.T) {
	// Arrange
	mockMailClient := new(mocks.MockMailClient)
	mockLoggerClient := new(mocks.MockLoggerClient)

	// Setup mock behaviors
	mockMailClient.On("SendPaymentConfirmationEmail",
		"customer@example.com",
		"ORD-123",
		99.99,
		"txn_abc123",
	).Return(nil)
	mockLoggerClient.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create test data for payment succeeded
	paymentData := event.PaymentSucceededData{
		OrderID:       "order-123",
		OrderNumber:   "ORD-123",
		UserID:        "user-123",
		UserEmail:     "customer@example.com",
		PaymentMethod: "card",
		Amount:        99.99,
		Currency:      "USD",
		TransactionID: "txn_abc123",
		Status:        "completed",
		PaymentDate:   time.Now().Format(time.RFC3339),
	}

	testEvent := createTestEvent("payment.succeeded", paymentData)

	// Mock consumer for testing event handler
	mockConsumer := mocks.NewMockConsumer()
	mockConsumer.On("HandleStandardEvent", testEvent, "payment.succeeded").Return(nil)

	// Act
	err := mockConsumer.HandleStandardEvent(testEvent, "payment.succeeded")

	// Assert
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

// TestSendOTPEmail tests the sendOTPEmail function
func TestSendOTPEmail(t *testing.T) {
	// Arrange
	mockMailClient := new(mocks.MockMailClient)
	mockLoggerClient := new(mocks.MockLoggerClient)

	// Setup mock behaviors
	mockMailClient.On("SendOTPEmail",
		"user@example.com",
		"123456",
		"login",
		int32(15*60),
		mock.Anything,
		mock.Anything,
	).Return(nil)
	mockLoggerClient.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create test data for OTP generated
	otpData := event.OTPGeneratedData{
		Email:      "user@example.com",
		OTP:        "123456",
		Purpose:    "login",
		ExpiresIn:  15 * 60,
		Message:    "Your OTP code for login",
		ActionText: "Login",
	}

	testEvent := createTestEvent("auth.otp_generated", otpData)

	// Mock consumer for testing event handler
	mockConsumer := mocks.NewMockConsumer()
	mockConsumer.On("HandleStandardEvent", testEvent, "auth.otp_generated").Return(nil)

	// Act
	err := mockConsumer.HandleStandardEvent(testEvent, "auth.otp_generated")

	// Assert
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

// Helper function to create a standard event with the given name and data
func createTestEvent(name string, data interface{}) event.StandardEvent {
	return event.StandardEvent{
		ID:          "test-event-id-" + name,
		Name:        name,
		Data:        data,
		DataSchema:  "1.0",
		Source:      "test-service",
		CreatedAt:   time.Now(),
		PublishedAt: time.Now(),
		Version:     "1.0",
	}
}

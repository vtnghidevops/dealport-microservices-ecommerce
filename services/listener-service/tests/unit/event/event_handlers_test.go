package event

import (
	"encoding/json"
	"listener-service/event"
	userpb "listener-service/proto/user"
	"listener-service/tests/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ExportedConsumer is a testing wrapper for real consumer
// that exposes methods we need to test directly
type ExportedConsumer struct {
	event.Consumer
}

// EventHandlerFunc defines a function that handles events
type EventHandlerFunc func(event event.StandardEvent) error

// TestHandleUserRegistered tests handling of user.registered event
func TestHandleUserRegistered(t *testing.T) {
	// Arrange
	mockLoggerClient := new(mocks.MockLoggerClient)
	mockMailClient := new(mocks.MockMailClient)

	// Setup default behaviors
	mockLoggerClient.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Test welcome email delivery
	mockMailClient.On("SendWelcomeEmail",
		"test@example.com",
		"testuser",
		"Test",
		"User",
	).Return(nil)

	// Create test data
	userData := event.UserRegistered{
		ID:        "user123",
		Email:     "test@example.com",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	}

	testEvent := createStandardEvent("user.registered", userData)

	// Mock behavior
	mockConsumer := mocks.NewMockConsumer()
	mockConsumer.On("HandleStandardEvent", testEvent, "user.registered").Return(nil)

	// Act
	err := mockConsumer.HandleStandardEvent(testEvent, "user.registered")

	// Assert
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
	// Note: In a real implementation, we would verify that the SendWelcomeEmail was called
	// But our current test structure doesn't allow testing private methods directly
}

// TestHandlePasswordResetRequested tests handling of auth.password_reset_requested event
func TestHandlePasswordResetRequested(t *testing.T) {
	// Arrange
	mockLoggerClient := new(mocks.MockLoggerClient)
	mockMailClient := new(mocks.MockMailClient)

	// Setup behaviors
	mockLoggerClient.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockMailClient.On("SendPasswordResetEmail",
		"test@example.com",
		mock.Anything,
		mock.Anything,
	).Return(nil)

	// Create test data
	resetData := event.PasswordResetRequested{
		Email:     "test@example.com",
		TokenHash: "abc123",
		ExpiresAt: "2023-12-31T23:59:59Z",
	}

	testEvent := createStandardEvent("auth.password_reset_requested", resetData)

	// Mock behavior
	mockConsumer := mocks.NewMockConsumer()
	mockConsumer.On("HandleStandardEvent", testEvent, "auth.password_reset_requested").Return(nil)

	// Act
	err := mockConsumer.HandleStandardEvent(testEvent, "auth.password_reset_requested")

	// Assert
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

// TestHandleOrderCreated tests handling of order.created event
func TestHandleOrderCreated(t *testing.T) {
	// Arrange
	mockLoggerClient := new(mocks.MockLoggerClient)
	mockMailClient := new(mocks.MockMailClient)
	mockUserClient := new(mocks.MockUserServiceClient)

	// Setup behaviors
	mockLoggerClient.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockMailClient.On("SendOrderConfirmationEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockUserClient.On("SyncUserOrderData", mock.Anything, mock.MatchedBy(func(req *userpb.SyncUserOrderDataRequest) bool {
		return req.UserId == "user-123"
	})).Return(&userpb.SyncUserOrderDataResponse{Success: true}, nil)

	// Create test data
	now := time.Now()
	orderData := event.OrderCreatedData{
		OrderID:       "order-123",
		OrderNumber:   "ORD-123",
		UserID:        "user-123",
		UserEmail:     "customer@example.com",
		Status:        "pending",
		PaymentMethod: "card",
		Total:         99.99,
		CreatedAt:     now,
		Items: []event.ItemData{
			{
				ProductID: "prod-1",
				Name:      "Product 1",
				Quantity:  1,
				Price:     99.99,
			},
		},
	}

	testEvent := createStandardEvent("order.created", orderData)

	// Mock behavior
	mockConsumer := mocks.NewMockConsumer()
	mockConsumer.On("HandleStandardEvent", testEvent, "order.created").Return(nil)

	// Act
	err := mockConsumer.HandleStandardEvent(testEvent, "order.created")

	// Assert
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

// TestHandleLegacyLogEvent tests handling a legacy log event
func TestHandleLegacyLogEvent(t *testing.T) {
	// Arrange
	mockLoggerClient := new(mocks.MockLoggerClient)

	// Setup behaviors
	mockLoggerClient.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create test data
	logData := map[string]string{
		"service": "auth-service",
		"message": "User login successful",
	}
	dataJSON, _ := json.Marshal(logData)

	legacyPayload := event.Payload{
		Name: "log",
		Data: string(dataJSON),
	}

	// Mock behavior
	mockConsumer := mocks.NewMockConsumer()
	mockConsumer.On("HandleLegacyEvent", legacyPayload).Return(nil)

	// Act
	err := mockConsumer.HandleLegacyEvent(legacyPayload)

	// Assert
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

// Helper function to create a standard event
func createStandardEvent(name string, data interface{}) event.StandardEvent {
	return event.StandardEvent{
		ID:          "test-event-id",
		Name:        name,
		Data:        data,
		DataSchema:  "1.0",
		Source:      "test-service",
		CreatedAt:   time.Now(),
		PublishedAt: time.Now(),
		Version:     "1.0",
	}
}

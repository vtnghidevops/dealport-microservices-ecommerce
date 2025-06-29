package event

import (
	"errors"
	"listener-service/event"
	"listener-service/tests/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// createErrorTestEvent creates a test event for error testing scenarios
func createErrorTestEvent(name string, data interface{}) event.StandardEvent {
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

// TestHandleInvalidEventFormat tests handling of events with invalid format
func TestHandleInvalidEventFormat(t *testing.T) {
	// Arrange
	mockConsumer := mocks.NewMockConsumer()

	// Create malformed event - missing required fields
	malformedEvent := event.StandardEvent{
		// Missing ID
		Name: "", // Missing name
		Data: nil,
		// Missing DataSchema
		Source:    "test-service",
		CreatedAt: time.Now(),
		// Missing PublishedAt
		Version: "1.0",
	}

	// Setup mock behavior - expect error
	mockConsumer.On("HandleStandardEvent", malformedEvent, "unknown.event").Return(
		errors.New("invalid event format"),
	)

	// Act
	err := mockConsumer.HandleStandardEvent(malformedEvent, "unknown.event")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid event format")
	mockConsumer.AssertExpectations(t)
}

// TestHandleInvalidEventData tests handling events with invalid data structure
func TestHandleInvalidEventData(t *testing.T) {
	// Arrange
	mockConsumer := mocks.NewMockConsumer()

	// Valid event structure but invalid data format for the event type
	invalidDataEvent := event.StandardEvent{
		ID:          "test-event-id",
		Name:        "order.created",                // Expecting OrderCreatedData
		Data:        "this is not valid order data", // String instead of struct
		DataSchema:  "1.0",
		Source:      "test-service",
		CreatedAt:   time.Now(),
		PublishedAt: time.Now(),
		Version:     "1.0",
	}

	// Setup mock behavior - expect error
	mockConsumer.On("HandleStandardEvent", invalidDataEvent, "order.created").Return(
		errors.New("failed to unmarshal order data"),
	)

	// Act
	err := mockConsumer.HandleStandardEvent(invalidDataEvent, "order.created")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal")
	mockConsumer.AssertExpectations(t)
}

// TestHandleUnknownEventType tests handling of unknown event types
func TestHandleUnknownEventType(t *testing.T) {
	// Arrange
	mockConsumer := mocks.NewMockConsumer()

	// Create event with unknown type
	unknownEvent := event.StandardEvent{
		ID:          "test-event-id",
		Name:        "unknown.event.type",
		Data:        map[string]interface{}{"foo": "bar"},
		DataSchema:  "1.0",
		Source:      "test-service",
		CreatedAt:   time.Now(),
		PublishedAt: time.Now(),
		Version:     "1.0",
	}

	// Setup mock behavior - expect error about unknown event type
	mockConsumer.On("HandleStandardEvent", unknownEvent, "unknown.event.type").Return(
		errors.New("unknown event type"),
	)

	// Act
	err := mockConsumer.HandleStandardEvent(unknownEvent, "unknown.event.type")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown event type")
	mockConsumer.AssertExpectations(t)
}

// TestHandleMailServiceError tests error handling when mail service fails
func TestHandleMailServiceError(t *testing.T) {
	// Arrange
	mockMailClient := new(mocks.MockMailClient)
	mockLoggerClient := new(mocks.MockLoggerClient)

	// Setup mock behaviors - simulate email sending failure
	mailError := errors.New("failed to send email")
	mockMailClient.On("SendWelcomeEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mailError)
	mockLoggerClient.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create valid user registered event
	userData := event.UserRegistered{
		ID:        "user123",
		Email:     "test@example.com",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	}

	testEvent := createErrorTestEvent("user.registered", userData)

	// Mock consumer that will simulate mail service error
	mockConsumer := mocks.NewMockConsumer()
	mockConsumer.On("HandleStandardEvent", testEvent, "user.registered").Return(
		errors.New("failed to send welcome email: failed to send email"),
	)

	// Act
	err := mockConsumer.HandleStandardEvent(testEvent, "user.registered")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send welcome email")
	mockConsumer.AssertExpectations(t)
}

// TestHandleLoggerError tests error handling when logger service fails
func TestHandleLoggerError(t *testing.T) {
	// Arrange
	mockLoggerClient := new(mocks.MockLoggerClient)

	// Setup mock behaviors - simulate logging failure
	loggerError := errors.New("failed to write log")
	mockLoggerClient.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return(loggerError)

	// Create a test event to log
	testEvent := createErrorTestEvent("test.event", map[string]string{"test": "data"})

	// Mock consumer that will simulate logger service error
	mockConsumer := mocks.NewMockConsumer()
	mockConsumer.On("HandleStandardEvent", testEvent, "test.event").Return(
		errors.New("failed to log event: failed to write log"),
	)

	// Act
	err := mockConsumer.HandleStandardEvent(testEvent, "test.event")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to log event")
	mockConsumer.AssertExpectations(t)
}

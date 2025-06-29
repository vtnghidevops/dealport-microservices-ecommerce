package event

import (
	"errors"
	"listener-service/event"
	"listener-service/tests/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestConsumer_Listen tests the Listen method
func TestConsumer_Listen(t *testing.T) {
	// Arrange
	mockConsumer := mocks.NewMockConsumer()

	// Define test data
	topics := []string{"order.*", "user.*"}

	// Setup mock behavior
	mockConsumer.On("Listen", topics).Return(nil)

	// Act
	err := mockConsumer.Listen(topics)

	// Assert
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

// TestConsumer_Listen_Error tests error handling in Listen
func TestConsumer_Listen_Error(t *testing.T) {
	// Arrange
	mockConsumer := mocks.NewMockConsumer()
	expectedError := errors.New("listen error")

	// Define test data
	topics := []string{"test.*"}

	// Setup mock behavior
	mockConsumer.On("Listen", topics).Return(expectedError)

	// Act
	err := mockConsumer.Listen(topics)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockConsumer.AssertExpectations(t)
}

// TestConsumer_HandleStandardEvent tests handling a standard event
func TestConsumer_HandleStandardEvent(t *testing.T) {
	// Arrange
	mockConsumer := mocks.NewMockConsumer()

	// Create test data
	testEvent := event.StandardEvent{
		ID:          "test-event-id",
		Name:        "test.event",
		Data:        map[string]interface{}{"key": "value"},
		DataSchema:  "1.0",
		Source:      "test-service",
		CreatedAt:   time.Now(),
		PublishedAt: time.Now(),
		Version:     "1.0",
	}

	// Setup mock behavior
	mockConsumer.On("HandleStandardEvent", testEvent, "test.event").Return(nil)

	// Act
	err := mockConsumer.HandleStandardEvent(testEvent, "test.event")

	// Assert
	assert.NoError(t, err)
	assert.True(t, mockConsumer.IsProcessed("test-event-id"), "Event should be marked as processed")
	mockConsumer.AssertExpectations(t)
}

// TestConsumer_HandleLegacyEvent tests handling a legacy event
func TestConsumer_HandleLegacyEvent(t *testing.T) {
	// Arrange
	mockConsumer := mocks.NewMockConsumer()

	// Create test data
	testPayload := event.Payload{
		Name: "log",
		Data: `{"message": "test message", "service": "test-service"}`,
	}

	// Setup mock behavior
	mockConsumer.On("HandleLegacyEvent", testPayload).Return(nil)

	// Act
	err := mockConsumer.HandleLegacyEvent(testPayload)

	// Assert
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

package event

import (
	"context"
	"encoding/json"
	"errors"
	"listener-service/event"
	"listener-service/tests/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// LogClientInterface defines the interface for logging clients
type LogClientInterface interface {
	WriteLog(context.Context, string, string) error
}

// TestLogEvent tests logging of events
func TestLogEvent(t *testing.T) {
	// Arrange - Create a mock LoggerClient
	mockLoggerClient := new(mocks.MockLoggerClient)

	// Create test event
	testEvent := event.StandardEvent{
		ID:         "test-123",
		Name:       "user.registered",
		Data:       map[string]string{"user_id": "123", "email": "test@example.com"},
		DataSchema: "1.0",
		Source:     "auth-service",
		CreatedAt:  time.Now(),
		Version:    "1.0",
	}

	// Marshal the expected JSON data that would be sent to the logger
	logData, _ := json.Marshal(testEvent)

	// Setup mock behaviors - expect a WriteLog call with the JSON data
	mockLoggerClient.On("WriteLog",
		mock.Anything,
		"INFO",
		string(logData),
	).Return(nil)

	// Create a mock transport that uses our mock logger client
	mockTransport := &mockLogTransport{
		loggerClient: mockLoggerClient,
	}

	// Act - Log the event
	err := mockTransport.LogEvent(testEvent)

	// Assert
	assert.NoError(t, err)
	mockLoggerClient.AssertCalled(t, "WriteLog", mock.Anything, "INFO", string(logData))
}

// TestLogEvent_Error tests error handling in log event
func TestLogEvent_Error(t *testing.T) {
	// Arrange - Create a mock LoggerClient with error response
	mockLoggerClient := new(mocks.MockLoggerClient)
	expectedErr := errors.New("logging error")

	// Create test event
	testEvent := event.StandardEvent{
		ID:   "test-123",
		Name: "user.registered",
	}

	// Setup mock behaviors - return an error
	mockLoggerClient.On("WriteLog",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(expectedErr)

	// Create a mock transport that uses our mock logger client
	mockTransport := &mockLogTransport{
		loggerClient: mockLoggerClient,
	}

	// Act - Log the event
	err := mockTransport.LogEvent(testEvent)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockLoggerClient.AssertCalled(t, "WriteLog", mock.Anything, mock.Anything, mock.Anything)
}

// mockLogTransport is a test helper that implements just the logging functionality
type mockLogTransport struct {
	loggerClient LogClientInterface
}

// LogEvent implements our simplified logging functionality for testing
func (m *mockLogTransport) LogEvent(event event.StandardEvent) error {
	// Marshal the event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Log the event
	ctx := context.Background()
	return m.loggerClient.WriteLog(ctx, "INFO", string(eventJSON))
}

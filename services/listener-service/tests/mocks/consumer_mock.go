package mocks

import (
	"listener-service/event"

	"github.com/stretchr/testify/mock"
)

// ConsumerInterface defines an interface for the Consumer
// so we can properly mock it without reflection
type ConsumerInterface interface {
	Listen(topics []string) error
	HandleStandardEvent(event event.StandardEvent, routingKey string) error
	HandleLegacyEvent(payload event.Payload) error
}

// MockConsumer implements the ConsumerInterface for testing
type MockConsumer struct {
	mock.Mock

	// For convenience in testing
	ProcessedMessages map[string]bool
}

// NewMockConsumer creates a new MockConsumer
func NewMockConsumer() *MockConsumer {
	return &MockConsumer{
		ProcessedMessages: make(map[string]bool),
	}
}

// Listen mocks the Listen method
func (m *MockConsumer) Listen(topics []string) error {
	args := m.Called(topics)
	return args.Error(0)
}

// HandleStandardEvent mocks the HandleStandardEvent method
func (m *MockConsumer) HandleStandardEvent(e event.StandardEvent, routingKey string) error {
	args := m.Called(e, routingKey)

	// For convenience in testing, track processed message IDs
	if e.ID != "" {
		m.ProcessedMessages[e.ID] = true
	}

	return args.Error(0)
}

// HandleLegacyEvent mocks the HandleLegacyEvent method
func (m *MockConsumer) HandleLegacyEvent(payload event.Payload) error {
	args := m.Called(payload)
	return args.Error(0)
}

// IsProcessed checks if a message has been processed
func (m *MockConsumer) IsProcessed(messageID string) bool {
	return m.ProcessedMessages[messageID]
}

package helpers

import (
	"listener-service/event"
	"sync"

	"github.com/stretchr/testify/mock"
)

// MockConsumer is a mock implementation of the event Consumer interface
type MockConsumer struct {
	mock.Mock
	processedEvents map[string]bool
	mu              sync.Mutex
}

// NewMockConsumer creates a new mock consumer for testing
func NewMockConsumer() *MockConsumer {
	return &MockConsumer{
		processedEvents: make(map[string]bool),
	}
}

// Listen mocks the Listen method
func (m *MockConsumer) Listen(topics []string) error {
	args := m.Called(topics)
	return args.Error(0)
}

// HandleStandardEvent mocks the HandleStandardEvent method
func (m *MockConsumer) HandleStandardEvent(e event.StandardEvent, topic string) error {
	args := m.Called(e, topic)

	// Mark event as processed regardless of error
	m.mu.Lock()
	defer m.mu.Unlock()
	m.processedEvents[e.ID] = true

	return args.Error(0)
}

// HandleLegacyEvent mocks the HandleLegacyEvent method
func (m *MockConsumer) HandleLegacyEvent(payload event.Payload) error {
	args := m.Called(payload)
	return args.Error(0)
}

// IsProcessed checks if an event has been processed
func (m *MockConsumer) IsProcessed(eventID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.processedEvents[eventID]
}

// MockEventProcessor is a more detailed mock for testing specific event handlers
type MockEventProcessor struct {
	mock.Mock
	HandlerFunc func(event event.StandardEvent) error
}

// ProcessEvent mocks the processing of a specific event
func (m *MockEventProcessor) ProcessEvent(e event.StandardEvent) error {
	args := m.Called(e)

	// If handler func is set, call it
	if m.HandlerFunc != nil {
		return m.HandlerFunc(e)
	}

	return args.Error(0)
}

// SetHandler sets a custom handler function for an event
func (m *MockEventProcessor) SetHandler(handlerFunc func(event event.StandardEvent) error) {
	m.HandlerFunc = handlerFunc
}

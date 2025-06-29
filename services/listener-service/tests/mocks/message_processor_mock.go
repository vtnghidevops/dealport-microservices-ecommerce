package mocks

import (
	"sync"
	"time"

	"github.com/stretchr/testify/mock"
)

// MessageProcessorInterface defines an interface for the MessageProcessor
// so we can properly mock it without reflection
type MessageProcessorInterface interface {
	IsProcessed(messageID string) bool
	MarkProcessed(messageID string)
}

// MockMessageProcessor is a mock implementation of MessageProcessorInterface
type MockMessageProcessor struct {
	mock.Mock
	processedMessages map[string]time.Time
	mutex             sync.RWMutex
}

// NewMockMessageProcessor creates a new MockMessageProcessor
func NewMockMessageProcessor() *MockMessageProcessor {
	return &MockMessageProcessor{
		processedMessages: make(map[string]time.Time),
	}
}

// IsProcessed mocks the IsProcessed method
func (m *MockMessageProcessor) IsProcessed(messageID string) bool {
	args := m.Called(messageID)
	return args.Bool(0)
}

// MarkProcessed mocks the MarkProcessed method
func (m *MockMessageProcessor) MarkProcessed(messageID string) {
	m.Called(messageID)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.processedMessages[messageID] = time.Now()
}

// GetProcessedMessageIDs returns the IDs of all processed messages
func (m *MockMessageProcessor) GetProcessedMessageIDs() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var ids []string
	for id := range m.processedMessages {
		ids = append(ids, id)
	}
	return ids
}

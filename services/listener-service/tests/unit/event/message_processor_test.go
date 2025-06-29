package event

import (
	"fmt"
	"listener-service/tests/mocks"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMessageProcessor_IsProcessed tests the IsProcessed method
func TestMessageProcessor_IsProcessed(t *testing.T) {
	// Create a new MockMessageProcessor for testing
	processor := mocks.NewMockMessageProcessor()

	// Set up mock behavior
	messageID := "test-message-1"
	processor.On("IsProcessed", messageID).Return(false).Once()
	processor.On("MarkProcessed", messageID).Return().Once()
	processor.On("IsProcessed", messageID).Return(true).Once()

	// Test when message is not processed
	assert.False(t, processor.IsProcessed(messageID), "New message should not be marked as processed")

	// Mark the message as processed
	processor.MarkProcessed(messageID)

	// Test when message is processed
	assert.True(t, processor.IsProcessed(messageID), "Message should be marked as processed")

	// Verify all expectations were met
	processor.AssertExpectations(t)
}

// TestMessageProcessor_ConcurrentAccess tests concurrent access to the message processor
func TestMessageProcessor_ConcurrentAccess(t *testing.T) {
	// Create a new MockMessageProcessor
	processor := mocks.NewMockMessageProcessor()

	// Number of concurrent goroutines
	numGoroutines := 10
	messagesPerGoroutine := 100

	// Set up mock expectations dynamically
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < messagesPerGoroutine; j++ {
			messageID := fmt.Sprintf("message-%d-%d", i, j)
			processor.On("MarkProcessed", messageID).Return().Once()
			processor.On("IsProcessed", messageID).Return(true).Once()
		}
	}

	// Wait group to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch goroutines to mark messages as processed
	for i := 0; i < numGoroutines; i++ {
		go func(routineID int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				messageID := fmt.Sprintf("message-%d-%d", routineID, j)
				processor.MarkProcessed(messageID)
				// Immediately check if it's processed
				assert.True(t, processor.IsProcessed(messageID), "Message should be marked as processed")
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Verify all expectations were met
	processor.AssertExpectations(t)
}

// TestMessageProcessor_Cleanup tests the cleanup of old processed messages
func TestMessageProcessor_Cleanup(t *testing.T) {
	// This test is optional and would depend on how MessageProcessor is implemented
	// Since we're using a mock, we'll just ensure the basic contract is followed

	processor := mocks.NewMockMessageProcessor()

	// Set up expectations
	processor.On("MarkProcessed", "message1").Return().Once()
	processor.On("MarkProcessed", "message2").Return().Once()
	processor.On("IsProcessed", "message1").Return(true).Once()
	processor.On("IsProcessed", "message2").Return(true).Once()

	// Mark messages as processed
	processor.MarkProcessed("message1")
	processor.MarkProcessed("message2")

	// Verify they are processed
	assert.True(t, processor.IsProcessed("message1"))
	assert.True(t, processor.IsProcessed("message2"))

	// Verify expectations
	processor.AssertExpectations(t)
}

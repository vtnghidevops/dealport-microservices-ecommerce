package mocks

import (
	"checkout-service/internal/event"

	"github.com/stretchr/testify/mock"
)

// MockEventEmitter mocks the event emitter for integration tests
type MockEventEmitter struct {
	mock.Mock
}

// EmitOrderCreated mocks emitting order created event
func (m *MockEventEmitter) EmitOrderCreated(data event.OrderEventData) error {
	args := m.Called(data)
	return args.Error(0)
}

// EmitPaymentSucceeded mocks emitting payment succeeded event
func (m *MockEventEmitter) EmitPaymentSucceeded(data event.PaymentEventData) error {
	args := m.Called(data)
	return args.Error(0)
}

// EmitPaymentFailed mocks emitting payment failed event
func (m *MockEventEmitter) EmitPaymentFailed(data event.PaymentEventData) error {
	args := m.Called(data)
	return args.Error(0)
}

// EmitOrderStatusChanged mocks emitting order status changed event
func (m *MockEventEmitter) EmitOrderStatusChanged(data event.OrderEventData) error {
	args := m.Called(data)
	return args.Error(0)
}

// GetEmittedEvents returns a list of emitted events for verification
func (m *MockEventEmitter) GetEmittedEvents() []struct {
	EventName string
	Data      interface{}
} {
	// This is just a helper method for tests, not part of the actual interface
	// We can add implementation in specific tests if needed
	return nil
}

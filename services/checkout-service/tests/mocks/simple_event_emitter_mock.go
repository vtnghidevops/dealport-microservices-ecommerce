package mocks

import (
	"checkout-service/internal/event"
)

// SimpleEventEmitter mocks the event Emitter for testing
type SimpleEventEmitter struct {
	// Track method calls
	EmitOrderCreatedCalled       bool
	EmitOrderStatusChangedCalled bool
	EmitPaymentSucceededCalled   bool
	EmitPaymentFailedCalled      bool

	// Store emitted events for verification
	LastOrderCreatedEvent       event.OrderEventData
	LastOrderStatusChangedEvent event.OrderEventData
	LastPaymentSucceededEvent   event.PaymentEventData
	LastPaymentFailedEvent      event.PaymentEventData

	// For injecting errors
	EmitOrderCreatedError       error
	EmitOrderStatusChangedError error
	EmitPaymentSucceededError   error
	EmitPaymentFailedError      error
}

// NewSimpleEventEmitter creates a new mock event emitter
func NewSimpleEventEmitter() *SimpleEventEmitter {
	return &SimpleEventEmitter{}
}

// EmitOrderCreated mocks emitting an order created event
func (m *SimpleEventEmitter) EmitOrderCreated(data event.OrderEventData) error {
	m.EmitOrderCreatedCalled = true
	m.LastOrderCreatedEvent = data
	return m.EmitOrderCreatedError
}

// EmitOrderStatusChanged mocks emitting an order status changed event
func (m *SimpleEventEmitter) EmitOrderStatusChanged(data event.OrderEventData) error {
	m.EmitOrderStatusChangedCalled = true
	m.LastOrderStatusChangedEvent = data
	return m.EmitOrderStatusChangedError
}

// EmitPaymentSucceeded mocks emitting a payment succeeded event
func (m *SimpleEventEmitter) EmitPaymentSucceeded(data event.PaymentEventData) error {
	m.EmitPaymentSucceededCalled = true
	m.LastPaymentSucceededEvent = data
	return m.EmitPaymentSucceededError
}

// EmitPaymentFailed mocks emitting a payment failed event
func (m *SimpleEventEmitter) EmitPaymentFailed(data event.PaymentEventData) error {
	m.EmitPaymentFailedCalled = true
	m.LastPaymentFailedEvent = data
	return m.EmitPaymentFailedError
}

// EmitPaymentCompleted is for compatibility with existing tests
// It maps to EmitPaymentSucceeded
func (m *SimpleEventEmitter) EmitPaymentCompleted(data event.PaymentEventData) error {
	return m.EmitPaymentSucceeded(data)
}

package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockEventEmitter is a mock implementation of EventEmitter
type MockEventEmitter struct {
	mock.Mock
}

// EmitUserActivityEvent mocks the EmitUserActivityEvent method
func (m *MockEventEmitter) EmitUserActivityEvent(ctx context.Context, action, userID, message string, metadata map[string]interface{}) error {
	args := m.Called(ctx, action, userID, message, metadata)
	return args.Error(0)
}

// EmitProfileUpdatedEvent mocks the EmitProfileUpdatedEvent method
func (m *MockEventEmitter) EmitProfileUpdatedEvent(ctx context.Context, userID string, updatedFields map[string]interface{}) error {
	args := m.Called(ctx, userID, updatedFields)
	return args.Error(0)
}

// EmitLogoutEvent mocks the EmitLogoutEvent method
func (m *MockEventEmitter) EmitLogoutEvent(ctx context.Context, userID, email string) error {
	args := m.Called(ctx, userID, email)
	return args.Error(0)
}

// Close mocks the Close method
func (m *MockEventEmitter) Close() error {
	args := m.Called()
	return args.Error(0)
}

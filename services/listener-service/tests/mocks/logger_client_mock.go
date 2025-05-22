package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockLoggerClient is a mock for the LoggerClient
type MockLoggerClient struct {
	mock.Mock
}

// WriteLog mocks the WriteLog method
func (m *MockLoggerClient) WriteLog(ctx context.Context, name, data string) error {
	args := m.Called(ctx, name, data)
	return args.Error(0)
}

// LogUserActivity mocks the LogUserActivity method
func (m *MockLoggerClient) LogUserActivity(ctx context.Context, action, userID, message string, metadata map[string]interface{}) error {
	args := m.Called(ctx, action, userID, message, metadata)
	return args.Error(0)
}

// Close mocks the Close method
func (m *MockLoggerClient) Close() error {
	args := m.Called()
	return args.Error(0)
} 
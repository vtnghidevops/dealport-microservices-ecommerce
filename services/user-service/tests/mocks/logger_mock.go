package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockLoggerClient is a mock of LoggerClient
type MockLoggerClient struct {
	mock.Mock
}

// LogUserActivity mocks the LogUserActivity method
func (m *MockLoggerClient) LogUserActivity(ctx context.Context, action, userID, message string, metadata map[string]interface{}) error {
	args := m.Called(ctx, action, userID, message, metadata)
	return args.Error(0)
}

// GetUserActivityLogs mocks the GetUserActivityLogs method
func (m *MockLoggerClient) GetUserActivityLogs(ctx context.Context, userID string, actionType string) (interface{}, error) {
	args := m.Called(ctx, userID, actionType)
	return args.Get(0), args.Error(1)
}

// Close mocks the Close method
func (m *MockLoggerClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

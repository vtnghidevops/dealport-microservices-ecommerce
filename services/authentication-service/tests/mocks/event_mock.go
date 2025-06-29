package mocks

import (
	"context"

	"authentication-service/internal/event"

	"github.com/stretchr/testify/mock"
)

// MockEventEmitter là mock của event.Emitter interface
type MockEventEmitter struct {
	mock.Mock
}

// AsEmitter returns the MockEventEmitter as an *event.Emitter
// This is a hack for testing to satisfy the interface requirement
func (m *MockEventEmitter) AsEmitter() *event.Emitter {
	return &event.Emitter{}
}

// EmitUserRegistered mocks phát sự kiện user.registered
func (m *MockEventEmitter) EmitUserRegistered(userData event.UserRegisteredData) error {
	args := m.Called(userData)
	return args.Error(0)
}

// EmitLoginFailed mocks phát sự kiện auth.login.failed
func (m *MockEventEmitter) EmitLoginFailed(email, reason string, metadata map[string]interface{}) error {
	args := m.Called(email, reason, metadata)
	return args.Error(0)
}

// EmitLoginSucceeded mocks phát sự kiện auth.login.succeeded
func (m *MockEventEmitter) EmitLoginSucceeded(userID, email, role string) error {
	args := m.Called(userID, email, role)
	return args.Error(0)
}

// EmitOTPGenerated mocks phát sự kiện auth.otp.generated
func (m *MockEventEmitter) EmitOTPGenerated(
	email string,
	otp string,
	purpose string,
	expiresIn int,
	message string,
	action string,
) error {
	args := m.Called(email, otp, purpose, expiresIn, message, action)
	return args.Error(0)
}

// Emit mocks phát một sự kiện tự do
func (m *MockEventEmitter) Emit(ctx context.Context, eventName string, data interface{}) error {
	args := m.Called(ctx, eventName, data)
	return args.Error(0)
}

// Push mocks đẩy một sự kiện vào hàng đợi
func (m *MockEventEmitter) Push(eventName string, data interface{}) error {
	args := m.Called(eventName, data)
	return args.Error(0)
}

// EmitLoginSuccess mocks phát sự kiện log.INFO.user.login_success
func (m *MockEventEmitter) EmitLoginSuccess(userID, email string, metadata map[string]interface{}) error {
	args := m.Called(userID, email, metadata)
	return args.Error(0)
}

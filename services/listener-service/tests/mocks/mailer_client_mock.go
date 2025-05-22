package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockMailClient is a mock for the MailClient
type MockMailClient struct {
	mock.Mock
}

// SendEmail mocks the SendEmail method
func (m *MockMailClient) SendEmail(email, subject, message string) error {
	args := m.Called(email, subject, message)
	return args.Error(0)
}

// SendTemplateEmail mocks the SendTemplateEmail method
func (m *MockMailClient) SendTemplateEmail(email, subject, template string, variables map[string]string) error {
	args := m.Called(email, subject, template, variables)
	return args.Error(0)
}

// SendOTPEmail mocks the SendOTPEmail method
func (m *MockMailClient) SendOTPEmail(email, otp, purpose string, expiresIn int32, message, actionText string) error {
	args := m.Called(email, otp, purpose, expiresIn, message, actionText)
	return args.Error(0)
}

// SendWelcomeEmail mocks the SendWelcomeEmail method
func (m *MockMailClient) SendWelcomeEmail(email, firstName, lastName string) error {
	args := m.Called(email, firstName, lastName)
	return args.Error(0)
}

// SendPasswordResetEmail mocks the SendPasswordResetEmail method
func (m *MockMailClient) SendPasswordResetEmail(email, tokenHash, expiresAt string) error {
	args := m.Called(email, tokenHash, expiresAt)
	return args.Error(0)
}

// SendPasswordChangedEmail mocks the SendPasswordChangedEmail method
func (m *MockMailClient) SendPasswordChangedEmail(email, changedAt string) error {
	args := m.Called(email, changedAt)
	return args.Error(0)
}

// SendOrderConfirmationEmail mocks the SendOrderConfirmationEmail method
func (m *MockMailClient) SendOrderConfirmationEmail(email, orderID, orderNumber, customerName, orderDate, status,
	paymentMethod, total, items, shippingAddress, shippingName, shippingPhone string) error {
	args := m.Called(email, orderID, orderNumber, customerName, orderDate, status,
		paymentMethod, total, items, shippingAddress, shippingName, shippingPhone)
	return args.Error(0)
}

// Close mocks the Close method
func (m *MockMailClient) Close() error {
	args := m.Called()
	return args.Error(0)
} 
package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockMailClient mô phỏng MailClient
type MockMailClient struct {
	mock.Mock
}

// SendMail mocks gửi email
func (m *MockMailClient) SendMail(to, subject, content string) error {
	args := m.Called(to, subject, content)
	return args.Error(0)
}

// SendRegistrationOTP mocks gửi OTP đăng ký
func (m *MockMailClient) SendRegistrationOTP(to, otp string, expiresIn int) error {
	args := m.Called(to, otp, expiresIn)
	return args.Error(0)
}

// SendPasswordResetOTP mocks gửi OTP reset mật khẩu
func (m *MockMailClient) SendPasswordResetOTP(to, otp string, expiresIn int) error {
	args := m.Called(to, otp, expiresIn)
	return args.Error(0)
}

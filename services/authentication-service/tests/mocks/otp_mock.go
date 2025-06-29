package mocks

import (
	"authentication-service/internal/util"

	"github.com/stretchr/testify/mock"
)

// MockOTPManager mô phỏng OTPManager trong util package
type MockOTPManager struct {
	mock.Mock
}

// AsOTPManager returns the MockOTPManager as a *util.OTPManager
// This is a hack for testing to satisfy the interface requirement
func (m *MockOTPManager) AsOTPManager() *util.OTPManager {
	return &util.OTPManager{}
}

// GenerateOTP mocks tạo OTP mới
func (m *MockOTPManager) GenerateOTP(email string, purpose util.OTPPurpose) (string, error) {
	args := m.Called(email, string(purpose))
	return args.String(0), args.Error(1)
}

// VerifyOTP mocks xác thực OTP
func (m *MockOTPManager) VerifyOTP(email, otp string, purpose util.OTPPurpose) (bool, error) {
	args := m.Called(email, otp, string(purpose))
	return args.Bool(0), args.Error(1)
}

// GetRemainingTime mocks lấy thời gian còn lại của OTP
func (m *MockOTPManager) GetRemainingTime(email string) (int, error) {
	args := m.Called(email)
	return args.Int(0), args.Error(1)
}

// CleanupExpiredOTPs mocks dọn dẹp OTP hết hạn
func (m *MockOTPManager) CleanupExpiredOTPs() {
	m.Called()
}

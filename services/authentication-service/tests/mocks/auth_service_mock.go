package mocks

import (
	"context"

	"authentication-service/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockAuthService mocks the AuthService interface
type MockAuthService struct {
	mock.Mock
}

// Login mocks the Login method
func (m *MockAuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenDetails, *domain.User, error) {
	args := m.Called(ctx, req)

	var td *domain.TokenDetails
	var user *domain.User

	if args.Get(0) != nil {
		td = args.Get(0).(*domain.TokenDetails)
	}

	if args.Get(1) != nil {
		user = args.Get(1).(*domain.User)
	}

	return td, user, args.Error(2)
}

// Register mocks the Register method
func (m *MockAuthService) Register(ctx context.Context, req *domain.RegisterRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

// VerifyRegistration mocks the VerifyRegistration method
func (m *MockAuthService) VerifyRegistration(ctx context.Context, req *domain.VerifyRegistrationRequest) (*domain.TokenDetails, *domain.User, error) {
	args := m.Called(ctx, req)

	var td *domain.TokenDetails
	var user *domain.User

	if args.Get(0) != nil {
		td = args.Get(0).(*domain.TokenDetails)
	}

	if args.Get(1) != nil {
		user = args.Get(1).(*domain.User)
	}

	return td, user, args.Error(2)
}

// ValidateToken mocks the ValidateToken method
func (m *MockAuthService) ValidateToken(ctx context.Context, tokenString string) (*domain.TokenMetadata, error) {
	args := m.Called(ctx, tokenString)

	var tm *domain.TokenMetadata

	if args.Get(0) != nil {
		tm = args.Get(0).(*domain.TokenMetadata)
	}

	return tm, args.Error(1)
}

// RefreshToken mocks the RefreshToken method
func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenDetails, error) {
	args := m.Called(ctx, refreshToken)

	var td *domain.TokenDetails

	if args.Get(0) != nil {
		td = args.Get(0).(*domain.TokenDetails)
	}

	return td, args.Error(1)
}

// RequestPasswordReset mocks the RequestPasswordReset method
func (m *MockAuthService) RequestPasswordReset(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

// VerifyPasswordReset mocks the VerifyPasswordReset method
func (m *MockAuthService) VerifyPasswordReset(ctx context.Context, req *domain.VerifyPasswordResetRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

// UpdatePassword mocks the UpdatePassword method
func (m *MockAuthService) UpdatePassword(ctx context.Context, req *domain.UpdatePasswordRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// GetUserByID mocks the GetUserByID method
func (m *MockAuthService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)

	var user *domain.User

	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User)
	}

	return user, args.Error(1)
}

// RequestOTP mocks the RequestOTP method
func (m *MockAuthService) RequestOTP(ctx context.Context, email string, purpose string) error {
	args := m.Called(ctx, email, purpose)
	return args.Error(0)
}

// VerifyOTP mocks the VerifyOTP method
func (m *MockAuthService) VerifyOTP(ctx context.Context, req *domain.VerifyOTPRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

// CheckAccountExists mocks the CheckAccountExists method
func (m *MockAuthService) CheckAccountExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

// Logout mocks the Logout method
func (m *MockAuthService) Logout(ctx context.Context, userID, email string) error {
	args := m.Called(ctx, userID, email)
	return args.Error(0)
}

// LogoutFromAllDevices mocks the LogoutFromAllDevices method
func (m *MockAuthService) LogoutFromAllDevices(ctx context.Context, userID, email string) error {
	args := m.Called(ctx, userID, email)
	return args.Error(0)
}

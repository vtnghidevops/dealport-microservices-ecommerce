package mocks

import (
	authpb "broker-service/proto/auth"
	"context"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockAuthServiceClient là mock cho AuthServiceClient
type MockAuthServiceClient struct {
	mock.Mock
}

// Đảm bảo rằng MockAuthServiceClient triển khai AuthServiceClient interface
var _ authpb.AuthServiceClient = (*MockAuthServiceClient)(nil)

// Register mocks the Register method
func (m *MockAuthServiceClient) Register(ctx context.Context, in *authpb.RegisterRequest, opts ...grpc.CallOption) (*authpb.RegisterResponse, error) {
	args := m.Called(ctx, in)

	// Trả về nil nếu không có response được cấu hình
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*authpb.RegisterResponse), args.Error(1)
}

// VerifyRegistration mocks the VerifyRegistration method
func (m *MockAuthServiceClient) VerifyRegistration(ctx context.Context, in *authpb.VerifyRegistrationRequest, opts ...grpc.CallOption) (*authpb.VerifyRegistrationResponse, error) {
	args := m.Called(ctx, in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*authpb.VerifyRegistrationResponse), args.Error(1)
}

// Login mocks the Login method
func (m *MockAuthServiceClient) Login(ctx context.Context, in *authpb.LoginRequest, opts ...grpc.CallOption) (*authpb.LoginResponse, error) {
	args := m.Called(ctx, in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*authpb.LoginResponse), args.Error(1)
}

// Validate mocks the Validate method
func (m *MockAuthServiceClient) Validate(ctx context.Context, in *authpb.ValidateRequest, opts ...grpc.CallOption) (*authpb.ValidateResponse, error) {
	args := m.Called(ctx, in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*authpb.ValidateResponse), args.Error(1)
}

// RefreshToken mocks the RefreshToken method
func (m *MockAuthServiceClient) RefreshToken(ctx context.Context, in *authpb.RefreshTokenRequest, opts ...grpc.CallOption) (*authpb.RefreshTokenResponse, error) {
	args := m.Called(ctx, in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*authpb.RefreshTokenResponse), args.Error(1)
}

// RequestPasswordReset mocks the RequestPasswordReset method
func (m *MockAuthServiceClient) RequestPasswordReset(ctx context.Context, in *authpb.PasswordResetRequest, opts ...grpc.CallOption) (*authpb.PasswordResetResponse, error) {
	args := m.Called(ctx, in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*authpb.PasswordResetResponse), args.Error(1)
}

// VerifyPasswordReset mocks the VerifyPasswordReset method
func (m *MockAuthServiceClient) VerifyPasswordReset(ctx context.Context, in *authpb.VerifyPasswordResetRequest, opts ...grpc.CallOption) (*authpb.VerifyPasswordResetResponse, error) {
	args := m.Called(ctx, in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*authpb.VerifyPasswordResetResponse), args.Error(1)
}

// UpdatePassword mocks the UpdatePassword method
func (m *MockAuthServiceClient) UpdatePassword(ctx context.Context, in *authpb.UpdatePasswordRequest, opts ...grpc.CallOption) (*authpb.UpdatePasswordResponse, error) {
	args := m.Called(ctx, in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*authpb.UpdatePasswordResponse), args.Error(1)
}

// RequestOTP mocks the RequestOTP method
func (m *MockAuthServiceClient) RequestOTP(ctx context.Context, in *authpb.OTPRequest, opts ...grpc.CallOption) (*authpb.OTPResponse, error) {
	args := m.Called(ctx, in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*authpb.OTPResponse), args.Error(1)
}

// VerifyOTP mocks the VerifyOTP method
func (m *MockAuthServiceClient) VerifyOTP(ctx context.Context, in *authpb.VerifyOTPRequest, opts ...grpc.CallOption) (*authpb.VerifyOTPResponse, error) {
	args := m.Called(ctx, in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*authpb.VerifyOTPResponse), args.Error(1)
}

// CheckAccountExists mocks the CheckAccountExists method
func (m *MockAuthServiceClient) CheckAccountExists(ctx context.Context, in *authpb.CheckAccountExistsRequest, opts ...grpc.CallOption) (*authpb.CheckAccountExistsResponse, error) {
	args := m.Called(ctx, in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*authpb.CheckAccountExistsResponse), args.Error(1)
}

// Logout mocks the Logout method
func (m *MockAuthServiceClient) Logout(ctx context.Context, in *authpb.LogoutRequest, opts ...grpc.CallOption) (*authpb.LogoutResponse, error) {
	args := m.Called(ctx, in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*authpb.LogoutResponse), args.Error(1)
}

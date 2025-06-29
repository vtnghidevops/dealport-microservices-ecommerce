package mocks

import (
	"context"

	authpb "broker-service/proto/auth"

	"github.com/stretchr/testify/mock"
)

// AuthClientAdapter là một adapter để sử dụng MockAuthServiceClient trong tests
type AuthClientAdapter struct {
	Mock *MockAuthServiceClient
}

// NewAuthClientAdapter tạo một adapter mới đóng vai trò là AuthServiceClient
func NewAuthClientAdapter(mock *MockAuthServiceClient) authpb.AuthServiceClient {
	return mock
}

// Tạo một function helper để khởi tạo MockAuthServiceClient
// và thêm mock behavior cho các gRPC calls phổ biến
func SetupMockAuthClient() *MockAuthServiceClient {
	mockClient := new(MockAuthServiceClient)
	return mockClient
}

// AddRegisterMock cấu hình mock cho Register gRPC call
func AddRegisterMock(mockClient *MockAuthServiceClient, email, password, firstName, lastName, username, userID string) {
	mockClient.On(
		"Register",
		context.Background(),
		&authpb.RegisterRequest{
			Email:     email,
			Password:  password,
			FirstName: firstName,
			LastName:  lastName,
			Username:  username,
		},
		mock.Anything,
	).Return(&authpb.RegisterResponse{
		Success: true,
		Message: "Registration successful. Please check your email for verification code.",
		UserId:  userID,
	}, nil).Once()
}

// AddLoginMock cấu hình mock cho Login gRPC call
func AddLoginMock(mockClient *MockAuthServiceClient, email, password, userID, accessToken, refreshToken string) {
	mockClient.On(
		"Login",
		context.Background(),
		&authpb.LoginRequest{
			Email:    email,
			Password: password,
		},
		mock.Anything,
	).Return(&authpb.LoginResponse{
		Success:      true,
		Message:      "Logged in successfully",
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserInfo: &authpb.UserInfo{
			Email: email,
		},
	}, nil).Once()
}

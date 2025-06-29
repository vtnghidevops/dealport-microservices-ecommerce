package service

import (
	"context"

	"authentication-service/internal/domain"
)

// AuthService defines authentication service operations
type AuthService interface {
	// Register creates a new user account
	Register(ctx context.Context, req *domain.RegisterRequest) (string, error)

	// VerifyRegistration verifies registration OTP and completes registration
	VerifyRegistration(ctx context.Context, req *domain.VerifyRegistrationRequest) (*domain.TokenDetails, *domain.User, error)

	// Login authenticates a user and returns tokens
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenDetails, *domain.User, error)

	// Logout logs out a user
	Logout(ctx context.Context, userID, email string) error

	// LogoutFromAllDevices logs out a user from all devices
	LogoutFromAllDevices(ctx context.Context, userID, email string) error

	// ValidateToken validates an access token and returns claims
	ValidateToken(ctx context.Context, token string) (*domain.TokenMetadata, error)

	// RefreshToken refreshes an access token using a refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenDetails, error)

	// RequestPasswordReset initiates password reset process
	RequestPasswordReset(ctx context.Context, email string) error

	// VerifyPasswordReset verifies password reset OTP
	VerifyPasswordReset(ctx context.Context, req *domain.VerifyPasswordResetRequest) (string, error)

	// UpdatePassword updates user's password using reset token or current password
	UpdatePassword(ctx context.Context, req *domain.UpdatePasswordRequest) error

	// RequestOTP generates and sends a new OTP for a specific purpose
	RequestOTP(ctx context.Context, email string, purpose string) error

	// VerifyOTP verifies an OTP code for a specific purpose
	VerifyOTP(ctx context.Context, req *domain.VerifyOTPRequest) (string, error)

	// CheckAccountExists checks if an account with the given email exists
	CheckAccountExists(ctx context.Context, email string) (bool, error)

	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}

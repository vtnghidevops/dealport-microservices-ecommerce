package service

import (
	"context"

	"authentication-service/internal/domain"
)

// AuthService defines authentication service operations
type AuthService interface {
	// Register creates a new user account
	Register(ctx context.Context, req *domain.RegisterRequest) (string, error)

	// Login authenticates a user and returns tokens
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenDetails, *domain.User, error)

	// ValidateToken validates an access token and returns claims
	ValidateToken(ctx context.Context, token string) (*domain.TokenMetadata, error)

	// RefreshToken refreshes an access token using a refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenDetails, error)

	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}

package repository

import (
	"context"

	"authentication-service/internal/domain"
)

// UserRepository interface defines methods to interact with user data
type UserRepository interface {
	// Create a new user
	CreateUser(ctx context.Context, user *domain.User) error

	// Get a user by ID
	GetUserByID(ctx context.Context, id string) (*domain.User, error)

	// Get a user by email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)

	// Update user details
	UpdateUser(ctx context.Context, user *domain.User) error

	// Update user's refresh token
	UpdateRefreshToken(ctx context.Context, userID, refreshToken string) error

	// Logout from all devices
	LogoutFromAllDevices(ctx context.Context, userID string) error

	// Check if a user exists by email
	UserExists(ctx context.Context, email string) (bool, error)
}

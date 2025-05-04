package service

import (
	"context"

	"user-service/internal/domain"
)

// UserService defines the operations for working with users
type UserService interface {
	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id string) (*domain.User, error)

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)

	// GetUsers retrieves a list of users with pagination and filtering
	GetUsers(ctx context.Context, filter *domain.UserFilter) ([]*domain.User, int, error)

	// CreateUser creates a new user
	CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)

	// UpdateUser updates user information
	UpdateUser(ctx context.Context, req *domain.UpdateUserRequest) (*domain.User, error)

	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id string) error

	// SearchUsers searches for users based on criteria
	SearchUsers(ctx context.Context, params *domain.SearchParams) ([]*domain.User, int, error)

	// CreateAddress creates a new address for a user
	CreateAddress(ctx context.Context, userID string, address *domain.Address) error

	// UpdateAddress updates an address
	UpdateAddress(ctx context.Context, address *domain.Address) error

	// DeleteAddress deletes an address
	DeleteAddress(ctx context.Context, id string, userID string) error

	// SetDefaultAddress sets an address as default
	SetDefaultAddress(ctx context.Context, id string, userID string) error

	// AddToWishlist adds a product to a user's wishlist
	AddToWishlist(ctx context.Context, req *domain.AddToWishlistRequest) error

	// RemoveFromWishlist removes a product from a user's wishlist
	RemoveFromWishlist(ctx context.Context, req *domain.RemoveFromWishlistRequest) error

	// GetWishlist retrieves a user's wishlist
	GetWishlist(ctx context.Context, req *domain.GetWishlistRequest) (*domain.GetWishlistResponse, error)

	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error

	// RequestPasswordReset initiates a password reset request
	RequestPasswordReset(ctx context.Context, email string) error

	// ValidateCredentials validates user login credentials
	ValidateCredentials(ctx context.Context, email, password string) (*domain.User, error)

	// LogoutUser logs out a user
	LogoutUser(ctx context.Context, userID string) error

	// LogUserActivity logs a user activity
	LogUserActivity(ctx context.Context, action, userID, message string, metadata map[string]interface{}) error

	// GetUserActivityLogs retrieves user activity logs
	GetUserActivityLogs(ctx context.Context, userID string, actionType string) (interface{}, error)
}

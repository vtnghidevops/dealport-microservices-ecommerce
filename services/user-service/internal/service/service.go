package service

import (
	"context"

	"user-service/internal/domain"
)

// UserService defines the interface for user operations
type UserService interface {
	// User Operations
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUsers(ctx context.Context, filter *domain.UserFilter) ([]*domain.User, int, error)
	CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	UpdateUser(ctx context.Context, req *domain.UpdateUserRequest) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	SearchUsers(ctx context.Context, params *domain.SearchParams) ([]*domain.User, int, error)

	// Address Operations
	CreateAddress(ctx context.Context, userID string, address *domain.Address) error
	UpdateAddress(ctx context.Context, address *domain.Address) error
	DeleteAddress(ctx context.Context, id string, userID string) error
	SetDefaultAddress(ctx context.Context, id string, userID string) error

	// Wishlist Operations
	AddToWishlist(ctx context.Context, req *domain.AddToWishlistRequest) error
	RemoveFromWishlist(ctx context.Context, req *domain.RemoveFromWishlistRequest) error
	GetWishlist(ctx context.Context, req *domain.GetWishlistRequest) (*domain.GetWishlistResponse, error)

	// Authentication related operations
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	RequestPasswordReset(ctx context.Context, email string) error
	ValidateCredentials(ctx context.Context, email, password string) (*domain.User, error)
	LogoutUser(ctx context.Context, userID string) error

	// Activity Tracking
	LogUserActivity(ctx context.Context, action, userID, message string, metadata map[string]interface{}) error
	GetUserActivityLogs(ctx context.Context, userID string, actionType string) (interface{}, error)

	// Admin Dashboard Operations
	GetUserStatistics(ctx context.Context) (*domain.UserStatistics, error)
	GetUserActivityChart(ctx context.Context, days int) (*domain.UserActivityChart, error)

	// Order Data Synchronization
	SyncUserOrderData(ctx context.Context, userID string, orderCount int, totalSpend float64) error
}

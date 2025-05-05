package repository

import (
	"context"
	"time"

	"user-service/internal/domain"
)

// UserRepository interface defines methods for user data access
type UserRepository interface {
	// Create a new user
	CreateUser(ctx context.Context, user *domain.User) error

	// Get a user by ID
	GetUserByID(ctx context.Context, id string) (*domain.User, error)

	// Get a user by email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)

	// Get all users with optional filtering
	GetUsers(ctx context.Context, filter *domain.UserFilter) ([]*domain.User, int, error)

	// Update user details
	UpdateUser(ctx context.Context, user *domain.User) error

	// Delete a user (soft delete)
	DeleteUser(ctx context.Context, id string) error

	// Search for users
	SearchUsers(ctx context.Context, params *domain.SearchParams) ([]*domain.User, int, error)

	// Check if a user exists by email
	UserExists(ctx context.Context, email string) (bool, error)

	// Create a new address for a user
	CreateAddress(ctx context.Context, address *domain.Address) error

	// Get all addresses for a user
	GetAddresses(ctx context.Context, userID string) ([]*domain.Address, error)

	// Update an address
	UpdateAddress(ctx context.Context, address *domain.Address) error

	// Delete an address
	DeleteAddress(ctx context.Context, id string, userID string) error

	// Set an address as default
	SetDefaultAddress(ctx context.Context, id string, userID string) error

	// Wishlist methods
	AddToWishlist(ctx context.Context, userID string, productID int, notes string) error
	RemoveFromWishlist(ctx context.Context, userID string, productID int) error
	GetWishlist(ctx context.Context, userID string) ([]*domain.WishlistItem, error)

	// Order and spending related methods
	GetUserWithOrderCount(ctx context.Context, userID string) (int, error)
	GetUserTotalSpend(ctx context.Context, userID string) (float64, error)
	UpdateUserOrderCount(ctx context.Context, userID string, count int) error
	UpdateUserTotalSpend(ctx context.Context, userID string, amount float64) error
	GetRepeatCustomers(ctx context.Context) ([]*domain.User, error)
	GetNewUsersCount(ctx context.Context, since time.Time) (int, error)
	GetActiveUsersCount(ctx context.Context, since time.Time) (int, error)

	// Statistics methods
	GetNewUsersSince(ctx context.Context, filter *domain.UserFilter, since time.Time) ([]*domain.User, error)
	GetActiveUsersSince(ctx context.Context, filter *domain.UserFilter, since time.Time) ([]*domain.User, error)
	GetUsersWithOrderCount(ctx context.Context, filter *domain.UserFilter, minOrders int) ([]*domain.User, error)
	GetUserActivityCountForDay(ctx context.Context, date time.Time) (int, error)
}

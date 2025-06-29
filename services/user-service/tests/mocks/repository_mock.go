package mocks

import (
	"context"
	"time"

	"user-service/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock of UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

// CreateUser mocks the CreateUser method
func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// GetUserByID mocks the GetUserByID method
func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)

	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User)
	}

	return user, args.Error(1)
}

// GetUserByEmail mocks the GetUserByEmail method
func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)

	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User)
	}

	return user, args.Error(1)
}

// GetUsers mocks the GetUsers method
func (m *MockUserRepository) GetUsers(ctx context.Context, filter *domain.UserFilter) ([]*domain.User, int, error) {
	args := m.Called(ctx, filter)

	var users []*domain.User
	if args.Get(0) != nil {
		users = args.Get(0).([]*domain.User)
	}

	return users, args.Int(1), args.Error(2)
}

// UpdateUser mocks the UpdateUser method
func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// DeleteUser mocks the DeleteUser method
func (m *MockUserRepository) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// SearchUsers mocks the SearchUsers method
func (m *MockUserRepository) SearchUsers(ctx context.Context, params *domain.SearchParams) ([]*domain.User, int, error) {
	args := m.Called(ctx, params)

	var users []*domain.User
	if args.Get(0) != nil {
		users = args.Get(0).([]*domain.User)
	}

	return users, args.Int(1), args.Error(2)
}

// UserExists mocks the UserExists method
func (m *MockUserRepository) UserExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

// CreateAddress mocks the CreateAddress method
func (m *MockUserRepository) CreateAddress(ctx context.Context, address *domain.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

// GetAddresses mocks the GetAddresses method
func (m *MockUserRepository) GetAddresses(ctx context.Context, userID string) ([]*domain.Address, error) {
	args := m.Called(ctx, userID)

	var addresses []*domain.Address
	if args.Get(0) != nil {
		addresses = args.Get(0).([]*domain.Address)
	}

	return addresses, args.Error(1)
}

// UpdateAddress mocks the UpdateAddress method
func (m *MockUserRepository) UpdateAddress(ctx context.Context, address *domain.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

// DeleteAddress mocks the DeleteAddress method
func (m *MockUserRepository) DeleteAddress(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

// SetDefaultAddress mocks the SetDefaultAddress method
func (m *MockUserRepository) SetDefaultAddress(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

// AddToWishlist mocks the AddToWishlist method
func (m *MockUserRepository) AddToWishlist(ctx context.Context, userID string, productID int, notes string) error {
	args := m.Called(ctx, userID, productID, notes)
	return args.Error(0)
}

// RemoveFromWishlist mocks the RemoveFromWishlist method
func (m *MockUserRepository) RemoveFromWishlist(ctx context.Context, userID string, productID int) error {
	args := m.Called(ctx, userID, productID)
	return args.Error(0)
}

// GetWishlist mocks the GetWishlist method
func (m *MockUserRepository) GetWishlist(ctx context.Context, userID string) ([]*domain.WishlistItem, error) {
	args := m.Called(ctx, userID)

	var items []*domain.WishlistItem
	if args.Get(0) != nil {
		items = args.Get(0).([]*domain.WishlistItem)
	}

	return items, args.Error(1)
}

// GetUserWithOrderCount mocks the GetUserWithOrderCount method
func (m *MockUserRepository) GetUserWithOrderCount(ctx context.Context, userID string) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

// GetUserTotalSpend mocks the GetUserTotalSpend method
func (m *MockUserRepository) GetUserTotalSpend(ctx context.Context, userID string) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

// UpdateUserOrderCount mocks the UpdateUserOrderCount method
func (m *MockUserRepository) UpdateUserOrderCount(ctx context.Context, userID string, count int) error {
	args := m.Called(ctx, userID, count)
	return args.Error(0)
}

// UpdateUserTotalSpend mocks the UpdateUserTotalSpend method
func (m *MockUserRepository) UpdateUserTotalSpend(ctx context.Context, userID string, amount float64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

// GetRepeatCustomers mocks the GetRepeatCustomers method
func (m *MockUserRepository) GetRepeatCustomers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)

	var users []*domain.User
	if args.Get(0) != nil {
		users = args.Get(0).([]*domain.User)
	}

	return users, args.Error(1)
}

// GetNewUsersCount mocks the GetNewUsersCount method
func (m *MockUserRepository) GetNewUsersCount(ctx context.Context, since time.Time) (int, error) {
	args := m.Called(ctx, since)
	return args.Int(0), args.Error(1)
}

// GetActiveUsersCount mocks the GetActiveUsersCount method
func (m *MockUserRepository) GetActiveUsersCount(ctx context.Context, since time.Time) (int, error) {
	args := m.Called(ctx, since)
	return args.Int(0), args.Error(1)
}

// GetNewUsersSince mocks the GetNewUsersSince method
func (m *MockUserRepository) GetNewUsersSince(ctx context.Context, filter *domain.UserFilter, since time.Time) ([]*domain.User, error) {
	args := m.Called(ctx, filter, since)

	var users []*domain.User
	if args.Get(0) != nil {
		users = args.Get(0).([]*domain.User)
	}

	return users, args.Error(1)
}

// GetActiveUsersSince mocks the GetActiveUsersSince method
func (m *MockUserRepository) GetActiveUsersSince(ctx context.Context, filter *domain.UserFilter, since time.Time) ([]*domain.User, error) {
	args := m.Called(ctx, filter, since)

	var users []*domain.User
	if args.Get(0) != nil {
		users = args.Get(0).([]*domain.User)
	}

	return users, args.Error(1)
}

// GetUsersWithOrderCount mocks the GetUsersWithOrderCount method
func (m *MockUserRepository) GetUsersWithOrderCount(ctx context.Context, filter *domain.UserFilter, minOrders int) ([]*domain.User, error) {
	args := m.Called(ctx, filter, minOrders)

	var users []*domain.User
	if args.Get(0) != nil {
		users = args.Get(0).([]*domain.User)
	}

	return users, args.Error(1)
}

// GetUserActivityCountForDay mocks the GetUserActivityCountForDay method
func (m *MockUserRepository) GetUserActivityCountForDay(ctx context.Context, date time.Time) (int, error) {
	args := m.Called(ctx, date)
	return args.Int(0), args.Error(1)
}

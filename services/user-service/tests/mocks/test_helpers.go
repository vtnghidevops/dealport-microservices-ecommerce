package mocks

import (
	"time"

	"user-service/internal/domain"
	"user-service/internal/service"

	"github.com/stretchr/testify/mock"
)

// CreateMockUser creates a mock user with the given details
func CreateMockUser(id, email, firstName, lastName string) *domain.User {
	phone := "1234567890"
	role := "customer"
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()

	return &domain.User{
		ID:           id,
		Email:        email,
		PasswordHash: "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK", // hash for "password"
		FirstName:    firstName,
		LastName:     lastName,
		Phone:        &phone,
		Role:         role,
		Active:       true,
		Status:       "active",
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Addresses:    []domain.Address{},
	}
}

// CreateMockUserWithAddresses creates a mock user with the given addresses
func CreateMockUserWithAddresses(id, email, firstName, lastName string, addresses []domain.Address) *domain.User {
	user := CreateMockUser(id, email, firstName, lastName)
	user.Addresses = addresses
	return user
}

// CreateTestUserService creates a mock repository and a user service for testing
func CreateTestUserService() (*MockUserRepository, service.UserService) {
	mockRepo := new(MockUserRepository)
	// Use nil event emitter for testing since we're mocking and don't need actual event publishing
	userService := service.NewUserService(mockRepo, nil)
	return mockRepo, userService
}

// SetupBasicUserMocks sets up basic mocks for user operations
func SetupBasicUserMocks(mockRepo *MockUserRepository, user *domain.User) {
	// Setup basic expectations for common operations
	mockRepo.On("GetUserByID", mock.Anything, user.ID).Return(user, nil).Maybe()
	mockRepo.On("GetUserByEmail", mock.Anything, user.Email).Return(user, nil).Maybe()
	mockRepo.On("GetWishlist", mock.Anything, user.ID).Return([]*domain.WishlistItem{}, nil).Maybe()
}

// SetupWishlistMocks sets up mocks for wishlist operations
func SetupWishlistMocks(mockRepo *MockUserRepository, userID string, items []*domain.WishlistItem) {
	mockRepo.On("GetWishlist", mock.Anything, userID).Return(items, nil).Maybe()
}

// CreateMockAddress creates a mock address with the given details
func CreateMockAddress(id, userID string, isDefault bool) *domain.Address {
	return &domain.Address{
		ID:          id,
		UserID:      userID,
		Name:        "Home Address",
		Phone:       "1234567890",
		Line1:       "123 Main St",
		City:        "New York",
		State:       "NY",
		PostalCode:  "10001",
		Country:     "USA",
		IsDefault:   isDefault,
		AddressType: "shipping",
	}
}

// SetupAddressMocks sets up mocks for address operations
func SetupAddressMocks(mockRepo *MockUserRepository, userID string, addresses []domain.Address) {
	mockRepo.On("GetAddresses", mock.Anything, userID).Return(addresses, nil).Maybe()
}

// CreateMockWishlistItem creates a mock wishlist item
func CreateMockWishlistItem(id int, userID string, productID int, productName string, price float64, notes string) *domain.WishlistItem {
	return &domain.WishlistItem{
		ID:        "wish" + string(rune(id+'0')),
		UserID:    userID,
		ProductID: productID,
		AddedAt:   time.Now().Add(-24 * time.Hour),
		Notes:     notes,
	}
}

package unit

import (
	"context"
	"errors"
	"testing"

	"user-service/internal/domain"
	"user-service/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGetUserByID tests the GetUserByID function
func TestGetUserByID(t *testing.T) {
	// Setup using helper functions
	mockRepo, userService := mocks.CreateTestUserService()

	// Create a mock user
	mockUser := mocks.CreateMockUser("user123", "test@example.com", "Test", "User")

	// Test case 1: Successfully get user
	t.Run("Successfully get user", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, "user123").Return(mockUser, nil).Once()
		mockRepo.On("GetWishlist", mock.Anything, "user123").Return([]*domain.WishlistItem{}, nil).Once()

		// Execute
		user, err := userService.GetUserByID(context.Background(), "user123")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, mockUser.ID, user.ID)
		assert.Equal(t, mockUser.Email, user.Email)
		assert.Empty(t, user.PasswordHash)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: User not found
	t.Run("User not found", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, "nonexistent").Return(nil, errors.New("user not found")).Once()

		// Execute
		user, err := userService.GetUserByID(context.Background(), "nonexistent")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestCreateUser tests the CreateUser function
func TestCreateUser(t *testing.T) {
	// Setup using helper functions
	mockRepo, userService := mocks.CreateTestUserService()

	// Test case 1: Successfully create user
	t.Run("Successfully create user", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("UserExists", mock.Anything, "new@example.com").Return(false, nil).Once()
		mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		// Create user request
		req := &domain.CreateUserRequest{
			Email:     "new@example.com",
			Password:  "secure_password",
			FirstName: "New",
			LastName:  "User",
			Role:      "customer",
		}

		// Execute
		user, err := userService.CreateUser(context.Background(), req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Email, user.Email)
		assert.Equal(t, req.FirstName, user.FirstName)
		assert.Equal(t, req.LastName, user.LastName)
		assert.Equal(t, req.Role, user.Role)
		assert.Empty(t, user.PasswordHash)
		assert.True(t, user.Active)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Email already exists
	t.Run("Email already exists", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("UserExists", mock.Anything, "existing@example.com").Return(true, nil).Once()

		// Create user request with existing email
		req := &domain.CreateUserRequest{
			Email:     "existing@example.com",
			Password:  "secure_password",
			FirstName: "Existing",
			LastName:  "User",
			Role:      "customer",
		}

		// Execute
		user, err := userService.CreateUser(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "already exists")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdateUser tests the UpdateUser function
func TestUpdateUser(t *testing.T) {
	// Create a mock user
	mockUser := mocks.CreateMockUser("user123", "test@example.com", "Test", "User")

	// Test case 1: Successfully update user
	t.Run("Successfully update user", func(t *testing.T) {
		// Create a test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, "user123").Return(mockUser, nil).Once()
		mockRepo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()
		// Allow GetWishlist to be called multiple times
		mockRepo.On("GetWishlist", mock.Anything, "user123").Return([]*domain.WishlistItem{}, nil).Maybe()

		// Create update request
		req := &domain.UpdateUserRequest{
			ID:        "user123",
			FirstName: "Updated",
			LastName:  "User",
			Phone:     "1234567890",
		}

		// Execute
		user, err := userService.UpdateUser(context.Background(), req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.FirstName, user.FirstName)
		assert.Equal(t, req.LastName, user.LastName)
		assert.Equal(t, req.Phone, *user.Phone)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Empty user ID
	t.Run("Empty user ID", func(t *testing.T) {
		// Create a test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Create update request with empty ID
		req := &domain.UpdateUserRequest{
			ID:        "", // Empty ID
			FirstName: "Updated",
			LastName:  "User",
		}

		// Execute
		user, err := userService.UpdateUser(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "ID is required")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: User not found
	t.Run("User not found", func(t *testing.T) {
		// Create a test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, "nonexistent").Return(nil, errors.New("user not found")).Once()

		// Create update request for non-existent user
		req := &domain.UpdateUserRequest{
			ID:        "nonexistent",
			FirstName: "Updated",
			LastName:  "User",
		}

		// Execute
		user, err := userService.UpdateUser(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestDeleteUser tests the DeleteUser function
func TestDeleteUser(t *testing.T) {
	// Setup using helper functions
	mockRepo, userService := mocks.CreateTestUserService()

	// Test case 1: Successfully delete user
	t.Run("Successfully delete user", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("DeleteUser", mock.Anything, "user123").Return(nil).Once()

		// Execute
		err := userService.DeleteUser(context.Background(), "user123")

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: User not found
	t.Run("User not found", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("DeleteUser", mock.Anything, "nonexistent").Return(errors.New("user not found")).Once()

		// Execute
		err := userService.DeleteUser(context.Background(), "nonexistent")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

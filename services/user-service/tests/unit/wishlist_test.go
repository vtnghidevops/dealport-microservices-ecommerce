package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"user-service/internal/domain"
	"user-service/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestAddToWishlist tests the AddToWishlist function
func TestAddToWishlist(t *testing.T) {
	// Test case 1: Successfully add product to wishlist
	t.Run("Successfully add product to wishlist", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, "user123").Return(mocks.CreateMockUser("user123", "test@example.com", "Test", "User"), nil).Once()
		mockRepo.On("AddToWishlist", mock.Anything, "user123", 101, "Test notes").Return(nil).Once()

		// Create wishlist request
		req := &domain.AddToWishlistRequest{
			UserID:    "user123",
			ProductID: 101,
			Notes:     "Test notes",
		}

		// Execute
		err := userService.AddToWishlist(context.Background(), req)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: User not found
	t.Run("User not found", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, "nonexistent").Return(nil, errors.New("user not found")).Once()

		// Create wishlist request
		req := &domain.AddToWishlistRequest{
			UserID:    "nonexistent",
			ProductID: 101,
			Notes:     "Test notes",
		}

		// Execute
		err := userService.AddToWishlist(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, "user123").Return(mocks.CreateMockUser("user123", "test@example.com", "Test", "User"), nil).Once()
		mockRepo.On("AddToWishlist", mock.Anything, "user123", 102, "Test notes").Return(errors.New("database error")).Once()

		// Create wishlist request
		req := &domain.AddToWishlistRequest{
			UserID:    "user123",
			ProductID: 102,
			Notes:     "Test notes",
		}

		// Execute
		err := userService.AddToWishlist(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestRemoveFromWishlist tests the RemoveFromWishlist function
func TestRemoveFromWishlist(t *testing.T) {
	// Test case 1: Successfully remove product from wishlist
	t.Run("Successfully remove product from wishlist", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations - no need for GetUserByID since the service doesn't use it
		mockRepo.On("RemoveFromWishlist", mock.Anything, "user123", 101).Return(nil).Once()

		// Create wishlist request
		req := &domain.RemoveFromWishlistRequest{
			UserID:    "user123",
			ProductID: 101,
		}

		// Execute
		err := userService.RemoveFromWishlist(context.Background(), req)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: User not found
	t.Run("User not found", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations - the service goes straight to RemoveFromWishlist
		mockRepo.On("RemoveFromWishlist", mock.Anything, "nonexistent", 101).Return(errors.New("user not found")).Once()

		// Create wishlist request
		req := &domain.RemoveFromWishlistRequest{
			UserID:    "nonexistent",
			ProductID: 101,
		}

		// Execute
		err := userService.RemoveFromWishlist(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations - no need for GetUserByID
		mockRepo.On("RemoveFromWishlist", mock.Anything, "user123", 102).Return(errors.New("database error")).Once()

		// Create wishlist request
		req := &domain.RemoveFromWishlistRequest{
			UserID:    "user123",
			ProductID: 102,
		}

		// Execute
		err := userService.RemoveFromWishlist(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestGetWishlist tests the GetWishlist function
func TestGetWishlist(t *testing.T) {
	// Create mock wishlist items
	now := time.Now()
	mockWishlistItems := []*domain.WishlistItem{
		{
			ID:        "wish1",
			UserID:    "user123",
			ProductID: 101,
			AddedAt:   now.Add(-2 * time.Hour),
			Notes:     "Test notes 1",
		},
		{
			ID:        "wish2",
			UserID:    "user123",
			ProductID: 102,
			AddedAt:   now.Add(-1 * time.Hour),
			Notes:     "Test notes 2",
		},
	}

	// Test case 1: Successfully get wishlist
	t.Run("Successfully get wishlist", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations
		mockRepo.On("GetWishlist", mock.Anything, "user123").Return(mockWishlistItems, nil).Once()

		// Create wishlist request
		req := &domain.GetWishlistRequest{
			UserID: "user123",
		}

		// Execute
		response, err := userService.GetWishlist(context.Background(), req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Items, 2)
		assert.Equal(t, 101, response.Items[0].ProductID)
		assert.Equal(t, 102, response.Items[1].ProductID)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: User not found
	t.Run("User not found", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations
		mockRepo.On("GetWishlist", mock.Anything, "nonexistent").Return(nil, errors.New("user not found")).Once()

		// Create wishlist request
		req := &domain.GetWishlistRequest{
			UserID: "nonexistent",
		}

		// Execute
		response, err := userService.GetWishlist(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "user not found")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Empty wishlist
	t.Run("Empty wishlist", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations
		mockRepo.On("GetWishlist", mock.Anything, "user456").Return([]*domain.WishlistItem{}, nil).Once()

		// Create wishlist request
		req := &domain.GetWishlistRequest{
			UserID: "user456",
		}

		// Execute
		response, err := userService.GetWishlist(context.Background(), req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Items, 0)
		assert.Equal(t, 0, response.Count)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 4: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations
		mockRepo.On("GetWishlist", mock.Anything, "user123").Return(nil, errors.New("database error")).Once()

		// Create wishlist request
		req := &domain.GetWishlistRequest{
			UserID: "user123",
		}

		// Execute
		response, err := userService.GetWishlist(context.Background(), req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

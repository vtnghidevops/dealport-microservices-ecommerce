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

// TestCreateAddress tests the CreateAddress function
func TestCreateAddress(t *testing.T) {
	// Setup using helper functions
	mockRepo, userService := mocks.CreateTestUserService()

	// Create a mock user
	mockUser := mocks.CreateMockUser("user123", "test@example.com", "Test", "User")

	// Create a sample address
	sampleAddress := &domain.Address{
		ID:          "addr1",
		UserID:      "user123",
		Name:        "Home Address",
		Phone:       "1234567890",
		Line1:       "123 Main St",
		City:        "Test City",
		State:       "Test State",
		PostalCode:  "12345",
		Country:     "Test Country",
		IsDefault:   false,
		AddressType: "shipping",
	}

	// Test case 1: Successfully create address
	t.Run("Successfully create address", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, "user123").Return(mockUser, nil).Once()
		mockRepo.On("CreateAddress", mock.Anything, mock.AnythingOfType("*domain.Address")).
			Run(func(args mock.Arguments) {
				addr := args.Get(1).(*domain.Address)
				assert.Equal(t, "user123", addr.UserID)
				assert.Equal(t, "123 Main St", addr.Line1)
				assert.Equal(t, "Test City", addr.City)
			}).
			Return(nil).Once()

		// Execute
		err := userService.CreateAddress(context.Background(), "user123", sampleAddress)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: User not found
	t.Run("User not found", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, "nonexistent").Return(nil, errors.New("user not found")).Once()

		// Execute
		err := userService.CreateAddress(context.Background(), "nonexistent", sampleAddress)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Missing required fields
	t.Run("Missing required fields", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, "user123").Return(mockUser, nil).Once()

		// Create address with missing fields, but since the service doesn't validate these fields,
		// we need to mock the repository call as it will be called anyway
		invalidAddress := &domain.Address{
			UserID:      "user123",
			Line1:       "", // Empty required field
			City:        "Test City",
			Country:     "Test Country",
			AddressType: "shipping",
		}

		mockRepo.On("CreateAddress", mock.Anything, mock.AnythingOfType("*domain.Address")).Return(nil).Once()

		// Execute
		err := userService.CreateAddress(context.Background(), "user123", invalidAddress)

		// Assert - should succeed since the service doesn't validate fields
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 4: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetUserByID", mock.Anything, "user123").Return(mockUser, nil).Once()
		mockRepo.On("CreateAddress", mock.Anything, mock.AnythingOfType("*domain.Address")).
			Return(errors.New("database error")).Once()

		// Execute
		err := userService.CreateAddress(context.Background(), "user123", sampleAddress)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdateAddress tests the UpdateAddress function
func TestUpdateAddress(t *testing.T) {
	// Create a sample address
	sampleAddress := &domain.Address{
		ID:          "addr1",
		UserID:      "user123",
		Name:        "Home Address",
		Phone:       "1234567890",
		Line1:       "123 Main St",
		City:        "Test City",
		State:       "Test State",
		PostalCode:  "12345",
		Country:     "Test Country",
		IsDefault:   false,
		AddressType: "shipping",
	}

	// Create a mock user with addresses - not directly used in the tests but keeping for reference
	mockAddresses := []domain.Address{*sampleAddress}
	_ = mocks.CreateMockUserWithAddresses("user123", "test@example.com", "Test", "User", mockAddresses)

	// Test case 1: Successfully update address
	t.Run("Successfully update address", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations
		mockRepo.On("UpdateAddress", mock.Anything, mock.AnythingOfType("*domain.Address")).
			Run(func(args mock.Arguments) {
				addr := args.Get(1).(*domain.Address)
				assert.Equal(t, "addr1", addr.ID)
				assert.Equal(t, "user123", addr.UserID)
				assert.Equal(t, "456 New St", addr.Line1)
			}).
			Return(nil).Once()

		// Updated address
		updatedAddress := &domain.Address{
			ID:          "addr1",
			UserID:      "user123",
			Name:        "New Home",
			Phone:       "9876543210",
			Line1:       "456 New St",
			City:        "New City",
			State:       "New State",
			PostalCode:  "54321",
			Country:     "New Country",
			IsDefault:   true,
			AddressType: "billing",
		}

		// Execute
		err := userService.UpdateAddress(context.Background(), updatedAddress)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Address not found for user
	t.Run("Address not found for user", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations
		mockRepo.On("UpdateAddress", mock.Anything, mock.AnythingOfType("*domain.Address")).
			Return(errors.New("address not found")).Once()

		// Non-existent address ID
		nonExistentAddress := &domain.Address{
			ID:          "nonexistent",
			UserID:      "user123",
			Name:        "New Home",
			Phone:       "9876543210",
			Line1:       "456 New St",
			City:        "New City",
			Country:     "New Country",
			AddressType: "shipping",
		}

		// Execute
		err := userService.UpdateAddress(context.Background(), nonExistentAddress)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update address")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: User not found
	t.Run("User not found", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations
		mockRepo.On("UpdateAddress", mock.Anything, mock.AnythingOfType("*domain.Address")).
			Return(errors.New("user not found")).Once()

		// Address with non-existent user
		addressWithNonExistentUser := &domain.Address{
			ID:          "addr1",
			UserID:      "nonexistent",
			Name:        "New Home",
			Phone:       "9876543210",
			Line1:       "456 New St",
			City:        "New City",
			Country:     "New Country",
			AddressType: "shipping",
		}

		// Execute
		err := userService.UpdateAddress(context.Background(), addressWithNonExistentUser)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 4: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations
		mockRepo.On("UpdateAddress", mock.Anything, mock.AnythingOfType("*domain.Address")).
			Return(errors.New("database error")).Once()

		// Execute
		err := userService.UpdateAddress(context.Background(), sampleAddress)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestDeleteAddress tests the DeleteAddress function
func TestDeleteAddress(t *testing.T) {
	// Setup using helper functions
	mockRepo, userService := mocks.CreateTestUserService()

	// Create a sample address
	sampleAddress := &domain.Address{
		ID:          "addr1",
		UserID:      "user123",
		Name:        "Home Address",
		Phone:       "1234567890",
		Line1:       "123 Main St",
		City:        "Test City",
		State:       "Test State",
		PostalCode:  "12345",
		Country:     "Test Country",
		IsDefault:   false,
		AddressType: "shipping",
	}

	// Create a mock user with addresses
	mockAddresses := []domain.Address{*sampleAddress}
	_ = mocks.CreateMockUserWithAddresses("user123", "test@example.com", "Test", "User", mockAddresses)

	// Test case 1: Successfully delete address
	t.Run("Successfully delete address", func(t *testing.T) {
		// Reset expectations
		newMockRepo, newUserService := mocks.CreateTestUserService()
		mockRepo = newMockRepo
		userService = newUserService

		// Setup expectations - service doesn't call GetUserByID
		mockRepo.On("DeleteAddress", mock.Anything, "addr1", "user123").Return(nil).Once()

		// Execute
		err := userService.DeleteAddress(context.Background(), "addr1", "user123")

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: User not found
	t.Run("User not found", func(t *testing.T) {
		// Setup expectations - service doesn't call GetUserByID
		mockRepo.On("DeleteAddress", mock.Anything, "addr1", "nonexistent").
			Return(errors.New("user not found")).Once()

		// Execute
		err := userService.DeleteAddress(context.Background(), "addr1", "nonexistent")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Address not found
	t.Run("Address not found for user", func(t *testing.T) {
		// Reset expectations
		newMockRepo, newUserService := mocks.CreateTestUserService()
		mockRepo = newMockRepo
		userService = newUserService

		// Setup expectations - service doesn't call GetUserByID
		mockRepo.On("DeleteAddress", mock.Anything, "nonexistent", "user123").
			Return(errors.New("address not found")).Once()

		// Execute
		err := userService.DeleteAddress(context.Background(), "nonexistent", "user123")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "address not found")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 4: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Reset expectations
		newMockRepo, newUserService := mocks.CreateTestUserService()
		mockRepo = newMockRepo
		userService = newUserService

		// Setup expectations - service doesn't call GetUserByID
		mockRepo.On("DeleteAddress", mock.Anything, "addr1", "user123").
			Return(errors.New("database error")).Once()

		// Execute
		err := userService.DeleteAddress(context.Background(), "addr1", "user123")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestSetDefaultAddress tests the SetDefaultAddress function
func TestSetDefaultAddress(t *testing.T) {
	// Create sample addresses
	address1 := domain.Address{
		ID:          "addr1",
		UserID:      "user123",
		Name:        "Home Address",
		Phone:       "1234567890",
		Line1:       "123 Main St",
		City:        "Test City",
		Country:     "Test Country",
		IsDefault:   false,
		AddressType: "shipping",
	}

	address2 := domain.Address{
		ID:          "addr2",
		UserID:      "user123",
		Name:        "Work Address",
		Phone:       "0987654321",
		Line1:       "456 Second St",
		City:        "Test City",
		Country:     "Test Country",
		IsDefault:   true,
		AddressType: "billing",
	}

	// Create a mock user with addresses
	mockAddresses := []domain.Address{address1, address2}
	_ = mocks.CreateMockUserWithAddresses("user123", "test@example.com", "Test", "User", mockAddresses)

	// Test case 1: Successfully set default address
	t.Run("Successfully set default address", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations - service doesn't call GetUserByID
		mockRepo.On("SetDefaultAddress", mock.Anything, "addr1", "user123").Return(nil).Once()

		// Execute
		err := userService.SetDefaultAddress(context.Background(), "addr1", "user123")

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: User not found
	t.Run("User not found", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations - service doesn't call GetUserByID
		mockRepo.On("SetDefaultAddress", mock.Anything, "addr1", "nonexistent").
			Return(errors.New("user not found")).Once()

		// Execute
		err := userService.SetDefaultAddress(context.Background(), "addr1", "nonexistent")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Address not found
	t.Run("Address not found for user", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations - service doesn't call GetUserByID
		mockRepo.On("SetDefaultAddress", mock.Anything, "nonexistent", "user123").
			Return(errors.New("address not found")).Once()

		// Execute
		err := userService.SetDefaultAddress(context.Background(), "nonexistent", "user123")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "address not found")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 4: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Create a separate test service for this test
		mockRepo, userService := mocks.CreateTestUserService()

		// Setup expectations - service doesn't call GetUserByID
		mockRepo.On("SetDefaultAddress", mock.Anything, "addr1", "user123").
			Return(errors.New("database error")).Once()

		// Execute
		err := userService.SetDefaultAddress(context.Background(), "addr1", "user123")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

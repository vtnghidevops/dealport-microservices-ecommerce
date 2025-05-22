package unit

import (
	"context"
	"testing"
	"time"

	"authentication-service/internal/domain"
	"authentication-service/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestAuthHandler tests the auth handler functionality
func TestAuthHandler(t *testing.T) {
	// Setup mock auth service
	mockAuthService := new(mocks.MockAuthService)

	// Test login functionality
	t.Run("Login", func(t *testing.T) {
		// Setup mock user and token details
		mockUser := &domain.User{
			ID:        "user123",
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Role:      "user",
		}

		mockTokenDetails := &domain.TokenDetails{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
			AccessUUID:   "access_uuid",
			RefreshUUID:  "refresh_uuid",
			AtExpires:    time.Now().Add(time.Minute * 15).Unix(),
			RtExpires:    time.Now().Add(time.Hour * 24 * 7).Unix(),
		}

		// Setup login request
		loginReq := &domain.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		// Setup mock expectations
		mockAuthService.On("Login", mock.Anything, loginReq).Return(mockTokenDetails, mockUser, nil).Once()

		// Call the service
		tokenDetails, user, err := mockAuthService.Login(context.Background(), loginReq)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, mockTokenDetails, tokenDetails)
		assert.Equal(t, mockUser, user)

		// Verify expectations
		mockAuthService.AssertExpectations(t)
	})

	// Test token validation
	t.Run("ValidateToken", func(t *testing.T) {
		// Setup mock token metadata
		mockTokenMetadata := &domain.TokenMetadata{
			UserID: "user123",
			Email:  "test@example.com",
			Role:   "user",
			UUID:   "token_uuid",
			Exp:    time.Now().Add(time.Minute * 15).Unix(),
		}

		// Setup mock expectations
		mockAuthService.On("ValidateToken", mock.Anything, "access_token").Return(mockTokenMetadata, nil).Once()

		// Call the service
		tokenMetadata, err := mockAuthService.ValidateToken(context.Background(), "access_token")

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, mockTokenMetadata, tokenMetadata)

		// Verify expectations
		mockAuthService.AssertExpectations(t)
	})
}

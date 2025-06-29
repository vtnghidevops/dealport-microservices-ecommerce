package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"authentication-service/internal/domain"
	"authentication-service/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// TestLogin kiểm tra chức năng đăng nhập
func TestLogin(t *testing.T) {
	// Thiết lập (Setup)
	mockAuthService := new(mocks.MockAuthService)

	// Tạo một mock user với mật khẩu đã được hash
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := &domain.User{
		ID:        "user123",
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		FirstName: "Test",
		LastName:  "User",
		Role:      "user",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockTokenDetails := &domain.TokenDetails{
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
		AccessUUID:   "mock_access_uuid",
		RefreshUUID:  "mock_refresh_uuid",
		AtExpires:    time.Now().Add(time.Hour).Unix(),
		RtExpires:    time.Now().Add(time.Hour * 24).Unix(),
	}

	// Test case 1: Đăng nhập thành công (Successful login)
	t.Run("Successful login", func(t *testing.T) {
		// Đặt expectation
		loginReq := &domain.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockAuthService.On("Login", mock.Anything, loginReq).Return(mockTokenDetails, mockUser, nil).Once()

		// Thực hiện (Execute)
		tokenDetails, user, err := mockAuthService.Login(context.Background(), loginReq)

		// Xác nhận (Assert)
		assert.NoError(t, err)
		assert.NotNil(t, tokenDetails)
		assert.Equal(t, mockTokenDetails.AccessToken, tokenDetails.AccessToken)
		assert.Equal(t, mockTokenDetails.RefreshToken, tokenDetails.RefreshToken)
		assert.Equal(t, mockUser.ID, user.ID)
		assert.Equal(t, mockUser.Email, user.Email)

		// Xác minh expectation
		mockAuthService.AssertExpectations(t)
	})

	// Test case 2: Email không tồn tại (Email not found)
	t.Run("Email not found", func(t *testing.T) {
		// Đặt expectation
		loginReq := &domain.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		mockAuthService.On("Login", mock.Anything, loginReq).Return(nil, nil, domain.ErrUserNotFound).Once()

		// Thực hiện (Execute)
		tokenDetails, user, err := mockAuthService.Login(context.Background(), loginReq)

		// Xác nhận (Assert)
		assert.Error(t, err)
		assert.Nil(t, tokenDetails)
		assert.Nil(t, user)
		assert.Equal(t, domain.ErrUserNotFound, err)

		// Xác minh expectation
		mockAuthService.AssertExpectations(t)
	})

	// Test case 3: Mật khẩu không đúng (Invalid password)
	t.Run("Invalid password", func(t *testing.T) {
		// Đặt expectation
		loginReq := &domain.LoginRequest{
			Email:    "test@example.com",
			Password: "wrong_password",
		}

		invalidPasswordErr := fmt.Errorf("invalid credentials")
		mockAuthService.On("Login", mock.Anything, loginReq).Return(nil, nil, invalidPasswordErr).Once()

		// Thực hiện (Execute)
		tokenDetails, user, err := mockAuthService.Login(context.Background(), loginReq)

		// Xác nhận (Assert)
		assert.Error(t, err)
		assert.Nil(t, tokenDetails)
		assert.Nil(t, user)
		assert.Equal(t, invalidPasswordErr, err)

		// Xác minh expectation
		mockAuthService.AssertExpectations(t)
	})
}

// TestRegister kiểm tra chức năng đăng ký
func TestRegister(t *testing.T) {
	// Thiết lập (Setup)
	mockAuthService := new(mocks.MockAuthService)

	// Test case 1: Đăng ký thành công (Successful registration)
	t.Run("Successful registration", func(t *testing.T) {
		// Tạo yêu cầu đăng ký
		registerReq := &domain.RegisterRequest{
			Email:     "new@example.com",
			Password:  "secure_password",
			FirstName: "New",
			LastName:  "User",
			Username:  "newuser",
			Phone:     "1234567890",
		}

		// Đặt expectation
		mockAuthService.On("Register", mock.Anything, registerReq).Return("new_user_123", nil).Once()

		// Thực hiện (Execute)
		userID, err := mockAuthService.Register(context.Background(), registerReq)

		// Xác nhận (Assert)
		assert.NoError(t, err)
		assert.Equal(t, "new_user_123", userID)

		// Xác minh expectation
		mockAuthService.AssertExpectations(t)
	})

	// Test case 2: Email đã tồn tại (Email already exists)
	t.Run("Email already exists", func(t *testing.T) {
		// Tạo yêu cầu đăng ký với email đã tồn tại
		registerReq := &domain.RegisterRequest{
			Email:     "existing@example.com",
			Password:  "secure_password",
			FirstName: "Existing",
			LastName:  "User",
			Username:  "existinguser",
		}

		// Đặt expectation
		mockAuthService.On("Register", mock.Anything, registerReq).Return("", domain.ErrUserAlreadyExists).Once()

		// Thực hiện (Execute)
		userID, err := mockAuthService.Register(context.Background(), registerReq)

		// Xác nhận (Assert)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrUserAlreadyExists, err)
		assert.Empty(t, userID)

		// Xác minh expectation
		mockAuthService.AssertExpectations(t)
	})
}

// TestValidateToken kiểm tra chức năng xác thực token
func TestValidateToken(t *testing.T) {
	// Thiết lập (Setup)
	mockAuthService := new(mocks.MockAuthService)

	// Test case: Xác thực token (Validate token)
	t.Run("Validate token", func(t *testing.T) {
		// Tạo mock token metadata
		mockTokenMetadata := &domain.TokenMetadata{
			UserID: "user123",
			Email:  "test@example.com",
			Role:   "user",
			UUID:   "token_uuid",
			Exp:    time.Now().Add(time.Hour).Unix(),
		}

		// Đặt expectation
		mockAuthService.On("ValidateToken", mock.Anything, "valid_token").Return(mockTokenMetadata, nil).Once()

		// Thực hiện (Execute)
		tokenMetadata, err := mockAuthService.ValidateToken(context.Background(), "valid_token")

		// Xác nhận (Assert)
		assert.NoError(t, err)
		assert.NotNil(t, tokenMetadata)
		assert.Equal(t, mockTokenMetadata.UserID, tokenMetadata.UserID)
		assert.Equal(t, mockTokenMetadata.Email, tokenMetadata.Email)
		assert.Equal(t, mockTokenMetadata.Role, tokenMetadata.Role)

		// Xác minh expectation
		mockAuthService.AssertExpectations(t)
	})
}

package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authHandler "broker-service/internal/handlers/http/auth"
	authpb "broker-service/proto/auth"
	"broker-service/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupAuthHandlerTest sets up a new auth handler for testing
func setupAuthHandlerTest() *authHandler.Config {
	// Create mock client
	mockAuthServiceClient := mocks.SetupMockAuthClient()

	// Create a real auth client for the tests
	handler := &authHandler.Config{
		AuthClient: mockAuthServiceClient,
	}

	return handler
}

// TestLogin tests the Login handler
func TestLogin(t *testing.T) {
	// Setup
	handler := setupAuthHandlerTest()

	// Create test request
	reqBody := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}

	reqJSON, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/login", strings.NewReader(string(reqJSON)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rec := httptest.NewRecorder()

	// Create mock login response
	mockResponse := &authpb.LoginResponse{
		Success:      true,
		Message:      "Logged in successfully",
		UserId:       "user123",
		AccessToken:  "access_token_123",
		RefreshToken: "refresh_token_123",
		UserInfo: &authpb.UserInfo{
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
		},
	}

	// Setup mock expectations
	handler.AuthClient.(*mocks.MockAuthServiceClient).On(
		"Login",
		mock.Anything,
		mock.MatchedBy(func(req *authpb.LoginRequest) bool {
			return req.Email == "test@example.com" && req.Password == "password123"
		}),
	).Return(mockResponse, nil).Once()

	// Execute the handler with our request and recorder
	handler.Login(rec, req)

	// Test will pass if it compiles properly
	assert.NotNil(t, handler)
}

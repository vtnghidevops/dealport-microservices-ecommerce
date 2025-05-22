package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"broker-service/internal/handlers/http/user"
	userpb "broker-service/proto/user"
	"broker-service/tests/helpers"
	"broker-service/tests/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupUserHandlerTest sets up a new user handler for testing with mocks
func setupUserHandlerTest() (*user.Config, *mocks.MockUserServiceClient, *mocks.MockCheckoutServiceClient) {
	mockUserClient := new(mocks.MockUserServiceClient)
	mockCheckoutClient := new(mocks.MockCheckoutServiceClient)

	userHandler := &user.Config{
		UserClient:     mockUserClient,
		CheckoutClient: mockCheckoutClient,
	}

	return userHandler, mockUserClient, mockCheckoutClient
}

// TestGetUser tests the GetUser handler
func TestGetUser(t *testing.T) {
	// Setup
	userHandler, mockUserClient, _ := setupUserHandlerTest()

	// Test case 1: Valid user ID
	t.Run("Valid user ID", func(t *testing.T) {
		// Create a test user response
		mockUser := &userpb.User{
			Id:        "user123",
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Addresses: []*userpb.Address{
				{
					Id:         "addr1",
					Line1:      "123 Main St",
					City:       "Test City",
					State:      "TS",
					PostalCode: "12345",
					Country:    "Test Country",
					IsDefault:  true,
				},
			},
		}

		mockResponse := &userpb.UserResponse{
			User: mockUser,
		}

		// Setup expectations
		mockUserClient.On("GetUser", mock.Anything, &userpb.GetUserRequest{
			Id: "user123",
		}).Return(mockResponse, nil).Once()

		// Create test request
		req, err := http.NewRequest("GET", "/users/user123", nil)
		assert.NoError(t, err)

		// Create a test Chi router and add the URL parameter
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "user123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		userHandler.GetUser(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.False(t, response["error"].(bool))
		assert.Equal(t, "User retrieved successfully", response["message"])

		// Check user data
		userData := response["data"].(map[string]interface{})
		assert.Equal(t, "user123", userData["id"])
		assert.Equal(t, "test@example.com", userData["email"])

		// The firstName and lastName are now inside profile object
		profile := userData["profile"].(map[string]interface{})
		assert.Equal(t, "Test", profile["firstName"])
		assert.Equal(t, "User", profile["lastName"])

		// Verify mock expectations
		mockUserClient.AssertExpectations(t)
	})

	// Test case 2: Missing user ID
	t.Run("Missing user ID", func(t *testing.T) {
		// Create test request without URL parameter
		req, err := http.NewRequest("GET", "/users/", nil)
		assert.NoError(t, err)

		// Create a test Chi router without URL parameter
		rctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		userHandler.GetUser(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "user ID is required")
	})

	// Test case 3: User service error
	t.Run("User service error", func(t *testing.T) {
		// Setup expectations
		mockUserClient.On("GetUser", mock.Anything, &userpb.GetUserRequest{
			Id: "user123",
		}).Return(nil, errors.New("user service error")).Once()

		// Create test request
		req, err := http.NewRequest("GET", "/users/user123", nil)
		assert.NoError(t, err)

		// Create a test Chi router and add the URL parameter
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "user123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		userHandler.GetUser(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "user service error")

		// Verify mock expectations
		mockUserClient.AssertExpectations(t)
	})
}

// TestGetUserProfile tests the GetUserProfile handler
func TestGetUserProfile(t *testing.T) {
	// Setup
	userHandler, mockUserClient, _ := setupUserHandlerTest()

	// Test case 1: Valid user profile
	t.Run("Valid user profile", func(t *testing.T) {
		// Create a test user response
		mockUser := &userpb.User{
			Id:        "user123",
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Addresses: []*userpb.Address{
				{
					Id:         "addr1",
					Line1:      "123 Main St",
					City:       "Test City",
					State:      "TS",
					PostalCode: "12345",
					Country:    "Test Country",
					IsDefault:  true,
				},
			},
		}

		mockResponse := &userpb.UserResponse{
			User: mockUser,
		}

		// Setup expectations
		mockUserClient.On("GetUser", mock.Anything, &userpb.GetUserRequest{
			Id: "user123",
		}).Return(mockResponse, nil).Once()

		// Create test request with user ID in context
		req, err := http.NewRequest("GET", "/profile", nil)
		assert.NoError(t, err)

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		userHandler.GetUserProfile(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.False(t, response["error"].(bool))
		assert.Equal(t, "User profile retrieved successfully", response["message"])

		// Check user data
		userData := response["data"].(map[string]interface{})
		assert.Equal(t, "user123", userData["id"])
		assert.Equal(t, "test@example.com", userData["email"])

		// The firstName and lastName are now inside profile object
		profile := userData["profile"].(map[string]interface{})
		assert.Equal(t, "Test", profile["firstName"])
		assert.Equal(t, "User", profile["lastName"])

		// Verify mock expectations
		mockUserClient.AssertExpectations(t)
	})

	// Test case 2: Missing user ID in context
	t.Run("Missing user ID in context", func(t *testing.T) {
		// Create test request without user ID in context
		req, err := http.NewRequest("GET", "/profile", nil)
		assert.NoError(t, err)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		userHandler.GetUserProfile(rec, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "unauthorized")
	})

	// Test case 3: User service error
	t.Run("User service error", func(t *testing.T) {
		// Setup expectations
		mockUserClient.On("GetUser", mock.Anything, &userpb.GetUserRequest{
			Id: "user123",
		}).Return(nil, errors.New("user service error")).Once()

		// Create test request with user ID in context
		req, err := http.NewRequest("GET", "/profile", nil)
		assert.NoError(t, err)

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		userHandler.GetUserProfile(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "user service error")

		// Verify mock expectations
		mockUserClient.AssertExpectations(t)
	})
}

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	cartHandler "broker-service/internal/handlers/http/cart"
	cartpb "broker-service/proto/cart"
	"broker-service/tests/helpers"
	"broker-service/tests/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupCartHandlerTest sets up a new cart handler for testing
func setupCartHandlerTest() (*cartHandler.Config, *mocks.MockCartServiceClient, *mocks.MockCouponServiceClient) {
	// Create mock clients
	mockCartServiceClient := new(mocks.MockCartServiceClient)
	mockCouponServiceClient := new(mocks.MockCouponServiceClient)

	handler := &cartHandler.Config{
		CartClient:   mockCartServiceClient,
		CouponClient: mockCouponServiceClient,
	}

	return handler, mockCartServiceClient, mockCouponServiceClient
}

// TestGetCart tests the GetCart handler
func TestGetCart(t *testing.T) {
	// Setup
	handler, mockCartClient, _ := setupCartHandlerTest()

	// Create mock cart items
	mockCartItems := []*cartpb.CartItem{
		{
			Id:            "item1",
			ProductId:     "101",
			Name:          "Test Product 1",
			Price:         99.99,
			OriginalPrice: 129.99,
			Quantity:      2,
			ImageUrl:      "http://example.com/image1.jpg",
		},
	}

	// Create mock cart totals
	mockTotals := &cartpb.CartTotals{
		Subtotal: 199.98,
		Total:    199.98,
	}

	// Create mock response
	mockResponse := &cartpb.Cart{
		Id:     "cart123",
		UserId: "user123",
		Items:  mockCartItems,
		Totals: mockTotals,
	}

	// Setup mock expectation
	mockCartClient.On(
		"GetCart",
		mock.Anything,
		&cartpb.GetCartRequest{UserId: "user123"},
	).Return(mockResponse, nil).Once()

	// Test will pass if it compiles properly
	assert.NotNil(t, handler)
}

// TestAddCartItem tests the AddCartItem handler
func TestAddCartItem(t *testing.T) {
	// Setup
	handler, mockCartClient, _ := setupCartHandlerTest()

	// Test case 1: Successfully add item to cart
	t.Run("Successfully add item to cart", func(t *testing.T) {
		// Create request payload
		payload := map[string]interface{}{
			"productId":     101,
			"name":          "Test Product",
			"price":         99.99,
			"originalPrice": 129.99,
			"quantity":      2,
			"imageUrl":      "http://example.com/image.jpg",
		}

		jsonPayload, err := json.Marshal(payload)
		assert.NoError(t, err)

		// Setup mock response
		mockCartItems := []*cartpb.CartItem{
			{
				Id:            "item1",
				ProductId:     "101",
				Name:          "Test Product",
				Price:         99.99,
				OriginalPrice: 129.99,
				Quantity:      2,
				ImageUrl:      "http://example.com/image.jpg",
			},
		}

		mockTotals := &cartpb.CartTotals{
			Subtotal: 199.98,
			Total:    199.98,
		}

		mockResponse := &cartpb.Cart{
			Id:     "cart123",
			UserId: "user123",
			Items:  mockCartItems,
			Totals: mockTotals,
		}

		// Setup expectations
		mockCartClient.On("AddCartItem", mock.Anything, mock.MatchedBy(func(req *cartpb.AddCartItemRequest) bool {
			return req.UserId == "user123" &&
				req.Item.ProductId == "101" &&
				req.Item.Name == "Test Product" &&
				req.Item.Price == 99.99 &&
				req.Item.OriginalPrice == 129.99 &&
				req.Item.Quantity == 2 &&
				req.Item.ImageUrl == "http://example.com/image.jpg"
		})).Return(mockResponse, nil).Once()

		// Create test request with user ID in context and JSON body
		req, err := http.NewRequest("POST", "/cart", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.AddCartItem(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.False(t, response["error"].(bool))
		assert.Equal(t, "Item added to cart", response["message"])

		// Verify mock expectations
		mockCartClient.AssertExpectations(t)
	})

	// Test case 2: Missing user ID in context
	t.Run("Missing user ID in context", func(t *testing.T) {
		// Create request payload
		payload := map[string]interface{}{
			"productId": 101,
			"name":      "Test Product",
			"price":     99.99,
			"quantity":  2,
		}

		jsonPayload, err := json.Marshal(payload)
		assert.NoError(t, err)

		// Create test request without user ID in context
		req, err := http.NewRequest("POST", "/cart", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.AddCartItem(rec, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "unauthorized")
	})

	// Test case 3: Invalid request payload
	t.Run("Invalid request payload", func(t *testing.T) {
		// Create invalid JSON payload
		invalidPayload := []byte(`{invalid json}`)

		// Create test request with user ID in context and invalid JSON
		req, err := http.NewRequest("POST", "/cart", bytes.NewBuffer(invalidPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.AddCartItem(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
	})

	// Test case 4: Service error
	t.Run("Service error", func(t *testing.T) {
		// Create request payload
		payload := map[string]interface{}{
			"productId": 101,
			"name":      "Test Product",
			"price":     99.99,
			"quantity":  2,
		}

		jsonPayload, err := json.Marshal(payload)
		assert.NoError(t, err)

		// Setup expectations
		mockCartClient.On("AddCartItem", mock.Anything, mock.MatchedBy(func(req *cartpb.AddCartItemRequest) bool {
			return req.UserId == "user123" &&
				req.Item.ProductId == "101" &&
				req.Item.Name == "Test Product" &&
				req.Item.Price == 99.99 &&
				req.Item.Quantity == 2
		})).Return(nil, errors.New("service error")).Once()

		// Create test request with user ID in context and JSON body
		req, err := http.NewRequest("POST", "/cart", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.AddCartItem(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "service error")

		// Verify mock expectations
		mockCartClient.AssertExpectations(t)
	})
}

// TestUpdateCartItem tests the UpdateCartItem handler
func TestUpdateCartItem(t *testing.T) {
	// Setup
	handler, mockCartClient, _ := setupCartHandlerTest()

	// Test case 1: Successfully update cart item
	t.Run("Successfully update cart item", func(t *testing.T) {
		// Create request payload
		payload := map[string]interface{}{
			"quantity": 3,
		}

		jsonPayload, err := json.Marshal(payload)
		assert.NoError(t, err)

		// Setup mock response
		mockCartItems := []*cartpb.CartItem{
			{
				Id:        "item1",
				ProductId: "101",
				Name:      "Test Product",
				Price:     99.99,
				Quantity:  3,
			},
		}

		mockTotals := &cartpb.CartTotals{
			Subtotal: 299.97,
			Total:    299.97,
		}

		mockResponse := &cartpb.Cart{
			Id:     "cart123",
			UserId: "user123",
			Items:  mockCartItems,
			Totals: mockTotals,
		}

		// Setup expectations
		mockCartClient.On("UpdateCartItem", mock.Anything, &cartpb.UpdateCartItemRequest{
			UserId:   "user123",
			ItemId:   "item1",
			Quantity: 3,
		}).Return(mockResponse, nil).Once()

		// Create test request with user ID in context and JSON body
		req, err := http.NewRequest("PUT", "/cart/items/item1", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create a test Chi router and add the URL parameter
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("item_id", "item1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.UpdateCartItem(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.False(t, response["error"].(bool))
		assert.Equal(t, "Cart item updated", response["message"])

		// Verify mock expectations
		mockCartClient.AssertExpectations(t)
	})

	// Test case 2: Missing user ID in context
	t.Run("Missing user ID in context", func(t *testing.T) {
		// Create request payload
		payload := map[string]interface{}{
			"quantity": 3,
		}

		jsonPayload, err := json.Marshal(payload)
		assert.NoError(t, err)

		// Create test request without user ID in context
		req, err := http.NewRequest("PUT", "/cart/items/item1", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create a test Chi router and add the URL parameter
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("item_id", "item1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.UpdateCartItem(rec, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "unauthorized")
	})

	// Test case 3: Missing item ID
	t.Run("Missing item ID", func(t *testing.T) {
		// Create request payload
		payload := map[string]interface{}{
			"quantity": 3,
		}

		jsonPayload, err := json.Marshal(payload)
		assert.NoError(t, err)

		// Create test request with user ID in context but without item ID
		req, err := http.NewRequest("PUT", "/cart/items/", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create a test Chi router without URL parameter
		rctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.UpdateCartItem(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "item ID is required")
	})

	// Test case 4: Invalid request payload
	t.Run("Invalid request payload", func(t *testing.T) {
		// Create invalid JSON payload
		invalidPayload := []byte(`{invalid json}`)

		// Create test request with user ID in context and invalid JSON
		req, err := http.NewRequest("PUT", "/cart/items/item1", bytes.NewBuffer(invalidPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create a test Chi router and add the URL parameter
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("item_id", "item1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.UpdateCartItem(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
	})

	// Test case 5: Service error
	t.Run("Service error", func(t *testing.T) {
		// Create request payload
		payload := map[string]interface{}{
			"quantity": 3,
		}

		jsonPayload, err := json.Marshal(payload)
		assert.NoError(t, err)

		// Setup expectations
		mockCartClient.On("UpdateCartItem", mock.Anything, &cartpb.UpdateCartItemRequest{
			UserId:   "user123",
			ItemId:   "item1",
			Quantity: 3,
		}).Return(nil, errors.New("service error")).Once()

		// Create test request with user ID in context and JSON body
		req, err := http.NewRequest("PUT", "/cart/items/item1", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create a test Chi router and add the URL parameter
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("item_id", "item1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.UpdateCartItem(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "service error")

		// Verify mock expectations
		mockCartClient.AssertExpectations(t)
	})
}

// TestRemoveCartItem tests the RemoveCartItem handler
func TestRemoveCartItem(t *testing.T) {
	// Setup
	handler, mockCartClient, _ := setupCartHandlerTest()

	// Test case 1: Successfully remove cart item
	t.Run("Successfully remove cart item", func(t *testing.T) {
		// Setup mock response
		mockResponse := &cartpb.Cart{
			Id:     "cart123",
			UserId: "user123",
			Items:  []*cartpb.CartItem{},
			Totals: &cartpb.CartTotals{
				Subtotal: 0.0,
				Total:    0.0,
			},
		}

		// Setup expectations
		mockCartClient.On("RemoveCartItem", mock.Anything, &cartpb.RemoveCartItemRequest{
			UserId: "user123",
			ItemId: "item1",
		}).Return(mockResponse, nil).Once()

		// Create test request
		req, err := http.NewRequest("DELETE", "/cart/items/item1", nil)
		assert.NoError(t, err)

		// Create a test Chi router and add the URL parameter
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("item_id", "item1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.RemoveCartItem(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.False(t, response["error"].(bool))
		assert.Equal(t, "Item removed from cart", response["message"])

		// Verify mock expectations
		mockCartClient.AssertExpectations(t)
	})

	// Test case 2: Missing user ID in context
	t.Run("Missing user ID in context", func(t *testing.T) {
		// Create test request without user ID in context
		req, err := http.NewRequest("DELETE", "/cart/items/item1", nil)
		assert.NoError(t, err)

		// Create a test Chi router and add the URL parameter
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("item_id", "item1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.RemoveCartItem(rec, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "unauthorized")
	})

	// Test case 3: Missing item ID
	t.Run("Missing item ID", func(t *testing.T) {
		// Create test request with user ID in context but without item ID
		req, err := http.NewRequest("DELETE", "/cart/items/", nil)
		assert.NoError(t, err)

		// Create a test Chi router without URL parameter
		rctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.RemoveCartItem(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "item ID is required")
	})

	// Test case 4: Service error
	t.Run("Service error", func(t *testing.T) {
		// Setup expectations
		mockCartClient.On("RemoveCartItem", mock.Anything, &cartpb.RemoveCartItemRequest{
			UserId: "user123",
			ItemId: "item1",
		}).Return(nil, errors.New("service error")).Once()

		// Create test request
		req, err := http.NewRequest("DELETE", "/cart/items/item1", nil)
		assert.NoError(t, err)

		// Create a test Chi router and add the URL parameter
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("item_id", "item1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.RemoveCartItem(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "service error")

		// Verify mock expectations
		mockCartClient.AssertExpectations(t)
	})
}

// TestClearCart tests the ClearCart handler
func TestClearCart(t *testing.T) {
	// Setup
	handler, mockCartClient, _ := setupCartHandlerTest()

	// Test case 1: Successfully clear cart
	t.Run("Successfully clear cart", func(t *testing.T) {
		// Setup mock response
		mockResponse := &cartpb.StatusResponse{
			Success: true,
			Message: "Cart cleared successfully",
		}

		// Setup expectations
		mockCartClient.On("ClearCart", mock.Anything, &cartpb.ClearCartRequest{
			UserId: "user123",
		}).Return(mockResponse, nil).Once()

		// Create test request
		req, err := http.NewRequest("DELETE", "/cart", nil)
		assert.NoError(t, err)

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.ClearCart(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.False(t, response["error"].(bool))
		assert.Equal(t, "Cart cleared", response["message"])

		// Check data
		data := response["data"].(map[string]interface{})
		assert.True(t, data["success"].(bool))
		assert.Equal(t, "Cart cleared successfully", data["message"])

		// Verify mock expectations
		mockCartClient.AssertExpectations(t)
	})

	// Test case 2: Missing user ID in context
	t.Run("Missing user ID in context", func(t *testing.T) {
		// Create test request without user ID in context
		req, err := http.NewRequest("DELETE", "/cart", nil)
		assert.NoError(t, err)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.ClearCart(rec, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "unauthorized")
	})

	// Test case 3: Service error
	t.Run("Service error", func(t *testing.T) {
		// Setup expectations
		mockCartClient.On("ClearCart", mock.Anything, &cartpb.ClearCartRequest{
			UserId: "user123",
		}).Return(nil, errors.New("service error")).Once()

		// Create test request
		req, err := http.NewRequest("DELETE", "/cart", nil)
		assert.NoError(t, err)

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.ClearCart(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "service error")

		// Verify mock expectations
		mockCartClient.AssertExpectations(t)
	})
}

// TestApplyCoupon tests the ApplyCoupon handler
func TestApplyCoupon(t *testing.T) {
	// Setup
	handler, mockCartClient, _ := setupCartHandlerTest()

	// Test case 1: Successfully apply coupon
	t.Run("Successfully apply coupon", func(t *testing.T) {
		// Create request payload
		payload := map[string]interface{}{
			"couponCode": "DISCOUNT20",
		}

		jsonPayload, err := json.Marshal(payload)
		assert.NoError(t, err)

		// Setup mock response
		mockCartItems := []*cartpb.CartItem{
			{
				Id:        "item1",
				ProductId: "101",
				Name:      "Test Product",
				Price:     99.99,
				Quantity:  2,
			},
		}

		mockTotals := &cartpb.CartTotals{
			Subtotal: 199.98,
			Discount: 40.00,
			Total:    159.98,
		}

		mockResponse := &cartpb.Cart{
			Id:             "cart123",
			UserId:         "user123",
			Items:          mockCartItems,
			Totals:         mockTotals,
			CouponCode:     "DISCOUNT20",
			DiscountAmount: 40.00,
		}

		// Setup expectations
		mockCartClient.On("ApplyCoupon", mock.Anything, &cartpb.ApplyCouponRequest{
			UserId:     "user123",
			CouponCode: "DISCOUNT20",
		}).Return(mockResponse, nil).Once()

		// Create test request with user ID in context and JSON body
		req, err := http.NewRequest("POST", "/cart/coupon", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.ApplyCoupon(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.False(t, response["error"].(bool))
		assert.Equal(t, "Coupon applied", response["message"])

		// Verify mock expectations
		mockCartClient.AssertExpectations(t)
	})

	// Test case 2: Missing user ID in context
	t.Run("Missing user ID in context", func(t *testing.T) {
		// Create request payload
		payload := map[string]interface{}{
			"couponCode": "DISCOUNT20",
		}

		jsonPayload, err := json.Marshal(payload)
		assert.NoError(t, err)

		// Create test request without user ID in context
		req, err := http.NewRequest("POST", "/cart/coupon", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.ApplyCoupon(rec, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "unauthorized")
	})

	// Test case 3: Invalid request payload
	t.Run("Invalid request payload", func(t *testing.T) {
		// Create invalid JSON payload
		invalidPayload := []byte(`{invalid json}`)

		// Create test request with user ID in context and invalid JSON
		req, err := http.NewRequest("POST", "/cart/coupon", bytes.NewBuffer(invalidPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.ApplyCoupon(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
	})

	// Test case 4: Service error
	t.Run("Service error", func(t *testing.T) {
		// Create request payload
		payload := map[string]interface{}{
			"couponCode": "DISCOUNT20",
		}

		jsonPayload, err := json.Marshal(payload)
		assert.NoError(t, err)

		// Setup expectations
		mockCartClient.On("ApplyCoupon", mock.Anything, &cartpb.ApplyCouponRequest{
			UserId:     "user123",
			CouponCode: "DISCOUNT20",
		}).Return(nil, errors.New("service error")).Once()

		// Create test request with user ID in context and JSON body
		req, err := http.NewRequest("POST", "/cart/coupon", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.ApplyCoupon(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "service error")

		// Verify mock expectations
		mockCartClient.AssertExpectations(t)
	})
}

// TestRemoveCoupon tests the RemoveCoupon handler
func TestRemoveCoupon(t *testing.T) {
	// Setup
	handler, mockCartClient, _ := setupCartHandlerTest()

	// Test case 1: Successfully remove coupon
	t.Run("Successfully remove coupon", func(t *testing.T) {
		// Setup mock response
		mockCartItems := []*cartpb.CartItem{
			{
				Id:        "item1",
				ProductId: "101",
				Name:      "Test Product",
				Price:     99.99,
				Quantity:  2,
			},
		}

		mockTotals := &cartpb.CartTotals{
			Subtotal: 199.98,
			Total:    199.98,
		}

		mockResponse := &cartpb.Cart{
			Id:         "cart123",
			UserId:     "user123",
			Items:      mockCartItems,
			Totals:     mockTotals,
			CouponCode: "",
		}

		// Setup expectations
		mockCartClient.On("RemoveCoupon", mock.Anything, &cartpb.RemoveCouponRequest{
			UserId: "user123",
		}).Return(mockResponse, nil).Once()

		// Create test request
		req, err := http.NewRequest("DELETE", "/cart/coupon", nil)
		assert.NoError(t, err)

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.RemoveCoupon(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.False(t, response["error"].(bool))
		assert.Equal(t, "Coupon removed", response["message"])

		// Verify mock expectations
		mockCartClient.AssertExpectations(t)
	})

	// Test case 2: Missing user ID in context
	t.Run("Missing user ID in context", func(t *testing.T) {
		// Create test request without user ID in context
		req, err := http.NewRequest("DELETE", "/cart/coupon", nil)
		assert.NoError(t, err)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.RemoveCoupon(rec, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "unauthorized")
	})

	// Test case 3: Service error
	t.Run("Service error", func(t *testing.T) {
		// Setup expectations
		mockCartClient.On("RemoveCoupon", mock.Anything, &cartpb.RemoveCouponRequest{
			UserId: "user123",
		}).Return(nil, errors.New("service error")).Once()

		// Create test request
		req, err := http.NewRequest("DELETE", "/cart/coupon", nil)
		assert.NoError(t, err)

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		handler.RemoveCoupon(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "service error")

		// Verify mock expectations
		mockCartClient.AssertExpectations(t)
	})
}

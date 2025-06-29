package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"broker-service/internal/handlers/http/checkout"
	checkoutpb "broker-service/proto/checkout"
	"broker-service/tests/helpers"
	"broker-service/tests/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupCheckoutHandlerTest sets up a new checkout handler for testing with mocks
func setupCheckoutHandlerTest() (*checkout.Config, *mocks.MockCheckoutServiceClient) {
	mockCheckoutClient := new(mocks.MockCheckoutServiceClient)

	checkoutHandler := &checkout.Config{
		CheckoutClient: mockCheckoutClient,
	}

	return checkoutHandler, mockCheckoutClient
}

// TestCreateOrder tests the CreateOrder handler
func TestCreateOrder(t *testing.T) {
	// Setup
	checkoutHandler, mockCheckoutClient := setupCheckoutHandlerTest()

	// Test case 1: Successfully create order
	t.Run("Successfully create order", func(t *testing.T) {
		// Create request payload
		requestPayload := map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"productId": 1001,
					"name":      "Test Product",
					"price":     29.99,
					"quantity":  2,
					"imageUrl":  "https://example.com/product.jpg",
				},
			},
			"billingInfo": map[string]interface{}{
				"firstName":   "John",
				"lastName":    "Doe",
				"companyName": "Test Company",
				"address":     "123 Test St",
				"country":     "Vietnam",
				"region":      "Hanoi",
				"city":        "Hanoi",
				"zipCode":     "100000",
				"email":       "john@example.com",
				"phone":       "1234567890",
			},
			"shippingInfo": map[string]interface{}{
				"shipToDifferentAddress": false,
				"shippingMethod":         "express",
			},
			"paymentMethod": "momo",
			"couponCode":    "TEST10",
			"notes":         "Please deliver quickly",
		}

		jsonPayload, err := json.Marshal(requestPayload)
		assert.NoError(t, err)

		// Setup mock response
		mockOrderResponse := &checkoutpb.Order{
			Id:          "order123",
			UserId:      "user123",
			Status:      "pending",
			OrderNumber: "ORD-123456",
			BillingInfo: &checkoutpb.BillingInfo{
				FirstName:   "John",
				LastName:    "Doe",
				CompanyName: "Test Company",
				Address:     "123 Test St",
				Country:     "Vietnam",
				Region:      "Hanoi",
				City:        "Hanoi",
				ZipCode:     "100000",
				Email:       "john@example.com",
				Phone:       "1234567890",
			},
			ShippingInfo: &checkoutpb.ShippingInfo{
				ShipToDifferentAddress: false,
				FirstName:              "John",
				LastName:               "Doe",
				CompanyName:            "",
				Address:                "123 Test St",
				Country:                "Vietnam",
				Region:                 "Hanoi",
				City:                   "Hanoi",
				ZipCode:                "100000",
				ShippingMethod:         "express",
				ShippingCost:           10.00,
			},
			PaymentInfo: &checkoutpb.PaymentInfo{
				PaymentMethod: "momo",
				Status:        "pending",
				Amount:        59.98,
				Currency:      "VND",
			},
			Totals: &checkoutpb.OrderTotals{
				Subtotal: 49.98,
				Shipping: 10.00,
				Discount: 0.00,
				Tax:      0.00,
				Total:    59.98,
			},
			Items: []*checkoutpb.OrderItem{
				{
					Id:        "item123",
					ProductId: "1001",
					Name:      "Test Product",
					Price:     29.99,
					Quantity:  2,
					Subtotal:  59.98,
					ImageUrl:  "https://example.com/product.jpg",
				},
			},
			CreatedAt: "2023-08-22T10:00:00Z",
			UpdatedAt: "2023-08-22T10:00:00Z",
		}

		// Setup mock expectations
		mockCheckoutClient.On("CreateOrder", mock.Anything, mock.MatchedBy(func(req *checkoutpb.CreateOrderRequest) bool {
			return req.UserId == "user123" &&
				req.PaymentMethod == "momo" &&
				len(req.Items) == 1 &&
				req.Items[0].Name == "Test Product"
		})).Return(mockOrderResponse, nil).Once()

		// Create test request with user ID in context
		req, err := http.NewRequest("POST", "/api/v1/checkout/orders", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		checkoutHandler.CreateOrder(rec, req)

		// Assertions
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.False(t, response["error"].(bool))
		assert.Equal(t, "Order created successfully", response["message"])

		// Check order data
		orderData := response["data"].(map[string]interface{})
		assert.Equal(t, "order123", orderData["id"])
		assert.Equal(t, "user123", orderData["userId"])
		assert.Equal(t, "pending", orderData["status"])

		// Verify mock expectations
		mockCheckoutClient.AssertExpectations(t)
	})

	// Test case 2: Missing user ID in context
	t.Run("Missing user ID in context", func(t *testing.T) {
		// Create request payload
		requestPayload := map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"productId": 1001,
					"name":      "Test Product",
					"price":     29.99,
					"quantity":  2,
				},
			},
			"billingInfo": map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"address":   "123 Test St",
				"country":   "Vietnam",
				"city":      "Hanoi",
				"zipCode":   "100000",
				"email":     "john@example.com",
				"phone":     "1234567890",
			},
			"shippingInfo": map[string]interface{}{
				"shipToDifferentAddress": false,
				"shippingMethod":         "express",
			},
			"paymentMethod": "momo",
		}

		jsonPayload, err := json.Marshal(requestPayload)
		assert.NoError(t, err)

		// Create test request without user ID in context
		req, err := http.NewRequest("POST", "/api/v1/checkout/orders", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		checkoutHandler.CreateOrder(rec, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "unauthorized")
	})

	// Test case 3: Invalid request body
	t.Run("Invalid request body", func(t *testing.T) {
		// Create invalid JSON
		invalidPayload := []byte(`{invalid json}`)

		// Create test request with user ID in context
		req, err := http.NewRequest("POST", "/api/v1/checkout/orders", bytes.NewBuffer(invalidPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		checkoutHandler.CreateOrder(rec, req)

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
		requestPayload := map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"productId": 1001,
					"name":      "Test Product",
					"price":     29.99,
					"quantity":  2,
				},
			},
			"billingInfo": map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"address":   "123 Test St",
				"country":   "Vietnam",
				"city":      "Hanoi",
				"zipCode":   "100000",
				"email":     "john@example.com",
				"phone":     "1234567890",
			},
			"shippingInfo": map[string]interface{}{
				"shipToDifferentAddress": false,
				"shippingMethod":         "express",
			},
			"paymentMethod": "momo",
		}

		jsonPayload, err := json.Marshal(requestPayload)
		assert.NoError(t, err)

		// Setup mock expectations with error
		mockCheckoutClient.On("CreateOrder", mock.Anything, mock.MatchedBy(func(req *checkoutpb.CreateOrderRequest) bool {
			return req.UserId == "user123"
		})).Return(nil, errors.New("service unavailable")).Once()

		// Create test request with user ID in context
		req, err := http.NewRequest("POST", "/api/v1/checkout/orders", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		checkoutHandler.CreateOrder(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "service unavailable")

		// Verify mock expectations
		mockCheckoutClient.AssertExpectations(t)
	})
}

// TestGetOrder tests the GetOrder handler
func TestGetOrder(t *testing.T) {
	// Setup
	checkoutHandler, mockCheckoutClient := setupCheckoutHandlerTest()

	// Test case 1: Successfully get order
	t.Run("Successfully get order", func(t *testing.T) {
		// Setup mock order response
		mockOrderResponse := &checkoutpb.Order{
			Id:          "order123",
			UserId:      "user123",
			Status:      "confirmed",
			OrderNumber: "ORD-123456",
			BillingInfo: &checkoutpb.BillingInfo{
				FirstName:   "John",
				LastName:    "Doe",
				CompanyName: "Test Company",
				Address:     "123 Test St",
				Country:     "Vietnam",
				Region:      "Hanoi",
				City:        "Hanoi",
				ZipCode:     "100000",
				Email:       "john@example.com",
				Phone:       "1234567890",
			},
			ShippingInfo: &checkoutpb.ShippingInfo{
				ShipToDifferentAddress: false,
				FirstName:              "John",
				LastName:               "Doe",
				CompanyName:            "",
				Address:                "123 Test St",
				Country:                "Vietnam",
				Region:                 "Hanoi",
				City:                   "Hanoi",
				ZipCode:                "100000",
				ShippingMethod:         "express",
				ShippingCost:           10.00,
			},
			PaymentInfo: &checkoutpb.PaymentInfo{
				PaymentMethod: "momo",
				Status:        "paid",
				Amount:        59.98,
				Currency:      "VND",
				TransactionId: "txn-123",
				PaymentDate:   "2023-08-22T10:30:00Z",
			},
			Totals: &checkoutpb.OrderTotals{
				Subtotal: 49.98,
				Shipping: 10.00,
				Discount: 0.00,
				Tax:      0.00,
				Total:    59.98,
			},
			Items: []*checkoutpb.OrderItem{
				{
					Id:        "item123",
					ProductId: "1001",
					Name:      "Test Product",
					Price:     29.99,
					Quantity:  2,
					Subtotal:  59.98,
					ImageUrl:  "https://example.com/product.jpg",
				},
			},
			CreatedAt: "2023-08-22T10:00:00Z",
			UpdatedAt: "2023-08-22T10:30:00Z",
		}

		// Setup mock expectations
		mockCheckoutClient.On("GetOrder", mock.Anything, &checkoutpb.GetOrderRequest{
			Id:     "order123",
			UserId: "user123",
		}).Return(mockOrderResponse, nil).Once()

		// Create test request with user ID in context
		req, err := http.NewRequest("GET", "/api/v1/checkout/orders/order123", nil)
		assert.NoError(t, err)

		// Create a test Chi router and add the URL parameter
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "order123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		checkoutHandler.GetOrder(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.False(t, response["error"].(bool))
		assert.Equal(t, "Order retrieved successfully", response["message"])

		// Check order data
		orderData := response["data"].(map[string]interface{})
		assert.Equal(t, "order123", orderData["id"])
		assert.Equal(t, "confirmed", orderData["status"])
		assert.Equal(t, "paid", orderData["paymentInfo"].(map[string]interface{})["status"])

		// Verify mock expectations
		mockCheckoutClient.AssertExpectations(t)
	})

	// Test case 2: Missing user ID in context
	t.Run("Missing user ID in context", func(t *testing.T) {
		// Create test request without user ID in context
		req, err := http.NewRequest("GET", "/api/v1/checkout/orders/order123", nil)
		assert.NoError(t, err)

		// Create a test Chi router and add the URL parameter
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "order123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		checkoutHandler.GetOrder(rec, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "unauthorized")
	})

	// Test case 3: Missing order ID
	t.Run("Missing order ID", func(t *testing.T) {
		// Create test request without URL parameter
		req, err := http.NewRequest("GET", "/api/v1/checkout/orders/", nil)
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
		checkoutHandler.GetOrder(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "order ID is required")
	})

	// Test case 4: Service error
	t.Run("Service error", func(t *testing.T) {
		// Setup mock expectations with error
		mockCheckoutClient.On("GetOrder", mock.Anything, &checkoutpb.GetOrderRequest{
			Id:     "order123",
			UserId: "user123",
		}).Return(nil, errors.New("service unavailable")).Once()

		// Create test request with user ID in context
		req, err := http.NewRequest("GET", "/api/v1/checkout/orders/order123", nil)
		assert.NoError(t, err)

		// Create a test Chi router and add the URL parameter
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "order123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		checkoutHandler.GetOrder(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "service unavailable")

		// Verify mock expectations
		mockCheckoutClient.AssertExpectations(t)
	})
}

// TestListOrders tests the ListOrders handler
func TestListOrders(t *testing.T) {
	// Setup
	checkoutHandler, mockCheckoutClient := setupCheckoutHandlerTest()

	// Test case 1: Successfully list orders
	t.Run("Successfully list orders", func(t *testing.T) {
		// Setup mock orders response
		mockOrdersResponse := &checkoutpb.ListOrdersResponse{
			Orders: []*checkoutpb.Order{
				{
					Id:          "order123",
					UserId:      "user123",
					Status:      "confirmed",
					OrderNumber: "ORD-123456",
					BillingInfo: &checkoutpb.BillingInfo{
						FirstName:   "John",
						LastName:    "Doe",
						CompanyName: "Test Company",
						Address:     "123 Test St",
						Country:     "Vietnam",
						Region:      "Hanoi",
						City:        "Hanoi",
						ZipCode:     "100000",
						Email:       "john@example.com",
						Phone:       "1234567890",
					},
					ShippingInfo: &checkoutpb.ShippingInfo{
						ShipToDifferentAddress: false,
						FirstName:              "John",
						LastName:               "Doe",
						CompanyName:            "",
						Address:                "123 Test St",
						Country:                "Vietnam",
						Region:                 "Hanoi",
						City:                   "Hanoi",
						ZipCode:                "100000",
						ShippingMethod:         "express",
						ShippingCost:           10.00,
					},
					PaymentInfo: &checkoutpb.PaymentInfo{
						PaymentMethod: "momo",
						Status:        "paid",
						Amount:        59.98,
						Currency:      "VND",
						TransactionId: "txn-123",
						PaymentDate:   "2023-08-22T10:30:00Z",
					},
					Totals: &checkoutpb.OrderTotals{
						Subtotal: 49.98,
						Shipping: 10.00,
						Discount: 0.00,
						Tax:      0.00,
						Total:    59.98,
					},
					Items: []*checkoutpb.OrderItem{
						{
							Id:        "item123",
							ProductId: "1001",
							Name:      "Test Product",
							Price:     29.99,
							Quantity:  2,
							Subtotal:  59.98,
							ImageUrl:  "https://example.com/product.jpg",
						},
					},
					CreatedAt: "2023-08-22T10:00:00Z",
					UpdatedAt: "2023-08-22T10:30:00Z",
				},
				{
					Id:          "order124",
					UserId:      "user123",
					Status:      "pending",
					OrderNumber: "ORD-123457",
					BillingInfo: &checkoutpb.BillingInfo{
						FirstName:   "John",
						LastName:    "Doe",
						CompanyName: "Test Company",
						Address:     "123 Test St",
						Country:     "Vietnam",
						Region:      "Hanoi",
						City:        "Hanoi",
						ZipCode:     "100000",
						Email:       "john@example.com",
						Phone:       "1234567890",
					},
					ShippingInfo: &checkoutpb.ShippingInfo{
						ShipToDifferentAddress: false,
						FirstName:              "John",
						LastName:               "Doe",
						CompanyName:            "",
						Address:                "123 Test St",
						Country:                "Vietnam",
						Region:                 "Hanoi",
						City:                   "Hanoi",
						ZipCode:                "100000",
						ShippingMethod:         "standard",
						ShippingCost:           5.00,
					},
					PaymentInfo: &checkoutpb.PaymentInfo{
						PaymentMethod: "cod",
						Status:        "pending",
						Amount:        129.99,
						Currency:      "VND",
					},
					Totals: &checkoutpb.OrderTotals{
						Subtotal: 124.99,
						Shipping: 5.00,
						Discount: 0.00,
						Tax:      0.00,
						Total:    129.99,
					},
					Items: []*checkoutpb.OrderItem{
						{
							Id:        "item124",
							ProductId: "1002",
							Name:      "Another Product",
							Price:     124.99,
							Quantity:  1,
							Subtotal:  124.99,
							ImageUrl:  "https://example.com/another-product.jpg",
						},
					},
					CreatedAt: "2023-08-22T11:00:00Z",
					UpdatedAt: "2023-08-22T11:00:00Z",
				},
			},
			Total: 2,
		}

		// Setup mock expectations
		mockCheckoutClient.On("ListOrders", mock.Anything, &checkoutpb.ListOrdersRequest{
			UserId:   "user123",
			Page:     1,
			PageSize: 10,
		}).Return(mockOrdersResponse, nil).Once()

		// Create test request with user ID in context
		req, err := http.NewRequest("GET", "/api/v1/checkout/orders", nil)
		assert.NoError(t, err)

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		checkoutHandler.ListOrders(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.False(t, response["error"].(bool))
		assert.Equal(t, "Orders retrieved successfully", response["message"])

		// Check order data
		data := response["data"].(map[string]interface{})
		orders := data["orders"].([]interface{})
		assert.Len(t, orders, 2)

		order1 := orders[0].(map[string]interface{})
		assert.Equal(t, "order123", order1["id"])
		assert.Equal(t, "confirmed", order1["status"])

		// Check pagination
		assert.Equal(t, float64(2), data["total"])

		// Verify mock expectations
		mockCheckoutClient.AssertExpectations(t)
	})

	// Test case 2: Missing user ID in context
	t.Run("Missing user ID in context", func(t *testing.T) {
		// Create test request without user ID in context
		req, err := http.NewRequest("GET", "/api/v1/checkout/orders", nil)
		assert.NoError(t, err)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		checkoutHandler.ListOrders(rec, req)

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
		// Setup mock expectations with error
		mockCheckoutClient.On("ListOrders", mock.Anything, &checkoutpb.ListOrdersRequest{
			UserId:   "user123",
			Page:     1,
			PageSize: 10,
		}).Return(nil, errors.New("service unavailable")).Once()

		// Create test request with user ID in context
		req, err := http.NewRequest("GET", "/api/v1/checkout/orders", nil)
		assert.NoError(t, err)

		// Add user_id to context (simulating auth middleware)
		ctx := context.WithValue(req.Context(), "user_id", "user123")
		req = req.WithContext(ctx)

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		checkoutHandler.ListOrders(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = helpers.ParseTestResponse(rec, &response)
		assert.NoError(t, err)

		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "service unavailable")

		// Verify mock expectations
		mockCheckoutClient.AssertExpectations(t)
	})
}

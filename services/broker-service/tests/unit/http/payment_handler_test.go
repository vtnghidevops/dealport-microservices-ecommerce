package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"broker-service/internal/handlers/http/payment"
	paymentpb "broker-service/proto/payment"
	"broker-service/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupPaymentHandlerTest sets up a new payment handler for testing with mocks
func setupPaymentHandlerTest() (*payment.Config, *mocks.MockPaymentServiceClient) {
	mockClient := new(mocks.MockPaymentServiceClient)

	paymentHandler := &payment.Config{
		CheckoutClient: mockClient,
	}

	return paymentHandler, mockClient
}

// TestHandleMomoCallback tests the HandleMomoCallback handler
func TestHandleMomoCallback(t *testing.T) {
	// Setup
	paymentHandler, mockClient := setupPaymentHandlerTest()

	// Test case 1: Successful callback processing
	t.Run("Successful callback processing", func(t *testing.T) {
		// Setup form data
		form := url.Values{}
		form.Add("orderId", "order123")
		form.Add("transId", "trans456")
		form.Add("resultCode", "0")
		form.Add("message", "Success")

		// Create test request with form data
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/callback", bytes.NewBufferString(form.Encode()))
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// Setup mock response
		mockResponse := &paymentpb.PaymentVerifyResponse{
			Success: true,
			Message: "Payment processed successfully",
			OrderId: "order123",
		}

		// Setup mock expectations with parameter matching
		mockClient.On("ProcessMomoCallback", mock.Anything, mock.MatchedBy(func(req *paymentpb.MomoCallbackRequest) bool {
			return req.Params["orderId"] == "order123" &&
				req.Params["transId"] == "trans456" &&
				req.Params["resultCode"] == "0" &&
				req.Params["message"] == "Success"
		})).Return(mockResponse, nil).Once()

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.HandleMomoCallback(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response payment.MomoCallbackResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, "Payment processed successfully", response.StatusMessage)
		assert.Equal(t, "order123", response.OrderID)

		// Verify mock expectations
		mockClient.AssertExpectations(t)
	})

	// Test case 2: Failed callback processing
	t.Run("Failed callback processing", func(t *testing.T) {
		// Setup form data
		form := url.Values{}
		form.Add("orderId", "order123")
		form.Add("transId", "trans456")
		form.Add("resultCode", "99")
		form.Add("message", "Failed")

		// Create test request with form data
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/callback", bytes.NewBufferString(form.Encode()))
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// Setup mock response
		mockResponse := &paymentpb.PaymentVerifyResponse{
			Success: false,
			Message: "Payment verification failed",
			OrderId: "order123",
		}

		// Setup mock expectations with parameter matching
		mockClient.On("ProcessMomoCallback", mock.Anything, mock.MatchedBy(func(req *paymentpb.MomoCallbackRequest) bool {
			return req.Params["orderId"] == "order123" &&
				req.Params["transId"] == "trans456" &&
				req.Params["resultCode"] == "99" &&
				req.Params["message"] == "Failed"
		})).Return(mockResponse, nil).Once()

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.HandleMomoCallback(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response payment.MomoCallbackResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "Payment verification failed", response.StatusMessage)
		assert.Equal(t, "order123", response.OrderID)

		// Verify mock expectations
		mockClient.AssertExpectations(t)
	})

	// Test case 3: Error calling checkout service
	t.Run("Error calling checkout service", func(t *testing.T) {
		// Setup form data
		form := url.Values{}
		form.Add("orderId", "order123")
		form.Add("transId", "trans456")
		form.Add("resultCode", "0")
		form.Add("message", "Success")

		// Create test request with form data
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/callback", bytes.NewBufferString(form.Encode()))
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// Setup mock expectations with error
		mockClient.On("ProcessMomoCallback", mock.Anything, mock.MatchedBy(func(req *paymentpb.MomoCallbackRequest) bool {
			return req.Params["orderId"] == "order123"
		})).Return(nil, errors.New("service error")).Once()

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.HandleMomoCallback(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response payment.MomoCallbackResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
		assert.Equal(t, "Internal server error processing callback", response.StatusMessage)

		// Verify mock expectations
		mockClient.AssertExpectations(t)
	})

	// Test case 4: No callback parameters
	t.Run("No callback parameters", func(t *testing.T) {
		// Create test request with empty form
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/callback", nil)
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.HandleMomoCallback(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response payment.MomoCallbackResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Contains(t, response.StatusMessage, "Failed to parse parameters")

		// No mock expectations for this case
	})
}

// TestCreateMomoPayment tests the CreateMomoPayment handler
func TestCreateMomoPayment(t *testing.T) {
	// Setup
	paymentHandler, mockClient := setupPaymentHandlerTest()

	// Test case 1: Successful payment creation
	t.Run("Successful payment creation", func(t *testing.T) {
		// Create request payload
		requestPayload := payment.CreateMomoPaymentRequest{
			OrderID:   "order123",
			Amount:    100000,
			OrderInfo: "Test payment",
			ReturnURL: "https://example.com/return",
		}

		jsonPayload, err := json.Marshal(requestPayload)
		assert.NoError(t, err)

		// Create test request
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/create", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Setup mock response
		mockResponse := &paymentpb.MomoPaymentResponse{
			Success:       true,
			Message:       "Payment URL created successfully",
			PaymentUrl:    "https://test-payment.momo.vn/pay/order123",
			OrderId:       "order123",
			RequestId:     "req123",
			TransactionId: "trans123",
		}

		// Setup mock expectations
		mockClient.On("CreateMomoPayment", mock.Anything, mock.MatchedBy(func(req *paymentpb.MomoPaymentRequest) bool {
			return req.PaymentInfo.OrderId == "order123" &&
				req.PaymentInfo.Amount == 100000 &&
				req.OrderInfo == "Test payment" &&
				req.PaymentInfo.ReturnUrl == "https://example.com/return"
		})).Return(mockResponse, nil).Once()

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.CreateMomoPayment(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "https://test-payment.momo.vn/pay/order123", data["paymentUrl"])
		assert.Equal(t, "order123", data["orderId"])
		assert.Equal(t, "req123", data["requestId"])
		assert.Equal(t, "trans123", data["transactionId"])
		assert.Equal(t, float64(100000), data["amount"])

		// Verify mock expectations
		mockClient.AssertExpectations(t)
	})

	// Test case 2: Invalid request format
	t.Run("Invalid request format", func(t *testing.T) {
		// Create invalid JSON
		invalidPayload := []byte(`{invalid json}`)

		// Create test request
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/create", bytes.NewBuffer(invalidPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.CreateMomoPayment(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check for the message directly since error field might not exist
		message, exists := response["message"]
		assert.True(t, exists)
		assert.Equal(t, "Invalid request format", message)

		// No mock expectations for this case
	})

	// Test case 3: Missing required fields
	t.Run("Missing required fields", func(t *testing.T) {
		// Create incomplete request payload
		requestPayload := map[string]interface{}{
			"orderId":   "",
			"amount":    0,
			"orderInfo": "Test payment",
			"returnUrl": "",
		}

		jsonPayload, err := json.Marshal(requestPayload)
		assert.NoError(t, err)

		// Create test request
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/create", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.CreateMomoPayment(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check for the message directly since error field might not exist
		message, exists := response["message"]
		assert.True(t, exists)
		assert.Equal(t, "Missing required fields", message)

		// No mock expectations for this case
	})

	// Test case 4: Service error
	t.Run("Service error", func(t *testing.T) {
		// Create request payload
		requestPayload := payment.CreateMomoPaymentRequest{
			OrderID:   "order123",
			Amount:    100000,
			OrderInfo: "Test payment",
			ReturnURL: "https://example.com/return",
		}

		jsonPayload, err := json.Marshal(requestPayload)
		assert.NoError(t, err)

		// Create test request
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/create", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Setup mock expectations with error
		mockClient.On("CreateMomoPayment", mock.Anything, mock.MatchedBy(func(req *paymentpb.MomoPaymentRequest) bool {
			return req.PaymentInfo.OrderId == "order123"
		})).Return(nil, errors.New("service unavailable")).Once()

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.CreateMomoPayment(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check for the message directly since error field might not exist
		message, exists := response["message"]
		assert.True(t, exists)
		assert.Contains(t, message, "Failed to process payment: service unavailable")

		// Verify mock expectations
		mockClient.AssertExpectations(t)
	})

	// Test case 5: Failed payment creation
	t.Run("Failed payment creation", func(t *testing.T) {
		// Create request payload
		requestPayload := payment.CreateMomoPaymentRequest{
			OrderID:   "order123",
			Amount:    100000,
			OrderInfo: "Test payment",
			ReturnURL: "https://example.com/return",
		}

		jsonPayload, err := json.Marshal(requestPayload)
		assert.NoError(t, err)

		// Create test request
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/create", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Setup mock response for failure
		mockResponse := &paymentpb.MomoPaymentResponse{
			Success: false,
			Message: "Invalid payment request",
		}

		// Setup mock expectations
		mockClient.On("CreateMomoPayment", mock.Anything, mock.MatchedBy(func(req *paymentpb.MomoPaymentRequest) bool {
			return req.PaymentInfo.OrderId == "order123"
		})).Return(mockResponse, nil).Once()

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.CreateMomoPayment(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check for the message directly since error field might not exist
		message, exists := response["message"]
		assert.True(t, exists)
		assert.Equal(t, "Invalid payment request", message)

		// Verify mock expectations
		mockClient.AssertExpectations(t)
	})
}

// TestVerifyMomoPayment tests the VerifyMomoPayment handler
func TestVerifyMomoPayment(t *testing.T) {
	// Setup
	paymentHandler, mockClient := setupPaymentHandlerTest()

	// Test case 1: Successful payment verification
	t.Run("Successful payment verification", func(t *testing.T) {
		// Create request payload
		requestPayload := payment.VerifyMomoPaymentRequest{
			Params: map[string]string{
				"orderId":    "order123",
				"transId":    "trans456",
				"resultCode": "0",
				"message":    "Success",
			},
		}

		jsonPayload, err := json.Marshal(requestPayload)
		assert.NoError(t, err)

		// Create test request
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/verify", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Setup mock response
		mockResponse := &paymentpb.PaymentVerifyResponse{
			Success:       true,
			Message:       "Payment verified successfully",
			OrderId:       "order123",
			TransactionId: "trans456",
			Amount:        100000,
		}

		// Setup mock expectations
		mockClient.On("VerifyMomoPayment", mock.Anything, mock.MatchedBy(func(req *paymentpb.MomoVerifyRequest) bool {
			return req.Params["orderId"] == "order123" &&
				req.Params["transId"] == "trans456" &&
				req.Params["resultCode"] == "0"
		})).Return(mockResponse, nil).Once()

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.VerifyMomoPayment(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check data field exists and contains expected values
		dataVal, exists := response["data"]
		assert.True(t, exists, "Data field should exist in the response")
		if exists {
			data, ok := dataVal.(map[string]interface{})
			assert.True(t, ok, "Data field should be a map")
			if ok {
				assert.Equal(t, "order123", data["orderId"])
				assert.Equal(t, "trans456", data["transactionId"])
				assert.Equal(t, float64(100000), data["amount"])
			}
		}

		// Verify mock expectations
		mockClient.AssertExpectations(t)
	})

	// Test case 2: Invalid request format
	t.Run("Invalid request format", func(t *testing.T) {
		// Create invalid JSON
		invalidPayload := []byte(`{invalid json}`)

		// Create test request
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/verify", bytes.NewBuffer(invalidPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.VerifyMomoPayment(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check for the message directly
		message, exists := response["message"]
		assert.True(t, exists, "Message field should exist in the response")
		assert.Equal(t, "Invalid request format", message)

		// No mock expectations for this case
	})

	// Test case 3: Missing callback parameters
	t.Run("Missing callback parameters", func(t *testing.T) {
		// Create request payload with empty params
		requestPayload := payment.VerifyMomoPaymentRequest{
			Params: map[string]string{},
		}

		jsonPayload, err := json.Marshal(requestPayload)
		assert.NoError(t, err)

		// Create test request
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/verify", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.VerifyMomoPayment(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check for the message directly
		message, exists := response["message"]
		assert.True(t, exists, "Message field should exist in the response")
		assert.Equal(t, "Missing callback parameters", message)

		// No mock expectations for this case
	})

	// Test case 4: Service error
	t.Run("Service error", func(t *testing.T) {
		// Create request payload
		requestPayload := payment.VerifyMomoPaymentRequest{
			Params: map[string]string{
				"orderId":    "order123",
				"transId":    "trans456",
				"resultCode": "0",
				"message":    "Success",
			},
		}

		jsonPayload, err := json.Marshal(requestPayload)
		assert.NoError(t, err)

		// Create test request
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/verify", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Setup mock expectations with error
		mockClient.On("VerifyMomoPayment", mock.Anything, mock.MatchedBy(func(req *paymentpb.MomoVerifyRequest) bool {
			return req.Params["orderId"] == "order123"
		})).Return(nil, errors.New("service unavailable")).Once()

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.VerifyMomoPayment(rec, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check for the message directly
		message, exists := response["message"]
		assert.True(t, exists, "Message field should exist in the response")
		assert.Contains(t, message, "Failed to verify payment: service unavailable")

		// Verify mock expectations
		mockClient.AssertExpectations(t)
	})

	// Test case 5: Failed payment verification
	t.Run("Failed payment verification", func(t *testing.T) {
		// Create request payload
		requestPayload := payment.VerifyMomoPaymentRequest{
			Params: map[string]string{
				"orderId":    "order123",
				"transId":    "trans456",
				"resultCode": "99",
				"message":    "Failed",
			},
		}

		jsonPayload, err := json.Marshal(requestPayload)
		assert.NoError(t, err)

		// Create test request
		req, err := http.NewRequest("POST", "/api/v1/payment/momo/verify", bytes.NewBuffer(jsonPayload))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Setup mock response for failure
		mockResponse := &paymentpb.PaymentVerifyResponse{
			Success: false,
			Message: "Payment verification failed",
			OrderId: "order123",
		}

		// Setup mock expectations
		mockClient.On("VerifyMomoPayment", mock.Anything, mock.MatchedBy(func(req *paymentpb.MomoVerifyRequest) bool {
			return req.Params["orderId"] == "order123" &&
				req.Params["resultCode"] == "99"
		})).Return(mockResponse, nil).Once()

		// Create response recorder
		rec := httptest.NewRecorder()

		// Call the handler
		paymentHandler.VerifyMomoPayment(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		dataVal, ok := response["data"]
		assert.True(t, ok, "Data field should exist in the response")
		if ok {
			data, ok := dataVal.(map[string]interface{})
			assert.True(t, ok, "Data field should be a map")
			if ok {
				assert.Equal(t, false, data["success"])
				assert.Equal(t, "Payment verification failed", data["message"])
				assert.Equal(t, "order123", data["orderId"])
			}
		}

		// Verify mock expectations
		mockClient.AssertExpectations(t)
	})
}

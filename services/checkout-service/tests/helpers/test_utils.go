package helpers

import (
	"checkout-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

// CreateTestOrder creates a test order with default values
func CreateTestOrder(userID string) *domain.Order {
	return &domain.Order{
		ID:          uuid.New().String(),
		UserID:      userID,
		OrderNumber: "TEST-" + uuid.New().String()[0:8],
		Status:      "pending",
		Items: []domain.OrderItem{
			{
				ID:        uuid.New().String(),
				ProductID: "prod-1",
				Name:      "Test Product",
				Price:     19.99,
				Quantity:  2,
				ImageURL:  "https://example.com/image.jpg",
				Subtotal:  39.98,
			},
		},
		BillingInfo: domain.BillingInfo{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Phone:     "1234567890",
			Address:   "123 Test St",
			City:      "Test City",
			Region:    "Test Region",
			Country:   "Test Country",
			ZipCode:   "12345",
		},
		ShippingInfo: domain.ShippingInfo{
			ShipToDifferentAddress: false,
			ShippingMethod:         "standard",
			ShippingCost:           5.99,
		},
		PaymentInfo: domain.PaymentInfo{
			PaymentMethod: "card",
			Status:        "pending",
			Currency:      "USD",
		},
		Totals: domain.OrderTotals{
			Subtotal: 39.98,
			Shipping: 5.99,
			Tax:      4.00,
			Total:    49.97,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// SetupOrderService sets up a mock order service with its dependencies for testing
// Currently not used directly in tests - commented out but kept for reference
/*
func SetupOrderService() (domain.OrderService, *mocks.MockOrderRepository, *mocks.SimpleEventEmitter) {
	// Create the mocks
	mockRepo := mocks.NewMockOrderRepository()
	mockEmitter := mocks.NewSimpleEventEmitter()

	// Use the adapter to create the service with mock emitter
	orderService := NewOrderServiceWithAdapter(mockRepo, mockEmitter)

	return orderService, mockRepo, mockEmitter
}
*/

// CreateMockPaymentService creates a mock payment service for testing
// Currently not used directly in tests - commented out but kept for reference
/*
func CreateMockPaymentService() *mocks.MockPaymentService {
	return &mocks.MockPaymentService{}
}
*/

// CreateSuccessfulPaymentResult creates a successful payment result for testing
// Currently not used directly in tests - commented out but kept for reference
/*
func CreateSuccessfulPaymentResult(orderID string) domain.PaymentResult {
	return domain.PaymentResult{
		Success:       true,
		TransactionID: "txn_" + orderID,
		Status:        "completed",
		Message:       "Payment processed successfully",
	}
}
*/

// CreateFailedPaymentResult creates a failed payment result for testing
// Currently not used directly in tests - commented out but kept for reference
/*
func CreateFailedPaymentResult(errorMessage string) domain.PaymentResult {
	return domain.PaymentResult{
		Success:       false,
		TransactionID: "",
		Status:        "failed",
		Message:       errorMessage,
	}
}
*/

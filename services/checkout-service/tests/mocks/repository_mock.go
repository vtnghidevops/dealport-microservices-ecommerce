package mocks

import (
	"checkout-service/internal/domain"
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockOrderRepository is a mock implementation of the order repository interface
type MockOrderRepository struct {
	mock.Mock
}

// CreateOrder mocks the method for creating a new order
func (m *MockOrderRepository) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	args := m.Called(ctx, order)
	var result *domain.Order
	if args.Get(0) != nil {
		result = args.Get(0).(*domain.Order)
	}
	return result, args.Error(1)
}

// GetOrderByID mocks the method for retrieving an order by its ID
func (m *MockOrderRepository) GetOrderByID(ctx context.Context, orderID string) (*domain.Order, error) {
	args := m.Called(ctx, orderID)
	var order *domain.Order
	if args.Get(0) != nil {
		order = args.Get(0).(*domain.Order)
	}
	return order, args.Error(1)
}

// ListOrdersByUserID mocks the method for retrieving orders by user ID
func (m *MockOrderRepository) ListOrdersByUserID(ctx context.Context, userID string, skip, limit int) ([]*domain.Order, int, error) {
	args := m.Called(ctx, userID, skip, limit)
	return args.Get(0).([]*domain.Order), args.Get(1).(int), args.Error(2)
}

// ListAllOrders mocks the method for retrieving all orders
func (m *MockOrderRepository) ListAllOrders(ctx context.Context, skip, limit int) ([]*domain.Order, int, error) {
	args := m.Called(ctx, skip, limit)
	return args.Get(0).([]*domain.Order), args.Get(1).(int), args.Error(2)
}

// UpdateOrderStatus mocks the method for updating an order's status
func (m *MockOrderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}

// UpdatePaymentInfo mocks the method for updating payment information
func (m *MockOrderRepository) UpdatePaymentInfo(ctx context.Context, orderID string, paymentInfo domain.PaymentInfo) error {
	args := m.Called(ctx, orderID, paymentInfo)
	return args.Error(0)
}

// GetUserTotalSpend mocks the method for getting user's total spend
func (m *MockOrderRepository) GetUserTotalSpend(ctx context.Context, userID string) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

// GetUserOrderCount mocks the method for getting user's order count
func (m *MockOrderRepository) GetUserOrderCount(ctx context.Context, userID string) (int, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int), args.Error(1)
}

// CreateMockOrder is a helper method for creating orders with items for testing
func (m *MockOrderRepository) CreateMockOrder(userID string, status string) *domain.Order {
	now := time.Now()

	return &domain.Order{
		ID:          userID + "-order-id",
		UserID:      userID,
		Status:      status,
		OrderNumber: "ORD-" + userID[0:8],
		Items: []domain.OrderItem{
			{
				ID:        "item1",
				ProductID: "product1",
				Name:      "Test Product 1",
				Price:     29.99,
				Quantity:  2,
				Subtotal:  59.98,
			},
			{
				ID:        "item2",
				ProductID: "product2",
				Name:      "Test Product 2",
				Price:     49.99,
				Quantity:  1,
				Subtotal:  49.99,
			},
		},
		BillingInfo: domain.BillingInfo{
			FirstName: "John",
			LastName:  "Doe",
			Address:   "123 Test St",
			City:      "Test City",
			ZipCode:   "12345",
			Country:   "Test Country",
			Phone:     "1234567890",
			Email:     "test@example.com",
		},
		ShippingInfo: domain.ShippingInfo{
			FirstName:      "John",
			LastName:       "Doe",
			Address:        "123 Test St",
			City:           "Test City",
			ZipCode:        "12345",
			Country:        "Test Country",
			ShippingMethod: "standard",
			ShippingCost:   5.99,
		},
		PaymentInfo: domain.PaymentInfo{
			PaymentMethod: "credit_card",
			Status:        "pending",
			Amount:        126.95,
			Currency:      "USD",
		},
		Totals: domain.OrderTotals{
			Subtotal: 109.97,
			Tax:      10.99,
			Shipping: 5.99,
			Total:    126.95,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewMockOrderRepository creates a new mock order repository instance
func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{}
}

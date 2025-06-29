package mocks

import (
	"checkout-service/internal/domain"
	"context"
	"sync"
)

// SimpleOrderRepository mocks the OrderRepository interface for testing with simpler functionality
type SimpleOrderRepository struct {
	mutex           sync.RWMutex
	orders          map[string]*domain.Order
	userOrders      map[string][]*domain.Order
	userTotalSpends map[string]float64
	userOrderCounts map[string]int

	// For tracking method calls
	CreateOrderCalled        bool
	GetOrderByIDCalled       bool
	ListOrdersByUserIDCalled bool
	UpdateOrderStatusCalled  bool
	UpdatePaymentInfoCalled  bool
	GetUserTotalSpendCalled  bool
	GetUserOrderCountCalled  bool

	// For injecting errors
	CreateOrderError        error
	GetOrderByIDError       error
	ListOrdersByUserIDError error
	UpdateOrderStatusError  error
	UpdatePaymentInfoError  error
	GetUserTotalSpendError  error
	GetUserOrderCountError  error
}

// NewSimpleOrderRepository creates a new mock repository
func NewSimpleOrderRepository() *SimpleOrderRepository {
	return &SimpleOrderRepository{
		orders:          make(map[string]*domain.Order),
		userOrders:      make(map[string][]*domain.Order),
		userTotalSpends: make(map[string]float64),
		userOrderCounts: make(map[string]int),
	}
}

// CreateOrder mocks creating an order
func (m *SimpleOrderRepository) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	m.CreateOrderCalled = true

	if m.CreateOrderError != nil {
		return nil, m.CreateOrderError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Store the order
	m.orders[order.ID] = order

	// Add to user orders
	m.userOrders[order.UserID] = append(m.userOrders[order.UserID], order)

	// Update user spending
	m.userTotalSpends[order.UserID] += order.Totals.Total

	// Update user order count
	m.userOrderCounts[order.UserID]++

	return order, nil
}

// GetOrderByID mocks retrieving an order by ID
func (m *SimpleOrderRepository) GetOrderByID(ctx context.Context, orderID string) (*domain.Order, error) {
	m.GetOrderByIDCalled = true

	if m.GetOrderByIDError != nil {
		return nil, m.GetOrderByIDError
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	order, exists := m.orders[orderID]
	if !exists {
		return nil, domain.ErrOrderNotFound
	}

	return order, nil
}

// ListOrdersByUserID mocks listing orders for a user
func (m *SimpleOrderRepository) ListOrdersByUserID(ctx context.Context, userID string, skip, limit int) ([]*domain.Order, int, error) {
	m.ListOrdersByUserIDCalled = true

	if m.ListOrdersByUserIDError != nil {
		return nil, 0, m.ListOrdersByUserIDError
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	orders, exists := m.userOrders[userID]
	if !exists {
		return []*domain.Order{}, 0, nil
	}

	total := len(orders)

	// Apply pagination
	start := skip
	if start >= total {
		return []*domain.Order{}, total, nil
	}

	end := start + limit
	if end > total {
		end = total
	}

	return orders[start:end], total, nil
}

// ListAllOrders mocks listing all orders with pagination
func (m *SimpleOrderRepository) ListAllOrders(ctx context.Context, skip, limit int) ([]*domain.Order, int, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var allOrders []*domain.Order
	for _, order := range m.orders {
		allOrders = append(allOrders, order)
	}

	total := len(allOrders)

	// Apply pagination
	start := skip
	if start >= total {
		return []*domain.Order{}, total, nil
	}

	end := start + limit
	if end > total {
		end = total
	}

	return allOrders[start:end], total, nil
}

// UpdateOrderStatus mocks updating an order's status
func (m *SimpleOrderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	m.UpdateOrderStatusCalled = true

	if m.UpdateOrderStatusError != nil {
		return m.UpdateOrderStatusError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	order, exists := m.orders[orderID]
	if !exists {
		return domain.ErrOrderNotFound
	}

	order.Status = status
	return nil
}

// UpdatePaymentInfo mocks updating payment information
func (m *SimpleOrderRepository) UpdatePaymentInfo(ctx context.Context, orderID string, paymentInfo domain.PaymentInfo) error {
	m.UpdatePaymentInfoCalled = true

	if m.UpdatePaymentInfoError != nil {
		return m.UpdatePaymentInfoError
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	order, exists := m.orders[orderID]
	if !exists {
		return domain.ErrOrderNotFound
	}

	order.PaymentInfo = paymentInfo
	return nil
}

// GetUserTotalSpend mocks retrieving a user's total spend
func (m *SimpleOrderRepository) GetUserTotalSpend(ctx context.Context, userID string) (float64, error) {
	m.GetUserTotalSpendCalled = true

	if m.GetUserTotalSpendError != nil {
		return 0, m.GetUserTotalSpendError
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.userTotalSpends[userID], nil
}

// GetUserOrderCount mocks retrieving a user's order count
func (m *SimpleOrderRepository) GetUserOrderCount(ctx context.Context, userID string) (int, error) {
	m.GetUserOrderCountCalled = true

	if m.GetUserOrderCountError != nil {
		return 0, m.GetUserOrderCountError
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.userOrderCounts[userID], nil
}

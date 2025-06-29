package helpers

import (
	"checkout-service/internal/domain"
	"checkout-service/tests/mocks"
)

// NewOrderServiceWithSimpleMocks creates a new order service instance with simple mock dependencies
func NewOrderServiceWithSimpleMocks(mockRepo *mocks.SimpleOrderRepository, mockEmitter *mocks.SimpleEventEmitter) domain.OrderService {
	return NewOrderServiceWithAdapter(mockRepo, mockEmitter)
}

// NewTestOrderService creates a new order service with default simple mocks
func NewTestOrderService() (domain.OrderService, *mocks.SimpleOrderRepository, *mocks.SimpleEventEmitter) {
	mockRepo := mocks.NewSimpleOrderRepository()
	mockEmitter := mocks.NewSimpleEventEmitter()
	orderService := NewOrderServiceWithSimpleMocks(mockRepo, mockEmitter)
	return orderService, mockRepo, mockEmitter
}

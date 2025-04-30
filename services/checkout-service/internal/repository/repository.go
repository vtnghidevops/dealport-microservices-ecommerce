package repository

import (
	"checkout-service/internal/domain"
	"checkout-service/internal/repository/mongodb"
	"checkout-service/internal/service"
	"checkout-service/internal/transport/grpc"

	"go.mongodb.org/mongo-driver/mongo"
)

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *mongo.Database) domain.OrderRepository {
	return mongodb.NewOrderRepository(db)
}

// NewOrderService creates a new order service
func NewOrderService(repo domain.OrderRepository) domain.OrderService {
	return service.NewOrderService(repo)
}

// NewCheckoutServiceHandler creates a new checkout service handler
func NewCheckoutServiceHandler(orderService domain.OrderService) *grpc.CheckoutServiceHandler {
	return grpc.NewCheckoutServiceHandler(orderService)
}

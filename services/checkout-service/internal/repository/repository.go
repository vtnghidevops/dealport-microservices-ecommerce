package repository

import (
	"checkout-service/internal/domain"
	"checkout-service/internal/event"
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
func NewOrderService(repo domain.OrderRepository, emitter *event.Emitter) domain.OrderService {
	return service.NewOrderService(repo, emitter)
}

// NewCheckoutHandler creates a new checkout service handler
func NewCheckoutHandler(orderService domain.OrderService) *grpc.CheckoutHandler {
	return grpc.NewCheckoutHandler(orderService)
}

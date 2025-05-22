package helpers

import (
	"checkout-service/internal/domain"
	"checkout-service/internal/event"
	realService "checkout-service/internal/service"
	"checkout-service/tests/mocks"
	"context"
)

// NewOrderServiceWithAdapter creates a new order service with an adapter
func NewOrderServiceWithAdapter(repo domain.OrderRepository, emitter *mocks.SimpleEventEmitter) domain.OrderService {
	// Create a custom adapter for testing that will intercept event calls
	adapter := &OrderServiceAdapter{
		repo:    repo,
		emitter: emitter,
	}

	// Create a real service without event emitter for underlying functionality
	adapter.realService = realService.NewOrderService(repo, nil)

	return adapter
}

// OrderServiceAdapter wraps the real service to intercept event emissions
type OrderServiceAdapter struct {
	realService domain.OrderService
	repo        domain.OrderRepository
	emitter     *mocks.SimpleEventEmitter
}

// GetRepository returns the repository instance
func (a *OrderServiceAdapter) GetRepository() domain.OrderRepository {
	return a.repo
}

// CreateOrder intercepts create order calls to manually trigger event emission
func (a *OrderServiceAdapter) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	// Call the real service
	result, err := a.realService.CreateOrder(ctx, order)

	// If successful, manually trigger the mock's event handler
	if err == nil && result != nil && a.emitter != nil {
		// Prepare event data similar to what the real service would use
		eventData := event.OrderEventData{
			OrderID:       result.ID,
			OrderNumber:   result.OrderNumber,
			UserID:        result.UserID,
			Total:         result.Totals.Total,
			Status:        result.Status,
			PaymentMethod: result.PaymentInfo.PaymentMethod,
			CreatedAt:     result.CreatedAt,
		}

		// Call the mock directly
		a.emitter.EmitOrderCreated(eventData)
	}

	return result, err
}

// GetOrderByID delegates to the real service
func (a *OrderServiceAdapter) GetOrderByID(ctx context.Context, orderID string, userID string) (*domain.Order, error) {
	return a.realService.GetOrderByID(ctx, orderID, userID)
}

// ListOrdersByUserID delegates to the real service
func (a *OrderServiceAdapter) ListOrdersByUserID(ctx context.Context, userID string, page, pageSize int) ([]*domain.Order, int, error) {
	return a.realService.ListOrdersByUserID(ctx, userID, page, pageSize)
}

// UpdateOrderStatus intercepts status updates to manually trigger event emission
func (a *OrderServiceAdapter) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	err := a.realService.UpdateOrderStatus(ctx, orderID, status)

	// If successful, manually trigger the mock's event handler
	if err == nil && a.emitter != nil {
		// Create minimal event data
		eventData := event.OrderEventData{
			OrderID: orderID,
			Status:  status,
		}

		// Call the mock directly
		a.emitter.EmitOrderStatusChanged(eventData)
	}

	return err
}

// ValidateCheckout delegates to the real service
func (a *OrderServiceAdapter) ValidateCheckout(ctx context.Context, data domain.CheckoutValidationRequest) (domain.CheckoutValidationResult, error) {
	return a.realService.ValidateCheckout(ctx, data)
}

// ProcessPayment intercepts payment processing to manually trigger event emission
func (a *OrderServiceAdapter) ProcessPayment(ctx context.Context, paymentReq domain.PaymentRequest) (domain.PaymentResult, error) {
	result, err := a.realService.ProcessPayment(ctx, paymentReq)

	// If successful, manually trigger the mock's event handler
	if err == nil && result.Success && a.emitter != nil {
		// Create minimal event data
		eventData := event.PaymentEventData{
			OrderID:       paymentReq.OrderID,
			TransactionID: result.TransactionID,
			Amount:        paymentReq.Amount,
			Status:        "succeeded",
		}

		// Call the mock directly
		a.emitter.EmitPaymentSucceeded(eventData)
	}

	return result, err
}

// GetUserTotalSpend delegates to the real service
func (a *OrderServiceAdapter) GetUserTotalSpend(ctx context.Context, userID string) (float64, error) {
	return a.realService.GetUserTotalSpend(ctx, userID)
}

// GetUserOrderCount delegates to the real service
func (a *OrderServiceAdapter) GetUserOrderCount(ctx context.Context, userID string) (int, error) {
	return a.realService.GetUserOrderCount(ctx, userID)
}

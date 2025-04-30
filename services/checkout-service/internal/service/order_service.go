package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"checkout-service/internal/domain"

	"github.com/google/uuid"
)

// OrderService implements domain.OrderService
type OrderService struct {
	repo domain.OrderRepository
}

// NewOrderService creates a new order service
func NewOrderService(repo domain.OrderRepository) domain.OrderService {
	return &OrderService{
		repo: repo,
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	// Debug: In ra thông tin đơn hàng được gửi tới service
	log.Printf("OrderService.CreateOrder received: %+v", order)
	log.Printf("BillingInfo: %+v", order.BillingInfo)
	log.Printf("ShippingInfo: %+v", order.ShippingInfo)
	log.Printf("PaymentInfo: %+v", order.PaymentInfo)

	// Validate order data
	if order.UserID == "" {
		return nil, errors.Join(domain.ErrInvalidOrderData, fmt.Errorf("user ID is required"))
	}
	if len(order.Items) == 0 {
		return nil, errors.Join(domain.ErrInvalidOrderData, fmt.Errorf("order must have at least one item"))
	}

	// Calculate order totals
	subtotal := 0.0
	for i := range order.Items {
		// Generate ID for item if not present
		if order.Items[i].ID == "" {
			order.Items[i].ID = uuid.New().String()
		}

		// Calculate subtotal for item
		order.Items[i].Subtotal = order.Items[i].Price * float64(order.Items[i].Quantity)
		subtotal += order.Items[i].Subtotal
	}

	// Set default shipping cost if not provided
	shippingCost := order.ShippingInfo.ShippingCost
	if shippingCost <= 0 {
		shippingCost = 5.0 // Default shipping cost
	}

	// Calculate tax (e.g., 10% of subtotal)
	tax := subtotal * 0.1

	// Calculate discount based on coupon code
	discount := 0.0
	if order.CouponCode != "" {
		// In a real implementation, would validate the coupon code
		// For now, just apply a placeholder discount
		discount = subtotal * 0.05 // 5% discount
	}

	// Calculate total
	total := subtotal + shippingCost + tax - discount

	// Set the order totals
	order.Totals = domain.OrderTotals{
		Subtotal: subtotal,
		Shipping: shippingCost,
		Discount: discount,
		Tax:      tax,
		Total:    total,
	}

	// Debug: Trước khi gửi đến repo
	log.Printf("Order before saving to repo: %+v", order)
	log.Printf("BillingInfo before saving: %+v", order.BillingInfo)
	log.Printf("ShippingInfo before saving: %+v", order.ShippingInfo)
	log.Printf("PaymentInfo before saving: %+v", order.PaymentInfo)

	// Save the order
	savedOrder, err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		log.Printf("Error from repo.CreateOrder: %v", err)
		return nil, err
	}

	// Debug: Sau khi nhận kết quả từ repo
	log.Printf("Order after repo.CreateOrder: %+v", savedOrder)
	log.Printf("BillingInfo after repo: %+v", savedOrder.BillingInfo)
	log.Printf("ShippingInfo after repo: %+v", savedOrder.ShippingInfo)
	log.Printf("PaymentInfo after repo: %+v", savedOrder.PaymentInfo)

	return savedOrder, nil
}

// GetOrderByID retrieves an order by ID
func (s *OrderService) GetOrderByID(ctx context.Context, orderID string, userID string) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Verify that the order belongs to the user (security check)
	if order.UserID != userID {
		return nil, domain.ErrOrderNotFound
	}

	return order, nil
}

// ListOrdersByUserID retrieves orders for a user with pagination
func (s *OrderService) ListOrdersByUserID(ctx context.Context, userID string, page, pageSize int) ([]*domain.Order, int, error) {
	if userID == "" {
		log.Printf("ERROR: ListOrdersByUserID called with empty userID")
		return nil, 0, errors.Join(domain.ErrInvalidOrderData, fmt.Errorf("user ID is required"))
	}

	log.Printf("ListOrdersByUserID: Fetching orders for user %s (page=%d, pageSize=%d)", userID, page, pageSize)

	// Calculate skip for pagination
	skip := (page - 1) * pageSize
	if skip < 0 {
		skip = 0
	}

	orders, total, err := s.repo.ListOrdersByUserID(ctx, userID, skip, pageSize)
	if err != nil {
		log.Printf("ERROR: Failed to fetch orders for user %s: %v", userID, err)
		return nil, 0, err
	}

	log.Printf("SUCCESS: Found %d orders for user %s (total: %d)", len(orders), userID, total)
	for i, order := range orders {
		log.Printf("Order %d: ID=%s, OrderNumber=%s, Status=%s, CreatedAt=%v",
			i+1, order.ID, order.OrderNumber, order.Status, order.CreatedAt)
	}

	return orders, total, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"shipped":    true,
		"delivered":  true,
		"cancelled":  true,
	}

	if !validStatuses[status] {
		return errors.Join(domain.ErrInvalidOrderData, fmt.Errorf("invalid status: %s", status))
	}

	return s.repo.UpdateOrderStatus(ctx, orderID, status)
}

// ValidateCheckout validates checkout data
func (s *OrderService) ValidateCheckout(ctx context.Context, data domain.CheckoutValidationRequest) (domain.CheckoutValidationResult, error) {
	var result domain.CheckoutValidationResult
	var validationErrors []domain.ValidationError

	// Validate billing info
	if data.BillingInfo.FirstName == "" {
		validationErrors = append(validationErrors, domain.ValidationError{
			Field:   "billingInfo.firstName",
			Message: "First name is required",
		})
	}
	if data.BillingInfo.LastName == "" {
		validationErrors = append(validationErrors, domain.ValidationError{
			Field:   "billingInfo.lastName",
			Message: "Last name is required",
		})
	}
	if data.BillingInfo.Email == "" {
		validationErrors = append(validationErrors, domain.ValidationError{
			Field:   "billingInfo.email",
			Message: "Email is required",
		})
	}
	if data.BillingInfo.Phone == "" {
		validationErrors = append(validationErrors, domain.ValidationError{
			Field:   "billingInfo.phone",
			Message: "Phone is required",
		})
	}
	if data.BillingInfo.Address == "" {
		validationErrors = append(validationErrors, domain.ValidationError{
			Field:   "billingInfo.address",
			Message: "Address is required",
		})
	}
	if data.BillingInfo.City == "" {
		validationErrors = append(validationErrors, domain.ValidationError{
			Field:   "billingInfo.city",
			Message: "City is required",
		})
	}
	if data.BillingInfo.ZipCode == "" {
		validationErrors = append(validationErrors, domain.ValidationError{
			Field:   "billingInfo.zipCode",
			Message: "ZIP code is required",
		})
	}

	// If shipping to a different address, validate shipping info
	if data.ShippingInfo.ShipToDifferentAddress {
		if data.ShippingInfo.FirstName == "" {
			validationErrors = append(validationErrors, domain.ValidationError{
				Field:   "shippingInfo.firstName",
				Message: "First name is required for shipping",
			})
		}
		if data.ShippingInfo.LastName == "" {
			validationErrors = append(validationErrors, domain.ValidationError{
				Field:   "shippingInfo.lastName",
				Message: "Last name is required for shipping",
			})
		}
		if data.ShippingInfo.Address == "" {
			validationErrors = append(validationErrors, domain.ValidationError{
				Field:   "shippingInfo.address",
				Message: "Address is required for shipping",
			})
		}
		if data.ShippingInfo.City == "" {
			validationErrors = append(validationErrors, domain.ValidationError{
				Field:   "shippingInfo.city",
				Message: "City is required for shipping",
			})
		}
		if data.ShippingInfo.ZipCode == "" {
			validationErrors = append(validationErrors, domain.ValidationError{
				Field:   "shippingInfo.zipCode",
				Message: "ZIP code is required for shipping",
			})
		}
	}

	// Validate that there are items in the order
	if len(data.Items) == 0 {
		validationErrors = append(validationErrors, domain.ValidationError{
			Field:   "items",
			Message: "Order must have at least one item",
		})
	}

	// Validate payment method
	if data.PaymentMethod == "" {
		validationErrors = append(validationErrors, domain.ValidationError{
			Field:   "paymentMethod",
			Message: "Payment method is required",
		})
	} else {
		// Check if payment method is supported
		supportedPaymentMethods := map[string]bool{
			// "credit_card": true,
			"vnpay": true,
			// "stripe":      true,
			"cod":  true, // Cash on delivery
			"momo": true, // Cash on delivery
		}

		if !supportedPaymentMethods[data.PaymentMethod] {
			validationErrors = append(validationErrors, domain.ValidationError{
				Field:   "paymentMethod",
				Message: "Unsupported payment method",
			})
		}
	}

	// Set the result
	result.Valid = len(validationErrors) == 0
	result.Errors = validationErrors

	return result, nil
}

// ProcessPayment processes a payment for an order
func (s *OrderService) ProcessPayment(ctx context.Context, paymentReq domain.PaymentRequest) (domain.PaymentResult, error) {
	// In a real implementation, this would integrate with a payment gateway
	// For this example, we'll simulate a payment processor

	// Validate payment request
	if paymentReq.OrderID == "" {
		return domain.PaymentResult{}, errors.Join(domain.ErrInvalidOrderData, fmt.Errorf("order ID is required"))
	}
	if paymentReq.PaymentMethod == "" {
		return domain.PaymentResult{}, errors.Join(domain.ErrInvalidOrderData, fmt.Errorf("payment method is required"))
	}
	if paymentReq.Amount <= 0 {
		return domain.PaymentResult{}, errors.Join(domain.ErrInvalidOrderData, fmt.Errorf("payment amount must be positive"))
	}

	// Simulate payment processing
	paymentResult := domain.PaymentResult{
		Success:       true,
		TransactionID: uuid.New().String(),
		Status:        "completed",
		Message:       "Payment processed successfully",
	}

	// For certain payment methods, a redirect URL might be required
	if paymentReq.PaymentMethod == "paypal" {
		paymentResult.RedirectURL = fmt.Sprintf("https://paypal.com/checkout/%s", paymentResult.TransactionID)
	}

	// Update the order with payment information
	paymentInfo := domain.PaymentInfo{
		PaymentMethod: paymentReq.PaymentMethod,
		TransactionID: paymentResult.TransactionID,
		Status:        paymentResult.Status,
		Amount:        paymentReq.Amount,
		Currency:      paymentReq.Currency,
		PaymentDate:   time.Now().Format(time.RFC3339),
	}

	// Update the order status and payment info
	err := s.repo.UpdatePaymentInfo(ctx, paymentReq.OrderID, paymentInfo)
	if err != nil {
		return domain.PaymentResult{}, errors.Join(domain.ErrPaymentFailed, err)
	}

	// Update the order status to processing after payment
	err = s.repo.UpdateOrderStatus(ctx, paymentReq.OrderID, "processing")
	if err != nil {
		// Payment was processed but status update failed
		// This should be handled by a reconciliation process
		return paymentResult, errors.Join(domain.ErrPaymentFailed, err)
	}

	return paymentResult, nil
}

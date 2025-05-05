package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"checkout-service/internal/domain"
	"checkout-service/internal/event"
	"checkout-service/internal/utils"

	"github.com/google/uuid"
)

// OrderService implements domain.OrderService
type OrderService struct {
	repo         domain.OrderRepository
	eventEmitter *event.Emitter
}

// NewOrderService creates a new order service
func NewOrderService(repo domain.OrderRepository, eventEmitter *event.Emitter) domain.OrderService {
	return &OrderService{
		repo:         repo,
		eventEmitter: eventEmitter,
	}
}

// GetRepository returns the repository instance
// This is used for admin operations
func (s *OrderService) GetRepository() domain.OrderRepository {
	return s.repo
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	// Tạo ID cho đơn hàng nếu chưa có
	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	// Tạo mã đơn hàng nếu chưa có (VT-123456)
	if order.OrderNumber == "" {
		order.OrderNumber = fmt.Sprintf("VT-%s", utils.GenerateRandomString(6))
	}

	// Đặt thời gian tạo và cập nhật
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	// Đặt trạng thái mặc định dựa trên phương thức thanh toán
	if order.Status == "" {
		if order.PaymentInfo.PaymentMethod == "cod" {
			// Đơn hàng COD được chấp nhận ngay lập tức
			order.Status = "processing"
		} else {
			// Đơn hàng thanh toán online cần chờ thanh toán
			order.Status = "pending"
		}
	}

	// Đảm bảo payment_info có giá trị
	if order.PaymentInfo.Status == "" {
		// Trạng thái thanh toán luôn là pending ban đầu, bất kể phương thức thanh toán
		order.PaymentInfo.Status = "pending"
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

	// Save the order
	savedOrder, err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// Emit đơn hàng đã được tạo
	if s.eventEmitter != nil {
		// Chuẩn bị dữ liệu cho event
		items := make([]event.ItemData, 0, len(savedOrder.Items))
		for _, item := range savedOrder.Items {
			items = append(items, event.ItemData{
				ProductID: item.ProductID,
				Name:      item.Name,
				Quantity:  int(item.Quantity),
				Price:     item.Price,
			})
		}

		// Tạo event data
		eventData := event.OrderEventData{
			OrderID:       savedOrder.ID,
			OrderNumber:   savedOrder.OrderNumber,
			UserID:        savedOrder.UserID,
			UserEmail:     savedOrder.BillingInfo.Email,
			Status:        savedOrder.Status,
			PaymentMethod: savedOrder.PaymentInfo.PaymentMethod,
			Total:         savedOrder.Totals.Total,
			CreatedAt:     savedOrder.CreatedAt,
			Items:         items,
		}

		// Thêm thông tin giao hàng
		if savedOrder.ShippingInfo.ShipToDifferentAddress {
			eventData.ShippingName = fmt.Sprintf("%s %s", savedOrder.ShippingInfo.FirstName, savedOrder.ShippingInfo.LastName)
			eventData.ShippingAddress = savedOrder.ShippingInfo.Address
			eventData.ShippingPhone = savedOrder.BillingInfo.Phone
		} else {
			eventData.ShippingName = fmt.Sprintf("%s %s", savedOrder.BillingInfo.FirstName, savedOrder.BillingInfo.LastName)
			eventData.ShippingAddress = savedOrder.BillingInfo.Address
			eventData.ShippingPhone = savedOrder.BillingInfo.Phone
		}

		// Emit sự kiện order.created
		err = s.eventEmitter.EmitOrderCreated(eventData)
		if err != nil {
			log.Printf("Warning: Failed to emit order.created event: %v", err)
			// Tiếp tục xử lý ngay cả khi emit event thất bại
		}
	}

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
		// log.Printf("ERROR: ListOrdersByUserID called with empty userID")
		return nil, 0, errors.Join(domain.ErrInvalidOrderData, fmt.Errorf("user ID is required"))
	}

	// log.Printf("ListOrdersByUserID: Fetching orders for user %s (page=%d, pageSize=%d)", userID, page, pageSize)

	// Calculate skip for pagination
	skip := (page - 1) * pageSize
	if skip < 0 {
		skip = 0
	}

	orders, total, err := s.repo.ListOrdersByUserID(ctx, userID, skip, pageSize)
	if err != nil {
		// log.Printf("ERROR: Failed to fetch orders for user %s: %v", userID, err)
		return nil, 0, err
	}

	// log.Printf("SUCCESS: Found %d orders for user %s (total: %d)", len(orders), userID, total)
	// for i, order := range orders {
	// 	log.Printf("Order %d: ID=%s, OrderNumber=%s, Status=%s, CreatedAt=%v",
	// 		i+1, order.ID, order.OrderNumber, order.Status, order.CreatedAt)
	// }

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
		"paid":       true,
	}

	if !validStatuses[status] {
		return errors.Join(domain.ErrInvalidOrderData, fmt.Errorf("invalid status: %s", status))
	}

	// Lấy thông tin đơn hàng hiện tại để biết trạng thái trước
	currentOrder, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get current order: %w", err)
	}

	previousStatus := currentOrder.Status

	// Cập nhật trạng thái
	err = s.repo.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		return err
	}

	// Phát event status_changed
	if s.eventEmitter != nil {
		log.Printf("EVENT-DEBUG: Creating order.status_changed event for order %s", orderID)

		// Lấy thông tin đơn hàng đã cập nhật
		updatedOrder, err := s.repo.GetOrderByID(ctx, orderID)
		if err == nil && updatedOrder != nil {
			eventData := event.OrderEventData{
				OrderID:        updatedOrder.ID,
				OrderNumber:    updatedOrder.OrderNumber,
				UserID:         updatedOrder.UserID,
				UserEmail:      updatedOrder.BillingInfo.Email,
				Status:         status,
				PreviousStatus: previousStatus,
				PaymentMethod:  updatedOrder.PaymentInfo.PaymentMethod,
				Total:          updatedOrder.Totals.Total,
				CreatedAt:      updatedOrder.CreatedAt,
			}

			// Phát event không chặn
			go func() {
				if err := s.eventEmitter.EmitOrderStatusChanged(eventData); err != nil {
					log.Printf("EVENT-ERROR: Failed to emit order.status_changed event: %v", err)
				} else {
					log.Printf("EVENT-SUCCESS: Successfully emitted order.status_changed event for order ID %s", orderID)
				}
			}()
		} else {
			log.Printf("EVENT-ERROR: Could not get updated order data for status change event: %v", err)
		}
	} else {
		log.Printf("EVENT-ERROR: Cannot emit order.status_changed event - EventEmitter is nil")
	}

	return nil
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
	// Thêm debug log
	log.Printf("PAYMENT-DEBUG: ProcessPayment called with orderID=%s, method=%s, amount=%f",
		paymentReq.OrderID, paymentReq.PaymentMethod, paymentReq.Amount)

	// Check if eventEmitter is nil và log nó một cách rõ ràng
	if s.eventEmitter == nil {
		log.Printf("PAYMENT-ERROR: EventEmitter is nil in ProcessPayment! Events will not be sent!")
	} else {
		log.Printf("PAYMENT-DEBUG: EventEmitter is properly initialized")
	}

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

	// Xử lý đặc biệt cho COD (Cash on Delivery)
	if paymentReq.PaymentMethod == "cod" {
		log.Printf("PAYMENT-DEBUG: Processing COD payment for order %s", paymentReq.OrderID)

		// Với COD, thanh toán luôn thành công (sẽ thanh toán khi giao hàng)
		paymentResult := domain.PaymentResult{
			Success:       true,
			TransactionID: uuid.New().String(),
			Status:        "pending_delivery", // Trạng thái đặc biệt cho COD
			Message:       "COD order confirmed, payment will be collected upon delivery",
		}

		// Cập nhật thông tin thanh toán vào order
		paymentInfo := domain.PaymentInfo{
			PaymentMethod: paymentReq.PaymentMethod,
			TransactionID: paymentResult.TransactionID,
			Status:        paymentResult.Status, // Use consistent "pending_delivery" status
			Amount:        paymentReq.Amount,
			Currency:      paymentReq.Currency,
			PaymentDate:   time.Now().Format(time.RFC3339),
		}

		log.Printf("PAYMENT-DEBUG: Updating order %s with COD payment info", paymentReq.OrderID)

		// Cập nhật thông tin thanh toán
		err := s.repo.UpdatePaymentInfo(ctx, paymentReq.OrderID, paymentInfo)
		if err != nil {
			log.Printf("PAYMENT-ERROR: Failed to update payment info for COD: %v", err)
			return domain.PaymentResult{}, errors.Join(domain.ErrPaymentFailed, err)
		}

		// Cập nhật trạng thái đơn hàng thành đang xử lý
		log.Printf("PAYMENT-DEBUG: Updating order status to 'processing' for COD order")
		err = s.repo.UpdateOrderStatus(ctx, paymentReq.OrderID, "processing")
		if err != nil {
			log.Printf("PAYMENT-ERROR: Failed to update order status for COD: %v", err)
			return paymentResult, errors.Join(domain.ErrPaymentFailed, err)
		}

		// Phát event thanh toán thành công cho COD
		if s.eventEmitter != nil {
			log.Printf("PAYMENT-DEBUG: Emitting payment_succeeded event for COD order")
			// Lấy thông tin đơn hàng
			order, err := s.repo.GetOrderByID(ctx, paymentReq.OrderID)
			if err == nil && order != nil {
				eventData := event.PaymentEventData{
					OrderID:       order.ID,
					OrderNumber:   order.OrderNumber,
					UserID:        order.UserID,
					UserEmail:     order.BillingInfo.Email,
					PaymentMethod: paymentReq.PaymentMethod,
					Amount:        paymentReq.Amount,
					Currency:      paymentReq.Currency,
					TransactionID: paymentResult.TransactionID,
					Status:        "pending_delivery", // Use consistent status name
					PaymentDate:   time.Now().Format(time.RFC3339),
				}

				log.Printf("PAYMENT-DEBUG: COD event data created for order %s, email %s",
					order.ID, order.BillingInfo.Email)

				// Log debug dữ liệu chi tiết của event
				eventDataJSON, _ := json.Marshal(eventData)
				log.Printf("PAYMENT-DEBUG: COD event full data: %s", string(eventDataJSON))

				// Phát event ngay lập tức thay vì dùng goroutine - dễ theo dõi lỗi hơn
				err := s.eventEmitter.EmitPaymentSucceeded(eventData)
				if err != nil {
					log.Printf("PAYMENT-ERROR: Failed to emit COD payment_succeeded event: %v", err)
				} else {
					log.Printf("PAYMENT-SUCCESS: Successfully emitted COD payment_succeeded event")
				}
			} else {
				log.Printf("PAYMENT-ERROR: Failed to get order details for COD event: %v", err)
			}
		} else {
			log.Printf("PAYMENT-ERROR: Cannot emit COD payment event - EventEmitter is nil!")
		}

		return paymentResult, nil
	}

	// For online payments (MoMo, VNPay, etc.), we need to return a pending status since they'll be completed
	// via a callback after the user goes through the payment provider's flow
	paymentResult := domain.PaymentResult{
		Success:       true,
		TransactionID: uuid.New().String(),
		Status:        "pending", // Set as pending initially for online payments
		Message:       "Payment initiated. Please complete payment on the provider's site.",
	}

	// For certain payment methods, a redirect URL might be required
	if paymentReq.PaymentMethod == "paypal" {
		paymentResult.RedirectURL = fmt.Sprintf("https://paypal.com/checkout/%s", paymentResult.TransactionID)
	}

	// Update the order with payment information - keep status as pending
	paymentInfo := domain.PaymentInfo{
		PaymentMethod: paymentReq.PaymentMethod,
		TransactionID: paymentResult.TransactionID,
		Status:        "pending", // Keep as pending until callback confirms payment
		Amount:        paymentReq.Amount,
		Currency:      paymentReq.Currency,
		PaymentDate:   time.Now().Format(time.RFC3339),
	}

	log.Printf("PAYMENT-DEBUG: Updating order %s with payment info: method=%s, transID=%s, status=pending",
		paymentReq.OrderID, paymentInfo.PaymentMethod, paymentInfo.TransactionID)

	// Update the payment info but don't change order status yet
	err := s.repo.UpdatePaymentInfo(ctx, paymentReq.OrderID, paymentInfo)
	if err != nil {
		log.Printf("PAYMENT-DEBUG: Failed to update payment info: %v", err)
		return domain.PaymentResult{}, errors.Join(domain.ErrPaymentFailed, err)
	}

	// Don't emit a payment success event yet for online payments - wait for callback confirmation
	log.Printf("PAYMENT-DEBUG: Online payment initiated for order %s. Waiting for payment gateway callback.", paymentReq.OrderID)

	return paymentResult, nil
}

// GetUserTotalSpend retrieves the total amount spent by a user
func (s *OrderService) GetUserTotalSpend(ctx context.Context, userID string) (float64, error) {
	if userID == "" {
		return 0, errors.Join(domain.ErrInvalidOrderData, fmt.Errorf("user ID is required"))
	}

	log.Printf("GetUserTotalSpend: Calculating total spend for user %s", userID)

	totalSpend, err := s.repo.GetUserTotalSpend(ctx, userID)
	if err != nil {
		log.Printf("ERROR: Failed to get total spend for user %s: %v", userID, err)
		return 0, err
	}

	log.Printf("User %s has total spend of %.2f", userID, totalSpend)
	return totalSpend, nil
}

// GetUserOrderCount retrieves the number of orders placed by a user
func (s *OrderService) GetUserOrderCount(ctx context.Context, userID string) (int, error) {
	if userID == "" {
		return 0, errors.Join(domain.ErrInvalidOrderData, fmt.Errorf("user ID is required"))
	}

	log.Printf("GetUserOrderCount: Counting orders for user %s", userID)

	orderCount, err := s.repo.GetUserOrderCount(ctx, userID)
	if err != nil {
		log.Printf("ERROR: Failed to get order count for user %s: %v", userID, err)
		return 0, err
	}

	log.Printf("User %s has placed %d orders", userID, orderCount)
	return orderCount, nil
}

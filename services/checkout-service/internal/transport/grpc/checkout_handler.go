package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"checkout-service/internal/domain"
	pb "checkout-service/proto/checkout"
)

// CheckoutHandler implements the gRPC server interface
type CheckoutHandler struct {
	pb.UnimplementedCheckoutServiceServer
	orderService domain.OrderService
}

// NewCheckoutHandler creates a new gRPC handler
func NewCheckoutHandler(orderService domain.OrderService) *CheckoutHandler {
	return &CheckoutHandler{
		orderService: orderService,
	}
}

// CreateOrder handles creating a new order
func (h *CheckoutHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	// Debug: In ra thông tin request
	// log.Printf("CreateOrder Request received: %+v", req)
	// log.Printf("BillingInfo: %+v", req.BillingInfo)
	// log.Printf("ShippingInfo: %+v", req.ShippingInfo)
	log.Printf("EVENT-DEBUG: CreateOrder called with paymentMethod=%s", req.PaymentMethod)

	// Convert request to domain model
	items := make([]domain.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.OrderItem{
			ProductID: item.ProductId,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			ImageURL:  item.ImageUrl,
		})
	}

	orderData := &domain.Order{
		UserID: req.UserId,
		BillingInfo: domain.BillingInfo{
			FirstName:   req.BillingInfo.FirstName,
			LastName:    req.BillingInfo.LastName,
			CompanyName: req.BillingInfo.CompanyName,
			Address:     req.BillingInfo.Address,
			Country:     req.BillingInfo.Country,
			Region:      req.BillingInfo.Region,
			City:        req.BillingInfo.City,
			ZipCode:     req.BillingInfo.ZipCode,
			Email:       req.BillingInfo.Email,
			Phone:       req.BillingInfo.Phone,
		},
		ShippingInfo: domain.ShippingInfo{
			ShipToDifferentAddress: req.ShippingInfo.ShipToDifferentAddress,
			FirstName:              req.ShippingInfo.FirstName,
			LastName:               req.ShippingInfo.LastName,
			CompanyName:            req.ShippingInfo.CompanyName,
			Address:                req.ShippingInfo.Address,
			Country:                req.ShippingInfo.Country,
			Region:                 req.ShippingInfo.Region,
			City:                   req.ShippingInfo.City,
			ZipCode:                req.ShippingInfo.ZipCode,
			ShippingMethod:         req.ShippingInfo.ShippingMethod,
			ShippingCost:           req.ShippingInfo.ShippingCost,
		},
		PaymentInfo: domain.PaymentInfo{
			PaymentMethod: req.PaymentMethod,
		},
		Items:      items,
		CouponCode: req.CouponCode,
		Notes:      req.Notes,
	}

	// Debug: In ra domain model đã chuyển đổi
	// log.Printf("Domain Order created: %+v", orderData)
	// log.Printf("BillingInfo: %+v", orderData.BillingInfo)
	// log.Printf("ShippingInfo: %+v", orderData.ShippingInfo)
	// log.Printf("PaymentInfo: %+v", orderData.PaymentInfo)

	// Create order using service
	order, err := h.orderService.CreateOrder(ctx, orderData)
	if err != nil {
		log.Printf("EVENT-ERROR: Error creating order: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	// Debug: In ra kết quả từ service
	// log.Printf("Order created successfully: %+v", order)
	// log.Printf("BillingInfo after service: %+v", order.BillingInfo)
	// log.Printf("ShippingInfo after service: %+v", order.ShippingInfo)
	// log.Printf("PaymentInfo after service: %+v", order.PaymentInfo)

	// Convert domain model back to protobuf message
	protoOrder := convertOrderToProto(order)

	// Debug: In ra protobuf response
	// log.Printf("Proto response: %+v", protoOrder)
	// log.Printf("BillingInfo in response: %+v", protoOrder.BillingInfo)
	// log.Printf("ShippingInfo in response: %+v", protoOrder.ShippingInfo)
	// log.Printf("PaymentInfo in response: %+v", protoOrder.PaymentInfo)

	return protoOrder, nil
}

// GetOrder handles retrieving an existing order
func (h *CheckoutHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	order, err := h.orderService.GetOrderByID(ctx, req.Id, req.UserId)
	if err != nil {
		if err == domain.ErrOrderNotFound {
			return nil, status.Errorf(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}

	return convertOrderToProto(order), nil
}

// ListOrders handles retrieving a list of orders for a user
func (h *CheckoutHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	log.Printf("ListOrders gRPC handler called with UserId=%s, Page=%d, PageSize=%d", req.UserId, req.Page, req.PageSize)

	// Check if this is a request for all orders (admin case)
	isAdminRequest := req.UserId == ""
	if isAdminRequest {
		log.Printf("Admin request detected - will return all orders without user filtering")
	}

	// For normal user requests, verify user ID provided
	if !isAdminRequest && req.UserId == "" {
		log.Printf("ERROR: Empty user ID provided in non-admin ListOrders request")
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required for non-admin requests")
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10 // Default page size
		log.Printf("Using default page size: %d", pageSize)
	}

	page := int(req.Page)
	if page <= 0 {
		page = 1 // Default page
		log.Printf("Using default page: %d", page)
	}

	if isAdminRequest {
		// For admin - fetch all orders without user filtering
		log.Printf("Fetching all orders for admin (page=%d, pageSize=%d)", page, pageSize)

		// Calculate pagination parameters
		skip := (page - 1) * pageSize

		// Call a special method to get all orders
		orders, total, err := h.getAllOrders(ctx, skip, pageSize)
		if err != nil {
			log.Printf("ERROR: Failed to list all orders: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
		}

		log.Printf("Successfully retrieved %d orders for admin (total: %d)", len(orders), total)

		// Convert domain models to protobuf messages
		pbOrders := make([]*pb.Order, 0, len(orders))
		for i, order := range orders {
			protoOrder := convertOrderToProto(order)
			pbOrders = append(pbOrders, protoOrder)
			log.Printf("Converted order %d: ID=%s to proto message", i+1, order.ID)
		}

		response := &pb.ListOrdersResponse{
			Orders: pbOrders,
			Total:  int32(total),
		}
		return response, nil
	}

	// For normal users - get only their orders
	log.Printf("Calling orderService.ListOrdersByUserID with userID=%s, page=%d, pageSize=%d", req.UserId, page, pageSize)

	// Log recent orders for debugging - could be faster than a pagination query
	// Direct query to see most recent 5 orders in the system (for debugging only)
	allOrders, _, _ := h.orderService.ListOrdersByUserID(ctx, req.UserId, 1, 5)
	log.Printf("DEBUG: Most recent %d orders for user %s:", len(allOrders), req.UserId)
	for i, o := range allOrders {
		log.Printf("DEBUG: Order[%d]: ID=%s, OrderNumber=%s, Status=%s, CreatedAt=%v",
			i, o.ID, o.OrderNumber, o.Status, o.CreatedAt)
	}

	// Now execute the actual paginated query
	orders, total, err := h.orderService.ListOrdersByUserID(ctx, req.UserId, page, pageSize)
	if err != nil {
		log.Printf("ERROR: Failed to list orders: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}

	log.Printf("Successfully retrieved %d orders (total: %d)", len(orders), total)

	// Convert domain models to protobuf messages
	pbOrders := make([]*pb.Order, 0, len(orders))
	for i, order := range orders {
		protoOrder := convertOrderToProto(order)
		pbOrders = append(pbOrders, protoOrder)
		log.Printf("Converted order %d: ID=%s to proto message", i+1, order.ID)
	}

	response := &pb.ListOrdersResponse{
		Orders: pbOrders,
		Total:  int32(total),
	}
	log.Printf("Sending ListOrdersResponse with %d orders", len(pbOrders))

	return response, nil
}

// getAllOrders retrieves all orders regardless of user with pagination
// This is a helper method specifically for admin access
func (h *CheckoutHandler) getAllOrders(ctx context.Context, skip, limit int) ([]*domain.Order, int, error) {
	// This would typically call a method on the orderService, but we'll implement it directly
	// for now since we don't want to modify all the interfaces

	// We'll use direct repository access if available via the orderService
	// Check if orderService has a method to get the repository
	if repoAccessor, ok := h.orderService.(interface {
		GetRepository() domain.OrderRepository
	}); ok {
		repo := repoAccessor.GetRepository()
		// Now we can call a method on the repository to get all orders
		return repo.ListAllOrders(ctx, skip, limit)
	}

	// Fallback to a less efficient approach if repository is not directly accessible
	log.Printf("WARNING: Using inefficient approach to get all orders - consider adding GetRepository() to OrderService")

	// We'll return an empty list and error since this functionality should be properly
	// implemented in the domain layer
	return nil, 0, status.Errorf(codes.Unimplemented, "admin listing of all orders is not properly implemented")
}

// UpdateOrderStatus handles updating the status of an order
func (h *CheckoutHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.StatusResponse, error) {
	err := h.orderService.UpdateOrderStatus(ctx, req.OrderId, req.Status)
	if err != nil {
		if err == domain.ErrOrderNotFound {
			return nil, status.Errorf(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}

	return &pb.StatusResponse{
		Success: true,
		Message: "Order status updated successfully",
	}, nil
}

// ProcessPayment handles payment processing for an order
func (h *CheckoutHandler) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("EVENT-DEBUG: ProcessPayment gRPC handler called with orderID=%s, method=%s, amount=%f",
		req.OrderId, req.PaymentMethod, req.Amount)

	// If this is a COD payment, log it specifically
	if req.PaymentMethod == "cod" {
		log.Printf("EVENT-DEBUG: Processing COD payment for order %s", req.OrderId)
	}

	paymentInfo := domain.PaymentRequest{
		OrderID:       req.OrderId,
		PaymentMethod: req.PaymentMethod,
		Amount:        req.Amount,
		Currency:      req.Currency,
		ReturnURL:     req.ReturnUrl,
	}

	paymentResult, err := h.orderService.ProcessPayment(ctx, paymentInfo)
	if err != nil {
		log.Printf("EVENT-ERROR: Payment processing failed: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
	}

	log.Printf("EVENT-DEBUG: Payment processed successfully, transactionID=%s, status=%s",
		paymentResult.TransactionID, paymentResult.Status)

	return &pb.PaymentResponse{
		Success:       paymentResult.Success,
		TransactionId: paymentResult.TransactionID,
		Status:        paymentResult.Status,
		RedirectUrl:   paymentResult.RedirectURL,
		Message:       paymentResult.Message,
	}, nil
}

// ValidateCheckout handles validating checkout data
func (h *CheckoutHandler) ValidateCheckout(ctx context.Context, req *pb.ValidateCheckoutRequest) (*pb.ValidateCheckoutResponse, error) {
	// Convert request to domain model
	items := make([]domain.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.OrderItem{
			ProductID: item.ProductId,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			ImageURL:  item.ImageUrl,
		})
	}

	checkoutData := domain.CheckoutValidationRequest{
		Items: items,
		BillingInfo: domain.BillingInfo{
			FirstName:   req.BillingInfo.FirstName,
			LastName:    req.BillingInfo.LastName,
			CompanyName: req.BillingInfo.CompanyName,
			Address:     req.BillingInfo.Address,
			Country:     req.BillingInfo.Country,
			Region:      req.BillingInfo.Region,
			City:        req.BillingInfo.City,
			ZipCode:     req.BillingInfo.ZipCode,
			Email:       req.BillingInfo.Email,
			Phone:       req.BillingInfo.Phone,
		},
		ShippingInfo: domain.ShippingInfo{
			ShipToDifferentAddress: req.ShippingInfo.ShipToDifferentAddress,
			FirstName:              req.ShippingInfo.FirstName,
			LastName:               req.ShippingInfo.LastName,
			CompanyName:            req.ShippingInfo.CompanyName,
			Address:                req.ShippingInfo.Address,
			Country:                req.ShippingInfo.Country,
			Region:                 req.ShippingInfo.Region,
			City:                   req.ShippingInfo.City,
			ZipCode:                req.ShippingInfo.ZipCode,
			ShippingMethod:         req.ShippingInfo.ShippingMethod,
			ShippingCost:           req.ShippingInfo.ShippingCost,
		},
		PaymentMethod: req.PaymentMethod,
	}

	// Validate checkout data
	result, err := h.orderService.ValidateCheckout(ctx, checkoutData)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate checkout: %v", err)
	}

	// Convert domain validation errors to protobuf
	pbErrors := make([]*pb.ValidationError, 0, len(result.Errors))
	for _, err := range result.Errors {
		pbErrors = append(pbErrors, &pb.ValidationError{
			Field:   err.Field,
			Message: err.Message,
		})
	}

	return &pb.ValidateCheckoutResponse{
		Valid:  result.Valid,
		Errors: pbErrors,
	}, nil
}

// GetHealth returns the service health status
func (h *CheckoutHandler) GetHealth(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	}, nil
}

// GetUserTotalSpend handles retrieving the total amount spent by a user
func (h *CheckoutHandler) GetUserTotalSpend(ctx context.Context, req *pb.UserRequest) (*pb.TotalSpendResponse, error) {
	log.Printf("GetUserTotalSpend gRPC handler called with UserID=%s", req.UserId)

	if req.UserId == "" {
		log.Printf("ERROR: Empty user ID provided in GetUserTotalSpend request")
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
	}

	totalSpend, err := h.orderService.GetUserTotalSpend(ctx, req.UserId)
	if err != nil {
		log.Printf("ERROR: Failed to get total spend for user %s: %v", req.UserId, err)
		return nil, status.Errorf(codes.Internal, "failed to get total spend: %v", err)
	}

	log.Printf("Successfully retrieved total spend %.2f for user %s", totalSpend, req.UserId)
	return &pb.TotalSpendResponse{
		TotalSpend: totalSpend,
	}, nil
}

// GetUserOrderCount handles retrieving the number of orders placed by a user
func (h *CheckoutHandler) GetUserOrderCount(ctx context.Context, req *pb.UserRequest) (*pb.OrderCountResponse, error) {
	log.Printf("GetUserOrderCount gRPC handler called with UserID=%s", req.UserId)

	if req.UserId == "" {
		log.Printf("ERROR: Empty user ID provided in GetUserOrderCount request")
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
	}

	orderCount, err := h.orderService.GetUserOrderCount(ctx, req.UserId)
	if err != nil {
		log.Printf("ERROR: Failed to get order count for user %s: %v", req.UserId, err)
		return nil, status.Errorf(codes.Internal, "failed to get order count: %v", err)
	}

	log.Printf("Successfully retrieved order count %d for user %s", orderCount, req.UserId)
	return &pb.OrderCountResponse{
		OrderCount: int32(orderCount),
	}, nil
}

// Helper functions to convert between domain models and protobuf messages
func convertOrderToProto(order *domain.Order) *pb.Order {
	// Check if order is nil
	if order == nil {
		return nil
	}

	// Convert order items
	pbItems := make([]*pb.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		pbItems = append(pbItems, &pb.OrderItem{
			Id:        item.ID,
			ProductId: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Subtotal:  item.Subtotal,
			ImageUrl:  item.ImageURL,
		})
	}

	// Format timestamp fields as strings
	createdAt := order.CreatedAt.Format(time.RFC3339)
	updatedAt := order.UpdatedAt.Format(time.RFC3339)

	return &pb.Order{
		Id:          order.ID,
		UserId:      order.UserID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		Items:       pbItems,
		BillingInfo: &pb.BillingInfo{
			FirstName:   order.BillingInfo.FirstName,
			LastName:    order.BillingInfo.LastName,
			CompanyName: order.BillingInfo.CompanyName,
			Address:     order.BillingInfo.Address,
			Country:     order.BillingInfo.Country,
			Region:      order.BillingInfo.Region,
			City:        order.BillingInfo.City,
			ZipCode:     order.BillingInfo.ZipCode,
			Email:       order.BillingInfo.Email,
			Phone:       order.BillingInfo.Phone,
		},
		ShippingInfo: &pb.ShippingInfo{
			ShipToDifferentAddress: order.ShippingInfo.ShipToDifferentAddress,
			FirstName:              order.ShippingInfo.FirstName,
			LastName:               order.ShippingInfo.LastName,
			CompanyName:            order.ShippingInfo.CompanyName,
			Address:                order.ShippingInfo.Address,
			Country:                order.ShippingInfo.Country,
			Region:                 order.ShippingInfo.Region,
			City:                   order.ShippingInfo.City,
			ZipCode:                order.ShippingInfo.ZipCode,
			ShippingMethod:         order.ShippingInfo.ShippingMethod,
			ShippingCost:           order.ShippingInfo.ShippingCost,
		},
		PaymentInfo: &pb.PaymentInfo{
			PaymentMethod: order.PaymentInfo.PaymentMethod,
			TransactionId: order.PaymentInfo.TransactionID,
			Status:        order.PaymentInfo.Status,
			Amount:        order.PaymentInfo.Amount,
			Currency:      order.PaymentInfo.Currency,
			PaymentDate:   order.PaymentInfo.PaymentDate,
		},
		Totals: &pb.OrderTotals{
			Subtotal: order.Totals.Subtotal,
			Shipping: order.Totals.Shipping,
			Discount: order.Totals.Discount,
			Tax:      order.Totals.Tax,
			Total:    order.Totals.Total,
		},
		Notes:     order.Notes,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

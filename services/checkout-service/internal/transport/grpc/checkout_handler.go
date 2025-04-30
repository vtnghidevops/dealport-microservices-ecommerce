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

// CheckoutServiceHandler implements the gRPC server interface
type CheckoutServiceHandler struct {
	pb.UnimplementedCheckoutServiceServer
	orderService domain.OrderService
}

// NewCheckoutServiceHandler creates a new gRPC handler
func NewCheckoutServiceHandler(orderService domain.OrderService) *CheckoutServiceHandler {
	return &CheckoutServiceHandler{
		orderService: orderService,
	}
}

// CreateOrder handles creating a new order
func (h *CheckoutServiceHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	// Debug: In ra thông tin request
	log.Printf("CreateOrder Request received: %+v", req)
	log.Printf("BillingInfo: %+v", req.BillingInfo)
	log.Printf("ShippingInfo: %+v", req.ShippingInfo)
	log.Printf("PaymentMethod: %s", req.PaymentMethod)

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
	log.Printf("Domain Order created: %+v", orderData)
	log.Printf("BillingInfo: %+v", orderData.BillingInfo)
	log.Printf("ShippingInfo: %+v", orderData.ShippingInfo)
	log.Printf("PaymentInfo: %+v", orderData.PaymentInfo)

	// Create order using service
	order, err := h.orderService.CreateOrder(ctx, orderData)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	// Debug: In ra kết quả từ service
	log.Printf("Order created successfully: %+v", order)
	log.Printf("BillingInfo after service: %+v", order.BillingInfo)
	log.Printf("ShippingInfo after service: %+v", order.ShippingInfo)
	log.Printf("PaymentInfo after service: %+v", order.PaymentInfo)

	// Convert domain model back to protobuf message
	protoOrder := convertOrderToProto(order)

	// Debug: In ra protobuf response
	log.Printf("Proto response: %+v", protoOrder)
	log.Printf("BillingInfo in response: %+v", protoOrder.BillingInfo)
	log.Printf("ShippingInfo in response: %+v", protoOrder.ShippingInfo)
	log.Printf("PaymentInfo in response: %+v", protoOrder.PaymentInfo)

	return protoOrder, nil
}

// GetOrder handles retrieving an existing order
func (h *CheckoutServiceHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
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
func (h *CheckoutServiceHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	log.Printf("ListOrders gRPC handler called with UserId=%s, Page=%d, PageSize=%d", req.UserId, req.Page, req.PageSize)

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

	log.Printf("Calling orderService.ListOrdersByUserID with userID=%s, page=%d, pageSize=%d", req.UserId, page, pageSize)
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

// UpdateOrderStatus handles updating the status of an order
func (h *CheckoutServiceHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.StatusResponse, error) {
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
func (h *CheckoutServiceHandler) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.PaymentResponse, error) {
	paymentInfo := domain.PaymentRequest{
		OrderID:       req.OrderId,
		PaymentMethod: req.PaymentMethod,
		Amount:        req.Amount,
		Currency:      req.Currency,
		ReturnURL:     req.ReturnUrl,
	}

	paymentResult, err := h.orderService.ProcessPayment(ctx, paymentInfo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
	}

	return &pb.PaymentResponse{
		Success:       paymentResult.Success,
		TransactionId: paymentResult.TransactionID,
		Status:        paymentResult.Status,
		RedirectUrl:   paymentResult.RedirectURL,
		Message:       paymentResult.Message,
	}, nil
}

// ValidateCheckout handles validating checkout data
func (h *CheckoutServiceHandler) ValidateCheckout(ctx context.Context, req *pb.ValidateCheckoutRequest) (*pb.ValidateCheckoutResponse, error) {
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
func (h *CheckoutServiceHandler) GetHealth(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
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

package grpc

import (
	"cart-service/internal/domain"
	pb "cart-service/proto/cart"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the gRPC server for the cart service
type Server struct {
	pb.UnimplementedCartServiceServer
	cartService domain.CartService
}

// NewServer creates a new gRPC server
func NewServer(cartService domain.CartService) *Server {
	return &Server{
		cartService: cartService,
	}
}

// GetCart retrieves a cart for a user
func (s *Server) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.Cart, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Get cart from service
	cart, err := s.cartService.GetCart(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cart: %v", err)
	}

	// Convert domain cart to protobuf cart
	return domainCartToProto(cart), nil
}

// AddCartItem adds an item to the cart
func (s *Server) AddCartItem(ctx context.Context, req *pb.AddCartItemRequest) (*pb.Cart, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if req.Item == nil {
		return nil, status.Error(codes.InvalidArgument, "item is required")
	}

	// Convert protobuf item to domain item
	item := &domain.CartItem{
		ID:            uuid.New().String(),
		ProductID:     req.Item.ProductId,
		Name:          req.Item.Name,
		Price:         req.Item.Price,
		OriginalPrice: req.Item.OriginalPrice,
		Quantity:      int(req.Item.Quantity),
		ImageURL:      req.Item.ImageUrl,
	}

	// Add item to cart
	cart, err := s.cartService.AddCartItem(ctx, req.UserId, item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add item to cart: %v", err)
	}

	// Convert domain cart to protobuf cart
	return domainCartToProto(cart), nil
}

// UpdateCartItem updates an item in the cart
func (s *Server) UpdateCartItem(ctx context.Context, req *pb.UpdateCartItemRequest) (*pb.Cart, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if req.ItemId == "" {
		return nil, status.Error(codes.InvalidArgument, "item ID is required")
	}
	if req.Quantity < 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be non-negative")
	}

	// Update cart item
	cart, err := s.cartService.UpdateCartItem(ctx, req.UserId, req.ItemId, int(req.Quantity))
	if err != nil {
		if err == domain.ErrItemNotFound {
			return nil, status.Error(codes.NotFound, "item not found in cart")
		}
		if err == domain.ErrInvalidQuantity {
			return nil, status.Error(codes.InvalidArgument, "invalid quantity")
		}
		return nil, status.Errorf(codes.Internal, "failed to update cart item: %v", err)
	}

	// Convert domain cart to protobuf cart
	return domainCartToProto(cart), nil
}

// RemoveCartItem removes an item from the cart
func (s *Server) RemoveCartItem(ctx context.Context, req *pb.RemoveCartItemRequest) (*pb.Cart, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if req.ItemId == "" {
		return nil, status.Error(codes.InvalidArgument, "item ID is required")
	}

	// Remove cart item
	cart, err := s.cartService.RemoveCartItem(ctx, req.UserId, req.ItemId)
	if err != nil {
		if err == domain.ErrItemNotFound {
			return nil, status.Error(codes.NotFound, "item not found in cart")
		}
		return nil, status.Errorf(codes.Internal, "failed to remove cart item: %v", err)
	}

	// Convert domain cart to protobuf cart
	return domainCartToProto(cart), nil
}

// ClearCart removes all items from the cart
func (s *Server) ClearCart(ctx context.Context, req *pb.ClearCartRequest) (*pb.StatusResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Clear cart
	err := s.cartService.ClearCart(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to clear cart: %v", err)
	}

	// Return success response
	return &pb.StatusResponse{
		Success: true,
		Message: "Cart cleared successfully",
	}, nil
}

// ApplyCoupon applies a coupon to the cart
func (s *Server) ApplyCoupon(ctx context.Context, req *pb.ApplyCouponRequest) (*pb.Cart, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if req.CouponCode == "" {
		return nil, status.Error(codes.InvalidArgument, "coupon code is required")
	}

	// Apply coupon
	cart, err := s.cartService.ApplyCoupon(ctx, req.UserId, req.CouponCode)
	if err != nil {
		if err == domain.ErrCouponInvalid {
			return nil, status.Error(codes.InvalidArgument, "coupon code is invalid")
		}
		if err == domain.ErrCouponExpired {
			return nil, status.Error(codes.InvalidArgument, "coupon code has expired")
		}
		if err == domain.ErrInsufficientCart {
			return nil, status.Error(codes.FailedPrecondition, "cart total is below coupon minimum")
		}
		return nil, status.Errorf(codes.Internal, "failed to apply coupon: %v", err)
	}

	// Convert domain cart to protobuf cart
	return domainCartToProto(cart), nil
}

// RemoveCoupon removes a coupon from the cart
func (s *Server) RemoveCoupon(ctx context.Context, req *pb.RemoveCouponRequest) (*pb.Cart, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Remove coupon
	cart, err := s.cartService.RemoveCoupon(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove coupon: %v", err)
	}

	// Convert domain cart to protobuf cart
	return domainCartToProto(cart), nil
}

// GetHealth returns the health status of the service
func (s *Server) GetHealth(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Status:    "UP",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	}, nil
}

// domainCartToProto converts a domain cart to a protobuf cart
func domainCartToProto(cart *domain.Cart) *pb.Cart {
	if cart == nil {
		return nil
	}

	// Convert cart items
	items := make([]*pb.CartItem, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = &pb.CartItem{
			Id:            item.ID,
			ProductId:     item.ProductID,
			Name:          item.Name,
			Price:         item.Price,
			OriginalPrice: item.OriginalPrice,
			Quantity:      int32(item.Quantity),
			ImageUrl:      item.ImageURL,
		}
	}

	// Convert cart totals
	totals := &pb.CartTotals{
		Subtotal: cart.Totals.Subtotal,
		Shipping: cart.Totals.Shipping,
		Discount: cart.Totals.Discount,
		Tax:      cart.Totals.Tax,
		Total:    cart.Totals.Total,
	}

	// Return protobuf cart
	return &pb.Cart{
		Id:             cart.ID,
		UserId:         cart.UserID,
		Items:          items,
		Totals:         totals,
		CouponCode:     cart.CouponCode,
		DiscountAmount: cart.DiscountAmount,
		CreatedAt:      cart.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      cart.UpdatedAt.Format(time.RFC3339),
	}
}

// protoCartToDomain converts a protobuf cart to a domain cart
func protoCartToDomain(cart *pb.Cart) (*domain.Cart, error) {
	if cart == nil {
		return nil, fmt.Errorf("cart is nil")
	}

	// Parse timestamps
	createdAt, err := time.Parse(time.RFC3339, cart.CreatedAt)
	if err != nil {
		createdAt = time.Now()
	}
	updatedAt, err := time.Parse(time.RFC3339, cart.UpdatedAt)
	if err != nil {
		updatedAt = time.Now()
	}

	// Convert cart items
	items := make([]domain.CartItem, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = domain.CartItem{
			ID:            item.Id,
			ProductID:     item.ProductId,
			Name:          item.Name,
			Price:         item.Price,
			OriginalPrice: item.OriginalPrice,
			Quantity:      int(item.Quantity),
			ImageURL:      item.ImageUrl,
		}
	}

	// Convert cart totals
	totals := domain.CartTotals{
		Subtotal: cart.Totals.Subtotal,
		Shipping: cart.Totals.Shipping,
		Discount: cart.Totals.Discount,
		Tax:      cart.Totals.Tax,
		Total:    cart.Totals.Total,
	}

	// Return domain cart
	return &domain.Cart{
		ID:             cart.Id,
		UserID:         cart.UserId,
		Items:          items,
		Totals:         totals,
		CouponCode:     cart.CouponCode,
		DiscountAmount: cart.DiscountAmount,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}
 
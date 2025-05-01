package grpc

import (
	"cart-service/internal/domain"
	couponpb "cart-service/proto/coupon"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CouponServer implements the gRPC coupon service
type CouponServer struct {
	couponpb.UnimplementedCouponServiceServer
	couponService domain.CouponService
}

// NewCouponServer creates a new gRPC coupon server
func NewCouponServer(couponService domain.CouponService) *CouponServer {
	return &CouponServer{
		couponService: couponService,
	}
}

// GetCoupons returns a list of coupons with pagination
func (s *CouponServer) GetCoupons(ctx context.Context, req *couponpb.GetCouponsRequest) (*couponpb.GetCouponsResponse, error) {
	page := int(req.Page)
	if page < 1 {
		page = 1
	}

	limit := int(req.Limit)
	if limit < 1 || limit > 100 {
		limit = 10 // Default limit
	}

	coupons, total, err := s.couponService.GetCoupons(ctx, page, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get coupons: %v", err)
	}

	// Convert domain coupons to protobuf coupons
	pbCoupons := make([]*couponpb.Coupon, len(coupons))
	for i, coupon := range coupons {
		pbCoupons[i] = mapDomainCouponToProto(&coupon)
	}

	return &couponpb.GetCouponsResponse{
		Coupons: pbCoupons,
		Total:   int32(total),
		Page:    int32(page),
		Limit:   int32(limit),
	}, nil
}

// GetCouponByID returns a coupon by its ID
func (s *CouponServer) GetCouponByID(ctx context.Context, req *couponpb.GetCouponByIDRequest) (*couponpb.GetCouponResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Coupon ID is required")
	}

	coupon, err := s.couponService.GetCouponByID(ctx, req.Id)
	if err != nil {
		if err == domain.ErrCouponNotFound {
			return nil, status.Error(codes.NotFound, "Coupon not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get coupon: %v", err)
	}

	return &couponpb.GetCouponResponse{
		Coupon: mapDomainCouponToProto(coupon),
	}, nil
}

// GetCouponByCode returns a coupon by its code
func (s *CouponServer) GetCouponByCode(ctx context.Context, req *couponpb.GetCouponByCodeRequest) (*couponpb.GetCouponResponse, error) {
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "Coupon code is required")
	}

	coupon, err := s.couponService.GetCouponByCode(ctx, req.Code)
	if err != nil {
		if err == domain.ErrCouponNotFound {
			return nil, status.Error(codes.NotFound, "Coupon not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get coupon: %v", err)
	}

	return &couponpb.GetCouponResponse{
		Coupon: mapDomainCouponToProto(coupon),
	}, nil
}

// CreateCoupon creates a new coupon
func (s *CouponServer) CreateCoupon(ctx context.Context, req *couponpb.CreateCouponRequest) (*couponpb.GetCouponResponse, error) {
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "Coupon code is required")
	}

	if req.Discount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Discount must be greater than 0")
	}

	if req.DiscountType != "percentage" && req.DiscountType != "fixed" {
		return nil, status.Error(codes.InvalidArgument, "Discount type must be 'percentage' or 'fixed'")
	}

	// Parse time strings to time.Time
	validFrom := time.Now()
	validTo := validFrom.AddDate(1, 0, 0) // Default expiry is 1 year

	if req.ValidFrom != "" {
		var err error
		validFrom, err = time.Parse(time.RFC3339, req.ValidFrom)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid valid_from date format: %v", err)
		}
	}

	if req.ValidTo != "" {
		var err error
		validTo, err = time.Parse(time.RFC3339, req.ValidTo)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid valid_to date format: %v", err)
		}
	}

	// Create domain coupon from request
	coupon := &domain.Coupon{
		Code:           req.Code,
		Discount:       req.Discount,
		DiscountType:   req.DiscountType,
		MinOrderAmount: req.MinOrderAmount,
		MaxUsage:       int(req.MaxUsage),
		IsActive:       req.IsActive,
		Description:    req.Description,
		ValidFrom:      validFrom,
		ValidTo:        validTo,
	}

	createdCoupon, err := s.couponService.CreateCoupon(ctx, coupon)
	if err != nil {
		if err == domain.ErrCouponExists {
			return nil, status.Error(codes.AlreadyExists, "Coupon with this code already exists")
		}
		return nil, status.Errorf(codes.Internal, "Failed to create coupon: %v", err)
	}

	return &couponpb.GetCouponResponse{
		Coupon: mapDomainCouponToProto(createdCoupon),
	}, nil
}

// UpdateCoupon updates an existing coupon
func (s *CouponServer) UpdateCoupon(ctx context.Context, req *couponpb.UpdateCouponRequest) (*couponpb.GetCouponResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Coupon ID is required")
	}

	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "Coupon code is required")
	}

	if req.Discount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Discount must be greater than 0")
	}

	if req.DiscountType != "percentage" && req.DiscountType != "fixed" {
		return nil, status.Error(codes.InvalidArgument, "Discount type must be 'percentage' or 'fixed'")
	}

	// Parse time strings to time.Time if provided
	var validFrom time.Time
	var validTo time.Time

	if req.ValidFrom != "" {
		var err error
		validFrom, err = time.Parse(time.RFC3339, req.ValidFrom)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid valid_from date format: %v", err)
		}
	}

	if req.ValidTo != "" {
		var err error
		validTo, err = time.Parse(time.RFC3339, req.ValidTo)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid valid_to date format: %v", err)
		}
	}

	// Create domain coupon from request
	coupon := &domain.Coupon{
		ID:             req.Id,
		Code:           req.Code,
		Discount:       req.Discount,
		DiscountType:   req.DiscountType,
		MinOrderAmount: req.MinOrderAmount,
		MaxUsage:       int(req.MaxUsage),
		IsActive:       req.IsActive,
		Description:    req.Description,
		ValidFrom:      validFrom,
		ValidTo:        validTo,
	}

	updatedCoupon, err := s.couponService.UpdateCoupon(ctx, req.Id, coupon)
	if err != nil {
		if err == domain.ErrCouponNotFound {
			return nil, status.Error(codes.NotFound, "Coupon not found")
		}
		if err == domain.ErrCouponExists {
			return nil, status.Error(codes.AlreadyExists, "Coupon with this code already exists")
		}
		return nil, status.Errorf(codes.Internal, "Failed to update coupon: %v", err)
	}

	return &couponpb.GetCouponResponse{
		Coupon: mapDomainCouponToProto(updatedCoupon),
	}, nil
}

// DeleteCoupon deletes a coupon
func (s *CouponServer) DeleteCoupon(ctx context.Context, req *couponpb.DeleteCouponRequest) (*couponpb.DeleteCouponResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Coupon ID is required")
	}

	err := s.couponService.DeleteCoupon(ctx, req.Id)
	if err != nil {
		if err == domain.ErrCouponNotFound {
			return nil, status.Error(codes.NotFound, "Coupon not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to delete coupon: %v", err)
	}

	return &couponpb.DeleteCouponResponse{
		Success: true,
		Message: fmt.Sprintf("Coupon %s deleted successfully", req.Id),
	}, nil
}

// Helper function to map domain coupon to protobuf coupon
func mapDomainCouponToProto(coupon *domain.Coupon) *couponpb.Coupon {
	return &couponpb.Coupon{
		Id:             coupon.ID,
		Code:           coupon.Code,
		Discount:       coupon.Discount,
		DiscountType:   coupon.DiscountType,
		MinOrderAmount: coupon.MinOrderAmount,
		MaxUsage:       int32(coupon.MaxUsage),
		UsageCount:     int32(coupon.UsageCount),
		ValidFrom:      timestamppb.New(coupon.ValidFrom),
		ValidTo:        timestamppb.New(coupon.ValidTo),
		IsActive:       coupon.IsActive,
		Description:    coupon.Description,
		CreatedAt:      timestamppb.New(coupon.CreatedAt),
		UpdatedAt:      timestamppb.New(coupon.UpdatedAt),
	}
}

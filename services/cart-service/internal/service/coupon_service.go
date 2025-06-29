package service

import (
	"cart-service/internal/domain"
	"context"
	"math"
	"strings"
	"time"
)

// CouponService implements the domain.CouponService interface
type CouponService struct {
	// Repository for coupon data persistence
	couponRepo domain.CouponRepository
}

// NewCouponService creates a new CouponService
func NewCouponService(couponRepo domain.CouponRepository) *CouponService {
	return &CouponService{
		couponRepo: couponRepo,
	}
}

// ValidateCoupon checks if a coupon code is valid and returns the discount percentage
func (s *CouponService) ValidateCoupon(ctx context.Context, code string) (bool, float64, error) {
	// Normalize coupon code (uppercase)
	normalizedCode := strings.ToUpper(code)

	// Get the coupon from the repository
	coupon, err := s.couponRepo.GetCouponByCode(ctx, normalizedCode)
	if err != nil {
		if err == domain.ErrCouponNotFound {
			return false, 0, domain.ErrCouponInvalid
		}
		return false, 0, err
	}

	// Check if coupon is active
	if !coupon.IsActive {
		return false, 0, domain.ErrCouponInvalid
	}

	// Check if coupon has expired
	if time.Now().After(coupon.ValidTo) {
		return false, 0, domain.ErrCouponExpired
	}

	// Check if coupon has reached maximum usage
	if coupon.UsageCount >= coupon.MaxUsage {
		return false, 0, domain.ErrCouponInvalid
	}

	// Return discount based on type
	if coupon.DiscountType == "percentage" {
		return true, coupon.Discount / 100, nil // Convert 20% to 0.20
	}
	return true, coupon.Discount, nil // Fixed amount discount
}

// ApplyCoupon applies a coupon to the cart and returns the updated cart
func (s *CouponService) ApplyCoupon(ctx context.Context, cart *domain.Cart, code string) (*domain.Cart, error) {
	// Normalize coupon code (uppercase)
	normalizedCode := strings.ToUpper(code)

	// Validate coupon
	isValid, discountValue, err := s.ValidateCoupon(ctx, normalizedCode)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, domain.ErrCouponInvalid
	}

	// Get coupon from repository to check minimum order amount
	coupon, err := s.couponRepo.GetCouponByCode(ctx, normalizedCode)
	if err != nil && err != domain.ErrCouponNotFound {
		return nil, err
	}

	// If coupon not found in repository, return error (since we no longer use in-memory coupons)
	if coupon == nil {
		return nil, domain.ErrCouponNotFound
	}

	// Check for minimum cart value
	if cart.Totals.Subtotal < coupon.MinOrderAmount {
		return nil, domain.ErrInsufficientCart
	}

	// Store coupon code in cart
	cart.CouponCode = normalizedCode

	// Handle different types of coupons based on the discount type
	if coupon.DiscountType == "fixed" {
		// Fixed amount discount
		cart.DiscountAmount = coupon.Discount
	} else {
		// Percentage discount
		discountAmount := cart.Totals.Subtotal * discountValue
		// Round to 2 decimal places
		cart.DiscountAmount = math.Round(discountAmount*100) / 100
	}

	// Special handling for free shipping coupon (based on description or other criteria)
	if strings.Contains(strings.ToLower(coupon.Description), "free shipping") {
		cart.Totals.Shipping = "Free"
	}

	// Make sure discount is applied to cart totals
	cart.Totals.Discount = cart.DiscountAmount
	// Recalculate total with discount
	cart.Totals.Total = cart.Totals.Subtotal + cart.Totals.Tax - cart.DiscountAmount

	// Increment coupon usage
	if err := s.couponRepo.IncrementUsage(ctx, normalizedCode); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
		// log.Printf("Failed to increment coupon usage: %v", err)
	}

	// Return updated cart
	return cart, nil
}

// GetCoupons returns a paginated list of all coupons
func (s *CouponService) GetCoupons(ctx context.Context, page int, limit int) ([]domain.Coupon, int, error) {
	return s.couponRepo.GetCoupons(ctx, page, limit)
}

// GetCouponByID returns a coupon by its ID
func (s *CouponService) GetCouponByID(ctx context.Context, id string) (*domain.Coupon, error) {
	return s.couponRepo.GetCouponByID(ctx, id)
}

// GetCouponByCode returns a coupon by its code
func (s *CouponService) GetCouponByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	// Normalize code to uppercase
	normalizedCode := strings.ToUpper(code)
	return s.couponRepo.GetCouponByCode(ctx, normalizedCode)
}

// CreateCoupon creates a new coupon
func (s *CouponService) CreateCoupon(ctx context.Context, coupon *domain.Coupon) (*domain.Coupon, error) {
	// Normalize code to uppercase
	coupon.Code = strings.ToUpper(coupon.Code)

	// Validate coupon data
	if coupon.Code == "" {
		return nil, domain.ErrCouponInvalid
	}

	if coupon.Discount <= 0 {
		return nil, domain.ErrCouponInvalid
	}

	if coupon.DiscountType != "percentage" && coupon.DiscountType != "fixed" {
		return nil, domain.ErrCouponInvalid
	}

	if coupon.MaxUsage <= 0 {
		coupon.MaxUsage = 1000 // Default max usage if not specified
	}

	// Ensure valid date range
	now := time.Now()
	if coupon.ValidFrom.IsZero() {
		coupon.ValidFrom = now
	}

	if coupon.ValidTo.IsZero() {
		// Default expiry is 1 year from now if not specified
		coupon.ValidTo = now.AddDate(1, 0, 0)
	}

	if coupon.ValidFrom.After(coupon.ValidTo) {
		return nil, domain.ErrCouponInvalid
	}

	// Create coupon in repository
	return s.couponRepo.CreateCoupon(ctx, coupon)
}

// UpdateCoupon updates an existing coupon
func (s *CouponService) UpdateCoupon(ctx context.Context, id string, coupon *domain.Coupon) (*domain.Coupon, error) {
	// Ensure ID matches
	coupon.ID = id

	// Normalize code to uppercase
	coupon.Code = strings.ToUpper(coupon.Code)

	// Validate coupon data
	if coupon.Code == "" {
		return nil, domain.ErrCouponInvalid
	}

	if coupon.Discount <= 0 {
		return nil, domain.ErrCouponInvalid
	}

	if coupon.DiscountType != "percentage" && coupon.DiscountType != "fixed" {
		return nil, domain.ErrCouponInvalid
	}

	// Fetch existing coupon to retain some fields
	existing, err := s.couponRepo.GetCouponByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Ensure valid date range
	if coupon.ValidFrom.IsZero() {
		coupon.ValidFrom = existing.ValidFrom
	}

	if coupon.ValidTo.IsZero() {
		coupon.ValidTo = existing.ValidTo
	}

	if coupon.ValidFrom.After(coupon.ValidTo) {
		return nil, domain.ErrCouponInvalid
	}

	// Update coupon in repository
	return s.couponRepo.UpdateCoupon(ctx, coupon)
}

// DeleteCoupon deletes a coupon
func (s *CouponService) DeleteCoupon(ctx context.Context, id string) error {
	return s.couponRepo.DeleteCoupon(ctx, id)
}

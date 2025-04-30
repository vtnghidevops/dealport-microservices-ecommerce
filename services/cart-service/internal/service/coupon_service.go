package service

import (
	"cart-service/internal/domain"
	"context"
	"math"
	"strings"
)

// CouponService implements the domain.CouponService interface
type CouponService struct {
	// In a real implementation, this would be backed by a database
	// For now, we'll use a simple map of valid coupons and their discount percentages
	validCoupons map[string]float64
}

// NewCouponService creates a new CouponService
func NewCouponService() *CouponService {
	return &CouponService{
		validCoupons: map[string]float64{
			"WELCOME10":  0.10, // 10% discount
			"SUMMER20":   0.20, // 20% discount
			"FREESHIP":   0.00, // Free shipping (handled separately)
			"DISCOUNT50": 0.50, // 50% discount for testing
		},
	}
}

// ValidateCoupon checks if a coupon code is valid and returns the discount percentage
func (s *CouponService) ValidateCoupon(ctx context.Context, code string) (bool, float64, error) {
	// Normalize coupon code (uppercase)
	normalizedCode := strings.ToUpper(code)

	// Check if coupon exists
	discount, exists := s.validCoupons[normalizedCode]
	if !exists {
		return false, 0, ErrCouponInvalid
	}

	return true, discount, nil
}

// ApplyCoupon applies a coupon to the cart and returns the updated cart
func (s *CouponService) ApplyCoupon(ctx context.Context, cart *domain.Cart, code string) (*domain.Cart, error) {
	// Normalize coupon code (uppercase)
	normalizedCode := strings.ToUpper(code)

	// Validate coupon
	isValid, discountPercentage, err := s.ValidateCoupon(ctx, normalizedCode)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, ErrCouponInvalid
	}

	// Simple check for minimum cart value (example: $10 minimum for applying coupon)
	if cart.Totals.Subtotal < 10 && normalizedCode != "FREESHIP" {
		return nil, ErrInsufficientCart
	}

	// Store coupon code in cart
	cart.CouponCode = normalizedCode

	// Handle different types of coupons
	if normalizedCode == "FREESHIP" {
		// Free shipping coupon - no discount on subtotal
		cart.DiscountAmount = 0
		cart.Totals.Shipping = "Free"
	} else {
		// Regular percentage discount
		discountAmount := cart.Totals.Subtotal * discountPercentage
		// Round to 2 decimal places
		cart.DiscountAmount = math.Round(discountAmount*100) / 100
	}

	// Return updated cart
	return cart, nil
}

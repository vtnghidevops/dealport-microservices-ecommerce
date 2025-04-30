package domain

import (
	"context"
	"errors"
	"time"
)

// Error definitions
var (
	ErrItemNotFound     = errors.New("item not found in cart")
	ErrInvalidQuantity  = errors.New("invalid item quantity")
	ErrCouponInvalid    = errors.New("coupon code is invalid")
	ErrCouponExpired    = errors.New("coupon code has expired")
	ErrInsufficientCart = errors.New("cart total is below coupon minimum")
)

// CartItem represents a single item in a cart
type CartItem struct {
	ID            string  `json:"id"`
	ProductID     string  `json:"product_id"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"original_price"`
	Quantity      int     `json:"quantity"`
	ImageURL      string  `json:"image_url"`
}

// CartTotals represents the calculated totals for a cart
type CartTotals struct {
	Subtotal float64 `json:"subtotal"`
	Shipping string  `json:"shipping"`
	Discount float64 `json:"discount"`
	Tax      float64 `json:"tax"`
	Total    float64 `json:"total"`
}

// Cart represents a shopping cart
type Cart struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Items          []CartItem `json:"items"`
	Totals         CartTotals `json:"totals"`
	CouponCode     string     `json:"coupon_code"`
	DiscountAmount float64    `json:"discount_amount"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// CartRepository defines the interface for cart data persistence
type CartRepository interface {
	GetCart(ctx context.Context, userID string) (*Cart, error)
	SaveCart(ctx context.Context, cart *Cart) error
	DeleteCart(ctx context.Context, userID string) error
}

// CartService defines the interface for cart business logic
type CartService interface {
	GetCart(ctx context.Context, userID string) (*Cart, error)
	AddCartItem(ctx context.Context, userID string, item *CartItem) (*Cart, error)
	UpdateCartItem(ctx context.Context, userID, itemID string, quantity int) (*Cart, error)
	RemoveCartItem(ctx context.Context, userID, itemID string) (*Cart, error)
	ClearCart(ctx context.Context, userID string) error
	ApplyCoupon(ctx context.Context, userID, couponCode string) (*Cart, error)
	RemoveCoupon(ctx context.Context, userID string) (*Cart, error)
}

// CouponService defines the interface for coupon validation and application
type CouponService interface {
	ValidateCoupon(ctx context.Context, code string) (bool, float64, error)
	ApplyCoupon(ctx context.Context, cart *Cart, code string) (*Cart, error)
}


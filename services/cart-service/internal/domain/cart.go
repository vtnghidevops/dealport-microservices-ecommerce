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
	ErrCouponExists     = errors.New("coupon with this code already exists")
	ErrCouponNotFound   = errors.New("coupon not found")
)

// CartItem represents a single item in a cart
type CartItem struct {
	ID            string  `json:"id"`
	ProductID     string  `json:"productId"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"originalPrice"`
	Quantity      int     `json:"quantity"`
	ImageURL      string  `json:"imageUrl"`
}

// CartTotals represents the calculated totals for a cart
type CartTotals struct {
	Subtotal float64 `json:"subtotal"`
	Shipping string  `json:"shipping"`
	Discount float64 `json:"discount"`
	Tax      float64 `json:"tax"`
	Total    float64 `json:"total"`
}

// Coupon represents a discount coupon
type Coupon struct {
	ID             string    `json:"id"`
	Code           string    `json:"code"`
	Discount       float64   `json:"discount"`
	DiscountType   string    `json:"discountType"` // percentage or fixed
	MinOrderAmount float64   `json:"minOrderAmount"`
	MaxUsage       int       `json:"maxUsage"`
	UsageCount     int       `json:"usageCount"`
	ValidFrom      time.Time `json:"validFrom"`
	ValidTo        time.Time `json:"validTo"`
	IsActive       bool      `json:"isActive"`
	Description    string    `json:"description,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Cart represents a shopping cart
type Cart struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Items          []CartItem `json:"items"`
	Totals         CartTotals `json:"totals"`
	CouponCode     string     `json:"couponCode"`
	DiscountAmount float64    `json:"discountAmount"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

// CartRepository defines the interface for cart data persistence
type CartRepository interface {
	GetCart(ctx context.Context, userID string) (*Cart, error)
	SaveCart(ctx context.Context, cart *Cart) error
	DeleteCart(ctx context.Context, userID string) error
	RefreshCartTTL(ctx context.Context, userID string) error
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
	RefreshCartTTL(ctx context.Context, userID string) error
}

// CouponRepository defines the interface for coupon data persistence
type CouponRepository interface {
	GetCoupons(ctx context.Context, page int, limit int) ([]Coupon, int, error)
	GetCouponByID(ctx context.Context, id string) (*Coupon, error)
	GetCouponByCode(ctx context.Context, code string) (*Coupon, error)
	CreateCoupon(ctx context.Context, coupon *Coupon) (*Coupon, error)
	UpdateCoupon(ctx context.Context, coupon *Coupon) (*Coupon, error)
	DeleteCoupon(ctx context.Context, id string) error
	IncrementUsage(ctx context.Context, code string) error
}

// CouponService defines the interface for coupon validation and application
type CouponService interface {
	ValidateCoupon(ctx context.Context, code string) (bool, float64, error)
	ApplyCoupon(ctx context.Context, cart *Cart, code string) (*Cart, error)

	// New methods for coupon management
	GetCoupons(ctx context.Context, page int, limit int) ([]Coupon, int, error)
	GetCouponByID(ctx context.Context, id string) (*Coupon, error)
	GetCouponByCode(ctx context.Context, code string) (*Coupon, error)
	CreateCoupon(ctx context.Context, coupon *Coupon) (*Coupon, error)
	UpdateCoupon(ctx context.Context, id string, coupon *Coupon) (*Coupon, error)
	DeleteCoupon(ctx context.Context, id string) error
}

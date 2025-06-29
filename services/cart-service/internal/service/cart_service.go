package service

import (
	"cart-service/internal/domain"
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
)

var (
	ErrCartNotFound     = errors.New("cart not found")
	ErrItemNotFound     = errors.New("item not found in cart")
	ErrInvalidQuantity  = errors.New("invalid quantity")
	ErrCouponInvalid    = errors.New("coupon code is invalid")
	ErrCouponExpired    = errors.New("coupon code has expired")
	ErrInsufficientCart = errors.New("cart total is below coupon minimum")
)

// CartService implements the domain.CartService interface
type CartService struct {
	cartRepo      domain.CartRepository
	couponService domain.CouponService
}

// NewCartService creates a new CartService
func NewCartService(cartRepo domain.CartRepository, couponService domain.CouponService) *CartService {
	return &CartService{
		cartRepo:      cartRepo,
		couponService: couponService,
	}
}

// GetCart retrieves a cart for a given user
func (s *CartService) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Recalculate totals when getting the cart
	return s.calculateTotals(cart), nil
}

// AddCartItem adds a new item to the cart or updates quantity if it exists
func (s *CartService) AddCartItem(ctx context.Context, userID string, item *domain.CartItem) (*domain.Cart, error) {
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Generate a new ID if not provided
	if item.ID == "" {
		item.ID = uuid.New().String()
	}

	// Check if item already exists in cart
	found := false
	for i, cartItem := range cart.Items {
		if cartItem.ProductID == item.ProductID {
			// Update existing item quantity
			cart.Items[i].Quantity += item.Quantity
			found = true
			break
		}
	}

	// If item not found, add it to the cart
	if !found {
		cart.Items = append(cart.Items, *item)
	}

	// Recalculate totals
	cart = s.calculateTotals(cart)

	// Save updated cart
	if err := s.cartRepo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// UpdateCartItem updates the quantity of an item in the cart
func (s *CartService) UpdateCartItem(ctx context.Context, userID, itemID string, quantity int) (*domain.Cart, error) {
	if quantity < 0 {
		return nil, ErrInvalidQuantity
	}

	// Retrieve the original cart
	origCart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Make a deep copy of the cart to avoid modifying the original
	cart := &domain.Cart{
		ID:             origCart.ID,
		UserID:         origCart.UserID,
		Items:          make([]domain.CartItem, len(origCart.Items)),
		CouponCode:     origCart.CouponCode,
		DiscountAmount: origCart.DiscountAmount,
		Totals:         origCart.Totals,
		CreatedAt:      origCart.CreatedAt,
		UpdatedAt:      origCart.UpdatedAt,
	}

	// Copy each item
	copy(cart.Items, origCart.Items)

	// Find the item in the cart
	found := false
	for i, item := range cart.Items {
		if item.ID == itemID {
			if quantity == 0 {
				// Remove item from cart if quantity is 0
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			} else {
				// Update quantity
				cart.Items[i].Quantity = quantity
			}
			found = true
			break
		}
	}

	if !found {
		return nil, ErrItemNotFound
	}

	// Recalculate totals
	cart = s.calculateTotals(cart)

	// Save updated cart
	if err := s.cartRepo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// RemoveCartItem removes an item from the cart
func (s *CartService) RemoveCartItem(ctx context.Context, userID, itemID string) (*domain.Cart, error) {
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Find the item in the cart
	found := false
	for i, item := range cart.Items {
		if item.ID == itemID {
			// Remove item from cart
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return nil, ErrItemNotFound
	}

	// Recalculate totals
	cart = s.calculateTotals(cart)

	// Save updated cart
	if err := s.cartRepo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// ClearCart removes all items from the cart
func (s *CartService) ClearCart(ctx context.Context, userID string) error {
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	// Clear items
	cart.Items = []domain.CartItem{}

	// Clear coupon
	cart.CouponCode = ""
	cart.DiscountAmount = 0

	// Reset totals
	cart.Totals = domain.CartTotals{
		Subtotal: 0,
		Shipping: "Free",
		Discount: 0,
		Tax:      0,
		Total:    0,
	}

	// Save updated cart
	return s.cartRepo.SaveCart(ctx, cart)
}

// ApplyCoupon applies a coupon code to the cart
func (s *CartService) ApplyCoupon(ctx context.Context, userID, couponCode string) (*domain.Cart, error) {
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// First calculate totals without discount
	cart.CouponCode = ""
	cart.DiscountAmount = 0
	cart = s.calculateTotals(cart)

	// Apply coupon if the service is available
	if s.couponService != nil {
		cart, err = s.couponService.ApplyCoupon(ctx, cart, couponCode)
		if err != nil {
			return nil, err
		}
	} else {
		// Simple discount for testing if no coupon service
		cart.CouponCode = couponCode
		cart.DiscountAmount = cart.Totals.Subtotal * 0.1 // 10% discount
		cart = s.calculateTotals(cart)
	}

	// Save updated cart
	if err := s.cartRepo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// RemoveCoupon removes a coupon from the cart
func (s *CartService) RemoveCoupon(ctx context.Context, userID string) (*domain.Cart, error) {
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Remove coupon
	cart.CouponCode = ""
	cart.DiscountAmount = 0

	// Recalculate totals
	cart = s.calculateTotals(cart)

	// Save updated cart
	if err := s.cartRepo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// RefreshCartTTL refreshes the expiration time of a cart in Redis
// without modifying the cart data
func (s *CartService) RefreshCartTTL(ctx context.Context, userID string) error {
	// Use the repository's RefreshCartTTL method
	return s.cartRepo.RefreshCartTTL(ctx, userID)
}

// calculateTotals recalculates the cart totals
func (s *CartService) calculateTotals(cart *domain.Cart) *domain.Cart {
	// Calculate subtotal
	subtotal := 0.0
	for _, item := range cart.Items {
		subtotal += item.Price * float64(item.Quantity)
	}
	// Round subtotal to 2 decimal places
	subtotal = math.Round(subtotal*100) / 100

	// Set default shipping
	shipping := "Free"

	// Calculate tax (assuming 8% tax rate)
	taxRate := 0.08
	tax := subtotal * taxRate
	// Round tax to 2 decimal places
	tax = math.Round(tax*100) / 100

	// Apply discount if there's a coupon
	discount := cart.DiscountAmount
	// Round discount to 2 decimal places
	discount = math.Round(discount*100) / 100

	// Calculate total
	total := subtotal + tax - discount
	// Round total to 2 decimal places
	total = math.Round(total*100) / 100

	// Update cart totals
	cart.Totals = domain.CartTotals{
		Subtotal: subtotal,
		Shipping: shipping,
		Discount: discount,
		Tax:      tax,
		Total:    total,
	}

	return cart
}

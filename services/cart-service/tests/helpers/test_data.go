package helpers

import (
	"cart-service/internal/domain"
	"cart-service/internal/service"
	"cart-service/tests/mocks"
	"time"

	"github.com/google/uuid"
)

// CreateMockCartService tạo một cart service với các mocks repository
func CreateMockCartService() (*mocks.MockCartRepository, *mocks.MockCouponService, domain.CartService) {
	mockCartRepo := new(mocks.MockCartRepository)
	mockCouponService := new(mocks.MockCouponService)
	cartService := service.NewCartService(mockCartRepo, mockCouponService)
	return mockCartRepo, mockCouponService, cartService
}

// CreateEmptyCart tạo một giỏ hàng trống
func CreateEmptyCart(userID string) *domain.Cart {
	return &domain.Cart{
		ID:        uuid.New().String(),
		UserID:    userID,
		Items:     []domain.CartItem{},
		Totals: domain.CartTotals{
			Subtotal: 0,
			Shipping: "Free",
			Tax:      0,
			Total:    0,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateCartItem tạo một mục trong giỏ hàng
func CreateCartItem(productID, name string, price float64, quantity int) *domain.CartItem {
	return &domain.CartItem{
		ID:            uuid.New().String(),
		ProductID:     productID,
		Name:          name,
		Price:         price,
		OriginalPrice: price,
		Quantity:      quantity,
		ImageURL:      "https://example.com/images/" + productID + ".jpg",
	}
}

// CreateCartWithItems tạo một giỏ hàng với các sản phẩm
func CreateCartWithItems(userID string, items []domain.CartItem) *domain.Cart {
	cart := CreateEmptyCart(userID)
	cart.Items = items
	
	// Tính tổng giỏ hàng
	subtotal := 0.0
	for _, item := range items {
		subtotal += item.Price * float64(item.Quantity)
	}
	
	// Áp dụng phí vận chuyển và thuế
	tax := subtotal * 0.08 // Giả sử thuế 8%
	
	cart.Totals = domain.CartTotals{
		Subtotal: subtotal,
		Shipping: "Free",
		Tax:      tax,
		Total:    subtotal + tax,
	}
	
	return cart
}

// CreateSampleCoupon tạo một mã giảm giá mẫu
func CreateSampleCoupon(code string, discount float64, discountType string) *domain.Coupon {
	now := time.Now()
	return &domain.Coupon{
		ID:             uuid.New().String(),
		Code:           code,
		Discount:       discount,
		DiscountType:   discountType, // "percentage" hoặc "fixed"
		MinOrderAmount: 0,
		MaxUsage:       100,
		UsageCount:     0,
		ValidFrom:      now,
		ValidTo:        now.AddDate(0, 1, 0), // Hết hạn sau 1 tháng
		IsActive:       true,
		Description:    "Test coupon",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// CreateExpiredCoupon tạo một mã giảm giá đã hết hạn
func CreateExpiredCoupon(code string, discount float64, discountType string) *domain.Coupon {
	now := time.Now()
	return &domain.Coupon{
		ID:             uuid.New().String(),
		Code:           code,
		Discount:       discount,
		DiscountType:   discountType,
		MinOrderAmount: 0,
		MaxUsage:       100,
		UsageCount:     0,
		ValidFrom:      now.AddDate(0, -2, 0), // Bắt đầu 2 tháng trước
		ValidTo:        now.AddDate(0, -1, 0), // Hết hạn 1 tháng trước
		IsActive:       true,
		Description:    "Expired test coupon",
		CreatedAt:      now.AddDate(0, -2, 0),
		UpdatedAt:      now.AddDate(0, -2, 0),
	}
}

// CreateCouponWithMinimumOrder tạo mã giảm giá với yêu cầu đơn hàng tối thiểu
func CreateCouponWithMinimumOrder(code string, discount float64, discountType string, minOrderAmount float64) *domain.Coupon {
	coupon := CreateSampleCoupon(code, discount, discountType)
	coupon.MinOrderAmount = minOrderAmount
	return coupon
}

// ApplyCouponToCart áp dụng mã giảm giá vào giỏ hàng
func ApplyCouponToCart(cart *domain.Cart, coupon *domain.Coupon) *domain.Cart {
	// Tạo bản sao của giỏ hàng
	updatedCart := *cart
	
	// Áp dụng mã giảm giá
	updatedCart.CouponCode = coupon.Code
	
	// Tính giảm giá dựa trên loại
	if coupon.DiscountType == "percentage" {
		updatedCart.DiscountAmount = updatedCart.Totals.Subtotal * (coupon.Discount / 100)
	} else {
		updatedCart.DiscountAmount = coupon.Discount
	}
	
	// Cập nhật tổng
	updatedCart.Totals.Discount = updatedCart.DiscountAmount
	updatedCart.Totals.Total = updatedCart.Totals.Subtotal + updatedCart.Totals.Tax - updatedCart.DiscountAmount
	
	return &updatedCart
} 
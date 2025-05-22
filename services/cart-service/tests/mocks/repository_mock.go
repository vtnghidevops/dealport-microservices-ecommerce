package mocks

import (
	"cart-service/internal/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockCartRepository is a mock implementation of domain.CartRepository
type MockCartRepository struct {
	mock.Mock
}

// GetCart mocks the GetCart method
func (m *MockCartRepository) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	args := m.Called(ctx, userID)

	var cart *domain.Cart
	if args.Get(0) != nil {
		cart = args.Get(0).(*domain.Cart)
	}

	return cart, args.Error(1)
}

// SaveCart mocks the SaveCart method
func (m *MockCartRepository) SaveCart(ctx context.Context, cart *domain.Cart) error {
	args := m.Called(ctx, cart)
	return args.Error(0)
}

// DeleteCart mocks the DeleteCart method
func (m *MockCartRepository) DeleteCart(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockCouponRepository is a mock implementation of domain.CouponRepository
type MockCouponRepository struct {
	mock.Mock
}

// GetCoupons mocks the GetCoupons method
func (m *MockCouponRepository) GetCoupons(ctx context.Context, page int, limit int) ([]domain.Coupon, int, error) {
	args := m.Called(ctx, page, limit)

	var coupons []domain.Coupon
	if args.Get(0) != nil {
		coupons = args.Get(0).([]domain.Coupon)
	}

	return coupons, args.Int(1), args.Error(2)
}

// GetCouponByID mocks the GetCouponByID method
func (m *MockCouponRepository) GetCouponByID(ctx context.Context, id string) (*domain.Coupon, error) {
	args := m.Called(ctx, id)

	var coupon *domain.Coupon
	if args.Get(0) != nil {
		coupon = args.Get(0).(*domain.Coupon)
	}

	return coupon, args.Error(1)
}

// GetCouponByCode mocks the GetCouponByCode method
func (m *MockCouponRepository) GetCouponByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	args := m.Called(ctx, code)

	var coupon *domain.Coupon
	if args.Get(0) != nil {
		coupon = args.Get(0).(*domain.Coupon)
	}

	return coupon, args.Error(1)
}

// CreateCoupon mocks the CreateCoupon method
func (m *MockCouponRepository) CreateCoupon(ctx context.Context, coupon *domain.Coupon) (*domain.Coupon, error) {
	args := m.Called(ctx, coupon)

	var createdCoupon *domain.Coupon
	if args.Get(0) != nil {
		createdCoupon = args.Get(0).(*domain.Coupon)
	}

	return createdCoupon, args.Error(1)
}

// UpdateCoupon mocks the UpdateCoupon method
func (m *MockCouponRepository) UpdateCoupon(ctx context.Context, coupon *domain.Coupon) (*domain.Coupon, error) {
	args := m.Called(ctx, coupon)

	var updatedCoupon *domain.Coupon
	if args.Get(0) != nil {
		updatedCoupon = args.Get(0).(*domain.Coupon)
	}

	return updatedCoupon, args.Error(1)
}

// DeleteCoupon mocks the DeleteCoupon method
func (m *MockCouponRepository) DeleteCoupon(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// IncrementUsage mocks the IncrementUsage method
func (m *MockCouponRepository) IncrementUsage(ctx context.Context, code string) error {
	args := m.Called(ctx, code)
	return args.Error(0)
}

// MockCouponService is a mock implementation of domain.CouponService
type MockCouponService struct {
	mock.Mock
}

// ValidateCoupon mocks the ValidateCoupon method
func (m *MockCouponService) ValidateCoupon(ctx context.Context, code string) (bool, float64, error) {
	args := m.Called(ctx, code)
	return args.Bool(0), args.Get(1).(float64), args.Error(2)
}

// ApplyCoupon mocks the ApplyCoupon method
func (m *MockCouponService) ApplyCoupon(ctx context.Context, cart *domain.Cart, code string) (*domain.Cart, error) {
	args := m.Called(ctx, cart, code)

	var updatedCart *domain.Cart
	if args.Get(0) != nil {
		updatedCart = args.Get(0).(*domain.Cart)
	}

	return updatedCart, args.Error(1)
}

// GetCoupons mocks the GetCoupons method
func (m *MockCouponService) GetCoupons(ctx context.Context, page int, limit int) ([]domain.Coupon, int, error) {
	args := m.Called(ctx, page, limit)

	var coupons []domain.Coupon
	if args.Get(0) != nil {
		coupons = args.Get(0).([]domain.Coupon)
	}

	return coupons, args.Int(1), args.Error(2)
}

// GetCouponByID mocks the GetCouponByID method
func (m *MockCouponService) GetCouponByID(ctx context.Context, id string) (*domain.Coupon, error) {
	args := m.Called(ctx, id)

	var coupon *domain.Coupon
	if args.Get(0) != nil {
		coupon = args.Get(0).(*domain.Coupon)
	}

	return coupon, args.Error(1)
}

// GetCouponByCode mocks the GetCouponByCode method
func (m *MockCouponService) GetCouponByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	args := m.Called(ctx, code)

	var coupon *domain.Coupon
	if args.Get(0) != nil {
		coupon = args.Get(0).(*domain.Coupon)
	}

	return coupon, args.Error(1)
}

// CreateCoupon mocks the CreateCoupon method
func (m *MockCouponService) CreateCoupon(ctx context.Context, coupon *domain.Coupon) (*domain.Coupon, error) {
	args := m.Called(ctx, coupon)

	var createdCoupon *domain.Coupon
	if args.Get(0) != nil {
		createdCoupon = args.Get(0).(*domain.Coupon)
	}

	return createdCoupon, args.Error(1)
}

// UpdateCoupon mocks the UpdateCoupon method
func (m *MockCouponService) UpdateCoupon(ctx context.Context, id string, coupon *domain.Coupon) (*domain.Coupon, error) {
	args := m.Called(ctx, id, coupon)

	var updatedCoupon *domain.Coupon
	if args.Get(0) != nil {
		updatedCoupon = args.Get(0).(*domain.Coupon)
	}

	return updatedCoupon, args.Error(1)
}

// DeleteCoupon mocks the DeleteCoupon method
func (m *MockCouponService) DeleteCoupon(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

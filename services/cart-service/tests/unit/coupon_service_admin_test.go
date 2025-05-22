package unit

import (
	"cart-service/internal/domain"
	"cart-service/internal/service"
	"cart-service/tests/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGetCoupons tests the GetCoupons function
func TestGetCoupons(t *testing.T) {
	// Setup
	mockCouponRepo := new(mocks.MockCouponRepository)
	couponService := service.NewCouponService(mockCouponRepo)
	ctx := context.Background()

	// Create sample coupons
	coupons := []domain.Coupon{
		{
			ID:             "coupon1",
			Code:           "SUMMER20",
			Discount:       20.0,
			DiscountType:   "percentage",
			MinOrderAmount: 50.0,
			MaxUsage:       100,
			UsageCount:     10,
			IsActive:       true,
		},
		{
			ID:             "coupon2",
			Code:           "WELCOME10",
			Discount:       10.0,
			DiscountType:   "fixed",
			MinOrderAmount: 30.0,
			MaxUsage:       200,
			UsageCount:     50,
			IsActive:       true,
		},
	}

	// Test case 1: Successfully get coupons with pagination
	t.Run("Successfully get coupons", func(t *testing.T) {
		// Setup expectations
		mockCouponRepo.On("GetCoupons", ctx, 1, 10).Return(coupons, len(coupons), nil).Once()

		// Execute
		result, count, err := couponService.GetCoupons(ctx, 1, 10)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
		assert.Len(t, result, 2)
		assert.Equal(t, "SUMMER20", result[0].Code)
		assert.Equal(t, "WELCOME10", result[1].Code)

		// Verify expectations
		mockCouponRepo.AssertExpectations(t)
	})

	// Test case 2: No coupons found
	t.Run("No coupons found", func(t *testing.T) {
		// Setup expectations
		mockCouponRepo.On("GetCoupons", ctx, 2, 10).Return([]domain.Coupon{}, 0, nil).Once()

		// Execute
		result, count, err := couponService.GetCoupons(ctx, 2, 10)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
		assert.Empty(t, result)

		// Verify expectations
		mockCouponRepo.AssertExpectations(t)
	})

	// Test case 3: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Setup expectations
		mockCouponRepo.On("GetCoupons", ctx, 1, 20).Return(nil, 0, errors.New("database error")).Once()

		// Execute
		result, count, err := couponService.GetCoupons(ctx, 1, 20)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, 0, count)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockCouponRepo.AssertExpectations(t)
	})
}

// TestGetCouponByID tests the GetCouponByID function
func TestGetCouponByID(t *testing.T) {
	// Setup
	mockCouponRepo := new(mocks.MockCouponRepository)
	couponService := service.NewCouponService(mockCouponRepo)
	ctx := context.Background()

	// Create a sample coupon
	sampleCoupon := &domain.Coupon{
		ID:             "coupon1",
		Code:           "SUMMER20",
		Discount:       20.0,
		DiscountType:   "percentage",
		MinOrderAmount: 50.0,
		MaxUsage:       100,
		UsageCount:     10,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Test case 1: Successfully get coupon by ID
	t.Run("Successfully get coupon by ID", func(t *testing.T) {
		// Setup expectations
		mockCouponRepo.On("GetCouponByID", ctx, "coupon1").Return(sampleCoupon, nil).Once()

		// Execute
		result, err := couponService.GetCouponByID(ctx, "coupon1")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "coupon1", result.ID)
		assert.Equal(t, "SUMMER20", result.Code)
		assert.Equal(t, 20.0, result.Discount)

		// Verify expectations
		mockCouponRepo.AssertExpectations(t)
	})

	// Test case 2: Coupon not found
	t.Run("Coupon not found", func(t *testing.T) {
		// Setup expectations
		mockCouponRepo.On("GetCouponByID", ctx, "nonexistent").Return(nil, domain.ErrCouponNotFound).Once()

		// Execute
		result, err := couponService.GetCouponByID(ctx, "nonexistent")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrCouponNotFound, err)
		assert.Nil(t, result)

		// Verify expectations
		mockCouponRepo.AssertExpectations(t)
	})

	// Test case 3: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Setup expectations
		mockCouponRepo.On("GetCouponByID", ctx, "error").Return(nil, errors.New("database error")).Once()

		// Execute
		result, err := couponService.GetCouponByID(ctx, "error")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockCouponRepo.AssertExpectations(t)
	})
}

// TestGetCouponByCode tests the GetCouponByCode function
func TestGetCouponByCode(t *testing.T) {
	// Setup
	mockCouponRepo := new(mocks.MockCouponRepository)
	couponService := service.NewCouponService(mockCouponRepo)
	ctx := context.Background()

	// Create a sample coupon
	sampleCoupon := &domain.Coupon{
		ID:             "coupon1",
		Code:           "SUMMER20",
		Discount:       20.0,
		DiscountType:   "percentage",
		MinOrderAmount: 50.0,
		MaxUsage:       100,
		UsageCount:     10,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Test case 1: Successfully get coupon by code
	t.Run("Successfully get coupon by code", func(t *testing.T) {
		// Setup expectations - Notice code is uppercase in repo call
		mockCouponRepo.On("GetCouponByCode", ctx, "SUMMER20").Return(sampleCoupon, nil).Once()

		// Execute - Mixed case in service call
		result, err := couponService.GetCouponByCode(ctx, "Summer20")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "coupon1", result.ID)
		assert.Equal(t, "SUMMER20", result.Code)
		assert.Equal(t, 20.0, result.Discount)

		// Verify expectations
		mockCouponRepo.AssertExpectations(t)
	})

	// Test case 2: Coupon not found
	t.Run("Coupon not found", func(t *testing.T) {
		// Setup expectations
		mockCouponRepo.On("GetCouponByCode", ctx, "NONEXISTENT").Return(nil, domain.ErrCouponNotFound).Once()

		// Execute
		result, err := couponService.GetCouponByCode(ctx, "NONEXISTENT")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrCouponNotFound, err)
		assert.Nil(t, result)

		// Verify expectations
		mockCouponRepo.AssertExpectations(t)
	})
}

// TestUpdateCoupon tests the UpdateCoupon function
func TestUpdateCoupon(t *testing.T) {
	// Setup
	mockCouponRepo := new(mocks.MockCouponRepository)
	couponService := service.NewCouponService(mockCouponRepo)
	ctx := context.Background()

	// Create a sample coupon
	now := time.Now()
	updateCoupon := &domain.Coupon{
		ID:             "coupon1",
		Code:           "summer30", // lowercase for testing normalization
		Discount:       30.0,       // updated from 20 to 30
		DiscountType:   "percentage",
		MinOrderAmount: 60.0, // updated from 50 to 60
		MaxUsage:       150,  // updated from 100 to 150
		ValidFrom:      now.AddDate(0, -1, 0),
		ValidTo:        now.AddDate(0, 2, 0),
		IsActive:       true,
	}

	updatedCoupon := &domain.Coupon{
		ID:             "coupon1",
		Code:           "SUMMER30", // uppercase in result
		Discount:       30.0,
		DiscountType:   "percentage",
		MinOrderAmount: 60.0,
		MaxUsage:       150,
		UsageCount:     10, // unchanged
		ValidFrom:      now.AddDate(0, -1, 0),
		ValidTo:        now.AddDate(0, 2, 0),
		IsActive:       true,
		CreatedAt:      now.AddDate(0, -2, 0), // original creation date
		UpdatedAt:      now,                   // updated time
	}

	// Define the existing coupon outside the test cases
	existingCoupon := &domain.Coupon{
		ID:             "coupon1",
		Code:           "SUMMER20",
		Discount:       20.0,
		DiscountType:   "percentage",
		MinOrderAmount: 50.0,
		MaxUsage:       100,
		UsageCount:     10,
		ValidFrom:      now.AddDate(0, -2, 0),
		ValidTo:        now.AddDate(0, 1, 0),
		IsActive:       true,
		CreatedAt:      now.AddDate(0, -2, 0),
		UpdatedAt:      now.AddDate(0, -1, 0),
	}

	// Test case 1: Successfully update coupon
	t.Run("Successfully update coupon", func(t *testing.T) {
		// Setup expectations - Note that service transforms the coupon
		mockCouponRepo.On("GetCouponByID", ctx, "coupon1").Return(existingCoupon, nil).Once()
		mockCouponRepo.On("UpdateCoupon", ctx, mock.AnythingOfType("*domain.Coupon")).Return(updatedCoupon, nil).Once()

		// Execute
		result, err := couponService.UpdateCoupon(ctx, "coupon1", updateCoupon)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "coupon1", result.ID)
		assert.Equal(t, "SUMMER30", result.Code) // Code should be uppercase
		assert.Equal(t, 30.0, result.Discount)
		assert.Equal(t, 60.0, result.MinOrderAmount)
		assert.Equal(t, 150, result.MaxUsage)

		// Verify expectations
		mockCouponRepo.AssertExpectations(t)
	})

	// Test case 2: Invalid coupon data
	t.Run("Invalid coupon data - empty code", func(t *testing.T) {
		// Create an invalid coupon
		invalidCoupon := &domain.Coupon{
			ID:             "coupon1",
			Code:           "", // Empty code - invalid
			Discount:       30.0,
			DiscountType:   "percentage",
			MinOrderAmount: 60.0,
			IsActive:       true,
		}

		// Execute
		result, err := couponService.UpdateCoupon(ctx, "coupon1", invalidCoupon)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrCouponInvalid, err)
		assert.Nil(t, result)

		// Verify expectations - repo should not be called
		mockCouponRepo.AssertNotCalled(t, "UpdateCoupon")
	})

	// Test case 3: Invalid discount type
	t.Run("Invalid discount type", func(t *testing.T) {
		// Create a coupon with invalid discount type
		invalidTypeCoupon := &domain.Coupon{
			ID:             "coupon1",
			Code:           "SUMMER30",
			Discount:       30.0,
			DiscountType:   "invalid", // Neither "percentage" nor "fixed"
			MinOrderAmount: 60.0,
			IsActive:       true,
		}

		// Execute
		result, err := couponService.UpdateCoupon(ctx, "coupon1", invalidTypeCoupon)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrCouponInvalid, err)
		assert.Nil(t, result)

		// Verify expectations - repo should not be called
		mockCouponRepo.AssertNotCalled(t, "UpdateCoupon")
	})

	// Test case 4: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Setup expectations
		mockCouponRepo.On("GetCouponByID", ctx, "error").Return(existingCoupon, nil).Once()
		mockCouponRepo.On("UpdateCoupon", ctx, mock.AnythingOfType("*domain.Coupon")).Return(nil, errors.New("database error")).Once()

		// Execute
		result, err := couponService.UpdateCoupon(ctx, "error", updateCoupon)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockCouponRepo.AssertExpectations(t)
	})
}

// TestDeleteCoupon tests the DeleteCoupon function
func TestDeleteCoupon(t *testing.T) {
	// Setup
	mockCouponRepo := new(mocks.MockCouponRepository)
	couponService := service.NewCouponService(mockCouponRepo)
	ctx := context.Background()

	// Test case 1: Successfully delete coupon
	t.Run("Successfully delete coupon", func(t *testing.T) {
		// Setup expectations
		mockCouponRepo.On("DeleteCoupon", ctx, "coupon1").Return(nil).Once()

		// Execute
		err := couponService.DeleteCoupon(ctx, "coupon1")

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockCouponRepo.AssertExpectations(t)
	})

	// Test case 2: Coupon not found
	t.Run("Coupon not found", func(t *testing.T) {
		// Setup expectations
		mockCouponRepo.On("DeleteCoupon", ctx, "nonexistent").Return(domain.ErrCouponNotFound).Once()

		// Execute
		err := couponService.DeleteCoupon(ctx, "nonexistent")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrCouponNotFound, err)

		// Verify expectations
		mockCouponRepo.AssertExpectations(t)
	})

	// Test case 3: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Setup expectations
		mockCouponRepo.On("DeleteCoupon", ctx, "error").Return(errors.New("database error")).Once()

		// Execute
		err := couponService.DeleteCoupon(ctx, "error")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockCouponRepo.AssertExpectations(t)
	})
}

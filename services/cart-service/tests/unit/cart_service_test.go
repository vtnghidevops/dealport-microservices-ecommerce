package unit

import (
	"cart-service/internal/domain"
	"cart-service/internal/service"
	"cart-service/tests/mocks"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGetCart tests the GetCart function
func TestGetCart(t *testing.T) {
	// Setup
	mockCartRepo := new(mocks.MockCartRepository)
	mockCouponService := new(mocks.MockCouponService)

	// Create a sample cart
	sampleCart := &domain.Cart{
		ID:     "cart123",
		UserID: "user123",
		Items: []domain.CartItem{
			{
				ID:        "item1",
				ProductID: "prod1",
				Name:      "Product 1",
				Price:     19.99,
				Quantity:  2,
				ImageURL:  "http://example.com/product1.jpg",
			},
		},
		Totals: domain.CartTotals{
			Subtotal: 39.98,
			Shipping: "Free",
			Tax:      3.20,
			Total:    43.18,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Initialize the service
	cartService := service.NewCartService(mockCartRepo, mockCouponService)

	// Test case 1: Successfully get cart
	t.Run("Successfully get cart", func(t *testing.T) {
		// Setup expectations
		mockCartRepo.On("GetCart", mock.Anything, "user123").Return(sampleCart, nil).Once()

		// Execute
		cart, err := cartService.GetCart(context.Background(), "user123")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, cart)
		assert.Equal(t, "cart123", cart.ID)
		assert.Equal(t, "user123", cart.UserID)
		assert.Len(t, cart.Items, 1)
		assert.Equal(t, "Product 1", cart.Items[0].Name)

		// Verify expectations
		mockCartRepo.AssertExpectations(t)
	})

	// Test case 2: Error getting cart
	t.Run("Error getting cart", func(t *testing.T) {
		// Setup expectations
		mockCartRepo.On("GetCart", mock.Anything, "unknown_user").Return(nil, errors.New("cart not found")).Once()

		// Execute
		cart, err := cartService.GetCart(context.Background(), "unknown_user")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, cart)
		assert.Contains(t, err.Error(), "cart not found")

		// Verify expectations
		mockCartRepo.AssertExpectations(t)
	})
}

// TestAddCartItem tests the AddCartItem function
func TestAddCartItem(t *testing.T) {
	// Setup
	mockCartRepo := new(mocks.MockCartRepository)
	mockCouponService := new(mocks.MockCouponService)

	// Create a sample cart
	sampleCart := &domain.Cart{
		ID:     "cart123",
		UserID: "user123",
		Items:  []domain.CartItem{}, // Empty cart
		Totals: domain.CartTotals{
			Subtotal: 0,
			Shipping: "Free",
			Tax:      0,
			Total:    0,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create a sample item to add
	newItem := &domain.CartItem{
		ProductID: "prod1",
		Name:      "Product 1",
		Price:     19.99,
		Quantity:  2,
		ImageURL:  "http://example.com/product1.jpg",
	}

	// Initialize the service
	cartService := service.NewCartService(mockCartRepo, mockCouponService)

	// Test case 1: Add item to empty cart
	t.Run("Add item to empty cart", func(t *testing.T) {
		// Setup expectations
		mockCartRepo.On("GetCart", mock.Anything, "user123").Return(sampleCart, nil).Once()
		mockCartRepo.On("SaveCart", mock.Anything, mock.AnythingOfType("*domain.Cart")).Return(nil).Once()

		// Execute
		updatedCart, err := cartService.AddCartItem(context.Background(), "user123", newItem)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedCart)
		assert.Len(t, updatedCart.Items, 1)
		assert.Equal(t, "Product 1", updatedCart.Items[0].Name)
		assert.Equal(t, float64(19.99), updatedCart.Items[0].Price)
		assert.Equal(t, 2, updatedCart.Items[0].Quantity)
		assert.NotEmpty(t, updatedCart.Items[0].ID) // ID should be generated

		// Totals should be calculated
		assert.Equal(t, float64(39.98), updatedCart.Totals.Subtotal) // 19.99 * 2

		// Verify expectations
		mockCartRepo.AssertExpectations(t)
	})

	// Test case 2: Add item that already exists
	t.Run("Add existing item to cart", func(t *testing.T) {
		// Create a cart with an existing item
		cartWithItem := &domain.Cart{
			ID:     "cart123",
			UserID: "user123",
			Items: []domain.CartItem{
				{
					ID:        "item1",
					ProductID: "prod1", // Same product ID
					Name:      "Product 1",
					Price:     19.99,
					Quantity:  2,
					ImageURL:  "http://example.com/product1.jpg",
				},
			},
			Totals: domain.CartTotals{
				Subtotal: 39.98,
				Shipping: "Free",
				Tax:      3.20,
				Total:    43.18,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create the same item to add (same ProductID)
		sameItem := &domain.CartItem{
			ProductID: "prod1", // Same product ID
			Name:      "Product 1",
			Price:     19.99,
			Quantity:  3, // Different quantity
			ImageURL:  "http://example.com/product1.jpg",
		}

		// Setup expectations
		mockCartRepo.On("GetCart", mock.Anything, "user123").Return(cartWithItem, nil).Once()
		mockCartRepo.On("SaveCart", mock.Anything, mock.AnythingOfType("*domain.Cart")).Return(nil).Once()

		// Execute
		updatedCart, err := cartService.AddCartItem(context.Background(), "user123", sameItem)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedCart)
		assert.Len(t, updatedCart.Items, 1)               // Still just one item
		assert.Equal(t, 5, updatedCart.Items[0].Quantity) // 2 + 3 = 5

		// Totals should be recalculated
		assert.Equal(t, float64(99.95), updatedCart.Totals.Subtotal) // 19.99 * 5

		// Verify expectations
		mockCartRepo.AssertExpectations(t)
	})
}

// TestUpdateCartItem tests the UpdateCartItem function
func TestUpdateCartItem(t *testing.T) {
	// Setup
	mockCartRepo := new(mocks.MockCartRepository)
	mockCouponService := new(mocks.MockCouponService)

	// Initialize the service
	cartService := service.NewCartService(mockCartRepo, mockCouponService)

	// Test case 1: Update item quantity
	t.Run("Update item quantity", func(t *testing.T) {
		// Create a sample cart with items for this test case
		sampleCart := &domain.Cart{
			ID:     "cart123",
			UserID: "user123",
			Items: []domain.CartItem{
				{
					ID:        "item1",
					ProductID: "prod1",
					Name:      "Product 1",
					Price:     19.99,
					Quantity:  2,
					ImageURL:  "http://example.com/product1.jpg",
				},
				{
					ID:        "item2",
					ProductID: "prod2",
					Name:      "Product 2",
					Price:     29.99,
					Quantity:  1,
					ImageURL:  "http://example.com/product2.jpg",
				},
			},
			Totals: domain.CartTotals{
				Subtotal: 69.97, // (19.99 * 2) + 29.99
				Shipping: "Free",
				Tax:      5.60,
				Total:    75.57,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Setup expectations
		mockCartRepo.On("GetCart", mock.Anything, "user123").Return(sampleCart, nil).Once()
		mockCartRepo.On("SaveCart", mock.Anything, mock.AnythingOfType("*domain.Cart")).Return(nil).Once()

		// Execute - update quantity of first item from 2 to 4
		updatedCart, err := cartService.UpdateCartItem(context.Background(), "user123", "item1", 4)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedCart)
		assert.Len(t, updatedCart.Items, 2)

		// Find the updated item
		var updatedItem *domain.CartItem
		for _, item := range updatedCart.Items {
			if item.ID == "item1" {
				updatedItem = &item
				break
			}
		}

		assert.NotNil(t, updatedItem)
		assert.Equal(t, 4, updatedItem.Quantity)

		// Totals should be recalculated
		expectedSubtotal := (19.99 * 4) + 29.99 // 109.95
		assert.InDelta(t, expectedSubtotal, updatedCart.Totals.Subtotal, 0.01)

		// Verify expectations
		mockCartRepo.AssertExpectations(t)
	})

	// Test case 2: Remove item by setting quantity to 0
	t.Run("Remove item with quantity 0", func(t *testing.T) {
		fmt.Println("RUNNING TEST: Remove item with quantity 0")

		// Create a fresh sample cart with items for this test case
		sampleCart := &domain.Cart{
			ID:     "cart123",
			UserID: "user123",
			Items: []domain.CartItem{
				{
					ID:        "item1",
					ProductID: "prod1",
					Name:      "Product 1",
					Price:     19.99,
					Quantity:  2,
					ImageURL:  "http://example.com/product1.jpg",
				},
				{
					ID:        "item2",
					ProductID: "prod2",
					Name:      "Product 2",
					Price:     29.99,
					Quantity:  1,
					ImageURL:  "http://example.com/product2.jpg",
				},
			},
			Totals: domain.CartTotals{
				Subtotal: 69.97, // (19.99 * 2) + 29.99
				Shipping: "Free",
				Tax:      5.60,
				Total:    75.57,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Setup expectations with more strict checking
		mockCartRepo.On("GetCart", mock.Anything, "user123").Return(sampleCart, nil).Once()

		// Capture the argument passed to SaveCart
		var savedCart *domain.Cart
		mockCartRepo.On("SaveCart", mock.Anything, mock.MatchedBy(func(cart *domain.Cart) bool {
			savedCart = cart
			fmt.Printf("SaveCart called with cart having %d items\n", len(cart.Items))
			for i, item := range cart.Items {
				fmt.Printf("  Item %d: ID=%s, Name=%s\n", i, item.ID, item.Name)
			}
			return true
		})).Return(nil).Once()

		// Execute - set quantity of second item to 0 to remove it
		updatedCart, err := cartService.UpdateCartItem(context.Background(), "user123", "item2", 0)

		// Debug print statements
		fmt.Printf("Error: %v\n", err)
		if updatedCart != nil {
			fmt.Printf("Items in returned cart: %d\n", len(updatedCart.Items))
			for i, item := range updatedCart.Items {
				fmt.Printf("  Item %d: ID=%s, ProductID=%s, Name=%s\n", i, item.ID, item.ProductID, item.Name)
			}
			fmt.Printf("Subtotal: %.2f\n", updatedCart.Totals.Subtotal)
		} else {
			fmt.Printf("Updated cart is nil\n")
		}

		if savedCart != nil {
			fmt.Printf("SaveCart received cart with %d items\n", len(savedCart.Items))
			for i, item := range savedCart.Items {
				fmt.Printf("  Saved Item %d: ID=%s, ProductID=%s, Name=%s\n", i, item.ID, item.ProductID, item.Name)
			}
		} else {
			fmt.Printf("savedCart is nil\n")
		}

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedCart)

		// Print detailed error message if length assertion fails
		if len(updatedCart.Items) != 1 {
			fmt.Printf("ASSERTION FAILED: Expected 1 item, got %d items\n", len(updatedCart.Items))
			for i, item := range updatedCart.Items {
				fmt.Printf("  Item %d: ID=%s, Name=%s\n", i, item.ID, item.Name)
			}
		}

		assert.Len(t, updatedCart.Items, 1, "Cart should have exactly 1 item after removing item2")

		// Check if the remaining item is the correct one
		if len(updatedCart.Items) > 0 {
			if updatedCart.Items[0].ID != "item1" {
				fmt.Printf("ASSERTION FAILED: Expected item1, got %s\n", updatedCart.Items[0].ID)
			}
			assert.Equal(t, "item1", updatedCart.Items[0].ID, "The remaining item should be item1")
		}

		// Totals should be recalculated
		expectedSubtotal := 19.99 * 2 // Only the first item remains
		assert.InDelta(t, expectedSubtotal, updatedCart.Totals.Subtotal, 0.01, "Subtotal should be recalculated correctly")

		// Verify expectations
		mockCartRepo.AssertExpectations(t)
	})

	// Test case 3: Item not found
	t.Run("Item not found", func(t *testing.T) {
		// Create a fresh sample cart with items for this test case
		sampleCart := &domain.Cart{
			ID:     "cart123",
			UserID: "user123",
			Items: []domain.CartItem{
				{
					ID:        "item1",
					ProductID: "prod1",
					Name:      "Product 1",
					Price:     19.99,
					Quantity:  2,
					ImageURL:  "http://example.com/product1.jpg",
				},
				{
					ID:        "item2",
					ProductID: "prod2",
					Name:      "Product 2",
					Price:     29.99,
					Quantity:  1,
					ImageURL:  "http://example.com/product2.jpg",
				},
			},
			Totals: domain.CartTotals{
				Subtotal: 69.97, // (19.99 * 2) + 29.99
				Shipping: "Free",
				Tax:      5.60,
				Total:    75.57,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Setup expectations
		mockCartRepo.On("GetCart", mock.Anything, "user123").Return(sampleCart, nil).Once()

		// Execute with non-existent item ID
		updatedCart, err := cartService.UpdateCartItem(context.Background(), "user123", "nonexistent", 3)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, updatedCart)
		assert.Contains(t, err.Error(), "item not found")

		// Verify expectations
		mockCartRepo.AssertExpectations(t)
	})

	// Test case 4: Invalid quantity
	t.Run("Invalid quantity", func(t *testing.T) {
		// Execute with negative quantity
		updatedCart, err := cartService.UpdateCartItem(context.Background(), "user123", "item1", -1)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, updatedCart)
		assert.Contains(t, err.Error(), "invalid quantity")

		// Verify expectations - GetCart should not be called for invalid quantity
		mockCartRepo.AssertNotCalled(t, "GetCart")
	})
}

// TestRemoveCartItem tests the RemoveCartItem function
func TestRemoveCartItem(t *testing.T) {
	// Setup
	mockCartRepo := new(mocks.MockCartRepository)
	mockCouponService := new(mocks.MockCouponService)

	// Create a sample cart with items
	sampleCart := &domain.Cart{
		ID:     "cart123",
		UserID: "user123",
		Items: []domain.CartItem{
			{
				ID:        "item1",
				ProductID: "prod1",
				Name:      "Product 1",
				Price:     19.99,
				Quantity:  2,
				ImageURL:  "http://example.com/product1.jpg",
			},
			{
				ID:        "item2",
				ProductID: "prod2",
				Name:      "Product 2",
				Price:     29.99,
				Quantity:  1,
				ImageURL:  "http://example.com/product2.jpg",
			},
		},
		Totals: domain.CartTotals{
			Subtotal: 69.97, // (19.99 * 2) + 29.99
			Shipping: "Free",
			Tax:      5.60,
			Total:    75.57,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Initialize the service
	cartService := service.NewCartService(mockCartRepo, mockCouponService)

	// Test case 1: Successfully remove an item
	t.Run("Successfully remove item", func(t *testing.T) {
		// Setup expectations
		mockCartRepo.On("GetCart", mock.Anything, "user123").Return(sampleCart, nil).Once()
		mockCartRepo.On("SaveCart", mock.Anything, mock.AnythingOfType("*domain.Cart")).Return(nil).Once()

		// Execute - remove the first item
		updatedCart, err := cartService.RemoveCartItem(context.Background(), "user123", "item1")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedCart)
		assert.Len(t, updatedCart.Items, 1)
		assert.Equal(t, "item2", updatedCart.Items[0].ID) // Second item should remain

		// Totals should be recalculated
		expectedSubtotal := 29.99 // Only the second item remains
		assert.InDelta(t, expectedSubtotal, updatedCart.Totals.Subtotal, 0.01)

		// Verify expectations
		mockCartRepo.AssertExpectations(t)
	})

	// Test case 2: Item not found
	t.Run("Item not found", func(t *testing.T) {
		// Setup expectations
		mockCartRepo.On("GetCart", mock.Anything, "user123").Return(sampleCart, nil).Once()

		// Execute with non-existent item ID
		updatedCart, err := cartService.RemoveCartItem(context.Background(), "user123", "nonexistent")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, updatedCart)
		assert.Contains(t, err.Error(), "item not found")

		// Verify expectations
		mockCartRepo.AssertExpectations(t)
	})
}

// TestClearCart tests the ClearCart function
func TestClearCart(t *testing.T) {
	// Setup
	mockCartRepo := new(mocks.MockCartRepository)
	mockCouponService := new(mocks.MockCouponService)

	// Create a sample cart with items
	sampleCart := &domain.Cart{
		ID:     "cart123",
		UserID: "user123",
		Items: []domain.CartItem{
			{
				ID:        "item1",
				ProductID: "prod1",
				Name:      "Product 1",
				Price:     19.99,
				Quantity:  2,
				ImageURL:  "http://example.com/product1.jpg",
			},
		},
		CouponCode:     "DISCOUNT10",
		DiscountAmount: 4.00,
		Totals: domain.CartTotals{
			Subtotal: 39.98,
			Shipping: "Free",
			Discount: 4.00,
			Tax:      2.88,
			Total:    38.86,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Initialize the service
	cartService := service.NewCartService(mockCartRepo, mockCouponService)

	// Test case: Successfully clear cart
	t.Run("Successfully clear cart", func(t *testing.T) {
		// Setup expectations
		mockCartRepo.On("GetCart", mock.Anything, "user123").Return(sampleCart, nil).Once()
		mockCartRepo.On("SaveCart", mock.Anything, mock.AnythingOfType("*domain.Cart")).Return(nil).Once()

		// Execute
		err := cartService.ClearCart(context.Background(), "user123")

		// Assert
		assert.NoError(t, err)

		// Verify the saved cart had empty items and reset totals
		// Need to capture the cart that was passed to SaveCart
		mockCartRepo.AssertCalled(t, "SaveCart", mock.Anything, mock.MatchedBy(func(cart *domain.Cart) bool {
			return len(cart.Items) == 0 &&
				cart.CouponCode == "" &&
				cart.DiscountAmount == 0 &&
				cart.Totals.Subtotal == 0 &&
				cart.Totals.Total == 0
		}))

		// Verify expectations
		mockCartRepo.AssertExpectations(t)
	})
}

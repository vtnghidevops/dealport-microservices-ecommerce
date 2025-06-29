package security

import (
	"cart-service/internal/domain"
	"testing"
	"time"
)

// TestCartItemValidation tests validation of cart item inputs for security issues
func TestCartItemValidation(t *testing.T) {
	// Test cases for various security threats in cart item inputs
	tests := []struct {
		name    string
		item    domain.CartItem
		isValid bool
		reason  string
	}{
		{
			name: "Valid cart item",
			item: domain.CartItem{
				ID:        "item1",
				ProductID: "prod1",
				Name:      "Valid Product",
				Price:     19.99,
				Quantity:  2,
				ImageURL:  "https://example.com/image.jpg",
			},
			isValid: true,
			reason:  "",
		},
		{
			name: "Negative price",
			item: domain.CartItem{
				ID:        "item2",
				ProductID: "prod2",
				Name:      "Product with negative price",
				Price:     -10.50,
				Quantity:  1,
				ImageURL:  "https://example.com/image.jpg",
			},
			isValid: false,
			reason:  "negative price",
		},
		{
			name: "Excessive quantity",
			item: domain.CartItem{
				ID:        "item3",
				ProductID: "prod3",
				Name:      "Product with excessive quantity",
				Price:     25.99,
				Quantity:  10000, // Unreasonably high quantity
				ImageURL:  "https://example.com/image.jpg",
			},
			isValid: false,
			reason:  "excessive quantity",
		},
		{
			name: "XSS in product name",
			item: domain.CartItem{
				ID:        "item4",
				ProductID: "prod4",
				Name:      "Product <script>alert('XSS')</script>",
				Price:     15.99,
				Quantity:  1,
				ImageURL:  "https://example.com/image.jpg",
			},
			isValid: false,
			reason:  "contains script tags",
		},
		{
			name: "Malicious image URL",
			item: domain.CartItem{
				ID:        "item5",
				ProductID: "prod5",
				Name:      "Product with malicious URL",
				Price:     29.99,
				Quantity:  1,
				ImageURL:  "javascript:alert('XSS')",
			},
			isValid: false,
			reason:  "invalid image URL",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid, reason := validateCartItem(&test.item)
			if isValid != test.isValid {
				t.Errorf("Expected validity %v but got %v", test.isValid, isValid)
			}
			if !isValid && test.reason != "" && reason == "" {
				t.Errorf("Expected a reason for invalid item, but got none")
			}
		})
	}
}

// TestCouponCodeValidation tests validation of coupon codes for security issues
func TestCouponCodeValidation(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		isValid bool
		reason  string
	}{
		{
			name:    "Valid coupon code",
			code:    "SUMMER25",
			isValid: true,
			reason:  "",
		},
		{
			name:    "Coupon code with SQL injection",
			code:    "COUPON'; DROP TABLE coupons; --",
			isValid: false,
			reason:  "contains invalid characters",
		},
		{
			name:    "Coupon code with special characters",
			code:    "COUPON<>!@#",
			isValid: false,
			reason:  "contains invalid characters",
		},
		{
			name:    "Extremely long coupon code",
			code:    "SUMMER2023DISCOUNTVERYVERYLONGCODEFORATTACKSUMMER2023DISCOUNTVERYVERYLONGCODEFORATTACK",
			isValid: false,
			reason:  "exceeds maximum length",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid, reason := validateCouponCode(test.code)
			if isValid != test.isValid {
				t.Errorf("Expected validity %v but got %v", test.isValid, isValid)
			}
			if !isValid && test.reason != "" && reason == "" {
				t.Errorf("Expected a reason for invalid coupon code, but got none")
			}
		})
	}
}

// TestCouponValidityCheckBypass tests protection against coupon validity bypassing attempts
func TestCouponValidityCheckBypass(t *testing.T) {
	// Create expired coupon
	expiredCoupon := &domain.Coupon{
		ID:             "coupon1",
		Code:           "EXPIRED10",
		Discount:       10.0,
		DiscountType:   "percentage",
		MinOrderAmount: 20.0,
		ValidFrom:      time.Now().Add(-30 * 24 * time.Hour), // 30 days ago
		ValidTo:        time.Now().Add(-5 * 24 * time.Hour),  // 5 days ago
		IsActive:       true,
	}

	// Test various bypass attempts
	tests := []struct {
		name           string
		coupon         *domain.Coupon
		cartTotal      float64
		validationTime time.Time
		expectValid    bool
		reason         string
	}{
		{
			name:           "Expired coupon",
			coupon:         expiredCoupon,
			cartTotal:      50.0,
			validationTime: time.Now(),
			expectValid:    false,
			reason:         "coupon is expired",
		},
		{
			name: "Coupon below minimum order",
			coupon: &domain.Coupon{
				ID:             "coupon2",
				Code:           "MIN50",
				Discount:       15.0,
				DiscountType:   "percentage",
				MinOrderAmount: 50.0,
				ValidFrom:      time.Now().Add(-10 * 24 * time.Hour),
				ValidTo:        time.Now().Add(10 * 24 * time.Hour),
				IsActive:       true,
			},
			cartTotal:      30.0, // Below minimum
			validationTime: time.Now(),
			expectValid:    false,
			reason:         "below minimum order",
		},
		{
			name: "Inactive coupon",
			coupon: &domain.Coupon{
				ID:             "coupon3",
				Code:           "INACTIVE20",
				Discount:       20.0,
				DiscountType:   "percentage",
				MinOrderAmount: 20.0,
				ValidFrom:      time.Now().Add(-10 * 24 * time.Hour),
				ValidTo:        time.Now().Add(10 * 24 * time.Hour),
				IsActive:       false, // Inactive
			},
			cartTotal:      50.0,
			validationTime: time.Now(),
			expectValid:    false,
			reason:         "coupon is inactive",
		},
		{
			name: "Valid coupon",
			coupon: &domain.Coupon{
				ID:             "coupon4",
				Code:           "VALID25",
				Discount:       25.0,
				DiscountType:   "percentage",
				MinOrderAmount: 20.0,
				ValidFrom:      time.Now().Add(-10 * 24 * time.Hour),
				ValidTo:        time.Now().Add(10 * 24 * time.Hour),
				IsActive:       true,
			},
			cartTotal:      50.0,
			validationTime: time.Now(),
			expectValid:    true,
			reason:         "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid, reason := validateCouponForCart(test.coupon, test.cartTotal, test.validationTime)
			if isValid != test.expectValid {
				t.Errorf("Expected validity %v but got %v", test.expectValid, isValid)
			}
			if !isValid && test.reason != "" && !containsSubstring(reason, test.reason) {
				t.Errorf("Expected reason containing '%s' but got '%s'", test.reason, reason)
			}
		})
	}
}

// Helper functions for validation

// validateCartItem validates a cart item for security issues
func validateCartItem(item *domain.CartItem) (bool, string) {
	// Check for negative or zero price
	if item.Price <= 0 {
		return false, "negative price not allowed"
	}

	// Check for reasonable quantity
	if item.Quantity <= 0 || item.Quantity > 100 {
		return false, "quantity must be between 1 and 100"
	}

	// Check for XSS in name (simplified)
	if containsXSS(item.Name) {
		return false, "product name contains potentially harmful content"
	}

	// Validate image URL (simplified check)
	if !isValidImageURL(item.ImageURL) {
		return false, "invalid image URL"
	}

	return true, ""
}

// validateCouponCode validates a coupon code for security issues
func validateCouponCode(code string) (bool, string) {
	// Check for empty code
	if code == "" {
		return false, "coupon code cannot be empty"
	}

	// Check for reasonable length
	if len(code) > 20 {
		return false, "coupon code exceeds maximum length"
	}

	// Check for allowed characters (alphanumeric and some special chars)
	for _, c := range code {
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false, "coupon code contains invalid characters"
		}
	}

	return true, ""
}

// validateCouponForCart validates a coupon for a specific cart total at a specific time
func validateCouponForCart(coupon *domain.Coupon, cartTotal float64, validationTime time.Time) (bool, string) {
	// Check if coupon exists
	if coupon == nil {
		return false, "coupon not found"
	}

	// Check if coupon is active
	if !coupon.IsActive {
		return false, "coupon is inactive"
	}

	// Check if coupon is expired
	if validationTime.After(coupon.ValidTo) || validationTime.Before(coupon.ValidFrom) {
		return false, "coupon is expired or not yet valid"
	}

	// Check if cart total meets minimum order amount
	if cartTotal < coupon.MinOrderAmount {
		return false, "cart total is below coupon minimum order amount"
	}

	return true, ""
}

// Helper functions

// containsXSS checks if a string contains common XSS indicators (simplified)
func containsXSS(str string) bool {
	patterns := []string{"<script", "javascript:", "onerror=", "onload=", "onclick="}
	for _, pattern := range patterns {
		if containsSubstring(str, pattern) {
			return true
		}
	}
	return false
}

// isValidImageURL validates an image URL (simplified)
func isValidImageURL(url string) bool {
	// Reject URLs with javascript: protocol
	if containsSubstring(url, "javascript:") {
		return false
	}

	// Accept only http/https URLs
	if !containsSubstring(url, "http://") && !containsSubstring(url, "https://") {
		return false
	}

	return true
}

// containsSubstring checks if a string contains a substring (case insensitive)
func containsSubstring(s, substr string) bool {
	s, substr = toLowerCase(s), toLowerCase(substr)
	return s != "" && substr != "" && indexOf(s, substr) >= 0
}

// indexOf finds the index of substr in s
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// toLowerCase converts a string to lowercase (simple implementation)
func toLowerCase(s string) string {
	result := ""
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			result += string(c + 32)
		} else {
			result += string(c)
		}
	}
	return result
}

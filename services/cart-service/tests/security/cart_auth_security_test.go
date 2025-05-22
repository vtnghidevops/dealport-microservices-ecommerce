package security

import (
	"context"
	"testing"
	"time"
)

// TestCartOwnershipValidation tests prevention of unauthorized cart access
func TestCartOwnershipValidation(t *testing.T) {
	tests := []struct {
		name         string
		requestUser  string
		cartOwner    string
		isAuthorized bool
		description  string
	}{
		{
			name:         "Valid cart owner access",
			requestUser:  "user123",
			cartOwner:    "user123",
			isAuthorized: true,
			description:  "User accessing their own cart",
		},
		{
			name:         "Different user unauthorized access",
			requestUser:  "user456",
			cartOwner:    "user123",
			isAuthorized: false,
			description:  "User trying to access another user's cart",
		},
		{
			name:         "Empty user ID in request",
			requestUser:  "",
			cartOwner:    "user123",
			isAuthorized: false,
			description:  "Empty user ID in request",
		},
		{
			name:         "SQL injection attempt in user ID",
			requestUser:  "user123' OR '1'='1",
			cartOwner:    "user123",
			isAuthorized: false,
			description:  "SQL injection attempt in user ID",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Simulate authorization check
			isAuthorized := validateCartAccess(test.requestUser, test.cartOwner)
			if isAuthorized != test.isAuthorized {
				t.Errorf("Expected authorization %v but got %v for scenario: %s",
					test.isAuthorized, isAuthorized, test.description)
			}
		})
	}
}

// TestTokenBasedCartAccess tests security of token-based cart access
func TestTokenBasedCartAccess(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		token       *CartToken
		isValid     bool
		description string
	}{
		{
			name: "Valid token",
			token: &CartToken{
				UserID:    "user123",
				CartID:    "cart123",
				ExpiresAt: now.Add(15 * time.Minute),
				IssuedAt:  now.Add(-5 * time.Minute),
			},
			isValid:     true,
			description: "Valid token within expiration time",
		},
		{
			name: "Expired token",
			token: &CartToken{
				UserID:    "user123",
				CartID:    "cart123",
				ExpiresAt: now.Add(-5 * time.Minute),
				IssuedAt:  now.Add(-30 * time.Minute),
			},
			isValid:     false,
			description: "Token has expired",
		},
		{
			name: "Token from future",
			token: &CartToken{
				UserID:    "user123",
				CartID:    "cart123",
				ExpiresAt: now.Add(30 * time.Minute),
				IssuedAt:  now.Add(5 * time.Minute), // Issued in the future
			},
			isValid:     false,
			description: "Token issued in the future (clock manipulation attempt)",
		},
		{
			name: "Token with empty user ID",
			token: &CartToken{
				UserID:    "",
				CartID:    "cart123",
				ExpiresAt: now.Add(15 * time.Minute),
				IssuedAt:  now.Add(-5 * time.Minute),
			},
			isValid:     false,
			description: "Token missing user ID",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid := validateCartToken(test.token, now)
			if isValid != test.isValid {
				t.Errorf("Expected token validity %v but got %v for scenario: %s",
					test.isValid, isValid, test.description)
			}
		})
	}
}

// TestContextSecurityValidation tests security of user information in context
func TestContextSecurityValidation(t *testing.T) {
	tests := []struct {
		name         string
		setupContext func() context.Context
		expectUserID string
		isValid      bool
		description  string
	}{
		{
			name: "Valid user in context",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), userIDKey, "user123")
			},
			expectUserID: "user123",
			isValid:      true,
			description:  "Properly set user ID in context",
		},
		{
			name: "Missing user in context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectUserID: "",
			isValid:      false,
			description:  "No user ID in context",
		},
		{
			name: "Wrong type in context",
			setupContext: func() context.Context {
				// Using wrong type (int instead of string)
				return context.WithValue(context.Background(), userIDKey, 12345)
			},
			expectUserID: "",
			isValid:      false,
			description:  "Wrong data type for user ID in context",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := test.setupContext()
			userID, isValid := getUserFromContext(ctx)

			if isValid != test.isValid {
				t.Errorf("Expected validity %v but got %v for scenario: %s",
					test.isValid, isValid, test.description)
			}

			if userID != test.expectUserID {
				t.Errorf("Expected user ID %s but got %s for scenario: %s",
					test.expectUserID, userID, test.description)
			}
		})
	}
}

// Mock types and utility functions for testing

// CartToken represents a token for cart access
type CartToken struct {
	UserID    string
	CartID    string
	ExpiresAt time.Time
	IssuedAt  time.Time
}

// userIDKey is the key type for user ID in context
type userIDKeyType struct{}

var userIDKey = userIDKeyType{}

// validateCartAccess checks if the request user is authorized to access a cart
func validateCartAccess(requestUser, cartOwner string) bool {
	// Simple validation - requestUser must exactly match cartOwner
	// In real implementation, additional checks would be performed
	if requestUser == "" {
		return false
	}

	// Check for SQL injection attempts
	if containsSQL(requestUser) {
		return false
	}

	return requestUser == cartOwner
}

// validateCartToken validates a cart access token
func validateCartToken(token *CartToken, currentTime time.Time) bool {
	if token == nil {
		return false
	}

	// Validate token contents
	if token.UserID == "" || token.CartID == "" {
		return false
	}

	// Check token timing
	if token.ExpiresAt.Before(currentTime) {
		return false // Token expired
	}

	if token.IssuedAt.After(currentTime) {
		return false // Token from the future (possible clock manipulation)
	}

	return true
}

// getUserFromContext extracts and validates the user ID from context
func getUserFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	// Extract user ID from context
	userIDValue := ctx.Value(userIDKey)
	if userIDValue == nil {
		return "", false
	}

	// Type assertion
	userID, ok := userIDValue.(string)
	if !ok {
		return "", false
	}

	return userID, true
}

// containsSQL checks for SQL injection patterns
func containsSQL(s string) bool {
	sqlPatterns := []string{"'", ";", "--", "/*", "*/", "OR", "SELECT", "INSERT", "UPDATE", "DELETE", "DROP"}

	for _, pattern := range sqlPatterns {
		if contains(s, pattern) {
			return true
		}
	}

	return false
}

// Simple string contains helper
func contains(s, substr string) bool {
	return indexOf(s, substr) >= 0
}

// Simple string indexOf helper
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

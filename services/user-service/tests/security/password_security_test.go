package security

import (
	"strings"
	"testing"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// TestPasswordHashing tests password hashing functionality
func TestPasswordHashing(t *testing.T) {
	// Test case 1: Hash generation
	t.Run("Generate hash", func(t *testing.T) {
		password := "SecurePass123!"
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil {
			t.Errorf("Failed to generate password hash: %v", err)
		}

		if len(hash) == 0 {
			t.Error("Generated hash is empty")
		}
	})

	// Test case 2: Hash uniqueness
	t.Run("Hash uniqueness", func(t *testing.T) {
		password := "SecurePass123!"
		hash1, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		hash2, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if string(hash1) == string(hash2) {
			t.Error("Hash should be unique for each generation even with same password")
		}
	})

	// Test case 3: Hash verification
	t.Run("Verify hash", func(t *testing.T) {
		password := "SecurePass123!"
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		err := bcrypt.CompareHashAndPassword(hash, []byte(password))
		if err != nil {
			t.Errorf("Failed to verify password: %v", err)
		}

		err = bcrypt.CompareHashAndPassword(hash, []byte("WrongPassword"))
		if err == nil {
			t.Error("Verification should fail with incorrect password")
		}
	})
}

// TestPasswordRequirements tests password complexity requirements
func TestPasswordRequirements(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
		reason   string
	}{
		{"Valid password", "SecureP@ss123", true, ""},
		{"Too short", "Short1", false, "password too short"},
		{"No uppercase", "securepass123", false, "password requires at least one uppercase letter"},
		{"No lowercase", "SECUREPASS123", false, "password requires at least one lowercase letter"},
		{"No numbers", "SecurePassword", false, "password requires at least one number"},
		{"Common password", "MyPassword1", false, "password contains common password pattern"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			valid, reason := validatePasswordStrength(test.password)

			if valid != test.valid {
				t.Errorf("Expected validity %v but got %v", test.valid, valid)
			}

			if !valid && !strings.Contains(reason, test.reason) {
				t.Errorf("Expected reason to contain '%s', got '%s'", test.reason, reason)
			}
		})
	}
}

// validatePasswordStrength validates password complexity
func validatePasswordStrength(password string) (bool, string) {
	// Check length
	if len(password) < 8 {
		return false, "password too short, minimum 8 characters required"
	}

	// Check for uppercase, lowercase, and numbers
	var hasUpper, hasLower, hasNumber bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper {
		return false, "password requires at least one uppercase letter"
	}

	if !hasLower {
		return false, "password requires at least one lowercase letter"
	}

	if !hasNumber {
		return false, "password requires at least one number"
	}

	// Check for common passwords (simplified example)
	commonPasswords := []string{"password", "123456", "qwerty"}
	lowered := strings.ToLower(password)
	for _, common := range commonPasswords {
		if strings.Contains(lowered, common) {
			return false, "password contains common password pattern"
		}
	}

	return true, ""
}

// TestPasswordHistory tests prevention of password reuse
func TestPasswordHistory(t *testing.T) {
	// Simplified test for password history validation
	pastPasswords := []string{
		"$2a$10$abcdefghijklmnopqrstuvwxyz123456", // Mock bcrypt hashes
		"$2a$10$123456abcdefghijklmnopqrstuvwxyz",
	}

	t.Run("New password allowed", func(t *testing.T) {
		newPassword := "CompletelyNew123"
		newHash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

		// Check if password hash matches any in history
		isReused := false
		for _, oldHash := range pastPasswords {
			if err := bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(newPassword)); err == nil {
				isReused = true
				break
			}
		}

		if isReused {
			t.Error("New password should not match any in history")
		}

		// Verify new hash works with the password
		if err := bcrypt.CompareHashAndPassword(newHash, []byte(newPassword)); err != nil {
			t.Errorf("New password verification failed: %v", err)
		}
	})
}

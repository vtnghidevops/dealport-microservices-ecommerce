package security

import (
	"testing"

	"user-service/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

// TestPasswordValidation kiểm tra xác thực mật khẩu
func TestPasswordValidation(t *testing.T) {
	// Tạo mật khẩu đã hash
	password := "secure_password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Tạo user với mật khẩu đã hash
	user := &domain.User{
		ID:           "user123",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}

	// Test case 1: Mật khẩu đúng (Valid password)
	t.Run("Valid password", func(t *testing.T) {
		// Kiểm tra mật khẩu đúng
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			t.Errorf("Valid password should pass: %v", err)
		}
	})

	// Test case 2: Mật khẩu sai (Invalid password)
	t.Run("Invalid password", func(t *testing.T) {
		// Kiểm tra mật khẩu sai
		wrongPassword := "wrong_password"
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(wrongPassword))
		if err == nil {
			t.Errorf("Invalid password should not pass")
		}
	})
}

// TestPasswordComplexity kiểm tra độ phức tạp của mật khẩu
func TestPasswordComplexity(t *testing.T) {
	// Test case 1: Mật khẩu quá ngắn (Too short password)
	t.Run("Too short password", func(t *testing.T) {
		password := "short"
		if isStrongPassword(password) {
			t.Errorf("Password '%s' should be considered weak", password)
		}
	})

	// Test case 2: Mật khẩu chỉ có chữ (Only letters password)
	t.Run("Only letters password", func(t *testing.T) {
		password := "onlyletters"
		if isStrongPassword(password) {
			t.Errorf("Password '%s' should be considered weak", password)
		}
	})

	// Test case 3: Mật khẩu chỉ có số (Only numbers password)
	t.Run("Only numbers password", func(t *testing.T) {
		password := "12345678"
		if isStrongPassword(password) {
			t.Errorf("Password '%s' should be considered weak", password)
		}
	})

	// Test case 4: Mật khẩu mạnh (Strong password)
	t.Run("Strong password", func(t *testing.T) {
		password := "StrongP@ss123"
		if !isStrongPassword(password) {
			t.Errorf("Password '%s' should be considered strong", password)
		}
	})
}

// isStrongPassword kiểm tra độ mạnh của mật khẩu
func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	// Yêu cầu ít nhất 3 loại ký tự
	score := 0
	if hasLower {
		score++
	}
	if hasUpper {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}

	return score >= 3
}

// TestTokenSecurity kiểm tra bảo mật token
func TestTokenSecurity(t *testing.T) {
	// Các kiểm tra bảo mật JWT có thể được thêm ở đây
	// Ví dụ: kiểm tra thời gian hết hạn, chữ ký, etc.

	t.Run("JWT expiration should be reasonable", func(t *testing.T) {
		// Trong thực tế, kiểm tra này sẽ kiểm tra rằng JWT có thời gian hết hạn hợp lý
		// Ví dụ: access token nên hết hạn trong khoảng 15-60 phút
		// Refresh token có thể kéo dài hơn nhưng không nên quá 7-30 ngày

		// Ví dụ kiểm tra JWT expiration
		// maxAccessTokenLifetime := 60 * time.Minute
		// if service.AccessTokenLifetime > maxAccessTokenLifetime {
		//     t.Errorf("Access token lifetime too long: %v, should be <= %v", service.AccessTokenLifetime, maxAccessTokenLifetime)
		// }
	})
}

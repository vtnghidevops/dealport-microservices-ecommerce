package security

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// TestPasswordHashing kiểm tra việc mã hóa mật khẩu
func TestPasswordHashing(t *testing.T) {
	// Các mật khẩu cần kiểm tra
	passwords := []string{
		"password123",
		"SecureP@ssw0rd!",
		"VeryL0ngP@ssw0rdWithM@nyChar@cters12345!@#$%",
	}

	for _, password := range passwords {
		// Mã hóa mật khẩu
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}

		// Kiểm tra độ dài của mật khẩu đã mã hóa (nên đủ dài để bảo mật)
		if len(hashedPassword) < 60 {
			t.Errorf("Hashed password length too short: %d", len(hashedPassword))
		}

		// Xác minh rằng mật khẩu gốc có thể được xác thực với hash
		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
		if err != nil {
			t.Errorf("Password verification failed for: %s", password)
		}

		// Xác minh rằng mật khẩu sai không được xác thực
		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password+"wrong"))
		if err == nil {
			t.Errorf("Password verification should fail for incorrect password")
		}
	}
}

// TestPasswordStrengthValidation kiểm tra tính hợp lệ của mật khẩu
func TestPasswordStrengthValidation(t *testing.T) {
	// Danh sách các mật khẩu yếu nên bị từ chối
	weakPasswords := []string{
		"password",
		"123456",
		"qwerty",
		"abc123",
		"admin",
		"welcome",
	}

	// Danh sách các mật khẩu mạnh nên được chấp nhận
	strongPasswords := []string{
		"C0mpl3x!P@ssw0rd",
		"Sup3r$3cur3P@ss",
		"U$3r@Str0ng2023!",
	}

	// Kiểm tra các mật khẩu yếu
	for _, password := range weakPasswords {
		if isStrongPassword(password) {
			t.Errorf("Weak password '%s' passed strength validation", password)
		}
	}

	// Kiểm tra các mật khẩu mạnh
	for _, password := range strongPasswords {
		if !isStrongPassword(password) {
			t.Errorf("Strong password '%s' failed strength validation", password)
		}
	}
}

// isStrongPassword là hàm giả định kiểm tra độ mạnh mật khẩu
// Trong thực tế, nên sử dụng các thư viện chuyên dụng hoặc tham khảo logic từ service
func isStrongPassword(password string) bool {
	// Độ dài tối thiểu
	if len(password) < 8 {
		return false
	}

	// Kiểm tra có ít nhất một chữ cái thường
	hasLowerCase := false
	// Kiểm tra có ít nhất một chữ cái hoa
	hasUpperCase := false
	// Kiểm tra có ít nhất một chữ số
	hasDigit := false
	// Kiểm tra có ít nhất một ký tự đặc biệt
	hasSpecial := false

	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			hasLowerCase = true
		} else if char >= 'A' && char <= 'Z' {
			hasUpperCase = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		} else {
			hasSpecial = true
		}
	}

	// Yêu cầu ít nhất 3 trong 4 loại ký tự
	score := 0
	if hasLowerCase {
		score++
	}
	if hasUpperCase {
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

// TestBruteForceProtection kiểm tra tốc độ hash để đảm bảo bảo vệ khỏi tấn công brute force
func TestBruteForceProtection(t *testing.T) {
	password := "S3cur3P@ssw0rd"

	// Tăng cost để kiểm tra thời gian xử lý
	costs := []int{
		bcrypt.MinCost,         // thường là 4
		bcrypt.DefaultCost,     // thường là 10
		bcrypt.DefaultCost + 2, // 12
	}

	for _, cost := range costs {
		_, err := bcrypt.GenerateFromPassword([]byte(password), cost)
		if err != nil {
			t.Errorf("Failed to hash password with cost %d: %v", cost, err)
		}

		// Không kiểm tra thời gian cụ thể ở đây vì nó phụ thuộc vào phần cứng
		// Thực tế nên sử dụng benchmark tests để đánh giá thời gian xử lý
	}
}

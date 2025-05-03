package util

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

// OTPPurpose defines the purpose of OTP
type OTPPurpose string

// OTP purposes
const (
	OTPPurposeRegistration  OTPPurpose = "registration"
	OTPPurposePasswordReset OTPPurpose = "password_reset"
	OTPPurposeLogin         OTPPurpose = "login"
)

// OTPRecord represents a stored OTP
type OTPRecord struct {
	Email     string     // Email address
	OTP       string     // OTP code
	Purpose   OTPPurpose // Purpose of the OTP
	CreatedAt time.Time  // When the OTP was created
	ExpiresAt time.Time  // When the OTP expires
	Attempts  int        // Number of verification attempts
}

// OTPManager handles OTP generation, storage and verification
type OTPManager struct {
	otps        map[string]*OTPRecord // Map of email to OTP records
	mutex       sync.RWMutex          // Mutex for concurrent access
	otpLength   int                   // Length of OTP
	otpExpiry   time.Duration         // How long OTPs are valid for
	maxAttempts int                   // Maximum verification attempts
}

// NewOTPManager creates a new OTP manager
func NewOTPManager(otpLength int, otpExpiry time.Duration, maxAttempts int) *OTPManager {
	return &OTPManager{
		otps:        make(map[string]*OTPRecord),
		otpLength:   otpLength,
		otpExpiry:   otpExpiry,
		maxAttempts: maxAttempts,
	}
}

// GenerateOTP generates a new OTP for the given email and purpose
func (m *OTPManager) GenerateOTP(email string, purpose OTPPurpose) (string, error) {
	// Generate random OTP
	otp, err := m.generateRandomOTP()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Store the OTP
	m.otps[email] = &OTPRecord{
		Email:     email,
		OTP:       otp,
		Purpose:   purpose,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(m.otpExpiry),
		Attempts:  0,
	}

	return otp, nil
}

// VerifyOTP verifies the provided OTP for the email and purpose
func (m *OTPManager) VerifyOTP(email, otp string, purpose OTPPurpose) (bool, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Get the OTP record
	record, exists := m.otps[email]
	if !exists {
		return false, fmt.Errorf("no OTP found for email: %s", email)
	}

	// Check purpose
	if record.Purpose != purpose {
		return false, fmt.Errorf("OTP purpose mismatch")
	}

	// Check expiration
	if time.Now().After(record.ExpiresAt) {
		// Delete expired OTP
		delete(m.otps, email)
		return false, fmt.Errorf("OTP has expired")
	}

	// Increment attempts
	record.Attempts++

	// Check max attempts
	if record.Attempts > m.maxAttempts {
		delete(m.otps, email)
		return false, fmt.Errorf("maximum verification attempts exceeded")
	}

	// Verify OTP
	if record.OTP != otp {
		return false, fmt.Errorf("invalid OTP")
	}

	// OTP verified successfully, delete it to prevent reuse
	delete(m.otps, email)
	return true, nil
}

// GetRemainingTime gets the remaining time for an OTP in minutes
func (m *OTPManager) GetRemainingTime(email string) (int, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	record, exists := m.otps[email]
	if !exists {
		return 0, fmt.Errorf("no OTP found for email: %s", email)
	}

	remainingTime := time.Until(record.ExpiresAt)
	if remainingTime <= 0 {
		return 0, nil
	}

	return int(remainingTime.Minutes()) + 1, nil
}

// generateRandomOTP generates a random numeric OTP of specified length
func (m *OTPManager) generateRandomOTP() (string, error) {
	// For a 6-digit OTP, we need 6 bytes
	bytes := make([]byte, m.otpLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to digits
	for i := 0; i < m.otpLength; i++ {
		bytes[i] = '0' + bytes[i]%10
	}

	return string(bytes), nil
}

// CleanupExpiredOTPs removes expired OTPs (can be called periodically)
func (m *OTPManager) CleanupExpiredOTPs() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	for email, record := range m.otps {
		if now.After(record.ExpiresAt) {
			delete(m.otps, email)
		}
	}
}

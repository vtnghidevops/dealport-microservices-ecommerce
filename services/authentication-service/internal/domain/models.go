package domain

import (
	"errors"
	"time"
)

// Common errors
var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user with this email already exists")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
)

// User represents a user account
type User struct {
	ID           string     `json:"id" bson:"_id,omitempty" db:"id"`
	Email        string     `json:"email" bson:"email" db:"email"`
	Password     string     `json:"-" bson:"password_hash" db:"password_hash"`
	FirstName    string     `json:"first_name" bson:"first_name" db:"first_name"`
	LastName     string     `json:"last_name" bson:"last_name" db:"last_name"`
	Username     string     `json:"username" bson:"username" db:"username"`
	Phone        string     `json:"phone,omitempty" bson:"phone" db:"phone"`
	Role         string     `json:"role" bson:"role" db:"role"`
	Status       string     `json:"status" bson:"status" db:"status"`
	Active       bool       `json:"active" bson:"active" db:"active"`
	RefreshToken *string    `json:"-" bson:"refresh_token,omitempty" db:"refresh_token"`
	CreatedAt    time.Time  `json:"created_at" bson:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" bson:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"-" bson:"deleted_at,omitempty" db:"deleted_at"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Phone     string `json:"phone,omitempty"`
}

// VerifyRegistrationRequest represents a registration verification request
type VerifyRegistrationRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse represents the response after user registration
type RegisterResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// TokenDetails contains access and refresh tokens
type TokenDetails struct {
	// AccessToken  string `json:"access_token"`
	// RefreshToken string `json:"refresh_token"`
	// AccessUUID   string `json:"access_uuid"`
	// RefreshUUID  string `json:"refresh_uuid"`
	// AtExpires    int64  `json:"at_expires"`
	// RtExpires    int64  `json:"rt_expires"`
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

// TokenMetadata contains token metadata
type TokenMetadata struct {
	// UserID string `json:"user_id"`
	// Email  string `json:"email"`
	// Role   string `json:"role"`
	// UUID   string `json:"uuid"`
	// Exp    int64  `json:"exp"`
	UserID string
	Email  string
	Role   string
	UUID   string
	Exp    int64
}

// PasswordResetRequest represents a password reset request
type PasswordResetRequest struct {
	Email string `json:"email"`
}

// VerifyPasswordResetRequest represents a password reset verification request
type VerifyPasswordResetRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

// UpdatePasswordRequest represents a password update request
type UpdatePasswordRequest struct {
	Token           string `json:"token"`                      // Reset token or access token
	Password        string `json:"password"`                   // New password
	CurrentPassword string `json:"current_password,omitempty"` // Current password (when logged in)
	Email           string `json:"email,omitempty"`            // Email for verification with token
}

// VerifyOTPRequest represents a generic OTP verification request
type VerifyOTPRequest struct {
	Email   string `json:"email"`
	OTP     string `json:"otp"`
	Purpose string `json:"purpose"` // registration, password_reset, etc.
}

// OTP represents an OTP record in the database
type OTP struct {
	ID        string    `json:"id" bson:"_id,omitempty" db:"id"`
	Email     string    `json:"email" bson:"email" db:"email"`
	Code      string    `json:"code" bson:"code" db:"code"`
	Purpose   string    `json:"purpose" bson:"purpose" db:"purpose"`
	Used      bool      `json:"used" bson:"used" db:"used"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" db:"created_at"`
}

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

// User represents a user in the system
type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"-" db:"password_hash"` // Password hash never exported in JSON
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Username     string    `json:"username" db:"username"`
	Active       bool      `json:"active" db:"active"`
	Status       string    `json:"status" db:"status"`
	Role         string    `json:"role" db:"role"`
	RefreshToken *string   `json:"-" db:"refresh_token"` // Using pointer to handle NULL from DB
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// RegisterResponse represents the response after user registration
type RegisterResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// TokenDetails contains JWT access and refresh token details
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

// TokenMetadata contains token claim details
type TokenMetadata struct {
	UserID string
	Email  string
	Role   string
	UUID   string
	Exp    int64
}

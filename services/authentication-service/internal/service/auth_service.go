package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"authentication-service/internal/domain"
	"authentication-service/internal/repository"
)

// jwtCustomClaims contains custom claims data
type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	UUID   string `json:"uuid"`
	jwt.RegisteredClaims
}

// authService implements the AuthService interface
type authService struct {
	userRepo        repository.UserRepository
	accessSecret    string
	refreshSecret   string
	accessDuration  time.Duration
	refreshDuration time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo repository.UserRepository,
	accessSecret, refreshSecret string,
	accessDuration, refreshDuration time.Duration,
) AuthService {
	return &authService{
		userRepo:        userRepo,
		accessSecret:    accessSecret,
		refreshSecret:   refreshSecret,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

// Register creates a new user account
func (s *authService) Register(ctx context.Context, req *domain.RegisterRequest) (string, error) {
	// Check if user with the same email exists
	exists, err := s.userRepo.UserExists(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("failed to check existing user: %w", err)
	}

	if exists {
		return "", domain.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Debug log passwords
	log.Printf("Original password: %s", req.Password)
	log.Printf("Hashed password: %s", string(hashedPassword))

	// Create new user
	user := &domain.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Role:      "user",
		Status:    "active",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user to database
	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return user.ID, nil
}

// Login authenticates a user and returns JWT tokens
func (s *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenDetails, *domain.User, error) {
	fmt.Printf("DEBUG Login: Attempting login for email: %s\n", req.Email)

	// Retrieve user by email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		fmt.Printf("DEBUG Login: User not found: %v\n", err)
		return nil, nil, errors.New("invalid credentials")
	}

	fmt.Printf("DEBUG Login: User found with ID: %s\n", user.ID)

	// Compare passwords - using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		fmt.Printf("DEBUG Login: Password comparison failed: %v\n", err)
		fmt.Printf("DEBUG Login: Stored hash: %s\n", user.Password)
		fmt.Printf("DEBUG Login: Provided password: %s\n", req.Password)
		return nil, nil, errors.New("invalid credentials")
	}

	fmt.Printf("DEBUG Login: Password match successful\n")

	// Create token
	td, err := s.createToken(user.ID, user.Email, user.Role)
	if err != nil {
		fmt.Printf("DEBUG Login: Failed to create token: %v\n", err)
		return nil, nil, fmt.Errorf("failed to create token: %w", err)
	}

	// Print full tokens for debugging
	// fmt.Printf("DEBUG Login: FULL ACCESS TOKEN: %s\n", td.AccessToken)
	// fmt.Printf("DEBUG Login: FULL REFRESH TOKEN: %s\n", td.RefreshToken)

	// Store refresh token in the database
	fmt.Printf("DEBUG Login: Storing refresh token in database\n")
	err = s.userRepo.UpdateRefreshToken(ctx, user.ID, td.RefreshToken)
	if err != nil {
		fmt.Printf("DEBUG Login: Failed to store refresh token: %v\n", err)
		return nil, nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Don't return password
	user.Password = ""

	fmt.Printf("DEBUG Login: Login successful, tokens created\n")
	return td, user, nil
}

// ValidateToken validates an access token and returns claims
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*domain.TokenMetadata, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.accessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok {
		return nil, errors.New("failed to parse token claims")
	}

	// Get expiration time
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("failed to get expiration time: %w", err)
	}

	// Create token metadata
	metadata := &domain.TokenMetadata{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
		UUID:   claims.UUID,
		Exp:    exp.Unix(),
	}

	return metadata, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenDetails, error) {
	// Debug logs - print full token
	fmt.Printf("DEBUG RefreshToken: FULL TOKEN: %s\n", refreshToken)
	fmt.Printf("DEBUG RefreshToken: Token length: %d\n", len(refreshToken))
	fmt.Printf("DEBUG RefreshToken: Using refresh secret: %s\n", s.refreshSecret)

	// Parse token
	token, err := jwt.ParseWithClaims(refreshToken, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("DEBUG RefreshToken: Invalid signing method: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.refreshSecret), nil
	})

	if err != nil {
		fmt.Printf("DEBUG RefreshToken: Token parsing error: %v\n", err)
		// Try to decode the token without verification to see if structure is valid
		parts := strings.Split(refreshToken, ".")
		if len(parts) == 3 {
			fmt.Printf("DEBUG RefreshToken: Token has valid format with 3 parts\n")
			fmt.Printf("DEBUG RefreshToken: Header: %s\n", parts[0])
			fmt.Printf("DEBUG RefreshToken: Payload: %s\n", parts[1])
			fmt.Printf("DEBUG RefreshToken: Signature: %s\n", parts[2])

			// Try to decode the payload (middle part)
			if payload, err := base64.RawURLEncoding.DecodeString(parts[1]); err == nil {
				fmt.Printf("DEBUG RefreshToken: Token payload decoded: %s\n", string(payload))
			} else {
				fmt.Printf("DEBUG RefreshToken: Cannot decode payload: %v\n", err)
			}
		} else {
			fmt.Printf("DEBUG RefreshToken: Token format invalid, has %d parts instead of 3\n", len(parts))
		}
		return nil, err
	}

	// Check if token is valid
	if !token.Valid {
		fmt.Printf("DEBUG RefreshToken: Token is invalid\n")
		return nil, errors.New("invalid refresh token")
	}

	// Extract claims
	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok {
		fmt.Printf("DEBUG RefreshToken: Failed to parse token claims\n")
		return nil, errors.New("failed to parse token claims")
	}

	fmt.Printf("DEBUG RefreshToken: Token claims - UserID: %s, Email: %s, Role: %s\n",
		claims.UserID, claims.Email, claims.Role)

	// Get expiration time
	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("failed to get expiration time: %w", err)
	}
	fmt.Printf("DEBUG RefreshToken: Token expires at: %v\n", expirationTime.Time)

	// Verify token in the database
	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		fmt.Printf("DEBUG RefreshToken: User not found: %v\n", err)
		return nil, errors.New("user not found")
	}

	// Check if refresh token matches (handle nil pointer case)
	storedToken := ""
	if user.RefreshToken != nil {
		storedToken = *user.RefreshToken
	}

	fmt.Printf("DEBUG RefreshToken: Stored token in DB: %s\n",
		func() string {
			if storedToken == "" {
				return "NULL or EMPTY"
			}
			return storedToken[:min(20, len(storedToken))] + "..."
		}())

	if storedToken != refreshToken {
		fmt.Printf("DEBUG RefreshToken: Token mismatch! Stored != Provided\n")
		fmt.Printf("DEBUG RefreshToken: First 50 chars of stored token: %s\n",
			storedToken[:min(50, len(storedToken))])
		fmt.Printf("DEBUG RefreshToken: First 50 chars of provided token: %s\n",
			refreshToken[:min(50, len(refreshToken))])
		return nil, errors.New("refresh token has been revoked")
	}

	// Generate new tokens
	td, err := s.createToken(user.ID, user.Email, user.Role)
	if err != nil {
		fmt.Printf("DEBUG RefreshToken: Failed to generate token: %v\n", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update refresh token in the database
	fmt.Printf("DEBUG RefreshToken: Updating refresh token in database for user %s\n", user.ID)
	if err := s.userRepo.UpdateRefreshToken(ctx, user.ID, td.RefreshToken); err != nil {
		fmt.Printf("DEBUG RefreshToken: Failed to update refresh token in DB: %v\n", err)
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	fmt.Printf("DEBUG RefreshToken: Successfully refreshed token for user %s\n", user.ID)
	return td, nil
}

// GetUserByID retrieves a user by ID
func (s *authService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

// createToken generates access and refresh tokens
func (s *authService) createToken(userID, email, role string) (*domain.TokenDetails, error) {
	fmt.Printf("DEBUG createToken: Creating token for user ID: %s, email: %s, role: %s\n", userID, email, role)
	//fmt.Printf("DEBUG createToken: Using accessSecret: %s\n", s.accessSecret)
	//fmt.Printf("DEBUG createToken: Using refreshSecret: %s\n", s.refreshSecret)

	td := &domain.TokenDetails{
		AccessUUID:  uuid.New().String(),
		RefreshUUID: uuid.New().String(),
		AtExpires:   time.Now().Add(s.accessDuration).Unix(),
		RtExpires:   time.Now().Add(s.refreshDuration).Unix(),
	}

	fmt.Printf("DEBUG createToken: Token durations - Access: %v, Refresh: %v\n",
		s.accessDuration, s.refreshDuration)
	fmt.Printf("DEBUG createToken: Token expiry - Access: %v, Refresh: %v\n",
		time.Unix(td.AtExpires, 0), time.Unix(td.RtExpires, 0))

	// Create access token
	atClaims := jwtCustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		UUID:   td.AccessUUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.AtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(s.accessSecret))
	if err != nil {
		fmt.Printf("DEBUG createToken: Error signing access token: %v\n", err)
		return nil, err
	}
	td.AccessToken = accessToken

	// Create refresh token
	rtClaims := jwtCustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		UUID:   td.RefreshUUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.RtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(s.refreshSecret))
	if err != nil {
		fmt.Printf("DEBUG createToken: Error signing refresh token: %v\n", err)
		return nil, err
	}
	td.RefreshToken = refreshToken

	fmt.Printf("DEBUG createToken: Access token generated (first 20 chars): %s...\n",
		accessToken[:min(20, len(accessToken))])
	fmt.Printf("DEBUG createToken: Refresh token generated (first 20 chars): %s...\n",
		refreshToken[:min(20, len(refreshToken))])

	return td, nil
}

// Helper function for debug
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

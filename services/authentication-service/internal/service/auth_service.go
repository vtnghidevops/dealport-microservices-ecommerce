package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"authentication-service/internal/domain"
	"authentication-service/internal/event"
	"authentication-service/internal/repository"
	"authentication-service/internal/util"
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
	otpManager      *util.OTPManager
	mailClient      *util.MailClient
	eventEmitter    *event.Emitter
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo repository.UserRepository,
	accessSecret, refreshSecret string,
	accessDuration, refreshDuration time.Duration,
	otpManager *util.OTPManager,
	mailClient *util.MailClient,
	eventEmitter *event.Emitter,
) AuthService {
	return &authService{
		userRepo:        userRepo,
		accessSecret:    accessSecret,
		refreshSecret:   refreshSecret,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
		otpManager:      otpManager,
		mailClient:      mailClient,
		eventEmitter:    eventEmitter,
	}
}

// Register creates a new user account and sends OTP
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

	// Create new user with pending status
	user := &domain.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Phone:     req.Phone, // Lưu số điện thoại từ request
		Role:      "user",
		Status:    "pending", // Set as pending until OTP verification
		Active:    false,     // Not active until verified
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user to database
	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	// Get the user from database to verify status and active fields were saved correctly
	savedUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("Warning: Could not verify user creation: %v", err)
	} else {
		log.Printf("DEBUG Register: Saved user status=%s, active=%v", savedUser.Status, savedUser.Active)
		if savedUser.Status != "pending" || savedUser.Active != false {
			log.Printf("WARNING: User was not saved with correct status/active values. Updating...")
			savedUser.Status = "pending"
			savedUser.Active = false
			err = s.userRepo.UpdateUser(ctx, savedUser)
			if err != nil {
				log.Printf("Error fixing user status: %v", err)
			}
		}
	}

	// Generate OTP for email verification
	otp, err := s.otpManager.GenerateOTP(req.Email, util.OTPPurposeRegistration)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Get OTP expiration time in minutes
	expiresIn, _ := s.otpManager.GetRemainingTime(req.Email)

	// Send OTP via event emitter
	go func() {
		err := s.eventEmitter.EmitOTPGenerated(
			req.Email,
			otp,
			string(util.OTPPurposeRegistration),
			expiresIn,
			"We received a request to create an account with this email address.",
			"account registration",
		)
		if err != nil {
			log.Printf("Failed to emit OTP generated event: %v", err)

			// Fallback to direct mail client if event emission fails
			err = s.mailClient.SendRegistrationOTP(req.Email, otp, expiresIn)
			if err != nil {
				log.Printf("Failed to send registration OTP email via fallback: %v", err)
			}
		}
	}()

	return user.ID, nil
}

// VerifyRegistration verifies registration OTP and activates the user
func (s *authService) VerifyRegistration(ctx context.Context, req *domain.VerifyRegistrationRequest) (*domain.TokenDetails, *domain.User, error) {
	// Verify OTP
	valid, err := s.otpManager.VerifyOTP(req.Email, req.OTP, util.OTPPurposeRegistration)
	if err != nil {
		return nil, nil, fmt.Errorf("OTP verification failed: %w", err)
	}

	if !valid {
		return nil, nil, fmt.Errorf("invalid OTP")
	}

	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Lưu trạng thái trước khi cập nhật (chỉ để ghi log)
	wasActive := user.Active
	previousStatus := user.Status
	log.Printf("DEBUG VerifyRegistration: User %s current status=%s, active=%v", user.ID, previousStatus, wasActive)

	// Update user status to active
	user.Status = "active"
	user.Active = true
	user.UpdatedAt = time.Now()

	err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update user: %w", err)
	}

	// LUÔN phát sự kiện user.registered sau khi xác thực OTP thành công,
	// bất kể trạng thái active trước đó là gì
	if s.eventEmitter != nil {
		log.Printf("DEBUG VerifyRegistration: Emitting user.registered event for user %s", user.ID)
		// Create a display name from first name and last name
		displayName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)

		// Log thông tin trước khi tạo dữ liệu sự kiện
		log.Printf("DEBUG VerifyRegistration: User data - Phone: '%s', length: %d", user.Phone, len(user.Phone))

		userData := event.UserRegisteredData{
			ID:          user.ID,
			Email:       user.Email,
			Username:    user.Username,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			DisplayName: displayName, // Use the generated display name
			Phone:       user.Phone,  // Đảm bảo trường phone được thêm vào
			Role:        user.Role,
			Status:      user.Status,
			Active:      user.Active,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		}

		log.Printf("DEBUG VerifyRegistration: User data being sent in event: ID=%s, Email=%s, Phone=%s",
			userData.ID, userData.Email, userData.Phone)

		err := s.eventEmitter.EmitUserRegistered(userData)
		if err != nil {
			log.Printf("CRITICAL ERROR: Failed to emit user.registered event: %v", err)
		} else {
			log.Printf("SUCCESS: Event user.registered published with routing key user.registered")
		}
	} else {
		log.Printf("ERROR VerifyRegistration: Event emitter is nil, cannot emit user.registered event")
	}

	// Double-check after update to ensure status is correct
	updatedUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		log.Printf("DEBUG VerifyRegistration: After update, user status=%s, active=%v", updatedUser.Status, updatedUser.Active)
	}

	// Create tokens
	td, err := s.createToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create token: %w", err)
	}

	// Store refresh token
	err = s.userRepo.UpdateRefreshToken(ctx, user.ID, td.RefreshToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Don't return password
	user.Password = ""

	return td, user, nil
}

// RequestPasswordReset initiates password reset process
func (s *authService) RequestPasswordReset(ctx context.Context, email string) error {
	// Check if user exists
	exists, err := s.userRepo.UserExists(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if !exists {
		// Don't reveal that the user doesn't exist for security reasons
		return nil
	}

	// Generate OTP for password reset
	otp, err := s.otpManager.GenerateOTP(email, util.OTPPurposePasswordReset)
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Get OTP expiration time in minutes
	expiresIn, _ := s.otpManager.GetRemainingTime(email)

	// Create a token hash (essentially a placeholder since we're using OTP)
	tokenHash := fmt.Sprintf("%s:%s:%d", email, uuid.New().String(), time.Now().Unix())

	// Send OTP via event emitter
	go func() {
		// Emit password reset requested event
		err := s.eventEmitter.EmitPasswordResetRequested(email, tokenHash, time.Now().Add(time.Duration(expiresIn)*time.Minute))
		if err != nil {
			log.Printf("Failed to emit password reset requested event: %v", err)
		}

		// Also emit OTP generated event
		err = s.eventEmitter.EmitOTPGenerated(
			email,
			otp,
			string(util.OTPPurposePasswordReset),
			expiresIn,
			"We received a request to reset the password for your account.",
			"password reset",
		)
		if err != nil {
			log.Printf("Failed to emit OTP generated event: %v", err)

			// Fallback to direct mail client if event emission fails
			err = s.mailClient.SendPasswordResetOTP(email, otp, expiresIn)
			if err != nil {
				log.Printf("Failed to send password reset OTP email via fallback: %v", err)
			}
		}
	}()

	return nil
}

// VerifyPasswordReset verifies password reset OTP
func (s *authService) VerifyPasswordReset(ctx context.Context, req *domain.VerifyPasswordResetRequest) (string, error) {
	// Verify OTP
	valid, err := s.otpManager.VerifyOTP(req.Email, req.OTP, util.OTPPurposePasswordReset)
	if err != nil {
		return "", fmt.Errorf("OTP verification failed: %w", err)
	}

	if !valid {
		return "", fmt.Errorf("invalid OTP")
	}

	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Create a special JWT token for password reset
	// This token will have a short expiry time
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"type":    "password_reset",
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	})

	// Sign the token
	tokenString, err := token.SignedString([]byte(s.accessSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign reset token: %w", err)
	}

	return tokenString, nil
}

// UpdatePassword updates user's password
func (s *authService) UpdatePassword(ctx context.Context, req *domain.UpdatePasswordRequest) error {
	var userID string
	var isResetToken bool
	var userEmail string

	// Debug log for token
	fmt.Printf("DEBUG UpdatePassword: Received token: %s\n", req.Token)

	// Determine if this is a password reset or password change
	if req.CurrentPassword != "" {
		// This is a password change, validate token (must be access token)
		metadata, err := s.ValidateToken(ctx, req.Token)
		if err != nil {
			return fmt.Errorf("invalid token: %w", err)
		}
		userID = metadata.UserID
		isResetToken = false
	} else {
		// This is a password reset using OTP verification token
		// First try to parse as JWT (for backward compatibility)
		jwtToken, jwtErr := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.accessSecret), nil
		})

		if jwtErr == nil && jwtToken.Valid {
			// It's a valid JWT token
			claims, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				return fmt.Errorf("invalid token claims")
			}

			// Check if this is a password reset token
			tokenType, ok := claims["type"].(string)
			if !ok || tokenType != "password_reset" {
				return fmt.Errorf("invalid token type")
			}

			userID, ok = claims["user_id"].(string)
			if !ok {
				return fmt.Errorf("invalid user ID in token")
			}
			isResetToken = true
		} else {
			// Not a JWT token, try base64 decode (format: email:uuid:timestamp)
			fmt.Printf("DEBUG UpdatePassword: Not a JWT token, trying base64 decode\n")
			fmt.Printf("DEBUG UpdatePassword: Token length: %d\n", len(req.Token))

			// Log first few characters of the token for debugging
			if len(req.Token) > 0 {
				fmt.Printf("DEBUG UpdatePassword: Token starts with: %s\n", req.Token[0:min(20, len(req.Token))])
			} else {
				fmt.Printf("DEBUG UpdatePassword: Empty token received\n")
			}

			// Decode base64 string
			decodedBytes, err := base64.StdEncoding.DecodeString(req.Token)
			if err != nil {
				fmt.Printf("DEBUG UpdatePassword: Failed to decode base64: %v\n", err)
				return fmt.Errorf("invalid reset token: failed to decode base64: %w", err)
			}

			// Log decoded string for debugging
			decodedStr := string(decodedBytes)
			fmt.Printf("DEBUG UpdatePassword: Decoded token: %s\n", decodedStr[:min(30, len(decodedStr))])

			// Parse components
			parts := strings.Split(decodedStr, ":")
			fmt.Printf("DEBUG UpdatePassword: Token has %d parts\n", len(parts))
			if len(parts) != 3 {
				return fmt.Errorf("invalid reset token: token is malformed: expected 3 parts, got %d", len(parts))
			}

			// Extract parts
			tokenEmail := parts[0]
			// tokenUUID := parts[1]  // Not using UUID for now, but could be used for additional verification
			tokenTimestamp, err := strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid reset token: invalid timestamp: %w", err)
			}

			// Validate timestamp (token valid for 1 hour)
			if time.Now().Unix()-tokenTimestamp > 3600 {
				return fmt.Errorf("reset token has expired")
			}

			// Validate email with the one provided in the request
			fmt.Printf("DEBUG UpdatePassword: Token email: %s, Request email: %s\n", tokenEmail, req.Email)
			if strings.ToLower(tokenEmail) != strings.ToLower(req.Email) {
				return fmt.Errorf("email mismatch: token email does not match request email")
			}

			// Get user by email
			user, err := s.userRepo.GetUserByEmail(ctx, tokenEmail)
			if err != nil {
				return fmt.Errorf("failed to get user: %w", err)
			}

			userID = user.ID
			userEmail = tokenEmail
			isResetToken = true
		}
	}

	// Get user
	var user *domain.User
	var err error

	if userEmail != "" {
		// If we already have the user from email lookup above
		user, err = s.userRepo.GetUserByEmail(ctx, userEmail)
	} else {
		// Otherwise get by ID
		user, err = s.userRepo.GetUserByID(ctx, userID)
	}

	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// If this is a password change (not reset), verify current password
	if !isResetToken {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
		if err != nil {
			return fmt.Errorf("invalid current password")
		}
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user's password
	fmt.Printf("DEBUG UpdatePassword: Updating password for user ID: %s\n", user.ID)
	fmt.Printf("DEBUG UpdatePassword: Old password hash: %s\n", user.Password[:min(20, len(user.Password))]+"...")
	fmt.Printf("DEBUG UpdatePassword: New password hash: %s\n", string(hashedPassword[:min(20, len(hashedPassword))])+"...")

	// Store the old password hash to verify if it was updated
	oldPasswordHash := user.Password

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		fmt.Printf("DEBUG UpdatePassword: Failed to update user: %v\n", err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Get the user again to check if the password was updated
	updatedUser, err := s.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		fmt.Printf("DEBUG UpdatePassword: Failed to get updated user: %v\n", err)
		return fmt.Errorf("failed to verify password update: %w", err)
	}

	// Compare old and new password hashes
	if oldPasswordHash == updatedUser.Password {
		fmt.Printf("DEBUG UpdatePassword: Password not updated!\n")
		return fmt.Errorf("password was not updated in the database")
	} else {
		fmt.Printf("DEBUG UpdatePassword: Password updated successfully\n")
	}

	// Emit password changed event
	go func() {
		err := s.eventEmitter.EmitPasswordChanged(user.Email)
		if err != nil {
			log.Printf("Failed to emit password changed event: %v", err)
		}
	}()

	return nil
}

// Login authenticates a user and returns tokens
func (s *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenDetails, *domain.User, error) {
	// Check if email, password are provided
	if req.Email == "" || req.Password == "" {
		return nil, nil, fmt.Errorf("email and password are required")
	}

	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Không cần hiện chi tiết lỗi, chỉ log chung là đăng nhập thất bại
		s.emitLoginFailedEvent(req.Email, "invalid_credentials", nil)
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// Đăng nhập thất bại, emit event
		s.emitLoginFailedEvent(req.Email, "invalid_credentials", nil)
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Create tokens
	td, err := s.createToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create authentication token: %w", err)
	}

	// Store refresh token in the database
	if err := s.userRepo.UpdateRefreshToken(ctx, user.ID, td.RefreshToken); err != nil {
		return nil, nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Đăng nhập thành công, emit event
	metadata := map[string]interface{}{
		"role": user.Role,
	}
	if s.eventEmitter != nil {
		s.eventEmitter.EmitLoginSuccess(user.ID, user.Email, metadata)
	}

	// Mask password before returning user to client
	user.Password = ""
	return td, user, nil
}

// Helper function to emit login failed event
func (s *authService) emitLoginFailedEvent(email, reason string, metadata map[string]interface{}) {
	if s.eventEmitter != nil {
		s.eventEmitter.EmitLoginFailed(email, reason, metadata)
	}
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

// RequestOTP generates and sends a new OTP for a specific purpose
func (s *authService) RequestOTP(ctx context.Context, email string, purpose string) error {
	// Check if user exists
	exists, err := s.userRepo.UserExists(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	if !exists {
		// Don't reveal that the user doesn't exist for security reasons
		return nil
	}

	// Validate and convert the purpose string to OTPPurpose type
	var otpPurpose util.OTPPurpose
	switch purpose {
	case string(util.OTPPurposeRegistration):
		otpPurpose = util.OTPPurposeRegistration
	case string(util.OTPPurposePasswordReset):
		otpPurpose = util.OTPPurposePasswordReset
	case string(util.OTPPurposeLogin):
		otpPurpose = util.OTPPurposeLogin
	default:
		return fmt.Errorf("invalid OTP purpose: %s", purpose)
	}

	// Generate OTP for the requested purpose
	otp, err := s.otpManager.GenerateOTP(email, otpPurpose)
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Get OTP expiration time in minutes
	expiresIn, _ := s.otpManager.GetRemainingTime(email)

	// IMPORTANT: The code below uses direct HTTP calls to the mail service,
	// which conflicts with the RabbitMQ message queue approach.
	// We should use only one approach for email sending to prevent duplicates.
	// Currently, we've standardized on using RabbitMQ (eventEmitter) for all email communications.
	//
	// The direct HTTP call to mail service is commented out to avoid duplicate emails.
	// Instead, we should use the eventEmitter to publish events to RabbitMQ,
	// which will be picked up by the listener service and forwarded to the mail service.
	/*
		// Send OTP email based on purpose
		go func() {
			var mailErr error
			switch otpPurpose {
			case util.OTPPurposeRegistration:
				mailErr = s.mailClient.SendRegistrationOTP(email, otp, expiresIn)
			case util.OTPPurposePasswordReset:
				mailErr = s.mailClient.SendPasswordResetOTP(email, otp, expiresIn)
			default:
				// Generic OTP email
				mailErr = s.mailClient.SendGenericOTP(email, otp, string(otpPurpose), expiresIn)
			}

			if mailErr != nil {
				log.Printf("Failed to send OTP email for purpose %s: %v", purpose, mailErr)
			}
		}()
	*/

	// Use event emitter to publish the OTP generated event
	// This will be picked up by the listener service and forwarded to the mail service
	message := ""
	actionText := "Verify"

	switch otpPurpose {
	case util.OTPPurposeRegistration:
		message = "We received a request to create an account with this email address."
		actionText = "Complete Registration"
	case util.OTPPurposePasswordReset:
		message = "We received a request to reset the password for your account."
		actionText = "Reset Password"
	default:
		message = fmt.Sprintf("We received a request that requires verification for your account (%s).", purpose)
	}

	err = s.eventEmitter.EmitOTPGenerated(email, otp, string(otpPurpose), expiresIn, message, actionText)
	if err != nil {
		log.Printf("Failed to emit OTP generated event: %v", err)
		// Don't return the error to avoid revealing information about the internal system
	}

	return nil
}

// VerifyOTP verifies an OTP code for a specific purpose
func (s *authService) VerifyOTP(ctx context.Context, req *domain.VerifyOTPRequest) (string, error) {
	// Convert purpose string to OTPPurpose type
	var otpPurpose util.OTPPurpose
	switch req.Purpose {
	case string(util.OTPPurposeRegistration):
		otpPurpose = util.OTPPurposeRegistration
	case string(util.OTPPurposePasswordReset):
		otpPurpose = util.OTPPurposePasswordReset
	case string(util.OTPPurposeLogin):
		otpPurpose = util.OTPPurposeLogin
	default:
		return "", fmt.Errorf("invalid OTP purpose: %s", req.Purpose)
	}

	// Verify OTP
	valid, err := s.otpManager.VerifyOTP(req.Email, req.OTP, otpPurpose)
	if err != nil {
		return "", fmt.Errorf("OTP verification failed: %w", err)
	}

	if !valid {
		return "", fmt.Errorf("invalid OTP")
	}

	// Get user by email
	_, err = s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// For purposes that require generating a token
	if otpPurpose == util.OTPPurposeLogin {
		// Create a short-lived token for the login flow
		token := uuid.New().String()
		// Store this token somewhere if needed for verification

		return token, nil
	}

	// Generate a generic verification token
	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s:%d",
		req.Email, uuid.New().String(), time.Now().Unix())))

	return token, nil
}

// CheckAccountExists checks if an account with the given email exists
func (s *authService) CheckAccountExists(ctx context.Context, email string) (bool, error) {
	// Kiểm tra tham số
	if email == "" {
		return false, fmt.Errorf("email is required")
	}

	// Kiểm tra tài khoản tồn tại trong cơ sở dữ liệu
	exists, err := s.userRepo.UserExists(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

// createToken generates access and refresh tokens
func (s *authService) createToken(userID, email, role string) (*domain.TokenDetails, error) {
	td := &domain.TokenDetails{
		AccessUUID:  uuid.New().String(),
		RefreshUUID: uuid.New().String(),
		AtExpires:   time.Now().Add(s.accessDuration).Unix(),
		RtExpires:   time.Now().Add(s.refreshDuration).Unix(),
	}

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
		return nil, err
	}
	td.RefreshToken = refreshToken

	return td, nil
}

// Helper function for debug
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Logout logs out a user by invalidating their sessions
func (s *authService) Logout(ctx context.Context, userID, email string) error {
	// Invalidate refresh token by setting it to empty
	if err := s.userRepo.UpdateRefreshToken(ctx, userID, ""); err != nil {
		return fmt.Errorf("failed to invalidate refresh token: %w", err)
	}

	// Đăng xuất thành công, ghi log
	log.Printf("User %s (%s) logged out successfully from current session", userID, email)

	// Phát event log
	metadata := map[string]interface{}{
		"logout_type": "current_session",
	}
	if s.eventEmitter != nil {
		if err := s.eventEmitter.EmitLogout(userID, email, metadata); err != nil {
			log.Printf("Warning: Failed to emit logout event: %v", err)
			// Continue even if event emission fails
		}
	}

	return nil
}

// LogoutFromAllDevices logs out a user from all devices
func (s *authService) LogoutFromAllDevices(ctx context.Context, userID, email string) error {
	// Invalidate all refresh tokens
	if err := s.userRepo.LogoutFromAllDevices(ctx, userID); err != nil {
		return fmt.Errorf("failed to logout from all devices: %w", err)
	}

	// Đăng xuất thành công từ tất cả thiết bị, ghi log
	log.Printf("User %s (%s) logged out successfully from all devices", userID, email)

	// Phát event log với metadata chỉ rõ đã đăng xuất từ tất cả thiết bị
	metadata := map[string]interface{}{
		"logout_type": "all_devices",
	}
	if s.eventEmitter != nil {
		if err := s.eventEmitter.EmitLogout(userID, email, metadata); err != nil {
			log.Printf("Warning: Failed to emit logout event: %v", err)
			// Continue even if event emission fails
		}
	}

	return nil
}

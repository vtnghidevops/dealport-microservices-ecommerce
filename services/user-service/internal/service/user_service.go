package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"user-service/internal/domain"
	"user-service/internal/logging"
	"user-service/internal/repository"
)

// userService implements the UserService interface
type userService struct {
	userRepo     repository.UserRepository
	loggerClient *logging.LoggerClient
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, loggerClient *logging.LoggerClient) UserService {
	return &userService{
		userRepo:     userRepo,
		loggerClient: loggerClient,
	}
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Load wishlist items
	wishlistItems, err := s.userRepo.GetWishlist(ctx, id)
	if err != nil {
		// Log error but continue
		fmt.Printf("Warning: failed to get wishlist for user %s: %v\n", id, err)
	} else {
		// Convert wishlist items to product IDs array
		user.Wishlist = make([]int, len(wishlistItems))
		for i, item := range wishlistItems {
			user.Wishlist[i] = item.ProductID
		}
	}

	// Don't return password hash
	user.PasswordHash = ""

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Load wishlist items
	wishlistItems, err := s.userRepo.GetWishlist(ctx, user.ID)
	if err != nil {
		// Log error but continue
		fmt.Printf("Warning: failed to get wishlist for user %s: %v\n", user.ID, err)
	} else {
		// Convert wishlist items to product IDs array
		user.Wishlist = make([]int, len(wishlistItems))
		for i, item := range wishlistItems {
			user.Wishlist[i] = item.ProductID
		}
	}

	// Don't return password hash
	user.PasswordHash = ""

	return user, nil
}

// GetUsers retrieves a list of users with pagination and filtering
func (s *userService) GetUsers(ctx context.Context, filter *domain.UserFilter) ([]*domain.User, int, error) {
	users, total, err := s.userRepo.GetUsers(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	// Load wishlist for each user and remove password hashes
	for _, user := range users {
		// Load wishlist items
		wishlistItems, err := s.userRepo.GetWishlist(ctx, user.ID)
		if err != nil {
			// Log error but continue
			fmt.Printf("Warning: failed to get wishlist for user %s: %v\n", user.ID, err)
		} else {
			// Convert wishlist items to product IDs array
			user.Wishlist = make([]int, len(wishlistItems))
			for i, item := range wishlistItems {
				user.Wishlist[i] = item.ProductID
			}
		}

		// Don't return password hash
		user.PasswordHash = ""
	}

	return users, total, nil
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	// Check if user already exists
	exists, err := s.userRepo.UserExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password if provided
	var hashedPassword string
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		hashedPassword = string(hashed)
	}

	// Set default values
	id := uuid.New().String()

	// If ID is provided in request (e.g., from auth service), use it
	if req.ID != "" {
		id = req.ID
	}

	now := time.Now()

	if req.DisplayName == "" {
		req.DisplayName = fmt.Sprintf("%s %s", req.FirstName, req.LastName)
	}

	// Generate username if not provided (lastName + firstName)
	if req.Username == "" {
		req.Username = fmt.Sprintf("%s%s", req.LastName, req.FirstName)
	}

	// Default role
	if req.Role == "" {
		req.Role = "customer"
	}

	// Create user object
	user := &domain.User{
		ID:           id,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Username:     req.Username,
		DisplayName:  req.DisplayName,
		Addresses:    req.Addresses,
		Role:         req.Role,
		Status:       "active",
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Set pointer fields properly
	if req.Phone != "" {
		phone := req.Phone
		user.Phone = &phone
	}

	if req.ProfileImage != "" {
		profileImage := req.ProfileImage
		user.ProfileImage = &profileImage
	}

	// Create user in the database
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	user.PasswordHash = ""
	return user, nil

}

// UpdateUser updates user information
func (s *userService) UpdateUser(ctx context.Context, req *domain.UpdateUserRequest) (*domain.User, error) {
	// Validate ID
	if req.ID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get existing user
	user, err := s.userRepo.GetUserByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.Email != "" && req.Email != user.Email {
		user.Email = req.Email
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.Phone != "" {
		phone := req.Phone
		user.Phone = &phone
	}
	if req.ProfileImage != "" {
		profileImage := req.ProfileImage
		user.ProfileImage = &profileImage
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Gender != "" {
		gender := req.Gender
		user.Gender = &gender
	}
	if req.CartID != "" {
		cartID := req.CartID
		user.CartID = &cartID
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	user.Active = req.Active
	user.UpdatedAt = time.Now()

	// Update addresses if provided
	if len(req.Addresses) > 0 {
		// Replace existing addresses
		user.Addresses = req.Addresses
	}

	// Update payment methods if provided
	if len(req.PaymentMethods) > 0 {
		// Replace existing payment methods
		user.PaymentMethods = req.PaymentMethods
	}

	// Update user in the database
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Log profile update activity
	if s.loggerClient != nil {
		metadata := map[string]interface{}{
			"updated_fields": getUpdatedFields(req),
		}
		go s.loggerClient.LogUserActivity(context.Background(), "profile_updated", user.ID, "User profile updated", metadata)
	}

	user.PasswordHash = ""
	return user, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	if err := s.userRepo.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// SearchUsers searches for users based on criteria
func (s *userService) SearchUsers(ctx context.Context, params *domain.SearchParams) ([]*domain.User, int, error) {
	users, total, err := s.userRepo.SearchUsers(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	// Don't return password hashes
	for _, user := range users {
		user.PasswordHash = ""
	}

	return users, total, nil
}

// CreateAddress creates a new address for a user
func (s *userService) CreateAddress(ctx context.Context, userID string, address *domain.Address) error {
	// Ensure the user exists
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Set user ID for the address
	address.UserID = userID

	// Create address
	if err := s.userRepo.CreateAddress(ctx, address); err != nil {
		return fmt.Errorf("failed to create address: %w", err)
	}

	return nil
}

// UpdateAddress updates an address
func (s *userService) UpdateAddress(ctx context.Context, address *domain.Address) error {
	if err := s.userRepo.UpdateAddress(ctx, address); err != nil {
		return fmt.Errorf("failed to update address: %w", err)
	}
	return nil
}

// DeleteAddress deletes an address
func (s *userService) DeleteAddress(ctx context.Context, id string, userID string) error {
	if err := s.userRepo.DeleteAddress(ctx, id, userID); err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}
	return nil
}

// SetDefaultAddress sets an address as default
func (s *userService) SetDefaultAddress(ctx context.Context, id string, userID string) error {
	if err := s.userRepo.SetDefaultAddress(ctx, id, userID); err != nil {
		return fmt.Errorf("failed to set default address: %w", err)
	}
	return nil
}

// AddToWishlist adds a product to a user's wishlist
func (s *userService) AddToWishlist(ctx context.Context, req *domain.AddToWishlistRequest) error {
	// Validate input
	if req.UserID == "" {
		return errors.New("user ID is required")
	}

	// Check if user exists
	_, err := s.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Add to wishlist
	return s.userRepo.AddToWishlist(ctx, req.UserID, req.ProductID, req.Notes)
}

// RemoveFromWishlist removes a product from a user's wishlist
func (s *userService) RemoveFromWishlist(ctx context.Context, req *domain.RemoveFromWishlistRequest) error {
	// Validate input
	if req.UserID == "" {
		return errors.New("user ID is required")
	}

	// Remove from wishlist
	return s.userRepo.RemoveFromWishlist(ctx, req.UserID, req.ProductID)
}

// GetWishlist retrieves a user's wishlist
func (s *userService) GetWishlist(ctx context.Context, req *domain.GetWishlistRequest) (*domain.GetWishlistResponse, error) {
	// Validate input
	if req.UserID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get wishlist items
	items, err := s.userRepo.GetWishlist(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist: %w", err)
	}

	// Convert pointer slice to value slice
	wishlistItems := make([]domain.WishlistItem, len(items))
	for i, item := range items {
		wishlistItems[i] = *item
	}

	// Prepare response
	response := &domain.GetWishlistResponse{
		Items: wishlistItems,
		Count: len(wishlistItems),
	}

	return response, nil
}

// ChangePassword changes a user's password
func (s *userService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return errors.New("invalid current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	// Save to database
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Log password change activity
	if s.loggerClient != nil {
		go s.loggerClient.LogUserActivity(context.Background(), "password_changed", user.ID, "Password changed", nil)
	}

	return nil
}

// RequestPasswordReset initiates a password reset request
func (s *userService) RequestPasswordReset(ctx context.Context, email string) error {
	// Check if user exists
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Log password reset request activity
	if s.loggerClient != nil {
		go s.loggerClient.LogUserActivity(context.Background(), "password_reset_requested", user.ID, "Password reset requested", nil)
	}

	return nil
}

// ValidateCredentials validates user login credentials
func (s *userService) ValidateCredentials(ctx context.Context, email, password string) (*domain.User, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Don't return password hash
	user.PasswordHash = ""

	return user, nil
}

// LogoutUser logs out a user
func (s *userService) LogoutUser(ctx context.Context, userID string) error {
	// Get user to validate ID
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Log logout activity
	if s.loggerClient != nil {
		go s.loggerClient.LogUserActivity(context.Background(), "logout", user.ID, "User logged out", nil)
	}

	return nil
}

// LogUserActivity logs a user activity
func (s *userService) LogUserActivity(ctx context.Context, action, userID, message string, metadata map[string]interface{}) error {
	if s.loggerClient == nil {
		return errors.New("logger client not initialized")
	}

	return s.loggerClient.LogUserActivity(ctx, action, userID, message, metadata)
}

// GetUserActivityLogs retrieves user activity logs
func (s *userService) GetUserActivityLogs(ctx context.Context, userID string, actionType string) (interface{}, error) {
	if s.loggerClient == nil {
		return nil, errors.New("logger client not initialized")
	}

	return s.loggerClient.GetUserActivityLogs(ctx, userID, actionType)
}

// Helper function to get updated fields for logging
func getUpdatedFields(req *domain.UpdateUserRequest) map[string]interface{} {
	fields := make(map[string]interface{})

	if req.Email != "" {
		fields["email"] = req.Email
	}
	if req.FirstName != "" {
		fields["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		fields["last_name"] = req.LastName
	}
	if req.DisplayName != "" {
		fields["display_name"] = req.DisplayName
	}
	if req.Phone != "" {
		fields["phone"] = req.Phone
	}
	if req.ProfileImage != "" {
		fields["profile_image"] = req.ProfileImage
	}
	if req.Username != "" {
		fields["username"] = req.Username
	}
	if req.Gender != "" {
		fields["gender"] = req.Gender
	}
	if len(req.Addresses) > 0 {
		fields["addresses_updated"] = true
	}
	if len(req.PaymentMethods) > 0 {
		fields["payment_methods_updated"] = true
	}

	return fields
}

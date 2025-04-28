package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"user-service/internal/domain"
	"user-service/internal/repository"
)

// userService implements the UserService interface
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
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

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Set default values
	id := uuid.New().String()
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
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Username:     req.Username,
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

	// Don't return password hash
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
	existingUser, err := s.userRepo.GetUserByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Track if first name or last name changes to update display name if needed
	firstNameChanged := false
	lastNameChanged := false

	// Update fields that are provided
	if req.Email != "" {
		existingUser.Email = req.Email
	}

	if req.FirstName != "" {
		firstNameChanged = req.FirstName != existingUser.FirstName
		existingUser.FirstName = req.FirstName
	}

	if req.LastName != "" {
		lastNameChanged = req.LastName != existingUser.LastName
		existingUser.LastName = req.LastName
	}

	if req.DisplayName != "" {
		existingUser.DisplayName = req.DisplayName
	} else if firstNameChanged || lastNameChanged {
		// Update display name if first name or last name changed and display name was not provided
		existingUser.DisplayName = fmt.Sprintf("%s %s", existingUser.FirstName, existingUser.LastName)
	}

	// Update username if provided or if first name/last name changed and custom username wasn't set
	if req.Username != "" {
		existingUser.Username = req.Username
	} else if (firstNameChanged || lastNameChanged) && (existingUser.Username == "" ||
		existingUser.Username == fmt.Sprintf("%s%s", existingUser.LastName, existingUser.FirstName)) {
		// Only update automatically if username follows the default pattern or is empty
		existingUser.Username = fmt.Sprintf("%s%s", existingUser.LastName, existingUser.FirstName)
	}

	if req.Phone != "" {
		phone := req.Phone
		existingUser.Phone = &phone
	}

	if req.ProfileImage != "" {
		profileImage := req.ProfileImage
		existingUser.ProfileImage = &profileImage
	}

	if req.Role != "" {
		existingUser.Role = req.Role
	}

	// Update status if provided
	if req.Status != "" {
		existingUser.Status = req.Status
	}

	// Update active status
	existingUser.Active = req.Active

	// Update addresses if provided
	if len(req.Addresses) > 0 {
		existingUser.Addresses = req.Addresses
	}

	// Update timestamp
	existingUser.UpdatedAt = time.Now()

	// Update in the database
	if err := s.userRepo.UpdateUser(ctx, existingUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Don't return password hash
	existingUser.PasswordHash = ""

	return existingUser, nil
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

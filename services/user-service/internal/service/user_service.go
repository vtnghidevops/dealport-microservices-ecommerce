package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"user-service/internal/domain"
	"user-service/internal/event"
	"user-service/internal/repository"
)

// userService implements the UserService interface
type userService struct {
	userRepo     repository.UserRepository
	eventEmitter *event.EventEmitter
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, eventEmitter *event.EventEmitter) UserService {
	return &userService{
		userRepo:     userRepo,
		eventEmitter: eventEmitter,
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

	// // Hash password if provided
	// var hashedPassword string
	// if req.Password != "" {
	// 	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to hash password: %w", err)
	// 	}
	// 	hashedPassword = string(hashed)
	// }

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
		ID:          id,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Username:    req.Username,
		DisplayName: req.DisplayName,
		Addresses:   req.Addresses,
		Role:        req.Role,
		Status:      "active",
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
		// ID:           id,
		// Email:        req.Email,
		// PasswordHash: hashedPassword,
		// FirstName:    req.FirstName,
		// LastName:     req.LastName,
		// Username:     req.Username,
		// DisplayName:  req.DisplayName,
		// Addresses:    req.Addresses,
		// Role:         req.Role,
		// Status:       "active",
		// Active:       true,
		// CreatedAt:    now,
		// UpdatedAt:    now,
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

	// Profile update activity will be handled by event emitter below

	// Emit profile updated event
	if s.eventEmitter != nil {
		updatedFields := getUpdatedFields(req)
		go func() {
			if err := s.eventEmitter.EmitProfileUpdatedEvent(context.Background(), user.ID, updatedFields); err != nil {
				// Note: We don't return error here to avoid breaking the update flow
				fmt.Printf("Warning: Failed to emit profile updated event: %v\n", err)
			}
		}()
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

// // ChangePassword changes a user's password
// func (s *userService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
// 	// Get user
// 	user, err := s.userRepo.GetUserByID(ctx, userID)
// 	if err != nil {
// 		return fmt.Errorf("failed to get user: %w", err)
// 	}

// 	// Verify old password
// 	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
// 	if err != nil {
// 		return errors.New("invalid current password")
// 	}

// 	// Hash new password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
// 	if err != nil {
// 		return fmt.Errorf("failed to hash password: %w", err)
// 	}

// 	// Update password
// 	user.PasswordHash = string(hashedPassword)
// 	user.UpdatedAt = time.Now()

// 	// Save to database
// 	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
// 		return fmt.Errorf("failed to update user: %w", err)
// 	}

// 	// Password change activity will be handled by event emitter below

// 	// Emit password changed event
// 	if s.eventEmitter != nil {
// 		go func() {
// 			if err := s.eventEmitter.EmitPasswordChangedEvent(context.Background(), user.ID); err != nil {
// 				// Note: We don't return error here to avoid breaking the flow
// 				fmt.Printf("Warning: Failed to emit password changed event: %v\n", err)
// 			}
// 		}()
// 	}

// 	return nil
// }

// RequestPasswordReset initiates a password reset request
func (s *userService) RequestPasswordReset(ctx context.Context, email string) error {
	// Check if user exists
	_, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Password reset request activity will be handled by event emitter if needed
	// For now, just return nil as this is mainly for validation
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

	// Logout activity will be handled by event emitter below

	// Emit logout event
	if s.eventEmitter != nil {
		go func() {
			if err := s.eventEmitter.EmitLogoutEvent(context.Background(), user.ID, user.Email); err != nil {
				// Note: We don't return error here to avoid breaking the flow
				fmt.Printf("Warning: Failed to emit logout event: %v\n", err)
			}
		}()
	}

	return nil
}

// LogUserActivity logs a user activity using events
func (s *userService) LogUserActivity(ctx context.Context, action, userID, message string, metadata map[string]interface{}) error {
	// Use event emitter for logging instead of direct logger client
	if s.eventEmitter != nil {
		return s.eventEmitter.EmitUserActivityEvent(ctx, action, userID, message, metadata)
	}
	// Return nil instead of error to avoid breaking test flows when event emitter is not initialized
	return nil
}

// GetUserActivityLogs retrieves user activity logs
func (s *userService) GetUserActivityLogs(ctx context.Context, userID string, actionType string) (interface{}, error) {
	// This functionality is now handled by logger-service through events
	// For now, return empty result as this needs to be called via logger service API
	return nil, errors.New("activity logs are now handled by logger service - use logger service API directly")
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

// GetUserStatistics retrieves statistics for the admin dashboard
func (s *userService) GetUserStatistics(ctx context.Context) (*domain.UserStatistics, error) {
	// Get total users count
	filter := &domain.UserFilter{
		Page:  1,
		Limit: 1,
	}
	_, totalUsers, err := s.userRepo.GetUsers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}

	// Get new users in the last 7 days
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	newUsersCount, err := s.userRepo.GetNewUsersCount(ctx, sevenDaysAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to get new users count: %w", err)
	}

	// Get new users in the last 14 days (for calculating growth)
	fourteenDaysAgo := time.Now().AddDate(0, 0, -14)
	previousPeriodNewUsers, err := s.userRepo.GetNewUsersCount(ctx, fourteenDaysAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous period new users: %w", err)
	}
	previousPeriodNewUsers -= newUsersCount // Exclude the last 7 days

	// Calculate new user growth rate
	var newUserGrowth float64
	if previousPeriodNewUsers > 0 {
		newUserGrowth = float64(newUsersCount-previousPeriodNewUsers) / float64(previousPeriodNewUsers) * 100
	} else {
		newUserGrowth = 100 // If there were no users in the previous period, set growth to 100%
	}

	// Get active users (users who have logged in in the last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	activeUsersCount, err := s.userRepo.GetActiveUsersCount(ctx, thirtyDaysAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to get active users count: %w", err)
	}

	// Get the total active users from previous period
	sixtyDaysAgo := time.Now().AddDate(0, 0, -60)
	previousActiveUsersCount, err := s.userRepo.GetActiveUsersCount(ctx, sixtyDaysAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous active users count: %w", err)
	}
	previousActiveUsersCount -= activeUsersCount // Exclude the last 30 days

	// Calculate user growth rate
	var userGrowth float64
	if previousActiveUsersCount > 0 {
		userGrowth = float64(activeUsersCount-previousActiveUsersCount) / float64(previousActiveUsersCount) * 100
	} else {
		userGrowth = 100 // If there were no active users in the previous period, set growth to 100%
	}

	// Get repeat customers (users with more than one order)
	repeatCustomers, err := s.userRepo.GetRepeatCustomers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get repeat customers: %w", err)
	}

	// For visitor statistics, we'd typically use analytics data
	// For demo purposes, we'll use some reasonable estimates based on user counts
	visitors := totalUsers * 5                   // Assume 5 visitors per user on average
	visitorGrowth := userGrowth * 1.2            // Assume visitor growth is 20% higher than user growth
	shopVisitors := int(float64(visitors) * 0.7) // Assume 70% of visitors go to the shop

	// Calculate conversion rate (percentage of visitors who become active users)
	var conversionRate float64
	if visitors > 0 {
		conversionRate = float64(activeUsersCount) / float64(visitors) * 100
	}

	return &domain.UserStatistics{
		TotalUsers:      totalUsers,
		UserGrowth:      userGrowth,
		NewUsers:        newUsersCount,
		NewUserGrowth:   newUserGrowth,
		Visitors:        visitors,
		VisitorGrowth:   visitorGrowth,
		ActiveUsers:     activeUsersCount,
		RepeatCustomers: len(repeatCustomers),
		ShopVisitors:    shopVisitors,
		ConversionRate:  conversionRate,
	}, nil
}

// GetUserActivityChart retrieves data for the customer activity chart
func (s *userService) GetUserActivityChart(ctx context.Context, days int) (*domain.UserActivityChart, error) {
	if days <= 0 {
		days = 7 // Default to 7 days
	}

	// Default to activity chart type
	chartType := string(domain.ChartTypeActivity)

	// Validate chart type
	var selectedChartType domain.ChartType

	switch chartType {
	case string(domain.ChartTypeActivity):
		selectedChartType = domain.ChartTypeActivity
	case string(domain.ChartTypeNewUsers):
		selectedChartType = domain.ChartTypeNewUsers
	case string(domain.ChartTypeOrders):
		selectedChartType = domain.ChartTypeOrders
	case string(domain.ChartTypeRevenue):
		selectedChartType = domain.ChartTypeRevenue
	case string(domain.ChartTypeConversion):
		selectedChartType = domain.ChartTypeConversion
	default:
		// Default to activity chart
		selectedChartType = domain.ChartTypeActivity
	}

	// Create slice to hold chart data
	chartData := make([]domain.UserActivityChartPoint, 0, days)

	// Get current date
	now := time.Now()

	// Initialize response with default values
	result := &domain.UserActivityChart{
		ChartType:  selectedChartType,
		ChartData:  chartData,
		Title:      "User Activity",
		YAxisLabel: "Count",
	}

	// Track aggregate values for analysis
	var totalCount int
	var totalValue float64
	maxCount := 0
	minCount := -1 // -1 means uninitialized
	maxValue := 0.0
	minValue := -1.0 // -1 means uninitialized

	// For each day in the requested period
	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)

		// Format the day and date
		day := date.Format("Mon")
		dateStr := date.Format("2006-01-02")

		dataPoint := domain.UserActivityChartPoint{
			Day:  day,
			Date: dateStr,
		}

		// Generate data based on chart type
		switch selectedChartType {
		case domain.ChartTypeActivity:
			// Get user activity count for this day
			count, err := s.userRepo.GetUserActivityCountForDay(ctx, date)
			if err != nil {
				return nil, fmt.Errorf("failed to get user activity for day %s: %w", dateStr, err)
			}
			dataPoint.Count = count

			// Calculate session duration as secondary value (mocked for demo purposes)
			// In a real implementation, this would come from analytics data
			dataPoint.Value = float64(20 + rand.Intn(40)) // Average session duration in minutes

		case domain.ChartTypeNewUsers:
			// Get new users registered on this day
			startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
			endOfDay := startOfDay.Add(24 * time.Hour)

			// Get new users count for this specific day
			count, err := s.getNewUsersForTimeRange(ctx, startOfDay, endOfDay)
			if err != nil {
				return nil, fmt.Errorf("failed to get new users for day %s: %w", dateStr, err)
			}
			dataPoint.Count = count

			// Set retention rate as secondary value (mocked for demo purposes)
			dataPoint.Value = 70.0 + float64(rand.Intn(30)) // Retention percentage

		case domain.ChartTypeOrders:
			// For orders, we'd typically query the checkout service
			// But for demo purposes, we'll generate realistic order data
			baseOrderCount := 50
			weekday := date.Weekday()

			// Weekends typically have more orders
			if weekday == time.Saturday {
				baseOrderCount = 80
			} else if weekday == time.Sunday {
				baseOrderCount = 70
			} else if weekday == time.Friday {
				baseOrderCount = 65
			}

			// Add some randomness
			variance := baseOrderCount / 5
			orderCount := baseOrderCount + rand.Intn(variance*2) - variance

			dataPoint.Count = orderCount

			// Set average order value as secondary value
			dataPoint.Value = 50.0 + rand.Float64()*100.0 // Average order value in $

		case domain.ChartTypeRevenue:
			// For revenue, we'd typically query the checkout service
			// But for demo purposes, we'll generate realistic revenue data
			baseOrderCount := 50
			weekday := date.Weekday()

			// Weekends typically have more orders
			if weekday == time.Saturday {
				baseOrderCount = 80
			} else if weekday == time.Sunday {
				baseOrderCount = 70
			} else if weekday == time.Friday {
				baseOrderCount = 65
			}

			// Add some randomness
			variance := baseOrderCount / 5
			orderCount := baseOrderCount + rand.Intn(variance*2) - variance

			// Average order value with some randomness
			avgOrderValue := 50.0 + rand.Float64()*100.0

			// Calculate revenue
			revenue := float64(orderCount) * avgOrderValue

			dataPoint.Count = orderCount
			dataPoint.Value = revenue

		case domain.ChartTypeConversion:
			// For conversion rates, we'd typically get from analytics
			// But for demo purposes, we'll generate realistic conversion data

			// Get number of visitors (mocked)
			visitors := 1000 + rand.Intn(500)

			// Get number of conversions (orders or signups)
			conversions := int(float64(visitors) * (0.02 + rand.Float64()*0.05))

			dataPoint.Count = conversions
			dataPoint.Value = float64(conversions) * 100.0 / float64(visitors) // Conversion rate as percentage
		}

		// Add to chart data
		chartData = append(chartData, dataPoint)

		// Update aggregates for analysis
		totalCount += dataPoint.Count
		totalValue += dataPoint.Value

		if maxCount < dataPoint.Count {
			maxCount = dataPoint.Count
		}
		if minCount == -1 || minCount > dataPoint.Count {
			minCount = dataPoint.Count
		}

		if maxValue < dataPoint.Value {
			maxValue = dataPoint.Value
		}
		if minValue == -1 || minValue > dataPoint.Value {
			minValue = dataPoint.Value
		}
	}

	// Calculate growth rate (comparing first half to second half of the period)
	var growthRate float64
	if days > 1 {
		midpoint := days / 2
		var firstHalfCount, secondHalfCount int

		for i, point := range chartData {
			if i < midpoint {
				firstHalfCount += point.Count
			} else {
				secondHalfCount += point.Count
			}
		}

		if firstHalfCount > 0 {
			growthRate = float64(secondHalfCount-firstHalfCount) * 100.0 / float64(firstHalfCount)
		}
	}

	// Set chart metadata based on chart type
	switch selectedChartType {
	case domain.ChartTypeActivity:
		result.Title = "User Activity"
		result.YAxisLabel = "Active Users"
		result.Description = "Daily active users over the selected period"
	case domain.ChartTypeNewUsers:
		result.Title = "New User Registrations"
		result.YAxisLabel = "New Users"
		result.Description = "New user sign-ups over the selected period"
	case domain.ChartTypeOrders:
		result.Title = "Order Volume"
		result.YAxisLabel = "Orders"
		result.Description = "Number of orders placed over the selected period"
	case domain.ChartTypeRevenue:
		result.Title = "Revenue"
		result.YAxisLabel = "Revenue ($)"
		result.Description = "Daily revenue over the selected period"
	case domain.ChartTypeConversion:
		result.Title = "Conversion Rate"
		result.YAxisLabel = "Conversion (%)"
		result.Description = "Daily visitor-to-customer conversion rate"
	}

	// Set statistical data
	result.ChartData = chartData
	result.TotalValue = totalValue
	result.AvgValue = totalValue / float64(len(chartData))
	result.MaxValue = maxValue
	if minValue >= 0 {
		result.MinValue = minValue
	}
	result.GrowthRate = growthRate

	return result, nil
}

// Helper method to get new users for a specific time range
func (s *userService) getNewUsersForTimeRange(ctx context.Context, start, end time.Time) (int, error) {
	// Mock implementation - trong triển khai thực tế sẽ truy vấn database
	// Trong demo này, chúng ta sẽ sử dụng dữ liệu giả lập

	// Generate patterns based on day of week - higher on weekends, lower on weekdays
	dayOfWeek := start.Weekday()
	count := 0

	switch dayOfWeek {
	case time.Saturday, time.Sunday:
		// Weekend has more sign-ups
		count = 30 + rand.Intn(20)
	case time.Friday:
		// Friday has moderate-high sign-ups
		count = 25 + rand.Intn(15)
	case time.Monday:
		// Monday has moderate sign-ups
		count = 20 + rand.Intn(10)
	default:
		// Tuesday-Thursday have moderate-low sign-ups
		count = 15 + rand.Intn(10)
	}

	return count, nil
}

// SyncUserOrderData synchronizes user order statistics
func (s *userService) SyncUserOrderData(ctx context.Context, userID string, orderCount int, totalSpend float64) error {
	// Validate userID
	if userID == "" {
		return errors.New("user ID is required")
	}

	// Check if user exists
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Update order count
	err = s.userRepo.UpdateUserOrderCount(ctx, userID, orderCount)
	if err != nil {
		return fmt.Errorf("failed to update order count: %w", err)
	}

	// Note: totalSpend is not stored in database, only used for logging/events
	// Order sync activity can be logged via events if needed
	// For now, we'll skip logging this internal sync operation
	return nil
}

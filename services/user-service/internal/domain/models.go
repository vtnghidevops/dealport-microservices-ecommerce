package domain

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           string     `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	FirstName    string     `json:"first_name" db:"first_name"`
	LastName     string     `json:"last_name" db:"last_name"`
	DisplayName  string     `json:"display_name" db:"display_name"`
	PasswordHash string     `json:"-" db:"-"`                                   // Marked as non-database field
	Phone        *string    `json:"phone,omitempty" db:"phone"`                 // Changed to pointer for NULL handling
	ProfileImage *string    `json:"profile_image,omitempty" db:"profile_image"` // Changed to pointer for NULL handling
	Addresses    []Address  `json:"addresses"`
	Role         string     `json:"role" db:"role"`
	Status       string     `json:"status" db:"status"` // 'active', 'inactive', 'suspended', 'pending'
	Active       bool       `json:"active" db:"active"` // Kept for backward compatibility
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin    *time.Time `json:"last_login,omitempty" db:"last_login"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	// Added fields to match frontend requirements
	Username       string          `json:"username,omitempty" db:"username"`
	DateOfBirth    *time.Time      `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Gender         *string         `json:"gender,omitempty" db:"gender"` // Changed to pointer for NULL handling
	PaymentMethods []PaymentMethod `json:"payment_methods,omitempty"`
	Wishlist       []int           `json:"wishlist,omitempty"`
	CartID         *string         `json:"cart_id,omitempty" db:"cart_id"` // Changed to pointer for NULL handling
	OrderCount     int             `json:"order_count,omitempty" db:"order_count"`
	TotalSpend     float64         `json:"total_spend,omitempty" db:"total_spend"`
}

// Address represents a user address
type Address struct {
	ID          string `json:"id" db:"id"`
	UserID      string `json:"user_id" db:"user_id"`
	Name        string `json:"name" db:"name"`   // Tên người nhận
	Phone       string `json:"phone" db:"phone"` // SĐT liên hệ
	Line1       string `json:"line1" db:"line1"` // Địa chỉ
	City        string `json:"city" db:"city"`
	State       string `json:"state" db:"state"`
	PostalCode  string `json:"postal_code" db:"postal_code"` // Renamed to match zipCode in frontend
	Country     string `json:"country" db:"country"`
	IsDefault   bool   `json:"is_default" db:"is_default"`
	AddressType string `json:"address_type" db:"address_type"` // Renamed to match type in frontend
}

// PaymentMethod represents a user payment method
type PaymentMethod struct {
	ID            string `json:"id" db:"id"`
	UserID        string `json:"user_id" db:"user_id"`
	Type          string `json:"type" db:"type"` // 'credit_card', 'paypal', 'bank_transfer', 'other'
	Provider      string `json:"provider,omitempty" db:"provider"`
	AccountNumber string `json:"account_number,omitempty" db:"account_number"`
	ExpiryDate    string `json:"expiry_date,omitempty" db:"expiry_date"`
	IsDefault     bool   `json:"is_default" db:"is_default"`
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	ID             string          `json:"id"` // Added to support direct ID from auth service
	Email          string          `json:"email" validate:"required,email"`
	Password       string          `json:"password" validate:"required,min=8"`
	FirstName      string          `json:"first_name" validate:"required"`
	LastName       string          `json:"last_name" validate:"required"`
	DisplayName    string          `json:"display_name"`
	Phone          string          `json:"phone"`
	ProfileImage   string          `json:"profile_image"`
	Addresses      []Address       `json:"addresses"`
	Role           string          `json:"role"`
	Username       string          `json:"username,omitempty"`
	DateOfBirth    *time.Time      `json:"date_of_birth,omitempty"`
	Gender         string          `json:"gender,omitempty"`
	PaymentMethods []PaymentMethod `json:"payment_methods,omitempty"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	ID             string          `json:"id" validate:"required"`
	Email          string          `json:"email" validate:"omitempty,email"`
	FirstName      string          `json:"first_name"`
	LastName       string          `json:"last_name"`
	DisplayName    string          `json:"display_name"`
	Phone          string          `json:"phone"`
	ProfileImage   string          `json:"profile_image"`
	Addresses      []Address       `json:"addresses"`
	Role           string          `json:"role"`
	Status         string          `json:"status"`
	Active         bool            `json:"active"`
	Username       string          `json:"username,omitempty"`
	DateOfBirth    *time.Time      `json:"date_of_birth,omitempty"`
	Gender         string          `json:"gender,omitempty"`
	PaymentMethods []PaymentMethod `json:"payment_methods,omitempty"`
	CartID         string          `json:"cart_id,omitempty"`
}

// UserFilter represents filter parameters for user queries
type UserFilter struct {
	Page      int    `json:"page" validate:"min=1"`
	Limit     int    `json:"limit" validate:"min=1,max=100"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order" validate:"oneof=asc desc"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	Active    *bool  `json:"active,omitempty"` // Added to filter by active status
}

// SearchParams represents search parameters
type SearchParams struct {
	Query string `json:"query" validate:"required"`
	Field string `json:"field"`
	Page  int    `json:"page" validate:"min=1"`
	Limit int    `json:"limit" validate:"min=1,max=100"`
}

// WishlistItem represents an item in a user's wishlist
type WishlistItem struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	ProductID int       `json:"product_id" db:"product_id"`
	AddedAt   time.Time `json:"added_at" db:"added_at"`
	Notes     string    `json:"notes,omitempty" db:"notes"`
}

// AddToWishlistRequest represents a request to add an item to wishlist
type AddToWishlistRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	ProductID int    `json:"product_id" validate:"required"`
	Notes     string `json:"notes,omitempty"`
}

// RemoveFromWishlistRequest represents a request to remove an item from wishlist
type RemoveFromWishlistRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	ProductID int    `json:"product_id" validate:"required"`
}

// GetWishlistRequest represents a request to get a user's wishlist
type GetWishlistRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

// GetWishlistResponse represents a response with wishlist items
type GetWishlistResponse struct {
	Items []WishlistItem `json:"items"`
	Count int            `json:"count"`
}

// UserStatistics represents the statistics for the admin dashboard
type UserStatistics struct {
	TotalUsers      int     `json:"total_users"`
	UserGrowth      float64 `json:"user_growth"`
	NewUsers        int     `json:"new_users"`
	NewUserGrowth   float64 `json:"new_user_growth"`
	Visitors        int     `json:"visitors"`
	VisitorGrowth   float64 `json:"visitor_growth"`
	ActiveUsers     int     `json:"active_users"`
	RepeatCustomers int     `json:"repeat_customers"`
	ShopVisitors    int     `json:"shop_visitors"`
	ConversionRate  float64 `json:"conversion_rate"`
}

// UserActivityChartPoint represents a data point for the user activity chart
type UserActivityChartPoint struct {
	Day   string  `json:"day"`             // Day of week (e.g., "Mon")
	Date  string  `json:"date"`            // Full date (e.g., "2023-07-10")
	Count int     `json:"count"`           // Primary count value (users/orders/etc)
	Value float64 `json:"value,omitempty"` // Secondary value (revenue/conversion/etc)
}

// ChartType defines the type of chart to retrieve
type ChartType string

// Chart type constants
const (
	ChartTypeActivity   ChartType = "activity"   // User logins/activities
	ChartTypeNewUsers   ChartType = "new_users"  // New user registrations
	ChartTypeOrders     ChartType = "orders"     // Order counts
	ChartTypeRevenue    ChartType = "revenue"    // Revenue/sales
	ChartTypeConversion ChartType = "conversion" // Conversion rates
)

// UserActivityChart represents the data for the customer activity chart
type UserActivityChart struct {
	ChartData   []UserActivityChartPoint `json:"chart_data"`
	ChartType   ChartType                `json:"chart_type"`
	Title       string                   `json:"title"`
	YAxisLabel  string                   `json:"y_axis_label"`
	Description string                   `json:"description,omitempty"`
	TotalValue  float64                  `json:"total_value,omitempty"` // Total value across all data points
	AvgValue    float64                  `json:"avg_value,omitempty"`   // Average value across all data points
	MaxValue    float64                  `json:"max_value,omitempty"`   // Maximum value in the data set
	MinValue    float64                  `json:"min_value,omitempty"`   // Minimum value in the data set
	GrowthRate  float64                  `json:"growth_rate,omitempty"` // Growth rate over the period
}

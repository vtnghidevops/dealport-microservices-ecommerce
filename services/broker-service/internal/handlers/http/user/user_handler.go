package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"broker-service/internal/util"
	checkoutpb "broker-service/proto/checkout"
	userpb "broker-service/proto/user"

	"github.com/go-chi/chi/v5"
)

type Config struct {
	UserClient     userpb.UserServiceClient
	CheckoutClient checkoutpb.CheckoutServiceClient
}

// GetUser retrieves a specific user by ID
func (c *Config) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL
	userID := chi.URLParam(r, "id")
	if userID == "" {
		util.ErrorJSON(w, errors.New("user ID is required"), http.StatusBadRequest)
		return
	}

	// Call user service via gRPC
	res, err := c.UserClient.GetUser(r.Context(), &userpb.GetUserRequest{
		Id: userID,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Format user data for frontend
	userData := formatUserForFrontend(res.User)

	// Remove recentlyViewed field from response
	delete(userData, "recentlyViewed")

	responseData := util.JsonResponse{
		Error:   false,
		Message: "User retrieved successfully",
		Data:    userData,
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// GetUserProfile retrieves the current user's profile (from token)
func (c *Config) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Call user service via gRPC
	res, err := c.UserClient.GetUser(r.Context(), &userpb.GetUserRequest{
		Id: userID,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Format user data in a way that matches frontend expectations
	userData := formatUserForFrontend(res.User)

	// Remove recentlyViewed field from response
	delete(userData, "recentlyViewed")

	responseData := util.JsonResponse{
		Error:   false,
		Message: "User profile retrieved successfully",
		Data:    userData,
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// UpdateUserProfile updates the current user's profile
func (c *Config) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var requestPayload struct {
		Email        string                   `json:"email"`
		FirstName    string                   `json:"firstName"` // Frontend uses camelCase
		LastName     string                   `json:"lastName"`  // Frontend uses camelCase
		Username     string                   `json:"username"`
		Phone        string                   `json:"phone"`
		ProfileImage string                   `json:"avatar"` // Frontend uses 'avatar'
		DateOfBirth  string                   `json:"dateOfBirth"`
		Gender       string                   `json:"gender"`
		Addresses    []map[string]interface{} `json:"addresses"`
		Preferences  map[string]interface{}   `json:"preferences"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Convert addresses from frontend format to protobuf format
	addresses := make([]*userpb.Address, 0)
	for _, addr := range requestPayload.Addresses {
		isDefault := false
		if isDefaultVal, ok := addr["isDefault"].(bool); ok {
			isDefault = isDefaultVal
		}

		addressType := "shipping"
		if typeVal, ok := addr["type"].(string); ok {
			addressType = typeVal
		}

		addresses = append(addresses, &userpb.Address{
			Id:          addr["id"].(string),
			Line1:       addr["street"].(string), // Frontend uses 'street'
			City:        addr["city"].(string),
			State:       addr["state"].(string),
			PostalCode:  addr["zipCode"].(string), // Frontend uses 'zipCode'
			Country:     addr["country"].(string),
			IsDefault:   isDefault,
			AddressType: addressType,
		})
	}

	// Create update request with only fields that are available in the proto
	updateReq := &userpb.UpdateUserRequest{
		Id:           userID,
		Email:        requestPayload.Email,
		FirstName:    requestPayload.FirstName,
		LastName:     requestPayload.LastName,
		DisplayName:  requestPayload.FirstName + " " + requestPayload.LastName,
		Phone:        requestPayload.Phone,
		ProfileImage: requestPayload.ProfileImage,
		Addresses:    addresses,
		// Username field doesn't exist in the proto definition yet
	}

	// Add preferences if provided
	if len(requestPayload.Preferences) > 0 {
		// Handle preferences in a real implementation
		// This would require changes to the .proto file and service implementation
	}

	// Call user service via gRPC
	res, err := c.UserClient.UpdateUser(r.Context(), updateReq)

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "User profile updated successfully",
		Data: map[string]interface{}{
			"id": res.User.Id,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// Helper function to format user data for frontend
func formatUserForFrontend(user *userpb.User) map[string]interface{} {
	// Format addresses for frontend (camelCase and different field names)
	addresses := make([]map[string]interface{}, 0)
	for _, addr := range user.Addresses {
		addresses = append(addresses, map[string]interface{}{
			"id":        addr.Id,
			"isDefault": addr.IsDefault,
			"name":      "",         // Default name if not available
			"phone":     "",         // Default phone if not available
			"street":    addr.Line1, // Frontend uses 'street'
			"city":      addr.City,
			"state":     addr.State,
			"zipCode":   addr.PostalCode, // Frontend uses 'zipCode'
			"country":   addr.Country,
			"type":      addr.AddressType, // Frontend uses 'type'
		})
	}

	// Create user profile structure expected by frontend
	profile := map[string]interface{}{
		"firstName":   user.FirstName,
		"lastName":    user.LastName,
		"avatar":      user.ProfileImage, // Frontend uses 'avatar'
		"phone":       user.Phone,
		"dateOfBirth": "",                  // Default value if not available
		"gender":      "prefer_not_to_say", // Default value if not available
	}

	// Default preferences
	preferences := map[string]interface{}{
		"newsletter":         false,
		"marketingEmails":    false,
		"orderNotifications": true,
		"twoFactorAuth":      false,
		"language":           "en",
		"currency":           "USD",
	}

	// Default user status based on active field
	status := "active"
	if !user.Active {
		status = "inactive"
	}

	// Use username if available, otherwise use display_name
	username := user.DisplayName
	if user.Username != "" {
		username = user.Username
	}

	// Get last login or use current time as default
	lastLogin := time.Now().Format(time.RFC3339)
	if user.LastLogin != "" {
		lastLogin = user.LastLogin
	}

	// Create complete user data structure
	userData := map[string]interface{}{
		"id":             user.Id,
		"email":          user.Email,
		"username":       username,
		"profile":        profile,
		"role":           user.Role,
		"status":         status,
		"addresses":      addresses,
		"paymentMethods": []interface{}{}, // Empty array for now
		"preferences":    preferences,
		"createdAt":      user.CreatedAt,
		"updatedAt":      user.UpdatedAt,
		"lastLogin":      lastLogin,
		"wishlist":       user.Wishlist,
		"cartId":         "", // Empty for now
		"orderCount":     user.OrderCount,
		"totalSpent":     user.TotalSpend,
	}

	return userData
}

// AddToWishlist adds a product to the user's wishlist
func (c *Config) AddToWishlist(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var requestPayload struct {
		ProductID int    `json:"product_id"`
		Notes     string `json:"notes,omitempty"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Call user service via gRPC
	res, err := c.UserClient.AddToWishlist(r.Context(), &userpb.AddToWishlistRequest{
		UserId:    userID,
		ProductId: int32(requestPayload.ProductID),
		Notes:     requestPayload.Notes,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Product added to wishlist",
		Data: map[string]interface{}{
			"success": res.Success,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// RemoveFromWishlist removes a product from the user's wishlist
func (c *Config) RemoveFromWishlist(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Get product ID from URL
	productIDStr := chi.URLParam(r, "product_id")
	if productIDStr == "" {
		util.ErrorJSON(w, errors.New("product ID is required"), http.StatusBadRequest)
		return
	}

	// Convert product ID to int
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		util.ErrorJSON(w, errors.New("invalid product ID"), http.StatusBadRequest)
		return
	}

	// Call user service via gRPC
	res, err := c.UserClient.RemoveFromWishlist(r.Context(), &userpb.RemoveFromWishlistRequest{
		UserId:    userID,
		ProductId: int32(productID),
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Product removed from wishlist",
		Data: map[string]interface{}{
			"success": res.Success,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// GetWishlist retrieves the user's wishlist
func (c *Config) GetWishlist(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Call user service via gRPC
	res, err := c.UserClient.GetWishlist(r.Context(), &userpb.GetWishlistRequest{
		UserId: userID,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Convert wishlist items to response format
	items := make([]map[string]interface{}, 0)
	for _, item := range res.Items {
		// Call product service to get product details
		// This will be implemented once we have access to the product service
		productItem := map[string]interface{}{
			"id":         item.Id,
			"user_id":    item.UserId,
			"product_id": item.ProductId,
			"added_at":   item.AddedAt,
			"notes":      item.Notes,
		}
		items = append(items, productItem)
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Wishlist retrieved successfully",
		Data: map[string]interface{}{
			"items": items,
			"count": res.Count,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// formatUsersForAdminPanel formats user data for the admin panel
func (c *Config) formatUsersForAdminPanel(users []*userpb.User) []map[string]interface{} {
	customers := make([]map[string]interface{}, 0)

	for _, user := range users {
		// Get order count from checkout service
		orderCount := 0
		orderCountRes, err := c.CheckoutClient.GetUserOrderCount(context.Background(), &checkoutpb.UserRequest{
			UserId: user.Id,
		})
		if err == nil {
			orderCount = int(orderCountRes.OrderCount)
		} else {
			fmt.Printf("Error getting order count for user %s: %v\n", user.Id, err)
		}

		// Get total spend from checkout service
		totalSpend := 0.0
		totalSpendRes, err := c.CheckoutClient.GetUserTotalSpend(context.Background(), &checkoutpb.UserRequest{
			UserId: user.Id,
		})
		if err == nil {
			totalSpend = totalSpendRes.TotalSpend
		} else {
			fmt.Printf("Error getting total spend for user %s: %v\n", user.Id, err)
		}

		// Determine status based on user properties
		status := CustomerStatus.ACTIVE
		if !user.Active {
			status = CustomerStatus.INACTIVE
		}

		// Determine VIP based on role
		if user.Role == "vip" {
			status = CustomerStatus.VIP
		}

		fullName := user.FirstName + " " + user.LastName
		if fullName == " " && user.DisplayName != "" {
			fullName = user.DisplayName
		}

		// Format date
		createdAt := ""
		if t, err := time.Parse(time.RFC3339, user.CreatedAt); err == nil {
			createdAt = t.Format("2006-01-02")
		} else {
			createdAt = user.CreatedAt
		}

		// Create customer object with data needed for frontend
		customer := map[string]interface{}{
			"id":          user.Id,
			"name":        fullName,
			"email":       user.Email,
			"status":      status,
			"date":        createdAt,
			"phone":       user.Phone,
			"image":       user.ProfileImage,
			"order_count": orderCount,
			"total_spend": totalSpend,
		}

		customers = append(customers, customer)
	}

	return customers
}

// CustomerStatus enum values matching the frontend
var CustomerStatus = struct {
	ACTIVE   string
	INACTIVE string
	VIP      string
}{
	ACTIVE:   "Active",
	INACTIVE: "Inactive",
	VIP:      "VIP",
}

// ListUsers retrieves a list of users with filtering and pagination for admin panel
func (c *Config) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestPayload struct {
		Page       int     `json:"page"`
		Limit      int     `json:"limit"`
		SearchTerm *string `json:"searchTerm"`
		Status     *string `json:"status"`
		SortBy     *string `json:"sortBy"`
		SortOrder  *string `json:"sortOrder"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Default values for pagination
	if requestPayload.Page <= 0 {
		requestPayload.Page = 1
	}
	if requestPayload.Limit <= 0 {
		requestPayload.Limit = 10
	}

	// Create request for user service
	request := &userpb.GetUsersRequest{
		Page:  int32(requestPayload.Page),
		Limit: int32(requestPayload.Limit),
	}

	// Add sorting if provided
	if requestPayload.SortBy != nil && *requestPayload.SortBy != "" {
		request.SortBy = *requestPayload.SortBy
	}
	if requestPayload.SortOrder != nil && *requestPayload.SortOrder != "" {
		request.SortOrder = *requestPayload.SortOrder
	}

	// Add status filter if provided
	if requestPayload.Status != nil && *requestPayload.Status != "" && *requestPayload.Status != "All" {
		// Convert frontend status to backend status
		switch *requestPayload.Status {
		case CustomerStatus.ACTIVE:
			request.Role = "user" // Filter by regular users
			request.Active = true // Only active users
		case CustomerStatus.INACTIVE:
			request.Active = false // Only inactive users
		case CustomerStatus.VIP:
			request.Role = "vip"  // Filter by VIP users
			request.Active = true // Only active VIP users
		}
	}

	// Add search if provided
	if requestPayload.SearchTerm != nil && *requestPayload.SearchTerm != "" {
		searchReq := &userpb.SearchUsersRequest{
			Query: *requestPayload.SearchTerm,
			Page:  int32(requestPayload.Page),
			Limit: int32(requestPayload.Limit),
		}

		// Call user service via gRPC with search
		res, err := c.UserClient.SearchUsers(r.Context(), searchReq)
		if err != nil {
			util.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		// Format users for frontend
		customers := c.formatUsersForAdminPanel(res.Users)

		responseData := util.JsonResponse{
			Error:   false,
			Message: "Users retrieved successfully",
			Data: map[string]interface{}{
				"customers": customers,
				"total":     res.Total,
				"page":      res.Page,
				"limit":     res.Limit,
			},
		}

		util.WriteJSON(w, http.StatusOK, responseData)
		return
	}

	// Call user service via gRPC for regular listing
	res, err := c.UserClient.GetUsers(r.Context(), request)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Format users for frontend
	customers := c.formatUsersForAdminPanel(res.Users)

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Users retrieved successfully",
		Data: map[string]interface{}{
			"customers": customers,
			"total":     res.Total,
			"page":      res.Page,
			"limit":     res.Limit,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// GetUserStatistics retrieves statistics for the admin dashboard
func (c *Config) GetUserStatistics(w http.ResponseWriter, r *http.Request) {
	// Call the user service via gRPC to get the statistics
	res, err := c.UserClient.GetUserStatistics(r.Context(), &userpb.GetUserStatisticsRequest{})
	if err != nil {
		// Log the error
		fmt.Printf("Error calling GetUserStatistics: %v\n", err)

		// Trả về lỗi
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Sử dụng dữ liệu thực từ user service với cấu trúc phù hợp cho frontend
	statistics := map[string]interface{}{
		"total_users":      res.TotalUsers,
		"user_growth":      res.UserGrowth,
		"new_users":        res.NewUsers,
		"new_user_growth":  res.NewUserGrowth,
		"visitors":         res.Visitors,
		"visitor_growth":   res.VisitorGrowth,
		"active_users":     res.ActiveUsers,
		"repeat_customers": res.RepeatCustomers,
		"shop_visitors":    res.ShopVisitors,
		"conversion_rate":  res.ConversionRate,
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "User statistics retrieved successfully",
		Data:    statistics,
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// GetUserActivityChart retrieves data for the customer activity chart
func (c *Config) GetUserActivityChart(w http.ResponseWriter, r *http.Request) {
	// Parse request to get parameters
	var requestPayload struct {
		Days      int    `json:"days"`
		ChartType string `json:"chart_type,omitempty"` // Optional chart type
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		// Default to 7 days if not specified
		requestPayload.Days = 7
	}

	// Validate and set defaults
	if requestPayload.Days <= 0 {
		requestPayload.Days = 7
	}

	// Call the user service via gRPC to get the chart data
	res, err := c.UserClient.GetUserActivityChart(r.Context(), &userpb.GetUserActivityChartRequest{
		Days: int32(requestPayload.Days),
	})

	if err != nil {
		// Log the error
		fmt.Printf("Error calling GetUserActivityChart: %v\n", err)

		// Trả về lỗi
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Convert the proto response to a slice of maps for JSON output
	chartData := make([]map[string]interface{}, 0, len(res.ChartData))

	// Calculate statistics based on actual chart data
	var totalCount int
	var maxCount int
	var minCount = -1

	for _, point := range res.ChartData {
		dataPoint := map[string]interface{}{
			"day":   point.Day,
			"date":  point.Date,
			"count": point.Count,
		}

		// Add a mocked secondary value based on the chart type
		// This is synthetic data created on the fly for the frontend
		switch requestPayload.ChartType {
		case "revenue":
			// For revenue, add average order value
			avgOrderValue := float64(30 + (point.Count % 70))
			dataPoint["value"] = float64(point.Count) * avgOrderValue
		case "conversion":
			// For conversion, calculate conversion rate
			totalVisitors := point.Count * (5 + (int32(point.Day[0]) % 10))
			conversionRate := float64(point.Count) * 100 / float64(totalVisitors)
			dataPoint["value"] = conversionRate
		case "orders":
			// For orders, add average order value
			dataPoint["value"] = 50.0 + float64(point.Count%50)
		default:
			// For activity and other types
			sessionTime := 15 + (point.Count % 30)
			dataPoint["value"] = float64(sessionTime)
		}

		chartData = append(chartData, dataPoint)

		// Update statistics
		totalCount += int(point.Count)
		if int(point.Count) > maxCount {
			maxCount = int(point.Count)
		}
		if minCount == -1 || int(point.Count) < minCount {
			minCount = int(point.Count)
		}
	}

	// Calculate average count
	var avgCount float64
	if len(res.ChartData) > 0 {
		avgCount = float64(totalCount) / float64(len(res.ChartData))
	}

	// Calculate growth (comparing first half to second half)
	var growthRate float64
	if len(res.ChartData) > 1 {
		midpoint := len(res.ChartData) / 2
		var firstHalfCount, secondHalfCount int

		for i, point := range res.ChartData {
			if i < midpoint {
				firstHalfCount += int(point.Count)
			} else {
				secondHalfCount += int(point.Count)
			}
		}

		if firstHalfCount > 0 {
			growthRate = float64(secondHalfCount-firstHalfCount) * 100.0 / float64(firstHalfCount)
		}
	}

	// Build a rich chart response for the frontend
	chartType := requestPayload.ChartType
	if chartType == "" {
		chartType = "activity" // Default chart type
	}

	// Generate appropriate title and labels based on chart type
	title := "User Activity"
	yAxisLabel := "Count"
	description := "Daily user activity over the selected period"

	switch chartType {
	case "activity":
		title = "User Activity"
		yAxisLabel = "Active Users"
		description = "Daily active users over the selected period"
	case "new_users":
		title = "New User Registrations"
		yAxisLabel = "New Users"
		description = "New user sign-ups over the selected period"
	case "orders":
		title = "Order Volume"
		yAxisLabel = "Orders"
		description = "Number of orders placed over the selected period"
	case "revenue":
		title = "Revenue"
		yAxisLabel = "Revenue ($)"
		description = "Daily revenue over the selected period"
	case "conversion":
		title = "Conversion Rate"
		yAxisLabel = "Conversion (%)"
		description = "Daily visitor-to-customer conversion rate"
	}

	chartResponse := map[string]interface{}{
		"chart_data":   chartData,
		"chart_type":   chartType,
		"title":        title,
		"y_axis_label": yAxisLabel,
		"description":  description,
		"total_count":  totalCount,
		"avg_count":    avgCount,
		"max_count":    maxCount,
		"min_count":    minCount,
		"growth_rate":  growthRate,
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "User activity chart data retrieved successfully",
		Data:    chartResponse,
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// UpdateCustomerStatus updates a customer's status (active, inactive, vip)
func (c *Config) UpdateCustomerStatus(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL
	userID := chi.URLParam(r, "id")
	if userID == "" {
		util.ErrorJSON(w, errors.New("user ID is required"), http.StatusBadRequest)
		return
	}

	// Get new status from request body
	var requestPayload struct {
		Status string `json:"status"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate status
	if requestPayload.Status != CustomerStatus.ACTIVE &&
		requestPayload.Status != CustomerStatus.INACTIVE &&
		requestPayload.Status != CustomerStatus.VIP {
		util.ErrorJSON(w, errors.New("invalid status value"), http.StatusBadRequest)
		return
	}

	// First get current user info
	userRes, err := c.UserClient.GetUser(r.Context(), &userpb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Prepare update data
	user := userRes.User
	updateReq := &userpb.UpdateUserRequest{
		Id:          userID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DisplayName: user.DisplayName,
		Phone:       user.Phone,
		Role:        user.Role,
	}

	// Update status based on new value
	switch requestPayload.Status {
	case CustomerStatus.ACTIVE:
		updateReq.Active = true
		updateReq.Role = "user" // Ensure not VIP
	case CustomerStatus.INACTIVE:
		updateReq.Active = false
	case CustomerStatus.VIP:
		updateReq.Active = true
		updateReq.Role = "vip"
	}

	// Call user service to update
	_, err = c.UserClient.UpdateUser(r.Context(), updateReq)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Customer status updated successfully",
		Data: map[string]interface{}{
			"id":     userID,
			"status": requestPayload.Status,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// DeleteCustomer deletes a customer
func (c *Config) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL
	userID := chi.URLParam(r, "id")
	if userID == "" {
		util.ErrorJSON(w, errors.New("user ID is required"), http.StatusBadRequest)
		return
	}

	// Call user service to delete user
	res, err := c.UserClient.DeleteUser(r.Context(), &userpb.DeleteUserRequest{
		Id: userID,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Customer deleted successfully",
		Data: map[string]interface{}{
			"success": res.Success,
			"id":      userID,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

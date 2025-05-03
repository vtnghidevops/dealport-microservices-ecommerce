package user

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"broker-service/internal/util"
	userpb "broker-service/proto/user"

	"github.com/go-chi/chi/v5"
)

type Config struct {
	UserClient userpb.UserServiceClient
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
		"lastLogin":      time.Now().Format(time.RFC3339), // Default value
		"wishlist":       user.Wishlist,                   // Use wishlist from user object
		"cartId":         "",                              // Empty for now
		"orderCount":     0,                               // Default value
		"totalSpent":     0,                               // Default value
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

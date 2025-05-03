package cart

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"broker-service/internal/util"
	cartpb "broker-service/proto/cart"
	couponpb "broker-service/proto/coupon"

	"github.com/go-chi/chi/v5"
)

type Config struct {
	CartClient   cartpb.CartServiceClient
	CouponClient couponpb.CouponServiceClient
}

// GetCart retrieves a user's cart
func (c *Config) GetCart(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Debug log for userID
	fmt.Printf("GetCart handler: User ID from context: '%s'\n", userID)

	// Call cart service via gRPC
	res, err := c.CartClient.GetCart(r.Context(), &cartpb.GetCartRequest{
		UserId: userID,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Cart retrieved successfully",
		Data:    formatCartResponse(res),
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// AddCartItem adds an item to the user's cart
func (c *Config) AddCartItem(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Debug log for userID
	fmt.Printf("AddCartItem handler: User ID from context: '%s'\n", userID)

	var requestPayload struct {
		ProductID     int64   `json:"productId"`
		Name          string  `json:"name"`
		Price         float64 `json:"price"`
		OriginalPrice float64 `json:"originalPrice,omitempty"`
		Quantity      int32   `json:"quantity"`
		ImageURL      string  `json:"imageUrl,omitempty"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Call cart service via gRPC
	res, err := c.CartClient.AddCartItem(r.Context(), &cartpb.AddCartItemRequest{
		UserId: userID,
		Item: &cartpb.CartItem{
			ProductId:     strconv.FormatInt(requestPayload.ProductID, 10),
			Name:          requestPayload.Name,
			Price:         requestPayload.Price,
			OriginalPrice: requestPayload.OriginalPrice,
			Quantity:      requestPayload.Quantity,
			ImageUrl:      requestPayload.ImageURL,
		},
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Debug the response from cart service
	fmt.Printf("AddCartItem: Response cart ID: %s, User ID: %s, Items count: %d\n",
		res.Id, res.UserId, len(res.Items))

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Item added to cart",
		Data:    formatCartResponse(res),
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// UpdateCartItem updates an item in the user's cart
func (c *Config) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Get item ID from URL
	itemID := chi.URLParam(r, "item_id")
	if itemID == "" {
		util.ErrorJSON(w, errors.New("item ID is required"), http.StatusBadRequest)
		return
	}

	var requestPayload struct {
		Quantity int32 `json:"quantity"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Call cart service via gRPC
	res, err := c.CartClient.UpdateCartItem(r.Context(), &cartpb.UpdateCartItemRequest{
		UserId:   userID,
		ItemId:   itemID,
		Quantity: requestPayload.Quantity,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Cart item updated",
		Data:    formatCartResponse(res),
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// RemoveCartItem removes an item from the user's cart
func (c *Config) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Get item ID from URL
	itemID := chi.URLParam(r, "item_id")
	if itemID == "" {
		util.ErrorJSON(w, errors.New("item ID is required"), http.StatusBadRequest)
		return
	}

	// Call cart service via gRPC
	res, err := c.CartClient.RemoveCartItem(r.Context(), &cartpb.RemoveCartItemRequest{
		UserId: userID,
		ItemId: itemID,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Item removed from cart",
		Data:    formatCartResponse(res),
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// ClearCart removes all items from the user's cart
func (c *Config) ClearCart(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Call cart service via gRPC
	res, err := c.CartClient.ClearCart(r.Context(), &cartpb.ClearCartRequest{
		UserId: userID,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Cart cleared",
		Data: map[string]interface{}{
			"success": res.Success,
			"message": res.Message,
		},
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// ApplyCoupon applies a coupon to the user's cart
func (c *Config) ApplyCoupon(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var requestPayload struct {
		CouponCode string `json:"couponCode"`
	}

	err := util.ReadJSON(w, r, &requestPayload)
	if err != nil {
		util.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Call cart service via gRPC
	res, err := c.CartClient.ApplyCoupon(r.Context(), &cartpb.ApplyCouponRequest{
		UserId:     userID,
		CouponCode: requestPayload.CouponCode,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Coupon applied",
		Data:    formatCartResponse(res),
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// RemoveCoupon removes a coupon from the user's cart
func (c *Config) RemoveCoupon(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context - would be set by auth middleware
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Call cart service via gRPC
	res, err := c.CartClient.RemoveCoupon(r.Context(), &cartpb.RemoveCouponRequest{
		UserId: userID,
	})

	if err != nil {
		util.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	responseData := util.JsonResponse{
		Error:   false,
		Message: "Coupon removed",
		Data:    formatCartResponse(res),
	}

	util.WriteJSON(w, http.StatusOK, responseData)
}

// Helper function to format cart response for the frontend
func formatCartResponse(cart *cartpb.Cart) map[string]interface{} {
	// Format cart items
	items := make([]map[string]interface{}, 0)
	for _, item := range cart.Items {
		// Convert ProductId string to number if possible
		productId, err := strconv.ParseInt(item.ProductId, 10, 64)
		if err != nil {
			productId = 0 // Default value if conversion fails
		}

		cartItem := map[string]interface{}{
			"id":            item.Id,
			"productId":     productId,
			"name":          item.Name,
			"price":         item.Price,
			"originalPrice": item.OriginalPrice,
			"quantity":      item.Quantity,
			"imageUrl":      item.ImageUrl,
		}
		items = append(items, cartItem)
	}

	// Format cart totals
	totals := map[string]interface{}{
		"subtotal": cart.Totals.Subtotal,
		"shipping": cart.Totals.Shipping,
		"discount": cart.Totals.Discount,
		"tax":      cart.Totals.Tax,
		"total":    cart.Totals.Total,
	}

	// Create complete cart data structure with camelCase keys for frontend
	cartData := map[string]interface{}{
		"id":             cart.Id,
		"userId":         cart.UserId,
		"items":          items,
		"totals":         totals,
		"couponCode":     cart.CouponCode,
		"discountAmount": cart.DiscountAmount,
		"itemCount":      len(cart.Items),
		"createdAt":      cart.CreatedAt,
		"updatedAt":      cart.UpdatedAt,
	}

	return cartData
}

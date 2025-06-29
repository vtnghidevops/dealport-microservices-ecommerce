package cart

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	couponpb "broker-service/proto/coupon"

	"github.com/go-chi/chi/v5"
)

// NewCouponHandler creates a new coupon handler
func NewCouponHandler(couponClient couponpb.CouponServiceClient) *Config {
	return &Config{
		CouponClient: couponClient,
	}
}

// CouponResponse represents the response for coupon operations
type CouponResponse struct {
	ID             string    `json:"id"`
	Code           string    `json:"code"`
	Discount       float64   `json:"discount"`
	DiscountType   string    `json:"discountType"`
	MinOrderAmount float64   `json:"minOrderAmount"`
	MaxUsage       int       `json:"maxUsage"`
	UsageCount     int       `json:"usageCount"`
	ValidFrom      time.Time `json:"validFrom"`
	ValidTo        time.Time `json:"validTo"`
	IsActive       bool      `json:"isActive"`
	Description    string    `json:"description,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// convertProtoCouponToResponse converts a proto coupon message to a CouponResponse
func convertProtoCouponToResponse(coupon *couponpb.Coupon) CouponResponse {
	return CouponResponse{
		ID:             coupon.Id,
		Code:           coupon.Code,
		Discount:       coupon.Discount,
		DiscountType:   coupon.DiscountType,
		MinOrderAmount: coupon.MinOrderAmount,
		MaxUsage:       int(coupon.MaxUsage),
		UsageCount:     int(coupon.UsageCount),
		ValidFrom:      coupon.ValidFrom.AsTime(),
		ValidTo:        coupon.ValidTo.AsTime(),
		IsActive:       coupon.IsActive,
		Description:    coupon.Description,
		CreatedAt:      coupon.CreatedAt.AsTime(),
		UpdatedAt:      coupon.UpdatedAt.AsTime(),
	}
}

// GetCoupons handles GET /coupons
func (c *Config) GetCoupons(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Call the coupon service
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create request
	req := &couponpb.GetCouponsRequest{
		Page:  int32(page),
		Limit: int32(limit),
	}

	// Call gRPC service
	resp, err := c.CouponClient.GetCoupons(ctx, req)
	if err != nil {
		log.Printf("Error getting coupons: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to get coupons")
		return
	}

	// Convert response
	coupons := make([]CouponResponse, 0, len(resp.Coupons))
	for _, coupon := range resp.Coupons {
		coupons = append(coupons, convertProtoCouponToResponse(coupon))
	}

	sendJSONResponse(w, http.StatusOK, envelope{
		"data": coupons,
		"meta": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": resp.Total,
		},
	})
}

// GetCouponByID handles GET /coupons/{id}
func (c *Config) GetCouponByID(w http.ResponseWriter, r *http.Request) {
	// Get coupon ID from URL
	id := chi.URLParam(r, "id")
	if id == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Coupon ID is required")
		return
	}

	// Call the coupon service
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create request
	req := &couponpb.GetCouponByIDRequest{
		Id: id,
	}

	// Call gRPC service
	resp, err := c.CouponClient.GetCouponByID(ctx, req)
	if err != nil {
		log.Printf("Error getting coupon by ID: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to get coupon")
		return
	}

	if resp.Coupon == nil {
		sendErrorResponse(w, http.StatusNotFound, "Coupon not found")
		return
	}

	// Convert response
	coupon := convertProtoCouponToResponse(resp.Coupon)

	sendJSONResponse(w, http.StatusOK, envelope{"data": coupon})
}

// GetCouponByCode handles GET /coupons/code/{code}
func (c *Config) GetCouponByCode(w http.ResponseWriter, r *http.Request) {
	// Get coupon code from URL
	code := chi.URLParam(r, "code")
	if code == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Coupon code is required")
		return
	}

	// Call the coupon service
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create request
	req := &couponpb.GetCouponByCodeRequest{
		Code: code,
	}

	// Call gRPC service
	resp, err := c.CouponClient.GetCouponByCode(ctx, req)
	if err != nil {
		log.Printf("Error getting coupon by code: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to get coupon")
		return
	}

	if resp.Coupon == nil {
		sendErrorResponse(w, http.StatusNotFound, "Coupon not found")
		return
	}

	// Convert response
	coupon := convertProtoCouponToResponse(resp.Coupon)

	sendJSONResponse(w, http.StatusOK, envelope{"data": coupon})
}

// CreateCouponRequest represents the request to create a coupon
type CreateCouponRequest struct {
	Code           string  `json:"code"`
	Discount       float64 `json:"discount"`
	DiscountType   string  `json:"discountType"`
	MinOrderAmount float64 `json:"minOrderAmount"`
	MaxUsage       int     `json:"maxUsage"`
	ValidFrom      string  `json:"validFrom,omitempty"`
	ValidTo        string  `json:"validTo,omitempty"`
	IsActive       bool    `json:"isActive"`
	Description    string  `json:"description,omitempty"`
}

// CreateCoupon handles POST /coupons
func (c *Config) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	var req CreateCouponRequest

	// Parse request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}
	fmt.Println("Full Conpon receive in broker: ", req)
	// Validate required fields
	if req.Code == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Coupon code is required")
		return
	}

	if req.Discount <= 0 {
		sendErrorResponse(w, http.StatusBadRequest, "Discount must be greater than 0")
		return
	}

	if req.DiscountType != "percentage" && req.DiscountType != "fixed" {
		sendErrorResponse(w, http.StatusBadRequest, "Discount type must be 'percentage' or 'fixed'")
		return
	}

	// Call the coupon service
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create request
	grpcReq := &couponpb.CreateCouponRequest{
		Code:           req.Code,
		Discount:       req.Discount,
		DiscountType:   req.DiscountType,
		MinOrderAmount: req.MinOrderAmount,
		MaxUsage:       int32(req.MaxUsage),
		ValidFrom:      req.ValidFrom,
		ValidTo:        req.ValidTo,
		IsActive:       req.IsActive,
		Description:    req.Description,
	}

	// Call gRPC service
	resp, err := c.CouponClient.CreateCoupon(ctx, grpcReq)
	if err != nil {
		log.Printf("Error creating coupon: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to create coupon")
		return
	}

	if resp.Coupon == nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to create coupon")
		return
	}

	// Convert response
	coupon := convertProtoCouponToResponse(resp.Coupon)

	sendJSONResponse(w, http.StatusCreated, envelope{"data": coupon})
}

// UpdateCoupon handles PUT /coupons/{id}
func (c *Config) UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	// Get coupon ID from URL
	id := chi.URLParam(r, "id")
	if id == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Coupon ID is required")
		return
	}

	var req CreateCouponRequest

	// Parse request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Call the coupon service
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create request
	grpcReq := &couponpb.UpdateCouponRequest{
		Id:             id,
		Code:           req.Code,
		Discount:       req.Discount,
		DiscountType:   req.DiscountType,
		MinOrderAmount: req.MinOrderAmount,
		MaxUsage:       int32(req.MaxUsage),
		ValidFrom:      req.ValidFrom,
		ValidTo:        req.ValidTo,
		IsActive:       req.IsActive,
		Description:    req.Description,
	}

	// Call gRPC service
	resp, err := c.CouponClient.UpdateCoupon(ctx, grpcReq)
	if err != nil {
		log.Printf("Error updating coupon: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to update coupon")
		return
	}

	if resp.Coupon == nil {
		sendErrorResponse(w, http.StatusNotFound, "Coupon not found")
		return
	}

	// Convert response
	coupon := convertProtoCouponToResponse(resp.Coupon)

	sendJSONResponse(w, http.StatusOK, envelope{"data": coupon})
}

// DeleteCoupon handles DELETE /coupons/{id}
func (c *Config) DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	// Get coupon ID from URL
	id := chi.URLParam(r, "id")
	if id == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Coupon ID is required")
		return
	}

	// Call the coupon service
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create request
	req := &couponpb.DeleteCouponRequest{
		Id: id,
	}

	// Call gRPC service
	resp, err := c.CouponClient.DeleteCoupon(ctx, req)
	if err != nil {
		log.Printf("Error deleting coupon: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete coupon")
		return
	}

	sendJSONResponse(w, http.StatusOK, envelope{
		"success": resp.Success,
		"message": resp.Message,
	})
}

// Helper functions
type envelope map[string]interface{}

func sendJSONResponse(w http.ResponseWriter, status int, data envelope) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error sending JSON response: %v", err)
	}
}

func sendErrorResponse(w http.ResponseWriter, status int, message string) {
	sendJSONResponse(w, status, envelope{"error": message})
}

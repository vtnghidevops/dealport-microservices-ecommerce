package grpc

import (
	pb "broker-service/proto/product"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// BannerGrpcHandler handles banner-related HTTP requests using gRPC
type BannerGrpcHandler struct {
	client *ProductClient
}

// NewBannerGrpcHandler creates a new banner gRPC handler
func NewBannerGrpcHandler(client *ProductClient) *BannerGrpcHandler {
	return &BannerGrpcHandler{
		client: client,
	}
}

// GetAllBanners handles GET /api/v1/banners
func (h *BannerGrpcHandler) GetAllBanners(w http.ResponseWriter, r *http.Request) {
	// Extract page and page size from query parameters
	page := 1
	pageSize := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	// Extract filter parameters from the query string
	filters := make(map[string]string)
	for key, values := range r.URL.Query() {
		if key != "page" && key != "page_size" && len(values) > 0 {
			filters[key] = values[0]
		}
	}

	// Call the gRPC client
	resp, err := h.client.ListBanners(page, pageSize, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Banners retrieved successfully",
		"data":    resp.Banners,
		"total":   resp.Total,
	})
}

// GetBannerByID handles GET /api/v1/banners/{id}
func (h *BannerGrpcHandler) GetBannerByID(w http.ResponseWriter, r *http.Request) {
	// Get the banner ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid banner ID", http.StatusBadRequest)
		return
	}

	// Call the gRPC client
	banner, err := h.client.GetBannerByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Banner retrieved successfully",
		"data":    banner,
	})
}

// GetBannersByType handles GET /api/v1/banners/type/{type}
func (h *BannerGrpcHandler) GetBannersByType(w http.ResponseWriter, r *http.Request) {
	// Get the banner type from the URL
	bannerType := chi.URLParam(r, "type")
	if bannerType == "" {
		http.Error(w, "Invalid banner type", http.StatusBadRequest)
		return
	}

	// Get the limit parameter
	limit := 5
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Call the gRPC client
	resp, err := h.client.GetBannersByType(bannerType, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Banners retrieved successfully",
		"data":    resp.Banners,
	})
}

// CreateBanner handles POST /api/v1/banners
func (h *BannerGrpcHandler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var bannerReq BannerRequest
	err := json.NewDecoder(r.Body).Decode(&bannerReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if bannerReq.Title == "" || bannerReq.ImageURL == "" {
		http.Error(w, "Title and imageUrl are required", http.StatusBadRequest)
		return
	}

	// Create the gRPC request
	req := &pb.Banner{
		Title:           bannerReq.Title,
		Subtitle:        bannerReq.Subtitle,
		Description:     bannerReq.Description,
		Discount:        bannerReq.Discount,
		HighlightText:   bannerReq.HighlightText,
		ImageUrl:        bannerReq.ImageURL,
		LinkUrl:         bannerReq.LinkURL,
		ActionText:      bannerReq.ActionText,
		BackgroundColor: bannerReq.BackgroundColor,
		TextColor:       bannerReq.TextColor,
		AnimationType:   bannerReq.AnimationType,
		IsActive:        bannerReq.IsActive,
		Priority:        int32(bannerReq.Priority),
		Type:            bannerReq.Type,
	}

	// Handle optional product and category IDs
	if bannerReq.ProductID > 0 {
		req.ProductId = int32(bannerReq.ProductID)
	}
	if bannerReq.CategoryID > 0 {
		req.CategoryId = int32(bannerReq.CategoryID)
	}

	// Call the gRPC client
	resp, err := h.client.CreateBanner(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Banner created successfully",
		"data":    resp.Banner,
	})
}

// UpdateBanner handles PUT /api/v1/banners/{id}
func (h *BannerGrpcHandler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	// Get the banner ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid banner ID", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var bannerReq BannerRequest
	err = json.NewDecoder(r.Body).Decode(&bannerReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if bannerReq.Title == "" || bannerReq.ImageURL == "" {
		http.Error(w, "Title and imageUrl are required", http.StatusBadRequest)
		return
	}

	// Create the gRPC request
	req := &pb.Banner{
		Id:              int32(id),
		Title:           bannerReq.Title,
		Subtitle:        bannerReq.Subtitle,
		Description:     bannerReq.Description,
		Discount:        bannerReq.Discount,
		HighlightText:   bannerReq.HighlightText,
		ImageUrl:        bannerReq.ImageURL,
		LinkUrl:         bannerReq.LinkURL,
		ActionText:      bannerReq.ActionText,
		BackgroundColor: bannerReq.BackgroundColor,
		TextColor:       bannerReq.TextColor,
		AnimationType:   bannerReq.AnimationType,
		IsActive:        bannerReq.IsActive,
		Priority:        int32(bannerReq.Priority),
		Type:            bannerReq.Type,
	}

	// Handle optional product and category IDs
	if bannerReq.ProductID > 0 {
		req.ProductId = int32(bannerReq.ProductID)
	}
	if bannerReq.CategoryID > 0 {
		req.CategoryId = int32(bannerReq.CategoryID)
	}

	// Call the gRPC client
	resp, err := h.client.UpdateBanner(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": resp.Message,
	})
}

// DeleteBanner handles DELETE /api/v1/banners/{id}
func (h *BannerGrpcHandler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	// Get the banner ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid banner ID", http.StatusBadRequest)
		return
	}

	// Call the gRPC client
	resp, err := h.client.DeleteBanner(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": resp.Message,
	})
}

// BannerRequest represents the JSON structure for banner requests
type BannerRequest struct {
	Title           string `json:"title"`
	Subtitle        string `json:"subtitle,omitempty"`
	Description     string `json:"description,omitempty"`
	Discount        string `json:"discount,omitempty"`
	HighlightText   string `json:"highlightText,omitempty"`
	ImageURL        string `json:"imageUrl"`
	LinkURL         string `json:"linkUrl,omitempty"`
	ActionText      string `json:"actionText,omitempty"`
	BackgroundColor string `json:"backgroundColor,omitempty"`
	TextColor       string `json:"textColor,omitempty"`
	AnimationType   string `json:"animationType,omitempty"`
	IsActive        bool   `json:"isActive"`
	Priority        int    `json:"priority,omitempty"`
	Type            string `json:"type"` // hero, promotional, category, seasonal, product
	ProductID       int    `json:"productId,omitempty"`
	CategoryID      int    `json:"categoryId,omitempty"`
}

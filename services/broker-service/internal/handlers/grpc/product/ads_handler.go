package grpc

import (
	pb "broker-service/proto/product"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// AdsGrpcHandler handles ads-related HTTP requests using gRPC
type AdsGrpcHandler struct {
	client *ProductClient
}

// NewAdsGrpcHandler creates a new ads gRPC handler
func NewAdsGrpcHandler(client *ProductClient) *AdsGrpcHandler {
	return &AdsGrpcHandler{
		client: client,
	}
}

// GetAdsByLocation handles GET /api/v1/ads/placement/{location}
func (h *AdsGrpcHandler) GetAdsByLocation(w http.ResponseWriter, r *http.Request) {
	// Get the location from the URL
	location := chi.URLParam(r, "location")
	if location == "" {
		http.Error(w, "Invalid location", http.StatusBadRequest)
		return
	}

	// Call the gRPC client
	resp, err := h.client.GetAdsByLocation(location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Ads retrieved successfully",
		"data":    resp.Ads,
	})
}

// GetAllAdsPlacements handles GET /api/v1/ads/placements
func (h *AdsGrpcHandler) GetAllAdsPlacements(w http.ResponseWriter, r *http.Request) {
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

	// Call the gRPC client
	resp, err := h.client.ListAdsPlacements(page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Ads placements retrieved successfully",
		"data":    resp.Placements,
		"total":   resp.Total,
	})
}

// GetAdsPlacementByID handles GET /api/v1/ads/placements/{id}
func (h *AdsGrpcHandler) GetAdsPlacementByID(w http.ResponseWriter, r *http.Request) {
	// Get the placement ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid placement ID", http.StatusBadRequest)
		return
	}

	// Call the gRPC client
	placement, err := h.client.GetAdsPlacementByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Ads placement retrieved successfully",
		"data":    placement,
	})
}

// CreateAdsPlacement handles POST /api/v1/ads/placements
func (h *AdsGrpcHandler) CreateAdsPlacement(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var placementReq AdsPlacementRequest
	err := json.NewDecoder(r.Body).Decode(&placementReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if placementReq.Location == "" || placementReq.ReferenceType == "" || placementReq.ReferenceID <= 0 {
		http.Error(w, "Location, referenceType, and referenceId are required", http.StatusBadRequest)
		return
	}

	// Create the gRPC request
	req := &pb.AdsPlacement{
		Location:       placementReq.Location,
		ReferenceType:  placementReq.ReferenceType,
		ReferenceId:    int32(placementReq.ReferenceID),
		DisplayOrder:   int32(placementReq.DisplayOrder),
		CustomTitle:    placementReq.CustomTitle,
		CustomImageUrl: placementReq.CustomImageURL,
		UiSettings:     placementReq.UISettings,
		IsActive:       placementReq.IsActive,
	}

	// Call the gRPC client
	resp, err := h.client.CreateAdsPlacement(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Ads placement created successfully",
		"data":    resp.Placement,
	})
}

// UpdateAdsPlacement handles PUT /api/v1/ads/placements/{id}
func (h *AdsGrpcHandler) UpdateAdsPlacement(w http.ResponseWriter, r *http.Request) {
	// Get the placement ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid placement ID", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var placementReq AdsPlacementRequest
	err = json.NewDecoder(r.Body).Decode(&placementReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if placementReq.Location == "" || placementReq.ReferenceType == "" || placementReq.ReferenceID <= 0 {
		http.Error(w, "Location, referenceType, and referenceId are required", http.StatusBadRequest)
		return
	}

	// Create the gRPC request
	req := &pb.AdsPlacement{
		Id:             int32(id),
		Location:       placementReq.Location,
		ReferenceType:  placementReq.ReferenceType,
		ReferenceId:    int32(placementReq.ReferenceID),
		DisplayOrder:   int32(placementReq.DisplayOrder),
		CustomTitle:    placementReq.CustomTitle,
		CustomImageUrl: placementReq.CustomImageURL,
		UiSettings:     placementReq.UISettings,
		IsActive:       placementReq.IsActive,
	}

	// Call the gRPC client
	resp, err := h.client.UpdateAdsPlacement(req)
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

// DeleteAdsPlacement handles DELETE /api/v1/ads/placements/{id}
func (h *AdsGrpcHandler) DeleteAdsPlacement(w http.ResponseWriter, r *http.Request) {
	// Get the placement ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid placement ID", http.StatusBadRequest)
		return
	}

	// Call the gRPC client
	resp, err := h.client.DeleteAdsPlacement(id)
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

// AdsPlacementRequest represents the JSON structure for ads placement requests
type AdsPlacementRequest struct {
	Location       string `json:"location"`
	ReferenceType  string `json:"referenceType"`
	ReferenceID    int    `json:"referenceId"`
	DisplayOrder   int    `json:"displayOrder"`
	CustomTitle    string `json:"customTitle,omitempty"`
	CustomImageURL string `json:"customImageUrl,omitempty"`
	UISettings     string `json:"uiSettings,omitempty"` // JSON as string
	IsActive       bool   `json:"isActive"`
}

package grpc

import (
	pb "broker-service/proto/product"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// CategoryGrpcHandler handles category-related HTTP requests using gRPC
type CategoryGrpcHandler struct {
	client *ProductClient
}

// NewCategoryGrpcHandler creates a new category gRPC handler
func NewCategoryGrpcHandler(client *ProductClient) *CategoryGrpcHandler {
	return &CategoryGrpcHandler{
		client: client,
	}
}

// GetAllCategories handles GET /api/v1/categories
func (h *CategoryGrpcHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	// Extract filter parameters from the query string
	filters := make(map[string]string)
	for key, values := range r.URL.Query() {
		if key != "page" && key != "page_size" && len(values) > 0 {
			filters[key] = values[0]
		}
	}

	// Call the gRPC client
	resp, err := h.client.ListCategories(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": http.StatusOK,
		"data":   resp.Categories,
	})
}

// GetCategoryByID handles GET /api/v1/categories/{id}
func (h *CategoryGrpcHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	// Get the category ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Call the gRPC client
	category, err := h.client.GetCategoryByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": http.StatusOK,
		"data":   category,
	})
}

// GetCategoryBySlug handles GET /api/v1/categories/slug/{slug}
func (h *CategoryGrpcHandler) GetCategoryBySlug(w http.ResponseWriter, r *http.Request) {
	// Get the category slug from the URL
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Invalid category slug", http.StatusBadRequest)
		return
	}

	// Call the gRPC client
	category, err := h.client.GetCategoryBySlug(slug)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": http.StatusOK,
		"data":   category,
	})
}

// CreateCategory handles POST /api/v1/categories
func (h *CategoryGrpcHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var category CategoryRequest
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if category.Name == "" || category.Slug == "" {
		http.Error(w, "Name and slug are required", http.StatusBadRequest)
		return
	}

	// Create the gRPC request
	req := &pb.Category{
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ImageUrl:    category.ImageURL,
		IsActive:    category.IsActive,
		IsVisible:   category.IsVisible,
	}

	// Call the gRPC client
	resp, err := h.client.CreateCategory(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  http.StatusCreated,
		"message": "Category created successfully",
		"data":    resp.Category,
	})
}

// UpdateCategory handles PUT /api/v1/categories/{id}
func (h *CategoryGrpcHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	// Get the category ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var category CategoryRequest
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if category.Name == "" || category.Slug == "" {
		http.Error(w, "Name and slug are required", http.StatusBadRequest)
		return
	}

	// Create the gRPC request
	req := &pb.Category{
		Id:          int32(id),
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ImageUrl:    category.ImageURL,
		IsActive:    category.IsActive,
		IsVisible:   category.IsVisible,
	}

	// Call the gRPC client
	_, err = h.client.UpdateCategory(req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated category to include in response
	updatedCategory, err := h.client.GetCategoryByID(id)
	if err != nil {
		// Already updated, just return success without the data
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusOK,
			"message": "Category updated successfully",
		})
		return
	}

	// Return the response with updated category
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Category updated successfully",
		"data":    updatedCategory,
	})
}

// DeleteCategory handles DELETE /api/v1/categories/{id}
func (h *CategoryGrpcHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Get the category ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Call the gRPC client
	_, err = h.client.DeleteCategory(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Category deleted successfully",
	})
}

// CategoryRequest represents the JSON structure for category requests
type CategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"imageUrl,omitempty"`
	IsActive    bool   `json:"isActive"`
	IsVisible   bool   `json:"isVisible"`
}

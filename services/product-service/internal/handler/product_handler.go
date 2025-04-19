// cmd/api/handlers.go
package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"product-service/data"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Config is the application config
type Config struct {
	Database *sql.DB
	Models   data.Models
}

// Response is the JSON response format
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

// Helper methods for JSON response handling
func (app *Config) writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := Response{
		Status:  statusCode,
		Message: err.Error(),
	}

	return app.writeJSON(w, statusCode, payload)
}

// HealthCheck is a simple health check endpoint
func (app *Config) HealthCheck(w http.ResponseWriter, r *http.Request) {
	payload := Response{
		Status:  http.StatusOK,
		Message: "Product service is healthy and running",
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// GetAllProducts returns all products with filtering and pagination
func (app *Config) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if pageSize <= 0 {
		pageSize = 10
	}

	// Build filters from query parameters
	filters := make(map[string]string)

	if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
		filters["category_id"] = categoryID
	}

	if categorySlug := r.URL.Query().Get("category_slug"); categorySlug != "" {
		filters["category_slug"] = categorySlug
	}

	if brand := r.URL.Query().Get("brand"); brand != "" {
		filters["brand"] = brand
	}

	if productType := r.URL.Query().Get("type"); productType != "" {
		filters["type"] = productType
	}

	if orderBy := r.URL.Query().Get("order_by"); orderBy != "" {
		filters["order_by"] = orderBy
	}

	if orderDir := r.URL.Query().Get("order_dir"); orderDir != "" {
		filters["order_dir"] = orderDir
	}

	// Get products from database
	products, totalCount, err := app.Models.Product.GetAll(page, pageSize, filters)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := (totalCount + pageSize - 1) / pageSize

	// Build response
	meta := PaginationMeta{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalItems:  totalCount,
		TotalPages:  totalPages,
	}

	payload := Response{
		Status: http.StatusOK,
		Data:   products,
		Meta:   meta,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// GetProductByID returns a product by ID
func (app *Config) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		app.errorJSON(w, errors.New("missing product ID"), http.StatusBadRequest)
		return
	}

	product, err := app.Models.Product.GetByID(id)
	if err != nil {
		if err.Error() == "product not found" {
			app.errorJSON(w, err, http.StatusNotFound)
			return
		}
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := Response{
		Status: http.StatusOK,
		Data:   product,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// GetProductBySlug returns a product by slug
func (app *Config) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		app.errorJSON(w, errors.New("missing product slug"), http.StatusBadRequest)
		return
	}

	product, err := app.Models.Product.GetBySlug(slug)
	if err != nil {
		if err.Error() == "product not found" {
			app.errorJSON(w, err, http.StatusNotFound)
			return
		}
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := Response{
		Status: http.StatusOK,
		Data:   product,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// CreateProduct adds a new product
func (app *Config) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product data.Product

	err := app.readJSON(w, r, &product)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if product.Name == "" || product.Slug == "" || product.Price <= 0 || product.CategoryID == "" {
		app.errorJSON(w, errors.New("missing required fields"), http.StatusBadRequest)
		return
	}

	id, err := app.Models.Product.Insert(product)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Get the newly created product
	createdProduct, err := app.Models.Product.GetByID(id)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := Response{
		Status:  http.StatusCreated,
		Message: "Product created successfully",
		Data:    createdProduct,
	}

	app.writeJSON(w, http.StatusCreated, payload)
}

// UpdateProduct updates an existing product
func (app *Config) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		app.errorJSON(w, errors.New("missing product ID"), http.StatusBadRequest)
		return
	}

	var product data.Product

	err := app.readJSON(w, r, &product)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Ensure ID in URL matches product ID
	product.ID = id

	// Validate required fields
	if product.Name == "" || product.Slug == "" || product.Price <= 0 || product.CategoryID == "" {
		app.errorJSON(w, errors.New("missing required fields"), http.StatusBadRequest)
		return
	}

	err = app.Models.Product.Update(product)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Get updated product
	updatedProduct, err := app.Models.Product.GetByID(id)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := Response{
		Status:  http.StatusOK,
		Message: "Product updated successfully",
		Data:    updatedProduct,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// PatchProduct partially updates an existing product
func (app *Config) PatchProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		app.errorJSON(w, errors.New("missing product ID"), http.StatusBadRequest)
		return
	}

	// First, get the existing product
	existingProduct, err := app.Models.Product.GetByID(id)
	if err != nil {
		if err.Error() == "product not found" {
			app.errorJSON(w, err, http.StatusNotFound)
			return
		}
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Read the partial update
	var partialUpdate map[string]interface{}
	err = app.readJSON(w, r, &partialUpdate)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Apply updates to the existing product
	// Handle basic fields
	if name, ok := partialUpdate["name"].(string); ok && name != "" {
		existingProduct.Name = name
	}

	if slug, ok := partialUpdate["slug"].(string); ok && slug != "" {
		existingProduct.Slug = slug
	}

	if description, ok := partialUpdate["description"].(string); ok {
		existingProduct.Description = description
	}

	if productType, ok := partialUpdate["type"].(string); ok {
		existingProduct.Type = productType
	}

	if price, ok := partialUpdate["price"].(float64); ok && price > 0 {
		existingProduct.Price = price
	}

	if originalPrice, ok := partialUpdate["original_price"].(float64); ok {
		existingProduct.OriginalPrice = originalPrice
	}

	if discount, ok := partialUpdate["discount"].(float64); ok {
		existingProduct.Discount = discount
	}

	if imageURL, ok := partialUpdate["image_url"].(string); ok {
		existingProduct.ImageURL = imageURL
	}

	if categoryID, ok := partialUpdate["category_id"].(string); ok && categoryID != "" {
		existingProduct.CategoryID = categoryID
	}

	if categorySlug, ok := partialUpdate["category_slug"].(string); ok && categorySlug != "" {
		existingProduct.CategorySlug = categorySlug
	}

	if stockQuantity, ok := partialUpdate["stock_quantity"].(float64); ok {
		existingProduct.StockQuantity = int(stockQuantity)
	}

	if brand, ok := partialUpdate["brand"].(string); ok {
		existingProduct.Brand = brand
	}

	// Handle complex fields if present
	if features, ok := partialUpdate["features"]; ok {
		featuresJSON, err := json.Marshal(features)
		if err == nil {
			existingProduct.Features = featuresJSON
		}
	}

	if shippingInfo, ok := partialUpdate["shipping_info"]; ok {
		shippingJSON, err := json.Marshal(shippingInfo)
		if err == nil {
			existingProduct.ShippingInfo = shippingJSON
		}
	}

	// Handle images and tags if they are completely replaced
	if images, ok := partialUpdate["images"].([]interface{}); ok {
		newImages := []data.ProductImage{}
		for _, img := range images {
			if imgMap, ok := img.(map[string]interface{}); ok {
				image := data.ProductImage{}

				if url, ok := imgMap["url"].(string); ok && url != "" {
					image.URL = url
				}

				if isPrimary, ok := imgMap["is_primary"].(bool); ok {
					image.IsPrimary = isPrimary
				}

				if displayOrder, ok := imgMap["display_order"].(float64); ok {
					image.DisplayOrder = int(displayOrder)
				}

				newImages = append(newImages, image)
			}
		}

		if len(newImages) > 0 {
			existingProduct.Images = newImages
		}
	}

	if tags, ok := partialUpdate["tags"].([]interface{}); ok {
		newTags := []string{}
		for _, tag := range tags {
			if tagStr, ok := tag.(string); ok && tagStr != "" {
				newTags = append(newTags, tagStr)
			}
		}

		existingProduct.Tags = newTags
	}

	// Update the product
	err = app.Models.Product.Update(*existingProduct)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Get updated product
	updatedProduct, err := app.Models.Product.GetByID(id)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := Response{
		Status:  http.StatusOK,
		Message: "Product updated successfully",
		Data:    updatedProduct,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// DeleteProduct removes a product
func (app *Config) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		app.errorJSON(w, errors.New("missing product ID"), http.StatusBadRequest)
		return
	}

	err := app.Models.Product.Delete(id)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := Response{
		Status:  http.StatusOK,
		Message: "Product deleted successfully",
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// GetProductReviews returns reviews for a product
func (app *Config) GetProductReviews(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		app.errorJSON(w, errors.New("missing product ID"), http.StatusBadRequest)
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if pageSize <= 0 {
		pageSize = 10
	}

	reviews, totalCount, err := app.Models.Product.GetProductReviews(id, page, pageSize)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := (totalCount + pageSize - 1) / pageSize

	// Build response
	meta := PaginationMeta{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalItems:  totalCount,
		TotalPages:  totalPages,
	}

	payload := Response{
		Status: http.StatusOK,
		Data:   reviews,
		Meta:   meta,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// AddProductReview adds a new review for a product
func (app *Config) AddProductReview(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if productID == "" {
		app.errorJSON(w, errors.New("missing product ID"), http.StatusBadRequest)
		return
	}

	var review data.ProductReview

	err := app.readJSON(w, r, &review)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Set the product ID from URL
	review.ProductID = productID

	// Validate required fields
	if review.UserID == "" || review.Rating < 1 || review.Rating > 5 {
		app.errorJSON(w, errors.New("invalid review data"), http.StatusBadRequest)
		return
	}

	id, err := app.Models.Product.AddProductReview(review)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := Response{
		Status:  http.StatusCreated,
		Message: "Review added successfully",
		Data:    map[string]string{"id": id},
	}

	app.writeJSON(w, http.StatusCreated, payload)
}

// GetAllCategories returns all categories
func (app *Config) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for filtering
	filters := make(map[string]string)

	if parentID := r.URL.Query().Get("parent_id"); parentID != "" {
		filters["parent_id"] = parentID
	}

	if isActive := r.URL.Query().Get("is_active"); isActive != "" {
		filters["is_active"] = isActive
	}

	if isVisible := r.URL.Query().Get("is_visible"); isVisible != "" {
		filters["is_visible"] = isVisible
	}

	if level := r.URL.Query().Get("level"); level != "" {
		filters["level"] = level
	}

	if orderBy := r.URL.Query().Get("order_by"); orderBy != "" {
		filters["order_by"] = orderBy
	}

	if orderDir := r.URL.Query().Get("order_dir"); orderDir != "" {
		filters["order_dir"] = orderDir
	}

	// Get categories
	categories, err := app.Models.Category.GetAllCategories(filters)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := Response{
		Status: http.StatusOK,
		Data:   categories,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// GetCategoryByID returns a category by ID
func (app *Config) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		app.errorJSON(w, errors.New("missing category ID"), http.StatusBadRequest)
		return
	}

	category, err := app.Models.Category.GetCategoryByID(id)
	if err != nil {
		if err.Error() == "category not found" {
			app.errorJSON(w, err, http.StatusNotFound)
			return
		}
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := Response{
		Status: http.StatusOK,
		Data:   category,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// GetCategoryBySlug returns a category by slug
func (app *Config) GetCategoryBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		app.errorJSON(w, errors.New("missing category slug"), http.StatusBadRequest)
		return
	}

	category, err := app.Models.Category.GetCategoryBySlug(slug)
	if err != nil {
		if err.Error() == "category not found" {
			app.errorJSON(w, err, http.StatusNotFound)
			return
		}
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := Response{
		Status: http.StatusOK,
		Data:   category,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

// GetCategoryTree returns the full category hierarchy
func (app *Config) GetCategoryTree(w http.ResponseWriter, r *http.Request) {
	categoryTree, err := app.Models.Category.GetCategoryTree()
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := Response{
		Status: http.StatusOK,
		Data:   categoryTree,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

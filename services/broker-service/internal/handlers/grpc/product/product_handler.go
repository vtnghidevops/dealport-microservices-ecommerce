package grpc

import (
	pb "broker-service/proto/product"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// ProductHandler handles product-related HTTP requests using gRPC
type ProductHandler struct {
	client *ProductClient
}

// NewProductHandler creates a new product gRPC handler
func NewProductHandler(client *ProductClient) *ProductHandler {
	return &ProductHandler{
		client: client,
	}
}

// GetAllProducts handles GET /api/v1/products
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	// Extract page and page size from query parameters
	page := 1
	pageSize := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Check for page_size first, then fall back to limit if page_size is not provided
	pageSizeStr := r.URL.Query().Get("page_size")
	if pageSizeStr == "" {
		pageSizeStr = r.URL.Query().Get("limit") // Support both page_size and limit
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	// Extract filter parameters from the query string
	filters := make(map[string]string)
	for key, values := range r.URL.Query() {
		if key != "page" && key != "page_size" && key != "limit" && len(values) > 0 {
			filters[key] = values[0]
		}
	}

	// Call the gRPC client
	resp, err := h.client.ListProducts(page, pageSize, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response with pagination metadata
	meta := map[string]interface{}{
		"current_page": page,
		"page_size":    pageSize,
		"total_items":  resp.Total,
		"total_pages":  (resp.Total + int32(pageSize) - 1) / int32(pageSize),
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Products retrieved successfully",
		"data":    resp.Products,
		"meta":    meta,
	})
}

// GetProductByID handles GET /api/v1/products/{id}
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Call the gRPC client
	product, err := h.client.GetProduct(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Product retrieved successfully",
		"data":    product,
	})
}

// GetProductBySlug handles GET /api/v1/products/slug/{slug}
func (h *ProductHandler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	// Get the product slug from the URL
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Invalid product slug", http.StatusBadRequest)
		return
	}

	// Call the gRPC client
	product, err := h.client.GetProductBySlug(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Product retrieved successfully",
		"data":    product,
	})
}

// CreateProduct handles POST /api/v1/products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var productRequest map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to protobuf Product
	product := &pb.Product{}

	// Process required fields
	if name, ok := productRequest["name"].(string); ok {
		product.Name = name
	} else {
		http.Error(w, "Name is required and must be a string", http.StatusBadRequest)
		return
	}

	if price, ok := productRequest["price"].(float64); ok {
		product.Price = price
	} else if priceStr, ok := productRequest["price"].(string); ok {
		// Try to parse string as float
		if parsedPrice, err := strconv.ParseFloat(priceStr, 64); err == nil {
			product.Price = parsedPrice
		} else {
			http.Error(w, "Price must be a valid number", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Price is required and must be a number", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if product.Name == "" || product.Price <= 0 {
		http.Error(w, "Name and price (greater than 0) are required", http.StatusBadRequest)
		return
	}

	// Process optional fields
	if description, ok := productRequest["description"].(string); ok {
		product.Description = description
	}

	if slug, ok := productRequest["slug"].(string); ok {
		product.Slug = slug
	}

	if imageUrl, ok := productRequest["imageUrl"].(string); ok {
		product.ImageUrl = imageUrl
	}

	if categoryId, ok := productRequest["categoryId"].(float64); ok {
		product.CategoryId = int32(categoryId)
	}

	if categorySlug, ok := productRequest["categorySlug"].(string); ok {
		product.CategorySlug = categorySlug
	}

	if stockQuantity, ok := productRequest["stockQuantity"].(float64); ok {
		product.StockQuantity = int32(stockQuantity)
	}

	if originalPrice, ok := productRequest["originalPrice"].(float64); ok {
		product.OriginalPrice = originalPrice
	}

	if discount, ok := productRequest["discount"].(float64); ok {
		product.Discount = discount
	}

	if brand, ok := productRequest["brand"].(string); ok {
		product.Brand = brand
	}

	if typeStr, ok := productRequest["type"].(string); ok {
		product.Type = typeStr
	}

	if orders, ok := productRequest["orders"].(float64); ok {
		product.Orders = int32(orders)
	}

	// Handle arrays of strings
	if tags, ok := productRequest["tags"].([]interface{}); ok {
		for _, tag := range tags {
			if tagStr, ok := tag.(string); ok {
				product.Tags = append(product.Tags, tagStr)
			}
		}
	}

	if categories, ok := productRequest["categories"].([]interface{}); ok {
		for _, category := range categories {
			if catStr, ok := category.(string); ok {
				product.Categories = append(product.Categories, catStr)
			}
		}
	}

	if features, ok := productRequest["features"].([]interface{}); ok {
		for _, feature := range features {
			if featureStr, ok := feature.(string); ok {
				product.Features = append(product.Features, featureStr)
			}
		}
	}

	if imgSlider, ok := productRequest["imgSlider"].([]interface{}); ok {
		for _, img := range imgSlider {
			if imgStr, ok := img.(string); ok {
				// Không remove domain - giữ nguyên URL
				product.ImgSlider = append(product.ImgSlider, imgStr)
				log.Println("DEBUG PATCH: Added image URL (keeping domain):", imgStr)
			} else {
				// Cần type assertion khi append interface{} vào []string
				product.ImgSlider = append(product.ImgSlider, fmt.Sprintf("%v", img))
				log.Println("DEBUG PATCH: Added non-string image value:", img)
			}
		}
	}

	// Handle shippingInfo (nested object)
	if shippingInfo, ok := productRequest["shippingInfo"].(map[string]interface{}); ok {
		product.ShippingInfo = &pb.ShippingInfo{}

		if courier, ok := shippingInfo["courier"].(string); ok {
			product.ShippingInfo.Courier = courier
		}

		if local, ok := shippingInfo["local"].(string); ok {
			product.ShippingInfo.Local = local
		}

		if ups, ok := shippingInfo["ups"].(string); ok {
			product.ShippingInfo.Ups = ups
		}

		if global, ok := shippingInfo["global"].(string); ok {
			product.ShippingInfo.Global = global
		}
	}

	// Handle images array (array of objects)
	if images, ok := productRequest["images"].([]interface{}); ok {
		for _, imgData := range images {
			if imgObj, ok := imgData.(map[string]interface{}); ok {
				img := &pb.ProductImage{}

				if url, ok := imgObj["url"].(string); ok {
					img.Url = url
				}

				// Handle field name differences (is_primary in JSON vs isPrimary in Go struct)
				var isPrimary bool
				if val, ok := imgObj["is_primary"].(bool); ok {
					isPrimary = val
				} else if val, ok := imgObj["isPrimary"].(bool); ok {
					isPrimary = val
				}
				img.IsPrimary = isPrimary

				// Handle field name differences for display_order
				if order, ok := imgObj["display_order"].(float64); ok {
					img.DisplayOrder = int32(order)
				} else if order, ok := imgObj["displayOrder"].(float64); ok {
					img.DisplayOrder = int32(order)
				}

				product.Images = append(product.Images, img)
			}
		}
	}

	// Handle reviewsAvg
	if reviewsAvg, ok := productRequest["reviewsAvg"].(map[string]interface{}); ok {
		product.ReviewsAvg = &pb.ProductRating{}

		if rating, ok := reviewsAvg["rating"].(float64); ok {
			product.ReviewsAvg.AverageRating = rating
		} else if rating, ok := reviewsAvg["average_rating"].(float64); ok {
			product.ReviewsAvg.AverageRating = rating
		}

		if count, ok := reviewsAvg["count"].(float64); ok {
			product.ReviewsAvg.Count = int32(count)
		}
	}

	// Handle uiMetadata as JSON string
	if uiMetadata, ok := productRequest["uiMetadata"].(map[string]interface{}); ok {
		// Convert the map back to JSON string
		uiMetadataBytes, err := json.Marshal(uiMetadata)
		if err == nil {
			product.UiMetadata = string(uiMetadataBytes)
		}
	}

	// Call the gRPC client
	resp, err := h.client.CreateProduct(product)
	if err != nil {
		http.Error(w, "Error creating product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Product created successfully",
		"data":    resp.Product,
	})
}

// UpdateProduct handles PUT /api/v1/products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	log.Printf("🔄 Received PUT request for product ID: %d", id)

	// Parse the request body
	var product pb.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set the ID from the URL
	product.Id = int32(id)

	// Call the gRPC client
	resp, err := h.client.UpdateProduct(&product)
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

// PatchProduct handles PATCH /api/v1/products/{id}
func (h *ProductHandler) PatchProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Parse the request body into a map
	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// DEBUG: Log all update fields
	log.Println("DEBUG PATCH: Received updates with fields:", getMapKeys(updates))

	// Log updates for debugging
	if imgSlider, ok := updates["imgSlider"]; ok {
		log.Println("DEBUG PATCH: Updating imgSlider for product", id, "with value: ", imgSlider)

		// Normalize image URLs in imgSlider by removing domain if present
		if imgSliderArr, ok := imgSlider.([]interface{}); ok {
			log.Println("DEBUG PATCH: Original imgSlider array length:", len(imgSliderArr))
			normalizedSlider := make([]interface{}, 0, len(imgSliderArr))
			for _, img := range imgSliderArr {
				if imgStr, ok := img.(string); ok {
					// Không remove domain - giữ nguyên URL
					normalizedSlider = append(normalizedSlider, imgStr)
					log.Println("DEBUG PATCH: Added image URL (keeping domain):", imgStr)
				} else {
					normalizedSlider = append(normalizedSlider, img)
					log.Println("DEBUG PATCH: Added non-string image value:", img)
				}
			}
			updates["imgSlider"] = normalizedSlider
			log.Println("DEBUG PATCH: Normalized imgSlider:", normalizedSlider)
		} else {
			log.Println("DEBUG PATCH: imgSlider is not an array, type:", reflect.TypeOf(imgSlider))
		}
	} else {
		log.Println("DEBUG PATCH: No imgSlider field in updates")
	}

	// Call the gRPC client's PatchProduct method directly
	updatedProduct, err := h.client.PatchProduct(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated product
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Product updated successfully",
		"data":    updatedProduct,
	})
}

// Helper function to get map keys for debugging
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// DeleteProduct handles DELETE /api/v1/products/{id}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Call the gRPC client
	resp, err := h.client.DeleteProduct(id)
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

// GetProductReviews handles GET /api/v1/products/{id}/reviews
func (h *ProductHandler) GetProductReviews(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	idStr := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

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
	resp, err := h.client.GetProductReviews(productID, page, pageSize)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare pagination metadata
	meta := map[string]interface{}{
		"current_page": page,
		"page_size":    pageSize,
		"total_items":  resp.Total,
		"total_pages":  (resp.Total + int32(pageSize) - 1) / int32(pageSize),
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Product reviews retrieved successfully",
		"data":    resp.Reviews,
		"meta":    meta,
	})
}

// AddProductReview handles POST /api/v1/products/{id}/reviews
func (h *ProductHandler) AddProductReview(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	idStr := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var reviewReq ProductReviewRequest
	err = json.NewDecoder(r.Body).Decode(&reviewReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if reviewReq.UserName == "" || reviewReq.Rating <= 0 || reviewReq.Rating > 5 {
		http.Error(w, "userName and rating (1-5) are required", http.StatusBadRequest)
		return
	}

	// Create the gRPC request
	review := &pb.ProductReview{
		ProductId:  int32(productID),
		UserName:   reviewReq.UserName,
		Rating:     int32(reviewReq.Rating),
		ReviewText: reviewReq.ReviewText,
	}

	// Call the gRPC client
	resp, err := h.client.AddProductReview(review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Product review added successfully",
		"data": map[string]interface{}{
			"id": resp.Id,
		},
	})
}

// UpdateProductReview handles PUT /api/v1/products/{id}/reviews/{reviewId}
func (h *ProductHandler) UpdateProductReview(w http.ResponseWriter, r *http.Request) {
	// Get the product ID and review ID from the URL
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	reviewIDStr := chi.URLParam(r, "reviewId")
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var reviewReq ProductReviewRequest
	err = json.NewDecoder(r.Body).Decode(&reviewReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if reviewReq.UserName == "" || reviewReq.Rating <= 0 || reviewReq.Rating > 5 {
		http.Error(w, "userName and rating (1-5) are required", http.StatusBadRequest)
		return
	}

	// Create the gRPC request
	review := &pb.ProductReview{
		Id:         int32(reviewID),
		ProductId:  int32(productID),
		UserName:   reviewReq.UserName,
		Rating:     int32(reviewReq.Rating),
		ReviewText: reviewReq.ReviewText,
	}

	// Call the gRPC client
	resp, err := h.client.UpdateProductReview(review)
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

// UploadProductImage handles POST /api/v1/products/{id}/images
func (h *ProductHandler) UploadProductImage(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	idStr := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("ERROR: Invalid product ID: %v", err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Parse the multipart form
	err = r.ParseMultipartForm(32 << 20) // 32MB limit
	if err != nil {
		log.Printf("ERROR: Invalid form data: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("ERROR: No image file provided: %v", err)
		http.Error(w, "No image file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file content
	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Printf("ERROR: Error reading file: %v", err)
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Check if the image should be set as primary
	isPrimary := false
	if r.FormValue("isPrimary") == "true" {
		isPrimary = true
	}

	// Get the content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg" // Default
	}

	log.Printf("Uploading image: filename=%s, size=%d, contentType=%s, isPrimary=%v",
		header.Filename, len(fileData), contentType, isPrimary)

	// Call the gRPC client
	resp, err := h.client.UploadProductImage(productID, fileData, header.Filename, contentType, isPrimary)
	if err != nil {
		log.Printf("ERROR: Failed to upload image: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Image uploaded successfully: %s", resp.Url)

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Product image uploaded successfully",
		"data": map[string]interface{}{
			"url": resp.Url,
		},
	})
}

// DeleteProductImage handles DELETE /api/v1/products/{id}/images/{imageId}
func (h *ProductHandler) DeleteProductImage(w http.ResponseWriter, r *http.Request) {
	// Get the image ID from the URL
	imageIDStr := chi.URLParam(r, "imageId")
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	// Call the gRPC client
	resp, err := h.client.DeleteProductImage(imageID)
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

// SetPrimaryProductImage handles PUT /api/v1/products/{id}/images/{imageId}/primary
func (h *ProductHandler) SetPrimaryProductImage(w http.ResponseWriter, r *http.Request) {
	// Get the product ID and image ID from the URL
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	imageIDStr := chi.URLParam(r, "imageId")
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	// Call the gRPC client
	resp, err := h.client.SetPrimaryProductImage(productID, imageID)
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

// GetTopRatedTestimonials handles GET /api/v1/testimonials
func (h *ProductHandler) GetTopRatedTestimonials(w http.ResponseWriter, r *http.Request) {
	// Get the limit parameter
	limit := 5
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Call the gRPC client
	resp, err := h.client.GetTopRatedTestimonials(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Testimonials retrieved successfully",
		"data":    resp.Testimonials,
	})
}

// GetProductHealth handles GET /api/v1/health
func (h *ProductHandler) GetProductHealth(w http.ResponseWriter, r *http.Request) {
	status, err := h.client.CheckHealth()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    status,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// ProxyProductImage handles GET /api/products/images/{imageFile}
func (h *ProductHandler) ProxyProductImage(w http.ResponseWriter, r *http.Request) {
	// Get the image filename from URL
	imageFile := chi.URLParam(r, "imageFile")

	// Call the gRPC client to get the image data
	imageData, contentType, err := h.client.GetProductImageFile(imageFile)
	if err != nil {
		// If there's an error, return a 404 Not Found
		http.NotFound(w, r)
		return
	}

	// Set the content type header
	w.Header().Set("Content-Type", contentType)

	// Set cache control headers to improve performance
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours

	// Write the image data to the response
	w.Write(imageData)
}

// ProductReviewRequest represents the JSON structure for product review requests
type ProductReviewRequest struct {
	UserName   string  `json:"userName"`
	Rating     float64 `json:"rating"`
	ReviewText string  `json:"reviewText,omitempty"`
}

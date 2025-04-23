package products

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"product-service/internal/domain"
	"product-service/internal/handler"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// GetAll returns all products with pagination and filtering
func GetAll(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get page and page_size from query parameters
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("page_size")

		page := 1
		pageSize := 10

		if pageStr != "" {
			pageVal, err := strconv.Atoi(pageStr)
			if err == nil && pageVal > 0 {
				page = pageVal
			}
		}

		if pageSizeStr != "" {
			pageSizeVal, err := strconv.Atoi(pageSizeStr)
			if err == nil && pageSizeVal > 0 {
				pageSize = pageSizeVal
			}
		}

		// Build filters from query parameters
		filters := make(map[string]string)

		if category := r.URL.Query().Get("category"); category != "" {
			// Try to convert category to int if it looks like a numeric ID
			if _, err := strconv.Atoi(category); err == nil {
				filters["category_id"] = category
			} else {
				filters["category_slug"] = category
			}
		}

		if brand := r.URL.Query().Get("brand"); brand != "" {
			filters["brand"] = brand
		}

		if productType := r.URL.Query().Get("type"); productType != "" {
			filters["type"] = productType
		}

		// Get products from service
		products, totalCount, err := app.ProductService.GetAllProducts(page, pageSize, filters)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Remove ui_metadata from products that aren't top-sale type
		for i := range products {
			if products[i].Type != "top-sale" {
				products[i].UIMetadata = nil
			}
		}

		// Calculate pagination metadata
		totalPages := (totalCount + pageSize - 1) / pageSize
		meta := handler.PaginationMeta{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalItems:  totalCount,
			TotalPages:  totalPages,
		}

		payload := handler.Response{
			Status: http.StatusOK,
			Data:   products,
			Meta:   meta,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// GetByID returns a product by ID
func GetByID(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from URL parameters
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing product ID"), http.StatusBadRequest)
			return
		}

		// Convert string ID to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid product ID format"), http.StatusBadRequest)
			return
		}

		// Get product from service
		product, err := app.ProductService.GetProductByID(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusNotFound)
			return
		}

		// Only include ui_metadata for top-sale products
		if product.Type != "top-sale" {
			product.UIMetadata = nil
		}

		payload := handler.Response{
			Status: http.StatusOK,
			Data:   product,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// GetBySlug returns a product by slug
func GetBySlug(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get slug from URL parameters
		slug := chi.URLParam(r, "slug")
		if slug == "" {
			app.ErrorResponse(w, errors.New("missing product slug"), http.StatusBadRequest)
			return
		}

		// Get product from service
		product, err := app.ProductService.GetProductBySlug(slug)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusNotFound)
			return
		}

		// Only include ui_metadata for top-sale products
		if product.Type != "top-sale" {
			product.UIMetadata = nil
		}

		payload := handler.Response{
			Status: http.StatusOK,
			Data:   product,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// Create creates a new product
func Create(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product domain.Product

		err := app.ParseJSON(w, r, &product)
		fmt.Println("Product received:", product.Name)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		// Validate required fields
		if product.Name == "" || product.Slug == "" || product.Price <= 0 || product.CategoryID == 0 {
			app.ErrorResponse(w, errors.New("missing required fields"), http.StatusBadRequest)
			return
		}

		// Process imgSlider to create Images array if not provided
		if len(product.ImgSlider) > 0 && len(product.Images) == 0 {
			// Convert imgSlider to Images
			for i, url := range product.ImgSlider {
				product.Images = append(product.Images, domain.ProductImage{
					URL:          url,
					IsPrimary:    i == 0, // First image is primary
					DisplayOrder: i,
				})
			}

			// Set primary image URL if not already set
			if product.ImageURL == "" && len(product.ImgSlider) > 0 {
				product.ImageURL = product.ImgSlider[0]
			}
		}

		id, err := app.ProductService.CreateProduct(&product)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Get the newly created product
		createdProduct, err := app.ProductService.GetProductByID(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Create response with modified structure
		responseData := map[string]interface{}{
			"id":             createdProduct.ID,
			"name":           createdProduct.Name,
			"slug":           createdProduct.Slug,
			"description":    createdProduct.Description,
			"price":          createdProduct.Price,
			"original_price": createdProduct.OriginalPrice,
			"discount":       createdProduct.Discount,
			"category_id":    createdProduct.CategoryID,
			"category_slug":  createdProduct.CategorySlug,
			"stock_quantity": createdProduct.StockQuantity,
			"type":           createdProduct.Type,
			"brand":          createdProduct.Brand,
			"features":       createdProduct.Features,
			"shipping_info":  createdProduct.ShippingInfo,
			"created_at":     createdProduct.CreatedAt,
			"updated_at":     createdProduct.UpdatedAt,
		}

		// Extract image URLs into imageSlider array and find primary image
		imageSlider := []string{}
		var primaryImageURL string

		if createdProduct.Images != nil {
			for _, img := range createdProduct.Images {
				imageSlider = append(imageSlider, img.URL)
				if img.IsPrimary {
					primaryImageURL = img.URL
				}
			}

			// If no primary image found but we have images, use the first one
			if primaryImageURL == "" && len(imageSlider) > 0 {
				primaryImageURL = imageSlider[0]
			}
		}

		responseData["imageSlider"] = imageSlider
		responseData["image_url"] = primaryImageURL
		// Handle review ratings
		if createdProduct.ReviewsAvg.Count > 0 {
			// Round the average rating to 1 decimal place (làm tròn)
			roundedAvg := math.Round(createdProduct.ReviewsAvg.AverageRating*10) / 10
			responseData["reviewsAvg"] = map[string]interface{}{
				"count":          createdProduct.ReviewsAvg.Count,
				"average_rating": roundedAvg,
			}
		}

		payload := handler.Response{
			Status:  http.StatusCreated,
			Message: "Product created successfully",
			Data:    responseData,
		}

		app.SuccessResponse(w, http.StatusCreated, payload)
	}
}

// Update updates an existing product
func Update(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from URL parameter
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing product ID"), http.StatusBadRequest)
			return
		}

		// Convert string ID to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid product ID format"), http.StatusBadRequest)
			return
		}

		var product domain.Product

		err = app.ParseJSON(w, r, &product)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		// Set ID from URL
		product.ID = id

		// Validate required fields
		if product.Name == "" || product.Slug == "" || product.Price <= 0 || product.CategoryID == 0 {
			app.ErrorResponse(w, errors.New("missing required fields"), http.StatusBadRequest)
			return
		}

		// Check and preserve ReviewsAvg if not provided
		if product.ReviewsAvg.Count == 0 && product.ReviewsAvg.AverageRating == 0 {
			// Get existing product to preserve reviewsAvg value
			existingProduct, err := app.ProductService.GetProductByID(id)
			if err == nil {
				product.ReviewsAvg = existingProduct.ReviewsAvg
			}
		}

		err = app.ProductService.UpdateProduct(&product)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Get updated product
		updatedProduct, err := app.ProductService.GetProductByID(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Product updated successfully",
			Data:    updatedProduct,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// Patch applies partial updates to an existing product
func Patch(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from URL parameter
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing product ID"), http.StatusBadRequest)
			return
		}

		// Convert string ID to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid product ID format"), http.StatusBadRequest)
			return
		}

		// Get existing product
		existingProduct, err := app.ProductService.GetProductByID(id)
		if err != nil {
			if errors.Is(err, domain.ErrProductNotFound) {
				app.ErrorResponse(w, errors.New("product not found"), http.StatusNotFound)
				return
			}
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Parse JSON request (allows partial fields)
		var updates map[string]interface{}
		err = app.ParseJSON(w, r, &updates)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		// Apply updates to the existing product
		for field, value := range updates {
			switch field {
			case "name":
				if name, ok := value.(string); ok && name != "" {
					existingProduct.Name = name
				}
			case "description":
				if description, ok := value.(string); ok {
					existingProduct.Description = description
				}
			case "slug":
				if slug, ok := value.(string); ok && slug != "" {
					existingProduct.Slug = slug
				}
			case "price":
				if price, ok := value.(float64); ok && price > 0 {
					existingProduct.Price = price
				}
			case "original_price", "originalPrice":
				if originalPrice, ok := value.(float64); ok {
					existingProduct.OriginalPrice = originalPrice
				}
			case "discount":
				if discount, ok := value.(float64); ok {
					existingProduct.Discount = discount
				}
			case "category_id", "categoryId":
				if catIDFloat, ok := value.(float64); ok {
					existingProduct.CategoryID = int(catIDFloat)
				}
			case "category_slug", "categorySlug":
				if catSlug, ok := value.(string); ok && catSlug != "" {
					existingProduct.CategorySlug = catSlug
				}
			case "stock_quantity", "stockQuantity":
				if stockQuantity, ok := value.(float64); ok {
					existingProduct.StockQuantity = int(stockQuantity)
				}
			case "type":
				if productType, ok := value.(string); ok {
					existingProduct.Type = productType
				}
			case "brand":
				if brand, ok := value.(string); ok {
					existingProduct.Brand = brand
				}
			case "features":
				if features, ok := value.([]interface{}); ok {
					stringFeatures := make([]string, 0, len(features))
					for _, f := range features {
						if featureStr, ok := f.(string); ok {
							stringFeatures = append(stringFeatures, featureStr)
						}
					}
					existingProduct.Features = stringFeatures
				}
			case "shipping_info", "shippingInfo":
				if shippingInfo, ok := value.(map[string]interface{}); ok {
					existingProduct.ShippingInfo = domain.ShippingInfo{
						Courier: getString(shippingInfo, "courier"),
						Local:   getString(shippingInfo, "local"),
						Ups:     getString(shippingInfo, "ups"),
						Global:  getString(shippingInfo, "global"),
					}
				}
			case "images", "imageSlider", "imgSlider":
				// Handle images update
				if imageURLs, ok := value.([]interface{}); ok && len(imageURLs) > 0 {
					images := make([]domain.ProductImage, 0, len(imageURLs))

					// Get current images for reference
					currentImages, err := app.ProductService.GetProductImages(id)
					if err != nil {
						// Handle error gracefully, create new images array if can't get current images
						currentImages = []domain.ProductImage{}
					}

					// Create a map to track existing images by URL for easier lookup
					existingImagesByURL := make(map[string]domain.ProductImage)
					for _, img := range currentImages {
						existingImagesByURL[img.URL] = img
					}

					// Process each image URL
					for i, imgVal := range imageURLs {
						if url, ok := imgVal.(string); ok && url != "" {
							// Check if this image already exists
							var imageID int

							// If we already have this image, reuse its ID
							if existingImg, found := existingImagesByURL[url]; found {
								imageID = existingImg.ID
							} else {
								// New image ID will be assigned by database (AUTO_INCREMENT)
								imageID = 0
							}

							// Create new image object
							images = append(images, domain.ProductImage{
								ID:           imageID,
								ProductID:    id,
								URL:          url,
								IsPrimary:    i == 0, // First image is primary
								DisplayOrder: i,
								CreatedAt:    time.Now(),
							})
						}
					}

					// Update the product's images
					if len(images) > 0 {
						existingProduct.Images = images

						// Update image_url to the primary image
						for _, img := range images {
							if img.IsPrimary {
								existingProduct.ImageURL = img.URL
								break
							}
						}

						// If no primary image, use the first one
						if existingProduct.ImageURL == "" && len(images) > 0 {
							existingProduct.ImageURL = images[0].URL
						}

						// Update imgSlider as well
						imgSlider := make([]string, 0, len(images))
						for _, img := range images {
							imgSlider = append(imgSlider, img.URL)
						}
						existingProduct.ImgSlider = imgSlider
					}
				}
			case "tags":
				if tags, ok := value.([]interface{}); ok {
					stringTags := make([]string, 0, len(tags))
					for _, t := range tags {
						if tagStr, ok := t.(string); ok {
							stringTags = append(stringTags, tagStr)
						}
					}
					existingProduct.Tags = stringTags
				}
			case "orders":
				if orders, ok := value.(float64); ok {
					existingProduct.Orders = int(orders)
				}
			case "ui_metadata", "uiMetadata":
				// UI metadata is handled as raw JSON
				existingProduct.UIMetadata = nil // Clear existing metadata

				// If type is map or array, convert to JSON
				switch v := value.(type) {
				case map[string]interface{}:
					// Handle as JSON object
					jsonBytes, err := json.Marshal(v)
					if err == nil {
						existingProduct.UIMetadata = jsonBytes
					}
				case []interface{}:
					// Handle as JSON array
					jsonBytes, err := json.Marshal(v)
					if err == nil {
						existingProduct.UIMetadata = jsonBytes
					}
				}
			}
		}

		// Update the product
		err = app.ProductService.UpdateProduct(existingProduct)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Get the updated product
		updatedProduct, err := app.ProductService.GetProductByID(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Product updated successfully",
			Data:    updatedProduct,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// Helper function to extract string value from map
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

// Delete removes a product
func Delete(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from URL parameter
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing product ID"), http.StatusBadRequest)
			return
		}

		// Convert string ID to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid product ID format"), http.StatusBadRequest)
			return
		}

		err = app.ProductService.DeleteProduct(id)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Product deleted successfully",
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// GetReviews returns all reviews for a product
func GetReviews(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get product ID from URL parameter
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing product ID"), http.StatusBadRequest)
			return
		}

		// Convert string ID to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid product ID format"), http.StatusBadRequest)
			return
		}

		// Get page and page_size from query parameters
		pageStr := r.URL.Query().Get("page")
		pageSizeStr := r.URL.Query().Get("page_size")

		page := 1
		pageSize := 10

		if pageStr != "" {
			pageVal, err := strconv.Atoi(pageStr)
			if err == nil && pageVal > 0 {
				page = pageVal
			}
		}

		if pageSizeStr != "" {
			pageSizeVal, err := strconv.Atoi(pageSizeStr)
			if err == nil && pageSizeVal > 0 {
				pageSize = pageSizeVal
			}
		}

		// Get reviews from service
		reviews, totalCount, err := app.ProductService.GetProductReviews(id, page, pageSize)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		// Calculate pagination metadata
		totalPages := (totalCount + pageSize - 1) / pageSize
		meta := handler.PaginationMeta{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalItems:  totalCount,
			TotalPages:  totalPages,
		}

		payload := handler.Response{
			Status: http.StatusOK,
			Data:   reviews,
			Meta:   meta,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// AddReview adds a new review for a product
func AddReview(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get product ID from URL parameter
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			app.ErrorResponse(w, errors.New("missing product ID"), http.StatusBadRequest)
			return
		}

		// Convert string ID to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid product ID format"), http.StatusBadRequest)
			return
		}

		var review domain.ProductReview

		err = app.ParseJSON(w, r, &review)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		review.ProductID = id

		// Validate required fields
		if review.UserID == "" || review.Rating < 1 || review.Rating > 5 {
			app.ErrorResponse(w, errors.New("missing required fields or invalid rating"), http.StatusBadRequest)
			return
		}

		reviewID, err := app.ProductService.AddProductReview(&review)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusCreated,
			Message: "Review added successfully",
			Data: map[string]interface{}{
				"id": reviewID,
			},
		}

		app.SuccessResponse(w, http.StatusCreated, payload)
	}
}

// UpdateReview updates an existing review for a product
func UpdateReview(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get IDs from URL parameters
		productIDStr := chi.URLParam(r, "id")
		reviewIDStr := chi.URLParam(r, "reviewId")

		if productIDStr == "" || reviewIDStr == "" {
			app.ErrorResponse(w, errors.New("missing product ID or review ID"), http.StatusBadRequest)
			return
		}

		// Convert string IDs to int
		productID, err := strconv.Atoi(productIDStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid product ID format"), http.StatusBadRequest)
			return
		}

		reviewID, err := strconv.Atoi(reviewIDStr)
		if err != nil {
			app.ErrorResponse(w, errors.New("invalid review ID format"), http.StatusBadRequest)
			return
		}

		var reviewUpdate domain.ProductReview

		err = app.ParseJSON(w, r, &reviewUpdate)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		// Set IDs from URL
		reviewUpdate.ProductID = productID
		reviewUpdate.ID = reviewID

		// Validate required fields
		if reviewUpdate.UserID == "" || reviewUpdate.Rating < 1 || reviewUpdate.Rating > 5 {
			app.ErrorResponse(w, errors.New("missing required fields or invalid rating"), http.StatusBadRequest)
			return
		}

		err = app.ProductService.UpdateProductReview(&reviewUpdate)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status:  http.StatusOK,
			Message: "Review updated successfully",
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

// GetTopRatedTestimonials returns random top-rated reviews for HappyCustomers section
func GetTopRatedTestimonials(app *handler.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get limit from query parameter, default to 7 if not provided
		limitStr := r.URL.Query().Get("limit")
		limit := 7

		if limitStr != "" {
			limitVal, err := strconv.Atoi(limitStr)
			if err == nil && limitVal > 0 {
				limit = limitVal
			}
		}

		// Get random top-rated reviews
		testimonials, err := app.ProductService.GetRandomTopRatedReviews(limit)
		if err != nil {
			app.ErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		payload := handler.Response{
			Status: http.StatusOK,
			Data:   testimonials,
		}

		app.SuccessResponse(w, http.StatusOK, payload)
	}
}

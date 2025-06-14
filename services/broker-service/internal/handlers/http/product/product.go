package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	productGrpc "broker-service/internal/handlers/grpc/product"

	"github.com/go-chi/chi/v5"
)

// Config holds dependencies for HTTP proxy handlers
type Config struct {
	ProductClient *productGrpc.ProductClient
}

// ProductServiceURL is the URL for the product service
// Making it a var instead of const allows tests to override it
var ProductServiceURL = "http://product-service:8082"

// ProxyRequest forwards the request to the product service
func ProxyRequest(w http.ResponseWriter, r *http.Request, endpoint string) {
	// Get the URL path parameters
	chiParams := chi.RouteContext(r.Context()).URLParams

	// Build the target URL
	// Use the global ProductServiceURL, which might be overridden in tests
	targetURL := ProductServiceURL + endpoint

	// Replace path parameters in the URL
	for i, key := range chiParams.Keys {
		targetURL = strings.Replace(targetURL, "{"+key+"}", chiParams.Values[i], -1)
	}

	// Make sure to forward query parameters
	if r.URL.RawQuery != "" {
		// Check if we need to rename 'limit' to 'page_size' for product service
		if strings.Contains(endpoint, "/products") && !strings.Contains(r.URL.RawQuery, "page_size=") && strings.Contains(r.URL.RawQuery, "limit=") {
			// Replace 'limit=' with 'page_size=' for product service compatibility
			rawQuery := strings.Replace(r.URL.RawQuery, "limit=", "page_size=", 1)
			targetURL = targetURL + "?" + rawQuery
		} else {
			// Keep original query parameters for all other endpoints
			targetURL = targetURL + "?" + r.URL.RawQuery
		}

		log.Printf("Forwarding request to %s with query: %s", targetURL, r.URL.RawQuery)
	}

	// Read the request body
	var requestBody []byte
	if r.Body != nil {
		requestBody, _ = io.ReadAll(r.Body)
		r.Body.Close()
		// Recreate the body for further use
		r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	}

	// Create a new request to the product service
	req, err := http.NewRequest(r.Method, targetURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error creating proxy request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// Send the request to the product service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending proxy request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response from the product service
	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy headers
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Copy body
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}

// Product Proxy Handlers

// GetAllProducts handles GET /api/v1/products
func (c *Config) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/products")
}

// CreateProduct handles POST /api/v1/products
func (c *Config) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/products")
}

// GetProductByID handles GET /api/v1/products/{id}
func (c *Config) GetProductByID(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/products/{id}")
}

// GetProductBySlug handles GET /api/v1/products/slug/{slug}
func (c *Config) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/products/slug/{slug}")
}

// UpdateProduct handles PUT /api/v1/products/{id}
func (c *Config) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/products/{id}")
}

// PatchProduct handles PATCH /api/v1/products/{id}
func (c *Config) PatchProduct(w http.ResponseWriter, r *http.Request) {
	// Parse the request body first to filter blob URLs
	var updates map[string]interface{}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Close the original body
	r.Body.Close()

	// Parse the body into a map
	if err := json.Unmarshal(body, &updates); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Log updates for debugging
	log.Println("DEBUG PATCH HTTP: Received updates with fields:", getMapKeys(updates))

	// Filter imgSlider URLs - remove blob: URLs
	if imgSlider, ok := updates["imgSlider"]; ok {
		log.Println("DEBUG PATCH HTTP: Found imgSlider in request")

		if imgSliderArr, ok := imgSlider.([]interface{}); ok {
			log.Println("DEBUG PATCH HTTP: Original imgSlider array length:", len(imgSliderArr))
			normalizedSlider := make([]interface{}, 0, len(imgSliderArr))

			for _, img := range imgSliderArr {
				if imgStr, ok := img.(string); ok {
					// Filter out blob: and data: URLs
					if strings.HasPrefix(imgStr, "blob:") || strings.HasPrefix(imgStr, "data:") {
						log.Println("WARNING HTTP: Removing blob/data URL from PATCH request:", imgStr[:15]+"...")
						continue
					}

					// Nếu URL chứa tham số truy vấn (presigned URL), cắt bỏ tất cả tham số truy vấn
					if strings.Contains(imgStr, "?") {
						// Nếu URL chứa đường dẫn /images/products/, lấy chỉ đường dẫn cơ bản
						if strings.Contains(imgStr, "/images/products/") {
							urlPattern := regexp.MustCompile(`(/images/products/[^?]+)`)
							matches := urlPattern.FindStringSubmatch(imgStr)
							if len(matches) > 0 {
								imgStr = matches[1]
								log.Println("DEBUG PATCH HTTP: Normalized presigned URL to:", imgStr)
							}
						} else {
							// Cắt bỏ tất cả tham số truy vấn
							imgStr = strings.Split(imgStr, "?")[0]
							log.Println("DEBUG PATCH HTTP: Removed query parameters from URL:", imgStr)
						}
					}

					// Keep valid URLs
					normalizedSlider = append(normalizedSlider, imgStr)
					log.Println("DEBUG PATCH HTTP: Added image URL:", imgStr)
				} else {
					// Keep non-string values
					normalizedSlider = append(normalizedSlider, img)
				}
			}

			// Replace with filtered array
			updates["imgSlider"] = normalizedSlider
			log.Println("DEBUG PATCH HTTP: Filtered imgSlider array length:", len(normalizedSlider))
		}
	}

	// Process images array if present
	if images, ok := updates["images"].([]interface{}); ok {
		log.Println("DEBUG PATCH HTTP: Found images array in request with length:", len(images))
		normalizedImages := make([]interface{}, 0, len(images))

		for _, imgData := range images {
			if imgObj, ok := imgData.(map[string]interface{}); ok {
				// Process each image object
				if url, ok := imgObj["url"].(string); ok {
					// Filter out blob URLs
					if strings.HasPrefix(url, "blob:") || strings.HasPrefix(url, "data:") {
						log.Println("WARNING HTTP: Removing blob/data URL from images array:", url[:15]+"...")
						continue
					}

					// Normalize presigned URLs
					if strings.Contains(url, "?") {
						if strings.Contains(url, "/images/products-api/") {
							urlPattern := regexp.MustCompile(`(/images/products-api/[^?]+)`)
							matches := urlPattern.FindStringSubmatch(url)
							if len(matches) > 0 {
								imgObj["url"] = matches[1]
								log.Println("DEBUG PATCH HTTP: Normalized MinIO image URL in images array to:", matches[1])
							}
						} else if strings.Contains(url, "/images/products/") {
							urlPattern := regexp.MustCompile(`(/images/products/[^?]+)`)
							matches := urlPattern.FindStringSubmatch(url)
							if len(matches) > 0 {
								imgObj["url"] = matches[1]
								log.Println("DEBUG PATCH HTTP: Normalized local storage image URL in images array to:", matches[1])
							}
						} else {
							imgObj["url"] = strings.Split(url, "?")[0]
							log.Println("DEBUG PATCH HTTP: Removed query parameters from image URL in images array:", imgObj["url"])
						}
					}
				}
				normalizedImages = append(normalizedImages, imgObj)
			} else {
				normalizedImages = append(normalizedImages, imgData)
			}
		}

		updates["images"] = normalizedImages
		log.Println("DEBUG PATCH HTTP: Normalized images array length:", len(normalizedImages))
	}

	// Check imageUrl for blob: URLs
	if imageUrl, ok := updates["imageUrl"].(string); ok {
		if strings.HasPrefix(imageUrl, "blob:") || strings.HasPrefix(imageUrl, "data:") {
			log.Println("WARNING HTTP: Removing blob/data imageUrl from PATCH request:", imageUrl[:15]+"...")
			delete(updates, "imageUrl")
		} else if strings.Contains(imageUrl, "?") {
			// Nếu URL chứa tham số truy vấn (presigned URL), cắt bỏ tất cả tham số truy vấn
			if strings.Contains(imageUrl, "/images/products-api/") {
				urlPattern := regexp.MustCompile(`(/images/products-api/[^?]+)`)
				matches := urlPattern.FindStringSubmatch(imageUrl)
				if len(matches) > 0 {
					updates["imageUrl"] = matches[1]
					log.Println("DEBUG PATCH HTTP: Normalized imageUrl MinIO presigned URL to:", matches[1])
				}
			} else if strings.Contains(imageUrl, "/images/products/") {
				urlPattern := regexp.MustCompile(`(/images/products/[^?]+)`)
				matches := urlPattern.FindStringSubmatch(imageUrl)
				if len(matches) > 0 {
					updates["imageUrl"] = matches[1]
					log.Println("DEBUG PATCH HTTP: Normalized imageUrl local storage URL to:", matches[1])
				}
			} else {
				// Cắt bỏ tất cả tham số truy vấn
				updates["imageUrl"] = strings.Split(imageUrl, "?")[0]
				log.Println("DEBUG PATCH HTTP: Removed query parameters from imageUrl:", updates["imageUrl"])
			}
		}
	}

	// Create new JSON body
	newBody, err := json.Marshal(updates)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	// Create new body reader
	r.Body = io.NopCloser(bytes.NewBuffer(newBody))

	// Continue with proxying the modified request
	ProxyRequest(w, r, "/api/v1/products/{id}")
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
func (c *Config) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/products/{id}")
}

// GetProductReviews handles GET /api/v1/products/{id}/reviews
func (c *Config) GetProductReviews(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/products/{id}/reviews")
}

// AddProductReview handles POST /api/v1/products/{id}/reviews
func (c *Config) AddProductReview(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/products/{id}/reviews")
}

// UpdateProductReview handles PUT /api/v1/products/{id}/reviews/{reviewId}
func (c *Config) UpdateProductReview(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/products/{id}/reviews/{reviewId}")
}

// UploadProductImage handles POST /api/v1/products/{id}/images
func (c *Config) UploadProductImage(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/products/{id}/images")
}

// DeleteProductImage handles DELETE /api/v1/products/{id}/images/{imageId}
func (c *Config) DeleteProductImage(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/products/{id}/images/{imageId}")
}

// SetPrimaryProductImage handles PUT /api/v1/products/{id}/images/{imageId}/primary
func (c *Config) SetPrimaryProductImage(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/api/v1/products/{id}/images/{imageId}/primary")
}

// GetProductHealth handles GET /api/v1/products/health
func (c *Config) GetProductHealth(w http.ResponseWriter, r *http.Request) {
	ProxyRequest(w, r, "/health")
}

// GetProductImage handles GET /api/products/images/{filename}
func (c *Config) GetProductImage(w http.ResponseWriter, r *http.Request) {
	// Get the filename from the URL
	filename := chi.URLParam(r, "filename")
	if filename == "" {
		http.Error(w, "Missing image filename", http.StatusBadRequest)
		return
	}

	// Extract the full image path from the request
	fullPath := r.URL.Path

	// Determine image source based on path pattern
	// MinIO images use the special /images/products-api/ pattern
	// Local storage images use /images/products/ pattern
	isMinIOImage := strings.HasPrefix(filename, "products-api/") ||
		strings.Contains(fullPath, "/images/products-api/")

	log.Printf("Image request for %s (MinIO storage: %v), full path: %s", filename, isMinIOImage, fullPath)

	// For MinIO images, always try to get a presigned URL
	// For local storage images, try to get the file data directly
	var imageData []byte
	var contentType string
	var presignedURL string
	var redirectToURL bool
	var err error

	// For direct HTTP access to local files, check if this is a direct access to /images/products/
	// This handles case where frontend directly accesses http://localhost:58080/images/products/...
	if strings.Contains(fullPath, "/images/products/") && !isMinIOImage {
		// Just pass through the request for local files - no need to get from product service
		log.Printf("Direct access to local storage image: %s", fullPath)

		// Extract the filename part
		parts := strings.Split(fullPath, "/images/products/")
		if len(parts) == 2 {
			localFilename := parts[1]
			localPath := filepath.Join("uploads", "products", localFilename)

			// Attempt to read the file directly
			fileData, err := os.ReadFile(localPath)
			if err == nil {
				// Determine content type
				contentType := http.DetectContentType(fileData)

				// Set headers and return the file
				w.Header().Set("Content-Type", contentType)
				w.Header().Set("Cache-Control", "public, max-age=86400") // 24 hours cache
				w.WriteHeader(http.StatusOK)
				w.Write(fileData)
				return
			} else {
				log.Printf("Failed to read local file %s directly: %v, falling back to product service", localPath, err)
			}
		}
	}

	// If direct file access didn't work or isn't applicable, use the product service
	imageData, contentType, presignedURL, redirectToURL, err = c.ProductClient.GetProductImageFile(filename)
	if err != nil {
		log.Printf("Error getting product image file: %v", err)
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	// For MinIO images, we should always get a presigned URL and redirect
	if isMinIOImage && presignedURL == "" {
		log.Printf("WARNING: MinIO image detected but no presigned URL returned for %s", filename)
	}

	// If we have a presigned URL and should redirect, do a redirect
	if redirectToURL && presignedURL != "" {
		// Log only the base URL without query params for security
		baseURL := strings.Split(presignedURL, "?")[0]
		log.Printf("Redirecting to MinIO URL: %s (credentials omitted)", baseURL)

		// Use temporary redirect (302) instead of Found (302) to better handle CORS
		// This ensures the browser keeps the original request headers
		http.Redirect(w, r, presignedURL, http.StatusFound)
		return
	}

	// Set content type header
	w.Header().Set("Content-Type", contentType)

	// Set Cache-Control header for better performance
	w.Header().Set("Cache-Control", "public, max-age=86400") // 24 hours cache

	// Write the image data to the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(imageData)
	if err != nil {
		log.Printf("Error writing image data to response: %v", err)
	}
}

package http

import (
	"bytes"
	"io"
	"log"
	"net/http"
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
	ProxyRequest(w, r, "/api/v1/products/{id}")
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

	// Use gRPC client to get the image data
	imageData, contentType, err := c.ProductClient.GetProductImageFile(filename)
	if err != nil {
		log.Printf("Error getting product image file: %v", err)
		http.Error(w, "Image not found", http.StatusNotFound)
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

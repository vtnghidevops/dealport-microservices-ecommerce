package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

// Set this environment variable to run integration tests
const integrationTestFlag = "RUN_INTEGRATION_TESTS"

// Response is the API response format
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// TestProduct is a product model for testing
type TestProduct struct {
	ID            string   `json:"id,omitempty"`
	Name          string   `json:"name"`
	Slug          string   `json:"slug"`
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	CategoryID    string   `json:"category_id"`
	CategorySlug  string   `json:"category_slug"`
	StockQuantity int      `json:"stock_quantity"`
	Brand         string   `json:"brand,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

// Skip integrations tests unless flag is set
func skipUnlessIntegration(t *testing.T) {
	if os.Getenv(integrationTestFlag) == "" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS environment variable to run.")
	}
}

// Base URL for API tests
const baseURL = "http://localhost:8082/api/v1"

func TestHealthCheckIntegration(t *testing.T) {
	skipUnlessIntegration(t)

	resp, err := http.Get("http://localhost:8082/health")
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response Response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response.Status != http.StatusOK {
		t.Errorf("Expected response status %d, got %d", http.StatusOK, response.Status)
	}
}

func TestProductCRUD(t *testing.T) {
	skipUnlessIntegration(t)

	// Create a unique product for testing
	timestamp := time.Now().Unix()
	testProduct := TestProduct{
		Name:          fmt.Sprintf("Test Product %d", timestamp),
		Slug:          fmt.Sprintf("test-product-%d", timestamp),
		Description:   "This is a test product",
		Price:         99.99,
		CategoryID:    "some-category-id", // Make sure this exists in your database
		CategorySlug:  "electronics",
		StockQuantity: 100,
		Brand:         "Test Brand",
		Tags:          []string{"test", "integration"},
	}

	// 1. Create product
	productID := createProduct(t, testProduct)

	// 2. Get product by ID
	getProductByID(t, productID)

	// 3. Update product
	testProduct.ID = productID
	testProduct.Description = "Updated description"
	updateProduct(t, testProduct)

	// 4. Delete product
	deleteProduct(t, productID)
}

func createProduct(t *testing.T, product TestProduct) string {
	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshaling product: %v", err)
	}

	resp, err := http.Post(baseURL+"/products", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Error creating product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", resp.Status)
	}

	var response Response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Extract product ID from response
	productData := response.Data.(map[string]interface{})
	return productData["id"].(string)
}

func getProductByID(t *testing.T, id string) {
	resp, err := http.Get(fmt.Sprintf("%s/products/%s", baseURL, id))
	if err != nil {
		t.Fatalf("Error getting product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response Response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response.Status != http.StatusOK {
		t.Errorf("Expected response status %d, got %d", http.StatusOK, response.Status)
	}
}

func updateProduct(t *testing.T, product TestProduct) {
	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Error marshaling product: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/products/%s", baseURL, product.ID), bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error updating product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response Response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response.Status != http.StatusOK {
		t.Errorf("Expected response status %d, got %d", http.StatusOK, response.Status)
	}
}

func deleteProduct(t *testing.T, id string) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/products/%s", baseURL, id), nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error deleting product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response Response
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response.Status != http.StatusOK {
		t.Errorf("Expected response status %d, got %d", http.StatusOK, response.Status)
	}
}

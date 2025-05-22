package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	productHandler "broker-service/internal/handlers/http/product"
	"broker-service/tests/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// init function adds some debugging information
func init() {
	fmt.Println("Starting test initialization")
}

// setupProductHandlerTest sets up a new product handler for testing
func setupProductHandlerTest() *productHandler.Config {
	fmt.Println("Setting up product handler test")

	// Create mock client
	mockProductServiceClient := mocks.SetupMockProductClient()
	fmt.Println("Created mock client")

	// Set up the mock for GetProductImageFile
	mocks.AddGetProductImageFileMock(mockProductServiceClient, "test.jpg")
	fmt.Println("Added mock for GetProductImageFile")

	// Create a real product client for the tests
	mockProductClient := mocks.NewProductClientMockAdapter(mockProductServiceClient)
	fmt.Println("Created product client adapter")

	handler := &productHandler.Config{
		ProductClient: mockProductClient,
	}
	fmt.Println("Configured handler")

	return handler
}

// TestMain is used for test setup and teardown
func TestMain(m *testing.M) {
	fmt.Println("Test setup")

	// Run the tests
	code := m.Run()

	fmt.Println("Test teardown")

	// Exit with the same code as the test
	os.Exit(code)
}

// TestGetProductImage tests the GetProductImage handler
func TestGetProductImage(t *testing.T) {
	fmt.Println("Running TestGetProductImage")

	// Setup
	handler := setupProductHandlerTest()
	fmt.Println("Handler setup complete")

	// Create test request
	req, err := http.NewRequest("GET", "/api/v1/products/images/test.jpg", nil)
	assert.NoError(t, err)
	fmt.Println("Created request")

	// Create a test Chi router and add the URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("filename", "test.jpg")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	fmt.Println("Added route context")

	// Create response recorder
	rec := httptest.NewRecorder()
	fmt.Println("Created response recorder")

	// Call the handler
	fmt.Println("Calling handler")
	handler.GetProductImage(rec, req)
	fmt.Println("Handler called")

	// Verify response
	fmt.Println("Verifying response")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "image/jpeg", rec.Header().Get("Content-Type"))
	assert.Equal(t, "public, max-age=86400", rec.Header().Get("Cache-Control"))
	fmt.Println("Verification complete")
}

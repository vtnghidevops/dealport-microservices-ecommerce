package unit

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"product-service/internal/domain"
	"product-service/internal/service"
	"product-service/tests/mocks"
)

// TestGetProductByID tests the GetProductByID function
func TestGetProductByID(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)

	// Create a sample product
	mockProduct := &domain.Product{
		ID:            123,
		Name:          "Test Product",
		Description:   "Test description",
		Slug:          "test-product",
		Price:         99.99,
		CategoryID:    1,
		CategorySlug:  "electronics",
		StockQuantity: 100,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Initialize the service
	productService := service.NewProductService(mockRepo)

	// Test case 1: Successfully get a product
	t.Run("Successfully get product", func(t *testing.T) {
		// Set expectation
		mockRepo.On("GetProductByID", 123).Return(mockProduct, nil).Once()

		// Execute
		product, err := productService.GetProductByID(123)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, mockProduct.ID, product.ID)
		assert.Equal(t, mockProduct.Name, product.Name)
		assert.Equal(t, mockProduct.Price, product.Price)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Product not found
	t.Run("Product not found", func(t *testing.T) {
		// Set expectation
		mockRepo.On("GetProductByID", 456).Return(nil, domain.ErrProductNotFound).Once()

		// Execute
		product, err := productService.GetProductByID(456)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, domain.ErrProductNotFound, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestCreateProduct tests the CreateProduct function
func TestCreateProduct(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo)

	// Test case 1: Successfully create a product
	t.Run("Successfully create product", func(t *testing.T) {
		// Create product request
		product := &domain.Product{
			Name:          "New Product",
			Description:   "New product description",
			Price:         149.99,
			CategoryID:    2,
			StockQuantity: 50,
		}

		// Set expectation
		mockRepo.On("CreateProduct", mock.AnythingOfType("*domain.Product")).Return(789, nil).Once()

		// Execute
		id, err := productService.CreateProduct(product)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 789, id)
		assert.NotEmpty(t, product.Slug) // Slug should be generated

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Validation error - empty name
	t.Run("Empty product name", func(t *testing.T) {
		// Create invalid product request
		product := &domain.Product{
			Name:          "", // Empty name
			Description:   "Invalid product",
			Price:         149.99,
			CategoryID:    2,
			StockQuantity: 50,
		}

		// Execute
		id, err := productService.CreateProduct(product)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.Equal(t, domain.ErrInvalidProduct, err)

		// Verify expectations - CreateProduct shouldn't be called
		mockRepo.AssertNotCalled(t, "CreateProduct")
	})

	// Test case 3: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Create product request
		product := &domain.Product{
			Name:          "Error Product",
			Description:   "This will cause an error",
			Price:         149.99,
			CategoryID:    2,
			StockQuantity: 50,
		}

		// Set expectation - repository returns an error
		mockRepo.On("CreateProduct", mock.AnythingOfType("*domain.Product")).Return(0, errors.New("database error")).Once()

		// Execute
		id, err := productService.CreateProduct(product)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdateProduct tests the UpdateProduct function
func TestUpdateProduct(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo)

	// Test case 1: Successfully update a product
	t.Run("Successfully update product", func(t *testing.T) {
		// Create product to update
		product := &domain.Product{
			ID:            123,
			Name:          "Updated Product",
			Description:   "Updated description",
			Price:         199.99,
			CategoryID:    3,
			StockQuantity: 75,
		}

		// Set expectation
		mockRepo.On("UpdateProduct", mock.AnythingOfType("*domain.Product")).Return(nil).Once()

		// Execute
		err := productService.UpdateProduct(product)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, product.Slug) // Slug should be generated

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Validation error - missing ID
	t.Run("Missing product ID", func(t *testing.T) {
		// Create invalid product request
		product := &domain.Product{
			ID:            0, // Missing ID
			Name:          "Invalid Update",
			Description:   "Invalid update request",
			Price:         199.99,
			CategoryID:    3,
			StockQuantity: 75,
		}

		// Execute
		err := productService.UpdateProduct(product)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidProduct, err)

		// Verify expectations - UpdateProduct shouldn't be called
		mockRepo.AssertNotCalled(t, "UpdateProduct")
	})

	// Test case 3: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Create product to update
		product := &domain.Product{
			ID:            456,
			Name:          "Error Update",
			Description:   "This will cause an error",
			Price:         199.99,
			CategoryID:    3,
			StockQuantity: 75,
		}

		// Set expectation - repository returns an error
		mockRepo.On("UpdateProduct", mock.AnythingOfType("*domain.Product")).Return(errors.New("database error")).Once()

		// Execute
		err := productService.UpdateProduct(product)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestDeleteProduct tests the DeleteProduct function
func TestDeleteProduct(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo)

	// Test case 1: Successfully delete a product
	t.Run("Successfully delete product", func(t *testing.T) {
		// Set expectation
		mockRepo.On("DeleteProduct", 123).Return(nil).Once()

		// Execute
		err := productService.DeleteProduct(123)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Product not found
	t.Run("Product not found", func(t *testing.T) {
		// Set expectation
		mockRepo.On("DeleteProduct", 456).Return(domain.ErrProductNotFound).Once()

		// Execute
		err := productService.DeleteProduct(456)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrProductNotFound, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestGetProductReviews tests the GetProductReviews function
func TestGetProductReviews(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)

	// Create sample reviews
	mockReviews := []*domain.ProductReview{
		{
			ID:        1,
			ProductID: 123,
			UserID:    "user1",
			UserName:  "John Doe",
			Rating:    4.5,
			Comment:   "Great product!",
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			ProductID: 123,
			UserID:    "user2",
			UserName:  "Jane Smith",
			Rating:    5.0,
			Comment:   "Excellent quality!",
			CreatedAt: time.Now(),
		},
	}

	// Initialize the service
	productService := service.NewProductService(mockRepo)

	// Test case 1: Successfully get product reviews
	t.Run("Successfully get product reviews", func(t *testing.T) {
		// Set expectation
		mockRepo.On("GetProductReviews", 123, 1, 10).Return(mockReviews, len(mockReviews), nil).Once()

		// Execute
		reviews, count, err := productService.GetProductReviews(123, 1, 10)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, len(mockReviews), count)
		assert.Len(t, reviews, 2)
		assert.Equal(t, mockReviews[0].UserName, reviews[0].UserName)
		assert.Equal(t, mockReviews[1].Rating, reviews[1].Rating)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Default page/pageSize values
	t.Run("Default pagination values", func(t *testing.T) {
		// Set expectation - negative values should be corrected to defaults
		mockRepo.On("GetProductReviews", 123, 1, 10).Return(mockReviews, len(mockReviews), nil).Once()

		// Execute with negative values
		reviews, count, err := productService.GetProductReviews(123, -1, -5)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, len(mockReviews), count)
		assert.Len(t, reviews, 2)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestGetProductBySlug tests the GetProductBySlug function
func TestGetProductBySlug(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)

	// Create a sample product
	mockProduct := &domain.Product{
		ID:            123,
		Name:          "Test Product",
		Description:   "Test description",
		Slug:          "test-product",
		Price:         99.99,
		CategoryID:    1,
		CategorySlug:  "electronics",
		StockQuantity: 100,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Initialize the service
	productService := service.NewProductService(mockRepo)

	// Test case 1: Successfully get a product by slug
	t.Run("Successfully get product by slug", func(t *testing.T) {
		// Set expectation
		mockRepo.On("GetProductBySlug", "test-product").Return(mockProduct, nil).Once()

		// Execute
		product, err := productService.GetProductBySlug("test-product")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, mockProduct.ID, product.ID)
		assert.Equal(t, mockProduct.Name, product.Name)
		assert.Equal(t, mockProduct.Slug, product.Slug)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Product not found by slug
	t.Run("Product not found by slug", func(t *testing.T) {
		// Set expectation
		mockRepo.On("GetProductBySlug", "non-existent-product").Return(nil, domain.ErrProductNotFound).Once()

		// Execute
		product, err := productService.GetProductBySlug("non-existent-product")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, domain.ErrProductNotFound, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestGetAllProducts tests the GetAllProducts function
func TestGetAllProducts(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)

	// Create sample products
	mockProducts := []*domain.Product{
		{
			ID:            1,
			Name:          "Product 1",
			Slug:          "product-1",
			Price:         99.99,
			CategoryID:    1,
			StockQuantity: 100,
		},
		{
			ID:            2,
			Name:          "Product 2",
			Slug:          "product-2",
			Price:         149.99,
			CategoryID:    2,
			StockQuantity: 50,
		},
	}

	// Initialize the service
	productService := service.NewProductService(mockRepo)

	// Test case 1: Successfully get all products
	t.Run("Successfully get all products", func(t *testing.T) {
		// Filters
		filters := map[string]string{"category": "electronics"}

		// Set expectation
		mockRepo.On("GetAllProducts", 1, 10, filters).Return(mockProducts, 2, nil).Once()

		// Execute
		products, count, err := productService.GetAllProducts(1, 10, filters)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, products)
		assert.Equal(t, 2, count)
		assert.Equal(t, 2, len(products))
		assert.Equal(t, mockProducts[0].ID, products[0].ID)
		assert.Equal(t, mockProducts[1].ID, products[1].ID)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Zero or negative page/pageSize
	t.Run("Zero or negative page/pageSize", func(t *testing.T) {
		// Filters
		filters := map[string]string{}

		// Set expectation - note the service should default to page 1, pageSize 10
		mockRepo.On("GetAllProducts", 1, 10, filters).Return(mockProducts, 2, nil).Once()

		// Execute with zero values
		products, count, err := productService.GetAllProducts(0, 0, filters)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, products)
		assert.Equal(t, 2, count)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Filters
		filters := map[string]string{}

		// Set expectation
		mockRepo.On("GetAllProducts", 1, 10, filters).Return(nil, 0, errors.New("database error")).Once()

		// Execute
		products, count, err := productService.GetAllProducts(1, 10, filters)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Equal(t, 0, count)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestGetProductImages tests the GetProductImages function
func TestGetProductImages(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)

	// Create sample product images
	mockImages := []domain.ProductImage{
		{
			ID:           1,
			ProductID:    123,
			URL:          "http://example.com/image1.jpg",
			IsPrimary:    true,
			DisplayOrder: 1,
		},
		{
			ID:           2,
			ProductID:    123,
			URL:          "http://example.com/image2.jpg",
			IsPrimary:    false,
			DisplayOrder: 2,
		},
	}

	// Initialize the service
	productService := service.NewProductService(mockRepo)

	// Test case 1: Successfully get product images
	t.Run("Successfully get product images", func(t *testing.T) {
		// Set expectation
		mockRepo.On("GetProductImages", 123).Return(mockImages, nil).Once()

		// Execute
		images, err := productService.GetProductImages(123)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, images)
		assert.Equal(t, 2, len(images))
		assert.Equal(t, mockImages[0].ID, images[0].ID)
		assert.Equal(t, mockImages[1].ID, images[1].ID)
		assert.True(t, images[0].IsPrimary)
		assert.False(t, images[1].IsPrimary)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: No images found
	t.Run("No images found", func(t *testing.T) {
		// Set expectation
		mockRepo.On("GetProductImages", 456).Return([]domain.ProductImage{}, nil).Once()

		// Execute
		images, err := productService.GetProductImages(456)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, images)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Set expectation
		mockRepo.On("GetProductImages", 789).Return(nil, errors.New("database error")).Once()

		// Execute
		images, err := productService.GetProductImages(789)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, images)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestUploadProductImage tests the UploadProductImage function
func TestUploadProductImage(t *testing.T) {
	// Skip this test as upload is handled by the broker-service
	t.Skip("Skipping as file upload is actually handled by the broker-service")

	// Setup temporary uploads directory for testing
	tempDir := "./uploads/products"
	defer func() {
		// Clean up after tests
		os.RemoveAll(tempDir)
	}()

	mockRepo := new(mocks.MockProductRepository)

	// Create a sample product
	mockProduct := &domain.Product{
		ID:            123,
		Name:          "Test Product",
		Description:   "Test description",
		Slug:          "test-product",
		Price:         99.99,
		CategoryID:    1,
		StockQuantity: 100,
	}

	// Create sample product images
	mockImages := []domain.ProductImage{
		{
			ID:           1,
			ProductID:    123,
			URL:          "http://example.com/image1.jpg",
			IsPrimary:    true,
			DisplayOrder: 1,
		},
	}

	// Sample file content
	fileContent := []byte("This is a test image file content")

	// Initialize the service
	productService := service.NewProductService(mockRepo)

	// Test case 1: Successfully upload product image
	t.Run("Successfully upload product image", func(t *testing.T) {
		// Create mock file upload
		mockFile := mocks.NewMockFileUpload("test-image.jpg", int64(len(fileContent)), fileContent)
		mockFile.On("Open").Return(nil, nil)
		mockFile.On("Filename").Return("test-image.jpg")
		mockFile.On("Size").Return(int64(len(fileContent))).Maybe()

		// Set expectations
		mockRepo.On("GetProductByID", 123).Return(mockProduct, nil).Once()
		mockRepo.On("GetProductImages", 123).Return(mockImages, nil).Once()
		mockRepo.On("AddProductImage", 123, mock.AnythingOfType("string"), false, 1).Return(2, nil).Once()

		// Execute
		url, err := productService.UploadProductImage(123, mockFile, false)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, url)
		assert.Contains(t, url, "/products/") // Check that URL contains expected path

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockFile.AssertExpectations(t)
	})

	// Test case 2: Product not found
	t.Run("Product not found", func(t *testing.T) {
		// Create mock file upload
		mockFile := mocks.NewMockFileUpload("test-image.jpg", int64(len(fileContent)), fileContent)

		// Set expectations
		mockRepo.On("GetProductByID", 456).Return(nil, domain.ErrProductNotFound).Once()

		// Execute
		url, err := productService.UploadProductImage(456, mockFile, false)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, url)
		assert.Equal(t, domain.ErrProductNotFound, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Error getting product images
	t.Run("Error getting product images", func(t *testing.T) {
		// Create mock file upload
		mockFile := mocks.NewMockFileUpload("test-image.jpg", int64(len(fileContent)), fileContent)
		mockFile.On("Open").Return(nil, nil)
		mockFile.On("Filename").Return("test-image.jpg")
		mockFile.On("Size").Return(int64(len(fileContent))).Maybe()

		// Set expectations
		mockRepo.On("GetProductByID", 123).Return(mockProduct, nil).Once()
		mockRepo.On("GetProductImages", 123).Return(nil, errors.New("database error")).Once()

		// Execute
		url, err := productService.UploadProductImage(123, mockFile, false)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, url)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestDeleteProductImage tests the DeleteProductImage function
func TestDeleteProductImage(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo)

	// Test case 1: Successfully delete product image
	t.Run("Successfully delete product image", func(t *testing.T) {
		// Set expectation
		mockRepo.On("DeleteProductImage", 1).Return(nil).Once()

		// Execute
		err := productService.DeleteProductImage(1)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Set expectation
		mockRepo.On("DeleteProductImage", 2).Return(errors.New("database error")).Once()

		// Execute
		err := productService.DeleteProductImage(2)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdateProductImageOrder tests the UpdateProductImageOrder function
func TestUpdateProductImageOrder(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo)

	// Test case 1: Successfully update product image order
	t.Run("Successfully update product image order", func(t *testing.T) {
		// Set expectation
		mockRepo.On("UpdateProductImageOrder", 1, 3).Return(nil).Once()

		// Execute
		err := productService.UpdateProductImageOrder(1, 3)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Set expectation
		mockRepo.On("UpdateProductImageOrder", 2, 5).Return(errors.New("database error")).Once()

		// Execute
		err := productService.UpdateProductImageOrder(2, 5)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestSetPrimaryProductImage tests the SetPrimaryProductImage function
func TestSetPrimaryProductImage(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo)

	// Test case 1: Successfully set primary product image
	t.Run("Successfully set primary product image", func(t *testing.T) {
		// Set expectation
		mockRepo.On("SetPrimaryProductImage", 123, 1).Return(nil).Once()

		// Execute
		err := productService.SetPrimaryProductImage(123, 1)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Set expectation
		mockRepo.On("SetPrimaryProductImage", 456, 2).Return(errors.New("database error")).Once()

		// Execute
		err := productService.SetPrimaryProductImage(456, 2)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

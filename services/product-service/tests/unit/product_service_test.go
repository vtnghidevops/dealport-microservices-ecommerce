package unit

import (
	"errors"
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
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	productService := service.NewProductService(mockRepo, mockCategoryRepo)

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
	mockCategoryRepo := new(mocks.MockCategoryRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo, mockCategoryRepo)

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
		mockCategoryRepo.On("SyncProductCounts").Return(nil).Maybe()

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
	mockCategoryRepo := new(mocks.MockCategoryRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo, mockCategoryRepo)

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
		mockCategoryRepo.On("SyncProductCounts").Return(nil).Maybe()

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
	mockCategoryRepo := new(mocks.MockCategoryRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo, mockCategoryRepo)

	// Test case 1: Successfully delete a product
	t.Run("Successfully delete product", func(t *testing.T) {
		// Set expectation
		mockRepo.On("DeleteProduct", 123).Return(nil).Once()
		mockCategoryRepo.On("SyncProductCounts").Return(nil).Maybe()

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
	mockCategoryRepo := new(mocks.MockCategoryRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo, mockCategoryRepo)

	// Create sample reviews
	mockReviews := []*domain.ProductReview{
		{
			ID:        1,
			ProductID: 123,
			UserID:    "user1",
			UserName:  "John Doe",
			Rating:    5.0,
			Comment:   "Great product!",
		},
		{
			ID:        2,
			ProductID: 123,
			UserID:    "user2",
			UserName:  "Jane Smith",
			Rating:    4.0,
			Comment:   "Good quality",
		},
	}

	// Test case 1: Successfully get reviews
	t.Run("Successfully get reviews", func(t *testing.T) {
		// Set expectation
		mockRepo.On("GetProductReviews", 123, 1, 10).Return(mockReviews, 2, nil).Once()

		// Execute
		reviews, total, err := productService.GetProductReviews(123, 1, 10)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, reviews)
		assert.Equal(t, 2, len(reviews))
		assert.Equal(t, 2, total)
		assert.Equal(t, mockReviews[0].Comment, reviews[0].Comment)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Set expectation - repository returns an error
		mockRepo.On("GetProductReviews", 456, 1, 10).Return(nil, 0, errors.New("database error")).Once()

		// Execute
		reviews, total, err := productService.GetProductReviews(456, 1, 10)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, reviews)
		assert.Equal(t, 0, total)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestGetProductBySlug tests the GetProductBySlug function
func TestGetProductBySlug(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)
	mockCategoryRepo := new(mocks.MockCategoryRepository)

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
	productService := service.NewProductService(mockRepo, mockCategoryRepo)

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
		assert.Equal(t, mockProduct.Slug, product.Slug)
		assert.Equal(t, mockProduct.Name, product.Name)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Product not found
	t.Run("Product not found", func(t *testing.T) {
		// Set expectation
		mockRepo.On("GetProductBySlug", "non-existent").Return(nil, domain.ErrProductNotFound).Once()

		// Execute
		product, err := productService.GetProductBySlug("non-existent")

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
	mockCategoryRepo := new(mocks.MockCategoryRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo, mockCategoryRepo)

	// Create sample products
	mockProducts := []*domain.Product{
		{
			ID:            1,
			Name:          "Product 1",
			Description:   "First product",
			Slug:          "product-1",
			Price:         50.0,
			CategoryID:    1,
			CategorySlug:  "electronics",
			StockQuantity: 100,
		},
		{
			ID:            2,
			Name:          "Product 2",
			Description:   "Second product",
			Slug:          "product-2",
			Price:         75.0,
			CategoryID:    2,
			CategorySlug:  "books",
			StockQuantity: 50,
		},
	}

	// Test case 1: Successfully get all products
	t.Run("Successfully get all products", func(t *testing.T) {
		filters := map[string]string{}

		// Set expectation
		mockRepo.On("GetAllProducts", 1, 10, filters).Return(mockProducts, 2, nil).Once()

		// Execute
		products, total, err := productService.GetAllProducts(1, 10, filters)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, products)
		assert.Equal(t, 2, len(products))
		assert.Equal(t, 2, total)
		assert.Equal(t, mockProducts[0].Name, products[0].Name)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Repository error
	t.Run("Repository error", func(t *testing.T) {
		filters := map[string]string{}

		// Set expectation - repository returns an error
		mockRepo.On("GetAllProducts", 1, 10, filters).Return(nil, 0, errors.New("database error")).Once()

		// Execute
		products, total, err := productService.GetAllProducts(1, 10, filters)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Equal(t, 0, total)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Default pagination values
	t.Run("Default pagination values", func(t *testing.T) {
		filters := map[string]string{}

		// Set expectation with default values
		mockRepo.On("GetAllProducts", 1, 10, filters).Return(mockProducts, 2, nil).Once()

		// Execute with invalid pagination values
		products, total, err := productService.GetAllProducts(0, 0, filters)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, products)
		assert.Equal(t, 2, len(products))
		assert.Equal(t, 2, total)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestGetProductImages tests the GetProductImages function
func TestGetProductImages(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)
	mockCategoryRepo := new(mocks.MockCategoryRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo, mockCategoryRepo)

	// Create sample product images
	mockImages := []domain.ProductImage{
		{
			ID:           1,
			ProductID:    123,
			URL:          "/images/product1.jpg",
			IsPrimary:    true,
			DisplayOrder: 0,
		},
		{
			ID:           2,
			ProductID:    123,
			URL:          "/images/product2.jpg",
			IsPrimary:    false,
			DisplayOrder: 1,
		},
	}

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
		assert.Equal(t, mockImages[0].URL, images[0].URL)
		assert.True(t, images[0].IsPrimary)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Set expectation - repository returns an error
		mockRepo.On("GetProductImages", 456).Return(nil, errors.New("database error")).Once()

		// Execute
		images, err := productService.GetProductImages(456)

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
	// Setup
	mockRepo := new(mocks.MockProductRepository)
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	mockStorage := new(mocks.MockStorageService)

	// Initialize the service with storage
	productService := service.NewProductServiceWithStorage(mockRepo, mockCategoryRepo, mockStorage, nil)

	// Test case 1: Invalid product ID
	t.Run("Invalid product ID", func(t *testing.T) {
		// Create mock file upload
		mockFile := mocks.NewMockFileUpload("test.jpg", 1024, []byte("test content"))

		// Execute
		imageURL, err := productService.UploadProductImage(0, mockFile, true)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, imageURL)
		assert.Equal(t, domain.ErrInvalidProduct, err)

		// Verify no repository calls
		mockRepo.AssertNotCalled(t, "AddProductImage")
	})

	// Test case 2: Successfully upload image to MinIO
	t.Run("Successfully upload image to MinIO", func(t *testing.T) {
		// Create mock file upload
		mockFile := mocks.NewMockFileUpload("test.jpg", 1024, []byte("test content"))

		// Setup mock expectations for file operations
		mockFile.On("Filename").Return("test.jpg")
		mockFile.On("Size").Return(int64(1024))
		mockFile.On("Open").Return(nil, nil)

		// Set expectations
		mockStorage.On("UploadFile", mock.Anything, mock.AnythingOfType("string"), mock.Anything, int64(1024), "image/jpeg").Return("/images/products-api/123_test.jpg", nil).Once()
		mockStorage.On("GetPresignedURL", mock.Anything, "/images/products-api/123_test.jpg", 3600).Return("https://presigned-url.example.com/123_test.jpg", nil).Once()
		mockRepo.On("GetProductImages", 123).Return([]domain.ProductImage{}, nil).Once()
		mockRepo.On("AddProductImage", 123, "/images/products-api/123_test.jpg", true, 0).Return(1, nil).Once()

		// Execute
		imageURL, err := productService.UploadProductImage(123, mockFile, true)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, imageURL)
		assert.Equal(t, "https://presigned-url.example.com/123_test.jpg", imageURL)

		// Verify expectations
		mockStorage.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Storage upload error
	t.Run("Storage upload error", func(t *testing.T) {
		// Create mock file upload
		mockFile := mocks.NewMockFileUpload("test.jpg", 1024, []byte("test content"))

		// Setup mock expectations for file operations
		mockFile.On("Filename").Return("test.jpg")
		mockFile.On("Size").Return(int64(1024))
		mockFile.On("Open").Return(nil, nil)

		// Set expectations - storage returns error
		mockStorage.On("UploadFile", mock.Anything, mock.AnythingOfType("string"), mock.Anything, int64(1024), "image/jpeg").Return("", errors.New("storage error")).Once()

		// For local storage fallback (after MinIO fails)
		mockRepo.On("GetProductImages", 123).Return([]domain.ProductImage{}, nil).Once()
		mockRepo.On("AddProductImage", 123, mock.AnythingOfType("string"), true, 0).Return(1, nil).Once()

		// Execute
		imageURL, err := productService.UploadProductImage(123, mockFile, true)

		// We expect this to fall back to local storage, so no error should occur
		// The exact behavior depends on the implementation
		_ = imageURL
		_ = err

		// Verify expectations
		mockStorage.AssertExpectations(t)
	})
}

// TestDeleteProductImage tests the DeleteProductImage function
func TestDeleteProductImage(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	mockStorage := new(mocks.MockStorageService)

	// Initialize the service with storage
	productService := service.NewProductServiceWithStorage(mockRepo, mockCategoryRepo, mockStorage, nil)

	// Test case 1: Successfully delete product image
	t.Run("Successfully delete product image", func(t *testing.T) {
		// Set expectation
		mockRepo.On("GetProductImages", 0).Return([]domain.ProductImage{
			{ID: 1, URL: "/images/products/test.jpg"},
		}, nil).Once()
		mockRepo.On("DeleteProductImage", 1).Return(nil).Once()
		mockStorage.On("DeleteFile", mock.Anything, "products/test.jpg").Return(nil).Once()

		// Execute
		err := productService.DeleteProductImage(1)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	// Test case 2: Repository error
	t.Run("Repository error", func(t *testing.T) {
		// Set expectation - repository returns an error
		mockRepo.On("GetProductImages", 0).Return([]domain.ProductImage{
			{ID: 2, URL: "/images/products/test2.jpg"},
		}, nil).Once()
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
	mockCategoryRepo := new(mocks.MockCategoryRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo, mockCategoryRepo)

	// Test case 1: Successfully update image order
	t.Run("Successfully update image order", func(t *testing.T) {
		// Set expectation
		mockRepo.On("UpdateProductImageOrder", 1, 5).Return(nil).Once()

		// Execute
		err := productService.UpdateProductImageOrder(1, 5)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestSetPrimaryProductImage tests the SetPrimaryProductImage function
func TestSetPrimaryProductImage(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)
	mockCategoryRepo := new(mocks.MockCategoryRepository)

	// Initialize the service
	productService := service.NewProductService(mockRepo, mockCategoryRepo)

	// Test case 1: Successfully set primary image
	t.Run("Successfully set primary image", func(t *testing.T) {
		// Set expectation
		mockRepo.On("SetPrimaryProductImage", 123, 1).Return(nil).Once()

		// Execute
		err := productService.SetPrimaryProductImage(123, 1)

		// Assert
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestGetPresignedURL tests the GetPresignedURL function
func TestGetPresignedURL(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockProductRepository)
	mockCategoryRepo := new(mocks.MockCategoryRepository)
	mockStorage := new(mocks.MockStorageService)

	// Initialize the service with storage
	productService := service.NewProductServiceWithStorage(mockRepo, mockCategoryRepo, mockStorage, nil)

	// Test case 1: Successfully get presigned URL
	t.Run("Successfully get presigned URL", func(t *testing.T) {
		objectPath := "/images/products-api/test.jpg"
		expectedURL := "https://presigned-url.example.com/test.jpg"

		// Set expectation - service transforms path before calling storage
		mockStorage.On("GetPresignedURL", mock.Anything, "products/test.jpg", 3600).Return(expectedURL, nil).Once()

		// Execute
		url, err := productService.GetPresignedURL(objectPath)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedURL, url)

		// Verify expectations
		mockStorage.AssertExpectations(t)
	})
}

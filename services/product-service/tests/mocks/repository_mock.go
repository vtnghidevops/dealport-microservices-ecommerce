package mocks

import (
	"product-service/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock of ProductRepository interface
type MockProductRepository struct {
	mock.Mock
}

// GetProductByID mocks the GetProductByID method
func (m *MockProductRepository) GetProductByID(id int) (*domain.Product, error) {
	args := m.Called(id)

	var product *domain.Product
	if args.Get(0) != nil {
		product = args.Get(0).(*domain.Product)
	}

	return product, args.Error(1)
}

// GetProductBySlug mocks the GetProductBySlug method
func (m *MockProductRepository) GetProductBySlug(slug string) (*domain.Product, error) {
	args := m.Called(slug)

	var product *domain.Product
	if args.Get(0) != nil {
		product = args.Get(0).(*domain.Product)
	}

	return product, args.Error(1)
}

// GetAllProducts mocks the GetAllProducts method
func (m *MockProductRepository) GetAllProducts(page, pageSize int, filters map[string]string) ([]*domain.Product, int, error) {
	args := m.Called(page, pageSize, filters)

	var products []*domain.Product
	if args.Get(0) != nil {
		products = args.Get(0).([]*domain.Product)
	}

	return products, args.Int(1), args.Error(2)
}

// CreateProduct mocks the CreateProduct method
func (m *MockProductRepository) CreateProduct(product *domain.Product) (int, error) {
	args := m.Called(product)
	return args.Int(0), args.Error(1)
}

// UpdateProduct mocks the UpdateProduct method
func (m *MockProductRepository) UpdateProduct(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

// DeleteProduct mocks the DeleteProduct method
func (m *MockProductRepository) DeleteProduct(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetProductImages mocks the GetProductImages method
func (m *MockProductRepository) GetProductImages(productID int) ([]domain.ProductImage, error) {
	args := m.Called(productID)

	var images []domain.ProductImage
	if args.Get(0) != nil {
		images = args.Get(0).([]domain.ProductImage)
	}

	return images, args.Error(1)
}

// AddProductImage mocks the AddProductImage method
func (m *MockProductRepository) AddProductImage(productID int, imageURL string, isPrimary bool, displayOrder int) (int, error) {
	args := m.Called(productID, imageURL, isPrimary, displayOrder)
	return args.Int(0), args.Error(1)
}

// DeleteProductImage mocks the DeleteProductImage method
func (m *MockProductRepository) DeleteProductImage(imageID int) error {
	args := m.Called(imageID)
	return args.Error(0)
}

// UpdateProductImageOrder mocks the UpdateProductImageOrder method
func (m *MockProductRepository) UpdateProductImageOrder(imageID int, displayOrder int) error {
	args := m.Called(imageID, displayOrder)
	return args.Error(0)
}

// SetPrimaryProductImage mocks the SetPrimaryProductImage method
func (m *MockProductRepository) SetPrimaryProductImage(productID int, imageID int) error {
	args := m.Called(productID, imageID)
	return args.Error(0)
}

// GetProductTags mocks the GetProductTags method
func (m *MockProductRepository) GetProductTags(productID int) ([]string, error) {
	args := m.Called(productID)

	var tags []string
	if args.Get(0) != nil {
		tags = args.Get(0).([]string)
	}

	return tags, args.Error(1)
}

// GetProductReviews mocks the GetProductReviews method
func (m *MockProductRepository) GetProductReviews(productID int, page, pageSize int) ([]*domain.ProductReview, int, error) {
	args := m.Called(productID, page, pageSize)

	var reviews []*domain.ProductReview
	if args.Get(0) != nil {
		reviews = args.Get(0).([]*domain.ProductReview)
	}

	return reviews, args.Int(1), args.Error(2)
}

// AddProductReview mocks the AddProductReview method
func (m *MockProductRepository) AddProductReview(review *domain.ProductReview) (int, error) {
	args := m.Called(review)
	return args.Int(0), args.Error(1)
}

// UpdateProductReview mocks the UpdateProductReview method
func (m *MockProductRepository) UpdateProductReview(review *domain.ProductReview) error {
	args := m.Called(review)
	return args.Error(0)
}

// DeleteProductReview mocks the DeleteProductReview method
func (m *MockProductRepository) DeleteProductReview(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetRandomTopRatedReviews mocks the GetRandomTopRatedReviews method
func (m *MockProductRepository) GetRandomTopRatedReviews(limit int) ([]*domain.Testimonial, error) {
	args := m.Called(limit)

	var testimonials []*domain.Testimonial
	if args.Get(0) != nil {
		testimonials = args.Get(0).([]*domain.Testimonial)
	}

	return testimonials, args.Error(1)
}

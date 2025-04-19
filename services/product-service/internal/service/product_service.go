package service

import (
	"product-service/internal/domain"
)

// ProductService implements the domain.ProductService interface
type ProductService struct {
	repo domain.ProductRepository
}

// NewProductService creates a new instance of ProductService
func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// GetByID returns a product by ID
func (s *ProductService) GetByID(id string) (*domain.Product, error) {
	return s.repo.GetByID(id)
}

// GetBySlug returns a product by slug
func (s *ProductService) GetBySlug(slug string) (*domain.Product, error) {
	return s.repo.GetBySlug(slug)
}

// List returns products with optional filtering and pagination
func (s *ProductService) List(page, pageSize int, filters map[string]string) ([]*domain.Product, int, error) {
	// Default values
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	return s.repo.List(page, pageSize, filters)
}

// Create adds a new product
func (s *ProductService) Create(product *domain.Product) (string, error) {
	return s.repo.Create(product)
}

// Update updates an existing product
func (s *ProductService) Update(product *domain.Product) error {
	return s.repo.Update(product)
}

// Delete removes a product
func (s *ProductService) Delete(id string) error {
	return s.repo.Delete(id)
}

// GetProductReviews returns reviews for a product
func (s *ProductService) GetProductReviews(productID string, page, pageSize int) ([]*domain.ProductReview, int, error) {
	// Default values
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	return s.repo.GetProductReviews(productID, page, pageSize)
}

// AddProductReview adds a review for a product
func (s *ProductService) AddProductReview(review *domain.ProductReview) (string, error) {
	return s.repo.AddProductReview(review)
}

// DeleteProductReview deletes a review
func (s *ProductService) DeleteProductReview(id string) error {
	return s.repo.DeleteProductReview(id)
}

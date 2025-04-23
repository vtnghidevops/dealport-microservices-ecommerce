package service

import (
	"product-service/internal/domain"
	"product-service/internal/util"
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

// GetProductByID returns a product by ID
func (s *ProductService) GetProductByID(id int) (*domain.Product, error) {
	return s.repo.GetProductByID(id)
}

// GetProductBySlug returns a product by slug
func (s *ProductService) GetProductBySlug(slug string) (*domain.Product, error) {
	return s.repo.GetProductBySlug(slug)
}

// GetAllProducts returns products with optional filtering and pagination
func (s *ProductService) GetAllProducts(page, pageSize int, filters map[string]string) ([]*domain.Product, int, error) {
	// Default values
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	return s.repo.GetAllProducts(page, pageSize, filters)
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(product *domain.Product) (int, error) {
	// Validate product fields here
	if product.Name == "" {
		return 0, domain.ErrInvalidProduct
	}

	// Create slug from name if not provided
	if product.Slug == "" {
		product.Slug = util.CreateSlug(product.Name)
	}

	return s.repo.CreateProduct(product)
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(product *domain.Product) error {
	// Validate product fields here
	if product.ID == 0 {
		return domain.ErrInvalidProduct
	}

	if product.Name == "" {
		return domain.ErrInvalidProduct
	}

	// Create slug from name if not provided
	if product.Slug == "" {
		product.Slug = util.CreateSlug(product.Name)
	}

	return s.repo.UpdateProduct(product)
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}

// GetProductImages returns all images for a product
func (s *ProductService) GetProductImages(productID int) ([]domain.ProductImage, error) {
	return s.repo.GetProductImages(productID)
}

// GetProductReviews returns reviews for a product with pagination
func (s *ProductService) GetProductReviews(productID int, page, pageSize int) ([]*domain.ProductReview, int, error) {
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
func (s *ProductService) AddProductReview(review *domain.ProductReview) (int, error) {
	return s.repo.AddProductReview(review)
}

// UpdateProductReview updates an existing review
func (s *ProductService) UpdateProductReview(review *domain.ProductReview) error {
	return s.repo.UpdateProductReview(review)
}

// DeleteProductReview deletes a review by ID
func (s *ProductService) DeleteProductReview(id int) error {
	return s.repo.DeleteProductReview(id)
}

// GetRandomTopRatedReviews returns random top-rated reviews for HappyCustomers section
func (s *ProductService) GetRandomTopRatedReviews(limit int) ([]*domain.Testimonial, error) {
	// Default value for limit
	if limit <= 0 {
		limit = 7 // Default number of testimonials to show
	}

	return s.repo.GetRandomTopRatedReviews(limit)
}

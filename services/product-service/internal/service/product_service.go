package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"product-service/internal/domain"
	"product-service/internal/util"
	"time"
)

// ProductService implements the domain.ProductService interface
type ProductService struct {
	productRepo domain.ProductRepository
}

// NewProductService creates a new instance of ProductService
func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: repo,
	}
}

// GetProductByID returns a product by ID
func (s *ProductService) GetProductByID(id int) (*domain.Product, error) {
	return s.productRepo.GetProductByID(id)
}

// GetProductBySlug returns a product by slug
func (s *ProductService) GetProductBySlug(slug string) (*domain.Product, error) {
	return s.productRepo.GetProductBySlug(slug)
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

	return s.productRepo.GetAllProducts(page, pageSize, filters)
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

	return s.productRepo.CreateProduct(product)
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

	return s.productRepo.UpdateProduct(product)
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(id int) error {
	return s.productRepo.DeleteProduct(id)
}

// GetProductImages returns all images for a product
func (s *ProductService) GetProductImages(productID int) ([]domain.ProductImage, error) {
	return s.productRepo.GetProductImages(productID)
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

	return s.productRepo.GetProductReviews(productID, page, pageSize)
}

// AddProductReview adds a review for a product
func (s *ProductService) AddProductReview(review *domain.ProductReview) (int, error) {
	return s.productRepo.AddProductReview(review)
}

// UpdateProductReview updates an existing review
func (s *ProductService) UpdateProductReview(review *domain.ProductReview) error {
	return s.productRepo.UpdateProductReview(review)
}

// DeleteProductReview deletes a review by ID
func (s *ProductService) DeleteProductReview(id int) error {
	return s.productRepo.DeleteProductReview(id)
}

// GetRandomTopRatedReviews returns random top-rated reviews for HappyCustomers section
func (s *ProductService) GetRandomTopRatedReviews(limit int) ([]*domain.Testimonial, error) {
	// Default value for limit
	if limit <= 0 {
		limit = 7 // Default number of testimonials to show
	}

	return s.productRepo.GetRandomTopRatedReviews(limit)
}

// UploadProductImage uploads a product image and returns the URL
func (s *ProductService) UploadProductImage(productID int, file *multipart.FileHeader, isPrimary bool) (string, error) {
	// Check if product exists
	product, err := s.GetProductByID(productID)
	if err != nil {
		return "", err
	}
	if product == nil {
		return "", domain.ErrProductNotFound
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "./uploads/products"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create uploads directory: %w", err)
	}

	// Generate unique filename
	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d_%d%s", productID, time.Now().UnixNano(), fileExt)
	filePath := filepath.Join(uploadsDir, fileName)

	// Save the file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// Get current display order
	images, err := s.GetProductImages(productID)
	if err != nil {
		return "", err
	}

	// Determine display order (last position by default)
	displayOrder := 0
	if len(images) > 0 {
		displayOrder = len(images)
	}

	// Create relative URL path for database
	imageURL := fmt.Sprintf("/api/products/images/%s", fileName)

	// If this is set as primary and there are existing images,
	// we need to update other images to non-primary
	if isPrimary && len(images) > 0 {
		// Add the image first
		imageID, err := s.productRepo.AddProductImage(productID, imageURL, isPrimary, displayOrder)
		if err != nil {
			// Try to clean up the file
			os.Remove(filePath)
			return "", err
		}

		// Then set it as primary (which will handle updating other images)
		if err := s.productRepo.SetPrimaryProductImage(productID, imageID); err != nil {
			return "", err
		}
	} else {
		// Regular insert
		_, err = s.productRepo.AddProductImage(productID, imageURL, isPrimary, displayOrder)
		if err != nil {
			// Try to clean up the file
			os.Remove(filePath)
			return "", err
		}
	}

	return imageURL, nil
}

// DeleteProductImage deletes a product image by ID
func (s *ProductService) DeleteProductImage(imageID int) error {
	// Here we would need to:
	// 1. Get the image to find its file path
	// 2. Delete the file from storage
	// 3. Delete the database record

	// This is simplified as we don't have an API to get a single image by ID yet
	// In a complete implementation, we would get the image URL, parse the filename,
	// and delete the corresponding file before removing the database record

	return s.productRepo.DeleteProductImage(imageID)
}

// UpdateProductImageOrder updates the display order of a product image
func (s *ProductService) UpdateProductImageOrder(imageID int, displayOrder int) error {
	return s.productRepo.UpdateProductImageOrder(imageID, displayOrder)
}

// SetPrimaryProductImage sets the primary image for a product
func (s *ProductService) SetPrimaryProductImage(productID int, imageID int) error {
	return s.productRepo.SetPrimaryProductImage(productID, imageID)
}



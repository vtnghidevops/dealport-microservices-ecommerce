package service

import (
	"fmt"
	"io"
	"log"
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

	// Call repository layer
	err := s.productRepo.UpdateProduct(product)
	if err != nil {
		return err
	}

	return nil
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
func (s *ProductService) UploadProductImage(productID int, file domain.FileUpload, isPrimary bool) (string, error) {
	// Log the start of the upload process
	log.Printf("Starting image upload for product ID %d, isPrimary: %v", productID, isPrimary)

	// Check if product exists
	product, err := s.GetProductByID(productID)
	if err != nil {
		log.Printf("ERROR: Failed to find product ID %d: %v", productID, err)
		return "", err
	}
	if product == nil {
		log.Printf("ERROR: Product ID %d not found", productID)
		return "", domain.ErrProductNotFound
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "./uploads/products"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Printf("ERROR: Failed to create uploads directory: %v", err)
		return "", fmt.Errorf("failed to create uploads directory: %w", err)
	}

	// Generate unique filename
	fileExt := filepath.Ext(file.Filename())
	fileName := fmt.Sprintf("%d_%d%s", productID, time.Now().UnixNano(), fileExt)
	filePath := filepath.Join(uploadsDir, fileName)

	log.Printf("Generated filename: %s", fileName)

	// Save the file
	src, err := file.Open()
	if err != nil {
		log.Printf("ERROR: Failed to open uploaded file: %v", err)
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	log.Printf("Successfully opened source file for reading")

	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("ERROR: Failed to create destination file: %v", err)
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	log.Printf("Successfully created destination file for writing")

	written, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("ERROR: Failed to copy file: %v", err)
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	log.Printf("Successfully copied %d bytes to destination file", written)

	// Get current display order
	images, err := s.GetProductImages(productID)
	if err != nil {
		log.Printf("ERROR: Failed to get product images: %v", err)
		return "", err
	}

	// Determine display order (last position by default)
	displayOrder := 0
	if len(images) > 0 {
		displayOrder = len(images)
	}

	// Create relative URL path for database
	imageURL := fmt.Sprintf("/api/products/images/%s", fileName)
	fmt.Printf("Image URL: %s, display order: %d", imageURL, displayOrder)

	// Create fully qualified URL for response - make sure we provide full URL
	// This is what will be returned to the client after upload

	imgBaseURL := os.Getenv("ECOMMERCE_IMG_URL")
	if imgBaseURL == "" {
		// Kiểm tra môi trường để quyết định URL mặc định
		_, isLocalDev := os.LookupEnv("LOCAL_DEV")
		if isLocalDev {
			// Đang ở môi trường phát triển cục bộ
			imgBaseURL = "http://localhost:58082" // Sử dụng cổng local của product-service
		} else {
			// Môi trường sản xuất hoặc staging
			imgBaseURL = "https://api.deploy.io.vn" // Fallback nếu không có biến môi trường
		}
	}
	fullImageURL := fmt.Sprintf("%s/api/products/images/%s", imgBaseURL, fileName)
	fmt.Printf("Full image URL for response: %s", fullImageURL)

	// If this is set as primary and there are existing images,
	// we need to update other images to non-primary
	if isPrimary && len(images) > 0 {
		// Add the image first
		imageID, err := s.productRepo.AddProductImage(productID, imageURL, isPrimary, displayOrder)
		if err != nil {
			// Try to clean up the file
			os.Remove(filePath)
			log.Printf("ERROR: Failed to add product image to database: %v", err)
			return "", err
		}

		log.Printf("Added image to database with ID %d", imageID)

		// Then set it as primary (which will handle updating other images)
		if err := s.productRepo.SetPrimaryProductImage(productID, imageID); err != nil {
			log.Printf("ERROR: Failed to set image as primary: %v", err)
			return "", err
		}

		log.Printf("Successfully set image as primary")
	} else {
		// Regular insert
		imageID, err := s.productRepo.AddProductImage(productID, imageURL, isPrimary, displayOrder)
		if err != nil {
			// Try to clean up the file
			os.Remove(filePath)
			log.Printf("ERROR: Failed to add product image to database: %v", err)
			return "", err
		}

		log.Printf("Added image to database with ID %d", imageID)
	}

	log.Printf("Image upload completed successfully")
	return fullImageURL, nil
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

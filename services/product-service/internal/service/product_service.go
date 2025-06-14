package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"product-service/internal/cache"
	"product-service/internal/domain"
	"product-service/internal/storage"
	"product-service/internal/util"
)

// ProductService implements the domain.ProductService interface
type ProductService struct {
	productRepo    domain.ProductRepository
	storageService storage.StorageService
	urlCache       cache.ImageURLCache
}

// NewProductService creates a new ProductService with local file storage
func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: repo,
	}
}

// NewProductServiceWithStorage creates a new ProductService with the specified storage service
func NewProductServiceWithStorage(repo domain.ProductRepository, storageService storage.StorageService, urlCache cache.ImageURLCache) *ProductService {
	return &ProductService{
		productRepo:    repo,
		storageService: storageService,
		urlCache:       urlCache,
	}
}

// transformProductURLs converts MinIO URLs to presigned URLs for a product
func (s *ProductService) transformProductURLs(ctx context.Context, product *domain.Product) {
	if s.storageService == nil {
		return
	}

	log.Printf("Transforming URLs for product ID %d", product.ID)

	// For proper URL transformation
	if product.ImageURL != "" {
		log.Printf("Original main image URL: %s", product.ImageURL)

		// Handle both URL patterns
		isMinioURL := strings.Contains(product.ImageURL, "/images/products-api/") ||
			strings.Contains(product.ImageURL, "minioapi.deploy.io.vn")

		if isMinioURL {
			if presignedURL, err := s.storageService.GetPresignedURL(ctx, product.ImageURL, 3600); err == nil {
				log.Printf("Transformed main image URL to presigned URL")
				product.ImageURL = presignedURL
			} else {
				log.Printf("Error transforming main image URL: %v", err)
			}
		} else {
			log.Printf("Main image URL is not a MinIO URL, keeping as is")
		}
	}

	// Transform image slider URLs
	for i, url := range product.ImgSlider {
		if url != "" {
			isMinioURL := strings.Contains(url, "/images/products-api/") ||
				strings.Contains(url, "minioapi.deploy.io.vn")

			if isMinioURL {
				if presignedURL, err := s.storageService.GetPresignedURL(ctx, url, 3600); err == nil {
					product.ImgSlider[i] = presignedURL
					log.Printf("Transformed slider image URL at index %d", i)
				}
			}
		}
	}

	// Transform images array URLs
	for i := range product.Images {
		if product.Images[i].URL != "" {
			isMinioURL := strings.Contains(product.Images[i].URL, "/images/products-api/") ||
				strings.Contains(product.Images[i].URL, "minioapi.deploy.io.vn")

			if isMinioURL {
				if presignedURL, err := s.storageService.GetPresignedURL(ctx, product.Images[i].URL, 3600); err == nil {
					product.Images[i].URL = presignedURL
					log.Printf("Transformed product image URL for image ID %d", product.Images[i].ID)
				}
			}
		}
	}
}

// GetProductByID returns a product by ID
func (s *ProductService) GetProductByID(id int) (*domain.Product, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	// Transform any MinIO URLs
	if s.storageService != nil {
		s.transformProductURLs(context.Background(), product)
	}

	return product, nil
}

// GetProductBySlug returns a product by slug
func (s *ProductService) GetProductBySlug(slug string) (*domain.Product, error) {
	product, err := s.productRepo.GetProductBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Transform any MinIO URLs
	if s.storageService != nil {
		s.transformProductURLs(context.Background(), product)
	}

	return product, nil
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

	products, totalCount, err := s.productRepo.GetAllProducts(page, pageSize, filters)
	if err != nil {
		return nil, 0, err
	}

	// Transform any MinIO URLs if we have storage service
	if s.storageService != nil {
		ctx := context.Background()
		// Process each product to transform URLs
		for _, product := range products {
			s.transformProductURLs(ctx, product)
		}
	}

	return products, totalCount, nil
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

// transformImageURL converts a MinIO URL to a presigned URL
func (s *ProductService) transformImageURL(ctx context.Context, url string) string {
	if s.storageService == nil || url == "" {
		return url
	}

	log.Printf("Transforming image URL: %s", url)

	// Check if this is a MinIO URL by either pattern
	isMinioURL := strings.Contains(url, "/images/products-api/") ||
		strings.Contains(url, "minioapi.deploy.io.vn") ||
		strings.Contains(url, "X-Amz-Algorithm=AWS")

	if !isMinioURL {
		log.Printf("Image URL is not a MinIO URL, keeping as is")
		return url
	}

	// Generate presigned URL for MinIO object
	presignedURL, err := s.storageService.GetPresignedURL(ctx, url, 3600)
	if err != nil {
		log.Printf("Error generating presigned URL for %s: %v", url, err)
		return url
	}

	log.Printf("Transformed URL to presigned URL")
	return presignedURL
}

// GetProductImages returns all images for a product
func (s *ProductService) GetProductImages(productID int) ([]domain.ProductImage, error) {
	images, err := s.productRepo.GetProductImages(productID)
	if err != nil {
		return nil, err
	}

	// If we have MinIO storage, transform the URLs
	if s.storageService != nil {
		ctx := context.Background()

		// Transform each URL if it's a MinIO URL
		for i := range images {
			images[i].URL = s.transformImageURL(ctx, images[i].URL)
		}
	}

	return images, nil
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

// UploadProductImage adds an image to a product
func (s *ProductService) UploadProductImage(productID int, file domain.FileUpload, isPrimary bool) (string, error) {
	// Validate productID
	if productID <= 0 {
		return "", domain.ErrInvalidProduct
	}

	// Generate a unique filename based on product ID and timestamp
	fileName := fmt.Sprintf("%d_%d%s", productID, time.Now().UnixNano(), filepath.Ext(file.Filename()))

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		log.Printf("ERROR: Failed to open uploaded file: %v", err)
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	log.Printf("Successfully opened source file for reading")

	var imageURL string
	var fullImageURL string
	useLocalStorage := false

	// Check if we have MinIO storage available and use it by default
	if s.storageService != nil {
		log.Printf("Attempting to use MinIO storage for upload")
		// Use MinIO storage
		// Store in products/ folder in MinIO but will be accessible as /images/products-api/
		objectName := fmt.Sprintf("products/%s", fileName)

		// Log the object name for debugging
		log.Printf("MinIO object name: %s", objectName)

		// Create background context for storage operations
		ctx := context.Background()

		// Upload to MinIO
		relativeURL, err := s.storageService.UploadFile(
			ctx,
			objectName,
			src,
			file.Size(),
			getContentType(fileName),
		)
		if err != nil {
			log.Printf("ERROR: Failed to upload file to MinIO: %v", err)
			log.Printf("Falling back to local storage")
			useLocalStorage = true
		} else {
			// IMPORTANT: The relativeURL should now already come back with the correct prefix (/images/products-api/)
			// from the MinioStorage.UploadFile method
			imageURL = relativeURL
			log.Printf("Stored MinIO image URL in database: %s", imageURL)

			// Verify that the URL starts with the MinIO pattern
			if !strings.HasPrefix(imageURL, "/images/products-api/") {
				// If for some reason it doesn't, fix it
				imageURL = fmt.Sprintf("/images/products-api/%s", fileName)
				log.Printf("Corrected MinIO image URL format to: %s", imageURL)
			}

			// Generate presigned URL for response
			log.Printf("Generating presigned URL for MinIO object")
			presignedURL, err := s.storageService.GetPresignedURL(
				ctx,
				imageURL,
				3600, // 1 hour expiry
			)
			if err != nil {
				log.Printf("ERROR: Failed to generate presigned URL: %v", err)
				// Fall back to relative URL
				fullImageURL = imageURL
			} else {
				fullImageURL = presignedURL
				// Cache the presigned URL if cache is available
				if s.urlCache != nil {
					s.urlCache.Set(imageURL, presignedURL)
					log.Printf("Cached presigned URL for %s", imageURL)
				}
			}
		}
	} else {
		useLocalStorage = true
	}

	// Fall back to local storage if needed or if MinIO is not available
	if useLocalStorage {
		log.Printf("Using local file storage")

		// Reset file pointer to beginning if possible
		if seeker, ok := src.(io.Seeker); ok {
			_, err = seeker.Seek(0, io.SeekStart)
			if err != nil {
				log.Printf("ERROR: Failed to reset file pointer: %v", err)
				// We'll create a new uploads directory and try our best
			}
		}

		// Create uploads directory if it doesn't exist
		uploadsDir := "./uploads/products"
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			log.Printf("ERROR: Failed to create uploads directory: %v", err)
			return "", fmt.Errorf("failed to create uploads directory: %w", err)
		}

		// Generate file path
		filePath := filepath.Join(uploadsDir, fileName)

		// Create destination file
		dst, err := os.Create(filePath)
		if err != nil {
			log.Printf("ERROR: Failed to create destination file: %v", err)
			return "", fmt.Errorf("failed to create destination file: %w", err)
		}
		defer dst.Close()

		// Copy the file
		_, err = io.Copy(dst, src)
		if err != nil {
			log.Printf("ERROR: Failed to copy file data: %v", err)
			return "", fmt.Errorf("failed to copy file data: %w", err)
		}

		// Set URL for database - always use /images/products/ prefix for local storage
		imageURL = fmt.Sprintf("/images/products/%s", fileName)
		fullImageURL = imageURL

		log.Printf("Successfully saved file to %s", filePath)
	}

	log.Printf("Image URL to store in database: %s", imageURL)
	log.Printf("Full image URL for response: %s", fullImageURL)

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

	// If this is set as primary and there are existing images,
	// we need to update other images to non-primary
	if isPrimary && len(images) > 0 {
		// Add the image first
		imageID, err := s.productRepo.AddProductImage(productID, imageURL, isPrimary, displayOrder)
		if err != nil {
			// Clean up from MinIO if using it
			if s.storageService != nil && !useLocalStorage {
				objectName := fmt.Sprintf("products/%s", fileName)
				s.storageService.DeleteFile(context.Background(), objectName)
			}
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
			// Clean up from MinIO if using it
			if s.storageService != nil && !useLocalStorage {
				objectName := fmt.Sprintf("products/%s", fileName)
				s.storageService.DeleteFile(context.Background(), objectName)
			}
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
	// Get the image URL before deleting
	images, err := s.productRepo.GetProductImages(0) // 0 means all products
	if err != nil {
		return err
	}

	var imageURL string
	for _, img := range images {
		if img.ID == imageID {
			imageURL = img.URL
			break
		}
	}

	// Delete from database first
	err = s.productRepo.DeleteProductImage(imageID)
	if err != nil {
		return err
	}

	// Create background context for storage operations
	ctx := context.Background()

	// Delete from MinIO if we found the image URL
	if imageURL != "" {
		// Extract object name from URL
		if strings.HasPrefix(imageURL, "/images/") {
			objectName := strings.TrimPrefix(imageURL, "/images/")
			err = s.storageService.DeleteFile(ctx, objectName)
			if err != nil {
				log.Printf("WARNING: Failed to delete image from MinIO: %v", err)
				// Don't return error as the database record is already deleted
			}
		}
	}

	return nil
}

// Helper function to determine content type based on file extension
func getContentType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// UpdateProductImageOrder updates the display order of a product image
func (s *ProductService) UpdateProductImageOrder(imageID int, displayOrder int) error {
	return s.productRepo.UpdateProductImageOrder(imageID, displayOrder)
}

// SetPrimaryProductImage sets the primary image for a product
func (s *ProductService) SetPrimaryProductImage(productID int, imageID int) error {
	return s.productRepo.SetPrimaryProductImage(productID, imageID)
}

// GetPresignedURL generates a presigned URL for the given object name
func (s *ProductService) GetPresignedURL(objectPath string) (string, error) {
	if s.storageService == nil {
		return "", fmt.Errorf("storage service not initialized")
	}

	ctx := context.Background()

	// Ensure correct path format for MinIO
	log.Printf("GetPresignedURL: Original path: %s", objectPath)

	// Handle different path patterns
	if strings.HasPrefix(objectPath, "/images/products-api/") {
		// Format: /images/products-api/filename.png -> products/filename.png
		objectPath = "products/" + strings.TrimPrefix(objectPath, "/images/products-api/")
		log.Printf("GetPresignedURL: Converted from /images/products-api/ to: %s", objectPath)
	} else if strings.Contains(objectPath, "/images/products-api/") {
		// Format: .../images/products-api/filename.png -> products/filename.png
		parts := strings.Split(objectPath, "/images/products-api/")
		if len(parts) > 1 {
			objectPath = "products/" + parts[1]
			log.Printf("GetPresignedURL: Extracted from URL with /images/products-api/: %s", objectPath)
		}
	} else if strings.HasPrefix(objectPath, "/images/products/") {
		// Format: /images/products/filename.png -> products/filename.png
		objectPath = "products/" + strings.TrimPrefix(objectPath, "/images/products/")
		log.Printf("GetPresignedURL: Converted from /images/products/ to: %s", objectPath)
	} else if strings.Contains(objectPath, "/images/products/") {
		// Format: .../images/products/filename.png -> products/filename.png
		parts := strings.Split(objectPath, "/images/products/")
		if len(parts) > 1 {
			objectPath = "products/" + parts[1]
			log.Printf("GetPresignedURL: Extracted from URL with /images/products/: %s", objectPath)
		}
	} else if !strings.HasPrefix(objectPath, "products/") {
		// If no prefix, add products/ if it appears to be a filename
		if !strings.Contains(objectPath, "/") {
			objectPath = "products/" + objectPath
			log.Printf("GetPresignedURL: Added products/ prefix to filename: %s", objectPath)
		}
	}

	// Use MinIO storage to get presigned URL
	presignedURL, err := s.storageService.GetPresignedURL(ctx, objectPath, 3600)
	if err != nil {
		log.Printf("Error generating presigned URL: %v", err)
		return "", err
	}

	// Log URL (truncate to avoid overly long logs)
	if len(presignedURL) > 50 {
		log.Printf("Generated presigned URL (truncated): %s...", presignedURL[:50])
	} else {
		log.Printf("Generated presigned URL: %s", presignedURL)
	}

	return presignedURL, nil
}

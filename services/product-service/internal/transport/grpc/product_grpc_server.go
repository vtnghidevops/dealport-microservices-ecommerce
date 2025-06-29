// services/product-service/internal/transport/grpc/grpc_server.go
package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"product-service/internal/domain"
	"product-service/internal/service"
	pb "product-service/proto/product"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GrpcServer represents the gRPC server for product service
type GrpcServer struct {
	pb.UnimplementedProductServiceServer
	productService domain.ProductService
}

// NewGrpcServer creates a new gRPC server with the provided product service
func NewGrpcServer(productSvc domain.ProductService) *GrpcServer {
	log.Println("Creating new gRPC server for product service")
	return &GrpcServer{
		productService: productSvc,
	}
}

// GetProduct implements the GetProduct RPC method
func (s *GrpcServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	product, err := s.productService.GetProductByID(int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return convertDomainProductToProto(product), nil
}

// GetProductBySlug implements the GetProductBySlug RPC method
func (s *GrpcServer) GetProductBySlug(ctx context.Context, req *pb.GetProductBySlugRequest) (*pb.Product, error) {
	product, err := s.productService.GetProductBySlug(req.Slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by slug: %w", err)
	}

	return convertDomainProductToProto(product), nil
}

// ListProducts implements the ListProducts RPC method
func (s *GrpcServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	filters := make(map[string]string)
	for k, v := range req.Filters {
		filters[k] = v
	}

	products, total, err := s.productService.GetAllProducts(int(req.Page), int(req.PageSize), filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	protoProducts := make([]*pb.Product, 0, len(products))
	for _, product := range products {
		protoProducts = append(protoProducts, convertDomainProductToProto(product))
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(total),
	}, nil
}

// CreateProduct implements the CreateProduct RPC method
func (s *GrpcServer) CreateProduct(ctx context.Context, req *pb.Product) (*pb.CreateProductResponse, error) {
	domainProduct := convertProtoToDomainProduct(req)

	id, err := s.productService.CreateProduct(domainProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Get the created product
	createdProduct, err := s.productService.GetProductByID(id)
	if err != nil {
		return nil, fmt.Errorf("product created but failed to retrieve: %w", err)
	}

	return &pb.CreateProductResponse{
		Id:      int32(id),
		Product: convertDomainProductToProto(createdProduct),
	}, nil
}

// UpdateProduct implements the UpdateProduct RPC method
func (s *GrpcServer) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.StatusResponse, error) {
	// Convert to domain model
	domainProduct := convertProtoToDomainProduct(req)

	// Update the product in the database
	err := s.productService.UpdateProduct(domainProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// No need to get the product again just to return a status
	return &pb.StatusResponse{
		Success: true,
		Message: "product updated successfully",
	}, nil
}

// PatchProduct implements the PatchProduct RPC method
func (s *GrpcServer) PatchProduct(ctx context.Context, req *pb.PatchProductRequest) (*pb.Product, error) {
	// Get existing product
	existingProduct, err := s.productService.GetProductByID(int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Parse JSON patch updates
	var updates map[string]interface{}
	err = json.Unmarshal([]byte(req.UpdatesJson), &updates)
	if err != nil {
		return nil, fmt.Errorf("invalid patch data: %w", err)
	}

	// We'll use regex patterns directly where needed

	// Apply updates to the existing product
	for field, value := range updates {
		switch field {
		case "name":
			existingProduct.Name = getString(value)
		case "type":
			existingProduct.Type = getString(value)
		case "slug":
			existingProduct.Slug = getString(value)
		case "description":
			existingProduct.Description = getString(value)
		case "price":
			existingProduct.Price = getFloat64(value)
		case "original_price", "originalPrice":
			existingProduct.OriginalPrice = getFloat64(value)
		case "discount":
			existingProduct.Discount = getFloat64(value)
		case "category_id", "categoryId":
			existingProduct.CategoryID = getInt(value)
		case "category_slug", "categorySlug":
			existingProduct.CategorySlug = getString(value)
		case "stock_quantity", "stockQuantity":
			existingProduct.StockQuantity = getInt(value)
		case "brand":
			existingProduct.Brand = getString(value)
		case "features":
			if featuresArray, ok := value.([]interface{}); ok {
				features := make([]string, 0, len(featuresArray))
				for _, feat := range featuresArray {
					if feature, ok := feat.(string); ok {
						features = append(features, feature)
					}
				}
				existingProduct.Features = features
			}
		case "imgSlider", "imageSlider", "images":
			// Add handling for imgSlider field
			if imgArray, ok := value.([]interface{}); ok {
				imgSlider := make([]string, 0, len(imgArray))

				// Convert array elements to strings
				for _, imgVal := range imgArray {
					if imgURL, ok := imgVal.(string); ok && imgURL != "" {
						// Normalize URL:
						normalizedURL := imgURL

						// Case 1: MinIO presigned URL - extract /images/products-api/
						if (strings.HasPrefix(imgURL, "http://") || strings.HasPrefix(imgURL, "https://")) && strings.Contains(imgURL, "/images/products-api/") {
							pattern := regexp.MustCompile(`(/images/products-api/[^?]+)`)
							matches := pattern.FindStringSubmatch(imgURL)
							if len(matches) > 0 {
								normalizedURL = matches[1]
								log.Printf("Normalized MinIO URL from %s to %s", imgURL, normalizedURL)
							}
							// Case 2: Local storage URL - extract /images/products/
						} else if (strings.HasPrefix(imgURL, "http://") || strings.HasPrefix(imgURL, "https://")) && strings.Contains(imgURL, "/images/products/") {
							pattern := regexp.MustCompile(`(/images/products/[^?]+)`)
							matches := pattern.FindStringSubmatch(imgURL)
							if len(matches) > 0 {
								normalizedURL = matches[1]
								log.Printf("Normalized local storage URL from %s to %s", imgURL, normalizedURL)
							}
						}

						imgSlider = append(imgSlider, normalizedURL)
					}
				}

				// Update the product's imgSlider field
				if len(imgSlider) > 0 {
					existingProduct.ImgSlider = imgSlider

					// Also create or update Images array from imgSlider
					images := make([]domain.ProductImage, 0, len(imgSlider))
					for i, url := range imgSlider {
						images = append(images, domain.ProductImage{
							ProductID:    existingProduct.ID,
							URL:          url,
							IsPrimary:    i == 0, // First image is primary
							DisplayOrder: i,
							CreatedAt:    time.Now(),
						})
					}

					if len(images) > 0 {
						existingProduct.Images = images
						// Update primary image URL
						existingProduct.ImageURL = images[0].URL
					}
				}
			}
			// Handle other fields as needed
		}
	}

	// Update the product
	err = s.productService.UpdateProduct(existingProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to patch product: %w", err)
	}

	// Get the updated product to return it
	updatedProduct, err := s.productService.GetProductByID(int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("product was updated but could not be retrieved: %w", err)
	}

	// Convert to proto product and return
	protoProduct := convertDomainProductToProto(updatedProduct)

	return protoProduct, nil
}

// Helper functions for type conversion
func getString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func getFloat64(v interface{}) float64 {
	val, ok := getFloat64Value(v)
	if !ok {
		return 0
	}
	return val
}

func getInt(v interface{}) int {
	val, ok := getIntValue(v)
	if !ok {
		return 0
	}
	return val
}

func getFloat64Value(v interface{}) (float64, bool) {
	switch value := v.(type) {
	case float64:
		return value, true
	case float32:
		return float64(value), true
	case int:
		return float64(value), true
	case int64:
		return float64(value), true
	case string:
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

func getIntValue(v interface{}) (int, bool) {
	switch value := v.(type) {
	case int:
		return value, true
	case int64:
		return int(value), true
	case float64:
		return int(value), true
	case float32:
		return int(value), true
	case string:
		if i, err := strconv.Atoi(value); err == nil {
			return i, true
		}
	}
	return 0, false
}

// DeleteProduct implements the DeleteProduct RPC method
func (s *GrpcServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.StatusResponse, error) {
	err := s.productService.DeleteProduct(int(req.Id))
	if err != nil {
		return &pb.StatusResponse{
			Success: false,
			Message: fmt.Sprintf("failed to delete product: %v", err),
		}, nil
	}

	return &pb.StatusResponse{
		Success: true,
		Message: "product deleted successfully",
	}, nil
}

// GetProductReviews implements the GetProductReviews RPC method
func (s *GrpcServer) GetProductReviews(ctx context.Context, req *pb.GetProductReviewsRequest) (*pb.GetProductReviewsResponse, error) {
	reviews, total, err := s.productService.GetProductReviews(int(req.ProductId), int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, fmt.Errorf("failed to get product reviews: %w", err)
	}

	protoReviews := make([]*pb.ProductReview, 0, len(reviews))
	for _, review := range reviews {
		protoReviews = append(protoReviews, &pb.ProductReview{
			Id:         int32(review.ID),
			ProductId:  int32(review.ProductID),
			UserName:   review.UserName,
			Email:      review.UserID,
			Rating:     int32(review.Rating),
			ReviewDate: review.CreatedAt.Format(time.RFC3339),
			ReviewText: review.Comment,
		})
	}

	return &pb.GetProductReviewsResponse{
		Reviews: protoReviews,
		Total:   int32(total),
	}, nil
}

// AddProductReview implements the AddProductReview RPC method
func (s *GrpcServer) AddProductReview(ctx context.Context, req *pb.ProductReview) (*pb.AddProductReviewResponse, error) {

	domainReview := &domain.ProductReview{
		ProductID: int(req.ProductId),
		UserName:  req.UserName,
		Rating:    float64(req.Rating),
		Comment:   req.ReviewText,
	}

	// Extract user ID from Email field if available (workaround for user ID)
	if req.Email != "" {
		domainReview.UserID = req.Email
		// Log to debug
		log.Printf("Using user ID from email field: %s", req.Email)
	}

	id, err := s.productService.AddProductReview(domainReview)
	if err != nil {
		return nil, fmt.Errorf("failed to add product review: %w", err)
	}

	return &pb.AddProductReviewResponse{
		Id: int32(id),
	}, nil
}

// UpdateProductReview implements the UpdateProductReview RPC method
func (s *GrpcServer) UpdateProductReview(ctx context.Context, req *pb.ProductReview) (*pb.StatusResponse, error) {

	domainReview := &domain.ProductReview{
		ID:        int(req.Id),
		ProductID: int(req.ProductId),
		UserName:  req.UserName,
		Rating:    float64(req.Rating),
		Comment:   req.ReviewText,
	}

	err := s.productService.UpdateProductReview(domainReview)
	if err != nil {
		return &pb.StatusResponse{
			Success: false,
			Message: fmt.Sprintf("failed to update product review: %v", err),
		}, nil
	}

	return &pb.StatusResponse{
		Success: true,
		Message: "product review updated successfully",
	}, nil
}

// GetTopRatedTestimonials implements the GetTopRatedTestimonials RPC method
func (s *GrpcServer) GetTopRatedTestimonials(ctx context.Context, req *pb.GetTopRatedTestimonialsRequest) (*pb.GetTopRatedTestimonialsResponse, error) {
	limit := 5 // Default limit
	if req.Limit > 0 {
		limit = int(req.Limit)
	}

	testimonials, err := s.productService.GetRandomTopRatedReviews(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top rated testimonials: %w", err)
	}

	protoTestimonials := make([]*pb.Testimonial, 0, len(testimonials))
	for _, t := range testimonials {
		protoTestimonials = append(protoTestimonials, &pb.Testimonial{
			Id:         int32(t.ProductID),
			ProductId:  int32(t.ProductID),
			UserName:   t.UserName,
			Avatar:     t.Avatar,
			Rating:     int32(t.Rating),
			ReviewText: t.Review,
		})
	}

	return &pb.GetTopRatedTestimonialsResponse{
		Testimonials: protoTestimonials,
	}, nil
}

// GetProductImages implements the GetProductImages RPC method
func (s *GrpcServer) GetProductImages(ctx context.Context, req *pb.GetProductImagesRequest) (*pb.GetProductImagesResponse, error) {
	images, err := s.productService.GetProductImages(int(req.ProductId))
	if err != nil {
		return nil, fmt.Errorf("failed to get product images: %w", err)
	}

	protoImages := make([]*pb.ProductImage, 0, len(images))
	for _, image := range images {
		protoImages = append(protoImages, &pb.ProductImage{
			Id:           int32(image.ID),
			ProductId:    int32(image.ProductID),
			Url:          image.URL,
			IsPrimary:    image.IsPrimary,
			DisplayOrder: int32(image.DisplayOrder),
			CreatedAt:    image.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.GetProductImagesResponse{
		Images: protoImages,
	}, nil
}

// UploadProductImage implements the UploadProductImage RPC method
func (s *GrpcServer) UploadProductImage(ctx context.Context, req *pb.UploadProductImageRequest) (*pb.UploadProductImageResponse, error) {
	// Create a temp file
	uploadsDir := "./uploads/products"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create uploads directory: %w", err)
	}

	fileExt := filepath.Ext(req.Filename)
	if fileExt == "" {
		fileExt = ".png" // Default extension if none provided
	}

	fileName := fmt.Sprintf("%d_%d%s", req.ProductId, time.Now().UnixNano(), fileExt)
	filePath := filepath.Join(uploadsDir, fileName)

	// Save the image data to the file
	err := os.WriteFile(filePath, req.ImageData, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	// Create an in-memory file instead of using customFileHeader
	// This avoids any potential type conversion issues
	tmpFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open saved file: %w", err)
	}
	defer tmpFile.Close()

	// Upload the image through the service
	// We don't use the direct file, the service will read and copy the file again
	// This is redundant but safer for now
	fileHeader := &multipart.FileHeader{
		Filename: fileName,
		Size:     int64(len(req.ImageData)),
		Header:   make(map[string][]string),
	}
	fileHeader.Header.Set("Content-Type", req.ContentType)

	// Use a new custom file wrapper that 100% implements domain.FileUpload
	fileUpload := &grpcFileUpload{
		filename: fileName,
		size:     int64(len(req.ImageData)),
		path:     filePath,
	}

	// Upload the image through the service
	imageURL, err := s.productService.UploadProductImage(int(req.ProductId), fileUpload, req.IsPrimary)
	if err != nil {
		// Clean up the temp file
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to process image upload: %w", err)
	}

	return &pb.UploadProductImageResponse{
		Url: imageURL,
	}, nil
}

// grpcFileUpload is a simplified implementation of domain.FileUpload
// that just wraps an already saved file
type grpcFileUpload struct {
	filename string
	size     int64
	path     string
}

func (f *grpcFileUpload) Open() (multipart.File, error) {
	return os.Open(f.path)
}

func (f *grpcFileUpload) Filename() string {
	return f.filename
}

func (f *grpcFileUpload) Size() int64 {
	return f.size
}

// DeleteProductImage implements the DeleteProductImage RPC method
func (s *GrpcServer) DeleteProductImage(ctx context.Context, req *pb.DeleteProductImageRequest) (*pb.StatusResponse, error) {
	err := s.productService.DeleteProductImage(int(req.ImageId))
	if err != nil {
		return &pb.StatusResponse{
			Success: false,
			Message: fmt.Sprintf("failed to delete product image: %v", err),
		}, nil
	}

	return &pb.StatusResponse{
		Success: true,
		Message: "product image deleted successfully",
	}, nil
}

// SetPrimaryProductImage implements the SetPrimaryProductImage RPC method
func (s *GrpcServer) SetPrimaryProductImage(ctx context.Context, req *pb.SetPrimaryProductImageRequest) (*pb.StatusResponse, error) {
	err := s.productService.SetPrimaryProductImage(int(req.ProductId), int(req.ImageId))
	if err != nil {
		return &pb.StatusResponse{
			Success: false,
			Message: fmt.Sprintf("failed to set primary product image: %v", err),
		}, nil
	}

	return &pb.StatusResponse{
		Success: true,
		Message: "primary product image set successfully",
	}, nil
}

// GetProductImageFile implements the GetProductImageFile RPC method
func (s *GrpcServer) GetProductImageFile(ctx context.Context, req *pb.GetProductImageFileRequest) (*pb.GetProductImageFileResponse, error) {
	// This method now redirects to MinIO storage instead of using local files
	log.Printf("GetProductImageFile requested for: %s, redirecting to MinIO", req.Filename)

	// Extract the object name (only the filename portion)
	objectName := fmt.Sprintf("products/%s", req.Filename)

	// Determine content type based on file extension
	var contentType string
	switch strings.ToLower(filepath.Ext(req.Filename)) {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	default:
		contentType = "application/octet-stream"
	}

	// Get the product service with storage capabilities from product_service.go
	storageService, ok := s.productService.(*service.ProductService)
	if !ok {
		return nil, fmt.Errorf("storage service not available")
	}

	// Try to get a presigned URL for the file
	presignedURL, err := storageService.GetPresignedURL(objectName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// Return a response that tells the client to redirect to the presigned URL
	return &pb.GetProductImageFileResponse{
		ImageData:     []byte{}, // Empty image data
		ContentType:   contentType,
		PresignedUrl:  presignedURL, // New field added to proto
		RedirectToUrl: true,         // Tell client to redirect
	}, nil
}

// GetHealth implements the GetHealth RPC method
func (s *GrpcServer) GetHealth(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// GetPresignedURL generates a presigned URL for accessing object from storage
func (s *GrpcServer) GetPresignedURL(ctx context.Context, req *pb.PresignedURLRequest) (*pb.PresignedURLResponse, error) {
	log.Printf("GetPresignedURL request for path: %s", req.ObjectPath)

	// Use object_path if provided, otherwise construct from filename
	objectPath := req.ObjectPath
	if objectPath == "" && req.Filename != "" {
		// Assume it's a product image if only filename is provided
		objectPath = "/images/products/" + req.Filename
		log.Printf("Using constructed object path: %s", objectPath)
	}

	if objectPath == "" {
		log.Printf("ERROR: Missing object path or filename")
		return &pb.PresignedURLResponse{
			Success: false,
			Error:   "missing object path or filename",
		}, nil
	}

	// Call service method to get presigned URL
	presignedURL, err := s.productService.GetPresignedURL(objectPath)
	if err != nil {
		log.Printf("ERROR: Failed to generate presigned URL: %v", err)
		return &pb.PresignedURLResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Log success (truncated URL for logs)
	maxLogLength := 50
	if len(presignedURL) > maxLogLength {
		log.Printf("Generated presigned URL (truncated): %s...", presignedURL[:maxLogLength])
	} else {
		log.Printf("Generated presigned URL: %s", presignedURL)
	}

	return &pb.PresignedURLResponse{
		PresignedUrl: presignedURL,
		Success:      true,
	}, nil
}

// Helper function to convert domain Product to proto Product
func convertDomainProductToProto(product *domain.Product) *pb.Product {
	if product == nil {
		return nil
	}

	// Function to process image URLs
	addDomainToURL := func(url string) string {
		if url == "" {
			return ""
		}
		// If URL already has http:// or https://, keep as is (presigned URLs)
		if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
			return url
		}

		// For relative URLs, keep as is for frontend processing
		return url
	}

	protoProduct := &pb.Product{
		Id:            int32(product.ID),
		Type:          product.Type,
		Name:          product.Name,
		Description:   product.Description,
		Slug:          product.Slug,
		Price:         product.Price,
		ImageUrl:      product.ImageURL,
		CategoryId:    int32(product.CategoryID),
		CategorySlug:  product.CategorySlug,
		StockQuantity: int32(product.StockQuantity),
		OriginalPrice: product.OriginalPrice,
		Discount:      product.Discount,
		Brand:         product.Brand,
		Orders:        int32(product.Orders),
		CreatedAt:     product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     product.UpdatedAt.Format(time.RFC3339),
	}

	// Convert tags
	protoProduct.Tags = append([]string{}, product.Tags...)

	// Convert features
	protoProduct.Features = append([]string{}, product.Features...)

	// Convert img_slider - add domain to each URL
	for _, url := range product.ImgSlider {
		protoProduct.ImgSlider = append(protoProduct.ImgSlider, addDomainToURL(url))
	}

	// Convert categories
	protoProduct.Categories = append([]string{}, product.Categories...)

	// Convert shipping info
	protoProduct.ShippingInfo = &pb.ShippingInfo{
		Courier: product.ShippingInfo.Courier,
		Local:   product.ShippingInfo.Local,
		Ups:     product.ShippingInfo.Ups,
		Global:  product.ShippingInfo.Global,
	}

	// Convert reviews avg
	protoProduct.ReviewsAvg = &pb.ProductRating{
		AverageRating: product.ReviewsAvg.AverageRating,
		Count:         int32(product.ReviewsAvg.Count),
	}

	// Convert images - add domain to URL of each image
	protoImages := make([]*pb.ProductImage, 0, len(product.Images))
	for _, image := range product.Images {
		protoImages = append(protoImages, &pb.ProductImage{
			Id:           int32(image.ID),
			ProductId:    int32(image.ProductID),
			Url:          addDomainToURL(image.URL),
			IsPrimary:    image.IsPrimary,
			DisplayOrder: int32(image.DisplayOrder),
			CreatedAt:    image.CreatedAt.Format(time.RFC3339),
		})
	}
	protoProduct.Images = protoImages

	// Convert UI metadata to string if it exists
	if product.UIMetadata != nil {
		protoProduct.UiMetadata = string(product.UIMetadata)
	}

	return protoProduct
}

// Helper function to convert proto Product to domain Product
func convertProtoToDomainProduct(protoProduct *pb.Product) *domain.Product {
	if protoProduct == nil {
		return nil
	}

	// Parse dates
	var createdAt, updatedAt time.Time
	if protoProduct.CreatedAt != "" {
		createdAt, _ = time.Parse(time.RFC3339, protoProduct.CreatedAt)
	} else {
		createdAt = time.Now()
	}

	if protoProduct.UpdatedAt != "" {
		updatedAt, _ = time.Parse(time.RFC3339, protoProduct.UpdatedAt)
	} else {
		updatedAt = time.Now()
	}

	domainProduct := &domain.Product{
		ID:            int(protoProduct.Id),
		Type:          protoProduct.Type,
		Name:          protoProduct.Name,
		Description:   protoProduct.Description,
		Slug:          protoProduct.Slug,
		Price:         protoProduct.Price,
		ImageURL:      protoProduct.ImageUrl,
		CategoryID:    int(protoProduct.CategoryId),
		CategorySlug:  protoProduct.CategorySlug,
		StockQuantity: int(protoProduct.StockQuantity),
		OriginalPrice: protoProduct.OriginalPrice,
		Discount:      protoProduct.Discount,
		Brand:         protoProduct.Brand,
		Orders:        int(protoProduct.Orders),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	// Convert tags
	if len(protoProduct.Tags) > 0 {
		domainProduct.Tags = append([]string{}, protoProduct.Tags...)
	}

	// Convert features
	if len(protoProduct.Features) > 0 {
		domainProduct.Features = append([]string{}, protoProduct.Features...)
	}

	// Convert img_slider
	if len(protoProduct.ImgSlider) > 0 {
		domainProduct.ImgSlider = append([]string{}, protoProduct.ImgSlider...)
	}

	// Convert categories
	if len(protoProduct.Categories) > 0 {
		domainProduct.Categories = append([]string{}, protoProduct.Categories...)
	}

	// Convert shipping info
	if protoProduct.ShippingInfo != nil {
		domainProduct.ShippingInfo = domain.ShippingInfo{
			Courier: protoProduct.ShippingInfo.Courier,
			Local:   protoProduct.ShippingInfo.Local,
			Ups:     protoProduct.ShippingInfo.Ups,
			Global:  protoProduct.ShippingInfo.Global,
		}
	}

	// Convert images
	domainImages := make([]domain.ProductImage, 0, len(protoProduct.Images))
	for _, protoImage := range protoProduct.Images {
		var createdAt time.Time
		if protoImage.CreatedAt != "" {
			createdAt, _ = time.Parse(time.RFC3339, protoImage.CreatedAt)
		} else {
			createdAt = time.Now()
		}

		domainImages = append(domainImages, domain.ProductImage{
			ID:           int(protoImage.Id),
			ProductID:    int(protoImage.ProductId),
			URL:          protoImage.Url,
			IsPrimary:    protoImage.IsPrimary,
			DisplayOrder: int(protoImage.DisplayOrder),
			CreatedAt:    createdAt,
		})
	}
	domainProduct.Images = domainImages

	// Convert UI metadata if it exists
	if protoProduct.UiMetadata != "" {
		domainProduct.UIMetadata = json.RawMessage(protoProduct.UiMetadata)
	}

	return domainProduct
}

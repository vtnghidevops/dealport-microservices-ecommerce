package grpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	pb "broker-service/proto/product"
	"encoding/json"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProductServiceAddress is the address of the product service
// Get from environment variable or use default for Docker environment
func getProductServiceAddress() string {
	productServiceAddr := os.Getenv("PRODUCT_SERVICE_HOST")
	if productServiceAddr == "" {
		return "product-service:50053" // Use service name for Docker environment
	}
	return productServiceAddr
}

// ProductClient is a gRPC client for the product service
type ProductClient struct {
	productClient  pb.ProductServiceClient
	categoryClient pb.CategoryServiceClient
	bannerClient   pb.BannerServiceClient
	adsClient      pb.AdsServiceClient
	conn           *grpc.ClientConn
}

// NewProductClient creates a new gRPC client that connects to the product service
func NewProductClient() (*ProductClient, error) {
	// Get the address of the product service
	productServiceAddr := getProductServiceAddress()

	// Log the attempt to connect to the product service
	log.Printf("Attempting to connect to product service at %s", productServiceAddr)

	// Connect to the gRPC server (insecure for development environment)
	conn, err := grpc.Dial(productServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	client := &ProductClient{
		productClient:  pb.NewProductServiceClient(conn),
		categoryClient: pb.NewCategoryServiceClient(conn),
		bannerClient:   pb.NewBannerServiceClient(conn),
		adsClient:      pb.NewAdsServiceClient(conn),
		conn:           conn,
	}

	return client, nil
}

// Đóng kết nối gRPC
func (c *ProductClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CheckHealth kiểm tra trạng thái của product-service
func (c *ProductClient) CheckHealth() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.productClient.GetHealth(ctx, &pb.HealthRequest{})
	if err != nil {
		return "", fmt.Errorf("health check failed: %w", err)
	}

	return resp.Status, nil
}

// Product Service Methods

// GetProduct gets a product by ID
func (c *ProductClient) GetProduct(id int) (*pb.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.productClient.GetProduct(ctx, &pb.GetProductRequest{
		Id: int32(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return resp, nil
}

// GetProductBySlug gets a product by slug
func (c *ProductClient) GetProductBySlug(slug string) (*pb.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.productClient.GetProductBySlug(ctx, &pb.GetProductBySlugRequest{
		Slug: slug,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get product by slug: %w", err)
	}

	return resp, nil
}

// ListProducts gets a list of products with pagination and filters
func (c *ProductClient) ListProducts(page, pageSize int, filters map[string]string) (*pb.ListProductsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.productClient.ListProducts(ctx, &pb.ListProductsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Filters:  filters,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return resp, nil
}

// CreateProduct creates a new product
func (c *ProductClient) CreateProduct(product *pb.Product) (*pb.CreateProductResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.productClient.CreateProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return resp, nil
}

// UpdateProduct updates an existing product
func (c *ProductClient) UpdateProduct(product *pb.Product) (*pb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Log entire product for debugging
	productJSON, _ := json.MarshalIndent(product, "", "  ")
	fmt.Printf("product khi được call bởi client (dạng JSON): %s\n", string(productJSON))

	// Specific detailed logging for imgSlider field
	log.Printf("DEBUG CLIENT: UpdateProduct called with product ID %d", product.Id)
	log.Printf("DEBUG CLIENT: imgSlider field length: %d", len(product.ImgSlider))
	for i, img := range product.ImgSlider {
		log.Printf("DEBUG CLIENT: imgSlider[%d] = %s", i, img)
	}

	resp, err := c.productClient.UpdateProduct(ctx, product)
	if err != nil {
		log.Printf("DEBUG CLIENT: Error updating product: %v", err)
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	log.Printf("DEBUG CLIENT: Product update response: %+v", resp)
	return resp, nil
}

// DeleteProduct deletes a product
func (c *ProductClient) DeleteProduct(id int) (*pb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.productClient.DeleteProduct(ctx, &pb.DeleteProductRequest{
		Id: int32(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete product: %w", err)
	}

	return resp, nil
}

// GetProductReviews gets reviews for a product
func (c *ProductClient) GetProductReviews(productID, page, pageSize int) (*pb.GetProductReviewsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.productClient.GetProductReviews(ctx, &pb.GetProductReviewsRequest{
		ProductId: int32(productID),
		Page:      int32(page),
		PageSize:  int32(pageSize),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get product reviews: %w", err)
	}

	return resp, nil
}

// AddProductReview adds a review to a product
func (c *ProductClient) AddProductReview(review *pb.ProductReview) (*pb.AddProductReviewResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.productClient.AddProductReview(ctx, review)
	if err != nil {
		return nil, fmt.Errorf("failed to add product review: %w", err)
	}

	return resp, nil
}

// UpdateProductReview updates a product review
func (c *ProductClient) UpdateProductReview(review *pb.ProductReview) (*pb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.productClient.UpdateProductReview(ctx, review)
	if err != nil {
		return nil, fmt.Errorf("failed to update product review: %w", err)
	}

	return resp, nil
}

// UploadProductImage uploads an image for a product
func (c *ProductClient) UploadProductImage(productID int, fileData []byte, filename, contentType string, isPrimary bool) (*pb.UploadProductImageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Longer timeout for uploads
	defer cancel()

	resp, err := c.productClient.UploadProductImage(ctx, &pb.UploadProductImageRequest{
		ProductId:   int32(productID),
		ImageData:   fileData,
		Filename:    filename,
		ContentType: contentType,
		IsPrimary:   isPrimary,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload product image: %w", err)
	}

	return resp, nil
}

// DeleteProductImage deletes a product image
func (c *ProductClient) DeleteProductImage(imageID int) (*pb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.productClient.DeleteProductImage(ctx, &pb.DeleteProductImageRequest{
		ImageId: int32(imageID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete product image: %w", err)
	}

	return resp, nil
}

// SetPrimaryProductImage sets a product image as primary
func (c *ProductClient) SetPrimaryProductImage(productID, imageID int) (*pb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.productClient.SetPrimaryProductImage(ctx, &pb.SetPrimaryProductImageRequest{
		ProductId: int32(productID),
		ImageId:   int32(imageID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set primary product image: %w", err)
	}

	return resp, nil
}

// GetProductImageFileFunc is used for testing to override the real gRPC call
var GetProductImageFileFunc func(filename string) ([]byte, string, error)

// GetProductImageFile retrieves image file data by filename
func (c *ProductClient) GetProductImageFile(filename string) ([]byte, string, error) {
	// If we're in a test environment and the mock function is set, use it
	if GetProductImageFileFunc != nil {
		return GetProductImageFileFunc(filename)
	}

	// Otherwise make the real gRPC call
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Longer timeout for image data
	defer cancel()

	resp, err := c.productClient.GetProductImageFile(ctx, &pb.GetProductImageFileRequest{
		Filename: filename,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to get product image file: %w", err)
	}

	return resp.ImageData, resp.ContentType, nil
}

// GetTopRatedTestimonials gets top rated testimonials
func (c *ProductClient) GetTopRatedTestimonials(limit int) (*pb.GetTopRatedTestimonialsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.productClient.GetTopRatedTestimonials(ctx, &pb.GetTopRatedTestimonialsRequest{
		Limit: int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get top rated testimonials: %w", err)
	}

	return resp, nil
}

// Category Service Methods

// GetCategoryByID gets a category by ID
func (c *ProductClient) GetCategoryByID(id int) (*pb.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.categoryClient.GetCategoryByID(ctx, &pb.GetCategoryByIDRequest{
		Id: int32(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return resp, nil
}

// GetCategoryBySlug gets a category by slug
func (c *ProductClient) GetCategoryBySlug(slug string) (*pb.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.categoryClient.GetCategoryBySlug(ctx, &pb.GetCategoryBySlugRequest{
		Slug: slug,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}

	return resp, nil
}

// ListCategories gets all categories with filters
func (c *ProductClient) ListCategories(filters map[string]string) (*pb.ListCategoriesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.categoryClient.ListCategories(ctx, &pb.ListCategoriesRequest{
		Filters: filters,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return resp, nil
}

// CreateCategory creates a new category
func (c *ProductClient) CreateCategory(category *pb.Category) (*pb.CreateCategoryResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.categoryClient.CreateCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return resp, nil
}

// UpdateCategory updates an existing category
func (c *ProductClient) UpdateCategory(category *pb.Category) (*pb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.categoryClient.UpdateCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return resp, nil
}

// DeleteCategory deletes a category
func (c *ProductClient) DeleteCategory(id int) (*pb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.categoryClient.DeleteCategory(ctx, &pb.DeleteCategoryRequest{
		Id: int32(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete category: %w", err)
	}

	return resp, nil
}

// Banner Service Methods

// GetBannerByID gets a banner by ID
func (c *ProductClient) GetBannerByID(id int) (*pb.Banner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.bannerClient.GetBannerByID(ctx, &pb.GetBannerByIDRequest{
		Id: int32(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get banner: %w", err)
	}

	return resp, nil
}

// GetBannersByType gets banners by type
func (c *ProductClient) GetBannersByType(bannerType string, limit int) (*pb.GetBannersByTypeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.bannerClient.GetBannersByType(ctx, &pb.GetBannersByTypeRequest{
		Type:  bannerType,
		Limit: int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get banners by type: %w", err)
	}

	return resp, nil
}

// ListBanners gets all banners with pagination and filters
func (c *ProductClient) ListBanners(page, pageSize int, filters map[string]string) (*pb.ListBannersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.bannerClient.ListBanners(ctx, &pb.ListBannersRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Filters:  filters,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list banners: %w", err)
	}

	return resp, nil
}

// CreateBanner creates a new banner
func (c *ProductClient) CreateBanner(banner *pb.Banner) (*pb.CreateBannerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.bannerClient.CreateBanner(ctx, banner)
	if err != nil {
		return nil, fmt.Errorf("failed to create banner: %w", err)
	}

	return resp, nil
}

// UpdateBanner updates an existing banner
func (c *ProductClient) UpdateBanner(banner *pb.Banner) (*pb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.bannerClient.UpdateBanner(ctx, banner)
	if err != nil {
		return nil, fmt.Errorf("failed to update banner: %w", err)
	}

	return resp, nil
}

// DeleteBanner deletes a banner
func (c *ProductClient) DeleteBanner(id int) (*pb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.bannerClient.DeleteBanner(ctx, &pb.DeleteBannerRequest{
		Id: int32(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete banner: %w", err)
	}

	return resp, nil
}

// Ads Service Methods

// GetAdsByLocation gets ads for a specific location
func (c *ProductClient) GetAdsByLocation(location string) (*pb.GetAdsByLocationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.adsClient.GetAdsByLocation(ctx, &pb.GetAdsByLocationRequest{
		Location: location,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get ads by location: %w", err)
	}

	return resp, nil
}

// GetAdsPlacementByID gets an ads placement by ID
func (c *ProductClient) GetAdsPlacementByID(id int) (*pb.AdsPlacement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.adsClient.GetAdsPlacementByID(ctx, &pb.GetAdsPlacementByIDRequest{
		Id: int32(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get ads placement: %w", err)
	}

	return resp, nil
}

// ListAdsPlacements gets all ads placements with pagination
func (c *ProductClient) ListAdsPlacements(page, pageSize int) (*pb.ListAdsPlacementsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.adsClient.ListAdsPlacements(ctx, &pb.ListAdsPlacementsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list ads placements: %w", err)
	}

	return resp, nil
}

// CreateAdsPlacement creates a new ads placement
func (c *ProductClient) CreateAdsPlacement(placement *pb.AdsPlacement) (*pb.CreateAdsPlacementResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.adsClient.CreateAdsPlacement(ctx, placement)
	if err != nil {
		return nil, fmt.Errorf("failed to create ads placement: %w", err)
	}

	return resp, nil
}

// UpdateAdsPlacement updates an existing ads placement
func (c *ProductClient) UpdateAdsPlacement(placement *pb.AdsPlacement) (*pb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.adsClient.UpdateAdsPlacement(ctx, placement)
	if err != nil {
		return nil, fmt.Errorf("failed to update ads placement: %w", err)
	}

	return resp, nil
}

// DeleteAdsPlacement deletes an ads placement
func (c *ProductClient) DeleteAdsPlacement(id int) (*pb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.adsClient.DeleteAdsPlacement(ctx, &pb.DeleteAdsPlacementRequest{
		Id: int32(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to delete ads placement: %w", err)
	}

	return resp, nil
}

// PatchProduct updates specific fields of a product
func (c *ProductClient) PatchProduct(id int, updates map[string]interface{}) (*pb.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert updates map to JSON string
	updatesJSON, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updates: %w", err)
	}

	log.Printf("DEBUG CLIENT: PatchProduct called with product ID %d", id)
	log.Printf("DEBUG CLIENT: Updates JSON: %s", string(updatesJSON))

	// Create PatchProductRequest
	req := &pb.PatchProductRequest{
		Id:          int32(id),
		UpdatesJson: string(updatesJSON),
	}

	// Call gRPC method with updated return type
	product, err := c.productClient.PatchProduct(ctx, req)
	if err != nil {
		log.Printf("DEBUG CLIENT: Error patching product: %v", err)
		return nil, fmt.Errorf("failed to patch product: %w", err)
	}

	log.Printf("DEBUG CLIENT: Product patch response received with ID: %d", product.Id)
	return product, nil
}

var (
	productClientInstance *ProductClient
	productClientOnce     sync.Once
	productClientError    error
)

// GetProductClient trả về một singleton instance của ProductClient
func GetProductClient() (*ProductClient, error) {
	productClientOnce.Do(func() {
		productClientInstance, productClientError = NewProductClient()
		if productClientError != nil {
			log.Printf("Error creating ProductClient: %v", productClientError)
		} else {
			log.Println("ProductClient created successfully")
		}
	})
	return productClientInstance, productClientError
}

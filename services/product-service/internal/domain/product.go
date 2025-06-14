package domain

import (
	"encoding/json"
	"mime/multipart"
	"time"
)

// Product represents a product entity
type Product struct {
	ID            int     `json:"id"`
	Type          string  `json:"type"` // normal, trending, top-sale, new, limited
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Slug          string  `json:"slug"`
	Price         float64 `json:"price"`
	ImageURL      string  `json:"imageUrl"`
	CategoryID    int     `json:"categoryId"`
	CategorySlug  string  `json:"categorySlug"`
	StockQuantity int     `json:"stockQuantity"`

	OriginalPrice float64        `json:"originalPrice,omitempty"`
	Discount      float64        `json:"discount,omitempty"`
	Images        []ProductImage `json:"images,omitempty"`
	Categories    []string       `json:"categories,omitempty"`

	Brand      string        `json:"brand,omitempty"`
	Tags       []string      `json:"tags,omitempty"`
	ReviewsAvg ProductRating `json:"reviewsAvg,omitempty"`
	Orders     int           `json:"orders,omitempty"`
	ImgSlider  []string      `json:"imgSlider,omitempty"`

	Features     []string        `json:"features,omitempty"`
	ShippingInfo ShippingInfo    `json:"shippingInfo,omitempty"`
	UIMetadata   json.RawMessage `json:"uiMetadata,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShippingInfo struct {
	Courier string `json:"courier,omitempty"`
	Local   string `json:"local,omitempty"`
	Ups     string `json:"ups,omitempty"`
	Global  string `json:"global,omitempty"`
}

// ProductImage represents an image associated with a product
type ProductImage struct {
	ID           int       `json:"id"`
	ProductID    int       `json:"productId"`
	URL          string    `json:"url"`
	IsPrimary    bool      `json:"isPrimary"`
	DisplayOrder int       `json:"displayOrder"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ProductReviews represents a summary of product reviews
type ProductRating struct {
	AverageRating float64 `json:"rating"`
	Count         int     `json:"count"`
}

// ProductReview represents a review for a product
type ProductReview struct {
	ID         int       `json:"id"`
	ProductID  int       `json:"productId"`
	UserID     string    `json:"userId"`
	UserName   string    `json:"userName,omitempty"`
	Rating     float64   `json:"rating"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"createdAt"`
	UserAvatar string    `json:"userAvatar,omitempty"` // Added for HappyCustomers display
}

// Testimonial represents a customer testimonial for display in HappyCustomers
type Testimonial struct {
	ID        string  `json:"userId"`
	UserName  string  `json:"userName"`
	Avatar    string  `json:"avatar"`
	Review    string  `json:"review"`
	Rating    float64 `json:"rating"`
	ProductID int     `json:"productId,omitempty"`
}

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	GetProductByID(id int) (*Product, error)
	GetProductBySlug(slug string) (*Product, error)
	GetAllProducts(page, pageSize int, filters map[string]string) ([]*Product, int, error)
	CreateProduct(product *Product) (int, error)
	UpdateProduct(product *Product) error
	DeleteProduct(id int) error

	// Product image methods
	GetProductImages(productID int) ([]ProductImage, error)
	AddProductImage(productID int, imageURL string, isPrimary bool, displayOrder int) (int, error)
	DeleteProductImage(imageID int) error
	UpdateProductImageOrder(imageID int, displayOrder int) error
	SetPrimaryProductImage(productID int, imageID int) error

	// Product tag methods
	GetProductTags(productID int) ([]string, error)

	// Review methods
	GetProductReviews(productID int, page, pageSize int) ([]*ProductReview, int, error)
	AddProductReview(review *ProductReview) (int, error)
	UpdateProductReview(review *ProductReview) error
	DeleteProductReview(id int) error

	// Get random 5-star testimonials for HappyCustomers
	GetRandomTopRatedReviews(limit int) ([]*Testimonial, error)
}

// ProductService defines the interface for product business logic
type ProductService interface {
	GetProductByID(id int) (*Product, error)
	GetProductBySlug(slug string) (*Product, error)
	GetAllProducts(page, pageSize int, filters map[string]string) ([]*Product, int, error)
	CreateProduct(product *Product) (int, error)
	UpdateProduct(product *Product) error
	DeleteProduct(id int) error

	// Product image methods
	GetProductImages(productID int) ([]ProductImage, error)
	UploadProductImage(productID int, file FileUpload, isPrimary bool) (string, error)
	DeleteProductImage(imageID int) error
	UpdateProductImageOrder(imageID int, displayOrder int) error
	SetPrimaryProductImage(productID int, imageID int) error

	// Review methods
	GetProductReviews(productID int, page, pageSize int) ([]*ProductReview, int, error)
	AddProductReview(review *ProductReview) (int, error)
	UpdateProductReview(review *ProductReview) error
	DeleteProductReview(id int) error

	// Get random 5-star testimonials for HappyCustomers
	GetRandomTopRatedReviews(limit int) ([]*Testimonial, error)

	// Get presigned URL for accessing images from storage
	GetPresignedURL(objectPath string) (string, error)
}

// FileUpload defines an interface for file uploads to allow custom implementations
type FileUpload interface {
	Open() (multipart.File, error)
	Filename() string
	Size() int64
}

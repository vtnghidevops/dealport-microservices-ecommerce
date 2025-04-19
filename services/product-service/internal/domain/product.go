package domain

import (
	"encoding/json"
	"time"
)

// Product represents a product entity
type Product struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Slug          string          `json:"slug"`
	Description   string          `json:"description"`
	Type          string          `json:"type"`
	Price         float64         `json:"price"`
	OriginalPrice float64         `json:"original_price,omitempty"`
	Discount      float64         `json:"discount,omitempty"`
	ImageURL      string          `json:"image_url"`
	CategoryID    string          `json:"category_id"`
	CategorySlug  string          `json:"category_slug"`
	StockQuantity int             `json:"stock_quantity"`
	Brand         string          `json:"brand,omitempty"`
	Features      json.RawMessage `json:"features,omitempty"`
	ShippingInfo  json.RawMessage `json:"shipping_info,omitempty"`
	Images        []ProductImage  `json:"images,omitempty"`
	Tags          []string        `json:"tags,omitempty"`
	Reviews       *ProductReviews `json:"reviews,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// ProductImage represents an image associated with a product
type ProductImage struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	URL          string    `json:"url"`
	IsPrimary    bool      `json:"is_primary"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

// ProductReviews represents a summary of product reviews
type ProductReviews struct {
	AverageRating float64 `json:"average_rating"`
	Count         int     `json:"count"`
}

// ProductReview represents a review for a product
type ProductReview struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name,omitempty"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	GetByID(id string) (*Product, error)
	GetBySlug(slug string) (*Product, error)
	List(page, pageSize int, filters map[string]string) ([]*Product, int, error)
	Create(product *Product) (string, error)
	Update(product *Product) error
	Delete(id string) error

	// Product image methods
	GetProductImages(productID string) ([]ProductImage, error)

	// Product tag methods
	GetProductTags(productID string) ([]string, error)

	// Review methods
	GetProductReviews(productID string, page, pageSize int) ([]*ProductReview, int, error)
	AddProductReview(review *ProductReview) (string, error)
	DeleteProductReview(id string) error
}

// ProductService defines the interface for product business logic
type ProductService interface {
	GetByID(id string) (*Product, error)
	GetBySlug(slug string) (*Product, error)
	List(page, pageSize int, filters map[string]string) ([]*Product, int, error)
	Create(product *Product) (string, error)
	Update(product *Product) error
	Delete(id string) error

	// Review methods
	GetProductReviews(productID string, page, pageSize int) ([]*ProductReview, int, error)
	AddProductReview(review *ProductReview) (string, error)
	DeleteProductReview(id string) error
}

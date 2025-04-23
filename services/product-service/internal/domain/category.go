package domain

import (
	"time"
)

// Category represents a product category
type Category struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description,omitempty"`
	ImageURL     string    `json:"image_url,omitempty"`
	ProductCount int       `json:"product_count"`
	IsActive     bool      `json:"is_active"`
	IsVisible    bool      `json:"is_visible"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	GetCategoryByID(id int) (*Category, error)
	GetCategoryBySlug(slug string) (*Category, error)
	GetAllCategories(filters map[string]string) ([]*Category, error)
	CreateCategory(category *Category) (int, error)
	UpdateCategory(category *Category) error
	DeleteCategory(id int) error
	SyncProductCounts() error
}

// CategoryService defines the interface for category business logic
type CategoryService interface {
	GetCategoryByID(id int) (*Category, error)
	GetCategoryBySlug(slug string) (*Category, error)
	GetAllCategories(filters map[string]string) ([]*Category, error)
	CreateCategory(category *Category) (int, error)
	UpdateCategory(category *Category) error
	DeleteCategory(id int) error
	SyncProductCounts() error
}

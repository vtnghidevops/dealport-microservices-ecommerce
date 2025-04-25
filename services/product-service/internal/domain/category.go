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
	ImageURL     string    `json:"imageUrl,omitempty"`
	ProductCount int       `json:"productCount"`
	IsActive     bool      `json:"isActive"`
	IsVisible    bool      `json:"isVisible"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
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

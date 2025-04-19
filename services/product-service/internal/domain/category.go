package domain

import (
	"time"
)

// Category represents a product category
type Category struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	Slug            string              `json:"slug"`
	Description     string              `json:"description,omitempty"`
	ImageURL        string              `json:"image_url,omitempty"`
	Icon            string              `json:"icon,omitempty"`
	BannerURL       string              `json:"banner_url,omitempty"`
	ProductCount    int                 `json:"product_count"`
	IsActive        bool                `json:"is_active"`
	IsVisible       bool                `json:"is_visible"`
	DisplayOrder    int                 `json:"display_order"`
	MetaTitle       string              `json:"meta_title,omitempty"`
	MetaDescription string              `json:"meta_description,omitempty"`
	ParentID        string              `json:"parent_id,omitempty"`
	Level           int                 `json:"level"`
	Attributes      []CategoryAttribute `json:"attributes,omitempty"`
	Children        []*Category         `json:"children,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

// CategoryAttribute represents an attribute for a category
type CategoryAttribute struct {
	ID         string    `json:"id"`
	CategoryID string    `json:"category_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Required   bool      `json:"required"`
	Options    []string  `json:"options,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	GetByID(id string) (*Category, error)
	GetBySlug(slug string) (*Category, error)
	List(filters map[string]string) ([]*Category, error)
	GetCategoryTree() ([]*Category, error)
	GetCategoryAttributes(categoryID string) ([]CategoryAttribute, error)
}

// CategoryService defines the interface for category business logic
type CategoryService interface {
	GetByID(id string) (*Category, error)
	GetBySlug(slug string) (*Category, error)
	List(filters map[string]string) ([]*Category, error)
	GetCategoryTree() ([]*Category, error)
}

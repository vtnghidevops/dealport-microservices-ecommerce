package domain

import (
	"encoding/json"
	"time"
)

// AdsPlacement represents an advertisement placement in the homepage
type AdsPlacement struct {
	ID             int             `json:"id"`
	Location       string          `json:"location"`
	ReferenceType  string          `json:"reference_type"`
	ReferenceID    int             `json:"reference_id"`
	DisplayOrder   int             `json:"display_order"`
	CustomTitle    string          `json:"custom_title,omitempty"`
	CustomImageURL string          `json:"custom_image_url,omitempty"`
	UISettings     json.RawMessage `json:"ui_settings,omitempty"`
	IsActive       bool            `json:"is_active"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// AdsData represents the combined ad data with product/category information
type AdsData struct {
	ID           int             `json:"id"`
	Type         string          `json:"type"` // "product" or "category"
	Name         string          `json:"name"`
	Slug         string          `json:"slug"`
	ImageURL     string          `json:"image_url"`
	Price        float64         `json:"price,omitempty"`
	CategoryID   int             `json:"category_id,omitempty"`
	CategorySlug string          `json:"category_slug,omitempty"`
	CustomTitle  string          `json:"custom_title,omitempty"`
	UISettings   json.RawMessage `json:"ui_settings,omitempty"`
}

// AdsRepository defines the interface for advertisement data operations
type AdsRepository interface {
	// Get ads by location (banner, display, gaming, new_fashion)
	GetAdsByLocation(location string) ([]*AdsData, error)

	// CRUD operations for ad placements
	CreateAdsPlacement(ad *AdsPlacement) (int, error)
	GetAdsPlacementByID(id int) (*AdsPlacement, error)
	UpdateAdsPlacement(ad *AdsPlacement) error
	DeleteAdsPlacement(id int) error

	// Get all ad placements with pagination
	GetAllAdsPlacements(page, pageSize int) ([]*AdsPlacement, int, error)
}

// AdsService defines the interface for advertisement business logic
type AdsService interface {
	// Get ads by location (banner, display, gaming, new_fashion)
	GetAdsByLocation(location string) ([]*AdsData, error)

	// CRUD operations for ad placements
	CreateAdsPlacement(ad *AdsPlacement) (int, error)
	GetAdsPlacementByID(id int) (*AdsPlacement, error)
	UpdateAdsPlacement(ad *AdsPlacement) error
	DeleteAdsPlacement(id int) error

	// Get all ad placements with pagination
	GetAllAdsPlacements(page, pageSize int) ([]*AdsPlacement, int, error)
}

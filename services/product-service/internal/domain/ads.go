package domain

import (
	"encoding/json"
	"time"
)

// AdsPlacement represents an advertisement placement in the homepage
type AdsPlacement struct {
	ID             int             `json:"id"`
	Location       string          `json:"location"`
	ReferenceType  string          `json:"referenceType"`
	ReferenceID    int             `json:"referenceId"`
	DisplayOrder   int             `json:"displayOrder"`
	CustomTitle    string          `json:"customTitle,omitempty"`
	CustomImageURL string          `json:"customImageUrl,omitempty"`
	UISettings     json.RawMessage `json:"uiSettings,omitempty"`
	IsActive       bool            `json:"isActive"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

// AdsData represents the combined ad data with product/category information
type AdsData struct {
	ID           int             `json:"id"`
	Type         string          `json:"type"` // "product" or "category"
	Name         string          `json:"name"`
	Slug         string          `json:"slug"`
	ImageURL     string          `json:"imageUrl"`
	Price        float64         `json:"price,omitempty"`
	CategoryID   int             `json:"categoryId,omitempty"`
	CategorySlug string          `json:"categorySlug,omitempty"`
	CustomTitle  string          `json:"customTitle,omitempty"`
	UISettings   json.RawMessage `json:"uiSettings,omitempty"`
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

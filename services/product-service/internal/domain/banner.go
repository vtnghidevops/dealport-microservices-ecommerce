// internal/domain/banner.go
package domain

import (
	"time"
)

// Banner represents a promotional banner
type Banner struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Subtitle       string    `json:"subtitle,omitempty"`
	Description    string    `json:"description,omitempty"`
	Discount       string    `json:"discount,omitempty"`
	HighlightText  string    `json:"highlightText,omitempty"`
	ImageURL       string    `json:"imageUrl"`
	LinkURL        string    `json:"linkUrl,omitempty"`
	ActionText     string    `json:"actionText,omitempty"`
	BackgroundColor string   `json:"backgroundColor,omitempty"`
	TextColor      string    `json:"textColor,omitempty"`
	AnimationType  string    `json:"animationType,omitempty"`
	IsActive       bool      `json:"isActive"`
	Priority       int       `json:"priority,omitempty"`
	Type           string    `json:"type"` // hero, promotional, category, seasonal, product
	ProductID      *int      `json:"productId,omitempty"`
	CategoryID     *int      `json:"categoryId,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// Định nghĩa các loại banner
const (
	BannerTypeHero        = "hero"
	BannerTypePromotional = "promotional"
	BannerTypeCategory    = "category"
	BannerTypeSeasonal    = "seasonal"
	BannerTypeProduct     = "product"
	BannerTypeCollection  = "collection"
)


// BannerRepository định nghĩa interface cho thao tác dữ liệu banner
type BannerRepository interface {
	GetBannerByID(id int) (*Banner, error)
	GetBannersByType(bannerType string, limit int) ([]*Banner, error)
	GetAllBanners(page, pageSize int, filters map[string]string) ([]*Banner, int, error)
	CreateBanner(banner *Banner) (int, error)
	UpdateBanner(banner *Banner) error
	DeleteBanner(id int) error
}

// BannerService định nghĩa interface cho logic nghiệp vụ banner
type BannerService interface {
	GetBannerByID(id int) (*Banner, error)
	GetBannersByType(bannerType string, limit int) ([]*Banner, error)
	GetAllBanners(page, pageSize int, filters map[string]string) ([]*Banner, int, error)
	CreateBanner(banner *Banner) (int, error)
	UpdateBanner(banner *Banner) error
	DeleteBanner(id int) error
}
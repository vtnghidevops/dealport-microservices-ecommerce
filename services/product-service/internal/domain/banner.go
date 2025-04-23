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
	HighlightText  string    `json:"highlight_text,omitempty"`
	ImageURL       string    `json:"image_url"`
	LinkURL        string    `json:"link_url,omitempty"`
	ActionText     string    `json:"action_text,omitempty"`
	BackgroundColor string   `json:"background_color,omitempty"`
	TextColor      string    `json:"text_color,omitempty"`
	AnimationType  string    `json:"animation_type,omitempty"`
	IsActive       bool      `json:"is_active"`
	Priority       int       `json:"priority,omitempty"`
	Type           string    `json:"type"` // hero, promotional, category, seasonal, product
	ProductID      *int      `json:"product_id,omitempty"`
	CategoryID     *int      `json:"category_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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
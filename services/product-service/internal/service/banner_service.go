// internal/service/banner_service.go
package service

import (
	"product-service/internal/domain"
)

// BannerService implements the domain.BannerService interface
type BannerService struct {
	repo domain.BannerRepository
}

// NewBannerService creates a new instance of BannerService
func NewBannerService(repo domain.BannerRepository) *BannerService {
	return &BannerService{
		repo: repo,
	}
}

// GetBannerByID returns a banner by ID
func (s *BannerService) GetBannerByID(id int) (*domain.Banner, error) {
	return s.repo.GetBannerByID(id)
}

// GetBannersByType returns banners by type with optional limit
func (s *BannerService) GetBannersByType(bannerType string, limit int) ([]*domain.Banner, error) {
	return s.repo.GetBannersByType(bannerType, limit)
}

// GetAllBanners returns all banners with pagination and filtering
func (s *BannerService) GetAllBanners(page, pageSize int, filters map[string]string) ([]*domain.Banner, int, error) {
	return s.repo.GetAllBanners(page, pageSize, filters)
}

// CreateBanner creates a new banner
func (s *BannerService) CreateBanner(banner *domain.Banner) (int, error) {
	// Validation logic here
	if banner.Title == "" {
		return 0, domain.ErrInvalidBanner
	}

	if banner.ImageURL == "" {
		return 0, domain.ErrInvalidBanner
	}

	// Set default values if not provided
	if banner.ActionText == "" {
		banner.ActionText = "Shop Now"
	}

	if banner.BackgroundColor == "" {
		banner.BackgroundColor = "#1e3a8a"
	}

	if banner.TextColor == "" {
		banner.TextColor = "#ffffff"
	}

	if banner.AnimationType == "" {
		banner.AnimationType = "fade"
	}

	return s.repo.CreateBanner(banner)
}

// UpdateBanner updates an existing banner
func (s *BannerService) UpdateBanner(banner *domain.Banner) error {
	// Validation logic here
	if banner.ID == 0 {
		return domain.ErrInvalidBanner
	}

	if banner.Title == "" {
		return domain.ErrInvalidBanner
	}

	if banner.ImageURL == "" {
		return domain.ErrInvalidBanner
	}

	return s.repo.UpdateBanner(banner)
}

// DeleteBanner deletes a banner by ID
func (s *BannerService) DeleteBanner(id int) error {
	return s.repo.DeleteBanner(id)
}
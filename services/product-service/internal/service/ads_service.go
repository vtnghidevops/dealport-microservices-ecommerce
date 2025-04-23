package service

import (
	"product-service/internal/domain"
)

// AdsService implements the domain.AdsService interface
type AdsService struct {
	repo domain.AdsRepository
}

// NewAdsService creates a new instance of AdsService
func NewAdsService(repo domain.AdsRepository) *AdsService {
	return &AdsService{
		repo: repo,
	}
}

// GetAdsByLocation returns ads for a specific location
func (s *AdsService) GetAdsByLocation(location string) ([]*domain.AdsData, error) {
	ads, err := s.repo.GetAdsByLocation(location)
	if err != nil {
		return nil, err
	}

	return ads, nil
}

// CreateAdsPlacement creates a new ads placement
func (s *AdsService) CreateAdsPlacement(ad *domain.AdsPlacement) (int, error) {
	return s.repo.CreateAdsPlacement(ad)
}

// GetAdsPlacementByID returns an ads placement by ID
func (s *AdsService) GetAdsPlacementByID(id int) (*domain.AdsPlacement, error) {
	return s.repo.GetAdsPlacementByID(id)
}

// UpdateAdsPlacement updates an existing ads placement
func (s *AdsService) UpdateAdsPlacement(ad *domain.AdsPlacement) error {
	return s.repo.UpdateAdsPlacement(ad)
}

// DeleteAdsPlacement deletes an ads placement by ID
func (s *AdsService) DeleteAdsPlacement(id int) error {
	return s.repo.DeleteAdsPlacement(id)
}

// GetAllAdsPlacements returns all ads placements with pagination
func (s *AdsService) GetAllAdsPlacements(page, pageSize int) ([]*domain.AdsPlacement, int, error) {
	// Default values
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	return s.repo.GetAllAdsPlacements(page, pageSize)
}

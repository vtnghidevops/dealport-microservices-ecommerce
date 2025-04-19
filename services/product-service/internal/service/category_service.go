package service

import (
	"product-service/internal/domain"
)

// CategoryService implements the domain.CategoryService interface
type CategoryService struct {
	repo domain.CategoryRepository
}

// NewCategoryService creates a new instance of CategoryService
func NewCategoryService(repo domain.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

// GetByID returns a category by ID
func (s *CategoryService) GetByID(id string) (*domain.Category, error) {
	return s.repo.GetByID(id)
}

// GetBySlug returns a category by slug
func (s *CategoryService) GetBySlug(slug string) (*domain.Category, error) {
	return s.repo.GetBySlug(slug)
}

// List returns categories with optional filtering
func (s *CategoryService) List(filters map[string]string) ([]*domain.Category, error) {
	return s.repo.List(filters)
}

// GetCategoryTree returns the full category tree
func (s *CategoryService) GetCategoryTree() ([]*domain.Category, error) {
	return s.repo.GetCategoryTree()
}

// GetCategoryAttributes returns attributes for a specific category
func (s *CategoryService) GetCategoryAttributes(categoryID string) ([]domain.CategoryAttribute, error) {
	return s.repo.GetCategoryAttributes(categoryID)
}

package service

import (
	"product-service/internal/domain"
	"product-service/internal/util"
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

// GetCategoryByID returns a category by ID
func (s *CategoryService) GetCategoryByID(id int) (*domain.Category, error) {
	return s.repo.GetCategoryByID(id)
}

// GetCategoryBySlug returns a category by slug
func (s *CategoryService) GetCategoryBySlug(slug string) (*domain.Category, error) {
	return s.repo.GetCategoryBySlug(slug)
}

// GetAllCategories returns categories with optional filtering
func (s *CategoryService) GetAllCategories(filters map[string]string) ([]*domain.Category, error) {
	return s.repo.GetAllCategories(filters)
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(category *domain.Category) (int, error) {
	// Validate category fields here before passing to repo
	// For example:
	if category.Name == "" {
		return 0, domain.ErrInvalidCategory
	}

	// Create slug from name if not provided
	if category.Slug == "" {
		category.Slug = util.CreateSlug(category.Name)
	}

	return s.repo.CreateCategory(category)
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(category *domain.Category) error {
	// Validate category fields here before passing to repo
	if category.ID == 0 {
		return domain.ErrInvalidCategory
	}

	if category.Name == "" {
		return domain.ErrInvalidCategory
	}

	// Create slug from name if not provided
	if category.Slug == "" {
		category.Slug = util.CreateSlug(category.Name)
	}

	return s.repo.UpdateCategory(category)
}

// DeleteCategory deletes a category by ID
func (s *CategoryService) DeleteCategory(id int) error {
	return s.repo.DeleteCategory(id)
}

// SyncProductCounts updates product counts for all categories
func (s *CategoryService) SyncProductCounts() error {
	return s.repo.SyncProductCounts()
}

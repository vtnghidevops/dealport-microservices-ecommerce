package mocks

import (
	"product-service/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockCategoryRepository is a mock implementation of CategoryRepository interface
type MockCategoryRepository struct {
	mock.Mock
}

// GetCategoryByID mocks the GetCategoryByID method
func (m *MockCategoryRepository) GetCategoryByID(id int) (*domain.Category, error) {
	args := m.Called(id)

	if category := args.Get(0); category != nil {
		return category.(*domain.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

// GetCategoryBySlug mocks the GetCategoryBySlug method
func (m *MockCategoryRepository) GetCategoryBySlug(slug string) (*domain.Category, error) {
	args := m.Called(slug)

	if category := args.Get(0); category != nil {
		return category.(*domain.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

// GetAllCategories mocks the GetAllCategories method
func (m *MockCategoryRepository) GetAllCategories(filters map[string]string) ([]*domain.Category, error) {
	args := m.Called(filters)

	if categories := args.Get(0); categories != nil {
		return categories.([]*domain.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

// CreateCategory mocks the CreateCategory method
func (m *MockCategoryRepository) CreateCategory(category *domain.Category) (int, error) {
	args := m.Called(category)
	return args.Int(0), args.Error(1)
}

// UpdateCategory mocks the UpdateCategory method
func (m *MockCategoryRepository) UpdateCategory(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

// DeleteCategory mocks the DeleteCategory method
func (m *MockCategoryRepository) DeleteCategory(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// SyncProductCounts mocks the SyncProductCounts method
func (m *MockCategoryRepository) SyncProductCounts() error {
	args := m.Called()
	return args.Error(0)
}

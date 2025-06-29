package unit

import (
	"errors"
	"product-service/internal/domain"
	"product-service/internal/service"
	"product-service/tests/mocks"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// createMockCategoryService tạo một category service với repository mock
func createMockCategoryService() (*mocks.MockCategoryRepository, domain.CategoryService) {
	mockRepo := new(mocks.MockCategoryRepository)
	categoryService := service.NewCategoryService(mockRepo)
	return mockRepo, categoryService
}

// createSampleCategory tạo một danh mục mẫu
func createSampleCategory(id int, name string) *domain.Category {
	now := time.Now()
	return &domain.Category{
		ID:           id,
		Name:         name,
		Slug:         createSlugFromName(name),
		Description:  "Mô tả cho " + name,
		ImageURL:     "https://example.com/categories/" + string(rune(64+id)) + ".jpg",
		ProductCount: 10,
		IsActive:     true,
		IsVisible:    true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// createSlugFromName tạo slug từ tên
func createSlugFromName(name string) string {
	// Simple implementation for tests
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}

// TestGetCategoryByID kiểm tra chức năng lấy danh mục theo ID
func TestGetCategoryByID(t *testing.T) {
	// Setup
	mockRepo, categoryService := createMockCategoryService()

	// Tạo danh mục mẫu
	sampleCategory := createSampleCategory(1, "Electronics")

	// Test case 1: Lấy danh mục thành công
	t.Run("Get category successfully", func(t *testing.T) {
		// Đặt expectation
		mockRepo.On("GetCategoryByID", 1).Return(sampleCategory, nil).Once()

		// Thực hiện
		category, err := categoryService.GetCategoryByID(1)

		// Kiểm tra
		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, sampleCategory.ID, category.ID)
		assert.Equal(t, sampleCategory.Name, category.Name)
		assert.Equal(t, sampleCategory.Slug, category.Slug)
		assert.Equal(t, sampleCategory.ProductCount, category.ProductCount)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Danh mục không tồn tại
	t.Run("Category not found", func(t *testing.T) {
		// Đặt expectation
		mockRepo.On("GetCategoryByID", 999).Return(nil, domain.ErrCategoryNotFound).Once()

		// Thực hiện
		category, err := categoryService.GetCategoryByID(999)

		// Kiểm tra
		assert.Error(t, err)
		assert.Nil(t, category)
		assert.Equal(t, domain.ErrCategoryNotFound, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestGetCategoryBySlug kiểm tra chức năng lấy danh mục theo slug
func TestGetCategoryBySlug(t *testing.T) {
	// Setup
	mockRepo, categoryService := createMockCategoryService()

	// Tạo danh mục mẫu
	sampleCategory := createSampleCategory(1, "Electronics")

	// Test case 1: Lấy danh mục thành công
	t.Run("Get category by slug successfully", func(t *testing.T) {
		// Đặt expectation
		mockRepo.On("GetCategoryBySlug", "electronics").Return(sampleCategory, nil).Once()

		// Thực hiện
		category, err := categoryService.GetCategoryBySlug("electronics")

		// Kiểm tra
		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, sampleCategory.ID, category.ID)
		assert.Equal(t, sampleCategory.Name, category.Name)
		assert.Equal(t, "electronics", category.Slug)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Danh mục không tồn tại
	t.Run("Category slug not found", func(t *testing.T) {
		// Đặt expectation
		mockRepo.On("GetCategoryBySlug", "nonexistent").Return(nil, domain.ErrCategoryNotFound).Once()

		// Thực hiện
		category, err := categoryService.GetCategoryBySlug("nonexistent")

		// Kiểm tra
		assert.Error(t, err)
		assert.Nil(t, category)
		assert.Equal(t, domain.ErrCategoryNotFound, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestGetAllCategories kiểm tra chức năng lấy tất cả danh mục
func TestGetAllCategories(t *testing.T) {
	// Setup
	mockRepo, categoryService := createMockCategoryService()

	// Tạo danh sách danh mục mẫu
	sampleCategories := []*domain.Category{
		createSampleCategory(1, "Electronics"),
		createSampleCategory(2, "Clothing"),
		createSampleCategory(3, "Books"),
	}

	// Test case 1: Lấy tất cả danh mục thành công
	t.Run("Get all categories successfully", func(t *testing.T) {
		// Tạo filter rỗng
		filters := map[string]string{}

		// Đặt expectation
		mockRepo.On("GetAllCategories", filters).Return(sampleCategories, nil).Once()

		// Thực hiện
		categories, err := categoryService.GetAllCategories(filters)

		// Kiểm tra
		assert.NoError(t, err)
		assert.NotNil(t, categories)
		assert.Len(t, categories, 3)
		assert.Equal(t, "Electronics", categories[0].Name)
		assert.Equal(t, "Clothing", categories[1].Name)
		assert.Equal(t, "Books", categories[2].Name)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Lấy danh mục có filter
	t.Run("Get categories with filter", func(t *testing.T) {
		// Tạo filter
		filters := map[string]string{"active": "true"}

		// Kết quả lọc - chỉ danh mục active
		filteredCategories := []*domain.Category{
			createSampleCategory(1, "Electronics"),
			createSampleCategory(2, "Clothing"),
		}

		// Đặt expectation
		mockRepo.On("GetAllCategories", filters).Return(filteredCategories, nil).Once()

		// Thực hiện
		categories, err := categoryService.GetAllCategories(filters)

		// Kiểm tra
		assert.NoError(t, err)
		assert.NotNil(t, categories)
		assert.Len(t, categories, 2)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: Lỗi database
	t.Run("Database error", func(t *testing.T) {
		// Tạo filter rỗng
		filters := map[string]string{}

		// Đặt expectation - giả lập lỗi database
		mockRepo.On("GetAllCategories", filters).Return(nil, errors.New("database error")).Once()

		// Thực hiện
		categories, err := categoryService.GetAllCategories(filters)

		// Kiểm tra
		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestCreateCategory kiểm tra chức năng tạo danh mục mới
func TestCreateCategory(t *testing.T) {
	// Setup
	mockRepo, categoryService := createMockCategoryService()

	// Test case 1: Tạo danh mục thành công
	t.Run("Create category successfully", func(t *testing.T) {
		// Tạo danh mục để thêm
		newCategory := &domain.Category{
			Name:        "New Category",
			Description: "Description for new category",
			IsActive:    true,
			IsVisible:   true,
		}

		// Đặt expectation
		mockRepo.On("CreateCategory", mock.AnythingOfType("*domain.Category")).
			Run(func(args mock.Arguments) {
				// Kiểm tra category được truyền vào
				cat := args.Get(0).(*domain.Category)
				assert.Equal(t, "New Category", cat.Name)
				assert.Equal(t, "new-category", cat.Slug) // Slug được tạo tự động
			}).
			Return(123, nil).Once()

		// Thực hiện
		id, err := categoryService.CreateCategory(newCategory)

		// Kiểm tra
		assert.NoError(t, err)
		assert.Equal(t, 123, id)
		assert.Equal(t, "new-category", newCategory.Slug) // Slug được tạo tự động

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Lỗi validation - tên rỗng
	t.Run("Validation error - empty name", func(t *testing.T) {
		// Tạo danh mục không hợp lệ
		invalidCategory := &domain.Category{
			Name:        "", // Tên rỗng
			Description: "Some description",
			IsActive:    true,
		}

		// Thực hiện
		id, err := categoryService.CreateCategory(invalidCategory)

		// Kiểm tra
		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.Equal(t, domain.ErrInvalidCategory, err)

		// Verify expectations - không gọi repository.CreateCategory
		mockRepo.AssertNotCalled(t, "CreateCategory")
	})
}

// TestUpdateCategory kiểm tra chức năng cập nhật danh mục
func TestUpdateCategory(t *testing.T) {
	// Setup
	mockRepo, categoryService := createMockCategoryService()

	// Test case 1: Cập nhật danh mục thành công
	t.Run("Update category successfully", func(t *testing.T) {
		// Tạo danh mục để cập nhật
		updateCategory := &domain.Category{
			ID:          1,
			Name:        "Updated Category",
			Description: "Updated description",
			IsActive:    true,
			IsVisible:   true,
		}

		// Đặt expectation
		mockRepo.On("UpdateCategory", mock.AnythingOfType("*domain.Category")).
			Run(func(args mock.Arguments) {
				// Kiểm tra category được truyền vào
				cat := args.Get(0).(*domain.Category)
				assert.Equal(t, 1, cat.ID)
				assert.Equal(t, "Updated Category", cat.Name)
				assert.Equal(t, "updated-category", cat.Slug) // Slug được tạo tự động
			}).
			Return(nil).Once()

		// Thực hiện
		err := categoryService.UpdateCategory(updateCategory)

		// Kiểm tra
		assert.NoError(t, err)
		assert.Equal(t, "updated-category", updateCategory.Slug) // Slug được tạo tự động

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Lỗi validation - ID không hợp lệ
	t.Run("Validation error - invalid ID", func(t *testing.T) {
		// Tạo danh mục không hợp lệ
		invalidCategory := &domain.Category{
			ID:          0, // ID không hợp lệ
			Name:        "Some Category",
			Description: "Some description",
			IsActive:    true,
		}

		// Thực hiện
		err := categoryService.UpdateCategory(invalidCategory)

		// Kiểm tra
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidCategory, err)

		// Verify expectations - không gọi repository.UpdateCategory
		mockRepo.AssertNotCalled(t, "UpdateCategory")
	})

	// Test case 3: Lỗi validation - tên rỗng
	t.Run("Validation error - empty name", func(t *testing.T) {
		// Tạo danh mục không hợp lệ
		invalidCategory := &domain.Category{
			ID:          1,
			Name:        "", // Tên rỗng
			Description: "Some description",
			IsActive:    true,
		}

		// Thực hiện
		err := categoryService.UpdateCategory(invalidCategory)

		// Kiểm tra
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidCategory, err)

		// Verify expectations - không gọi repository.UpdateCategory
		mockRepo.AssertNotCalled(t, "UpdateCategory")
	})
}

// TestDeleteCategory kiểm tra chức năng xóa danh mục
func TestDeleteCategory(t *testing.T) {
	// Setup
	mockRepo, categoryService := createMockCategoryService()

	// Test case 1: Xóa danh mục thành công
	t.Run("Delete category successfully", func(t *testing.T) {
		// Đặt expectation
		mockRepo.On("DeleteCategory", 1).Return(nil).Once()

		// Thực hiện
		err := categoryService.DeleteCategory(1)

		// Kiểm tra
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Danh mục không tồn tại
	t.Run("Category not found", func(t *testing.T) {
		// Đặt expectation
		mockRepo.On("DeleteCategory", 999).Return(domain.ErrCategoryNotFound).Once()

		// Thực hiện
		err := categoryService.DeleteCategory(999)

		// Kiểm tra
		assert.Error(t, err)
		assert.Equal(t, domain.ErrCategoryNotFound, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

// TestSyncProductCounts kiểm tra chức năng đồng bộ số lượng sản phẩm trong danh mục
func TestSyncProductCounts(t *testing.T) {
	// Setup
	mockRepo, categoryService := createMockCategoryService()

	// Test case 1: Đồng bộ thành công
	t.Run("Sync product counts successfully", func(t *testing.T) {
		// Đặt expectation
		mockRepo.On("SyncProductCounts").Return(nil).Once()

		// Thực hiện
		err := categoryService.SyncProductCounts()

		// Kiểm tra
		assert.NoError(t, err)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Lỗi khi đồng bộ
	t.Run("Error during sync", func(t *testing.T) {
		// Đặt expectation
		mockRepo.On("SyncProductCounts").Return(errors.New("database error")).Once()

		// Thực hiện
		err := categoryService.SyncProductCounts()

		// Kiểm tra
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

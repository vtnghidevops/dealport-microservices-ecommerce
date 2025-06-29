package helpers

import (
	"product-service/internal/domain"
	"product-service/internal/service"
	"product-service/tests/mocks"
	"time"
)

// CreateMockProductService tạo một product service với repository mock
func CreateMockProductService() (*mocks.MockProductRepository, domain.ProductService) {
	mockRepo := new(mocks.MockProductRepository)
	productService := service.NewProductService(mockRepo)
	return mockRepo, productService
}

// CreateSampleProduct tạo một sản phẩm mẫu với ID và tên cho trước
func CreateSampleProduct(id int, name string) *domain.Product {
	now := time.Now()
	return &domain.Product{
		ID:            id,
		Name:          name,
		Description:   "Mô tả sản phẩm " + name,
		Slug:          createSlugFromName(name),
		Price:         99.99,
		CategoryID:    1,
		CategorySlug:  "electronics",
		StockQuantity: 100,
		Brand:         "TestBrand",
		Tags:          []string{"test", "sample"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// CreateSampleProducts tạo một danh sách sản phẩm mẫu
func CreateSampleProducts(count int) []*domain.Product {
	products := make([]*domain.Product, 0, count)
	for i := 1; i <= count; i++ {
		products = append(products, CreateSampleProduct(i, "Product "+string(rune(64+i))))
	}
	return products
}

// CreateSampleProductImage tạo một hình ảnh sản phẩm mẫu
func CreateSampleProductImage(id int, productID int, isPrimary bool) domain.ProductImage {
	return domain.ProductImage{
		ID:           id,
		ProductID:    productID,
		URL:          "https://example.com/products/" + string(rune(64+id)) + ".jpg",
		IsPrimary:    isPrimary,
		DisplayOrder: id,
		CreatedAt:    time.Now(),
	}
}

// CreateSampleProductImages tạo một danh sách hình ảnh sản phẩm mẫu
func CreateSampleProductImages(productID int, count int) []domain.ProductImage {
	images := make([]domain.ProductImage, 0, count)
	for i := 1; i <= count; i++ {
		isPrimary := i == 1 // Hình đầu tiên là hình chính
		images = append(images, CreateSampleProductImage(i, productID, isPrimary))
	}
	return images
}

// CreateSampleProductReview tạo một đánh giá sản phẩm mẫu
func CreateSampleProductReview(id int, productID int, userID string, rating float64) *domain.ProductReview {
	return &domain.ProductReview{
		ID:        id,
		ProductID: productID,
		UserID:    userID,
		UserName:  "User " + userID,
		Rating:    rating,
		Comment:   "Đánh giá sản phẩm " + string(rune(64+id)),
		CreatedAt: time.Now(),
	}
}

// CreateSampleProductReviews tạo một danh sách đánh giá sản phẩm mẫu
func CreateSampleProductReviews(productID int, count int) []*domain.ProductReview {
	reviews := make([]*domain.ProductReview, 0, count)
	for i := 1; i <= count; i++ {
		// Rating từ 1-5 sao
		rating := float64(i%5 + 1)
		reviews = append(reviews, CreateSampleProductReview(i, productID, "user"+string(rune(64+i)), rating))
	}
	return reviews
}

// CreateSampleTestimonials tạo một danh sách testimonials mẫu
func CreateSampleTestimonials(count int) []*domain.Testimonial {
	testimonials := make([]*domain.Testimonial, 0, count)
	for i := 1; i <= count; i++ {
		userID := "user" + string(rune(64+i))
		testimonials = append(testimonials, &domain.Testimonial{
			ID:        userID,
			UserName:  "User " + userID,
			Avatar:    "https://example.com/avatars/" + userID + ".jpg",
			Review:    "Trải nghiệm tuyệt vời với sản phẩm!",
			Rating:    5.0, // Testimonial thường là đánh giá cao
			ProductID: i,
		})
	}
	return testimonials
}

// createSlugFromName tạo slug từ tên sản phẩm
func createSlugFromName(name string) string {
	// Đơn giản hóa: chỉ thay thế khoảng trắng bằng gạch ngang và chuyển sang chữ thường
	slug := ""
	for _, ch := range name {
		if ch == ' ' {
			slug += "-"
		} else {
			// Chỉ giữ lại chữ cái và số, chuyển sang lowercase
			if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
				if ch >= 'A' && ch <= 'Z' {
					slug += string(ch + 32) // Chuyển sang chữ thường
				} else {
					slug += string(ch)
				}
			}
		}
	}
	return slug
}

// CreateMockFileUpload tạo một đối tượng giả lập để test upload file
func CreateMockFileUpload(filename string, size int64) *mocks.MockFileUpload {
	mockFileUpload := new(mocks.MockFileUpload)
	mockFileUpload.On("Filename").Return(filename)
	mockFileUpload.On("Size").Return(size)
	// Mock file sẽ được setup trong từng test case cụ thể
	return mockFileUpload
} 
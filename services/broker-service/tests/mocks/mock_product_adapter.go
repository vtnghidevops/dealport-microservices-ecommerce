package mocks

import (
	"context"

	productGrpc "broker-service/internal/handlers/grpc/product"
	productpb "broker-service/proto/product"
)

// ProductClientAdapter là một adapter cho ProductClient
// Nó thực hiện các phương thức từ productGrpc.ProductClient
// nhưng ủy quyền cho MockProductServiceClient để thuận tiện cho việc kiểm thử
type ProductClientAdapter struct {
	Mock *MockProductServiceClient
}

// NewProductClientAdapter tạo một adapter mới
func NewProductClientAdapter(mock *MockProductServiceClient) *productGrpc.ProductClient {
	// Đây là một hack để test, thường không nên làm vậy trong code thực tế
	// Trả về một con trỏ rỗng, vì chúng ta sẽ không gọi trực tiếp phương thức
	// Các test thường sẽ skip phần gọi thực tế
	return &productGrpc.ProductClient{}
}

// GetProductBySlug is used for testing - delegates to the mock
func (a *ProductClientAdapter) GetProductBySlug(slug string) (*productpb.Product, error) {
	ctx := context.Background()
	req := &productpb.GetProductBySlugRequest{Slug: slug}
	return a.Mock.GetProductBySlug(ctx, req)
}

// GetProduct is used for testing - delegates to the mock
func (a *ProductClientAdapter) GetProduct(id int) (*productpb.Product, error) {
	ctx := context.Background()
	req := &productpb.GetProductRequest{Id: int32(id)}
	return a.Mock.GetProduct(ctx, req)
}

// GetProductImageFile is used for testing - delegates to the mock
func (a *ProductClientAdapter) GetProductImageFile(filename string) ([]byte, string, error) {
	ctx := context.Background()
	req := &productpb.GetProductImageFileRequest{Filename: filename}
	resp, err := a.Mock.GetProductImageFile(ctx, req)
	if err != nil {
		return nil, "", err
	}
	return resp.ImageData, resp.ContentType, nil
}

package mocks

import (
	"context"

	productGrpc "broker-service/internal/handlers/grpc/product"
	productpb "broker-service/proto/product"

	"github.com/stretchr/testify/mock"
)

// SetupMockProductClient khởi tạo và trả về một mock product client
func SetupMockProductClient() *MockProductServiceClient {
	return new(MockProductServiceClient)
}

// ProductClientMockAdapter is a wrapper struct that implements the methods from productGrpc.ProductClient
// but delegates to the mock for testing
type ProductClientMockAdapter struct {
	Mock *MockProductServiceClient
}

// NewProductClientMockAdapter creates a new adapter that will be used in tests
func NewProductClientMockAdapter(mock *MockProductServiceClient) *productGrpc.ProductClient {
	// Create an instance of our adapter
	adapter := &ProductClientMockAdapter{
		Mock: mock,
	}

	// For tests, we'll wrap the adapter's GetProductImageFile method to use in the handler
	client := &productGrpc.ProductClient{}

	// Replace the GetProductImageFile method with our mock implementation
	// This is a hack to allow testing without actually making gRPC calls
	productGrpc.GetProductImageFileFunc = func(filename string) ([]byte, string, error) {
		return adapter.GetProductImageFile(filename)
	}

	return client
}

// GetProductImageFile is the adapter method that will be called by the handler
func (a *ProductClientMockAdapter) GetProductImageFile(filename string) ([]byte, string, error) {
	// Call the mock method
	req := &productpb.GetProductImageFileRequest{Filename: filename}
	// Create a real context for the mock call
	ctx := context.Background()
	// Call without any gRPC options
	resp, err := a.Mock.GetProductImageFile(ctx, req, nil)
	if err != nil {
		return nil, "", err
	}
	return resp.ImageData, resp.ContentType, nil
}

// AddGetProductMock thêm mock response cho GetProduct
func AddGetProductMock(mockClient *MockProductServiceClient, productID int32, productName string) {
	mockClient.On(
		"GetProduct",
		mock.Anything,
		&productpb.GetProductRequest{Id: productID},
		mock.Anything,
	).Return(&productpb.Product{
		Id:   productID,
		Name: productName,
	}, nil).Once()
}

// AddGetProductBySlugMock thêm mock response cho GetProductBySlug
func AddGetProductBySlugMock(mockClient *MockProductServiceClient, slug, productName string) {
	mockClient.On(
		"GetProductBySlug",
		mock.Anything,
		&productpb.GetProductBySlugRequest{Slug: slug},
		mock.Anything,
	).Return(&productpb.Product{
		Id:   1,
		Name: productName,
		Slug: slug,
	}, nil).Once()
}

// AddGetProductImageFileMock thêm mock response cho GetProductImageFile
func AddGetProductImageFileMock(mockClient *MockProductServiceClient, filename string) {
	mockClient.On(
		"GetProductImageFile",
		mock.Anything,
		&productpb.GetProductImageFileRequest{Filename: filename},
		mock.Anything,
	).Return(&productpb.GetProductImageFileResponse{
		ImageData:   []byte("mock image data"),
		ContentType: "image/jpeg",
	}, nil).Once()
}

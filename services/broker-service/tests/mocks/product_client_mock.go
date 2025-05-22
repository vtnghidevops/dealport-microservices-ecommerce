package mocks

import (
	"context"

	productpb "broker-service/proto/product"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockProductServiceClient is a mock for the ProductServiceClient interface
type MockProductServiceClient struct {
	mock.Mock
}

// GetProduct mocks the GetProduct method
func (m *MockProductServiceClient) GetProduct(ctx context.Context, req *productpb.GetProductRequest, opts ...grpc.CallOption) (*productpb.Product, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.Product
	if response != nil {
		responseValue = response.(*productpb.Product)
	}
	return responseValue, args.Error(1)
}

// GetProductBySlug mocks the GetProductBySlug method
func (m *MockProductServiceClient) GetProductBySlug(ctx context.Context, req *productpb.GetProductBySlugRequest, opts ...grpc.CallOption) (*productpb.Product, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.Product
	if response != nil {
		responseValue = response.(*productpb.Product)
	}
	return responseValue, args.Error(1)
}

// ListProducts mocks the ListProducts method
func (m *MockProductServiceClient) ListProducts(ctx context.Context, req *productpb.ListProductsRequest, opts ...grpc.CallOption) (*productpb.ListProductsResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.ListProductsResponse
	if response != nil {
		responseValue = response.(*productpb.ListProductsResponse)
	}
	return responseValue, args.Error(1)
}

// CreateProduct mocks the CreateProduct method
func (m *MockProductServiceClient) CreateProduct(ctx context.Context, req *productpb.Product, opts ...grpc.CallOption) (*productpb.CreateProductResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.CreateProductResponse
	if response != nil {
		responseValue = response.(*productpb.CreateProductResponse)
	}
	return responseValue, args.Error(1)
}

// UpdateProduct mocks the UpdateProduct method
func (m *MockProductServiceClient) UpdateProduct(ctx context.Context, req *productpb.Product, opts ...grpc.CallOption) (*productpb.StatusResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.StatusResponse
	if response != nil {
		responseValue = response.(*productpb.StatusResponse)
	}
	return responseValue, args.Error(1)
}

// PatchProduct mocks the PatchProduct method
func (m *MockProductServiceClient) PatchProduct(ctx context.Context, req *productpb.PatchProductRequest, opts ...grpc.CallOption) (*productpb.Product, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.Product
	if response != nil {
		responseValue = response.(*productpb.Product)
	}
	return responseValue, args.Error(1)
}

// DeleteProduct mocks the DeleteProduct method
func (m *MockProductServiceClient) DeleteProduct(ctx context.Context, req *productpb.DeleteProductRequest, opts ...grpc.CallOption) (*productpb.StatusResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.StatusResponse
	if response != nil {
		responseValue = response.(*productpb.StatusResponse)
	}
	return responseValue, args.Error(1)
}

// GetProductReviews mocks the GetProductReviews method
func (m *MockProductServiceClient) GetProductReviews(ctx context.Context, req *productpb.GetProductReviewsRequest, opts ...grpc.CallOption) (*productpb.GetProductReviewsResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.GetProductReviewsResponse
	if response != nil {
		responseValue = response.(*productpb.GetProductReviewsResponse)
	}
	return responseValue, args.Error(1)
}

// AddProductReview mocks the AddProductReview method
func (m *MockProductServiceClient) AddProductReview(ctx context.Context, req *productpb.ProductReview, opts ...grpc.CallOption) (*productpb.AddProductReviewResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.AddProductReviewResponse
	if response != nil {
		responseValue = response.(*productpb.AddProductReviewResponse)
	}
	return responseValue, args.Error(1)
}

// UpdateProductReview mocks the UpdateProductReview method
func (m *MockProductServiceClient) UpdateProductReview(ctx context.Context, req *productpb.ProductReview, opts ...grpc.CallOption) (*productpb.StatusResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.StatusResponse
	if response != nil {
		responseValue = response.(*productpb.StatusResponse)
	}
	return responseValue, args.Error(1)
}

// GetProductImages mocks the GetProductImages method
func (m *MockProductServiceClient) GetProductImages(ctx context.Context, req *productpb.GetProductImagesRequest, opts ...grpc.CallOption) (*productpb.GetProductImagesResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.GetProductImagesResponse
	if response != nil {
		responseValue = response.(*productpb.GetProductImagesResponse)
	}
	return responseValue, args.Error(1)
}

// UploadProductImage mocks the UploadProductImage method
func (m *MockProductServiceClient) UploadProductImage(ctx context.Context, req *productpb.UploadProductImageRequest, opts ...grpc.CallOption) (*productpb.UploadProductImageResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.UploadProductImageResponse
	if response != nil {
		responseValue = response.(*productpb.UploadProductImageResponse)
	}
	return responseValue, args.Error(1)
}

// DeleteProductImage mocks the DeleteProductImage method
func (m *MockProductServiceClient) DeleteProductImage(ctx context.Context, req *productpb.DeleteProductImageRequest, opts ...grpc.CallOption) (*productpb.StatusResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.StatusResponse
	if response != nil {
		responseValue = response.(*productpb.StatusResponse)
	}
	return responseValue, args.Error(1)
}

// SetPrimaryProductImage mocks the SetPrimaryProductImage method
func (m *MockProductServiceClient) SetPrimaryProductImage(ctx context.Context, req *productpb.SetPrimaryProductImageRequest, opts ...grpc.CallOption) (*productpb.StatusResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.StatusResponse
	if response != nil {
		responseValue = response.(*productpb.StatusResponse)
	}
	return responseValue, args.Error(1)
}

// GetProductImageFile mocks the GetProductImageFile method
func (m *MockProductServiceClient) GetProductImageFile(ctx context.Context, req *productpb.GetProductImageFileRequest, opts ...grpc.CallOption) (*productpb.GetProductImageFileResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.GetProductImageFileResponse
	if response != nil {
		responseValue = response.(*productpb.GetProductImageFileResponse)
	}
	return responseValue, args.Error(1)
}

// GetTopRatedTestimonials mocks the GetTopRatedTestimonials method
func (m *MockProductServiceClient) GetTopRatedTestimonials(ctx context.Context, req *productpb.GetTopRatedTestimonialsRequest, opts ...grpc.CallOption) (*productpb.GetTopRatedTestimonialsResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.GetTopRatedTestimonialsResponse
	if response != nil {
		responseValue = response.(*productpb.GetTopRatedTestimonialsResponse)
	}
	return responseValue, args.Error(1)
}

// GetHealth mocks the GetHealth method
func (m *MockProductServiceClient) GetHealth(ctx context.Context, req *productpb.HealthRequest, opts ...grpc.CallOption) (*productpb.HealthResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *productpb.HealthResponse
	if response != nil {
		responseValue = response.(*productpb.HealthResponse)
	}
	return responseValue, args.Error(1)
}

package mocks

import (
	"context"

	cartpb "broker-service/proto/cart"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockCartServiceClient is a mock for the CartServiceClient interface
type MockCartServiceClient struct {
	mock.Mock
}

// GetCart mocks the GetCart method
func (m *MockCartServiceClient) GetCart(ctx context.Context, req *cartpb.GetCartRequest, opts ...grpc.CallOption) (*cartpb.Cart, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *cartpb.Cart
	if response != nil {
		responseValue = response.(*cartpb.Cart)
	}
	return responseValue, args.Error(1)
}

// AddCartItem mocks the AddCartItem method
func (m *MockCartServiceClient) AddCartItem(ctx context.Context, req *cartpb.AddCartItemRequest, opts ...grpc.CallOption) (*cartpb.Cart, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *cartpb.Cart
	if response != nil {
		responseValue = response.(*cartpb.Cart)
	}
	return responseValue, args.Error(1)
}

// UpdateCartItem mocks the UpdateCartItem method
func (m *MockCartServiceClient) UpdateCartItem(ctx context.Context, req *cartpb.UpdateCartItemRequest, opts ...grpc.CallOption) (*cartpb.Cart, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *cartpb.Cart
	if response != nil {
		responseValue = response.(*cartpb.Cart)
	}
	return responseValue, args.Error(1)
}

// RemoveCartItem mocks the RemoveCartItem method
func (m *MockCartServiceClient) RemoveCartItem(ctx context.Context, req *cartpb.RemoveCartItemRequest, opts ...grpc.CallOption) (*cartpb.Cart, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *cartpb.Cart
	if response != nil {
		responseValue = response.(*cartpb.Cart)
	}
	return responseValue, args.Error(1)
}

// ClearCart mocks the ClearCart method
func (m *MockCartServiceClient) ClearCart(ctx context.Context, req *cartpb.ClearCartRequest, opts ...grpc.CallOption) (*cartpb.StatusResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *cartpb.StatusResponse
	if response != nil {
		responseValue = response.(*cartpb.StatusResponse)
	}
	return responseValue, args.Error(1)
}

// ApplyCoupon mocks the ApplyCoupon method
func (m *MockCartServiceClient) ApplyCoupon(ctx context.Context, req *cartpb.ApplyCouponRequest, opts ...grpc.CallOption) (*cartpb.Cart, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *cartpb.Cart
	if response != nil {
		responseValue = response.(*cartpb.Cart)
	}
	return responseValue, args.Error(1)
}

// RemoveCoupon mocks the RemoveCoupon method
func (m *MockCartServiceClient) RemoveCoupon(ctx context.Context, req *cartpb.RemoveCouponRequest, opts ...grpc.CallOption) (*cartpb.Cart, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *cartpb.Cart
	if response != nil {
		responseValue = response.(*cartpb.Cart)
	}
	return responseValue, args.Error(1)
}

// GetHealth mocks the GetHealth method
func (m *MockCartServiceClient) GetHealth(ctx context.Context, req *cartpb.HealthRequest, opts ...grpc.CallOption) (*cartpb.HealthResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *cartpb.HealthResponse
	if response != nil {
		responseValue = response.(*cartpb.HealthResponse)
	}
	return responseValue, args.Error(1)
}

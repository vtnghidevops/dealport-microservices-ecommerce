package mocks

import (
	"context"

	checkoutpb "broker-service/proto/checkout"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockCheckoutServiceClient is a mock for the CheckoutServiceClient interface
type MockCheckoutServiceClient struct {
	mock.Mock
}

// CreateOrder mocks the CreateOrder method
func (m *MockCheckoutServiceClient) CreateOrder(ctx context.Context, req *checkoutpb.CreateOrderRequest, opts ...grpc.CallOption) (*checkoutpb.Order, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *checkoutpb.Order
	if response != nil {
		responseValue = response.(*checkoutpb.Order)
	}
	return responseValue, args.Error(1)
}

// GetOrder mocks the GetOrder method
func (m *MockCheckoutServiceClient) GetOrder(ctx context.Context, req *checkoutpb.GetOrderRequest, opts ...grpc.CallOption) (*checkoutpb.Order, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *checkoutpb.Order
	if response != nil {
		responseValue = response.(*checkoutpb.Order)
	}
	return responseValue, args.Error(1)
}

// ListOrders mocks the ListOrders method
func (m *MockCheckoutServiceClient) ListOrders(ctx context.Context, req *checkoutpb.ListOrdersRequest, opts ...grpc.CallOption) (*checkoutpb.ListOrdersResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *checkoutpb.ListOrdersResponse
	if response != nil {
		responseValue = response.(*checkoutpb.ListOrdersResponse)
	}
	return responseValue, args.Error(1)
}

// UpdateOrderStatus mocks the UpdateOrderStatus method
func (m *MockCheckoutServiceClient) UpdateOrderStatus(ctx context.Context, req *checkoutpb.UpdateOrderStatusRequest, opts ...grpc.CallOption) (*checkoutpb.StatusResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *checkoutpb.StatusResponse
	if response != nil {
		responseValue = response.(*checkoutpb.StatusResponse)
	}
	return responseValue, args.Error(1)
}

// ProcessPayment mocks the ProcessPayment method
func (m *MockCheckoutServiceClient) ProcessPayment(ctx context.Context, req *checkoutpb.ProcessPaymentRequest, opts ...grpc.CallOption) (*checkoutpb.PaymentResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *checkoutpb.PaymentResponse
	if response != nil {
		responseValue = response.(*checkoutpb.PaymentResponse)
	}
	return responseValue, args.Error(1)
}

// ValidateCheckout mocks the ValidateCheckout method
func (m *MockCheckoutServiceClient) ValidateCheckout(ctx context.Context, req *checkoutpb.ValidateCheckoutRequest, opts ...grpc.CallOption) (*checkoutpb.ValidateCheckoutResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *checkoutpb.ValidateCheckoutResponse
	if response != nil {
		responseValue = response.(*checkoutpb.ValidateCheckoutResponse)
	}
	return responseValue, args.Error(1)
}

// GetUserTotalSpend mocks the GetUserTotalSpend method
func (m *MockCheckoutServiceClient) GetUserTotalSpend(ctx context.Context, req *checkoutpb.UserRequest, opts ...grpc.CallOption) (*checkoutpb.TotalSpendResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *checkoutpb.TotalSpendResponse
	if response != nil {
		responseValue = response.(*checkoutpb.TotalSpendResponse)
	}
	return responseValue, args.Error(1)
}

// GetUserOrderCount mocks the GetUserOrderCount method
func (m *MockCheckoutServiceClient) GetUserOrderCount(ctx context.Context, req *checkoutpb.UserRequest, opts ...grpc.CallOption) (*checkoutpb.OrderCountResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *checkoutpb.OrderCountResponse
	if response != nil {
		responseValue = response.(*checkoutpb.OrderCountResponse)
	}
	return responseValue, args.Error(1)
}

// GetHealth mocks the GetHealth method
func (m *MockCheckoutServiceClient) GetHealth(ctx context.Context, req *checkoutpb.HealthRequest, opts ...grpc.CallOption) (*checkoutpb.HealthResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *checkoutpb.HealthResponse
	if response != nil {
		responseValue = response.(*checkoutpb.HealthResponse)
	}
	return responseValue, args.Error(1)
}

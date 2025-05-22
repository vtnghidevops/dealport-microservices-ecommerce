package mocks

import (
	"context"

	couponpb "broker-service/proto/coupon"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockCouponServiceClient is a mock for the CouponServiceClient interface
type MockCouponServiceClient struct {
	mock.Mock
}

// GetCoupons mocks the GetCoupons method
func (m *MockCouponServiceClient) GetCoupons(ctx context.Context, req *couponpb.GetCouponsRequest, opts ...grpc.CallOption) (*couponpb.GetCouponsResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *couponpb.GetCouponsResponse
	if response != nil {
		responseValue = response.(*couponpb.GetCouponsResponse)
	}
	return responseValue, args.Error(1)
}

// GetCouponByID mocks the GetCouponByID method
func (m *MockCouponServiceClient) GetCouponByID(ctx context.Context, req *couponpb.GetCouponByIDRequest, opts ...grpc.CallOption) (*couponpb.GetCouponResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *couponpb.GetCouponResponse
	if response != nil {
		responseValue = response.(*couponpb.GetCouponResponse)
	}
	return responseValue, args.Error(1)
}

// GetCouponByCode mocks the GetCouponByCode method
func (m *MockCouponServiceClient) GetCouponByCode(ctx context.Context, req *couponpb.GetCouponByCodeRequest, opts ...grpc.CallOption) (*couponpb.GetCouponResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *couponpb.GetCouponResponse
	if response != nil {
		responseValue = response.(*couponpb.GetCouponResponse)
	}
	return responseValue, args.Error(1)
}

// CreateCoupon mocks the CreateCoupon method
func (m *MockCouponServiceClient) CreateCoupon(ctx context.Context, req *couponpb.CreateCouponRequest, opts ...grpc.CallOption) (*couponpb.GetCouponResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *couponpb.GetCouponResponse
	if response != nil {
		responseValue = response.(*couponpb.GetCouponResponse)
	}
	return responseValue, args.Error(1)
}

// UpdateCoupon mocks the UpdateCoupon method
func (m *MockCouponServiceClient) UpdateCoupon(ctx context.Context, req *couponpb.UpdateCouponRequest, opts ...grpc.CallOption) (*couponpb.GetCouponResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *couponpb.GetCouponResponse
	if response != nil {
		responseValue = response.(*couponpb.GetCouponResponse)
	}
	return responseValue, args.Error(1)
}

// DeleteCoupon mocks the DeleteCoupon method
func (m *MockCouponServiceClient) DeleteCoupon(ctx context.Context, req *couponpb.DeleteCouponRequest, opts ...grpc.CallOption) (*couponpb.DeleteCouponResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *couponpb.DeleteCouponResponse
	if response != nil {
		responseValue = response.(*couponpb.DeleteCouponResponse)
	}
	return responseValue, args.Error(1)
}

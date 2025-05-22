package mocks

import (
	"context"

	paymentpb "broker-service/proto/payment"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockPaymentServiceClient is a mock for the PaymentServiceClient interface
type MockPaymentServiceClient struct {
	mock.Mock
}

// CreateMomoPayment mocks the CreateMomoPayment method
func (m *MockPaymentServiceClient) CreateMomoPayment(ctx context.Context, req *paymentpb.MomoPaymentRequest, opts ...grpc.CallOption) (*paymentpb.MomoPaymentResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *paymentpb.MomoPaymentResponse
	if response != nil {
		responseValue = response.(*paymentpb.MomoPaymentResponse)
	}
	return responseValue, args.Error(1)
}

// VerifyMomoPayment mocks the VerifyMomoPayment method
func (m *MockPaymentServiceClient) VerifyMomoPayment(ctx context.Context, req *paymentpb.MomoVerifyRequest, opts ...grpc.CallOption) (*paymentpb.PaymentVerifyResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *paymentpb.PaymentVerifyResponse
	if response != nil {
		responseValue = response.(*paymentpb.PaymentVerifyResponse)
	}
	return responseValue, args.Error(1)
}

// ProcessMomoCallback mocks the ProcessMomoCallback method
func (m *MockPaymentServiceClient) ProcessMomoCallback(ctx context.Context, req *paymentpb.MomoCallbackRequest, opts ...grpc.CallOption) (*paymentpb.PaymentVerifyResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *paymentpb.PaymentVerifyResponse
	if response != nil {
		responseValue = response.(*paymentpb.PaymentVerifyResponse)
	}
	return responseValue, args.Error(1)
}

// CreateVnpayPayment mocks the CreateVnpayPayment method
func (m *MockPaymentServiceClient) CreateVnpayPayment(ctx context.Context, req *paymentpb.VnpayPaymentRequest, opts ...grpc.CallOption) (*paymentpb.VnpayPaymentResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *paymentpb.VnpayPaymentResponse
	if response != nil {
		responseValue = response.(*paymentpb.VnpayPaymentResponse)
	}
	return responseValue, args.Error(1)
}

// VerifyVnpayPayment mocks the VerifyVnpayPayment method
func (m *MockPaymentServiceClient) VerifyVnpayPayment(ctx context.Context, req *paymentpb.VnpayVerifyRequest, opts ...grpc.CallOption) (*paymentpb.PaymentVerifyResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *paymentpb.PaymentVerifyResponse
	if response != nil {
		responseValue = response.(*paymentpb.PaymentVerifyResponse)
	}
	return responseValue, args.Error(1)
}

// CreateMomoQRPayment mocks the CreateMomoQRPayment method
func (m *MockPaymentServiceClient) CreateMomoQRPayment(ctx context.Context, req *paymentpb.MomoQRPaymentRequest, opts ...grpc.CallOption) (*paymentpb.MomoQuickPayResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *paymentpb.MomoQuickPayResponse
	if response != nil {
		responseValue = response.(*paymentpb.MomoQuickPayResponse)
	}
	return responseValue, args.Error(1)
}

// CreateMomoPosPayment mocks the CreateMomoPosPayment method
func (m *MockPaymentServiceClient) CreateMomoPosPayment(ctx context.Context, req *paymentpb.MomoPosPaymentRequest, opts ...grpc.CallOption) (*paymentpb.MomoQuickPayResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *paymentpb.MomoQuickPayResponse
	if response != nil {
		responseValue = response.(*paymentpb.MomoQuickPayResponse)
	}
	return responseValue, args.Error(1)
}

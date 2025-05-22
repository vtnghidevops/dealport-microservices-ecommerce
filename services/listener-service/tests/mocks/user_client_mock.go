package mocks

import (
	"context"
	userpb "listener-service/proto/user"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockUserServiceClient is a mock for the UserServiceClient interface
type MockUserServiceClient struct {
	mock.Mock
}

// ProcessEvent mocks the ProcessEvent method
func (m *MockUserServiceClient) ProcessEvent(ctx context.Context, in *userpb.EventRequest, opts ...grpc.CallOption) (*userpb.EventResponse, error) {
	args := m.Called(ctx, in)
	response := args.Get(0)
	var responseValue *userpb.EventResponse
	if response != nil {
		responseValue = response.(*userpb.EventResponse)
	}
	return responseValue, args.Error(1)
}

// SyncUserOrderData mocks the SyncUserOrderData method
func (m *MockUserServiceClient) SyncUserOrderData(ctx context.Context, in *userpb.SyncUserOrderDataRequest, opts ...grpc.CallOption) (*userpb.SyncUserOrderDataResponse, error) {
	args := m.Called(ctx, in)
	response := args.Get(0)
	var responseValue *userpb.SyncUserOrderDataResponse
	if response != nil {
		responseValue = response.(*userpb.SyncUserOrderDataResponse)
	}
	return responseValue, args.Error(1)
}

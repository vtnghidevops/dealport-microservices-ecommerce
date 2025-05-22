package mocks

import (
	"context"

	userpb "broker-service/proto/user"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockUserServiceClient is a mock for the UserServiceClient interface
type MockUserServiceClient struct {
	mock.Mock
}

// GetUser mocks the GetUser method
func (m *MockUserServiceClient) GetUser(ctx context.Context, req *userpb.GetUserRequest, opts ...grpc.CallOption) (*userpb.UserResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *userpb.UserResponse
	if response != nil {
		responseValue = response.(*userpb.UserResponse)
	}
	return responseValue, args.Error(1)
}

// GetUsers mocks the GetUsers method
func (m *MockUserServiceClient) GetUsers(ctx context.Context, req *userpb.GetUsersRequest, opts ...grpc.CallOption) (*userpb.GetUsersResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *userpb.GetUsersResponse
	if response != nil {
		responseValue = response.(*userpb.GetUsersResponse)
	}
	return responseValue, args.Error(1)
}

// CreateUser mocks the CreateUser method
func (m *MockUserServiceClient) CreateUser(ctx context.Context, req *userpb.CreateUserRequest, opts ...grpc.CallOption) (*userpb.UserResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *userpb.UserResponse
	if response != nil {
		responseValue = response.(*userpb.UserResponse)
	}
	return responseValue, args.Error(1)
}

// UpdateUser mocks the UpdateUser method
func (m *MockUserServiceClient) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest, opts ...grpc.CallOption) (*userpb.UserResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *userpb.UserResponse
	if response != nil {
		responseValue = response.(*userpb.UserResponse)
	}
	return responseValue, args.Error(1)
}

// DeleteUser mocks the DeleteUser method
func (m *MockUserServiceClient) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest, opts ...grpc.CallOption) (*userpb.DeleteUserResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *userpb.DeleteUserResponse
	if response != nil {
		responseValue = response.(*userpb.DeleteUserResponse)
	}
	return responseValue, args.Error(1)
}

// SearchUsers mocks the SearchUsers method
func (m *MockUserServiceClient) SearchUsers(ctx context.Context, req *userpb.SearchUsersRequest, opts ...grpc.CallOption) (*userpb.GetUsersResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *userpb.GetUsersResponse
	if response != nil {
		responseValue = response.(*userpb.GetUsersResponse)
	}
	return responseValue, args.Error(1)
}

// GetWishlist mocks the GetWishlist method
func (m *MockUserServiceClient) GetWishlist(ctx context.Context, req *userpb.GetWishlistRequest, opts ...grpc.CallOption) (*userpb.GetWishlistResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *userpb.GetWishlistResponse
	if response != nil {
		responseValue = response.(*userpb.GetWishlistResponse)
	}
	return responseValue, args.Error(1)
}

// AddToWishlist mocks the AddToWishlist method
func (m *MockUserServiceClient) AddToWishlist(ctx context.Context, req *userpb.AddToWishlistRequest, opts ...grpc.CallOption) (*userpb.WishlistResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *userpb.WishlistResponse
	if response != nil {
		responseValue = response.(*userpb.WishlistResponse)
	}
	return responseValue, args.Error(1)
}

// RemoveFromWishlist mocks the RemoveFromWishlist method
func (m *MockUserServiceClient) RemoveFromWishlist(ctx context.Context, req *userpb.RemoveFromWishlistRequest, opts ...grpc.CallOption) (*userpb.WishlistResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *userpb.WishlistResponse
	if response != nil {
		responseValue = response.(*userpb.WishlistResponse)
	}
	return responseValue, args.Error(1)
}

// GetUserStatistics mocks the GetUserStatistics method
func (m *MockUserServiceClient) GetUserStatistics(ctx context.Context, req *userpb.GetUserStatisticsRequest, opts ...grpc.CallOption) (*userpb.GetUserStatisticsResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *userpb.GetUserStatisticsResponse
	if response != nil {
		responseValue = response.(*userpb.GetUserStatisticsResponse)
	}
	return responseValue, args.Error(1)
}

// GetUserActivityChart mocks the GetUserActivityChart method
func (m *MockUserServiceClient) GetUserActivityChart(ctx context.Context, req *userpb.GetUserActivityChartRequest, opts ...grpc.CallOption) (*userpb.GetUserActivityChartResponse, error) {
	args := m.Called(ctx, req, mock.Anything)
	response := args.Get(0)
	var responseValue *userpb.GetUserActivityChartResponse
	if response != nil {
		responseValue = response.(*userpb.GetUserActivityChartResponse)
	}
	return responseValue, args.Error(1)
}

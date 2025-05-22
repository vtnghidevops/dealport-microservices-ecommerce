package mocks

import (
	"context"
	"strconv"

	userpb "broker-service/proto/user"

	"github.com/stretchr/testify/mock"
)

// UserClientAdapter là một adapter để sử dụng MockUserServiceClient trong tests
type UserClientAdapter struct {
	Mock *MockUserServiceClient
}

// NewUserClientAdapter tạo một adapter mới đóng vai trò là UserServiceClient
func NewUserClientAdapter(mock *MockUserServiceClient) userpb.UserServiceClient {
	return mock
}

// SetupMockUserClient khởi tạo và trả về một mock client đã được cấu hình
func SetupMockUserClient() *MockUserServiceClient {
	mockClient := new(MockUserServiceClient)
	return mockClient
}

// AddGetUserMock cấu hình mock cho GetUser gRPC call
func AddGetUserMock(mockClient *MockUserServiceClient, userID, email, firstName, lastName string) {
	mockClient.On(
		"GetUser",
		context.Background(),
		&userpb.GetUserRequest{Id: userID},
		mock.Anything,
	).Return(&userpb.UserResponse{
		User: &userpb.User{
			Id:        userID,
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
			Addresses: []*userpb.Address{},
		},
	}, nil).Once()
}

// AddGetWishlistMock cấu hình mock cho GetWishlist gRPC call
func AddGetWishlistMock(mockClient *MockUserServiceClient, userID string) {
	// Chuyển đổi ID sản phẩm từ string sang int32
	productID := int32(123)

	mockClient.On(
		"GetWishlist",
		context.Background(),
		&userpb.GetWishlistRequest{UserId: userID},
		mock.Anything,
	).Return(&userpb.GetWishlistResponse{
		Message: "Wishlist retrieved successfully",
		Items: []*userpb.WishlistItem{
			{
				Id:        "wishlist_item_1",
				UserId:    userID,
				ProductId: productID,
				AddedAt:   "2023-07-01T12:00:00Z",
			},
		},
		Count: 1,
	}, nil).Once()
}

// AddAddToWishlistMock cấu hình mock cho AddToWishlist gRPC call
func AddAddToWishlistMock(mockClient *MockUserServiceClient, userID string, productIDStr string) {
	// Chuyển đổi ID sản phẩm từ string sang int32
	productID, _ := strconv.Atoi(productIDStr)

	mockClient.On(
		"AddToWishlist",
		context.Background(),
		&userpb.AddToWishlistRequest{
			UserId:    userID,
			ProductId: int32(productID),
		},
		mock.Anything,
	).Return(&userpb.WishlistResponse{
		Success: true,
		Message: "Product added to wishlist",
	}, nil).Once()
}

// AddRemoveFromWishlistMock cấu hình mock cho RemoveFromWishlist gRPC call
func AddRemoveFromWishlistMock(mockClient *MockUserServiceClient, userID string, productIDStr string) {
	// Chuyển đổi ID sản phẩm từ string sang int32
	productID, _ := strconv.Atoi(productIDStr)

	mockClient.On(
		"RemoveFromWishlist",
		context.Background(),
		&userpb.RemoveFromWishlistRequest{
			UserId:    userID,
			ProductId: int32(productID),
		},
		mock.Anything,
	).Return(&userpb.WishlistResponse{
		Success: true,
		Message: "Product removed from wishlist",
	}, nil).Once()
}

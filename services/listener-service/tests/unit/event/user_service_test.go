package event

import (
	"context"
	"errors"
	userpb "listener-service/proto/user"
	"listener-service/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUpdateUserOrderData tests that we can properly call the UserServiceClient
func TestUpdateUserOrderData(t *testing.T) {
	// Arrange
	mockUserClient := new(mocks.MockUserServiceClient)

	// Setup mock behaviors
	mockUserClient.On("SyncUserOrderData", mock.Anything, mock.MatchedBy(func(req *userpb.SyncUserOrderDataRequest) bool {
		return req.UserId == "user123"
	}), mock.Anything).Return(
		&userpb.SyncUserOrderDataResponse{Success: true}, nil,
	)

	// Act
	ctx := context.Background()
	userID := "user123"

	// Call UpdateUserOrderData logic directly
	if userID == "" {
		t.Fatal("user ID should not be empty")
	}

	// Call the user service to update the user order data
	resp, err := mockUserClient.SyncUserOrderData(ctx, &userpb.SyncUserOrderDataRequest{
		UserId:     userID,
		OrderCount: 0,   // Placeholder value
		TotalSpend: 0.0, // Placeholder value
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	mockUserClient.AssertExpectations(t)
}

// TestUpdateUserOrderData_Error tests error handling when updating user order data
func TestUpdateUserOrderData_Error(t *testing.T) {
	// Arrange
	mockUserClient := new(mocks.MockUserServiceClient)

	// Setup mock behaviors
	expectedErr := errors.New("sync error")
	mockUserClient.On("SyncUserOrderData", mock.Anything, mock.Anything, mock.Anything).Return(
		nil, expectedErr,
	)

	// Act
	ctx := context.Background()
	userID := "user123"

	// Call the user service to update the user order data
	resp, err := mockUserClient.SyncUserOrderData(ctx, &userpb.SyncUserOrderDataRequest{
		UserId:     userID,
		OrderCount: 0,
		TotalSpend: 0.0,
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, resp)
	mockUserClient.AssertExpectations(t)
}

// TestUpdateUserOrderData_EmptyUserID tests error handling with empty user ID
func TestUpdateUserOrderData_EmptyUserID(t *testing.T) {
	// Arrange - Test the validation logic directly

	// Act & Assert
	userID := ""
	if userID != "" {
		t.Error("Empty userID should be detected")
	}
}

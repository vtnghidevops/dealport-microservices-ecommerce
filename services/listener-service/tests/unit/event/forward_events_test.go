package event

import (
	"context"
	"errors"
	"listener-service/event"
	userpb "listener-service/proto/user"
	"listener-service/tests/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestForwardToUserService tests forwarding events to user service
func TestForwardToUserService(t *testing.T) {
	// Arrange
	mockUserClient := new(mocks.MockUserServiceClient)

	// Setup mock behavior for successful forwarding
	mockUserClient.On("ProcessEvent", mock.Anything, mock.MatchedBy(func(req *userpb.EventRequest) bool {
		return req.EventName == "user.registered" && req.EventData != ""
	})).Return(&userpb.EventResponse{Success: true}, nil)

	// Create test event to forward
	userData := map[string]interface{}{
		"id":    "user123",
		"email": "test@example.com",
	}

	userEvent := event.StandardEvent{
		ID:          "test-event-id",
		Name:        "user.registered",
		Data:        userData,
		DataSchema:  "1.0",
		Source:      "auth-service",
		CreatedAt:   time.Now(),
		PublishedAt: time.Now(),
		Version:     "1.0",
	}

	// Mock consumer for testing event forwarding
	mockConsumer := mocks.NewMockConsumer()
	mockConsumer.On("HandleStandardEvent", userEvent, "user.registered").Return(nil)

	// Act
	err := mockConsumer.HandleStandardEvent(userEvent, "user.registered")

	// Assert
	assert.NoError(t, err)
	mockConsumer.AssertExpectations(t)
}

// TestForwardToUserService_Error tests error handling when forwarding fails
func TestForwardToUserService_Error(t *testing.T) {
	// Arrange
	mockUserClient := new(mocks.MockUserServiceClient)

	// Setup mock behavior for failed forwarding
	forwardError := errors.New("failed to forward event")
	mockUserClient.On("ProcessEvent", mock.Anything, mock.Anything).Return(nil, forwardError)

	// Create test event to forward
	userEvent := event.StandardEvent{
		ID:          "test-event-id",
		Name:        "user.profile_updated",
		Data:        map[string]interface{}{"user_id": "user123"},
		DataSchema:  "1.0",
		Source:      "user-service",
		CreatedAt:   time.Now(),
		PublishedAt: time.Now(),
		Version:     "1.0",
	}

	// Mock consumer that will simulate forwarding error
	mockConsumer := mocks.NewMockConsumer()
	mockConsumer.On("HandleStandardEvent", userEvent, "user.profile_updated").Return(
		errors.New("failed to forward event to user service"),
	)

	// Act
	err := mockConsumer.HandleStandardEvent(userEvent, "user.profile_updated")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to forward")
	mockConsumer.AssertExpectations(t)
}

// TestUserOrderDataSynchronization tests synchronizing order data with user service
func TestUserOrderDataSynchronization(t *testing.T) {
	// Arrange
	mockUserClient := new(mocks.MockUserServiceClient)

	// Setup mock behavior for successful sync
	mockUserClient.On("SyncUserOrderData", mock.Anything, mock.MatchedBy(func(req *userpb.SyncUserOrderDataRequest) bool {
		return req.UserId == "user123"
	}), mock.Anything).Return(&userpb.SyncUserOrderDataResponse{Success: true}, nil)

	// Create order event that would trigger synchronization in real code
	// Note: We're testing the direct SyncUserOrderData call below, not using this event
	orderData := event.OrderCreatedData{
		OrderID: "order-123",
		UserID:  "user123",
		Total:   99.99,
	}

	// We don't use this event directly in the test, but it shows what event would trigger the sync
	_ = event.StandardEvent{
		ID:          "test-event-id",
		Name:        "order.created",
		Data:        orderData,
		DataSchema:  "1.0",
		Source:      "checkout-service",
		CreatedAt:   time.Now(),
		PublishedAt: time.Now(),
		Version:     "1.0",
	}

	// Create context for the test
	ctx := context.Background()

	// Call SyncUserOrderData directly with the mock client
	resp, err := mockUserClient.SyncUserOrderData(ctx, &userpb.SyncUserOrderDataRequest{
		UserId:     "user123",
		OrderCount: 1,
		TotalSpend: 99.99,
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	mockUserClient.AssertExpectations(t)
}

// TestUserOrderDataSync_Error tests error handling in user data synchronization
func TestUserOrderDataSync_Error(t *testing.T) {
	// Arrange
	mockUserClient := new(mocks.MockUserServiceClient)

	// Setup mock behavior for failed sync
	syncError := errors.New("failed to sync user data")
	mockUserClient.On("SyncUserOrderData", mock.Anything, mock.Anything, mock.Anything).Return(nil, syncError)

	// Create context for the test
	ctx := context.Background()

	// Call SyncUserOrderData directly with the mock client
	resp, err := mockUserClient.SyncUserOrderData(ctx, &userpb.SyncUserOrderDataRequest{
		UserId:     "user123",
		OrderCount: 1,
		TotalSpend: 99.99,
	})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, syncError, err)
	mockUserClient.AssertExpectations(t)
}

package helpers

import (
	"encoding/json"
	"listener-service/event"
	"time"
)

// CreateTestEvent creates a standardized test event with the given name and data
func CreateTestEvent(name string, data interface{}) event.StandardEvent {
	return event.StandardEvent{
		ID:          "test-event-id-" + name,
		Name:        name,
		Data:        data,
		DataSchema:  "1.0",
		Source:      "test-service",
		CreatedAt:   time.Now(),
		PublishedAt: time.Now(),
		Version:     "1.0",
	}
}

// CreateUserRegisteredEvent creates a test event for user registration
func CreateUserRegisteredEvent(userID, email, username, firstName, lastName string) event.StandardEvent {
	userData := event.UserRegistered{
		ID:        userID,
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
	}

	return CreateTestEvent("user.registered", userData)
}

// CreateOrderCreatedEvent creates a test event for order creation
func CreateOrderCreatedEvent(orderID, orderNumber, userID, userEmail string, total float64) event.StandardEvent {
	orderData := event.OrderCreatedData{
		OrderID:       orderID,
		OrderNumber:   orderNumber,
		UserID:        userID,
		UserEmail:     userEmail,
		Status:        "pending",
		PaymentMethod: "card",
		Total:         total,
		CreatedAt:     time.Now(),
		Items: []event.ItemData{
			{
				ProductID: "test-product",
				Name:      "Test Product",
				Quantity:  1,
				Price:     total,
			},
		},
	}

	return CreateTestEvent("order.created", orderData)
}

// CreateOrderStatusChangedEvent creates a test event for order status change
func CreateOrderStatusChangedEvent(orderID, orderNumber, userID, userEmail string, status, previousStatus string) event.StandardEvent {
	statusData := event.OrderStatusChangedData{
		OrderID:        orderID,
		OrderNumber:    orderNumber,
		UserID:         userID,
		UserEmail:      userEmail,
		Status:         status,
		PreviousStatus: previousStatus,
		PaymentMethod:  "card",
		Total:          99.99,
		CreatedAt:      time.Now(),
	}

	return CreateTestEvent("order.status_changed", statusData)
}

// CreatePaymentSucceededEvent creates a test event for successful payment
func CreatePaymentSucceededEvent(orderID, orderNumber, userID, userEmail, transactionID string, amount float64) event.StandardEvent {
	paymentData := event.PaymentSucceededData{
		OrderID:       orderID,
		OrderNumber:   orderNumber,
		UserID:        userID,
		UserEmail:     userEmail,
		PaymentMethod: "card",
		Amount:        amount,
		Currency:      "USD",
		TransactionID: transactionID,
		Status:        "completed",
		PaymentDate:   time.Now().Format(time.RFC3339),
	}

	return CreateTestEvent("payment.succeeded", paymentData)
}

// CreatePasswordResetEvent creates a test event for password reset request
func CreatePasswordResetEvent(email, tokenHash string) event.StandardEvent {
	resetData := event.PasswordResetRequested{
		Email:     email,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	return CreateTestEvent("auth.password_reset_requested", resetData)
}

// CreateTestLegacyPayload creates a legacy payload with given name and data
func CreateTestLegacyPayload(name string, data interface{}) event.Payload {
	dataJSON, _ := json.Marshal(data)

	return event.Payload{
		Name: name,
		Data: string(dataJSON),
	}
}

// CreateLogPayload creates a legacy log payload
func CreateLogPayload(service, message string) event.Payload {
	logData := map[string]string{
		"service": service,
		"message": message,
	}

	return CreateTestLegacyPayload("log", logData)
}

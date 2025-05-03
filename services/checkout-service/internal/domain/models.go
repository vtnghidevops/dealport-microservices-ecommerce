package domain

import (
	"context"
	"errors"
	"time"
)

// Common errors
var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrInvalidOrderData  = errors.New("invalid order data")
	ErrPaymentFailed     = errors.New("payment processing failed")
	ErrValidationFailed  = errors.New("validation failed")
	ErrDatabaseOperation = errors.New("database operation failed")
)

// OrderItem represents an item in an order
type OrderItem struct {
	ID        string  `json:"id" bson:"id"`
	ProductID string  `json:"productId" bson:"product_id"`
	Name      string  `json:"name" bson:"name"`
	Price     float64 `json:"price" bson:"price"`
	Quantity  int32   `json:"quantity" bson:"quantity"`
	Subtotal  float64 `json:"subtotal" bson:"subtotal"`
	ImageURL  string  `json:"imageUrl" bson:"image_url"`
}

// BillingInfo contains billing information for an order
type BillingInfo struct {
	FirstName   string `json:"firstName" bson:"first_name"`
	LastName    string `json:"lastName" bson:"last_name"`
	CompanyName string `json:"companyName" bson:"company_name"`
	Address     string `json:"address" bson:"address"`
	Country     string `json:"country" bson:"country"`
	Region      string `json:"region" bson:"region"`
	City        string `json:"city" bson:"city"`
	ZipCode     string `json:"zipCode" bson:"zip_code"`
	Email       string `json:"email" bson:"email"`
	Phone       string `json:"phone" bson:"phone"`
}

// ShippingInfo contains shipping information for an order
type ShippingInfo struct {
	ShipToDifferentAddress bool    `json:"shipToDifferentAddress" bson:"ship_to_different_address"`
	FirstName              string  `json:"firstName" bson:"first_name"`
	LastName               string  `json:"lastName" bson:"last_name"`
	CompanyName            string  `json:"companyName" bson:"company_name"`
	Address                string  `json:"address" bson:"address"`
	Country                string  `json:"country" bson:"country"`
	Region                 string  `json:"region" bson:"region"`
	City                   string  `json:"city" bson:"city"`
	ZipCode                string  `json:"zipCode" bson:"zip_code"`
	ShippingMethod         string  `json:"shippingMethod" bson:"shipping_method"`
	ShippingCost           float64 `json:"shippingCost" bson:"shipping_cost"`
}

// PaymentInfo contains payment information for an order
type PaymentInfo struct {
	PaymentMethod string  `json:"paymentMethod" bson:"payment_method"`
	TransactionID string  `json:"transactionId" bson:"transaction_id"`
	Status        string  `json:"status" bson:"status"`
	Amount        float64 `json:"amount" bson:"amount"`
	Currency      string  `json:"currency" bson:"currency"`
	PaymentDate   string  `json:"paymentDate" bson:"payment_date"`
}

// OrderTotals represents the financial totals for an order
type OrderTotals struct {
	Subtotal float64 `json:"subtotal" bson:"subtotal"`
	Shipping float64 `json:"shipping" bson:"shipping"`
	Discount float64 `json:"discount" bson:"discount"`
	Tax      float64 `json:"tax" bson:"tax"`
	Total    float64 `json:"total" bson:"total"`
}

// Order represents a complete order in the system
type Order struct {
	ID           string       `json:"id" bson:"_id"`
	UserID       string       `json:"userId" bson:"user_id"`
	OrderNumber  string       `json:"orderNumber" bson:"order_number"`
	Status       string       `json:"status" bson:"status"`
	Items        []OrderItem  `json:"items" bson:"items"`
	BillingInfo  BillingInfo  `json:"billingInfo" bson:"billing_info"`
	ShippingInfo ShippingInfo `json:"shippingInfo" bson:"shipping_info"`
	PaymentInfo  PaymentInfo  `json:"paymentInfo" bson:"payment_info"`
	Totals       OrderTotals  `json:"totals" bson:"totals"`
	CouponCode   string       `json:"couponCode" bson:"coupon_code"`
	Notes        string       `json:"notes" bson:"notes"`
	CreatedAt    time.Time    `json:"createdAt" bson:"created_at"`
	UpdatedAt    time.Time    `json:"updatedAt" bson:"updated_at"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// CheckoutValidationRequest is the data sent to validate a checkout
type CheckoutValidationRequest struct {
	Items         []OrderItem  `json:"items"`
	BillingInfo   BillingInfo  `json:"billingInfo"`
	ShippingInfo  ShippingInfo `json:"shippingInfo"`
	PaymentMethod string       `json:"paymentMethod"`
}

// CheckoutValidationResult is the result of a checkout validation
type CheckoutValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors"`
}

// PaymentRequest represents a payment processing request
type PaymentRequest struct {
	OrderID       string  `json:"orderId"`
	PaymentMethod string  `json:"paymentMethod"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	ReturnURL     string  `json:"returnUrl"`
}

// PaymentResult represents the result of a payment processing operation
type PaymentResult struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
	RedirectURL   string `json:"redirectUrl"`
	Message       string `json:"message"`
}

// OrderService defines the interface for order-related operations
type OrderService interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrderByID(ctx context.Context, orderID string, userID string) (*Order, error)
	ListOrdersByUserID(ctx context.Context, userID string, page, pageSize int) ([]*Order, int, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error
	ValidateCheckout(ctx context.Context, data CheckoutValidationRequest) (CheckoutValidationResult, error)
	ProcessPayment(ctx context.Context, paymentReq PaymentRequest) (PaymentResult, error)
	GetRepository() OrderRepository
}

// OrderRepository defines the interface for order data persistence
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*Order, error)
	ListOrdersByUserID(ctx context.Context, userID string, skip, limit int) ([]*Order, int, error)
	ListAllOrders(ctx context.Context, skip, limit int) ([]*Order, int, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error
	UpdatePaymentInfo(ctx context.Context, orderID string, paymentInfo PaymentInfo) error
}

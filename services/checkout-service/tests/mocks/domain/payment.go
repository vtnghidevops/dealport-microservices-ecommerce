package domain

import "time"

// Payment represents a payment record in the system
type Payment struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"order_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Method        string    `json:"method"`
	Status        string    `json:"status"`
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PaymentRequest represents a request to process a payment
type PaymentRequest struct {
	OrderID       string  `json:"order_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"payment_method"`
	CustomerID    string  `json:"customer_id"`
	Description   string  `json:"description"`
	ReturnURL     string  `json:"return_url,omitempty"`
}

// PaymentResult represents the result of a payment processing operation
type PaymentResult struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transaction_id,omitempty"`
	Status        string `json:"status"`
	RedirectURL   string `json:"redirect_url,omitempty"`
	Message       string `json:"message,omitempty"`
}

// RefundResult represents the result of a refund operation
type RefundResult struct {
	Success       bool    `json:"success"`
	RefundID      string  `json:"refund_id,omitempty"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	TransactionID string  `json:"transaction_id"`
	Message       string  `json:"message,omitempty"`
}

// PaymentVerificationResult represents the result of payment verification
type PaymentVerificationResult struct {
	Verified      bool    `json:"verified"`
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	Message       string  `json:"message,omitempty"`
}

package mocks

import (
	testDomain "checkout-service/tests/mocks/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockPaymentService mocks the payment service for integration tests
type MockPaymentService struct {
	mock.Mock
}

// ProcessPayment mocks the process payment method
func (m *MockPaymentService) ProcessPayment(ctx context.Context, paymentReq testDomain.PaymentRequest) (testDomain.PaymentResult, error) {
	args := m.Called(ctx, paymentReq)
	return args.Get(0).(testDomain.PaymentResult), args.Error(1)
}

// ValidatePaymentMethod mocks the validate payment method
func (m *MockPaymentService) ValidatePaymentMethod(paymentMethod string) error {
	args := m.Called(paymentMethod)
	return args.Error(0)
}

// GetPaymentByOrderID mocks the retrieval of payment by order ID
func (m *MockPaymentService) GetPaymentByOrderID(ctx context.Context, orderID string) (*testDomain.Payment, error) {
	args := m.Called(ctx, orderID)
	var payment *testDomain.Payment
	if args.Get(0) != nil {
		payment = args.Get(0).(*testDomain.Payment)
	}
	return payment, args.Error(1)
}

// RefundPayment mocks the refund payment method
func (m *MockPaymentService) RefundPayment(ctx context.Context, orderID string, amount float64, reason string) (*testDomain.RefundResult, error) {
	args := m.Called(ctx, orderID, amount, reason)
	var result *testDomain.RefundResult
	if args.Get(0) != nil {
		result = args.Get(0).(*testDomain.RefundResult)
	}
	return result, args.Error(1)
}

// VerifyPayment mocks the verify payment method
func (m *MockPaymentService) VerifyPayment(ctx context.Context, transactionID string) (*testDomain.PaymentVerificationResult, error) {
	args := m.Called(ctx, transactionID)
	var result *testDomain.PaymentVerificationResult
	if args.Get(0) != nil {
		result = args.Get(0).(*testDomain.PaymentVerificationResult)
	}
	return result, args.Error(1)
}

// CreateSuccessfulPaymentResult creates a successful payment result for testing
func (m *MockPaymentService) CreateSuccessfulPaymentResult(orderID string) testDomain.PaymentResult {
	return testDomain.PaymentResult{
		Success:       true,
		TransactionID: "txn_" + orderID,
		Status:        "completed",
		RedirectURL:   "",
		Message:       "Payment processed successfully",
	}
}

// CreateFailedPaymentResult creates a failed payment result for testing
func (m *MockPaymentService) CreateFailedPaymentResult(errorMessage string) testDomain.PaymentResult {
	return testDomain.PaymentResult{
		Success:       false,
		TransactionID: "",
		Status:        "failed",
		RedirectURL:   "",
		Message:       errorMessage,
	}
}

package grpc

import (
	"context"
	"errors"
	"log"

	"checkout-service/internal/payment"

	pb "checkout-service/proto/payment"
)

// PaymentHandler handles payment-related gRPC requests
type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	momoService  *payment.MomoService
	vnpayService *payment.VNPayService
	// Comment out QuickPay service as we focus on ATM wallet implementation
	// momoQuickPayService *payment.MomoQuickPayService
}

// NewPaymentHandler creates a new payment handler with initialized payment services
func NewPaymentHandler() *PaymentHandler {
	// Load payment configuration
	paymentConfig := payment.NewPaymentConfig()

	// Initialize payment services
	momoService := payment.NewMomoService(paymentConfig.MoMo)
	vnpayService := payment.NewVNPayService(paymentConfig.VNPay)
	// Comment out QuickPay service as we focus on ATM wallet implementation
	// momoQuickPayService := payment.NewMomoQuickPayService(paymentConfig.MoMo)

	return &PaymentHandler{
		momoService:  momoService,
		vnpayService: vnpayService,
		// momoQuickPayService: momoQuickPayService,
	}
}

// CreateMomoPayment creates a new MoMo payment
func (h *PaymentHandler) CreateMomoPayment(ctx context.Context, req *pb.MomoPaymentRequest) (*pb.MomoPaymentResponse, error) {
	if req.PaymentInfo == nil {
		return nil, errors.New("payment info is required")
	}

	log.Printf("Received MoMo payment request for order %s with amount %d",
		req.PaymentInfo.OrderId, req.PaymentInfo.Amount)

	// Call the MoMo service to create a payment
	result, err := h.momoService.CreatePayment(
		req.PaymentInfo.OrderId,
		req.PaymentInfo.Amount,
		req.OrderInfo,
		req.PaymentInfo.ReturnUrl,
	)

	if err != nil {
		log.Printf("Error creating MoMo payment: %v", err)
		return &pb.MomoPaymentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Return the response
	return &pb.MomoPaymentResponse{
		Success:       result.Success,
		Message:       result.Message,
		PaymentUrl:    result.PaymentURL,
		RequestId:     result.RequestID,
		TransactionId: result.TransactionID,
		OrderId:       result.OrderID,
	}, nil
}

// VerifyMomoPayment verifies a MoMo payment callback
func (h *PaymentHandler) VerifyMomoPayment(ctx context.Context, req *pb.MomoVerifyRequest) (*pb.PaymentVerifyResponse, error) {
	if req.Params == nil || len(req.Params) == 0 {
		return nil, errors.New("callback parameters are required")
	}

	log.Printf("Verifying MoMo payment callback with %d parameters", len(req.Params))

	// Call the MoMo service to verify the payment
	result, err := h.momoService.VerifyPayment(req.Params)

	if err != nil {
		log.Printf("Error verifying MoMo payment: %v", err)
		return &pb.PaymentVerifyResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Return the verification result
	return &pb.PaymentVerifyResponse{
		Success:       result.Success,
		Message:       result.Message,
		OrderId:       result.OrderID,
		TransactionId: result.TransactionID,
		Amount:        result.Amount,
	}, nil
}

// CreateVnpayPayment creates a new VNPAY payment
func (h *PaymentHandler) CreateVnpayPayment(ctx context.Context, req *pb.VnpayPaymentRequest) (*pb.VnpayPaymentResponse, error) {
	if req.PaymentInfo == nil {
		return nil, errors.New("payment info is required")
	}

	log.Printf("Received VNPAY payment request for order %s with amount %d",
		req.PaymentInfo.OrderId, req.PaymentInfo.Amount)

	// Call the VNPAY service to create a payment
	result, err := h.vnpayService.CreatePayment(
		req.PaymentInfo.OrderId,
		req.PaymentInfo.Amount,
		req.IpAddr,
		req.PaymentInfo.ReturnUrl,
	)

	if err != nil {
		log.Printf("Error creating VNPAY payment: %v", err)
		return &pb.VnpayPaymentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Return the response
	return &pb.VnpayPaymentResponse{
		Success:    result.Success,
		Message:    result.Message,
		PaymentUrl: result.PaymentURL,
		VnpTxnRef:  result.VnpTxnRef,
		OrderId:    result.OrderID,
	}, nil
}

// VerifyVnpayPayment verifies a VNPAY payment callback
func (h *PaymentHandler) VerifyVnpayPayment(ctx context.Context, req *pb.VnpayVerifyRequest) (*pb.PaymentVerifyResponse, error) {
	if req.Params == nil || len(req.Params) == 0 {
		return nil, errors.New("callback parameters are required")
	}

	log.Printf("Verifying VNPAY payment callback with %d parameters", len(req.Params))

	// Call the VNPAY service to verify the payment
	result, err := h.vnpayService.VerifyPayment(req.Params)

	if err != nil {
		log.Printf("Error verifying VNPAY payment: %v", err)
		return &pb.PaymentVerifyResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Return the verification result
	return &pb.PaymentVerifyResponse{
		Success:       result.Success,
		Message:       result.Message,
		OrderId:       result.OrderID,
		TransactionId: result.TransactionID,
		Amount:        result.Amount,
	}, nil
}

// Comment out QuickPay methods as we focus on ATM wallet implementation
/*
// CreateMomoQRPayment creates a MoMo QR Code payment
func (h *PaymentHandler) CreateMomoQRPayment(ctx context.Context, req *pb.MomoQRPaymentRequest) (*pb.MomoQuickPayResponse, error) {
	log.Printf("Received MoMo QR payment request for order %s with amount %d",
		req.OrderId, req.Amount)

	// Call the QuickPay service to create a QR payment
	result, err := h.momoQuickPayService.CreateQRPayment(
		req.OrderId,
		req.Amount,
		req.OrderInfo,
		req.ReturnUrl,
	)

	if err != nil {
		log.Printf("Error creating MoMo QR payment: %v", err)
		return &pb.MomoQuickPayResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Return the response
	return &pb.MomoQuickPayResponse{
		Success:       result.Success,
		Message:       result.Message,
		OrderId:       result.OrderID,
		RequestId:     result.RequestID,
		TransactionId: result.TransactionID,
		PaymentUrl:    result.PaymentURL,
		QrCodeUrl:     result.QrCodeURL,
		DeepLink:      result.DeepLink,
		Amount:        result.Amount,
	}, nil
}

// CreateMomoPosPayment creates a MoMo POS payment
func (h *PaymentHandler) CreateMomoPosPayment(ctx context.Context, req *pb.MomoPosPaymentRequest) (*pb.MomoQuickPayResponse, error) {
	log.Printf("Received MoMo POS payment request for order %s with amount %d",
		req.OrderId, req.Amount)

	// Call the QuickPay service to create a POS payment
	result, err := h.momoQuickPayService.CreatePosPayment(
		req.OrderId,
		req.Amount,
		req.OrderInfo,
		req.PaymentCode,
		req.ReturnUrl,
	)

	if err != nil {
		log.Printf("Error creating MoMo POS payment: %v", err)
		return &pb.MomoQuickPayResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Return the response
	return &pb.MomoQuickPayResponse{
		Success:       result.Success,
		Message:       result.Message,
		OrderId:       result.OrderID,
		RequestId:     result.RequestID,
		TransactionId: result.TransactionID,
		PaymentUrl:    result.PaymentURL,
		QrCodeUrl:     result.QrCodeURL,
		DeepLink:      result.DeepLink,
		Amount:        result.Amount,
	}, nil
}
*/

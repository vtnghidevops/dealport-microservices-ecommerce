package grpc

import (
	"context"
	"errors"
	"log"
	"time"

	"checkout-service/internal/domain"
	"checkout-service/internal/payment"

	pb "checkout-service/proto/payment"
)

// PaymentHandler handles payment-related gRPC requests
type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	momoService  *payment.MomoService
	vnpayService *payment.VNPayService
	orderService domain.OrderService
	// Comment out QuickPay service as we focus on ATM wallet implementation
	// momoQuickPayService *payment.MomoQuickPayService
}

// NewPaymentHandler creates a new payment handler with initialized payment services
func NewPaymentHandler(orderService domain.OrderService) *PaymentHandler {
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
		orderService: orderService,
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

// ProcessMomoCallback processes a MoMo payment callback and updates the order status
func (h *PaymentHandler) ProcessMomoCallback(ctx context.Context, req *pb.MomoCallbackRequest) (*pb.PaymentVerifyResponse, error) {
	log.Printf("CHECKOUT-MOMO-CALLBACK: Starting to process MoMo callback at %s", time.Now().Format(time.RFC3339))

	if req.Params == nil || len(req.Params) == 0 {
		log.Printf("CHECKOUT-MOMO-CALLBACK-ERROR: No callback parameters received")
		return nil, errors.New("callback parameters are required")
	}

	log.Printf("CHECKOUT-MOMO-CALLBACK: Processing MoMo payment callback with %d parameters", len(req.Params))

	// Log important parameters
	log.Printf("CHECKOUT-MOMO-CALLBACK: OrderId = %s", req.Params["orderId"])
	log.Printf("CHECKOUT-MOMO-CALLBACK: TransId = %s", req.Params["transId"])
	log.Printf("CHECKOUT-MOMO-CALLBACK: ResultCode = %s", req.Params["resultCode"])
	log.Printf("CHECKOUT-MOMO-CALLBACK: Message = %s", req.Params["message"])
	log.Printf("CHECKOUT-MOMO-CALLBACK: ExtraData = %s", req.Params["extraData"])

	// First verify the payment with MoMo
	log.Printf("CHECKOUT-MOMO-CALLBACK: Calling MoMo service to verify payment")
	verifyResult, err := h.momoService.VerifyPayment(req.Params)
	if err != nil {
		log.Printf("CHECKOUT-MOMO-CALLBACK-ERROR: Failed to verify payment: %v", err)
		return &pb.PaymentVerifyResponse{
			Success: false,
			Message: "Failed to verify payment: " + err.Error(),
		}, nil
	}

	log.Printf("CHECKOUT-MOMO-CALLBACK: MoMo verification result: success=%t, message=%s, orderId=%s",
		verifyResult.Success, verifyResult.Message, verifyResult.OrderID)

	if !verifyResult.Success {
		log.Printf("CHECKOUT-MOMO-CALLBACK-ERROR: Payment verification failed: %s", verifyResult.Message)
		return &pb.PaymentVerifyResponse{
			Success: false,
			Message: verifyResult.Message,
		}, nil
	}

	// Get the order details
	orderID := verifyResult.OrderID
	log.Printf("CHECKOUT-MOMO-CALLBACK: Getting order details for OrderID=%s", orderID)

	order, err := h.orderService.GetOrderByID(ctx, orderID, "")
	if err != nil {
		log.Printf("CHECKOUT-MOMO-CALLBACK-ERROR: Failed to get order details: %v", err)
		return &pb.PaymentVerifyResponse{
			Success: false,
			Message: "Failed to get order details: " + err.Error(),
			OrderId: orderID,
		}, nil
	}

	log.Printf("CHECKOUT-MOMO-CALLBACK: Found order %s with current status: %s, payment status: %s",
		orderID, order.Status, order.PaymentInfo.Status)

	// Update the payment info
	paymentInfo := domain.PaymentInfo{
		PaymentMethod: "momo",
		TransactionID: verifyResult.TransactionID,
		Status:        "completed", // Update status to completed
		Amount:        float64(verifyResult.Amount),
		Currency:      "VND",
		PaymentDate:   time.Now().Format(time.RFC3339),
	}

	log.Printf("CHECKOUT-MOMO-CALLBACK: Updating payment info to: method=%s, transactionId=%s, status=%s, amount=%f",
		paymentInfo.PaymentMethod, paymentInfo.TransactionID, paymentInfo.Status, paymentInfo.Amount)

	// Update payment info in the database directly using the order service
	err = h.orderService.GetRepository().UpdatePaymentInfo(ctx, orderID, paymentInfo)
	if err != nil {
		log.Printf("CHECKOUT-MOMO-CALLBACK-ERROR: Failed to update payment info: %v", err)
		return &pb.PaymentVerifyResponse{
			Success: false,
			Message: "Failed to update payment info: " + err.Error(),
			OrderId: orderID,
		}, nil
	}
	log.Printf("CHECKOUT-MOMO-CALLBACK: Successfully updated payment info in database")

	// Update order status to processing (payment completed)
	log.Printf("CHECKOUT-MOMO-CALLBACK: Updating order status from '%s' to 'processing'", order.Status)
	err = h.orderService.UpdateOrderStatus(ctx, orderID, "processing")
	if err != nil {
		log.Printf("CHECKOUT-MOMO-CALLBACK-ERROR: Failed to update order status: %v", err)
		return &pb.PaymentVerifyResponse{
			Success: false,
			Message: "Failed to update order status: " + err.Error(),
			OrderId: orderID,
		}, nil
	}
	log.Printf("CHECKOUT-MOMO-CALLBACK: Successfully updated order status to 'processing'")

	// Double-check the order was updated correctly
	updatedOrder, err := h.orderService.GetOrderByID(ctx, orderID, "")
	if err != nil {
		log.Printf("CHECKOUT-MOMO-CALLBACK-WARNING: Could not verify order updates: %v", err)
	} else {
		log.Printf("CHECKOUT-MOMO-CALLBACK: Order after update - Status: %s, Payment Status: %s",
			updatedOrder.Status, updatedOrder.PaymentInfo.Status)
	}

	log.Printf("CHECKOUT-MOMO-CALLBACK-SUCCESS: Payment completed successfully for order %s", orderID)

	return &pb.PaymentVerifyResponse{
		Success:       true,
		Message:       "Payment completed successfully",
		OrderId:       orderID,
		TransactionId: verifyResult.TransactionID,
		Amount:        verifyResult.Amount,
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

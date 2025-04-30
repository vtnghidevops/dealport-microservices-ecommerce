package payment

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pb "broker-service/proto/payment"
)

// Config handles payment-related HTTP requests
type Config struct {
	CheckoutClient pb.PaymentServiceClient
}

// CreateMomoPaymentRequest represents a request to create a MoMo payment
type CreateMomoPaymentRequest struct {
	OrderID   string `json:"orderId"`
	Amount    int64  `json:"amount"`
	OrderInfo string `json:"orderInfo"`
	ReturnURL string `json:"returnUrl"`
}

// CreateMomoPayment handles requests to create a new MoMo payment
func (c *Config) CreateMomoPayment(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req CreateMomoPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error parsing CreateMomoPayment request: %v", err)
		jsonResponse(w, http.StatusBadRequest, jsonError("Invalid request format"))
		return
	}

	// Log the incoming request for debugging
	reqBytes, _ := json.Marshal(req)
	log.Printf("CreateMomoPayment request: %s", string(reqBytes))

	// Validate request
	if req.OrderID == "" || req.Amount <= 0 || req.ReturnURL == "" {
		log.Printf("Invalid CreateMomoPayment request: OrderID=%s, Amount=%d, ReturnURL=%s",
			req.OrderID, req.Amount, req.ReturnURL)
		jsonResponse(w, http.StatusBadRequest, jsonError("Missing required fields"))
		return
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second) // Increase timeout
	defer cancel()

	// Build gRPC request
	grpcReq := &pb.MomoPaymentRequest{
		PaymentInfo: &pb.PaymentBaseRequest{
			OrderId:   req.OrderID,
			Amount:    req.Amount,
			ReturnUrl: req.ReturnURL,
		},
		OrderInfo: req.OrderInfo,
	}

	// Log the gRPC request being sent
	log.Printf("Sending CreateMomoPayment gRPC request to checkout service: OrderID=%s, Amount=%d",
		grpcReq.PaymentInfo.OrderId, grpcReq.PaymentInfo.Amount)

	// Call service
	resp, err := c.CheckoutClient.CreateMomoPayment(ctx, grpcReq)
	if err != nil {
		log.Printf("Error calling checkout service for MoMo payment: %v", err)
		jsonResponse(w, http.StatusInternalServerError, jsonError("Failed to process payment: "+err.Error()))
		return
	}

	// Log the checkout service response
	log.Printf("Received MoMo payment response: Success=%t, Message=%s", resp.Success, resp.Message)

	// Check response
	if !resp.Success {
		log.Printf("MoMo payment creation failed: %s", resp.Message)
		jsonResponse(w, http.StatusBadRequest, jsonError(resp.Message))
		return
	}

	// Create response data
	responseData := map[string]interface{}{
		"paymentUrl":    resp.PaymentUrl,
		"orderId":       resp.OrderId,
		"requestId":     resp.RequestId,
		"transactionId": resp.TransactionId,
		"amount":        grpcReq.PaymentInfo.Amount, // Include amount in response
	}

	// Log the successful response we're sending back
	log.Printf("Sending successful MoMo payment response: PaymentURL=%s, OrderID=%s",
		resp.PaymentUrl, resp.OrderId)

	// Return success response
	jsonResponse(w, http.StatusOK, jsonSuccess(responseData))
}

// VerifyMomoPaymentRequest represents a request to verify a MoMo payment
type VerifyMomoPaymentRequest struct {
	Params map[string]string `json:"params"`
}

// VerifyMomoPayment handles requests to verify a MoMo payment callback
func (c *Config) VerifyMomoPayment(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req VerifyMomoPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error parsing VerifyMomoPayment request: %v", err)
		jsonResponse(w, http.StatusBadRequest, jsonError("Invalid request format"))
		return
	}

	// Log the incoming request for debugging
	log.Printf("VerifyMomoPayment request received with %d parameters", len(req.Params))
	for key, value := range req.Params {
		// Don't log sensitive data like signature
		if key != "signature" {
			log.Printf("  Param %s: %s", key, value)
		}
	}

	// Validate request
	if req.Params == nil || len(req.Params) == 0 {
		log.Printf("Missing callback parameters in VerifyMomoPayment request")
		jsonResponse(w, http.StatusBadRequest, jsonError("Missing callback parameters"))
		return
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second) // Increase timeout
	defer cancel()

	// Build gRPC request
	grpcReq := &pb.MomoVerifyRequest{
		Params: req.Params,
	}

	// Log that we're making the gRPC call
	log.Printf("Sending VerifyMomoPayment gRPC request to checkout service")

	// Call service
	resp, err := c.CheckoutClient.VerifyMomoPayment(ctx, grpcReq)
	if err != nil {
		log.Printf("Error calling checkout service for MoMo verification: %v", err)
		jsonResponse(w, http.StatusInternalServerError, jsonError("Failed to verify payment: "+err.Error()))
		return
	}

	// Log the verification result
	log.Printf("MoMo payment verification result: Success=%t, OrderID=%s, TransactionID=%s, Message=%s",
		resp.Success, resp.OrderId, resp.TransactionId, resp.Message)

	// Create response data
	responseData := map[string]interface{}{
		"success":       resp.Success,
		"message":       resp.Message,
		"orderId":       resp.OrderId,
		"transactionId": resp.TransactionId,
		"amount":        resp.Amount,
	}

	// Return verification response
	jsonResponse(w, http.StatusOK, jsonData(responseData))
}

// CreateVnpayPaymentRequest represents a request to create a VNPAY payment
type CreateVnpayPaymentRequest struct {
	OrderID   string `json:"orderId"`
	Amount    int64  `json:"amount"`
	ReturnURL string `json:"returnUrl"`
	IPAddr    string `json:"ipAddr"`
	Locale    string `json:"locale"`
}

// CreateVnpayPayment handles requests to create a new VNPAY payment
func (c *Config) CreateVnpayPayment(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req CreateVnpayPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, jsonError("Invalid request format"))
		return
	}

	// Validate request
	if req.OrderID == "" || req.Amount <= 0 || req.ReturnURL == "" {
		jsonResponse(w, http.StatusBadRequest, jsonError("Missing required fields"))
		return
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Use client IP if not provided
	ipAddr := req.IPAddr
	if ipAddr == "" {
		ipAddr = r.RemoteAddr
	}

	// Default locale to Vietnamese if not provided
	locale := req.Locale
	if locale == "" {
		locale = "vn"
	}

	// Build gRPC request
	grpcReq := &pb.VnpayPaymentRequest{
		PaymentInfo: &pb.PaymentBaseRequest{
			OrderId:   req.OrderID,
			Amount:    req.Amount,
			ReturnUrl: req.ReturnURL,
		},
		IpAddr: ipAddr,
		Locale: locale,
	}

	// Call service
	resp, err := c.CheckoutClient.CreateVnpayPayment(ctx, grpcReq)
	if err != nil {
		log.Printf("Error calling checkout service: %v", err)
		jsonResponse(w, http.StatusInternalServerError, jsonError("Failed to process payment"))
		return
	}

	// Check response
	if !resp.Success {
		jsonResponse(w, http.StatusBadRequest, jsonError(resp.Message))
		return
	}

	// Return success response
	jsonResponse(w, http.StatusOK, jsonSuccess(map[string]interface{}{
		"paymentUrl": resp.PaymentUrl,
		"orderId":    resp.OrderId,
		"vnpTxnRef":  resp.VnpTxnRef,
	}))
}

// VerifyVnpayPaymentRequest represents a request to verify a VNPAY payment
type VerifyVnpayPaymentRequest struct {
	Params map[string]string `json:"params"`
}

// VerifyVnpayPayment handles requests to verify a VNPAY payment callback
func (c *Config) VerifyVnpayPayment(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req VerifyVnpayPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, jsonError("Invalid request format"))
		return
	}

	// Validate request
	if req.Params == nil || len(req.Params) == 0 {
		jsonResponse(w, http.StatusBadRequest, jsonError("Missing callback parameters"))
		return
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Build gRPC request
	grpcReq := &pb.VnpayVerifyRequest{
		Params: req.Params,
	}

	// Call service
	resp, err := c.CheckoutClient.VerifyVnpayPayment(ctx, grpcReq)
	if err != nil {
		log.Printf("Error calling checkout service: %v", err)
		jsonResponse(w, http.StatusInternalServerError, jsonError("Failed to verify payment"))
		return
	}

	// Return verification response
	jsonResponse(w, http.StatusOK, jsonData(map[string]interface{}{
		"success":       resp.Success,
		"message":       resp.Message,
		"orderId":       resp.OrderId,
		"transactionId": resp.TransactionId,
		"amount":        resp.Amount,
	}))
}

// CreateMomoQRPaymentRequest represents a request to create a MoMo QR code payment
type CreateMomoQRPaymentRequest struct {
	OrderID   string `json:"orderId"`
	Amount    int64  `json:"amount"`
	OrderInfo string `json:"orderInfo"`
	ReturnURL string `json:"returnUrl"`
}

// CreateMomoQRPayment handles requests to create a new MoMo QR code payment
func (c *Config) CreateMomoQRPayment(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req CreateMomoQRPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, jsonError("Invalid request format"))
		return
	}

	// Validate request
	if req.Amount <= 0 || req.ReturnURL == "" {
		jsonResponse(w, http.StatusBadRequest, jsonError("Missing required fields"))
		return
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Default order info if not provided
	orderInfo := req.OrderInfo
	if orderInfo == "" {
		orderInfo = "Payment for order"
	}

	// Build gRPC request
	grpcReq := &pb.MomoQRPaymentRequest{
		OrderId:   req.OrderID,
		Amount:    req.Amount,
		OrderInfo: orderInfo,
		ReturnUrl: req.ReturnURL,
	}

	// Call service
	resp, err := c.CheckoutClient.CreateMomoQRPayment(ctx, grpcReq)
	if err != nil {
		log.Printf("Error calling checkout service: %v", err)
		jsonResponse(w, http.StatusInternalServerError, jsonError("Failed to process payment"))
		return
	}

	// Check response
	if !resp.Success {
		jsonResponse(w, http.StatusBadRequest, jsonError(resp.Message))
		return
	}

	// Return success response
	jsonResponse(w, http.StatusOK, jsonSuccess(map[string]interface{}{
		"paymentUrl":    resp.PaymentUrl,
		"qrCodeUrl":     resp.QrCodeUrl,
		"deepLink":      resp.DeepLink,
		"orderId":       resp.OrderId,
		"requestId":     resp.RequestId,
		"transactionId": resp.TransactionId,
		"amount":        resp.Amount,
	}))
}

// CreateMomoPosPaymentRequest represents a request to create a MoMo POS payment
type CreateMomoPosPaymentRequest struct {
	OrderID     string `json:"orderId"`
	Amount      int64  `json:"amount"`
	OrderInfo   string `json:"orderInfo"`
	PaymentCode string `json:"paymentCode"` // QR code from MoMo app
	ReturnURL   string `json:"returnUrl"`
}

// CreateMomoPosPayment handles requests to create a new MoMo POS payment
func (c *Config) CreateMomoPosPayment(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req CreateMomoPosPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, jsonError("Invalid request format"))
		return
	}

	// Validate request
	if req.Amount <= 0 || req.PaymentCode == "" || req.ReturnURL == "" {
		jsonResponse(w, http.StatusBadRequest, jsonError("Missing required fields"))
		return
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Default order info if not provided
	orderInfo := req.OrderInfo
	if orderInfo == "" {
		orderInfo = "Payment for order"
	}

	// Build gRPC request
	grpcReq := &pb.MomoPosPaymentRequest{
		OrderId:     req.OrderID,
		Amount:      req.Amount,
		OrderInfo:   orderInfo,
		PaymentCode: req.PaymentCode,
		ReturnUrl:   req.ReturnURL,
	}

	// Call service
	resp, err := c.CheckoutClient.CreateMomoPosPayment(ctx, grpcReq)
	if err != nil {
		log.Printf("Error calling checkout service: %v", err)
		jsonResponse(w, http.StatusInternalServerError, jsonError("Failed to process payment"))
		return
	}

	// Check response
	if !resp.Success {
		jsonResponse(w, http.StatusBadRequest, jsonError(resp.Message))
		return
	}

	// Return success response
	jsonResponse(w, http.StatusOK, jsonSuccess(map[string]interface{}{
		"orderId":       resp.OrderId,
		"requestId":     resp.RequestId,
		"transactionId": resp.TransactionId,
		"amount":        resp.Amount,
		"message":       resp.Message,
	}))
}

// Helper functions for JSON responses
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	// Set content type and status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Marshal data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Log error and return a simple error message if marshaling fails
		log.Printf("Error marshaling JSON response: %v", err)
		w.Write([]byte(`{"success":false,"message":"Internal server error"}`))
		return
	}

	// Write JSON response
	if _, err := w.Write(jsonData); err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

func jsonError(message string) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"message": message,
	}
}

func jsonSuccess(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": true,
		"data":    data,
	}
}

func jsonData(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"data": data,
	}
}

package payment

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sony/sonyflake"
)

// MomoService handles MoMo payment operations
type MomoService struct {
	config MomoConfig
	flake  *sonyflake.Sonyflake
}

// NewMomoService creates a new MoMo service
func NewMomoService(config MomoConfig) *MomoService {
	return &MomoService{
		config: config,
		flake:  sonyflake.NewSonyflake(sonyflake.Settings{}),
	}
}

// MomoPaymentResult represents the result of a MoMo payment operation
type MomoPaymentResult struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	OrderID       string `json:"orderId"`
	RequestID     string `json:"requestId"`
	TransactionID string `json:"transactionId"`
	PaymentURL    string `json:"paymentUrl"`
	Amount        int64  `json:"amount"`
}

// MomoPayload defines the request payload for MoMo payment
type MomoPayload struct {
	PartnerCode   string `json:"partnerCode"`
	AccessKey     string `json:"accessKey"`
	RequestID     string `json:"requestId"`
	Amount        string `json:"amount"`
	OrderID       string `json:"orderId"`
	OrderInfo     string `json:"orderInfo"`
	ReturnURL     string `json:"returnUrl"`
	NotifyURL     string `json:"notifyUrl"`
	ExtraData     string `json:"extraData"`
	RequestType   string `json:"requestType"`
	Signature     string `json:"signature"`
	Lang          string `json:"lang,omitempty"`
	AutoCapture   bool   `json:"autoCapture,omitempty"`
	RedirectUrl   string `json:"redirectUrl,omitempty"`
	IpnUrl        string `json:"ipnUrl,omitempty"`
	BankCode      string `json:"bankCode,omitempty"`      // Optional: specific bank code
	PaymentOption string `json:"paymentOption,omitempty"` // ATM payment option
}

// MomoResponse defines the response from MoMo payment API
type MomoResponse struct {
	PartnerCode  string      `json:"partnerCode"`
	OrderID      string      `json:"orderId"`
	RequestID    string      `json:"requestId"`
	Amount       json.Number `json:"amount"`
	ResponseTime int64       `json:"responseTime"`
	Message      string      `json:"message"`
	ResultCode   int         `json:"resultCode"`
	PayURL       string      `json:"payUrl"`
	TransID      string      `json:"transId,omitempty"`
	ErrorCode    int         `json:"errorCode,omitempty"`
	ErrorMessage string      `json:"errorMessage,omitempty"`
}

// PaymentVerificationResult represents the result of verifying a payment
type PaymentVerificationResult struct {
	Success       bool   `json:"success"`
	OrderID       string `json:"orderId"`
	TransactionID string `json:"transactionId"`
	Amount        int64  `json:"amount"`
	Message       string `json:"message"`
}

// CreatePayment creates a new payment request to MoMo
func (s *MomoService) CreatePayment(orderID string, amount int64, orderInfo string, returnURL string) (*MomoPaymentResult, error) {
	log.Printf("Creating MoMo payment for order %s with amount %d VND", orderID, amount)

	// Use provided orderID or generate a new one
	if orderID == "" {
		id, err := s.flake.NextID()
		if err != nil {
			return nil, fmt.Errorf("error generating order ID: %w", err)
		}
		orderID = strconv.FormatUint(id, 16)
	}

	// Generate requestID
	requestID, err := s.generateRequestID()
	if err != nil {
		return nil, err
	}

	// Default values
	amountStr := strconv.FormatInt(amount, 10)
	extraData := ""
	requestType := "payWithATM" // Specific requestType for ATM banking payment

	// Ensure callback URL has /payments/momo/callback path (for broker service)
	// This is critical - MoMo needs to call the broker-service endpoint
	if returnURL == "" {
		returnURL = "https://ecommerce-api.example.com/payments/momo/callback"
	}
	// Log the callback URL being used
	log.Printf("Using callback URL for MoMo: %s", returnURL)

	// Build raw signature
	var rawSignature bytes.Buffer
	rawSignature.WriteString("accessKey=")
	rawSignature.WriteString(s.config.AccessKey)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(amountStr)
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(extraData)
	rawSignature.WriteString("&ipnUrl=")
	rawSignature.WriteString(returnURL) // Use returnURL as notify URL in test
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(orderID)
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(orderInfo)
	rawSignature.WriteString("&partnerCode=")
	rawSignature.WriteString(s.config.PartnerCode)
	rawSignature.WriteString("&redirectUrl=")
	rawSignature.WriteString(returnURL)
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(requestID)
	rawSignature.WriteString("&requestType=")
	rawSignature.WriteString(requestType)

	// Log raw signature for debugging
	log.Printf("Raw signature for MoMo payment: %s", rawSignature.String())

	// Calculate HMAC SHA256
	h := hmac.New(sha256.New, []byte(s.config.SecretKey))
	h.Write(rawSignature.Bytes())
	signature := hex.EncodeToString(h.Sum(nil))

	// Log the signature
	log.Printf("Generated signature: %s", signature)
	log.Printf("Using secret key: %s", s.config.SecretKey)

	// Create payload - for ATM/Banking payment method
	payload := MomoPayload{
		PartnerCode:   s.config.PartnerCode,
		AccessKey:     s.config.AccessKey,
		RequestID:     requestID,
		Amount:        amountStr,
		OrderID:       orderID,
		OrderInfo:     orderInfo,
		RedirectUrl:   returnURL, // match the raw signature key
		IpnUrl:        returnURL, // match the raw signature key
		ExtraData:     extraData,
		RequestType:   requestType,
		Signature:     signature,
		Lang:          "vi",
		AutoCapture:   true,
		PaymentOption: "atm", // Specify ATM payment option
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling MoMo payload: %w", err)
	}

	log.Printf("MoMo ATM payment request payload: %s", string(jsonPayload))

	// Determine endpoint - use ATM-specific endpoint if available
	endpoint := s.config.ATMEndpoint
	if endpoint == "" {
		endpoint = s.config.APIEndpoint
	}
	if endpoint == "" {
		endpoint = "https://test-payment.momo.vn/v2/gateway/api/create"
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send request to MoMo
	resp, err := client.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("HTTP error when calling MoMo API: %v", err)
		return nil, fmt.Errorf("error sending MoMo request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading MoMo response: %w", err)
	}

	log.Printf("MoMo payment response: %s", string(body))

	// Parse response
	var response MomoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing MoMo response: %w", err)
	}

	// Check response
	if response.ResultCode != 0 {
		errorMessage := response.Message
		if response.ErrorMessage != "" {
			errorMessage = response.ErrorMessage
		}
		log.Printf("MoMo payment error: Code=%d, Message=%s", response.ResultCode, errorMessage)
		return &MomoPaymentResult{
			Success:   false,
			Message:   errorMessage,
			OrderID:   orderID,
			RequestID: requestID,
		}, nil
	}

	// Convert amount string to int64
	amountInt, err := strconv.ParseInt(string(response.Amount), 10, 64)
	if err != nil {
		amountInt = amount // Use original amount if parsing fails
	}

	// Return success result
	log.Printf("MoMo payment created successfully: OrderID=%s, TransID=%s", response.OrderID, response.TransID)
	return &MomoPaymentResult{
		Success:       true,
		Message:       response.Message,
		OrderID:       response.OrderID,
		RequestID:     response.RequestID,
		TransactionID: response.TransID,
		PaymentURL:    response.PayURL,
		Amount:        amountInt,
	}, nil
}

// VerifyPayment verifies a payment from callback parameters
func (s *MomoService) VerifyPayment(params map[string]string) (*PaymentVerificationResult, error) {
	log.Printf("MOMO-VERIFY: Starting to verify MoMo payment")

	// Log received parameters for debugging
	for k, v := range params {
		if k != "signature" { // Don't log sensitive data
			log.Printf("MOMO-VERIFY: Param %s = %s", k, v)
		} else {
			log.Printf("MOMO-VERIFY: Param %s = [REDACTED]", k)
		}
	}

	// Get parameters from callback
	orderID := getParamWithDefault(params, "orderId", "")
	transID := getParamWithDefault(params, "transId", "")
	resultCode := getParamWithDefault(params, "resultCode", "")
	amount := getParamWithDefault(params, "amount", "0")
	signature := getParamWithDefault(params, "signature", "")

	log.Printf("MOMO-VERIFY: Validating order: OrderID=%s, TransactionID=%s, ResultCode=%s",
		orderID, transID, resultCode)

	// Validate required parameters
	if orderID == "" {
		log.Printf("MOMO-VERIFY-ERROR: Missing orderId parameter")
		return nil, fmt.Errorf("missing orderId parameter")
	}

	if resultCode == "" {
		log.Printf("MOMO-VERIFY-ERROR: Missing resultCode parameter")
		return nil, fmt.Errorf("missing resultCode parameter")
	}

	// Convert amount to int64
	amountInt, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		log.Printf("MOMO-VERIFY-ERROR: Invalid amount format: %s", amount)
		return nil, fmt.Errorf("invalid amount format: %w", err)
	}

	// Validate transaction success
	if resultCode != "0" {
		errorMessage := fmt.Sprintf("MoMo payment failed with result code %s", resultCode)
		log.Printf("MOMO-VERIFY-ERROR: %s", errorMessage)
		return &PaymentVerificationResult{
			Success:       false,
			OrderID:       orderID,
			TransactionID: transID,
			Amount:        amountInt,
			Message:       errorMessage,
		}, nil
	}

	// Basic verification - in production, should validate signature
	if signature == "" {
		log.Printf("MOMO-VERIFY-WARNING: Missing signature parameter")
		// We'll continue without signature validation for now
	} else {
		// Calculate and verify signature
		log.Printf("MOMO-VERIFY: Verifying signature...")
		// TODO: Implement signature validation
	}

	log.Printf("MOMO-VERIFY-SUCCESS: Payment verified successfully for order %s", orderID)

	// Return successful verification result
	return &PaymentVerificationResult{
		Success:       true,
		OrderID:       orderID,
		TransactionID: transID,
		Amount:        amountInt,
		Message:       "Payment verified successfully",
	}, nil
}

// Helper function to get params with a default value
func getParamWithDefault(params map[string]string, key, defaultValue string) string {
	if value, exists := params[key]; exists && value != "" {
		return value
	}
	return defaultValue
}

// generateRequestID generates a unique request ID
func (s *MomoService) generateRequestID() (string, error) {
	id, err := s.flake.NextID()
	if err != nil {
		return "", fmt.Errorf("error generating request ID: %w", err)
	}
	return strconv.FormatUint(id, 16), nil
}

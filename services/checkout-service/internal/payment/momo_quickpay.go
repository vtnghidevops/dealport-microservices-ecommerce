/*
This file contains MoMo QuickPay implementation which is not currently used.
We're using the standard captureWallet (ATM/Banking) implementation from momo.go instead.
To use this implementation, uncomment the relevant sections in the payment handler.
*/

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

	"github.com/sony/sonyflake"
)

// MomoQuickPayService handles MoMo QuickPay operations
type MomoQuickPayService struct {
	config MomoConfig
	flake  *sonyflake.Sonyflake
}

// NewMomoQuickPayService creates a new MoMo QuickPay service
func NewMomoQuickPayService(config MomoConfig) *MomoQuickPayService {
	return &MomoQuickPayService{
		config: config,
		flake:  sonyflake.NewSonyflake(sonyflake.Settings{}),
	}
}

// MomoQuickPayPayload defines the QuickPay request payload
type MomoQuickPayPayload struct {
	PartnerCode  string `json:"partnerCode"`
	AccessKey    string `json:"accessKey"`
	RequestID    string `json:"requestId"`
	Amount       string `json:"amount"`
	OrderID      string `json:"orderId"`
	OrderInfo    string `json:"orderInfo"`
	PartnerName  string `json:"partnerName"`
	StoreId      string `json:"storeId"`
	OrderGroupId string `json:"orderGroupId"`
	Lang         string `json:"lang"`
	AutoCapture  bool   `json:"autoCapture"`
	RedirectUrl  string `json:"redirectUrl"`
	IpnUrl       string `json:"ipnUrl"`
	ExtraData    string `json:"extraData"`
	PaymentCode  string `json:"paymentCode"`
	RequestType  string `json:"requestType"`
	Signature    string `json:"signature"`
}

// MomoQuickPayResponse defines the QuickPay response
type MomoQuickPayResponse struct {
	PartnerCode  string `json:"partnerCode"`
	OrderID      string `json:"orderId"`
	RequestID    string `json:"requestId"`
	Amount       int64  `json:"amount"`
	ResponseTime int64  `json:"responseTime"`
	Message      string `json:"message"`
	ResultCode   int    `json:"resultCode"`
	PayUrl       string `json:"payUrl,omitempty"`
	QrCodeUrl    string `json:"qrCodeUrl,omitempty"`
	DeepLink     string `json:"deeplink,omitempty"`
	TransID      string `json:"transId,omitempty"`
}

// MomoQuickPayResult represents the result of a QuickPay operation
type MomoQuickPayResult struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	OrderID       string `json:"orderId"`
	RequestID     string `json:"requestId"`
	TransactionID string `json:"transactionId"`
	PaymentURL    string `json:"paymentUrl"`
	QrCodeURL     string `json:"qrCodeUrl"`
	DeepLink      string `json:"deeplink"`
	Amount        int64  `json:"amount"`
}

// CreatePosPayment creates a new QuickPay payment for POS
func (s *MomoQuickPayService) CreatePosPayment(
	orderID string,
	amount int64,
	orderInfo string,
	paymentCode string,
	returnURL string,
) (*MomoQuickPayResult, error) {
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
	storeID := "Test Store"
	partnerName := s.config.PartnerName
	if partnerName == "" {
		partnerName = "MoMo Payment"
	}
	extraData := ""
	orderGroupID := ""
	autoCapture := true
	lang := "vi"

	// Build raw signature
	var rawSignature bytes.Buffer
	rawSignature.WriteString("accessKey=")
	rawSignature.WriteString(s.config.AccessKey)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(amountStr)
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(extraData)
	rawSignature.WriteString("&orderId=")
	rawSignature.WriteString(orderID)
	rawSignature.WriteString("&orderInfo=")
	rawSignature.WriteString(orderInfo)
	rawSignature.WriteString("&partnerCode=")
	rawSignature.WriteString(s.config.PartnerCode)
	rawSignature.WriteString("&paymentCode=")
	rawSignature.WriteString(paymentCode)
	rawSignature.WriteString("&requestId=")
	rawSignature.WriteString(requestID)

	// Log raw signature for debugging
	log.Printf("Raw signature: %s", rawSignature.String())

	// Calculate HMAC SHA256
	h := hmac.New(sha256.New, []byte(s.config.SecretKey))
	h.Write(rawSignature.Bytes())
	signature := hex.EncodeToString(h.Sum(nil))

	// Create payload
	payload := MomoQuickPayPayload{
		PartnerCode:  s.config.PartnerCode,
		AccessKey:    s.config.AccessKey,
		RequestID:    requestID,
		Amount:       amountStr,
		OrderID:      orderID,
		StoreId:      storeID,
		PartnerName:  partnerName,
		OrderGroupId: orderGroupID,
		AutoCapture:  autoCapture,
		Lang:         lang,
		OrderInfo:    orderInfo,
		RedirectUrl:  returnURL,
		IpnUrl:       returnURL, // use same URL for IPN in test
		ExtraData:    extraData,
		PaymentCode:  paymentCode,
		Signature:    signature,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling MoMo payload: %w", err)
	}

	log.Printf("MoMo QuickPay POS request payload: %s", string(jsonPayload))
	log.Printf("MoMo QuickPay POS signature: %s", signature)

	// Determine endpoint
	endpoint := s.config.PosEndpoint
	if endpoint == "" {
		endpoint = "https://test-payment.momo.vn/v2/gateway/api/pos"
	}

	// Send request to MoMo
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error sending MoMo request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading MoMo response: %w", err)
	}

	log.Printf("MoMo QuickPay POS response: %s", string(body))

	// Parse response
	var response MomoQuickPayResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing MoMo response: %w", err)
	}

	// Check response
	if response.ResultCode != 0 {
		return nil, fmt.Errorf("MoMo error: %s (code: %d)", response.Message, response.ResultCode)
	}

	// Return success result
	return &MomoQuickPayResult{
		Success:       true,
		Message:       response.Message,
		OrderID:       response.OrderID,
		RequestID:     response.RequestID,
		TransactionID: response.TransID,
		PaymentURL:    response.PayUrl,
		QrCodeURL:     response.QrCodeUrl,
		DeepLink:      response.DeepLink,
		Amount:        response.Amount,
	}, nil
}

// Create QR payment creates a new QR code payment
func (s *MomoQuickPayService) CreateQRPayment(
	orderID string,
	amount int64,
	orderInfo string,
	returnURL string,
) (*MomoQuickPayResult, error) {
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
	storeID := "Test Store"
	partnerName := s.config.PartnerName
	if partnerName == "" {
		partnerName = "MoMo Payment"
	}
	extraData := ""
	orderGroupID := ""
	autoCapture := true
	lang := "vi"
	requestType := "captureWallet"

	// Build raw signature
	var rawSignature bytes.Buffer
	rawSignature.WriteString("accessKey=")
	rawSignature.WriteString(s.config.AccessKey)
	rawSignature.WriteString("&amount=")
	rawSignature.WriteString(amountStr)
	rawSignature.WriteString("&extraData=")
	rawSignature.WriteString(extraData)
	rawSignature.WriteString("&ipnUrl=")
	rawSignature.WriteString(returnURL)
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
	log.Printf("Raw signature: %s", rawSignature.String())

	// Calculate HMAC SHA256
	h := hmac.New(sha256.New, []byte(s.config.SecretKey))
	h.Write(rawSignature.Bytes())
	signature := hex.EncodeToString(h.Sum(nil))

	// Create payload
	payload := MomoQuickPayPayload{
		PartnerCode:  s.config.PartnerCode,
		AccessKey:    s.config.AccessKey,
		RequestID:    requestID,
		Amount:       amountStr,
		OrderID:      orderID,
		StoreId:      storeID,
		PartnerName:  partnerName,
		OrderGroupId: orderGroupID,
		AutoCapture:  autoCapture,
		Lang:         lang,
		OrderInfo:    orderInfo,
		RedirectUrl:  returnURL,
		IpnUrl:       returnURL, // use same URL for IPN in test
		ExtraData:    extraData,
		RequestType:  requestType,
		Signature:    signature,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling MoMo payload: %w", err)
	}

	log.Printf("MoMo QuickPay QR request payload: %s", string(jsonPayload))
	log.Printf("MoMo QuickPay QR signature: %s", signature)

	// Determine endpoint
	endpoint := s.config.QREndpoint
	if endpoint == "" {
		endpoint = "https://test-payment.momo.vn/v2/gateway/api/create"
	}

	// Send request to MoMo
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error sending MoMo request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading MoMo response: %w", err)
	}

	log.Printf("MoMo QuickPay QR response: %s", string(body))

	// Parse response
	var response MomoQuickPayResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing MoMo response: %w", err)
	}

	// Check response
	if response.ResultCode != 0 {
		return nil, fmt.Errorf("MoMo error: %s (code: %d)", response.Message, response.ResultCode)
	}

	// Return success result
	return &MomoQuickPayResult{
		Success:       true,
		Message:       response.Message,
		OrderID:       response.OrderID,
		RequestID:     response.RequestID,
		TransactionID: response.TransID,
		PaymentURL:    response.PayUrl,
		QrCodeURL:     response.QrCodeUrl,
		DeepLink:      response.DeepLink,
		Amount:        response.Amount,
	}, nil
}

// generateRequestID generates a unique request ID
func (s *MomoQuickPayService) generateRequestID() (string, error) {
	id, err := s.flake.NextID()
	if err != nil {
		return "", fmt.Errorf("error generating request ID: %w", err)
	}
	return strconv.FormatUint(id, 16), nil
}

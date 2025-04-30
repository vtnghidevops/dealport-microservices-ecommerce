package payment

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// VNPayService handles VNPAY payment operations
type VNPayService struct {
	config VNPayConfig
}

// NewVNPayService creates a new VNPAY payment service
func NewVNPayService(config VNPayConfig) *VNPayService {
	return &VNPayService{
		config: config,
	}
}

// VNPayPaymentResult represents the result of a payment operation
type VNPayPaymentResult struct {
	Success    bool   `json:"success"`
	PaymentURL string `json:"paymentUrl"`
	OrderID    string `json:"orderId"`
	VnpTxnRef  string `json:"vnpTxnRef"`
	Message    string `json:"message"`
}



// CreatePayment creates a new payment request to VNPAY
func (s *VNPayService) CreatePayment(orderID string, amount int64, ipAddr string, returnURL string) (*VNPayPaymentResult, error) {
	// Create a unique transaction reference
	// Format: Order ID + Timestamp
	txnRef := orderID + "_" + GetCurrentDateString()

	// Create VNPAY payment URL
	vnpParams := make(map[string]string)
	vnpParams["vnp_Version"] = s.config.Version
	vnpParams["vnp_Command"] = s.config.Command
	vnpParams["vnp_TmnCode"] = s.config.MerchantID
	vnpParams["vnp_Amount"] = fmt.Sprintf("%d", amount*100) // VNPAY requires amount in VND * 100
	vnpParams["vnp_CreateDate"] = GetCurrentDateString()
	vnpParams["vnp_CurrCode"] = s.config.CurrencyCode
	vnpParams["vnp_IpAddr"] = ipAddr
	vnpParams["vnp_Locale"] = s.config.Locale
	vnpParams["vnp_OrderInfo"] = "Thanh toan don hang " + orderID
	vnpParams["vnp_OrderType"] = "other" // Order type (other/topup/billpayment/...)
	vnpParams["vnp_ReturnUrl"] = returnURL
	vnpParams["vnp_TxnRef"] = txnRef

	// Build query string
	var queryBuilder strings.Builder
	queryBuilder.WriteString(s.config.APIEndpoint)
	queryBuilder.WriteString("?")

	// Sort params by key for consistent signature
	var keys []string
	for k := range vnpParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build the query string
	var signData strings.Builder
	for i, k := range keys {
		if i > 0 {
			queryBuilder.WriteString("&")
			signData.WriteString("&")
		}

		// URL encode the parameter value
		value := url.QueryEscape(vnpParams[k])

		// Add to query string
		queryBuilder.WriteString(fmt.Sprintf("%s=%s", k, value))

		// Add to sign data
		signData.WriteString(fmt.Sprintf("%s=%s", k, vnpParams[k]))
	}

	// Calculate HMAC-SHA512 signature
	signature := s.generateVNPaySignature(signData.String())

	// Append signature to URL
	queryBuilder.WriteString(fmt.Sprintf("&vnp_SecureHash=%s", signature))

	// Return the payment URL and reference
	return &VNPayPaymentResult{
		Success:    true,
		PaymentURL: queryBuilder.String(),
		OrderID:    orderID,
		VnpTxnRef:  txnRef,
		Message:    "Payment URL generated successfully",
	}, nil
}

// VerifyPayment verifies a VNPAY payment callback
func (s *VNPayService) VerifyPayment(params map[string]string) (*PaymentVerificationResult, error) {
	// Extract secure hash from params
	secureHash := params["vnp_SecureHash"]
	if secureHash == "" {
		return nil, errors.New("missing secure hash in VNPAY callback")
	}

	// Make a copy of params without the secure hash for verification
	verifyParams := make(map[string]string)
	for k, v := range params {
		if k != "vnp_SecureHash" && k != "vnp_SecureHashType" {
			verifyParams[k] = v
		}
	}

	// Build the string to sign
	var signData strings.Builder

	// Sort params by key for consistent signature
	var keys []string
	for k := range verifyParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build the sign data string
	for i, k := range keys {
		if i > 0 {
			signData.WriteString("&")
		}
		signData.WriteString(fmt.Sprintf("%s=%s", k, verifyParams[k]))
	}

	// Calculate the expected signature
	expectedSignature := s.generateVNPaySignature(signData.String())

	// Verify signature
	if expectedSignature != secureHash {
		return nil, errors.New("invalid signature in VNPAY callback")
	}

	// Check response code
	responseCode := params["vnp_ResponseCode"]
	if responseCode != "00" {
		return &PaymentVerificationResult{
			Success:       false,
			OrderID:       extractOrderID(params["vnp_TxnRef"]),
			TransactionID: params["vnp_TransactionNo"],
			Amount:        parseVNPayAmount(params["vnp_Amount"]),
			Message:       getVNPayResponseMessage(responseCode),
		}, nil
	}

	// Payment was successful
	return &PaymentVerificationResult{
		Success:       true,
		OrderID:       extractOrderID(params["vnp_TxnRef"]),
		TransactionID: params["vnp_TransactionNo"],
		Amount:        parseVNPayAmount(params["vnp_Amount"]),
		Message:       "Thanh toán thành công",
	}, nil
}

// generateVNPaySignature generates a HMAC-SHA512 signature for VNPAY
func (s *VNPayService) generateVNPaySignature(data string) string {
	hmacObj := hmac.New(sha512.New, []byte(s.config.SecureHash))
	hmacObj.Write([]byte(data))
	return hex.EncodeToString(hmacObj.Sum(nil))
}

// Helper to generate MD5 hash (used in some VNPAY integrations)
func generateMD5Hash(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

// extractOrderID extracts the original order ID from VNPAY txnRef
func extractOrderID(txnRef string) string {
	parts := strings.Split(txnRef, "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return txnRef
}

// parseVNPayAmount parses the amount from VNPAY (amount is in VND * 100)
func parseVNPayAmount(amountStr string) int64 {
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return 0
	}
	return amount / 100 // Convert back to VND
}

// getVNPayResponseMessage returns a human-readable message for VNPAY response codes
func getVNPayResponseMessage(responseCode string) string {
	responses := map[string]string{
		"00": "Giao dịch thành công",
		"01": "Giao dịch đã tồn tại",
		"02": "Merchant không hợp lệ (kiểm tra lại vnp_TmnCode)",
		"03": "Dữ liệu gửi sang không đúng định dạng",
		"04": "Khởi tạo GD không thành công do Website đang bị tạm khóa",
		"05": "Giao dịch không thành công do: Quý khách nhập sai mật khẩu quá số lần quy định",
		"06": "Giao dịch không thành công do Quý khách nhập sai mật khẩu",
		"07": "Giao dịch bị nghi ngờ gian lận",
		"09": "Giao dịch không thành công do: Thẻ/Tài khoản của khách hàng bị khóa",
		"10": "Giao dịch không thành công do: Quý khách xác thực thông tin thẻ không đúng",
		"11": "Giao dịch không thành công do: Đã hết hạn chờ thanh toán",
		"12": "Giao dịch không thành công do: Thẻ/Tài khoản của khách hàng bị khóa",
		"13": "Giao dịch không thành công do Quý khách nhập sai mật khẩu",
		"24": "Giao dịch không thành công do: Khách hàng hủy giao dịch",
		"51": "Giao dịch không thành công do: Tài khoản không đủ số dư",
		"65": "Giao dịch không thành công do: Tài khoản của Quý khách đã vượt quá hạn mức giao dịch trong ngày",
		"75": "Ngân hàng thanh toán đang bảo trì",
		"79": "Giao dịch không thành công do: KH nhập sai mật khẩu thanh toán quá số lần quy định",
		"99": "Lỗi không xác định",
	}

	if message, exists := responses[responseCode]; exists {
		return message
	}
	return "Giao dịch thất bại với mã lỗi: " + responseCode
}

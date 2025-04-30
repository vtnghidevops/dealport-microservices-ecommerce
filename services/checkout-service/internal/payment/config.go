package payment

import (
	"os"
	"strconv"
	"time"
)

// PaymentConfig holds all payment gateway configurations
type PaymentConfig struct {
	MoMo  MomoConfig
	VNPay VNPayConfig
}

// MomoConfig holds MoMo payment gateway configuration
type MomoConfig struct {
	PartnerCode string
	AccessKey   string
	SecretKey   string
	APIEndpoint string
	PosEndpoint string
	QREndpoint  string
	ATMEndpoint string
	PartnerName string
	IsTestMode  bool
}

// VNPayConfig holds VNPAY payment gateway configuration
type VNPayConfig struct {
	MerchantID   string
	SecureHash   string
	APIEndpoint  string
	CurrencyCode string
	Locale       string
	IsTestMode   bool
	Version      string
	Command      string
}

// TEST MOMO ATM:
// Account: NGUYEN VAN A
// Card number: 9704 0000 0000 0018
// Card date: 03/07
// OTP: OTP

// TEST VNPAY:
// (vnp_TmnCode): CGXZLS0Z
// (vnp_HashSecret): XNBCJFAKAZQSGTARRLGCHVZWCIOIGSHN
// Ngân hàng: NCB
// Số thẻ: 9704198526191432198
// Tên chủ thẻ:NGUYEN VAN A
// Ngày phát hành:07/15
// Mật khẩu OTP:123456

// NewPaymentConfig creates a new payment configuration with values from environment variables
func NewPaymentConfig() *PaymentConfig {
	return &PaymentConfig{
		MoMo: MomoConfig{
			PartnerCode: getEnv("MOMO_PARTNER_CODE", "MOMOBKUN20180529"),
			AccessKey:   getEnv("MOMO_ACCESS_KEY", "klm05TvNBzhg7h7j"),
			SecretKey:   getEnv("MOMO_SECRET_KEY", "at67qH6mk8w5Y1nAyMoYKMWACiEi2bsa"),
			APIEndpoint: getEnv("MOMO_API_ENDPOINT", "https://test-payment.momo.vn/v2/gateway/api/create"),
			PosEndpoint: getEnv("MOMO_POS_ENDPOINT", "https://test-payment.momo.vn/v2/gateway/api/pos"),
			QREndpoint:  getEnv("MOMO_QR_ENDPOINT", "https://test-payment.momo.vn/v2/gateway/api/create"),
			ATMEndpoint: getEnv("MOMO_ATM_ENDPOINT", "https://test-payment.momo.vn/v2/gateway/api/create"),
			PartnerName: getEnv("MOMO_PARTNER_NAME", "Test MoMo"),
			IsTestMode:  getEnvAsBool("MOMO_TEST_MODE", true),
		},
		VNPay: VNPayConfig{
			MerchantID:   getEnv("VNPAY_MERCHANT_ID", "VNPAY"),
			SecureHash:   getEnv("VNPAY_SECURE_HASH", "secureHash"),
			APIEndpoint:  getEnv("VNPAY_API_ENDPOINT", "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html"),
			CurrencyCode: getEnv("VNPAY_CURRENCY_CODE", "VND"),
			Locale:       getEnv("VNPAY_LOCALE", "vn"),
			IsTestMode:   getEnvAsBool("VNPAY_TEST_MODE", true),
			Version:      getEnv("VNPAY_VERSION", "2.1.0"),
			Command:      getEnv("VNPAY_COMMAND", "pay"),
		},
	}
}

// Helper functions to read environment variables with fallback values
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolVal, err := strconv.ParseBool(value)
		if err == nil {
			return boolVal
		}
	}
	return fallback
}

// GetCurrentDateString returns formatted date string for VNPAY
func GetCurrentDateString() string {
	return time.Now().Format("20060102150405")
}

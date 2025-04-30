# Checkout Service

This service handles order checkout and payment processing for the e-commerce microservices project.

## Payment Gateways

The service currently supports the following payment gateways:

### 1. MoMo - ATM Banking Method (Direct Bank Transfer)

MoMo implementation uses the payWithATM method, allowing payment via ATM cards or direct bank transfers through the MoMo gateway.

Key features:
- Uses MoMo's Sandbox environment for testing
- Specifically configured for ATM banking payment method
- Implements full payment flow (create payment, redirect, verify payment)
- Properly handles callback verification with signature validation
- Captures detailed payment logs for troubleshooting

#### Test MoMo ATM Account Details:
- Account: NGUYEN VAN A
- Card number: 9704 0000 0000 0018
- Card date: 03/07
- OTP: OTP

### 2. VNPAY 

VNPAY implementation for standard bank card payments.

## Integration Notes

### How to Test MoMo ATM Payments

1. Start the checkout service
2. Create an order and select MoMo as the payment method
3. You will be redirected to the MoMo payment page
4. Select the ATM/Banking option
5. Enter the test card details provided above
6. Complete the payment using any OTP
7. After payment, you will be redirected back to the application with payment status

### Implementation Details

- The MoMo payment implementation is specifically configured for ATM banking payments
- MoMo QuickPay implementation has been commented out
- The `momo_quickpay.go` file has been renamed to `momo_quickpay.go.disabled` to prevent compilation

### Frontend Integration

The frontend communicates with the payment service through the following API endpoints:

1. Create MoMo Payment: `/api/v1/payments/momo/create`
2. Verify MoMo Payment: `/api/v1/payments/momo/verify`

## Configuration

Payment gateway configuration is loaded from environment variables with sensible defaults for testing:

```go
MoMo: MomoConfig{
    PartnerCode: getEnv("MOMO_PARTNER_CODE", "MOMOBKUN20180529"),
    AccessKey:   getEnv("MOMO_ACCESS_KEY", "klm05TvNBzhg7h7j"),
    SecretKey:   getEnv("MOMO_SECRET_KEY", "at67qH6mk8w5Y1nAyMoYKMWACiEi2bsa"),
    APIEndpoint: getEnv("MOMO_API_ENDPOINT", "https://test-payment.momo.vn/v2/gateway/api/create"),
    ATMEndpoint: getEnv("MOMO_ATM_ENDPOINT", "https://test-payment.momo.vn/v2/gateway/api/create"),
    PartnerName: getEnv("MOMO_PARTNER_NAME", "Test MoMo"),
    IsTestMode:  getEnvAsBool("MOMO_TEST_MODE", true),
},
```

## Troubleshooting

If you encounter issues with the MoMo ATM payment integration:

1. Check the payment logs for request/response details
2. Verify the configuration parameters are correctly set
3. Make sure the callback URL is accessible from the internet or properly tunneled
4. Check if the signature generation/validation is working correctly
5. Make sure you're entering the test card details exactly as specified 

# Testing MoMo ATM Payments

To test MoMo ATM payments, follow these steps:

1. Start the checkout service and frontend application
2. Create an order and proceed to checkout
3. Select MoMo as the payment method and click "Place Order"
4. You will be redirected to the MoMo payment page
5. Enter the following test card details:
   - Card number: 9704 0000 0000 0018
   - Card holder: NGUYEN VAN A
   - Card date: 03/07
   - When prompted for OTP, enter any code (e.g., 123456)
6. After successful payment, you will be redirected back to the success page

## Test Credentials

The implementation uses official MoMo sandbox credentials:
- Partner Code: MOMOBKUN20180529
- Access Key: klm05TvNBzhg7h7j
- Secret Key: at67qH6mk8w5Y1nAyMoYKMWACiEi2bsa

These credentials are suitable for the Sandbox environment only. 
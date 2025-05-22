# Broker Service Testing Plan

## Overview
Tài liệu này trình bày chiến lược kiểm thử cho Broker Service, dịch vụ đóng vai trò như API Gateway trong kiến trúc microservices, bao gồm:
- Unit Tests
- Integration Tests
- Mocks cho các client gRPC

## Vai trò của Broker Service
Broker Service đóng vai trò trung gian giữa client và các microservices khác trong hệ thống:
- Nhận các request HTTP từ client
- Chuyển tiếp request đến các service phù hợp thông qua gRPC
- Trả kết quả về cho client

## Phạm vi kiểm thử

### Unit Tests
Unit tests sẽ tập trung vào việc kiểm thử các handler HTTP của Broker Service:
1. **Auth Handlers**: Đăng ký, đăng nhập, xác thực token, refresh token
2. **User Handlers**: CRUD operations cho user
3. **Product Handlers**: Các handler liên quan đến sản phẩm
4. **Cart Handlers**: Quản lý giỏ hàng
5. **Checkout Handlers**: Xử lý thanh toán
6. **Payment Handlers**: Xử lý các phương thức thanh toán

### Integration Tests
Integration tests sẽ tập trung vào luồng xử lý end-to-end, bao gồm:
1. **Auth Flow**: Đăng ký, xác thực OTP, đăng nhập, refresh token
2. **User Flow**: Quản lý thông tin người dùng
3. **Shopping Flow**: Tìm kiếm sản phẩm, thêm vào giỏ hàng, thanh toán

## Phương pháp kiểm thử
Do Broker Service phụ thuộc vào nhiều service khác thông qua gRPC, chúng ta sẽ sử dụng mocks cho các client gRPC để kiểm thử độc lập.

### Mocking gRPC Clients
Chúng ta sẽ tạo các mock cho các client gRPC sau:
- AuthServiceClient
- UserServiceClient
- ProductServiceClient
- CartServiceClient
- CheckoutServiceClient

## Cấu trúc thư mục tests
```
tests/
├── unit/                           # Unit tests
│   ├── http/                       # Tests cho HTTP handlers
│   │   ├── auth_handler_test.go    # Tests cho auth handlers
│   │   ├── user_handler_test.go    # Tests cho user handlers
│   │   ├── product_handler_test.go # Tests cho product handlers
│   │   ├── cart_handler_test.go    # Tests cho cart handlers
│   │   └── checkout_handler_test.go # Tests cho checkout handlers
│   └── util/                       # Tests cho utility functions
│       └── json_util_test.go       # Tests cho JSON utilities
├── integration/                    # Integration tests
│   ├── auth_flow_test.go           # Tests cho auth flow
│   ├── user_flow_test.go           # Tests cho user flow
│   └── shopping_flow_test.go       # Tests cho shopping flow
├── mocks/                          # Mock objects
│   ├── auth_client_mock.go         # Mock cho auth client
│   ├── user_client_mock.go         # Mock cho user client
│   ├── product_client_mock.go      # Mock cho product client
│   ├── cart_client_mock.go         # Mock cho cart client
│   └── checkout_client_mock.go     # Mock cho checkout client
├── helpers/                        # Test helpers
│   └── http_test_utils.go          # Utilities cho HTTP testing
└── TESTING_PLAN.md                 # File này
```

## Kế hoạch triển khai

### Bước 1: Tạo Mock Objects
- Tạo thư mục `tests/mocks`
- Triển khai các mock cho các client gRPC

### Bước 2: Tạo Test Helpers
- Tạo thư mục `tests/helpers`
- Triển khai các utility function cho HTTP testing

### Bước 3: Triển khai Unit Tests
- Tạo thư mục `tests/unit`
- Triển khai các unit test cho các HTTP handler

### Bước 4: Triển khai Integration Tests
- Tạo thư mục `tests/integration`
- Triển khai các integration test cho các luồng chính

## Yêu cầu để chạy tests
- Go 1.18 hoặc cao hơn
- testify package
- gomock package (optional) 
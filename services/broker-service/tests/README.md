# Broker Service Tests

Hướng dẫn này mô tả cách thiết lập và chạy các tests cho Broker Service - dịch vụ đóng vai trò API Gateway trong hệ thống microservices.

## Cấu trúc thư mục tests

```
tests/
├── unit/                           # Unit tests
│   ├── http/                       # Tests cho HTTP handlers
│   │   └── auth_handler_test.go    # Tests cho auth handlers
│   └── util/                       # Tests cho utility functions
├── integration/                    # Integration tests
│   └── auth_flow_test.go           # Tests cho auth flow
├── mocks/                          # Mock objects
│   └── auth_client_mock.go         # Mock cho auth client
├── helpers/                        # Test helpers
│   └── http_test_utils.go          # Utilities cho HTTP testing
├── README.md                       # File này
└── TESTING_PLAN.md                 # Kế hoạch kiểm thử
```

## Yêu cầu

- Go 1.18 hoặc cao hơn
- Testify package: `github.com/stretchr/testify`

## Cài đặt dependencies

```bash
cd services/broker-service
go get github.com/stretchr/testify
```

## Chạy tests

### Chạy tất cả tests

```bash
cd services/broker-service
go test ./tests/... -v
```

### Chạy unit tests

```bash
cd services/broker-service
go test ./tests/unit/... -v
```

### Chạy integration tests

```bash
cd services/broker-service
go test ./tests/integration/... -v
```

### Chạy một file test cụ thể

```bash
cd services/broker-service
go test ./tests/unit/http/auth_handler_test.go -v
```

### Chạy một test case cụ thể

```bash
cd services/broker-service
go test ./tests/unit/http/auth_handler_test.go -run TestLogin -v
```

## Kiểm tra code coverage

```bash
cd services/broker-service
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Phát triển tests mới

1. **Unit Tests**:
   - Tạo file test mới trong thư mục `tests/unit/http/` hoặc `tests/unit/util/`
   - Sử dụng các mock trong thư mục `tests/mocks/`
   - Sử dụng các helper trong thư mục `tests/helpers/`

2. **Integration Tests**:
   - Tạo file test mới trong thư mục `tests/integration/`
   - Tests tích hợp nên tập trung vào luồng xử lý end-to-end

3. **Mocks**:
   - Nếu cần mock mới, tạo file trong thư mục `tests/mocks/`
   - Đảm bảo mock triển khai đúng interface của client gRPC

## Quy tắc viết tests

1. Mỗi test case nên kiểm tra một chức năng cụ thể
2. Sử dụng subtests để tổ chức các test cases liên quan
3. Sử dụng descriptive names cho test cases
4. Cấu trúc test theo mô hình AAA: Arrange, Act, Assert
5. Kiểm tra cả positive cases và negative cases
6. Cleanup sau mỗi test nếu cần

## Troubleshooting

- **Import cycle errors**: Đảm bảo các tests không import package của test
- **Mock failures**: Kiểm tra các expectations và arguments matcher 
# Cart Service Tests

## Overview

Thư mục này chứa các tests cho Cart Service trong ứng dụng e-commerce microservices. Tests được tổ chức theo các danh mục sau:

- **Unit Tests**: Tests cho từng phương thức của service một cách riêng biệt
- **Integration Tests**: Tests cho các luồng hoạt động hoàn chỉnh
- **Helpers**: Các hàm tiện ích để tạo dữ liệu test
- **Mocks**: Các implementation giả lập để sử dụng trong tests

## Cấu trúc Tests

```
tests/
├── unit/                  # Unit tests
│   ├── cart_service_test.go           # Tests cho chức năng giỏ hàng
│   ├── coupon_service_test.go         # Tests cho chức năng mã giảm giá (customer)
│   └── coupon_service_admin_test.go   # Tests cho chức năng mã giảm giá (admin)
├── integration/           # Integration tests
│   ├── cart_workflow_test.go          # Tests cho quy trình giỏ hàng hoàn chỉnh
│   ├── coupon_workflow_test.go        # Tests cho quy trình mã giảm giá
│   └── complex_scenarios_test.go      # Tests cho các tình huống phức tạp
├── helpers/               # Test helpers
│   └── test_data.go                   # Các hàm tạo dữ liệu test
├── mocks/                 # Mock implementations
│   └── repository_mock.go             # Các mock cho repository
├── README.md              # File này
└── TESTING_PLAN.md        # Kế hoạch kiểm thử
```

## Môi trường yêu cầu

- Go 1.19 hoặc cao hơn
- Các dependencies đã được liệt kê trong go.mod

## Chạy Tests

### Chạy tất cả tests

```bash
cd services/cart-service
go test ./tests/... -v
```

### Chạy theo danh mục

```bash
# Chạy unit tests
go test ./tests/unit/... -v

# Chạy integration tests
go test ./tests/integration/... -v
```

### Chạy một file test cụ thể

```bash
go test ./tests/integration/cart_workflow_test.go -v
```

### Chạy một test function cụ thể

```bash
go test ./tests/integration/... -run TestCartWorkflow -v
```

## Kỹ thuật Mock

Cart Service sử dụng mocks từ thư viện `github.com/stretchr/testify/mock`. Tham khảo `repository_mock.go` để xem cách cài đặt mock. Hiểu cách dùng các hàm như:

- `mock.On()`: Thiết lập expectation khi một hàm được gọi
- `mock.Return()`: Định nghĩa giá trị trả về khi hàm được gọi
- `mock.Run()`: Thực thi logic khi hàm được gọi
- `mock.AssertExpectations()`: Kiểm tra xem tất cả các expectations đã được thực hiện

## Test Coverage

Để kiểm tra độ phủ, chạy:

```bash
go test ./tests/... -cover

# Tạo profile chi tiết
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
``` 
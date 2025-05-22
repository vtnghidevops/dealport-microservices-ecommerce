# Cart Service Testing Plan

## Overview
Tài liệu này trình bày chiến lược kiểm thử cho Cart Service, bao gồm:
- Unit Tests (đã triển khai)
- Integration Tests (cần triển khai)
- Test Helpers và Mock Objects

## Phạm vi test hiện tại
Hiện tại, Cart Service đã có các unit tests sau:
- `cart_service_test.go`: Tests cho các chức năng cơ bản của giỏ hàng
- `coupon_service_test.go`: Tests cho chức năng quản lý mã giảm giá
- `coupon_service_admin_test.go`: Tests cho chức năng quản lý mã giảm giá (admin)

## Các phần còn thiếu

### Integration Tests
Chúng ta cần triển khai các integration tests để kiểm tra luồng hoạt động đầy đủ của giỏ hàng:

1. **Cart Workflow Test**:
   - Tạo giỏ hàng mới
   - Thêm sản phẩm vào giỏ hàng
   - Cập nhật số lượng sản phẩm
   - Xóa sản phẩm khỏi giỏ hàng
   - Xóa toàn bộ giỏ hàng

2. **Coupon Workflow Test**:
   - Tạo giỏ hàng và thêm sản phẩm
   - Áp dụng mã giảm giá hợp lệ
   - Kiểm tra giá trị giảm giá được tính đúng
   - Xóa mã giảm giá
   - Thử áp dụng mã giảm giá không hợp lệ

3. **Complex Cart Scenario Test**:
   - Tạo giỏ hàng với nhiều sản phẩm khác nhau
   - Kiểm tra tính toán tổng đơn hàng
   - Áp dụng và xóa mã giảm giá nhiều lần
   - Kiểm tra các trường hợp đặc biệt (giỏ hàng trống, mã giảm giá hết hạn, vv)

### Test helper và util
Cần tạo các helper function để đơn giản hóa việc tạo dữ liệu test như:
- Tạo giỏ hàng mẫu
- Tạo sản phẩm mẫu
- Tạo mã giảm giá mẫu

## Kế hoạch triển khai

### Bước 1: Cải thiện Test Helpers
- Tạo thư mục `tests/helpers` 
- Triển khai các helper function để tạo dữ liệu test
- Đảm bảo các mock đã có đầy đủ chức năng

### Bước 2: Triển khai Integration Tests
- Tạo thư mục `tests/integration`
- Triển khai các integration test theo kịch bản đã định nghĩa ở trên

### Bước 3: Kiểm tra Test Coverage
- Đánh giá độ phủ của tests đối với code
- Bổ sung tests cho các phần còn thiếu

## Cấu trúc thư mục tests
```
tests/
├── unit/                  # Unit tests (đã có)
│   ├── cart_service_test.go
│   ├── coupon_service_test.go
│   └── coupon_service_admin_test.go
├── integration/           # Integration tests (cần triển khai)
│   ├── cart_workflow_test.go
│   ├── coupon_workflow_test.go
│   └── complex_scenarios_test.go
├── helpers/               # Test helpers (cần triển khai)
│   └── test_data.go
├── mocks/                 # Mock implementations (đã có)
│   └── repository_mock.go
├── README.md              # Hướng dẫn chạy tests
└── TESTING_PLAN.md        # File này
```

## Yêu cầu
- Go 1.19 hoặc cao hơn
- Testify package
- Các dependencies khác đã được liệt kê trong go.mod

## Chạy Tests
```bash
# Chạy tất cả tests
go test ./tests/... -v

# Chạy riêng unit tests
go test ./tests/unit/... -v

# Chạy riêng integration tests
go test ./tests/integration/... -v
``` 
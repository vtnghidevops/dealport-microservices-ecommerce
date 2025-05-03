# Mail Service

Dịch vụ gửi email cho hệ thống e-commerce microservices.

## Chức năng

* Nhận sự kiện `email.send` từ RabbitMQ và gửi email tương ứng
* Hỗ trợ nhiều loại email: đăng ký, đặt lại mật khẩu, xác nhận đơn hàng, v.v.
* Sử dụng SMTP để gửi email
* Hỗ trợ cả định dạng HTML và văn bản thuần túy
* Cung cấp API HTTP để test trực tiếp với Postman

## Kiến trúc

### Luồng làm việc chính (Main workflow)

1. Các dịch vụ khác (user-service, order-service, v.v.) gửi sự kiện đến RabbitMQ với routing key `email.send`
2. Mail-service lắng nghe sự kiện `email.send` trên RabbitMQ
3. Khi nhận được sự kiện, mail-service phân tích và gửi email qua SMTP

```
Dịch vụ khác --> RabbitMQ --> Mail Service --> SMTP Server --> Email gửi tới người dùng
```

### Endpoints HTTP cho testing

* `POST /send`: Endpoint để test trực tiếp việc gửi email thông qua Postman
* `POST /handle`: Endpoint để test trực tiếp việc xử lý sự kiện (giống như qua RabbitMQ)
* `GET /ping`: Health check

## Cấu trúc dự án

```
mail-service/
├── cmd/
│   └── api/              # Entry point của dịch vụ
│       ├── main.go       # Khởi tạo server và RabbitMQ consumer
│       └── env.go        # Hỗ trợ đọc biến môi trường
├── internal/
│   ├── config/           # Cấu hình cho dịch vụ
│   └── mailer/           # Package xử lý việc gửi email
├── templates/            # Mẫu email HTML và văn bản
└── Dockerfile            # Cấu hình Docker
```

## Định dạng sự kiện

### Sự kiện RabbitMQ (`email.send`)

```json
{
  "name": "email.send",
  "data": {
    "type": "registration",  // hoặc "reset_password", "password_change", "order_confirmation"
    "to": "user@example.com",
    "subject": "Chào mừng đến với hệ thống",
    "from": "support@example.com",
    "from_name": "E-commerce Support",
    "template": "register.html.gohtml",  // tùy chọn
    "message": "Nội dung email",  // tùy chọn
    ... // dữ liệu khác cho template
  }
}
```

### Test qua HTTP (`POST /send`)

```json
{
  "to": "user@example.com",
  "subject": "Test Subject",
  "message": "Test Message",
  "from": "support@example.com",
  "from_name": "Support Team",
  "template": "mail.html.gohtml",
  "data": {
    "username": "testuser",
    "link": "https://example.com/verify?token=abc123"
  }
}
```

## Cài đặt và chạy

1. Đảm bảo các biến môi trường được cấu hình đúng (xem `.env.example`)
2. Chạy với Docker Compose hoặc riêng lẻ

```bash
# Với Docker Compose
docker-compose up -d

# Hoặc chạy riêng lẻ
go build -o mail-service ./cmd/api
./mail-service
```

## Kiểm tra

Sử dụng Postman để test việc gửi email qua endpoint `/send`:

1. Gửi POST request đến `http://localhost:9002/send`
2. Body: JSON với định dạng như mô tả trong phần "Test qua HTTP" 
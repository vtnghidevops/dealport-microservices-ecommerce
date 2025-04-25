# Hướng dẫn Migrate Product Service từ REST sang gRPC

Để hoàn thành việc migrate từ REST API sang gRPC cho product-service, hãy làm theo các bước sau:

## 1. Cài đặt Protobuf và protoc-gen-go

Trước tiên, bạn cần cài đặt công cụ Protobuf và gRPC:

```shell
# Cài đặt protoc compiler
brew install protobuf # hoặc cách khác tùy thuộc vào OS

# Cài đặt các plugin Go cho protoc
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

## 2. Tạo định nghĩa Proto

Tạo file `proto/product/product.proto` với định nghĩa service và message:

- Đã hoàn thành trong quá trình migration

## 3. Generate Go code từ Proto

Chạy script sau để tạo code từ proto file:

```shell
chmod +x scripts/generate_proto.sh
./scripts/generate_proto.sh
```

Hoặc chạy lệnh sau:

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/product/product.proto
```

## 4. Fix lỗi trong domain model

Cập nhật domain model để phản ánh các thay đổi trong ProductReview:

```go
// ProductReview struct trong domain/product.go
type ProductReview struct {
    ID         int       `json:"id"`
    ProductID  int       `json:"product_id"`
    UserName   string    `json:"user_name"`
    Email      string    `json:"email"`  // Thêm field này
    Avatar     string    `json:"avatar"` // Thêm field này
    Rating     int       `json:"rating"`
    ReviewDate time.Time `json:"review_date"`
    ReviewText string    `json:"review_text"`
}
```

## 5. Cập nhật main.go để kích hoạt gRPC server

Cập nhật file `cmd/api/main.go` để khởi động cả REST và gRPC server:

```go
// Thêm dòng này trong Register section
pb.RegisterProductServiceServer(grpcServer, productGrpcServer)
```

## 6. Triển khai gRPC client trong broker-service

Tạo file gRPC client trong broker-service để kết nối với product-service thông qua gRPC:

```go
// services/broker-service/internal/proxy/grpc/product_client.go
package grpc

import (
    "context"
    "log"
    pb "broker-service/proto/product"
    "google.golang.org/grpc"
)

const (
    productServiceAddress = "localhost:50051"
)

// ProductClient là wrapper cho gRPC client của product service
type ProductClient struct {
    client pb.ProductServiceClient
    conn   *grpc.ClientConn
}

// NewProductClient tạo client mới kết nối tới gRPC product service
func NewProductClient() (*ProductClient, error) {
    conn, err := grpc.Dial(productServiceAddress, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }

    client := pb.ProductServiceClient(conn)
    return &ProductClient{
        client: client,
        conn:   conn,
    }, nil
}

// Close đóng kết nối gRPC
func (c *ProductClient) Close() {
    if c.conn != nil {
        c.conn.Close()
    }
}

// GetProduct lấy thông tin sản phẩm theo ID
func (c *ProductClient) GetProduct(ctx context.Context, id int32) (*pb.Product, error) {
    return c.client.GetProduct(ctx, &pb.GetProductRequest{Id: id})
}

// Thêm các phương thức khác cho các operation của product service
```

## 7. Cập nhật các proxy handler trong broker-service

Thay thế các proxy REST handler bằng các handler sử dụng gRPC client:

```go
// services/broker-service/internal/proxy/http/product.go

// Ví dụ cho GetProductByID
func (c *Config) GetProductByID(w http.ResponseWriter, r *http.Request) {
    // Parse id from URL
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    // Tạo gRPC client
    grpcClient, err := grpc.NewProductClient()
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    defer grpcClient.Close()

    // Gọi gRPC service
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    product, err := grpcClient.GetProduct(ctx, int32(id))
    if err != nil {
        http.Error(w, "Error fetching product", http.StatusInternalServerError)
        return
    }

    // Chuyển đổi response từ gRPC thành JSON và trả về
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(product)
}
```

## 8. Kiểm tra và triển khai

1. Khởi động `product-service` để kích hoạt cả REST và gRPC server
2. Cập nhật `broker-service` để sử dụng gRPC client
3. Kiểm tra endpoints trong frontend để đảm bảo mọi thứ hoạt động như mong đợi

## Lưu ý

Quá trình migration này nên được thực hiện dần dần:

1. Đầu tiên, triển khai gRPC server song song với REST API
2. Sau đó, chuyển đổi từng endpoint trong broker-service sang sử dụng gRPC
3. Khi tất cả đã được chuyển đổi, có thể loại bỏ REST API trong product-service

Nếu gặp vấn đề trong quá trình migration, hãy kiểm tra logs và đảm bảo protobuf definitions đã được generate đúng cách. 
// // services/product-service/internal/transport/grpc_server.go
// package transport

// import (
// 	"context"
// 	"github.com/your-org/ecommerce/api/product"
// 	"github.com/your-org/ecommerce/services/product-service/internal/service"
// )

// type GrpcServer struct {
// 	product.UnimplementedProductServiceServer
// 	productService service.ProductService
// }

// func NewGrpcServer(productService service.ProductService) *GrpcServer {
// 	return &GrpcServer{
// 		productService: productService,
// 	}
// }

// func (s *GrpcServer) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.ProductResponse, error) {
// 	prod, err := s.productService.GetByID(ctx, req.Id)
// 	if err != nil {
// 		return nil, err
// 	}
	
// 	return mapProductToProto(prod), nil
// }

// // Other RPC implementations...;
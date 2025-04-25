// cmd/api/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"product-service/internal/domain"
	"product-service/internal/handler"
	"product-service/internal/repository/postgres"
	"product-service/internal/service"
	transportGrpc "product-service/internal/transport/grpc"
	transportHttp "product-service/internal/transport/http"
	"time"

	pb "product-service/proto/product"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"google.golang.org/grpc"
)

const (
	httpPort = "8082"
	grpcPort = "50051"
)

func main() {
	log.Printf("Starting product service - gRPC port: %s, HTTP port: %s", grpcPort, httpPort)

	// Connect to database
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to Postgres")
	}
	defer conn.Close()

	// Create repositories
	productRepo := postgres.NewProductRepository(conn)
	categoryRepo := postgres.NewCategoryRepository(conn)
	bannerRepo := postgres.NewBannerRepository(conn)
	adsRepo := postgres.NewAdsRepository(conn)

	// Create services
	productService := service.NewProductService(productRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	bannerService := service.NewBannerService(bannerRepo)
	adsService := service.NewAdsService(adsRepo)

	// Create a channel to catch errors
	errCh := make(chan error, 2)

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("Starting gRPC server on port %s...", grpcPort)
		errCh <- startGRPCServer(productService, categoryService, bannerService, adsService, grpcPort)
	}()

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on port %s...", httpPort)

		// Create handler config
		handlerConfig := &handler.Config{
			ProductService:  productService,
			CategoryService: categoryService,
			BannerService:   bannerService,
			AdsService:      adsService,
		}

		// Create HTTP server
		httpServer := transportHttp.NewServer(handlerConfig)

		// Start HTTP server
		errCh <- http.ListenAndServe(fmt.Sprintf(":%s", httpPort), httpServer.Routes())
	}()

	// Block until we get an error from one of the servers
	log.Fatalf("Server error: %v", <-errCh)
}

func startGRPCServer(
	productService domain.ProductService,
	categoryService domain.CategoryService,
	bannerService domain.BannerService,
	adsService domain.AdsService,
	port string,
) error {
	// Create TCP listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen for gRPC: %w", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Create and register product service server
	productGrpcServer := transportGrpc.NewGrpcServer(productService)
	categoryGrpcServer := transportGrpc.NewCategoryGrpcServer(categoryService)
	bannerGrpcServer := transportGrpc.NewBannerGrpcServer(bannerService)
	adsGrpcServer := transportGrpc.NewAdsGrpcServer(adsService)

	// Log before registration
	log.Printf("Registering services with gRPC server...")

	// Register the services with the gRPC server
	pb.RegisterProductServiceServer(grpcServer, productGrpcServer)
	pb.RegisterCategoryServiceServer(grpcServer, categoryGrpcServer)
	pb.RegisterBannerServiceServer(grpcServer, bannerGrpcServer)
	pb.RegisterAdsServiceServer(grpcServer, adsGrpcServer)

	log.Printf("gRPC server listening on port %s", port)

	// Start gRPC server (this blocks)
	return grpcServer.Serve(lis)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	// For development, provide a fallback if the environment variable is not set
	if dsn == "" {
		// dsn = "host=postgres-products port=5432 user=postgres-products password=password dbname=products sslmode=disable timezone=UTC connect_timeout=5"
		dsn = "host=localhost port=5433 user=postgres-products password=password dbname=products sslmode=disable timezone=UTC connect_timeout=5"
	}

	// Try to connect to the database with retries
	for i := 0; i < 10; i++ {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("PostgreSQL not ready yet...")
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Connected to PostgreSQL!")
		return connection
	}

	log.Println("Could not connect to PostgreSQL after multiple attempts")
	return nil
}

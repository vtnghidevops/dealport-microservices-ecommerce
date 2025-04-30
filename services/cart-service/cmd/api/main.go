package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cart-service/internal/config"
	"cart-service/internal/domain"
	"cart-service/internal/repository/redis"
	"cart-service/internal/service"
	transportGrpc "cart-service/internal/transport/grpc"
	pb "cart-service/proto/cart"

	redisClient "github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting cart service")

	// Load config
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to Redis
	redisConn := connectToRedis(cfg.Redis)
	defer redisConn.Close()

	// Create repository
	cartRepo := redis.NewCartRepository(redisConn)

	// Create services
	couponService := service.NewCouponService()
	cartService := service.NewCartService(cartRepo, couponService)

	// Start gRPC server
	go startGRPCServer(cartService, cfg.Server.GRPCPort)

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down cart service...")
	time.Sleep(1 * time.Second) // Allow some time for cleanup
}

func startGRPCServer(cartService domain.CartService, port string) {
	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register service
	cartGrpcServer := transportGrpc.NewServer(cartService)
	pb.RegisterCartServiceServer(grpcServer, cartGrpcServer)

	// Register reflection service for gRPC tools and clients
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on port %s", port)

	// Start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func connectToRedis(cfg config.RedisConfig) *redisClient.Client {
	// Create Redis client
	client := redisClient.NewClient(&redisClient.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping Redis
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully!")
	return client
}

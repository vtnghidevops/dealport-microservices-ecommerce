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

	"checkout-service/internal/config"
	"checkout-service/internal/repository"
	"checkout-service/internal/transport/grpc"
	pb_checkout "checkout-service/proto/checkout"
	pb_payment "checkout-service/proto/payment"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	grpc_server "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to MongoDB
	mongoClient, err := connectToMongoDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Create repository and service
	orderRepo := repository.NewOrderRepository(mongoClient.Database(cfg.MongoDB.Database))
	orderService := repository.NewOrderService(orderRepo)

	// Create gRPC server with recovery middleware
	opts := []grpc_server.ServerOption{
		grpc_server.UnaryInterceptor(recoverInterceptor),
	}
	grpcServer := grpc_server.NewServer(opts...)

	// Register checkout service handler
	checkoutHandler := repository.NewCheckoutServiceHandler(orderService)
	pb_checkout.RegisterCheckoutServiceServer(grpcServer, checkoutHandler)

	// Register payment service handler
	// This implementation uses MoMo's captureWallet (ATM/banking) method for payments
	// and VNPAY's standard payment gateway - QuickPay implementation is commented out
	paymentHandler := grpc.NewPaymentHandler()
	pb_payment.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	// Enable reflection for development tools
	reflection.Register(grpcServer)

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", cfg.Server.GRPCPort, err)
	}

	log.Printf("Checkout Service gRPC server starting on %s", cfg.Server.GRPCPort)

	// Start server in a goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(grpcServer)
}

// Connect to MongoDB
func connectToMongoDB(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(cfg.MongoDB.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to MongoDB at %s", cfg.MongoDB.URI)
	return client, nil
}

// Recovery interceptor for gRPC
func recoverInterceptor(ctx context.Context, req interface{}, info *grpc_server.UnaryServerInfo, handler grpc_server.UnaryHandler) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()
	return handler(ctx, req)
}

// Graceful shutdown
func gracefulShutdown(server *grpc_server.Server) {
	// Wait for interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	sig := <-shutdown

	log.Printf("Received %s signal, initiating graceful shutdown", sig)
	server.GracefulStop()
	log.Println("Checkout Service shutdown complete")
}

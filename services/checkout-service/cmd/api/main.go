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
	"checkout-service/internal/event"
	"checkout-service/internal/repository"
	"checkout-service/internal/service"
	"checkout-service/internal/transport/grpc"
	pb_checkout "checkout-service/proto/checkout"
	pb_payment "checkout-service/proto/payment"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	grpc_server "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Set up logger
	logger := log.New(os.Stdout, "checkout-service ", log.LstdFlags)

	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to MongoDB
	mongoClient, err := connectToMongoDB(cfg)
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Connect to RabbitMQ synchronously rather than in a goroutine
	logger.Printf("EVENT-DEBUG: Attempting to connect to RabbitMQ...")
	rabbitConn, err := connectToRabbitMQ(cfg, logger)
	if err != nil {
		logger.Printf("EVENT-ERROR: Failed to connect to RabbitMQ: %v", err)
		logger.Println("EVENT-WARNING: Continuing without event emitter - order events will not be published")
		rabbitConn = nil
	} else {
		logger.Printf("EVENT-DEBUG: Successfully connected to RabbitMQ!")
	}

	// Create event emitter if RabbitMQ connection is successful
	var eventEmitter *event.Emitter
	if rabbitConn != nil {
		logger.Printf("EVENT-DEBUG: Creating event emitter...")
		eventEmitter, err = event.NewEmitter(rabbitConn, logger)
		if err != nil {
			logger.Printf("EVENT-ERROR: Failed to create event emitter: %v", err)
			eventEmitter = nil
		} else {
			logger.Printf("EVENT-SUCCESS: Event emitter created successfully")
		}
	} else {
		logger.Printf("EVENT-ERROR: Cannot create event emitter - RabbitMQ connection is nil")
	}

	// Debug check if eventEmitter is nil
	if eventEmitter == nil {
		logger.Printf("EVENT-ERROR: eventEmitter is nil after setup - events will not be published!")
	} else {
		logger.Printf("EVENT-DEBUG: eventEmitter is properly initialized")
	}

	// Create repository
	orderRepo := repository.NewOrderRepository(mongoClient.Database(cfg.MongoDB.Database))

	// Create order service with event emitter
	orderService := service.NewOrderService(orderRepo, eventEmitter)

	// Debug check
	if eventEmitter == nil {
		logger.Printf("EVENT-ERROR: OrderService created with nil event emitter")
	} else {
		logger.Printf("EVENT-DEBUG: OrderService created with valid event emitter")
	}

	// Create gRPC server with recovery middleware
	opts := []grpc_server.ServerOption{
		grpc_server.UnaryInterceptor(recoverInterceptor),
	}
	grpcServer := grpc_server.NewServer(opts...)

	// Register checkout service handler
	checkoutHandler := grpc.NewCheckoutHandler(orderService)
	pb_checkout.RegisterCheckoutServiceServer(grpcServer, checkoutHandler)

	// Register payment service handler
	paymentHandler := grpc.NewPaymentHandler(orderService)
	pb_payment.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	// Enable reflection for development tools
	reflection.Register(grpcServer)

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.GRPCPort))
	if err != nil {
		logger.Fatalf("Failed to listen on %s: %v", cfg.Server.GRPCPort, err)
	}

	logger.Printf("Checkout Service gRPC server starting on %s", cfg.Server.GRPCPort)

	// Start server in a goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(grpcServer, rabbitConn, logger)
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

// Connect to RabbitMQ
func connectToRabbitMQ(cfg *config.Config, logger *log.Logger) (*amqp.Connection, error) {
	var counts int
	var conn *amqp.Connection
	var err error

	// Default RabbitMQ URL if not specified in config
	rabbitURL := "amqp://guest:guest@rabbitmq:5672"
	if cfg.RabbitMQ.URL != "" {
		rabbitURL = cfg.RabbitMQ.URL
	}

	// Try to connect to RabbitMQ with retry
	for {
		conn, err = amqp.Dial(rabbitURL)
		if err != nil {
			logger.Printf("RabbitMQ not ready yet: %v", err)
			counts++
		} else {
			logger.Printf("Connected to RabbitMQ at %s", rabbitURL)
			break
		}

		if counts > 5 {
			logger.Printf("Failed to connect to RabbitMQ after %d attempts, giving up", counts)
			return nil, err
		}

		logger.Println("Waiting 2 seconds to retry RabbitMQ connection...")
		time.Sleep(2 * time.Second)
	}

	return conn, nil
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
func gracefulShutdown(server *grpc_server.Server, rabbitConn *amqp.Connection, logger *log.Logger) {
	// Wait for interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	sig := <-shutdown

	logger.Printf("Received %s signal, initiating graceful shutdown", sig)

	// Stop gRPC server
	server.GracefulStop()

	// Close RabbitMQ connection if it exists
	if rabbitConn != nil {
		if err := rabbitConn.Close(); err != nil {
			logger.Printf("Error closing RabbitMQ connection: %v", err)
		} else {
			logger.Println("RabbitMQ connection closed")
		}
	}

	logger.Println("Checkout Service shutdown complete")
}

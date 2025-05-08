package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/internal/config"
	"logger-service/internal/domain"
	"logger-service/internal/handlers/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	// Thiết lập logger
	logger := log.New(os.Stdout, "[LOGGER-SVC] ", log.LstdFlags)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	logger.Printf("Starting logger service on port %s", cfg.Server.GRPCPort)

	// Kết nối đến MongoDB
	mongoClient, err := connectToMongo(cfg)
	if err != nil {
		logger.Panic(err)
	}
	client = mongoClient

	// Tạo context để đóng kết nối
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Đóng kết nối khi kết thúc
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Khởi tạo config với domain models
	app := &grpc.Config{
		Models: domain.New(client),
	}

	// Khởi tạo gRPC server
	server := app.NewGrpcServer()

	// Bắt đầu lắng nghe yêu cầu gRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.GRPCPort))
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	// Xử lý thoát graceful
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Println("Received termination signal, shutting down...")
		server.GracefulStop()
	}()

	// Khởi động gRPC server
	logger.Printf("gRPC server is running on port %s", cfg.Server.GRPCPort)
	if err := server.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}

func connectToMongo(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Manually build the connection string with credentials
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=logs",
		cfg.MongoDB.User,
		cfg.MongoDB.Password,
		cfg.MongoDB.Host,
		cfg.MongoDB.Port,
		cfg.MongoDB.Database)

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Debug connection string (don't include in production)
	log.Printf("Connecting to MongoDB with URI: %s", mongoURI)

	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Error connecting to MongoDB:", err)
		return nil, err
	}

	// Kiểm tra kết nối
	if err = c.Ping(ctx, nil); err != nil {
		log.Println("Error pinging MongoDB:", err)
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	// Ensure logs database and collection exist
	_, err = c.Database(cfg.MongoDB.Database).Collection("logs").CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Creating logs collection in %s database", cfg.MongoDB.Database)
		// Error is expected if collection doesn't exist yet
	}

	return c, nil
}

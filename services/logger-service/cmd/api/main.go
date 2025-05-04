package main

import (
	"context"
	"fmt"
	"log"
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

// Default values - can be overridden by environment variables
const (
	defaultMongoURL      = "mongodb://localhost:27018"
	defaultMongoUsername = "admin"
	defaultMongoPassword = "password"
	defaultMongoDatabase = "logs"
	defaultGRPCPort      = "50001"
)

var client *mongo.Client

func main() {
	// Thiết lập logger
	logger := log.New(os.Stdout, "[LOGGER-SVC] ", log.LstdFlags)

	// Get GRPC port from environment or use default
	gRpcPort := os.Getenv("GRPC_PORT")
	if gRpcPort == "" {
		gRpcPort = defaultGRPCPort
	}

	logger.Printf("Starting logger service on port %s", gRpcPort)

	// Kết nối đến MongoDB
	mongoClient, err := connectToMongo()
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
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
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
	logger.Printf("gRPC server is running on port %s", gRpcPort)
	if err := server.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// Get connection settings from environment variables
	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		mongoURL = defaultMongoURL
	}

	mongoUsername := os.Getenv("MONGO_USERNAME")
	if mongoUsername == "" {
		mongoUsername = defaultMongoUsername
	}

	mongoPassword := os.Getenv("MONGO_PASSWORD")
	if mongoPassword == "" {
		mongoPassword = defaultMongoPassword
	}

	mongoDatabase := os.Getenv("MONGO_DATABASE")
	if mongoDatabase == "" {
		mongoDatabase = defaultMongoDatabase
	}

	log.Printf("Connecting to MongoDB at %s with database %s", mongoURL, mongoDatabase)

	// Thiết lập các tùy chọn kết nối
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: mongoUsername,
		Password: mongoPassword,
	})

	// Kết nối
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to MongoDB:", err)
		return nil, err
	}

	// Kiểm tra kết nối
	if err = c.Ping(context.TODO(), nil); err != nil {
		log.Println("Error pinging MongoDB:", err)
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	// Ensure logs database and collection exist
	_, err = c.Database(mongoDatabase).Collection("logs").CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Creating logs collection in %s database", mongoDatabase)
		// Error is expected if collection doesn't exist yet
	}

	return c, nil
}

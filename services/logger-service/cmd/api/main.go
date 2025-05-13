package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/internal/config"
	"logger-service/internal/domain"
	"logger-service/internal/handlers/grpc"
	"logger-service/internal/migrations"
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

	// Process command line arguments for migrations
	if len(os.Args) > 1 {
		handleMigrationCommands(logger)
		return
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Kết nối đến MongoDB
	mongoClient, err := connectToMongo(cfg)
	if err != nil {
		logger.Panic(err)
	}
	client = mongoClient

	// Run migrations if AUTO_MIGRATE is set to true
	if os.Getenv("AUTO_MIGRATE") == "true" {
		logger.Println("AUTO_MIGRATE enabled, running migrations...")
		migrator := migrations.NewMongoDatabaseMigrator(client, cfg.MongoDB.Database)
		if err := migrator.RunMigrations(); err != nil {
			logger.Fatalf("Failed to run migrations: %v", err)
		}
	}

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

// handleMigrationCommands processes migration commands
func handleMigrationCommands(logger *log.Logger) {
	migrationCommand := os.Args[1]
	logger.Printf("Running migration command: %s", migrationCommand)

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Connect to MongoDB
	mongoClient, err := connectToMongo(cfg)
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB for migrations: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Create migrator
	migrator := migrations.NewMongoDatabaseMigrator(mongoClient, cfg.MongoDB.Database)

	// Process the command
	switch migrationCommand {
	case "migrate":
		// Run migrations
		if err := migrator.RunMigrations(); err != nil {
			logger.Fatalf("Migration failed: %v", err)
		}
		logger.Println("MongoDB migrations completed successfully")

	case "rollback":
		// Rollback the last migration
		if err := migrator.RollbackMigration(); err != nil {
			logger.Fatalf("Rollback failed: %v", err)
		}
		logger.Println("Rollback completed successfully")

	case "drop":
		// Drop the database
		if err := migrator.DropDatabase(); err != nil {
			logger.Fatalf("Drop failed: %v", err)
		}
		logger.Println("Database dropped successfully")

	case "status":
		// Get migration status
		status, err := migrator.GetMigrationStatus()
		if err != nil {
			logger.Fatalf("Failed to get migration status: %v", err)
		}
		logger.Printf("MongoDB migration status:\n%s", status)

	default:
		logger.Fatalf("Unknown migration command: %s", migrationCommand)
	}

	logger.Printf("Migration command '%s' completed successfully", migrationCommand)
}

func connectToMongo(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Manually build the connection string with credentials
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
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

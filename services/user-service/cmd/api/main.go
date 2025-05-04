package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"user-service/internal/config"
	"user-service/internal/logging"
	"user-service/internal/repository/postgres"
	"user-service/internal/service"
	grpcHandler "user-service/internal/transport/grpc"
	pb "user-service/proto/user"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set up logger
	logger := log.New(os.Stdout, "[USER-SVC] ", log.LstdFlags)
	logger.Printf("Starting user service on port %s", cfg.Server.Port)

	// Connect to database
	db, err := sqlx.Connect("postgres", cfg.PostgresConnectionString())
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Println("Connected to PostgreSQL database")

	// Ping the database to ensure connection
	if err := db.Ping(); err != nil {
		logger.Fatalf("Failed to ping database: %v", err)
	}

	// Get logger service host from environment or use default
	loggerHost := os.Getenv("LOGGER_SERVICE_HOST")
	if loggerHost == "" {
		loggerHost = "localhost:50001"
	}

	// Initialize the logger client
	var loggerClient *logging.LoggerClient
	loggerClient, err = logging.NewLoggerClient(loggerHost)
	if err != nil {
		logger.Printf("Warning: Failed to initialize logger client: %v", err)
		logger.Println("User activity logging will be disabled")
		loggerClient = nil
	} else {
		logger.Println("Connected to logger service")
		defer loggerClient.Close()
	}

	// Initialize repository
	userRepo := postgres.NewPostgresRepository(db)

	// Initialize service
	userService := service.NewUserService(userRepo, loggerClient)

	// Initialize gRPC handler
	userHandler := grpcHandler.NewUserHandler(userService, logger)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Start listening for gRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Println("Received termination signal, shutting down...")
		grpcServer.GracefulStop()
	}()

	// Start gRPC server
	logger.Printf("gRPC server is running on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}

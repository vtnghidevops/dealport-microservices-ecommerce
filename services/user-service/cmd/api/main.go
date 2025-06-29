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
	"user-service/internal/event"
	"user-service/internal/migrations"
	"user-service/internal/repository/postgres"
	"user-service/internal/service"
	grpcHandler "user-service/internal/transport/grpc"
	pb "user-service/proto/user"
)

func main() {
	// Process command line arguments for migrations
	if len(os.Args) > 1 {
		handleMigrationCommands()
		return
	}

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

	// Run migrations if AUTO_MIGRATE is set to true
	if os.Getenv("AUTO_MIGRATE") == "true" {
		logger.Println("AUTO_MIGRATE enabled, running migrations...")
		if err := migrations.RunMigrations(db); err != nil {
			logger.Fatalf("Failed to run migrations: %v", err)
		}
	}

	// Ping the database to ensure connection
	if err := db.Ping(); err != nil {
		logger.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize event emitter for RabbitMQ
	// Get RabbitMQ URL from environment
	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	eventEmitter, err := event.NewEventEmitter(rabbitURL, logger)
	if err != nil {
		logger.Printf("Warning: Failed to initialize event emitter: %v", err)
		logger.Println("Event publishing will be disabled")
		eventEmitter = nil
	} else {
		logger.Println("Event emitter initialized successfully")
		defer eventEmitter.Close()
	}

	// Initialize repository
	userRepo := postgres.NewPostgresRepository(db)

	// Initialize service
	userService := service.NewUserService(userRepo, eventEmitter)

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

// handleMigrationCommands processes the migration related commands
func handleMigrationCommands() {
	migrationCommand := os.Args[1]

	// Load config
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := sqlx.Connect("postgres", cfg.PostgresConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Process the command
	switch migrationCommand {
	case "migrate":
		// Run migrations up
		if err := migrations.RunMigrations(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "rollback":
		// Rollback the last migration
		if err := migrations.RollbackMigration(db); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
	case "drop":
		// Drop all tables (development only)
		if err := migrations.DropAllTables(db); err != nil {
			log.Fatalf("Drop failed: %v", err)
		}
	case "status":
		// Print current migration status
		version, dirty, err := migrations.GetMigrationStatus(db)
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
		log.Printf("Current migration version: %d, dirty: %t", version, dirty)
	default:
		log.Fatalf("Unknown migration command: %s", migrationCommand)
	}

	log.Printf("Migration command '%s' completed successfully", migrationCommand)
}

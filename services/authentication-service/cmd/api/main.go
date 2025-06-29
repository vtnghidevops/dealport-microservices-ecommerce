package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"authentication-service/internal/config"
	"authentication-service/internal/event"
	"authentication-service/internal/migrations"
	"authentication-service/internal/repository/postgres"
	"authentication-service/internal/service"
	grpcHandler "authentication-service/internal/transport/grpc"
	"authentication-service/internal/util"
	pb "authentication-service/proto/auth"
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
	logger := log.New(os.Stdout, "[AUTH-SVC] ", log.LstdFlags)
	logger.Printf("Starting authentication service on port %s", cfg.Server.Port)

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

	// Connect to RabbitMQ
	rabbitConn, err := util.ConnectToRabbitMQ(cfg.RabbitMQ.URL)
	if err != nil {
		logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	// logger.Println("Connected to RabbitMQ")

	// Initialize event emitter
	eventEmitter, err := event.NewEmitter(rabbitConn, logger)
	if err != nil {
		logger.Fatalf("Failed to create event emitter: %v", err)
	}
	// logger.Println("Event emitter initialized")

	// Add these debug logs after initializing event emitter (around line 90)
	// log.Printf("DEBUG STARTUP: Event emitter initialized - conn: %v, logger: %v", rabbitConn != nil, logger != nil)
	// Test event emission during startup
	// log.Printf("DEBUG STARTUP: Testing event emission during startup")
	time.Sleep(5 * time.Second) // Allow service to fully start
	go func() {
		// Create a test event
		testEvent := event.StandardEvent{
			ID:          uuid.New().String(),
			Name:        "test.event",
			Data:        map[string]string{"message": "Test event from auth service startup"},
			DataSchema:  "v1",
			Source:      "authentication-service-startup-test",
			CreatedAt:   time.Now(),
			PublishedAt: time.Now(),
			Version:     "1.0",
		}

		// Try to publish directly to test connection
		if ch, err := rabbitConn.Channel(); err == nil {
			defer ch.Close()

			// Declare exchange
			err := ch.ExchangeDeclare("logs_topic", "topic", true, false, false, false, nil)
			if err != nil {
				log.Printf("ERROR STARTUP: Failed to declare exchange: %v", err)
				return
			}

			// Convert to JSON
			jsonData, err := json.Marshal(testEvent)
			if err != nil {
				log.Printf("ERROR STARTUP: Failed to marshal test event: %v", err)
				return
			}

			// Try to publish
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err = ch.PublishWithContext(
				ctx,
				"logs_topic", // exchange
				"test.event", // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        jsonData,
				},
			)

			if err != nil {
				log.Printf("ERROR STARTUP: Failed to publish test event: %v", err)
			} else {
				log.Printf("SUCCESS STARTUP: Test event published successfully")
			}
		} else {
			log.Printf("ERROR STARTUP: Failed to create channel for test: %v", err)
		}
	}()

	// Initialize OTP manager
	otpManager := util.NewOTPManagerWithBypass(
		cfg.OTP.Length,
		cfg.OTP.Expiry,
		cfg.OTP.MaxAttempts,
		cfg.OTP.BypassEnabled,
		cfg.OTP.BypassCode,
	)

	// Log OTP configuration (hide bypass code in production)
	if cfg.OTP.BypassEnabled {
		logger.Printf("Initialized OTP manager with bypass ENABLED for %s environment", cfg.Server.Environment)
	} else {
		logger.Printf("Initialized OTP manager with bypass DISABLED for %s environment", cfg.Server.Environment)
	}
	logger.Printf("OTP expiry: %v, max attempts: %d", cfg.OTP.Expiry, cfg.OTP.MaxAttempts)

	// Initialize mail client
	//mailClient := util.NewMailClient(cfg.MailClient.BaseURL)
	//logger.Printf("Initialized mail client with base URL: %s", cfg.MailClient.BaseURL)

	// Start a periodic cleanup routine for expired OTPs
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			otpManager.CleanupExpiredOTPs()
		}
	}()

	// Initialize repository
	userRepo := postgres.NewPostgresRepository(db)

	// Initialize service
	authService := service.NewAuthService(
		userRepo,
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessDuration,
		cfg.JWT.RefreshDuration,
		otpManager,
		//mailClient,
		eventEmitter,
	)

	// Initialize gRPC handler
	authHandler := grpcHandler.NewAuthHandler(authService)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Start listening
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

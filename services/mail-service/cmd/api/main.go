package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"mail-service/internal/config"
	"mail-service/internal/grpc"
	pb "mail-service/proto"

	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Config holds configuration for mail service
type Config struct {
	Config config.Config // Service configuration
}

const (
	grpcPort string = "50057" // gRPC port for mail service
)

func main() {
	log.Println("Starting mail service...")

	// Load configuration
	cfg := config.LoadConfig()

	app := Config{
		Config: cfg,
	}

	// Create a logger for gRPC server
	grpcLogger := log.New(os.Stdout, "[MAIL-GRPC] ", log.LstdFlags)

	// Start gRPC server
	go app.serveGRPC(grpcLogger)

	log.Println("Mail service started successfully")
	log.Printf("gRPC server listening on port %s", grpcPort)

	// Keep the application running indefinitely
	select {}
}

// serveGRPC starts the gRPC server
func (app *Config) serveGRPC(logger *log.Logger) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		logger.Fatalf("Failed to listen for gRPC: %v", err)
	}
	defer lis.Close()

	// Create a new gRPC server
	s := grpclib.NewServer()

	// Create mail server
	mailServer := grpc.NewMailServer(app.Config, logger)

	// Register mail server
	pb.RegisterMailServiceServer(s, mailServer)

	// Register reflection service on gRPC server (useful for tools like grpcurl)
	reflection.Register(s)

	logger.Printf("gRPC Server started on port %s", grpcPort)
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve gRPC: %v", err)
	}
}

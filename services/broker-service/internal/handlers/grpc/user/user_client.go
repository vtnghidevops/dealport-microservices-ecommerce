package user

import (
	"context"
	"log"
	"os"
	"time"

	userpb "broker-service/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// GetUserClient returns a gRPC client for the user service
func GetUserClient() (userpb.UserServiceClient, error) {
	// Get user service host from environment variable or use default
	userHost := os.Getenv("USER_SERVICE_HOST")
	if userHost == "" {
		userHost = "ecommerce-user-service:50052" // Use actual service name for k8s environment
	}

	log.Printf("Attempting to connect to user service at %s", userHost)

	// Connection options with timeout and keepalive
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	// Set up connection to user service with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, userHost, opts...)
	if err != nil {
		log.Println("Error connecting to user service:", err)
		return nil, err
	}

	// Create user service client
	userClient := userpb.NewUserServiceClient(conn)
	log.Println("UserClient created successfully")

	return userClient, nil
}

// TestUserConnection tests if user service is actually reachable
func TestUserConnection(client userpb.UserServiceClient) bool {
	if client == nil {
		return false
	}

	log.Println("🔍 Testing actual gRPC connection to user service...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Try a simple request
	_, err := client.GetUsers(ctx, &userpb.GetUsersRequest{
		Page:  1,
		Limit: 1,
	})

	if err != nil {
		log.Printf("❌ User service gRPC test failed: %v", err)
		return false
	}

	log.Println("✅ User service gRPC connection verified!")
	return true
}

// NewUserHandler creates a new instance of the user handler
func NewUserHandler(client userpb.UserServiceClient) *UserHandler {
	return &UserHandler{
		Client: client,
	}
}

// UserHandler handles gRPC communication with the user service
type UserHandler struct {
	Client userpb.UserServiceClient
}

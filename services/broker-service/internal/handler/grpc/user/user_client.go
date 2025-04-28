package user

import (
	"log"
	"os"

	userpb "broker-service/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GetUserClient returns a gRPC client for the user service
func GetUserClient() (userpb.UserServiceClient, error) {
	// Get user service host from environment variable or use default
	userHost := os.Getenv("USER_SERVICE_HOST")
	if userHost == "" {
		userHost = "localhost:50052" // Use localhost for development
	}

	log.Printf("Attempting to connect to user service at %s", userHost)

	// Set up connection to user service - removed WithBlock() to prevent hanging
	conn, err := grpc.Dial(userHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Error connecting to user service:", err)
		return nil, err
	}

	// Create user service client
	userClient := userpb.NewUserServiceClient(conn)
	log.Println("UserClient created successfully")

	return userClient, nil
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

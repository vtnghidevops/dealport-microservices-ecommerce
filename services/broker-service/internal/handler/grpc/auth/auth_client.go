package auth

import (
	"log"
	"os"

	authpb "broker-service/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GetAuthClient returns a gRPC client for the authentication service
func GetAuthClient() (authpb.AuthServiceClient, error) {
	// Get auth service host from environment variable or use default
	authHost := os.Getenv("AUTH_SERVICE_HOST")
	if authHost == "" {
		authHost = "localhost:50051" // Use localhost for development
	}

	log.Printf("Attempting to connect to auth service at %s", authHost)

	// Set up connection to auth service - removed WithBlock() to prevent hanging
	conn, err := grpc.Dial(authHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Error connecting to authentication service:", err)
		return nil, err
	}

	// Create auth service client
	authClient := authpb.NewAuthServiceClient(conn)
	log.Println("AuthClient created successfully")

	return authClient, nil
}

// NewAuthHandler creates a new instance of the authentication handler
func NewAuthHandler(client authpb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		Client: client,
	}
}

// AuthHandler handles gRPC communication with the authentication service
type AuthHandler struct {
	Client authpb.AuthServiceClient
}

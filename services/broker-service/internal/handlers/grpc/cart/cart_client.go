package cart

import (
	"log"
	"os"

	cartpb "broker-service/proto/cart"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GetCartClient returns a gRPC client for the cart service
func GetCartClient() (cartpb.CartServiceClient, error) {
	// Get cart service host from environment variable or use default
	cartHost := os.Getenv("CART_SERVICE_HOST")
	if cartHost == "" {
		cartHost = "cart-service:50054" // Use service name for docker environment
	}

	log.Printf("Attempting to connect to cart service at %s", cartHost)

	// Set up connection to cart service
	conn, err := grpc.Dial(cartHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Error connecting to cart service:", err)
		return nil, err
	}

	// Create cart service client
	cartClient := cartpb.NewCartServiceClient(conn)
	log.Println("CartClient created successfully")

	return cartClient, nil
}

// CartHandler handles gRPC communication with the cart service
type CartHandler struct {
	Client cartpb.CartServiceClient
}

// NewCartHandler creates a new instance of the cart handler
func NewCartHandler(client cartpb.CartServiceClient) *CartHandler {
	return &CartHandler{
		Client: client,
	}
}

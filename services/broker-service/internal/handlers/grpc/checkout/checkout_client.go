package checkout

import (
	"log"
	"os"

	checkoutpb "broker-service/proto/checkout"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GetCheckoutClient returns a gRPC client for the checkout service
func GetCheckoutClient() (checkoutpb.CheckoutServiceClient, error) {
	// Get checkout service host from environment variable or use default
	checkoutHost := os.Getenv("CHECKOUT_SERVICE_HOST")
	if checkoutHost == "" {
		checkoutHost = "localhost:50055" // Default checkout service port
	}

	log.Printf("Attempting to connect to checkout service at %s", checkoutHost)

	// Set up connection to checkout service
	conn, err := grpc.Dial(checkoutHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Error connecting to checkout service:", err)
		return nil, err
	}

	// Create checkout service client
	checkoutClient := checkoutpb.NewCheckoutServiceClient(conn)
	log.Println("CheckoutClient created successfully")

	return checkoutClient, nil
}

// CheckoutHandler handles gRPC communication with the checkout service
type CheckoutHandler struct {
	Client checkoutpb.CheckoutServiceClient
}

// NewCheckoutHandler creates a new instance of the checkout handler
func NewCheckoutHandler(client checkoutpb.CheckoutServiceClient) *CheckoutHandler {
	return &CheckoutHandler{
		Client: client,
	}
}

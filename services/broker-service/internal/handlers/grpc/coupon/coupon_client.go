package coupon

import (
	"log"
	"os"

	couponpb "broker-service/proto/coupon"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GetCouponClient returns a gRPC client for the coupon service (which is part of cart service)
func GetCouponClient() (couponpb.CouponServiceClient, error) {
	// Get cart service host from environment variable or use default
	cartHost := os.Getenv("CART_SERVICE_HOST")
	if cartHost == "" {
		cartHost = "localhost:50054" // Default cart service port
	}

	log.Printf("Attempting to connect to cart service (for coupon operations) at %s", cartHost)

	// Set up connection to cart service
	conn, err := grpc.Dial(cartHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Error connecting to cart service for coupon operations:", err)
		return nil, err
	}

	// Create coupon service client
	couponClient := couponpb.NewCouponServiceClient(conn)
	log.Println("CouponClient created successfully")

	return couponClient, nil
}

// CouponHandler handles gRPC communication with the coupon service
type CouponHandler struct {
	Client couponpb.CouponServiceClient
}

// NewCouponHandler creates a new instance of the coupon handler
func NewCouponHandler(client couponpb.CouponServiceClient) *CouponHandler {
	return &CouponHandler{
		Client: client,
	}
}

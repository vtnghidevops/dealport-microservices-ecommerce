package payment

import (
	"log"
	"os"

	pb "broker-service/proto/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetPaymentClient() (pb.PaymentServiceClient, error) {
	// Get payment service host from environment variable or use default
	paymentHost := os.Getenv("CHECKOUT_SERVICE_HOST")
	if paymentHost == "" {
		paymentHost = "checkout-service:50055" // Use service name for docker environment
	}

	// Connect to the checkout service for payment handling
	conn, err := grpc.Dial(paymentHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Error connecting to checkout service:", err)
		return nil, err
	}

	// Create client
	c := pb.NewPaymentServiceClient(conn)
	return c, nil
}

package payment

import (
	"log"

	pb "broker-service/proto/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetPaymentClient() (pb.PaymentServiceClient, error) {
	// Use environment variable for service URL in a real deployment
	conn, err := grpc.Dial("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Error connecting to checkout service:", err)
		return nil, err
	}

	// Create client
	c := pb.NewPaymentServiceClient(conn)
	return c, nil
}

package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	grpcAuthHandler "broker-service/internal/handler/grpc/auth"
	grpcProductHandler "broker-service/internal/handler/grpc/product"
	grpcUserHandler "broker-service/internal/handler/grpc/user"
	httpAuthHandler "broker-service/internal/handler/http/auth"
	httpProductHandler "broker-service/internal/handler/http/product"
	httpUserHandler "broker-service/internal/handler/http/user"
	custommiddleware "broker-service/internal/middleware"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Rabbit             *amqp.Connection
	httpProductHandler *httpProductHandler.Config
	ProductHandler     *grpcProductHandler.ProductHandler
	CategoryHandler    *grpcProductHandler.CategoryGrpcHandler
	BannerHandler      *grpcProductHandler.BannerGrpcHandler
	AdsHandler         *grpcProductHandler.AdsGrpcHandler

	// Add auth and user handlers
	AuthHandler    *httpAuthHandler.Config
	UserHandler    *httpUserHandler.Config
	AuthMiddleware *custommiddleware.AuthMiddleware
}

const port string = "8080"

func main() {
	// try to connect to rabbitmq
	// rabbitConn, err := connect()
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	// defer rabbitConn.Close()

	// Initialize the gRPC clients
	productClient, err := grpcProductHandler.GetProductClient()
	if err != nil {
		log.Println("Error connecting to product service:", err)
		// Continue without product service functionality
	}

	// Connect to auth service
	authClient, err := grpcAuthHandler.GetAuthClient()
	if err != nil {
		log.Println("Error connecting to authentication service:", err)
		// Continue without auth service functionality
	}

	// Connect to user service
	userClient, err := grpcUserHandler.GetUserClient()
	if err != nil {
		log.Println("Error connecting to user service:", err)
		// Continue without user service functionality
	}

	// Initialize HTTP handlers
	authHttpHandler := &httpAuthHandler.Config{
		AuthClient: authClient,
	}

	userHttpHandler := &httpUserHandler.Config{
		UserClient: userClient,
	}

	// Initialize proxy handlers with productClient for image handling via gRPC
	httpProductHandler := &httpProductHandler.Config{
		ProductClient: productClient,
	}

	// Initialize auth middleware
	authMiddleware := &custommiddleware.AuthMiddleware{
		AuthClient: authClient,
	}

	// Initialize the application config
	app := Config{
		// Rabbit:      rabbitConn,
		httpProductHandler: httpProductHandler,
		ProductHandler:     grpcProductHandler.NewProductHandler(productClient),
		CategoryHandler:    grpcProductHandler.NewCategoryGrpcHandler(productClient),
		BannerHandler:      grpcProductHandler.NewBannerGrpcHandler(productClient),
		AdsHandler:         grpcProductHandler.NewAdsGrpcHandler(productClient),

		// Add auth and user handlers to config
		AuthHandler:    authHttpHandler,
		UserHandler:    userHttpHandler,
		AuthMiddleware: authMiddleware,
	}

	log.Printf("Starting broker service on port %s\n", port)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routers(),
	}

	// start http server
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}

}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

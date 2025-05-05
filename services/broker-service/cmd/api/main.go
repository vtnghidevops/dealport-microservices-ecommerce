package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	// communicate grpc between all service
	grpcAuthHandler "broker-service/internal/handlers/grpc/auth"
	grpcCartHandler "broker-service/internal/handlers/grpc/cart"
	grpcCheckoutHandler "broker-service/internal/handlers/grpc/checkout"
	grpcCouponHandler "broker-service/internal/handlers/grpc/coupon"
	grpcPaymentHandler "broker-service/internal/handlers/grpc/payment"
	grpcProductHandler "broker-service/internal/handlers/grpc/product"
	grpcUserHandler "broker-service/internal/handlers/grpc/user"

	// handle with http req from fe
	httpAuthHandler "broker-service/internal/handlers/http/auth"
	httpCartHandler "broker-service/internal/handlers/http/cart"
	httpCheckoutHandler "broker-service/internal/handlers/http/checkout"
	httpPaymentHandler "broker-service/internal/handlers/http/payment"
	httpProductHandler "broker-service/internal/handlers/http/product"
	httpUserHandler "broker-service/internal/handlers/http/user"
	custommiddleware "broker-service/internal/middleware"

	"broker-service/internal/event"

	"github.com/go-chi/chi/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Config is the application configuration
type Config struct {
	router             *chi.Mux
	eventEmitter       *event.Emitter
	Rabbit             *amqp.Connection
	httpProductHandler *httpProductHandler.Config
	ProductHandler     *grpcProductHandler.ProductHandler
	CategoryHandler    *grpcProductHandler.CategoryGrpcHandler
	BannerHandler      *grpcProductHandler.BannerGrpcHandler
	AdsHandler         *grpcProductHandler.AdsGrpcHandler

	// Auth and user handlers
	AuthHandler    *httpAuthHandler.Config
	UserHandler    *httpUserHandler.Config
	AuthMiddleware *custommiddleware.AuthMiddleware

	// Cart and checkout handlers
	CartHandler     *httpCartHandler.Config
	CheckoutHandler *httpCheckoutHandler.Config
	CouponHandler   *httpCartHandler.Config

	// Payment handler
	PaymentHandler *httpPaymentHandler.Config
}

func main() {
	// Create a new logger
	logger := log.New(os.Stdout, "[BROKER] ", log.LstdFlags)
	logger.Println("Starting broker service")

	// Get port from environment variables or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to RabbitMQ
	rabbitConn, err := connectToRabbitMQ()
	if err != nil {
		logger.Fatalf("Cannot connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	// Create an event emitter
	emitter, err := event.NewEventEmitter(rabbitConn)
	if err != nil {
		logger.Fatalf("Cannot create event emitter: %v", err)
	}

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

	// Connect to cart service
	cartClient, err := grpcCartHandler.GetCartClient()
	if err != nil {
		log.Println("Error connecting to cart service:", err)
		// Continue without cart service functionality
	}

	// Connect to checkout service
	checkoutClient, err := grpcCheckoutHandler.GetCheckoutClient()
	if err != nil {
		log.Println("Error connecting to checkout service:", err)
		// Continue without checkout service functionality
	}

	// Connect to payment service (uses checkout service underneath)
	paymentClient, err := grpcPaymentHandler.GetPaymentClient()
	if err != nil {
		log.Println("Error connecting to payment service:", err)
		// Continue without payment service functionality
	}

	// Connect to coupon service (which uses the same endpoint as cart service)
	couponClient, err := grpcCouponHandler.GetCouponClient()
	if err != nil {
		log.Println("Error connecting to coupon service:", err)
		// Continue without coupon service functionality
	}

	// Initialize HTTP handlers
	authHttpHandler := &httpAuthHandler.Config{
		AuthClient: authClient,
	}

	userHttpHandler := &httpUserHandler.Config{
		UserClient:     userClient,
		CheckoutClient: checkoutClient,
	}

	cartHttpHandler := &httpCartHandler.Config{
		CartClient: cartClient,
	}

	checkoutHttpHandler := &httpCheckoutHandler.Config{
		CheckoutClient: checkoutClient,
	}

	// Initialize proxy handlers with productClient for image handling via gRPC
	httpProductHandler := &httpProductHandler.Config{
		ProductClient: productClient,
	}

	// Initialize payment handler
	paymentHttpHandler := &httpPaymentHandler.Config{
		CheckoutClient: paymentClient,
	}

	// Initialize coupon handler
	couponHttpHandler := httpCartHandler.NewCouponHandler(couponClient)

	// Initialize auth middleware
	authMiddleware := &custommiddleware.AuthMiddleware{
		AuthClient: authClient,
	}

	// Create the application config
	app := Config{
		router:             chi.NewRouter(),
		eventEmitter:       emitter,
		Rabbit:             rabbitConn,
		httpProductHandler: httpProductHandler,
		ProductHandler:     grpcProductHandler.NewProductHandler(productClient),
		CategoryHandler:    grpcProductHandler.NewCategoryGrpcHandler(productClient),
		BannerHandler:      grpcProductHandler.NewBannerGrpcHandler(productClient),
		AdsHandler:         grpcProductHandler.NewAdsGrpcHandler(productClient),

		// Add auth and user handlers to config
		AuthHandler:    authHttpHandler,
		UserHandler:    userHttpHandler,
		AuthMiddleware: authMiddleware,

		// Add cart, checkout and coupon handlers to config
		CartHandler:     cartHttpHandler,
		CheckoutHandler: checkoutHttpHandler,
		CouponHandler:   couponHttpHandler,

		// Add payment handler to config
		PaymentHandler: paymentHttpHandler,
	}

	// Set up the routes
	app.routers()

	// Start the server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      app.routers(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Printf("Starting broker service on port %s\n", port)

	err = srv.ListenAndServe()
	if err != nil {
		logger.Println(err)
	}
}

func connectToRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost:5672")
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

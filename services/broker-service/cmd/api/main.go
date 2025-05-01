package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	// communicate grpc between all service
	grpcAuthHandler "broker-service/internal/handler/grpc/auth"
	grpcCartHandler "broker-service/internal/handler/grpc/cart"
	grpcCheckoutHandler "broker-service/internal/handler/grpc/checkout"
	grpcCouponHandler "broker-service/internal/handler/grpc/coupon"
	grpcPaymentHandler "broker-service/internal/handler/grpc/payment"
	grpcProductHandler "broker-service/internal/handler/grpc/product"
	grpcUserHandler "broker-service/internal/handler/grpc/user"

	// handle with http req from fe
	httpAuthHandler "broker-service/internal/handler/http/auth"
	httpCartHandler "broker-service/internal/handler/http/cart"
	httpCheckoutHandler "broker-service/internal/handler/http/checkout"
	httpPaymentHandler "broker-service/internal/handler/http/payment"
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
		UserClient: userClient,
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

		// Add cart, checkout and coupon handlers to config
		CartHandler:     cartHttpHandler,
		CheckoutHandler: checkoutHttpHandler,
		CouponHandler:   couponHttpHandler,

		// Add payment handler to config
		PaymentHandler: paymentHttpHandler,
	}

	log.Printf("Starting broker service on port %s\n", port)
	log.Println("Connected services:")
	if productClient != nil {
		log.Println("- Product service: Connected")
	}
	if authClient != nil {
		log.Println("- Auth service: Connected")
	}
	if userClient != nil {
		log.Println("- User service: Connected")
	}
	if cartClient != nil {
		log.Println("- Cart service: Connected")
	}
	if checkoutClient != nil {
		log.Println("- Checkout service: Connected")
	}
	if paymentClient != nil {
		log.Println("- Payment service: Connected")
	}

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

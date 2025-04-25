package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	grpcProductHandler "broker-service/internal/handler/grpc/product"
	httpProductHandler "broker-service/internal/handler/http/product"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Rabbit             *amqp.Connection
	httpProductHandler *httpProductHandler.Config
	ProductHandler     *grpcProductHandler.ProductHandler
	CategoryHandler    *grpcProductHandler.CategoryGrpcHandler
	BannerHandler      *grpcProductHandler.BannerGrpcHandler
	AdsHandler         *grpcProductHandler.AdsGrpcHandler
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

	// Initialize the gRPC client
	productClient, err := grpcProductHandler.GetProductClient()
	if err != nil {
		log.Fatal("Error connecting to product service: ", err)
	}

	// Initialize proxy handlers with productClient for image handling via gRPC
	httpProductHandler := &httpProductHandler.Config{
		ProductClient: productClient,
	}

	// Initialize the gRPC handlers for product service
	app := Config{
		// Rabbit:      rabbitConn,
		httpProductHandler: httpProductHandler,
		ProductHandler:     grpcProductHandler.NewProductHandler(productClient),
		CategoryHandler:    grpcProductHandler.NewCategoryGrpcHandler(productClient),
		BannerHandler:      grpcProductHandler.NewBannerGrpcHandler(productClient),
		AdsHandler:         grpcProductHandler.NewAdsGrpcHandler(productClient),
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

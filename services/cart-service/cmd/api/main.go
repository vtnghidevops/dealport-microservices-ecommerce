package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"cart-service/internal/config"
	"cart-service/internal/domain"
	"cart-service/internal/repository/redis"
	"cart-service/internal/service"
	transportGrpc "cart-service/internal/transport/grpc"
	pb "cart-service/proto/cart"
	couponpb "cart-service/proto/coupon"

	redisClient "github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Process command line arguments
	if len(os.Args) > 1 {
		handleCommands()
		return
	}

	log.Println("Starting cart service")

	// Load config
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to Redis
	redisConn := connectToRedis(cfg.Redis)
	defer redisConn.Close()

	// Auto-initialize Redis if enabled
	if os.Getenv("AUTO_REDIS_INIT") == "true" {
		log.Println("AUTO_REDIS_INIT enabled, initializing Redis...")
		if err := initRedis(); err != nil {
			log.Printf("Warning: Failed to initialize Redis: %v", err)
		}
	}

	// Create repositories
	cartRepo := redis.NewCartRepository(redisConn)
	couponRepo := redis.NewCouponRepository(redisConn)

	// Create services
	couponService := service.NewCouponService(couponRepo)
	cartService := service.NewCartService(cartRepo, couponService)

	// Start gRPC server
	go startGRPCServer(cartService, couponService, cfg.Server.GRPCPort)

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down cart service...")
	time.Sleep(1 * time.Second) // Allow some time for cleanup
}

// handleCommands processes command line arguments
func handleCommands() {
	command := os.Args[1]

	switch command {
	case "init-redis":
		// Initialize Redis
		if err := initRedis(); err != nil {
			log.Fatalf("Failed to initialize Redis: %v", err)
		}
		log.Println("Redis initialization completed successfully")
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

// initRedis runs the Redis initialization script
func initRedis() error {
	// Find the script path - should be in the migrations directory
	scriptPath := "migrations/init_redis.sh"

	// For container environment, we may need to look elsewhere
	containerScriptPath := "/app/migrations/init_redis.sh"

	// Check if the script exists in the first path
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		// If not, try the container path
		if _, err := os.Stat(containerScriptPath); os.IsNotExist(err) {
			return fmt.Errorf("Redis initialization script not found at %s or %s",
				scriptPath, containerScriptPath)
		}
		scriptPath = containerScriptPath
	}

	// Make the script executable
	if err := os.Chmod(scriptPath, 0755); err != nil {
		return fmt.Errorf("failed to make initialization script executable: %w", err)
	}

	// Get the directory containing the script
	scriptDir := filepath.Dir(scriptPath)

	// Set environment for the script
	env := os.Environ()
	if os.Getenv("ENVIRONMENT") == "" {
		env = append(env, "ENVIRONMENT="+os.Getenv("ENVIRONMENT"))
	}

	// Create and run the command
	cmd := exec.Command(scriptPath)
	cmd.Env = env
	cmd.Dir = scriptDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Running Redis initialization script: %s", scriptPath)
	return cmd.Run()
}

func startGRPCServer(cartService domain.CartService, couponService domain.CouponService, port string) {
	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	cartGrpcServer := transportGrpc.NewServer(cartService)
	pb.RegisterCartServiceServer(grpcServer, cartGrpcServer)

	// Register coupon service
	couponGrpcServer := transportGrpc.NewCouponServer(couponService)
	couponpb.RegisterCouponServiceServer(grpcServer, couponGrpcServer)

	// Register reflection service for gRPC tools and clients
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on port %s", port)

	// Start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func connectToRedis(cfg config.RedisConfig) *redisClient.Client {
	// Create Redis client
	client := redisClient.NewClient(&redisClient.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping Redis
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully!")
	return client
}

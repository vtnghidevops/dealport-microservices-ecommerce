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
	// First try the script approach
	if err := initRedisWithScript(); err != nil {
		log.Printf("Script initialization failed: %v", err)
		log.Println("Fallback to direct Redis commands...")
		return initRedisDirectly()
	}
	return nil
}

// initRedisWithScript uses the shell script
func initRedisWithScript() error {
	// Debug: Print current working directory
	wd, _ := os.Getwd()
	log.Printf("Current working directory: %s", wd)

	// Use absolute path for container environment - more reliable
	scriptPath := "/app/migrations/init_redis.sh"

	// Fallback to relative path if absolute doesn't exist
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		log.Printf("Script not found at absolute path: %s", scriptPath)
		scriptPath = "migrations/init_redis.sh"
		log.Printf("Trying relative path: %s", scriptPath)
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			log.Printf("Script not found at relative path: %s", scriptPath)
			return fmt.Errorf("Redis initialization script not found at %s or migrations/init_redis.sh", "/app/migrations/init_redis.sh")
		}
		log.Printf("Found script at relative path: %s", scriptPath)
	} else {
		log.Printf("Found script at absolute path: %s", scriptPath)
	}

	// Debug: Check current file permissions
	if info, err := os.Stat(scriptPath); err == nil {
		log.Printf("Script file info: %s, size: %d, mode: %s", scriptPath, info.Size(), info.Mode())
	}

	// Make the script executable
	log.Printf("Setting executable permissions for: %s", scriptPath)
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

	// Create and run the command - use sh to ensure compatibility
	cmd := exec.Command("sh", scriptPath)
	cmd.Env = env
	cmd.Dir = scriptDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Running Redis initialization script: sh %s", scriptPath)
	log.Printf("Script directory: %s", scriptDir)
	log.Printf("Command: %s", cmd.String())

	// Run and capture detailed error
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run Redis initialization script: %w", err)
	}

	log.Println("Redis initialization script completed successfully")
	return nil
}

// initRedisDirectly executes Redis commands directly via Go
func initRedisDirectly() error {
	log.Println("Initializing Redis directly from Go...")

	// Load config to get Redis connection details
	cfg, err := config.LoadConfig("")
	if err != nil {
		return fmt.Errorf("failed to load config for Redis init: %w", err)
	}

	// Create Redis client for initialization
	client := redisClient.NewClient(&redisClient.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Wait for Redis to be ready
	log.Println("Waiting for Redis to be ready...")
	for i := 0; i < 30; i++ {
		if err := client.Ping(ctx).Err(); err == nil {
			break
		}
		log.Printf("Redis not ready yet, waiting... (attempt %d/30)", i+1)
		time.Sleep(1 * time.Second)
	}

	// Test final connection
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis not ready after 30 seconds: %w", err)
	}

	log.Println("Redis is ready, running initialization commands...")

	// Run initialization commands
	pipe := client.Pipeline()
	pipe.FlushAll(ctx)
	pipe.Set(ctx, "init:status", "Initialization completed successfully", 24*time.Hour)
	pipe.Set(ctx, "init:timestamp", time.Now().Unix(), 24*time.Hour)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute Redis initialization commands: %w", err)
	}

	log.Println("Redis initialization completed successfully via Go")
	return nil
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

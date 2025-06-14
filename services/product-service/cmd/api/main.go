// cmd/api/main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"product-service/internal/cache"
	"product-service/internal/config"
	"product-service/internal/domain"
	"product-service/internal/migrations"
	"product-service/internal/repository/postgres"
	"product-service/internal/service"
	"product-service/internal/storage"
	transportGrpc "product-service/internal/transport/grpc"
	pb "product-service/proto/product"
	"strconv"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

func main() {
	// Process command line arguments for migrations
	if len(os.Args) > 1 {
		handleMigrationCommands()
		return
	}

	// Load config (Load configuration)
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Debug log about storage configuration
	log.Printf("DEBUG: Storage configuration - Provider: '%s'", cfg.Storage.Provider)
	if cfg.Storage.Provider == "minio" {
		log.Printf("DEBUG: MinIO config - Endpoint: %s, AccessKey: %s, BucketName: %s",
			cfg.Storage.MinIO.Endpoint,
			cfg.Storage.MinIO.AccessKeyID,
			cfg.Storage.MinIO.BucketName)
	} else {
		log.Printf("DEBUG: Using local file storage because Storage.Provider is not 'minio'")
	}

	log.Printf("Starting product service - gRPC port: %s, HTTP port: %s", cfg.Server.GRPCPort, cfg.Server.HTTPPort)

	// Connect to database
	conn := connectToDB(cfg.Database)
	if conn == nil {
		log.Panic("Could not connect to Postgres")
	}
	defer conn.Close()

	// Convert sql.DB to sqlx.DB for migrations
	dbx := sqlx.NewDb(conn, "pgx")

	// Run migrations if AUTO_MIGRATE is set to true
	if os.Getenv("AUTO_MIGRATE") == "true" {
		log.Println("AUTO_MIGRATE enabled, running migrations...")
		if err := migrations.RunMigrations(dbx); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	}

	// Create repositories
	productRepo := postgres.NewProductRepository(conn)
	categoryRepo := postgres.NewCategoryRepository(conn)
	bannerRepo := postgres.NewBannerRepository(conn)
	adsRepo := postgres.NewAdsRepository(conn)

	// Initialize storage service based on configuration
	var storageService storage.StorageService
	var urlCache cache.ImageURLCache

	if cfg.Storage.Provider == "minio" {
		log.Printf("Using MinIO storage with endpoint: %s", cfg.Storage.MinIO.Endpoint)

		// Initialize MinIO storage
		minioConfig := storage.MinioConfig{
			Endpoint:        cfg.Storage.MinIO.Endpoint,
			AccessKeyID:     cfg.Storage.MinIO.AccessKeyID,
			SecretAccessKey: cfg.Storage.MinIO.SecretAccessKey,
			UseSSL:          cfg.Storage.MinIO.UseSSL,
			BucketName:      cfg.Storage.MinIO.BucketName,
			Location:        cfg.Storage.MinIO.Location,
			BaseURL:         cfg.Storage.MinIO.BaseURL,
			PresignedTTL:    cfg.Storage.MinIO.PresignedTTL,
		}

		log.Printf("Attempting to connect to MinIO at %s", cfg.Storage.MinIO.Endpoint)
		log.Printf("Bucket: %s, UseSSL: %v", cfg.Storage.MinIO.BucketName, cfg.Storage.MinIO.UseSSL)
		log.Printf("BaseURL (for presigned URLs): %s", cfg.Storage.MinIO.BaseURL)

		// Create MinIO storage service
		minioStorage, err := storage.NewMinioStorage(minioConfig)
		if err != nil {
			log.Printf("ERROR: Failed to initialize MinIO client: %v", err)
			log.Printf("DEBUG: MinIO connection details - Endpoint: %s, AccessKey: %s, BucketName: %s",
				cfg.Storage.MinIO.Endpoint,
				cfg.Storage.MinIO.AccessKeyID,
				cfg.Storage.MinIO.BucketName)
			log.Printf("DEBUG: Check that MinIO service is running and accessible from this container/machine")
			log.Printf("DEBUG: Also verify network connectivity and firewall settings")
			log.Printf("Falling back to local file storage")
		} else {
			// Test connection to MinIO
			ctx := context.Background()
			err = minioStorage.TestConnection(ctx)
			if err != nil {
				log.Printf("ERROR: Failed to connect to MinIO: %v", err)
				log.Printf("DEBUG: MinIO connection test failed - details:")
				log.Printf("DEBUG: Endpoint: %s, UseSSL: %v", cfg.Storage.MinIO.Endpoint, cfg.Storage.MinIO.UseSSL)
				log.Printf("DEBUG: AccessKey: %s, BucketName: %s", cfg.Storage.MinIO.AccessKeyID, cfg.Storage.MinIO.BucketName)
				log.Printf("DEBUG: Network connectivity checks to try:")
				log.Printf("DEBUG: 1. Can you ping %s?", cfg.Storage.MinIO.Endpoint)
				log.Printf("DEBUG: 2. Can you access MinIO directly in browser: %s?", cfg.Storage.MinIO.BaseURL)
				log.Printf("DEBUG: 3. Check DNS resolution for %s", cfg.Storage.MinIO.Endpoint)
				log.Printf("DEBUG: 4. Check if your MinIO server is correctly configured and running")
				log.Printf("DEBUG: 5. Verify MinIO credentials are correct")

				// Check if bucket exists
				exists, bucketErr := minioStorage.BucketExists(ctx)
				if bucketErr != nil {
					log.Printf("ERROR: Failed to check if bucket exists: %v", bucketErr)
					log.Printf("DEBUG: This usually means connectivity issues to MinIO server")
				} else if !exists {
					log.Printf("DEBUG: MinIO server is accessible but bucket '%s' does not exist", cfg.Storage.MinIO.BucketName)
					log.Printf("Bucket '%s' does not exist, attempting to create it...", cfg.Storage.MinIO.BucketName)

					// Try to create the bucket
					createErr := minioStorage.CreateBucket(ctx, cfg.Storage.MinIO.BucketName, cfg.Storage.MinIO.Location)
					if createErr != nil {
						log.Printf("ERROR: Failed to create bucket: %v", createErr)
						log.Printf("DEBUG: Check MinIO user permissions - user needs CreateBucket rights")
					} else {
						log.Printf("Successfully created bucket '%s'", cfg.Storage.MinIO.BucketName)
					}
				}

				log.Printf("WARNING: Using MinIO storage but connection test failed. Service may fall back to local storage for uploads.")
			} else {
				log.Printf("Successfully connected to MinIO and verified bucket '%s' exists", cfg.Storage.MinIO.BucketName)
				log.Printf("MinIO storage is ready for use")
			}

			// Create URL cache for presigned URLs
			urlCache = cache.NewImageURLCache(30 * time.Minute)
			storageService = minioStorage
		}
	} else {
		log.Printf("Using local file storage because Storage.Provider is not 'minio'")
	}

	// Create services with storage
	var productService domain.ProductService
	if storageService != nil {
		productService = service.NewProductServiceWithStorage(productRepo, storageService, urlCache)
		log.Printf("Product service initialized with MinIO storage")
	} else {
		productService = service.NewProductService(productRepo)
		log.Printf("Product service initialized with local storage")
	}

	categoryService := service.NewCategoryService(categoryRepo)
	bannerService := service.NewBannerService(bannerRepo)
	adsService := service.NewAdsService(adsRepo)

	// Create a channel to catch errors
	errCh := make(chan error, 2)

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("Starting gRPC server on port %s...", cfg.Server.GRPCPort)
		errCh <- startGRPCServer(productService, categoryService, bannerService, adsService, cfg.Server.GRPCPort)
	}()

	// Block until we get an error from one of the servers
	log.Fatalf("Server error: %v", <-errCh)
}

// handleMigrationCommands processes the migration related commands
func handleMigrationCommands() {
	migrationCommand := os.Args[1]

	// Load config
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	conn := connectToDB(cfg.Database)
	if conn == nil {
		log.Fatalf("Could not connect to Postgres for migrations")
	}
	defer conn.Close()

	// Convert sql.DB to sqlx.DB for migrations
	dbx := sqlx.NewDb(conn, "pgx")

	// Process the command
	switch migrationCommand {
	case "migrate":
		// Run migrations up
		if err := migrations.RunMigrations(dbx); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "rollback":
		// Rollback the last migration
		if err := migrations.RollbackMigration(dbx); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
	case "drop":
		// Drop all tables (development only)
		if err := migrations.DropAllTables(dbx); err != nil {
			log.Fatalf("Drop failed: %v", err)
		}
	case "status":
		// Print current migration status
		version, dirty, err := migrations.GetMigrationStatus(dbx)
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
		log.Printf("Current migration version: %d, dirty: %t", version, dirty)
	case "force":
		// Force migration to specific version
		if len(os.Args) < 3 {
			log.Fatalf("force command requires a version number")
		}

		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}

		if err := migrations.ForceVersion(dbx, version); err != nil {
			log.Fatalf("Force failed: %v", err)
		}
	default:
		log.Fatalf("Unknown migration command: %s", migrationCommand)
	}

	log.Printf("Migration command '%s' completed successfully", migrationCommand)
}

func startGRPCServer(
	productService domain.ProductService,
	categoryService domain.CategoryService,
	bannerService domain.BannerService,
	adsService domain.AdsService,
	port string,
) error {
	// Create TCP listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen for gRPC: %w", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Create and register product service server
	productGrpcServer := transportGrpc.NewGrpcServer(productService)
	categoryGrpcServer := transportGrpc.NewCategoryGrpcServer(categoryService)
	bannerGrpcServer := transportGrpc.NewBannerGrpcServer(bannerService)
	adsGrpcServer := transportGrpc.NewAdsGrpcServer(adsService)

	// Log before registration
	log.Printf("Registering services with gRPC server...")

	// Register the services with the gRPC server
	pb.RegisterProductServiceServer(grpcServer, productGrpcServer)
	pb.RegisterCategoryServiceServer(grpcServer, categoryGrpcServer)
	pb.RegisterBannerServiceServer(grpcServer, bannerGrpcServer)
	pb.RegisterAdsServiceServer(grpcServer, adsGrpcServer)

	log.Printf("gRPC server listening on port %s", port)

	// Start gRPC server (this blocks)
	return grpcServer.Serve(lis)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB(dbCfg config.DatabaseConfig) *sql.DB {
	//dsn := os.Getenv("DSN")
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DBName, dbCfg.SSLMode,
	)

	// if dsn == "" {
	// 	dsn = "host=postgres-products port=5432 user=postgres-products password=password dbname=products sslmode=disable timezone=UTC connect_timeout=5"
	// 	//dsn = "host=localhost port=5433 user=postgres-products password=password dbname=products sslmode=disable timezone=UTC connect_timeout=5"
	// }

	// Try to connect to the database with retries
	for i := 0; i < 10; i++ {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("PostgreSQL not ready yet...")
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Connected to PostgreSQL!")
		return connection
	}

	log.Println("Could not connect to PostgreSQL after multiple attempts")
	return nil
}

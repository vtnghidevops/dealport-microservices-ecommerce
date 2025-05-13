// cmd/api/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"product-service/internal/config"
	"product-service/internal/domain"
	"product-service/internal/handler"
	"product-service/internal/migrations"
	"product-service/internal/repository/postgres"
	"product-service/internal/service"
	transportGrpc "product-service/internal/transport/grpc"
	transportHttp "product-service/internal/transport/http"
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

	// Load config (Tải cấu hình)
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
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

	// Create services
	productService := service.NewProductService(productRepo)
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

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on port %s...", cfg.Server.HTTPPort)

		// Create handler config
		handlerConfig := &handler.Config{
			ProductService:  productService,
			CategoryService: categoryService,
			BannerService:   bannerService,
			AdsService:      adsService,
		}

		// Create HTTP server
		httpServer := transportHttp.NewServer(handlerConfig)

		// Start HTTP server
		errCh <- http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.HTTPPort), httpServer.Routes())
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

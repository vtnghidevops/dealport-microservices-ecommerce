package migrations

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabaseMigrator handles MongoDB migrations
type MongoDatabaseMigrator struct {
	client *mongo.Client
	dbName string
	logger *log.Logger
}

// NewMongoDatabaseMigrator creates a new MongoDB migrator
func NewMongoDatabaseMigrator(client *mongo.Client, dbName string) *MongoDatabaseMigrator {
	return &MongoDatabaseMigrator{
		client: client,
		dbName: dbName,
		logger: log.New(os.Stdout, "[MONGO-SETUP] ", log.LstdFlags),
	}
}

// RunMigrations sets up MongoDB collections and indexes
func (m *MongoDatabaseMigrator) RunMigrations() error {
	m.logger.Println("Starting MongoDB setup...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := m.client.Database(m.dbName)

	// Check existing collections
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	// Check if orders collection exists
	ordersExists := false
	for _, name := range collections {
		if name == "orders" {
			ordersExists = true
			break
		}
	}

	// Create orders collection if it doesn't exist
	if !ordersExists {
		m.logger.Printf("Creating orders collection")
		if err := db.CreateCollection(ctx, "orders"); err != nil {
			return fmt.Errorf("failed to create orders collection: %w", err)
		}
	} else {
		m.logger.Printf("Orders collection already exists")
	}

	// Check if payments collection exists
	paymentsExists := false
	for _, name := range collections {
		if name == "payments" {
			paymentsExists = true
			break
		}
	}

	// Create payments collection if it doesn't exist
	if !paymentsExists {
		m.logger.Printf("Creating payments collection")
		if err := db.CreateCollection(ctx, "payments"); err != nil {
			return fmt.Errorf("failed to create payments collection: %w", err)
		}
	} else {
		m.logger.Printf("Payments collection already exists")
	}

	// Create indexes for orders collection
	ordersColl := db.Collection("orders")

	// User ID index
	m.logger.Printf("Creating user_id index on orders collection")
	userIDIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	}
	if _, err := ordersColl.Indexes().CreateOne(ctx, userIDIndex); err != nil {
		return fmt.Errorf("failed to create user_id index: %w", err)
	}

	// Order number index (unique)
	orderNumberIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "order_number", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := ordersColl.Indexes().CreateOne(ctx, orderNumberIndex); err != nil {
		return fmt.Errorf("failed to create order_number unique index: %w", err)
	}

	// Created at index (descending)
	createdAtIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: -1}},
	}
	if _, err := ordersColl.Indexes().CreateOne(ctx, createdAtIndex); err != nil {
		return fmt.Errorf("failed to create created_at index: %w", err)
	}

	// Status index
	statusIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}},
	}
	if _, err := ordersColl.Indexes().CreateOne(ctx, statusIndex); err != nil {
		return fmt.Errorf("failed to create status index: %w", err)
	}

	// Create indexes for payments collection
	paymentsColl := db.Collection("payments")

	// Order ID index
	orderIDIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "order_id", Value: 1}},
	}
	if _, err := paymentsColl.Indexes().CreateOne(ctx, orderIDIndex); err != nil {
		return fmt.Errorf("failed to create order_id index on payments: %w", err)
	}

	// Transaction ID index (unique, sparse)
	txnIDIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "transaction_id", Value: 1}},
		Options: options.Index().SetUnique(true).SetSparse(true),
	}
	if _, err := paymentsColl.Indexes().CreateOne(ctx, txnIDIndex); err != nil {
		return fmt.Errorf("failed to create transaction_id unique index: %w", err)
	}

	m.logger.Printf("Successfully created collections and indexes for checkout service")
	return nil
}

// RollbackMigration is a no-op in the simplified approach
func (m *MongoDatabaseMigrator) RollbackMigration(migrationName string) error {
	m.logger.Printf("Rollback functionality not needed with simplified approach")
	return nil
}

// GetMigrationStatus gets the status of collections
func (m *MongoDatabaseMigrator) GetMigrationStatus() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := m.client.Database(m.dbName)

	// Check if the collections exist
	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return "", fmt.Errorf("failed to list collections: %w", err)
	}

	// Count documents in collections
	ordersCount := int64(0)
	paymentsCount := int64(0)

	for _, name := range collections {
		if name == "orders" {
			ordersCount, err = db.Collection("orders").CountDocuments(ctx, bson.M{})
			if err != nil {
				m.logger.Printf("Error counting orders: %v", err)
			}
		}
		if name == "payments" {
			paymentsCount, err = db.Collection("payments").CountDocuments(ctx, bson.M{})
			if err != nil {
				m.logger.Printf("Error counting payments: %v", err)
			}
		}
	}

	return fmt.Sprintf("Checkout database has %d collections: orders(%d), payments(%d)",
		len(collections), ordersCount, paymentsCount), nil
}

// CheckDatabaseConnection checks if the MongoDB connection is working
func (m *MongoDatabaseMigrator) CheckDatabaseConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.client.Ping(ctx, nil)
}

// ConnectToMongoDB establishes a connection to MongoDB
func ConnectToMongoDB(mongoURI string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

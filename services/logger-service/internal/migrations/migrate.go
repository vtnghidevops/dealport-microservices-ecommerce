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
		logger: log.New(os.Stdout, "[MONGO-MIGRATE] ", log.LstdFlags),
	}
}

// RunMigrations runs all migrations from the JS files in the migrations directory
func (m *MongoDatabaseMigrator) RunMigrations() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	m.logger.Println("Starting MongoDB migrations...")

	// Get the current database
	db := m.client.Database(m.dbName)

	// Simplified approach: skip migrations tracking and just create logs collection directly
	m.logger.Printf("Setting up logs collection...")

	// Check if logs collection exists
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	// Check if logs collection exists
	logsExists := false
	for _, name := range collections {
		if name == "logs" {
			logsExists = true
			break
		}
	}

	// Create logs collection if it doesn't exist
	if !logsExists {
		m.logger.Printf("Creating 'logs' collection...")
		if err := db.CreateCollection(ctx, "logs"); err != nil {
			return fmt.Errorf("failed to create logs collection: %w", err)
		}
		m.logger.Printf("Collection 'logs' created successfully")
	} else {
		m.logger.Printf("Collection 'logs' already exists")
	}

	// Ensure indexes on logs collection
	m.logger.Printf("Creating indexes on logs collection...")
	indexOpts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	logsColl := db.Collection("logs")

	// Create timestamp index
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "timestamp", Value: -1},
		},
	}

	if _, err := logsColl.Indexes().CreateOne(ctx, indexModel, indexOpts); err != nil {
		return fmt.Errorf("failed to create timestamp index on logs collection: %w", err)
	}

	// Create service index
	serviceIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "service", Value: 1},
		},
	}

	if _, err := logsColl.Indexes().CreateOne(ctx, serviceIndex, indexOpts); err != nil {
		return fmt.Errorf("failed to create service index on logs collection: %w", err)
	}

	m.logger.Printf("Successfully created indexes on logs collection")
	m.logger.Println("MongoDB setup completed successfully")

	return nil
}

// RollbackMigration is no longer needed as we're not tracking migrations
func (m *MongoDatabaseMigrator) RollbackMigration() error {
	m.logger.Println("Rollback not needed with simplified approach")
	return nil
}

// DropDatabase drops the entire database
func (m *MongoDatabaseMigrator) DropDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	m.logger.Printf("Dropping database: %s", m.dbName)

	// Drop the database
	err := m.client.Database(m.dbName).Drop(ctx)
	if err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	m.logger.Println("Database dropped successfully")
	return nil
}

// GetMigrationStatus returns the status of the logs collection
func (m *MongoDatabaseMigrator) GetMigrationStatus() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get the current database
	db := m.client.Database(m.dbName)

	// Check if logs collection exists
	colNames, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return "", fmt.Errorf("failed to list collections: %w", err)
	}

	logsExists := false
	for _, name := range colNames {
		if name == "logs" {
			logsExists = true
			break
		}
	}

	if !logsExists {
		return "Logs collection has not been set up", nil
	}

	// Count logs
	count, err := db.Collection("logs").CountDocuments(ctx, bson.M{})
	if err != nil {
		return "", fmt.Errorf("failed to count logs: %w", err)
	}

	return fmt.Sprintf("Logs collection is set up and contains %d entries", count), nil
}

// CheckDatabaseConnection checks if the database connection is working
func (m *MongoDatabaseMigrator) CheckDatabaseConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	m.logger.Println("Checking database connection...")

	// Check connection by listing databases
	_, err := m.client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	m.logger.Println("Database connection is working")
	return nil
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

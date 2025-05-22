package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define the interfaces for MongoDB operations
// These interfaces allow for easier mocking in tests

// MongoCursorInterface defines the interface for a MongoDB cursor
type MongoCursorInterface interface {
	Next(ctx context.Context) bool
	Decode(val interface{}) error
	Close(ctx context.Context) error
}

// MongoSingleResultInterface defines the interface for a MongoDB single result
type MongoSingleResultInterface interface {
	Decode(val interface{}) error
}

// MongoCollectionInterface defines the interface for a MongoDB collection
type MongoCollectionInterface interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (MongoCursorInterface, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) MongoSingleResultInterface
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Drop(ctx context.Context) error
}

// MongoDatabaseInterface defines the interface for a MongoDB database
type MongoDatabaseInterface interface {
	Collection(name string, opts ...*options.CollectionOptions) MongoCollectionInterface
}

// MongoClientInterface defines the interface for a MongoDB client
type MongoClientInterface interface {
	Database(name string, opts ...*options.DatabaseOptions) MongoDatabaseInterface
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
}

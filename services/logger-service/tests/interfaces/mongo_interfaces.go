package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCursor định nghĩa interface cho MongoDB cursor
type MongoCursor interface {
	Next(context.Context) bool
	Decode(interface{}) error
	Close(context.Context) error
	Err() error
	All(context.Context, interface{}) error
	ID() int64
	RemainingBatchLength() int
	TryNext(context.Context) bool
}

// MongoSingleResult định nghĩa interface cho MongoDB single result
type MongoSingleResult interface {
	Decode(interface{}) error
	Err() error
}

// MongoDatabase định nghĩa interface cho MongoDB database
type MongoDatabase interface {
	Collection(string, ...*options.CollectionOptions) MongoCollection
}

// MongoCollection định nghĩa interface cho MongoDB collection
type MongoCollection interface {
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	Find(context.Context, interface{}, ...*options.FindOptions) (MongoCursor, error)
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) MongoSingleResult
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	Drop(context.Context) error
}

// MongoClient định nghĩa interface cho MongoDB client
type MongoClient interface {
	Database(string, ...*options.DatabaseOptions) MongoDatabase
}

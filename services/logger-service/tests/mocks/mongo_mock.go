package mocks

import (
	"context"
	"logger-service/tests/interfaces"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockMongoClient mô phỏng interfaces.MongoClientInterface
type MockMongoClient struct {
	mock.Mock
}

// Database mô phỏng phương thức Database của interfaces.MongoClientInterface
func (m *MockMongoClient) Database(name string, opts ...*options.DatabaseOptions) interfaces.MongoDatabaseInterface {
	args := m.Called(name, opts)
	return args.Get(0).(interfaces.MongoDatabaseInterface)
}

// Connect mô phỏng phương thức Connect của interfaces.MongoClientInterface
func (m *MockMongoClient) Connect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Disconnect mô phỏng phương thức Disconnect của interfaces.MongoClientInterface
func (m *MockMongoClient) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockMongoDatabase mô phỏng interfaces.MongoDatabaseInterface
type MockMongoDatabase struct {
	mock.Mock
}

// Collection mô phỏng phương thức Collection của interfaces.MongoDatabaseInterface
func (m *MockMongoDatabase) Collection(name string, opts ...*options.CollectionOptions) interfaces.MongoCollectionInterface {
	args := m.Called(name, opts)
	return args.Get(0).(interfaces.MongoCollectionInterface)
}

// MockMongoCollection mô phỏng interfaces.MongoCollectionInterface
type MockMongoCollection struct {
	mock.Mock
}

// InsertOne mô phỏng phương thức InsertOne của interfaces.MongoCollectionInterface
func (m *MockMongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

// Find mô phỏng phương thức Find của interfaces.MongoCollectionInterface
func (m *MockMongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (interfaces.MongoCursorInterface, error) {
	args := m.Called(ctx, filter, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(interfaces.MongoCursorInterface), args.Error(1)
}

// FindOne mô phỏng phương thức FindOne của interfaces.MongoCollectionInterface
func (m *MockMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) interfaces.MongoSingleResultInterface {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(interfaces.MongoSingleResultInterface)
}

// UpdateOne mô phỏng phương thức UpdateOne của interfaces.MongoCollectionInterface
func (m *MockMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

// DeleteOne mô phỏng phương thức DeleteOne của interfaces.MongoCollectionInterface
func (m *MockMongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

// Drop mô phỏng phương thức Drop của interfaces.MongoCollectionInterface
func (m *MockMongoCollection) Drop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockMongoCursor mô phỏng interfaces.MongoCursorInterface
type MockMongoCursor struct {
	mock.Mock
}

// Next mô phỏng phương thức Next của interfaces.MongoCursorInterface
func (m *MockMongoCursor) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

// Decode mô phỏng phương thức Decode của interfaces.MongoCursorInterface
func (m *MockMongoCursor) Decode(val interface{}) error {
	args := m.Called(val)
	return args.Error(0)
}

// Close mô phỏng phương thức Close của interfaces.MongoCursorInterface
func (m *MockMongoCursor) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockMongoSingleResult mô phỏng interfaces.MongoSingleResultInterface
type MockMongoSingleResult struct {
	mock.Mock
}

// Decode mô phỏng phương thức Decode của interfaces.MongoSingleResultInterface
func (m *MockMongoSingleResult) Decode(val interface{}) error {
	args := m.Called(val)
	return args.Error(0)
}

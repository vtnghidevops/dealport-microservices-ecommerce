package helpers

import (
	"errors"
	"logger-service/internal/domain"
	"logger-service/tests/mocks"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PatchMongoCollection is used to get the mock collection in tests
var PatchMongoCollection func() interface{}

// SetupMockCollection cài đặt các mock cho LogEntry operations trên một mock collection
func SetupMockCollection(collection *mocks.MockMongoCollection, mockResults map[string]interface{}) {
	// Mock cho Insert operation
	insertResult, hasInsertResult := mockResults["insertResult"].(*mongo.InsertOneResult)
	insertError, hasInsertError := mockResults["insertError"].(error)

	if hasInsertResult || hasInsertError {
		expectation := collection.On("InsertOne", mock.Anything, mock.AnythingOfType("domain.LogEntry"), mock.Anything)
		if hasInsertError {
			expectation.Return(nil, insertError)
		} else {
			expectation.Return(insertResult, nil)
		}
	}

	// Mock cho All operation (Find)
	findCursor, hasFindCursor := mockResults["findCursor"].(mocks.MockMongoCursor)
	findError, hasFindError := mockResults["findError"].(error)

	if hasFindCursor || hasFindError {
		expectation := collection.On("Find", mock.Anything, mock.Anything, mock.AnythingOfType("*options.FindOptions"))
		if hasFindError {
			expectation.Return(nil, findError)
		} else {
			expectation.Return(&findCursor, nil)
		}
	}

	// Setup cursor để trả về logs
	if mockLogs, ok := mockResults["logs"].([]*domain.LogEntry); ok && len(mockLogs) > 0 {
		cursor := new(mocks.MockMongoCursor)
		// Setup Next() để trả về true n lần, sau đó false
		for range mockLogs {
			cursor.On("Next", mock.Anything).Return(true).Once()
		}
		cursor.On("Next", mock.Anything).Return(false)

		// Setup Decode() để đưa log vào biến đích
		callCount := 0
		cursor.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
			if callCount < len(mockLogs) {
				arg := args.Get(0).(*domain.LogEntry)
				*arg = *mockLogs[callCount]
				callCount++
			}
		}).Return(nil)

		// Setup Close() để không có lỗi
		cursor.On("Close", mock.Anything).Return(nil)

		// Cài đặt Find() để trả về cursor
		collection.On("Find", mock.Anything, mock.Anything, mock.AnythingOfType("*options.FindOptions")).
			Return(cursor, nil)
	}
}

// SetupLoggerModels thiết lập mọi thứ cần thiết cho tests
func SetupLoggerModels() (*mocks.MockMongoClient, *mocks.MockMongoDatabase, *mocks.MockMongoCollection, domain.Models) {
	mockClient := new(mocks.MockMongoClient)
	mockDatabase := new(mocks.MockMongoDatabase)
	mockCollection := new(mocks.MockMongoCollection)

	// Setup mocks chain
	mockClient.On("Database", "logs", mock.Anything).Return(mockDatabase)
	mockDatabase.On("Collection", "logs", mock.Anything).Return(mockCollection)

	// Sử dụng SetDomainClientVar để đặt client
	realMongoClient := &mongo.Client{}
	SetDomainClientVar(realMongoClient)

	// Create models với client là mockClient
	models := domain.Models{
		LogEntry: domain.LogEntry{},
	}

	return mockClient, mockDatabase, mockCollection, models
}

// SetDomainClientVar đặt giá trị cho biến client trong domain package
// Đây là helper function dùng cho test để thay thế biến client trong package domain
func SetDomainClientVar(mongoClient interface{}) {
	// Just call domain.New with a real client and we'll mock the specific methods
	// This avoids type conversion issues
	domain.New(&mongo.Client{})
}

// CreateTestLogEntries tạo các bản ghi giả
func CreateTestLogEntries(count int) []*domain.LogEntry {
	entries := make([]*domain.LogEntry, count)
	for i := 0; i < count; i++ {
		entries[i] = &domain.LogEntry{
			ID:      primitive.NewObjectID().Hex(),
			Name:    "Test Log " + string(rune('A'+i)),
			Level:   "info",
			Service: "test-service",
			Message: "Test message " + string(rune('A'+i)),
		}
	}
	return entries
}

// CreateInsertResult tạo kết quả giả cho InsertOne
func CreateInsertResult(id interface{}) *mongo.InsertOneResult {
	return &mongo.InsertOneResult{
		InsertedID: id,
	}
}

// CreateUpdateResult tạo kết quả giả cho UpdateOne
func CreateUpdateResult(matched, modified int64) *mongo.UpdateResult {
	return &mongo.UpdateResult{
		MatchedCount:  matched,
		ModifiedCount: modified,
	}
}

// CreateTestError tạo lỗi giả cho test
func CreateTestError(message string) error {
	return errors.New(message)
}

// MongoClientAdapter wraps mock client to be used with domain.New
type MongoClientAdapter struct {
	mockClient *mocks.MockMongoClient
}

// Database forwards the call to the mock client
func (a *MongoClientAdapter) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	// Just return nil as we're not using the actual return value
	// The mock expectations will handle the calls
	return nil
}

// ModelsAdapter creates a domain.Models using a mock MongoDB client
func ModelsAdapter(mockClient *mocks.MockMongoClient) domain.Models {
	// Initialize with empty client
	models := domain.New(&mongo.Client{})

	// Set up database mocks
	mockDatabase := new(mocks.MockMongoDatabase)
	mockClient.On("Database", "logs", mock.Anything).Return(mockDatabase)

	// Set up collection mocks
	mockCollection := new(mocks.MockMongoCollection)
	mockDatabase.On("Collection", "logs", mock.Anything).Return(mockCollection)

	// Create a patch function to replace the client
	PatchMongoCollection = func() interface{} {
		return mockCollection
	}

	// Set up the domain GetMockCollection function to return our mock collection
	domain.GetMockCollection = func() interface{} {
		if PatchMongoCollection != nil {
			return PatchMongoCollection()
		}
		return nil
	}

	return models
}

// MockDatabase simulates a MongoDB database for testing
func MockDatabase(mockClient *mocks.MockMongoClient, dbName string) *mocks.MockMongoDatabase {
	mockDB := new(mocks.MockMongoDatabase)
	mockClient.On("Database", dbName, mock.Anything).Return(mockDB)
	return mockDB
}

// MockCollection simulates a MongoDB collection for testing
func MockCollection(mockDB *mocks.MockMongoDatabase, collectionName string) *mocks.MockMongoCollection {
	mockCollection := new(mocks.MockMongoCollection)
	mockDB.On("Collection", collectionName, mock.Anything).Return(mockCollection)
	return mockCollection
}

// MockCursor simulates a MongoDB cursor for testing
func MockCursor(docs []interface{}) *mocks.MockMongoCursor {
	mockCursor := new(mocks.MockMongoCursor)

	// Set up Next to return true for each document, then false
	for range docs {
		mockCursor.On("Next", mock.Anything).Return(true).Once()
	}
	mockCursor.On("Next", mock.Anything).Return(false)

	// Set up Decode to fill the provided value with each document in sequence
	var callCount int
	mockCursor.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		if callCount < len(docs) {
			// Copy the document to the provided argument
			// This would need actual implementation based on the document type
			callCount++
		}
	}).Return(nil)

	// Set up Close to return nil
	mockCursor.On("Close", mock.Anything).Return(nil)

	return mockCursor
}

// MockSingleResult simulates a MongoDB single result for testing
func MockSingleResult(doc interface{}, err error) *mocks.MockMongoSingleResult {
	mockResult := new(mocks.MockMongoSingleResult)

	if err != nil {
		mockResult.On("Decode", mock.Anything).Return(err)
	} else {
		mockResult.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
			// Copy the document to the provided argument
			// This would need actual implementation based on the document type
		}).Return(nil)
	}

	return mockResult
}

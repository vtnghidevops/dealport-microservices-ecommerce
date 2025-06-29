package service

import (
	"context"
	"logger-service/internal/domain"
	"logger-service/tests/helpers"
	"logger-service/tests/mocks"
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// mockClient được sử dụng làm biến global để thay thế client trong domain package
var mockClient *mongo.Client

func init() {
	// Tạo mock client và gán vào biến client trong domain package
	mockClient = &mongo.Client{}
	domain.New(mockClient) // Hàm này đặt biến client trong domain package
}

func setupTest() (*mocks.MockMongoClient, *mocks.MockMongoDatabase, *mocks.MockMongoCollection, domain.Models) {
	// Create mocks
	mockMongo := new(mocks.MockMongoClient)
	mockDatabase := new(mocks.MockMongoDatabase)
	mockCollection := new(mocks.MockMongoCollection)

	// Setup expectations for database and collection
	mockMongo.On("Database", "logs", mock.Anything).Return(mockDatabase)
	mockDatabase.On("Collection", "logs", mock.Anything).Return(mockCollection)

	// Tạo models sử dụng helper adapter
	models := helpers.ModelsAdapter(mockMongo)

	return mockMongo, mockDatabase, mockCollection, models
}

func TestLogEntry_Insert(t *testing.T) {
	_, _, mockCollection, models := setupTest()

	t.Run("Insert success", func(t *testing.T) {
		// Create a sample log
		logEntry := helpers.CreateSampleLogEntry(
			"Test Log",
			"info",
			"test-service",
			"test-action",
			"user123",
			"Test message",
		)

		// Setup expectation for InsertOne with specific handler - THIS IS THE KEY CHANGE
		mockID := primitive.NewObjectID()
		insertResult := helpers.CreateMongoInsertResult(mockID)

		// Instead of using the general matcher, we match exactly what will be passed
		mockCollection.On("InsertOne", context.TODO(), mock.MatchedBy(func(doc domain.LogEntry) bool {
			// Just check the basic fields to ensure it's our document
			return doc.Name == logEntry.Name &&
				doc.Level == logEntry.Level &&
				doc.Service == logEntry.Service &&
				doc.Action == logEntry.Action &&
				doc.UserID == logEntry.UserID
		}), mock.Anything).Return(insertResult, nil).Once()

		// Set up the mock collection to be returned
		domain.GetMockCollection = func() interface{} {
			return mockCollection
		}

		// Execute the method
		err := models.LogEntry.Insert(logEntry)

		// Assert
		assert.NoError(t, err)
		mockCollection.AssertExpectations(t)
	})

	t.Run("Insert error", func(t *testing.T) {
		// Create a sample log
		logEntry := helpers.CreateSampleLogEntry(
			"Failed Log",
			"error",
			"test-service",
			"test-action",
			"user123",
			"Test error message",
		)

		// Setup expectation for InsertOne to return an error
		mockCollection.On("InsertOne", context.TODO(), mock.MatchedBy(func(doc domain.LogEntry) bool {
			// Just check the basic fields to ensure it's our document
			return doc.Name == logEntry.Name &&
				doc.Level == logEntry.Level &&
				doc.Service == logEntry.Service &&
				doc.Action == logEntry.Action &&
				doc.UserID == logEntry.UserID
		}), mock.Anything).Return(nil, errors.New("database error")).Once()

		// Set up the mock collection to be returned
		domain.GetMockCollection = func() interface{} {
			return mockCollection
		}

		// Execute the method
		err := models.LogEntry.Insert(logEntry)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockCollection.AssertExpectations(t)
	})
}

func TestLogEntry_All(t *testing.T) {
	_, _, mockCollection, models := setupTest()

	t.Run("Get all logs success", func(t *testing.T) {
		// Create sample logs
		valueEntries := helpers.CreateSampleLogEntries(3, "test-service", "user123")

		// Convert value slice to pointer slice
		sampleLogs := make([]*domain.LogEntry, len(valueEntries))
		for i := range valueEntries {
			sampleLogs[i] = &valueEntries[i]
		}

		// Create and configure mock cursor
		mockCursor := new(mocks.MockMongoCursor)
		helpers.SetupMockCursor(mockCursor, sampleLogs)

		// Use mock.Anything for the third parameter to match []*options.FindOptions
		mockCollection.On("Find", context.TODO(), bson.D{}, mock.Anything).
			Return(mockCursor, nil).Once()

		// Set up the mock collection to be returned
		domain.GetMockCollection = func() interface{} {
			return mockCollection
		}

		// Execute the method
		logs, err := models.LogEntry.All()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, logs, 3)
		assert.Equal(t, sampleLogs[0].Name, logs[0].Name)
		assert.Equal(t, sampleLogs[1].Name, logs[1].Name)
		assert.Equal(t, sampleLogs[2].Name, logs[2].Name)
		mockCollection.AssertExpectations(t)
		mockCursor.AssertExpectations(t)
	})

	t.Run("Get all logs error", func(t *testing.T) {
		// Setup expectation for Find to return an error
		mockCollection.On("Find", context.TODO(), bson.D{}, mock.Anything).
			Return(nil, errors.New("database error")).Once()

		// Set up the mock collection to be returned
		domain.GetMockCollection = func() interface{} {
			return mockCollection
		}

		// Execute the method
		logs, err := models.LogEntry.All()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, logs)
		assert.Contains(t, err.Error(), "database error")
		mockCollection.AssertExpectations(t)
	})
}

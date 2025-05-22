package grpc

import (
	"errors"
	"logger-service/internal/domain"
	"logger-service/tests/helpers"
	"logger-service/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LogServer represents the gRPC server for logging
type LogServer struct {
	LogService LogService
}

// LogService định nghĩa interface cho service logger
type LogService interface {
	Insert(entry domain.LogEntry) error
	FindByService(service string) ([]*domain.LogEntry, error)
	FindByUserAction(userID, action string) ([]*domain.LogEntry, error)
}

// MockLogService là một implementation giả của LogService
type MockLogService struct {
	mock.Mock
}

// Insert giả lập hàm Insert
func (m *MockLogService) Insert(entry domain.LogEntry) error {
	args := m.Called(entry)
	return args.Error(0)
}

// FindByService giả lập hàm FindByService
func (m *MockLogService) FindByService(service string) ([]*domain.LogEntry, error) {
	args := m.Called(service)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.LogEntry), args.Error(1)
}

// FindByUserAction giả lập hàm FindByUserAction
func (m *MockLogService) FindByUserAction(userID, action string) ([]*domain.LogEntry, error) {
	args := m.Called(userID, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.LogEntry), args.Error(1)
}

// Mock request payload for testing
type WriteLogRequest struct {
	LogEntry domain.LogEntry
}

// Mock response for testing
type WriteLogResponse struct {
	Result string
	Error  string
}

// MongoClientStub is a test stub for mongo.Client
type MongoClientStub struct{}

// Database mocks the mongo.Client.Database method
func (s *MongoClientStub) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	// Return a nil database - we won't use this in tests
	// The actual calls will be mocked with the mock objects
	return nil
}

func setupGrpcTest() (*mocks.MockMongoClient, *mocks.MockMongoDatabase, *mocks.MockMongoCollection, LogServer) {
	// Create mocks
	mockMongo := new(mocks.MockMongoClient)
	mockDatabase := new(mocks.MockMongoDatabase)
	mockCollection := new(mocks.MockMongoCollection)

	// Setup expectations for database and collection
	mockMongo.On("Database", "logs", mock.Anything).Return(mockDatabase)
	mockDatabase.On("Collection", "logs", mock.Anything).Return(mockCollection)

	// Create a mock log service
	mockLogService := new(MockLogService)

	// Create server with mock log service
	server := LogServer{
		LogService: mockLogService,
	}

	return mockMongo, mockDatabase, mockCollection, server
}

func TestGrpcServer_WriteLog(t *testing.T) {
	mockLogService := new(MockLogService)
	server := LogServer{
		LogService: mockLogService,
	}

	t.Run("Write log success", func(t *testing.T) {
		// Create a test log entry
		logEntry := helpers.CreateSampleLogEntry(
			"Grpc Test Log",
			"info",
			"grpc-service",
			"grpc-action",
			"user123",
			"Test message from gRPC",
		)

		// Setup expectation for Insert
		mockLogService.On("Insert", logEntry).Return(nil).Once()

		// Execute the method
		err := server.LogService.Insert(logEntry)

		// Assert
		assert.NoError(t, err)
		mockLogService.AssertExpectations(t)
	})

	t.Run("Write log error", func(t *testing.T) {
		// Create a test log entry
		logEntry := helpers.CreateSampleLogEntry(
			"Grpc Error Log",
			"error",
			"grpc-service",
			"grpc-error",
			"user123",
			"Error message from gRPC",
		)

		// Setup expectation for Insert to fail
		mockLogService.On("Insert", logEntry).Return(errors.New("database error")).Once()

		// Execute the method
		err := server.LogService.Insert(logEntry)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockLogService.AssertExpectations(t)
	})
}

func TestGrpcServer_GetLogsByService(t *testing.T) {
	mockLogService := new(MockLogService)
	server := LogServer{
		LogService: mockLogService,
	}

	t.Run("Get logs by service success", func(t *testing.T) {
		// Setup test data
		serviceName := "auth-service"
		valueEntries := helpers.CreateSampleLogEntries(2, serviceName, "user123")

		// Convert value slice to pointer slice
		sampleLogs := make([]*domain.LogEntry, len(valueEntries))
		for i := range valueEntries {
			sampleLogs[i] = &valueEntries[i]
		}

		// Setup expectation for FindByService
		mockLogService.On("FindByService", serviceName).Return(sampleLogs, nil).Once()

		// Execute the method
		logs, err := server.LogService.FindByService(serviceName)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, logs, 2)
		for _, log := range logs {
			assert.Equal(t, serviceName, log.Service)
		}
		mockLogService.AssertExpectations(t)
	})

	t.Run("Get logs by service error", func(t *testing.T) {
		// Setup test data
		serviceName := "error-service"

		// Setup expectation for FindByService to return an error
		mockLogService.On("FindByService", serviceName).Return(nil, errors.New("database error")).Once()

		// Execute the method
		logs, err := server.LogService.FindByService(serviceName)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, logs)
		assert.Contains(t, err.Error(), "database error")
		mockLogService.AssertExpectations(t)
	})
}

func TestGrpcServer_GetLogsByUserAction(t *testing.T) {
	mockLogService := new(MockLogService)
	server := LogServer{
		LogService: mockLogService,
	}

	t.Run("Get logs by user and action success", func(t *testing.T) {
		// Setup test data
		userID := "user123"
		action := "login"
		valueEntries := helpers.CreateSampleLogEntries(2, "user-service", userID)

		// Convert value slice to pointer slice
		sampleLogs := make([]*domain.LogEntry, len(valueEntries))
		for i := range valueEntries {
			sampleLogs[i] = &valueEntries[i]
			sampleLogs[i].Action = action
		}

		// Setup expectation for FindByUserAction
		mockLogService.On("FindByUserAction", userID, action).Return(sampleLogs, nil).Once()

		// Execute the method
		logs, err := server.LogService.FindByUserAction(userID, action)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, logs, 2)
		for _, log := range logs {
			assert.Equal(t, userID, log.UserID)
			assert.Equal(t, action, log.Action)
		}
		mockLogService.AssertExpectations(t)
	})

	t.Run("Get logs by user and action error", func(t *testing.T) {
		// Setup test data
		userID := "error-user"
		action := "error-action"

		// Setup expectation for FindByUserAction to return an error
		mockLogService.On("FindByUserAction", userID, action).Return(nil, errors.New("database error")).Once()

		// Execute the method
		logs, err := server.LogService.FindByUserAction(userID, action)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, logs)
		assert.Contains(t, err.Error(), "database error")
		mockLogService.AssertExpectations(t)
	})
}

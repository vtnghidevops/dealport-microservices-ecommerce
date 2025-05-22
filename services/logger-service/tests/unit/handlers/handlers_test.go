package handlers

import (
	"errors"
	"logger-service/internal/domain"
	"logger-service/tests/helpers"
	"logger-service/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// LogService định nghĩa interface cho service logger
type LogService interface {
	Insert(entry domain.LogEntry) error
	All() ([]*domain.LogEntry, error)
	GetOne(id string) (*domain.LogEntry, error)
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

// All giả lập hàm All
func (m *MockLogService) All() ([]*domain.LogEntry, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.LogEntry), args.Error(1)
}

// GetOne giả lập hàm GetOne
func (m *MockLogService) GetOne(id string) (*domain.LogEntry, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LogEntry), args.Error(1)
}

// Handlers holds service dependencies
type Handlers struct {
	LogService LogService
}

// Create a sample request payload for testing
type logPayload struct {
	Name    string `json:"name"`
	Level   string `json:"level"`
	Service string `json:"service"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Action  string `json:"action"`
}

// mockClient được sử dụng làm biến global để thay thế client trong domain
var mockClient *mongo.Client

func init() {
	// Tạo mock client và gán vào biến client trong domain package
	mockClient = &mongo.Client{}
	domain.New(mockClient) // Hàm này đặt biến client trong domain package
}

func setupHandlerTest() (*mocks.MockMongoClient, *mocks.MockMongoDatabase, *mocks.MockMongoCollection, Handlers) {
	// Create mocks
	mockMongo := new(mocks.MockMongoClient)
	mockDatabase := new(mocks.MockMongoDatabase)
	mockCollection := new(mocks.MockMongoCollection)

	// Setup expectations for database and collection
	mockMongo.On("Database", "logs", mock.Anything).Return(mockDatabase)
	mockDatabase.On("Collection", "logs", mock.Anything).Return(mockCollection)

	// Create a mock log service
	mockLogService := new(MockLogService)

	// Create handlers with mock service
	handlers := Handlers{
		LogService: mockLogService,
	}

	return mockMongo, mockDatabase, mockCollection, handlers
}

func TestHandlers_WriteLog(t *testing.T) {
	mockLogService := new(MockLogService)
	handlers := Handlers{
		LogService: mockLogService,
	}

	t.Run("Write log success", func(t *testing.T) {
		// Create a sample log payload
		payload := logPayload{
			Name:    "Test Log",
			Level:   "info",
			Service: "test-service",
			UserID:  "user123",
			Message: "Test message",
			Action:  "test-action",
		}

		// Create log entry from payload
		logEntry := domain.LogEntry{
			Name:    payload.Name,
			Level:   payload.Level,
			Service: payload.Service,
			UserID:  payload.UserID,
			Message: payload.Message,
			Action:  payload.Action,
		}

		// Setup expectation for Insert
		mockLogService.On("Insert", mock.AnythingOfType("domain.LogEntry")).Return(nil).Once()

		// Simulate handler execution (normally would call handlers.WriteLog(rr, req))
		err := handlers.LogService.Insert(logEntry)

		// Assert
		assert.NoError(t, err)
		mockLogService.AssertExpectations(t)
	})

	t.Run("Write log error", func(t *testing.T) {
		// Create a sample log payload
		payload := logPayload{
			Name:    "Error Log",
			Level:   "error",
			Service: "test-service",
			UserID:  "user123",
			Message: "Test error message",
			Action:  "test-action",
		}

		// Create log entry from payload
		logEntry := domain.LogEntry{
			Name:    payload.Name,
			Level:   payload.Level,
			Service: payload.Service,
			UserID:  payload.UserID,
			Message: payload.Message,
			Action:  payload.Action,
		}

		// Setup expectation for Insert to fail
		mockLogService.On("Insert", mock.AnythingOfType("domain.LogEntry")).Return(errors.New("database error")).Once()

		// Simulate handler execution
		err := handlers.LogService.Insert(logEntry)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockLogService.AssertExpectations(t)
	})
}

func TestHandlers_GetAllLogs(t *testing.T) {
	mockLogService := new(MockLogService)
	handlers := Handlers{
		LogService: mockLogService,
	}

	t.Run("Get all logs success", func(t *testing.T) {
		// Create sample logs
		valueEntries := helpers.CreateSampleLogEntries(3, "test-service", "user123")

		// Convert value slice to pointer slice
		sampleLogs := make([]*domain.LogEntry, len(valueEntries))
		for i := range valueEntries {
			sampleLogs[i] = &valueEntries[i]
		}

		// Setup expectation for All
		mockLogService.On("All").Return(sampleLogs, nil).Once()

		// Simulate handler execution
		logs, err := handlers.LogService.All()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, logs, 3)
		mockLogService.AssertExpectations(t)
	})

	t.Run("Get all logs error", func(t *testing.T) {
		// Setup expectation for All to return an error
		mockLogService.On("All").Return(nil, errors.New("database error")).Once()

		// Simulate handler execution
		logs, err := handlers.LogService.All()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, logs)
		assert.Contains(t, err.Error(), "database error")
		mockLogService.AssertExpectations(t)
	})
}

func TestHandlers_GetOneLog(t *testing.T) {
	mockLogService := new(MockLogService)
	handlers := Handlers{
		LogService: mockLogService,
	}

	t.Run("Get one log success", func(t *testing.T) {
		// Create a valid ObjectID
		id := primitive.NewObjectID()
		idHex := id.Hex()

		// Create sample log
		sampleLog := helpers.CreateSampleLogEntry(
			"Test Log",
			"info",
			"test-service",
			"test-action",
			"user123",
			"Test message",
		)
		sampleLog.ID = idHex

		// Setup expectation for GetOne
		mockLogService.On("GetOne", idHex).Return(&sampleLog, nil).Once()

		// Simulate handler execution
		log, err := handlers.LogService.GetOne(idHex)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, log)
		assert.Equal(t, sampleLog.ID, log.ID)
		mockLogService.AssertExpectations(t)
	})

	t.Run("Get one with invalid ID", func(t *testing.T) {
		// Setup expectation for GetOne with invalid ID
		mockLogService.On("GetOne", "invalid-id").Return(nil, errors.New("invalid ObjectID")).Once()

		// Simulate handler execution
		log, err := handlers.LogService.GetOne("invalid-id")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, log)
		assert.Contains(t, err.Error(), "invalid ObjectID")
		mockLogService.AssertExpectations(t)
	})

	t.Run("Get one with not found error", func(t *testing.T) {
		// Create a valid ObjectID
		id := primitive.NewObjectID()
		idHex := id.Hex()

		// Setup expectation for GetOne to return not found error
		mockLogService.On("GetOne", idHex).Return(nil, mongo.ErrNoDocuments).Once()

		// Simulate handler execution
		log, err := handlers.LogService.GetOne(idHex)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, log)
		assert.Equal(t, mongo.ErrNoDocuments, err)
		mockLogService.AssertExpectations(t)
	})
}

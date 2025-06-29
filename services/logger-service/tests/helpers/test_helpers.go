package helpers

import (
	"logger-service/internal/domain"
	"logger-service/tests/mocks"
	"time"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateSampleLogEntry tạo một LogEntry mẫu cho mục đích kiểm thử
func CreateSampleLogEntry(name, level, service, action, userID, message string) domain.LogEntry {
	return domain.LogEntry{
		Name:      name,
		Level:     level,
		Service:   service,
		Action:    action,
		UserID:    userID,
		Message:   message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateSampleLogEntries tạo một danh sách LogEntry mẫu
func CreateSampleLogEntries(count int, service, userID string) []domain.LogEntry {
	entries := make([]domain.LogEntry, count)
	for i := 0; i < count; i++ {
		entries[i] = CreateSampleLogEntry(
			"Sample Log "+string(rune('A'+i)),
			"info",
			service,
			"test-action-"+string(rune('A'+i)),
			userID,
			"Test message "+string(rune('A'+i)),
		)
		entries[i].ID = primitive.NewObjectID().Hex()
	}
	return entries
}

// CreateMongoInsertResult tạo kết quả giả lập cho thao tác InsertOne
func CreateMongoInsertResult(id primitive.ObjectID) *mongo.InsertOneResult {
	return &mongo.InsertOneResult{
		InsertedID: id,
	}
}

// CreateMongoUpdateResult tạo kết quả giả lập cho thao tác UpdateOne
func CreateMongoUpdateResult(matched, modified int64, upsertedID interface{}) *mongo.UpdateResult {
	return &mongo.UpdateResult{
		MatchedCount:  matched,
		ModifiedCount: modified,
		UpsertedCount: 0,
		UpsertedID:    upsertedID,
	}
}

// SetupMockCursor cấu hình mock cursor để trả về các LogEntry khi gọi Next/Decode
func SetupMockCursor(mockCursor *mocks.MockMongoCursor, logs []*domain.LogEntry) {
	// Setup sequence of Next() calls - true for each log, then false when done
	for range logs {
		mockCursor.On("Next", mock.Anything).Return(true).Once()
	}
	mockCursor.On("Next", mock.Anything).Return(false).Once()

	// Setup Close() to return nil error
	mockCursor.On("Close", mock.Anything).Return(nil).Once()

	// Create a sequence of Decode() calls that will return each log in order
	callNumber := 0
	mockCursor.On("Decode", mock.AnythingOfType("*domain.LogEntry")).
		Run(func(args mock.Arguments) {
			if callNumber < len(logs) {
				arg := args.Get(0).(*domain.LogEntry)
				*arg = *logs[callNumber]
				callNumber++
			}
		}).Return(nil).Times(len(logs))
}

func CreateMongoDeleteResult(count int64) *mongo.DeleteResult {
	return &mongo.DeleteResult{
		DeletedCount: count,
	}
}

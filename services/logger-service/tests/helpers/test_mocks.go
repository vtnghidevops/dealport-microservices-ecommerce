package helpers

import (
	"context"
	"logger-service/internal/domain"
	"logger-service/tests/mocks"
)

// IsTestMode is whether mock testing is enabled
var IsTestMode = false

// TestInsertHandler is a function to mock Insert calls
var TestInsertHandler func(entry domain.LogEntry) error

// SetupMongoMockPatching sets up patching for MongoDB operations
func SetupMongoMockPatching(mockCollection *mocks.MockMongoCollection) {
	// Enable mock testing mode
	IsTestMode = true

	// Setup mock for Insert
	TestInsertHandler = func(entry domain.LogEntry) error {
		if mockCollection != nil {
			result, err := mockCollection.InsertOne(context.TODO(), entry, nil)
			if err != nil {
				return err
			}
			if result == nil {
				return nil
			}
			return nil
		}
		return nil
	}
}

// TeardownMongoMockPatching tears down patching for MongoDB operations
func TeardownMongoMockPatching() {
	IsTestMode = false
	TestInsertHandler = nil
}

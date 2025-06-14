package mocks

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

// MockStorageService is a mock implementation of storage.StorageService
type MockStorageService struct {
	mock.Mock
}

// UploadFile mocks the UploadFile method
func (m *MockStorageService) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	args := m.Called(ctx, objectName, reader, size, contentType)
	return args.String(0), args.Error(1)
}

// DeleteFile mocks the DeleteFile method
func (m *MockStorageService) DeleteFile(ctx context.Context, objectName string) error {
	args := m.Called(ctx, objectName)
	return args.Error(0)
}

// GetPresignedURL mocks the GetPresignedURL method
func (m *MockStorageService) GetPresignedURL(ctx context.Context, objectName string, expiry int) (string, error) {
	args := m.Called(ctx, objectName, expiry)
	return args.String(0), args.Error(1)
}

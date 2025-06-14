package storage

import (
	"context"
	"io"
)

// StorageService defines an interface for file storage operations
type StorageService interface {
	// UploadFile uploads a file to storage and returns the URL
	UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)

	// DeleteFile deletes a file from storage
	DeleteFile(ctx context.Context, objectName string) error

	// GetPresignedURL generates a presigned URL for accessing a file
	GetPresignedURL(ctx context.Context, objectName string, expiry int) (string, error)
}

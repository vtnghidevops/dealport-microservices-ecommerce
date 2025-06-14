package storage

import (
	"context"
	"log"
	"strings"
	"time"
)

// URLTransformerService handles transformation of image URLs
type URLTransformerService interface {
	TransformURL(ctx context.Context, path string) (string, error)
}

// URLTransformer implements the URLTransformerService interface
type URLTransformer struct {
	minioStorage    StorageService
	presignedURLTTL time.Duration
	minioURLPattern string
	localURLBase    string
}

// NewURLTransformer creates a new URL transformer
func NewURLTransformer(minioStorage StorageService, presignedURLTTL time.Duration) *URLTransformer {
	return &URLTransformer{
		minioStorage:    minioStorage,
		presignedURLTTL: presignedURLTTL,
		minioURLPattern: "/images/products-api/", // Pattern that identifies MinIO stored images
		localURLBase:    "/images/",              // Pattern for locally stored images
	}
}

// TransformURL converts an image path to the appropriate URL format
// For MinIO images (/images/products-api/), returns a presigned URL
// For local images (/images/), returns the original path
func (t *URLTransformer) TransformURL(ctx context.Context, path string) (string, error) {
	if path == "" {
		return "", nil
	}

	// Handle MinIO paths
	if strings.Contains(path, t.minioURLPattern) {
		// Extract the object name from the path
		objectName := strings.TrimPrefix(path, t.minioURLPattern)

		// Convert time.Duration to seconds (int)
		expirySeconds := int(t.presignedURLTTL.Seconds())

		// Generate presigned URL for MinIO object
		presignedURL, err := t.minioStorage.GetPresignedURL(ctx, objectName, expirySeconds)
		if err != nil {
			log.Printf("Error generating presigned URL for %s: %v", path, err)
			return path, err // Return original path on error
		}

		return presignedURL, nil
	}

	// Return original path for local images
	return path, nil
}

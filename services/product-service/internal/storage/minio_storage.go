package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorage implements the StorageService interface using MinIO
type MinioStorage struct {
	client       *minio.Client
	bucketName   string
	location     string
	baseURL      string
	presignedTTL time.Duration
}

// MinioConfig holds configuration for MinIO client
type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	Location        string
	BaseURL         string // Optional custom domain for MinIO
	PresignedTTL    int    // In seconds, default 3600 (1 hour)
}

// NewMinioStorage creates a new MinIO storage service
func NewMinioStorage(config MinioConfig) (*MinioStorage, error) {
	// Set default values if not provided
	if config.Location == "" {
		config.Location = "us-east-1"
	}

	if config.PresignedTTL <= 0 {
		config.PresignedTTL = 3600 // 1 hour default
	}

	log.Printf("DEBUG: Creating MinIO client with endpoint: %s", config.Endpoint)
	log.Printf("DEBUG: MinIO connection parameters - UseSSL: %v, AccessKey: %s, BucketName: %s",
		config.UseSSL, config.AccessKeyID, config.BucketName)

	// Check for empty configuration
	if config.Endpoint == "" {
		log.Printf("ERROR: MinIO endpoint is empty")
		return nil, fmt.Errorf("MinIO endpoint cannot be empty")
	}

	if config.AccessKeyID == "" || config.SecretAccessKey == "" {
		log.Printf("ERROR: MinIO credentials are missing - AccessKeyID or SecretAccessKey is empty")
		return nil, fmt.Errorf("MinIO credentials cannot be empty")
	}

	if config.BucketName == "" {
		log.Printf("ERROR: MinIO bucket name is empty")
		return nil, fmt.Errorf("MinIO bucket name cannot be empty")
	}

	// Try to resolve DNS of the endpoint before connecting
	log.Printf("DEBUG: Trying to resolve DNS for %s...", config.Endpoint)
	addrs, err := net.LookupHost(strings.Split(config.Endpoint, ":")[0])
	if err != nil {
		log.Printf("ERROR: DNS resolution failed for %s: %v", config.Endpoint, err)
		log.Printf("DEBUG: This indicates a network connectivity or DNS issue")
		return nil, fmt.Errorf("DNS resolution failed for MinIO endpoint: %w", err)
	}
	log.Printf("DEBUG: DNS resolution successful for %s: %v", config.Endpoint, addrs)

	// Initialize MinIO client
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		log.Printf("ERROR: Failed to create MinIO client: %v", err)
		log.Printf("DEBUG: Common issues include:")
		log.Printf("DEBUG: 1. Invalid endpoint format - should be 'domain:port' without http/https")
		log.Printf("DEBUG: 2. DNS resolution failure - check host file or DNS settings")
		log.Printf("DEBUG: 3. Network connectivity - firewall or routing issues")
		log.Printf("DEBUG: 4. SSL issues if UseSSL is true - certificate problems")
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	log.Printf("DEBUG: MinIO client created successfully, testing with API call...")

	// Test with a basic API call
	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		log.Printf("WARNING: Failed to list buckets: %v", err)
		log.Printf("DEBUG: Client created but server may not be reachable")
	} else {
		log.Printf("DEBUG: Successfully connected to MinIO server - Found %d buckets", len(buckets))
	}

	// Check if bucket exists, create if it doesn't
	exists, err := client.BucketExists(context.Background(), config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		log.Printf("Bucket %s does not exist, creating...", config.BucketName)
		err = client.MakeBucket(context.Background(), config.BucketName, minio.MakeBucketOptions{
			Region: config.Location,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Created bucket %s", config.BucketName)
	}

	// Create storage service
	return &MinioStorage{
		client:       client,
		bucketName:   config.BucketName,
		location:     config.Location,
		baseURL:      config.BaseURL,
		presignedTTL: time.Duration(config.PresignedTTL) * time.Second,
	}, nil
}

// UploadFile uploads a file to MinIO storage
func (s *MinioStorage) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	// Ensure object name doesn't start with slash
	objectName = strings.TrimPrefix(objectName, "/")

	log.Printf("UploadFile: Uploading to MinIO - object name: %s, size: %d, content-type: %s",
		objectName, size, contentType)

	// Upload the file
	info, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Printf("ERROR: Failed to upload file to MinIO: %v", err)
		return "", fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	log.Printf("Successfully uploaded %s of size %d to bucket %s", objectName, info.Size, s.bucketName)

	// Generate relative URL path with special MinIO marker
	// Use /images/products-api/ prefix for MinIO-stored images to distinguish from local storage
	// This helps differentiate from local /images/products/ folder in the public directory
	var relativeURL string

	// Ensure consistent URL pattern - always use /images/products-api/ for MinIO images
	if strings.HasPrefix(objectName, "products/") {
		relativeURL = fmt.Sprintf("/images/products-api/%s", strings.TrimPrefix(objectName, "products/"))
	} else {
		// If somehow the object doesn't have the products/ prefix, add it anyway
		relativeURL = fmt.Sprintf("/images/products-api/%s", filepath.Base(objectName))
	}

	// Important: The special products-api/ path will be used by frontend to identify MinIO images
	log.Printf("Stored image with relative URL: %s (MinIO storage)", relativeURL)
	return relativeURL, nil
}

// DeleteFile deletes a file from MinIO storage
func (s *MinioStorage) DeleteFile(ctx context.Context, objectName string) error {
	// Ensure object name doesn't start with slash
	objectName = strings.TrimPrefix(objectName, "/")

	// Extract the object name from the URL if it's a full URL
	if strings.HasPrefix(objectName, "images/products-api/") {
		// Convert from new pattern to actual object name in MinIO
		objectName = "products/" + strings.TrimPrefix(objectName, "images/products-api/")
	} else if strings.HasPrefix(objectName, "images/") {
		objectName = strings.TrimPrefix(objectName, "images/")
	}

	// Remove the file
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	log.Printf("Successfully deleted %s from bucket %s", objectName, s.bucketName)
	return nil
}

// GetPresignedURL generates a presigned URL for accessing a file
func (s *MinioStorage) GetPresignedURL(ctx context.Context, objectName string, expiry int) (string, error) {
	// Remove any leading slash
	objectName = strings.TrimPrefix(objectName, "/")

	// Log original object name for debugging
	log.Printf("GetPresignedURL: Original object name: %s", objectName)

	// The actual structure in MinIO:
	// - Bucket name: images
	// - Path within bucket: products/filename.png

	// Extract the object name from the URL if it's a relative URL
	if strings.HasPrefix(objectName, "images/products-api/") {
		// Convert from URL pattern to actual object name in MinIO
		objectName = "products/" + strings.TrimPrefix(objectName, "images/products-api/")
		log.Printf("GetPresignedURL: After products-api pattern conversion: %s", objectName)
	} else if strings.Contains(objectName, "/images/products-api/") {
		// Handle URLs starting with /images/products-api/
		parts := strings.Split(objectName, "/images/products-api/")
		if len(parts) > 1 {
			objectName = "products/" + parts[1]
			log.Printf("GetPresignedURL: After /images/products-api/ pattern conversion: %s", objectName)
		}
	} else if strings.HasPrefix(objectName, "images/products/") {
		// Handle old URL pattern - Already in correct format
		objectName = strings.TrimPrefix(objectName, "images/")
		log.Printf("GetPresignedURL: Converting products path: %s", objectName)
	} else if strings.Contains(objectName, "/images/products/") {
		// Handle URLs starting with /images/products/
		parts := strings.Split(objectName, "/images/products/")
		if len(parts) > 1 {
			objectName = "products/" + parts[1]
			log.Printf("GetPresignedURL: Converting with /images/products/ pattern: %s", objectName)
		}
	} else if strings.HasPrefix(objectName, "images/") {
		// Already in the right format
		objectName = strings.TrimPrefix(objectName, "images/")
		log.Printf("GetPresignedURL: Removed images/ prefix: %s", objectName)
	}

	// If expiry is not provided or invalid, use the default
	expiryDuration := s.presignedTTL
	if expiry > 0 {
		expiryDuration = time.Duration(expiry) * time.Second
	}

	// Log the final object name for debugging
	log.Printf("Generating presigned URL for object name: %s in bucket %s", objectName, s.bucketName)

	// Create request params with improved security settings
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(objectName)))

	// Generate presigned URL with security improvements
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, expiryDuration, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// Log only the base URL without query params for security
	baseURL := fmt.Sprintf("%s://%s%s", presignedURL.Scheme, presignedURL.Host, presignedURL.Path)
	log.Printf("Generated presigned URL for %s (credentials omitted from logs)", baseURL)

	// If baseURL is provided, replace the hostname in the presigned URL
	if s.baseURL != "" {
		// This is a simple way to replace the hostname, might need to be more sophisticated
		// depending on your URL structure
		urlStr := presignedURL.String()
		parsedBaseURL := fmt.Sprintf("%s://%s", presignedURL.Scheme, presignedURL.Host)
		return strings.Replace(urlStr, parsedBaseURL, s.baseURL, 1), nil
	}

	return presignedURL.String(), nil
}

// GetObjectName extracts the object name from a URL path
func (s *MinioStorage) GetObjectName(urlPath string) string {
	// Remove any leading slash
	urlPath = strings.TrimPrefix(urlPath, "/")

	// Handle special MinIO path pattern
	if strings.HasPrefix(urlPath, "images/products-api/") {
		// Convert from URL pattern to actual object name in MinIO
		return "products/" + strings.TrimPrefix(urlPath, "images/products-api/")
	}

	// If it starts with "images/", remove it
	if strings.HasPrefix(urlPath, "images/") {
		urlPath = strings.TrimPrefix(urlPath, "images/")
	}

	// Extract just the filename
	return filepath.Base(urlPath)
}

// BucketExists checks if the configured bucket exists
func (s *MinioStorage) BucketExists(ctx context.Context) (bool, error) {
	if s.client == nil {
		return false, fmt.Errorf("MinIO client is not initialized")
	}
	return s.client.BucketExists(ctx, s.bucketName)
}

// CreateBucket creates a new bucket if it doesn't exist
func (s *MinioStorage) CreateBucket(ctx context.Context, bucketName, location string) error {
	if s.client == nil {
		return fmt.Errorf("MinIO client is not initialized")
	}
	return s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
}

// TestConnection tests connectivity to the MinIO server
func (s *MinioStorage) TestConnection(ctx context.Context) error {
	log.Printf("Testing connection to MinIO server at %s, bucket %s...", s.client.EndpointURL().Host, s.bucketName)

	// Basic network connectivity test by trying to list buckets
	_, err := s.client.ListBuckets(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to list buckets: %v", err)

		// Check specific error type to provide more detailed guidance
		if strings.Contains(err.Error(), "connection refused") {
			log.Printf("ERROR: Connection refused - MinIO server is not running or not accessible")
			return fmt.Errorf("connection refused to MinIO server: %w", err)
		} else if strings.Contains(err.Error(), "no such host") {
			log.Printf("ERROR: Host not found - DNS resolution failure")
			return fmt.Errorf("host not found for MinIO endpoint: %w", err)
		} else if strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "timeout") {
			log.Printf("ERROR: Connection timeout - MinIO server is unreachable")
			return fmt.Errorf("connection timeout to MinIO server: %w", err)
		}

		// Try to ping endpoint to see if it's reachable
		parsed, _ := url.Parse(s.client.EndpointURL().String())
		hostname := parsed.Hostname()
		log.Printf("DEBUG: Trying to ping %s...", hostname)

		// Try to use TCP dial to test connection
		log.Printf("DEBUG: Trying TCP connection to %s...", s.client.EndpointURL().Host)
		timeout := 5 * time.Second
		conn, dialErr := net.DialTimeout("tcp", s.client.EndpointURL().Host, timeout)
		if dialErr != nil {
			log.Printf("ERROR: TCP dial failed: %v", dialErr)
			log.Printf("DEBUG: This indicates firewall or routing issues")
		} else {
			conn.Close()
			log.Printf("DEBUG: TCP connection successful to %s", s.client.EndpointURL().Host)
			log.Printf("DEBUG: Issue might be related to authentication or MinIO configuration")
		}

		return fmt.Errorf("failed to connect to MinIO server: %w", err)
	}

	// Continue with bucket check
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		log.Printf("ERROR: Bucket '%s' does not exist", s.bucketName)
		return fmt.Errorf("bucket '%s' does not exist", s.bucketName)
	}

	log.Printf("DEBUG: Successfully verified MinIO connection and bucket '%s' exists", s.bucketName)
	return nil
}

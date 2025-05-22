package security

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// TestImageUploadSecurity tests the security aspects of image uploads
func TestImageUploadSecurity(t *testing.T) {
	// Test file extension validation
	t.Run("File Extension Validation", func(t *testing.T) {
		validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
		invalidExtensions := []string{".exe", ".php", ".js", ".html", ".svg"}

		// Test valid extensions
		for _, ext := range validExtensions {
			filename := "test-image" + ext
			if !isValidImageExtension(filename) {
				t.Errorf("Extension %s should be valid, but was rejected", ext)
			}
		}

		// Test invalid extensions
		for _, ext := range invalidExtensions {
			filename := "test-image" + ext
			if isValidImageExtension(filename) {
				t.Errorf("Extension %s should be invalid, but was accepted", ext)
			}
		}
	})

	// Test file size validation
	t.Run("File Size Validation", func(t *testing.T) {
		// Test within size limit
		if !isValidFileSize(2 * 1024 * 1024) { // 2MB
			t.Error("2MB file should be valid, but was rejected")
		}

		// Test at max size limit
		if !isValidFileSize(5 * 1024 * 1024) { // 5MB
			t.Error("5MB file should be valid, but was rejected")
		}

		// Test exceeding size limit
		if isValidFileSize(6 * 1024 * 1024) { // 6MB
			t.Error("6MB file should be invalid, but was accepted")
		}
	})

	// Test file content validation (basic MIME type checking)
	t.Run("File Content Validation", func(t *testing.T) {
		// Valid JPG header
		jpgHeader := []byte{0xFF, 0xD8, 0xFF}
		if !hasValidImageHeader(jpgHeader) {
			t.Error("Valid JPG header should be accepted")
		}

		// Valid PNG header
		pngHeader := []byte{0x89, 0x50, 0x4E, 0x47}
		if !hasValidImageHeader(pngHeader) {
			t.Error("Valid PNG header should be accepted")
		}

		// Invalid header (text file)
		textHeader := []byte("<?php")
		if hasValidImageHeader(textHeader) {
			t.Error("PHP file header should be rejected")
		}
	})

	// Test path traversal prevention
	t.Run("Path Traversal Prevention", func(t *testing.T) {
		testCases := []struct {
			filename string
			isValid  bool
		}{
			{"normal.jpg", true},
			{"../../../etc/passwd.jpg", false},
			{"..\\..\\Windows\\system32\\config.jpg", false},
			{"image_name.jpg", true},
			{"/etc/passwd.jpg", false},
			{"C:\\Windows\\system32\\config.jpg", false},
		}

		for _, tc := range testCases {
			if isSafeFilename(tc.filename) != tc.isValid {
				t.Errorf("Expected %v for filename %s, got %v",
					tc.isValid, tc.filename, !tc.isValid)
			}
		}
	})
}

// isValidImageExtension validates file extensions
func isValidImageExtension(filename string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	ext := strings.ToLower(filepath.Ext(filename))

	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}

	return false
}

// isValidFileSize validates file size
func isValidFileSize(size int64) bool {
	// Maximum file size: 5MB
	maxSize := int64(5 * 1024 * 1024)
	return size > 0 && size <= maxSize
}

// hasValidImageHeader does basic validation of file content
func hasValidImageHeader(content []byte) bool {
	// This is a simplified version for testing
	// Real implementation would do proper content validation

	// Check for JPG header
	if len(content) >= 3 &&
		content[0] == 0xFF &&
		content[1] == 0xD8 &&
		content[2] == 0xFF {
		return true
	}

	// Check for PNG header
	if len(content) >= 4 &&
		content[0] == 0x89 &&
		content[1] == 0x50 &&
		content[2] == 0x4E &&
		content[3] == 0x47 {
		return true
	}

	// Check for GIF header
	if len(content) >= 6 &&
		bytes.Equal(content[:6], []byte("GIF87a")) ||
		bytes.Equal(content[:6], []byte("GIF89a")) {
		return true
	}

	return false
}

// isSafeFilename prevents path traversal attacks
func isSafeFilename(filename string) bool {
	// Check for path traversal patterns
	if strings.Contains(filename, "..") ||
		strings.Contains(filename, "/") ||
		strings.Contains(filename, "\\") {
		return false
	}

	// Ensure the filename is just a basename
	if filepath.Base(filename) != filename {
		return false
	}

	return true
}

package mocks

import (
	"bytes"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

// MockFileUpload is a mock implementation of the domain.FileUpload interface
type MockFileUpload struct {
	mock.Mock
	filename string
	size     int64
	content  []byte
}

// NewMockFileUpload creates a new MockFileUpload with the given parameters
func NewMockFileUpload(filename string, size int64, content []byte) *MockFileUpload {
	return &MockFileUpload{
		filename: filename,
		size:     size,
		content:  content,
	}
}

// Open provides a mock implementation of the Open method
func (m *MockFileUpload) Open() (multipart.File, error) {
	args := m.Called()

	if args.Get(0) != nil {
		return args.Get(0).(multipart.File), args.Error(1)
	}

	// If no mock behavior is specified, return a default implementation
	return &mockFile{bytes.NewReader(m.content)}, nil
}

// Filename provides a mock implementation of the Filename method
func (m *MockFileUpload) Filename() string {
	args := m.Called()

	// If a return value is specified in the test
	if args.Get(0) != nil {
		return args.String(0)
	}

	// Otherwise return the default value
	return m.filename
}

// Size provides a mock implementation of the Size method
func (m *MockFileUpload) Size() int64 {
	args := m.Called()

	// If a return value is specified in the test
	if ret, ok := args.Get(0).(int64); ok && ret != 0 {
		return ret
	}

	// Otherwise return the default value
	return m.size
}

// mockFile is a simple implementation of multipart.File interface
type mockFile struct {
	*bytes.Reader
}

// Close implements the Close method for multipart.File
func (m *mockFile) Close() error {
	return nil
}

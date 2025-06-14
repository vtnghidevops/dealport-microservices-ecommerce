package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockImageURLCache is a mock implementation of cache.ImageURLCache
type MockImageURLCache struct {
	mock.Mock
}

// Get mocks the Get method
func (m *MockImageURLCache) Get(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

// Set mocks the Set method
func (m *MockImageURLCache) Set(key, value string) {
	m.Called(key, value)
}

// Delete mocks the Delete method
func (m *MockImageURLCache) Delete(key string) {
	m.Called(key)
}

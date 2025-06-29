package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockConfig is a mock implementation of the config.Config struct
type MockConfig struct {
	mock.Mock
}

// No need for NewMockConfig anymore - we're creating proper configs in each test file

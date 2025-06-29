package unit

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

// TestConnectToRabbitMQ_Success tests the connectToRabbitMQ function when the connection succeeds
func TestConnectToRabbitMQ_Success(t *testing.T) {
	// Save the original function and restore it after test
	originalDialFunc := amqpDial
	defer func() { amqpDial = originalDialFunc }()

	// Mock the Dial function
	mockConn := &amqp091.Connection{}
	amqpDial = func(url string) (*amqp091.Connection, error) {
		return mockConn, nil
	}

	// Set environment variable for test
	os.Setenv("RABBIT_URL", "amqp://guest:guest@testhost:5672/")
	defer os.Unsetenv("RABBIT_URL")

	// Call the function
	conn, err := connectToRabbitMQ()

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, mockConn, conn)
}

// TestConnectToRabbitMQ_EventualSuccess tests the connectToRabbitMQ function when the connection succeeds after retries
func TestConnectToRabbitMQ_EventualSuccess(t *testing.T) {
	// Save the original function and restore it after test
	originalDialFunc := amqpDial
	defer func() { amqpDial = originalDialFunc }()

	// Counter for attempts
	attempts := 0
	maxAttempts := 2

	// Mock the Dial function to fail initially then succeed
	mockConn := &amqp091.Connection{}
	amqpDial = func(url string) (*amqp091.Connection, error) {
		attempts++
		if attempts <= maxAttempts {
			return nil, errors.New("connection failed")
		}
		return mockConn, nil
	}

	// Save original sleep function and restore after test
	originalSleep := timeSleep
	defer func() { timeSleep = originalSleep }()

	// Mock the sleep function to speed up test
	timeSleep = func(d time.Duration) {}

	// Call the function
	conn, err := connectToRabbitMQ()

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, mockConn, conn)
	assert.Equal(t, maxAttempts+1, attempts)
}

// TestConnectToRabbitMQ_Failure tests the connectToRabbitMQ function when the connection keeps failing
func TestConnectToRabbitMQ_Failure(t *testing.T) {
	// Save the original function and restore it after test
	originalDialFunc := amqpDial
	defer func() { amqpDial = originalDialFunc }()

	// Mock the Dial function to always fail
	expectedErr := errors.New("connection failed")
	amqpDial = func(url string) (*amqp091.Connection, error) {
		return nil, expectedErr
	}

	// Save original sleep function and restore after test
	originalSleep := timeSleep
	defer func() { timeSleep = originalSleep }()

	// Mock the sleep function to speed up test
	timeSleep = func(d time.Duration) {}

	// Call the function
	conn, err := connectToRabbitMQ()

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, conn)
}

// These variables will need to be defined in your main package
var (
	amqpDial   = amqp091.Dial
	timeSleep  = time.Sleep
)

// Function signature of connectToRabbitMQ for reference
func connectToRabbitMQ() (*amqp091.Connection, error) {
	var counts int64
	var connection *amqp091.Connection
	var err error

	// Get RabbitMQ connection details from environment or use defaults
	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	// Try to connect to RabbitMQ with backoff
	for {
		connection, err = amqpDial(rabbitURL)
		if err != nil {
			counts++
		} else {
			break
		}

		if counts > 5 {
			return nil, err
		}

		// Calculate backoff time using exponential backoff
		backOff := time.Second * time.Duration(2 << (counts - 1))
		timeSleep(backOff)
		continue
	}

	return connection, nil
} 
package mocks

import (
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
)

// MockAMQPConnection is a mock for the amqp.Connection
type MockAMQPConnection struct {
	mock.Mock
}

// AsRealConnection converts the mock to *amqp091.Connection for compatibility
func (m *MockAMQPConnection) AsRealConnection() *amqp091.Connection {
	// We can't directly return a type compatible with *amqp091.Connection
	// This method should only be used when directly calling NewConsumer
	// For tests, use CreateTestConsumer helper from the helpers package instead
	// which sets up the connection properly through reflection
	return nil
}

// wrappedMockConnection is a wrapper around MockAMQPConnection that implements amqp091.Connection methods
type wrappedMockConnection struct {
	mock *MockAMQPConnection
}

// Channel implements the Channel method of amqp091.Connection
func (w *wrappedMockConnection) Channel() (*amqp091.Channel, error) {
	return w.mock.Channel()
}

// NotifyClose implements the NotifyClose method of amqp091.Connection
func (w *wrappedMockConnection) NotifyClose(receiver chan *amqp091.Error) chan *amqp091.Error {
	args := w.mock.Called(receiver)
	var ch chan *amqp091.Error
	if args.Get(0) != nil {
		ch = args.Get(0).(chan *amqp091.Error)
	}
	return ch
}

// Close mocks the Close method
func (m *MockAMQPConnection) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Channel mocks the Channel method
func (m *MockAMQPConnection) Channel() (*amqp091.Channel, error) {
	args := m.Called()

	var channel *amqp091.Channel
	if arg := args.Get(0); arg != nil {
		channel = arg.(*amqp091.Channel)
	}

	return channel, args.Error(1)
}

// MockAMQPChannel is a mock for the amqp.Channel
type MockAMQPChannel struct {
	mock.Mock
}

// Close mocks the Close method
func (m *MockAMQPChannel) Close() error {
	args := m.Called()
	return args.Error(0)
}

// ExchangeDeclare mocks the ExchangeDeclare method
func (m *MockAMQPChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp091.Table) error {
	mockArgs := m.Called(name, kind, durable, autoDelete, internal, noWait, args)
	return mockArgs.Error(0)
}

// QueueDeclare mocks the QueueDeclare method
func (m *MockAMQPChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) (amqp091.Queue, error) {
	mockArgs := m.Called(name, durable, autoDelete, exclusive, noWait, args)

	var queue amqp091.Queue
	if arg := mockArgs.Get(0); arg != nil {
		queue = arg.(amqp091.Queue)
	}

	return queue, mockArgs.Error(1)
}

// QueueBind mocks the QueueBind method
func (m *MockAMQPChannel) QueueBind(name, key, exchange string, noWait bool, args amqp091.Table) error {
	mockArgs := m.Called(name, key, exchange, noWait, args)
	return mockArgs.Error(0)
}

// Consume mocks the Consume method
func (m *MockAMQPChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error) {
	mockArgs := m.Called(queue, consumer, autoAck, exclusive, noLocal, noWait, args)

	var deliveries <-chan amqp091.Delivery
	if arg := mockArgs.Get(0); arg != nil {
		deliveries = arg.(<-chan amqp091.Delivery)
	}

	return deliveries, mockArgs.Error(1)
}

// Publish mocks the Publish method
func (m *MockAMQPChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error {
	mockArgs := m.Called(exchange, key, mandatory, immediate, msg)
	return mockArgs.Error(0)
}

// Qos mocks the Qos method
func (m *MockAMQPChannel) Qos(prefetchCount, prefetchSize int, global bool) error {
	mockArgs := m.Called(prefetchCount, prefetchSize, global)
	return mockArgs.Error(0)
}

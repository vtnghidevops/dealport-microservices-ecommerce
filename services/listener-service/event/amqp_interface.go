package event

import "github.com/rabbitmq/amqp091-go"

// AMQPConnection is an interface for *amqp091.Connection
type AMQPConnection interface {
	Channel() (*amqp091.Channel, error)
	Close() error
}

// AMQPChannel is an interface for *amqp091.Channel
type AMQPChannel interface {
	ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp091.Table) error
	QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp091.Table) (amqp091.Queue, error)
	QueueBind(name string, key string, exchange string, noWait bool, args amqp091.Table) error
	Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error)
	Close() error
}

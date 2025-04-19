package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection * amqp.Connection
}

// sets up the RabbitMQ exchange that the emitter will use to publish messages.
func (e *Emitter) setup() error {
	channel,err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()
	return declareExchange(channel)

}

// Publishes a message (event) to the specified exchange with a given routing_key
func (e * Emitter) Push(event string, routing_key string) error {
	channel,err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()
	log.Println("Pushing to channel:", channel)
	err = channel.Publish(
		"logs_topic",
		routing_key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(event),
		},
	)

	if err != nil {
		return err
	}
	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
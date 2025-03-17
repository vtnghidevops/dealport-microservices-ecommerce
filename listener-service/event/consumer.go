package event

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"bytes"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
	queueName string
}

// Create a new consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn, 
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}


// setup a channel from connect and declare exchange 
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}
// Message consumer to recv
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Listen from a queue
// Topics is a list binding keys 
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	// bind queue with exchange
	for _, binding_key := range topics { 
		ch.QueueBind(
			q.Name,
			binding_key,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	// Start consuming messages from the queue.
	// Parameters explained:
	//   - q.Name: the name of the queue to consume from.
	//   - "" : consumer tag (auto-generated if empty).
	//   - true: enable auto-acknowledgment of messages.
	//   - false for exclusive, no-local, no-wait.
	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}


	// make a channel async to hold program is running
	forever := make(chan bool)
	// handle mess recv from queue
	go func() {
		for d := range messages {
			var payload Payload
			// json -> payload type
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	// Block forever (or until the program is terminated) to keep listening for messages.
	<-forever

	return nil
}

// recv and handle each message based on name attribute of payload
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// authenticate

	// you can have as many cases as you want, as long as you write the logic

	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service:9001/logs"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}
	
	return nil
}
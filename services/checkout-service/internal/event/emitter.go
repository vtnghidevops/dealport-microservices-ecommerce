package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// StandardEvent chuẩn hóa format cho mọi event
type StandardEvent struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Data        interface{} `json:"data,omitempty"`
	DataSchema  string      `json:"data_schema"`
	Source      string      `json:"source"`
	CreatedAt   time.Time   `json:"created_at"`
	PublishedAt time.Time   `json:"published_at"`
	Version     string      `json:"version"`
}

// OrderEventData chứa dữ liệu của event đơn hàng
type OrderEventData struct {
	OrderID         string     `json:"order_id"`
	OrderNumber     string     `json:"order_number"`
	UserID          string     `json:"user_id"`
	UserEmail       string     `json:"user_email"`
	Status          string     `json:"status"`
	PreviousStatus  string     `json:"previous_status,omitempty"`
	PaymentMethod   string     `json:"payment_method"`
	Total           float64    `json:"total"`
	CreatedAt       time.Time  `json:"created_at"`
	Items           []ItemData `json:"items,omitempty"`
	ShippingName    string     `json:"shipping_name,omitempty"`
	ShippingAddress string     `json:"shipping_address,omitempty"`
	ShippingPhone   string     `json:"shipping_phone,omitempty"`
}

// ItemData thông tin sản phẩm trong đơn hàng
type ItemData struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// PaymentEventData chứa dữ liệu của event thanh toán
type PaymentEventData struct {
	OrderID       string  `json:"order_id"`
	OrderNumber   string  `json:"order_number"`
	UserID        string  `json:"user_id"`
	UserEmail     string  `json:"user_email"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	TransactionID string  `json:"transaction_id"`
	Status        string  `json:"status"`
	PaymentDate   string  `json:"payment_date"`
}

// Emitter kết nối với RabbitMQ
type Emitter struct {
	conn   *amqp.Connection
	logger *log.Logger
}

// NewEmitter tạo emitter mới
func NewEmitter(conn *amqp.Connection, logger *log.Logger) (*Emitter, error) {
	if conn == nil {
		return nil, fmt.Errorf("RabbitMQ connection is nil")
	}

	emitter := &Emitter{
		conn:   conn,
		logger: logger,
	}

	// Đảm bảo exchange tồn tại
	if err := emitter.setupExchange(); err != nil {
		return nil, err
	}

	return emitter, nil
}

// Thiết lập exchange
func (e *Emitter) setupExchange() error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

// EmitOrderCreated phát event khi tạo đơn hàng mới
func (e *Emitter) EmitOrderCreated(data OrderEventData) error {
	e.logger.Printf("EVENT-DEBUG: Creating order.created event - order_id=%s", data.OrderID)

	event := StandardEvent{
		ID:         uuid.New().String(),
		Name:       "order.created",
		Data:       data,
		DataSchema: "v1",
		Source:     "checkout-service",
		CreatedAt:  time.Now(),
		Version:    "1.0",
	}

	// Log JSON data của event
	jsonData, _ := json.Marshal(data)
	e.logger.Printf("EVENT-DEBUG: order.created event payload: %s", string(jsonData))

	return e.publish("order.created", event)
}

// EmitPaymentSucceeded phát event khi thanh toán thành công
func (e *Emitter) EmitPaymentSucceeded(data PaymentEventData) error {
	e.logger.Printf("EVENT-DEBUG: Creating order.payment_succeeded event - order_id=%s, method=%s",
		data.OrderID, data.PaymentMethod)

	// Log chi tiết hơn cho các phương thức thanh toán cụ thể
	if data.PaymentMethod == "cod" {
		e.logger.Printf("EVENT-DEBUG: Special handling for COD payment event - order %s", data.OrderID)
	}

	// Tạo event ID duy nhất
	eventID := uuid.New().String()

	event := StandardEvent{
		ID:         eventID,
		Name:       "order.payment_succeeded",
		Data:       data,
		DataSchema: "v1",
		Source:     "checkout-service",
		CreatedAt:  time.Now(),
		Version:    "1.0",
	}

	// Log JSON data của event
	jsonData, _ := json.Marshal(data)
	e.logger.Printf("EVENT-DEBUG: order.payment_succeeded event payload: %s", string(jsonData))

	// Log thông tin trước khi gửi
	e.logger.Printf("EVENT-DEBUG: About to publish order.payment_succeeded event (ID: %s) for order %s",
		eventID, data.OrderID)

	// Publish sự kiện
	err := e.publish("order.payment_succeeded", event)
	if err != nil {
		e.logger.Printf("EVENT-ERROR: Failed to publish payment_succeeded event: %v", err)
		return err
	}

	e.logger.Printf("EVENT-SUCCESS: Successfully published payment_succeeded event (ID: %s) for order %s",
		eventID, data.OrderID)
	return nil
}

// EmitPaymentFailed phát event khi thanh toán thất bại
func (e *Emitter) EmitPaymentFailed(data PaymentEventData) error {
	e.logger.Printf("EVENT-DEBUG: Creating order.payment_failed event - order_id=%s, method=%s",
		data.OrderID, data.PaymentMethod)

	event := StandardEvent{
		ID:         uuid.New().String(),
		Name:       "order.payment_failed",
		Data:       data,
		DataSchema: "v1",
		Source:     "checkout-service",
		CreatedAt:  time.Now(),
		Version:    "1.0",
	}

	// Log JSON data của event
	jsonData, _ := json.Marshal(data)
	e.logger.Printf("EVENT-DEBUG: order.payment_failed event payload: %s", string(jsonData))

	return e.publish("order.payment_failed", event)
}

// EmitOrderStatusChanged phát event khi trạng thái đơn hàng thay đổi
func (e *Emitter) EmitOrderStatusChanged(data OrderEventData) error {
	e.logger.Printf("EVENT-DEBUG: Creating order.status_changed event - order_id=%s, status=%s",
		data.OrderID, data.Status)

	event := StandardEvent{
		ID:         uuid.New().String(),
		Name:       "order.status_changed",
		Data:       data,
		DataSchema: "v1",
		Source:     "checkout-service",
		CreatedAt:  time.Now(),
		Version:    "1.0",
	}

	// Log JSON data của event
	jsonData, _ := json.Marshal(data)
	e.logger.Printf("EVENT-DEBUG: order.status_changed event payload: %s", string(jsonData))

	return e.publish("order.status_changed", event)
}

// publish phương thức chung để phát event
func (e *Emitter) publish(routingKey string, event StandardEvent) error {
	e.logger.Printf("EVENT-DEBUG: Starting to publish event %s (ID: %s) with routing key %s",
		event.Name, event.ID, routingKey)

	ch, err := e.conn.Channel()
	if err != nil {
		errMsg := fmt.Errorf("failed to open channel: %w", err)
		e.logger.Printf("EVENT-ERROR: %v", errMsg)
		return errMsg
	}
	defer ch.Close()

	// Cập nhật thời gian phát
	event.PublishedAt = time.Now()

	// Chuyển event sang JSON
	jsonData, err := json.Marshal(event)
	if err != nil {
		errMsg := fmt.Errorf("failed to marshal event: %w", err)
		e.logger.Printf("EVENT-ERROR: %v", errMsg)
		return errMsg
	}

	// Log nội dung event để debug
	e.logger.Printf("EVENT-DEBUG: Event JSON data: %s", string(jsonData))

	// Tạo context với timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Phát message
	err = ch.PublishWithContext(
		ctx,
		"logs_topic", // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		errMsg := fmt.Errorf("failed to publish event: %w", err)
		e.logger.Printf("EVENT-ERROR: %v", errMsg)
		return errMsg
	}

	e.logger.Printf("EVENT-SUCCESS: Published event %s (ID: %s) with routing key %s",
		event.Name, event.ID, routingKey)
	return nil
}

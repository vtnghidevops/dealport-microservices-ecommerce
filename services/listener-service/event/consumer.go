package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userpb "listener-service/proto/user"
)

// StandardEvent represents a standardized message format for all events
type StandardEvent struct {
	ID          string      `json:"id"`             // Unique event identifier
	Name        string      `json:"name"`           // Event name/type (e.g., "user.registered")
	Data        interface{} `json:"data,omitempty"` // Event payload data
	DataSchema  string      `json:"data_schema"`    // Schema version for the data
	Source      string      `json:"source"`         // Source service that created the event
	CreatedAt   time.Time   `json:"created_at"`     // Timestamp when event was created
	PublishedAt time.Time   `json:"published_at"`   // Timestamp when event was published
	Version     string      `json:"version"`        // Event schema version
}

// UserRegistered represents the data structure for user.registered events
type UserRegistered struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// PasswordResetRequested represents the data structure for user.password_reset_requested events
type PasswordResetRequested struct {
	Email     string `json:"email"`
	TokenHash string `json:"token_hash"`
	ExpiresAt string `json:"expires_at"`
}

// PasswordChanged represents the data structure for user.password_changed events
type PasswordChanged struct {
	Email     string `json:"email"`
	ChangedAt string `json:"changed_at"`
}

// EmailData represents the data structure for email.send events
type EmailData struct {
	Type      string            `json:"type"`
	From      string            `json:"from,omitempty"`
	FromName  string            `json:"from_name,omitempty"`
	To        string            `json:"to"`
	Subject   string            `json:"subject"`
	Template  string            `json:"template,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
}

// OrderCreatedData represents the data for order.created events
type OrderCreatedData struct {
	OrderID         string     `json:"order_id"`
	OrderNumber     string     `json:"order_number"`
	UserID          string     `json:"user_id"`
	UserEmail       string     `json:"user_email"`
	Status          string     `json:"status"`
	PaymentMethod   string     `json:"payment_method"`
	Total           float64    `json:"total"`
	CreatedAt       time.Time  `json:"created_at"`
	Items           []ItemData `json:"items,omitempty"`
	ShippingName    string     `json:"shipping_name,omitempty"`
	ShippingAddress string     `json:"shipping_address,omitempty"`
	ShippingPhone   string     `json:"shipping_phone,omitempty"`
}

// OrderStatusChangedData represents the data for order.status_changed events
type OrderStatusChangedData struct {
	OrderID        string    `json:"order_id"`
	OrderNumber    string    `json:"order_number"`
	UserID         string    `json:"user_id"`
	UserEmail      string    `json:"user_email"`
	Status         string    `json:"status"`
	PreviousStatus string    `json:"previous_status,omitempty"`
	PaymentMethod  string    `json:"payment_method"`
	Total          float64   `json:"total"`
	CreatedAt      time.Time `json:"created_at"`
}

// PaymentSucceededData represents the data for order.payment_succeeded events
type PaymentSucceededData struct {
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

// ItemData represents an order item
type ItemData struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Legacy Payload format for backward compatibility
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// OTPGeneratedData represents the data for auth.otp_generated events
type OTPGeneratedData struct {
	Email      string `json:"email"`
	OTP        string `json:"otp"`
	Purpose    string `json:"purpose"`
	ExpiresIn  int    `json:"expires_in"`
	Message    string `json:"message"`
	ActionText string `json:"action_text"`
}

// Consumer handles message consumption from RabbitMQ
type Consumer struct {
	conn             *amqp.Connection
	queueName        string
	logger           *log.Logger
	userClient       userpb.UserServiceClient
	userConn         *grpc.ClientConn
	userMutex        sync.Mutex
	messageProcessor *MessageProcessor
}

// Create a new consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	// Get logger configuration from environment variables
	loggerPrefix := os.Getenv("LOGGER_PREFIX")
	if loggerPrefix == "" {
		loggerPrefix = "[LISTENER] "
	}

	logger := log.New(os.Stdout, loggerPrefix, log.LstdFlags)

	consumer := Consumer{
		conn:   conn,
		logger: logger,
		messageProcessor: &MessageProcessor{
			processedMessages: make(map[string]time.Time),
			mutex:             sync.RWMutex{},
			maxAge:            15 * time.Minute,
		},
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
	defer channel.Close()

	return declareExchange(channel)
}

// declareExchange declares the logs_topic exchange
func declareExchange(ch *amqp.Channel) error {
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

// declareRandomQueue declares a shared work queue to prevent duplicate processing
// ANTI-DUPLICATE: Uses fixed queue name so multiple listener instances share the same queue
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"", 
		true,                            // durable - survive broker restart
		false,                           // delete when unused
		false,                           // exclusive - allow multiple consumers
		false,                           // no-wait
		nil,                             // arguments
	)
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

	// Set up Dead Letter Exchange
	err = ch.ExchangeDeclare(
		"dead_letter_exchange", // name
		"fanout",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		consumer.logger.Printf("Warning: Failed to declare dead letter exchange: %v", err)
		// Continue anyway
	}

	// Declare dead letter queue
	dlq, err := ch.QueueDeclare(
		"dead_letter_queue", // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		consumer.logger.Printf("Warning: Failed to declare dead letter queue: %v", err)
		// Continue anyway
	} else {
		// Bind dead letter queue to dead letter exchange
		err = ch.QueueBind(
			dlq.Name,               // queue name
			"",                     // routing key (all messages)
			"dead_letter_exchange", // exchange
			false,                  // no-wait
			nil,                    // arguments
		)
		if err != nil {
			consumer.logger.Printf("Warning: Failed to bind dead letter queue: %v", err)
			// Continue anyway
		}
	}

	// Set up the main queue with the dead letter exchange
	args := amqp.Table{
		"x-dead-letter-exchange": "dead_letter_exchange",
	}

	// bind queue with exchange
	for _, binding_key := range topics {
		err = ch.QueueBind(
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

	// Set QoS (Quality of Service)
	err = ch.Qos(
		1,     // prefetch count (process one message at a time)
		0,     // prefetch size (no specific size limit)
		false, // global (false means apply to just this channel)
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Start consuming messages from the queue
	messages, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer
		false,  // auto-ack (false means manual acknowledgment)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		args,   // arguments
	)
	if err != nil {
		return err
	}

	// Create channel to listen for connection issues
	closeChannel := make(chan *amqp.Error)
	consumer.conn.NotifyClose(closeChannel)

	// make a channel async to hold program is running
	forever := make(chan bool)

	// handle mess recv from queue
	go func() {
		for {
			select {
			case err := <-closeChannel:
				consumer.logger.Printf("RabbitMQ connection closed: %v", err)
				forever <- true
				return
			case msg, ok := <-messages:
				if !ok {
					consumer.logger.Println("RabbitMQ channel closed")
					forever <- true
					return
				}

				// Process the message in a separate goroutine to not block the main loop
				go func(d amqp.Delivery) {
					// Try to process as a standard event first
					var stdEvent StandardEvent
					stdErr := json.Unmarshal(d.Body, &stdEvent)

					if stdErr == nil && stdEvent.Name != "" {
						// Process as a standard event
						eventName := stdEvent.Name
						consumer.logger.Printf("📨 Received standard event: %s", eventName)

						// Special log for user activity events
						if strings.HasPrefix(eventName, "log.INFO.user.") {
							actionType := strings.TrimPrefix(eventName, "log.INFO.user.")
							userID := ""

							// Try to extract user ID from data
							if data, ok := stdEvent.Data.(map[string]interface{}); ok {
								if uid, ok := data["user_id"].(string); ok {
									userID = uid
								}
							}

							if userID != "" {
								consumer.logger.Printf("👤 User activity detected: %s for user ID: %s", actionType, userID)
							} else {
								consumer.logger.Printf("👤 User activity detected: %s", actionType)
							}
						}

						err := consumer.handleStandardEvent(stdEvent, d.RoutingKey)
						if err != nil {
							consumer.logger.Printf("Error processing standard event: %v", err)
							// Nack and don't requeue to avoid infinite loop - goes to DLQ
							d.Nack(false, false)
						} else {
							// Acknowledge the message
							d.Ack(false)
						}
					} else {
						// Try legacy format
						var payload Payload
						legacyErr := json.Unmarshal(d.Body, &payload)

						if legacyErr == nil && payload.Name != "" {
							// Process as legacy event
							consumer.logger.Printf("Received legacy event: %s", payload.Name)
							err := consumer.handleLegacyEvent(payload)
							if err != nil {
								consumer.logger.Printf("Error processing legacy event: %v", err)
								// Nack and don't requeue to avoid infinite loop - goes to DLQ
								d.Nack(false, false)
							} else {
								// Acknowledge the message
								d.Ack(false)
							}
						} else {
							// Failed to parse in any format
							consumer.logger.Printf("Failed to parse message in any format: %v", d.Body)
							// Nack and don't requeue to avoid infinite loop - goes to DLQ
							d.Nack(false, false)
						}
					}
				}(msg)
			}
		}
	}()

	consumer.logger.Printf("Listening for messages on topics: %v", topics)
	// Block forever (or until the program is terminated) to keep listening for messages.
	<-forever

	return nil
}

// handleStandardEvent processes standardized event format based on event name
func (consumer *Consumer) handleStandardEvent(event StandardEvent, routingKey string) error {
	// Log ID của event để dễ dàng theo dõi
	consumer.logger.Printf("Received event: %s, ID: %s, Routing key: %s", event.Name, event.ID, routingKey)

	// Check duplicate
	mp := consumer.getMessageProcessor()
	if mp.IsProcessed(event.ID) {
		consumer.logger.Printf("Duplicate detected: Event %s (ID: %s) already processed, skipping...", event.Name, event.ID)
		return nil
	}

	// Mark processed
	mp.MarkProcessed(event.ID)
	consumer.logger.Printf("Processing: Event %s (ID: %s)", event.Name, event.ID)

	// Ghi log sự kiện vào logger-service
	err := consumer.logStandardEvent(event)
	if err != nil {
		consumer.logger.Printf("Warning: Failed to log event: %v", err)
	}

	// Process based on event name and routing key
	// Ưu tiên kiểm tra routing key trước vì đây là cách RabbitMQ định tuyến tin nhắn
	if routingKey == "auth.otp_generated" || event.Name == "auth.otp_generated" {
		// Send OTP verification email
		consumer.logger.Printf("Sending OTP verification email (from routing key: %s)", routingKey)
		err := consumer.sendOTPEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending OTP email: %v", err)
			return err
		}
		return nil
	}

	// Process based on event name if routing key handling didn't match
	switch event.Name {
	case "log.INFO.user.login_success":
		// Hiển thị thông báo đăng nhập thành công
		if data, ok := event.Data.(map[string]interface{}); ok {
			userID, _ := data["user_id"].(string)
			email, _ := data["email"].(string)
			consumer.logger.Printf("User login: %s (%s)", email, userID)
		}
		consumer.logger.Printf("User activity logged: login_success")
		return nil

	case "log.INFO.user.login_failed":
		// Hiển thị thông báo đăng nhập thất bại
		if data, ok := event.Data.(map[string]interface{}); ok {
			email, _ := data["email"].(string)
			consumer.logger.Printf("Failed login attempt for user: %s", email)
		}
		consumer.logger.Printf("User activity logged: login_failed")
		return nil

	case "log.INFO.user.registered", "log.INFO.user.profile_updated",
		"log.INFO.user.password_changed", "log.INFO.user.password_reset_requested":
		// Ghi log ngắn gọn
		if data, ok := event.Data.(map[string]interface{}); ok {
			action, _ := data["action"].(string)
			consumer.logger.Printf("User activity logged: %s", action)
		} else {
			consumer.logger.Printf("User activity logged: %s", event.Name)
		}
		return nil

	case "log.INFO.user.logout":
		// Special handling for logout events
		if data, ok := event.Data.(map[string]interface{}); ok {
			userID, _ := data["user_id"].(string)
			email, _ := data["email"].(string)
			logoutType := "current session"

			// Check if this was a logout from all devices
			if metadata, ok := data["metadata"].(map[string]interface{}); ok {
				if lt, ok := metadata["logout_type"].(string); ok && lt == "all_devices" {
					logoutType = "all devices"
				}
			}

			consumer.logger.Printf("User %s (%s) logged out from %s", userID, email, logoutType)
		} else {
			consumer.logger.Printf("User logged out (detailed info not available)")
		}
		consumer.logger.Printf("User activity logged: logout")
		return nil

	// Xử lý các event đơn hàng và chuyển đổi thành log events
	case "order.created":
		// Chuyển đổi thành log.INFO.order.created
		err := consumer.handleOrderCreatedEvent(event)
		if err != nil {
			consumer.logger.Printf("Error processing order.created event: %v", err)
			return err
		}
		consumer.logger.Printf("Order created event processed and logged")

		// Send order confirmation email
		consumer.logger.Printf("Sending order confirmation email")
		err = consumer.sendOrderConfirmationEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending order confirmation email: %v", err)
			// Continue even if email fails
		}
		return nil

	case "order.status_changed":
		// Chuyển đổi thành log.INFO.order.status_changed
		err := consumer.handleOrderStatusChangedEvent(event)
		if err != nil {
			consumer.logger.Printf("Error processing order.status_changed event: %v", err)
			return err
		}
		consumer.logger.Printf("Order status changed event processed and logged")

		// Send order status notification email
		consumer.logger.Printf("Sending order status update email")
		err = consumer.sendOrderStatusEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending order status email: %v", err)
			// Continue even if email fails
		}
		return nil

	case "order.payment_succeeded":
		// Chuyển đổi thành log.INFO.order.payment_succeeded
		err := consumer.handlePaymentSucceededEvent(event)
		if err != nil {
			consumer.logger.Printf("Error processing order.payment_succeeded event: %v", err)
			return err
		}
		consumer.logger.Printf("Payment succeeded event processed and logged")

		// Send payment success email
		consumer.logger.Printf("Sending payment success email")
		err = consumer.sendPaymentSuccessEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending payment success email: %v", err)
			// Continue even if email fails
		}
		return nil

	case "order.payment_failed":
		// Chuyển đổi thành log.INFO.order.payment_failed
		err := consumer.handlePaymentFailedEvent(event)
		if err != nil {
			consumer.logger.Printf("Error processing order.payment_failed event: %v", err)
			return err
		}
		consumer.logger.Printf("Payment failed event processed and logged")

		// Có thể gửi email thông báo thanh toán thất bại
		// (triển khai trong tương lai)
		return nil

	case "user.registered":
		// Forward to user-service to store the user
		consumer.logger.Printf("Forwarding new user registration to user-service")
		err := consumer.forwardToUserService(event)
		if err != nil {
			consumer.logger.Printf("Error forwarding to user-service: %v", err)
			return err
		}

		// Send welcome email
		consumer.logger.Printf("Sending welcome email to new user")
		err = consumer.sendWelcomeEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending welcome email: %v", err)
			// Continue even if email fails
		}

	case "auth.password_reset_requested":
		// Send password reset email
		consumer.logger.Printf("Sending password reset email")
		err := consumer.sendPasswordResetEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending password reset email: %v", err)
			return err
		}

	case "user.profile_updated":
		// Convert user.profile_updated to log format and log it
		consumer.logger.Printf("Converting user.profile_updated to log format")
		logEvent := StandardEvent{
			ID:         fmt.Sprintf("log_%s", event.ID),
			Name:       "log.INFO.user.profile_updated",
			Data:       event.Data,
			DataSchema: event.DataSchema,
			Source:     event.Source,
			CreatedAt:  time.Now(),
			Version:    event.Version,
		}

		// Log the converted event
		err := consumer.logStandardEvent(logEvent)
		if err != nil {
			consumer.logger.Printf("Error logging profile updated event: %v", err)
			return err
		}
		consumer.logger.Printf("Profile updated event logged successfully")

	case "auth.password_changed", "user.password_changed":
		// Convert user.password_changed to log format and log it
		if event.Name == "user.password_changed" {
			consumer.logger.Printf("Converting user.password_changed to log format")
			logEvent := StandardEvent{
				ID:         fmt.Sprintf("log_%s", event.ID),
				Name:       "log.INFO.user.password_changed",
				Data:       event.Data,
				DataSchema: event.DataSchema,
				Source:     event.Source,
				CreatedAt:  time.Now(),
				Version:    event.Version,
			}

			// Log the converted event
			err := consumer.logStandardEvent(logEvent)
			if err != nil {
				consumer.logger.Printf("Error logging password changed event: %v", err)
				return err
			}
			consumer.logger.Printf("Password changed event logged successfully")
		}

		// Send password changed notification email
		consumer.logger.Printf("Sending password changed notification email")
		err := consumer.sendPasswordChangedEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending password changed email: %v", err)
			return err
		}

	default:
		// Forward all other events to appropriate service
		consumer.logger.Printf("Forwarding event %s to appropriate service", event.Name)
	}

	return nil
}

// handleOrderCreatedEvent converts order.created event to log format and logs it
func (consumer *Consumer) handleOrderCreatedEvent(event StandardEvent) error {
	// Lấy dữ liệu từ event
	var orderData map[string]interface{}

	// Xử lý data dựa trên kiểu dữ liệu
	switch data := event.Data.(type) {
	case map[string]interface{}:
		orderData = data
	case string:
		if err := json.Unmarshal([]byte(data), &orderData); err != nil {
			return fmt.Errorf("failed to unmarshal order data string: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type for order event")
	}

	// Tạo event mới với định dạng log
	logEvent := StandardEvent{
		ID:         fmt.Sprintf("log_%s", event.ID),
		Name:       "log.INFO.order.created",
		Data:       orderData,
		DataSchema: event.DataSchema,
		Source:     event.Source,
		CreatedAt:  time.Now(),
		Version:    event.Version,
	}

	// Ghi log vào logger-service
	err := consumer.logStandardEvent(logEvent)
	if err != nil {
		return fmt.Errorf("failed to log order created event: %w", err)
	}

	// Hiển thị thông báo về đơn hàng mới
	orderID := ""
	orderNumber := ""
	userEmail := ""
	total := 0.0
	userID := ""

	if id, ok := orderData["order_id"].(string); ok {
		orderID = id
	}
	if num, ok := orderData["order_number"].(string); ok {
		orderNumber = num
	}
	if email, ok := orderData["user_email"].(string); ok {
		userEmail = email
	}
	if t, ok := orderData["total"].(float64); ok {
		total = t
	}
	if uid, ok := orderData["user_id"].(string); ok {
		userID = uid
	}

	consumer.logger.Printf("New order created: #%s (ID: %s) by %s - Total: $%.2f",
		orderNumber, orderID, userEmail, total)

	// Sync user order data
	if userID != "" {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := consumer.updateUserOrderData(ctx, userID); err != nil {
				consumer.logger.Printf("ERROR: Failed to sync user order data: %v", err)
			}
		}()
	}

	return nil
}

// handleOrderStatusChangedEvent converts order.status_changed event to log format and logs it
func (consumer *Consumer) handleOrderStatusChangedEvent(event StandardEvent) error {
	// Lấy dữ liệu từ event
	var orderData map[string]interface{}

	// Xử lý data dựa trên kiểu dữ liệu
	switch data := event.Data.(type) {
	case map[string]interface{}:
		orderData = data
	case string:
		if err := json.Unmarshal([]byte(data), &orderData); err != nil {
			return fmt.Errorf("failed to unmarshal order data string: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type for order event")
	}

	// Tạo event mới với định dạng log
	logEvent := StandardEvent{
		ID:         fmt.Sprintf("log_%s", event.ID),
		Name:       "log.INFO.order.status_changed",
		Data:       orderData,
		DataSchema: event.DataSchema,
		Source:     event.Source,
		CreatedAt:  time.Now(),
		Version:    event.Version,
	}

	// Ghi log vào logger-service
	err := consumer.logStandardEvent(logEvent)
	if err != nil {
		return fmt.Errorf("failed to log order status changed event: %w", err)
	}

	// Hiển thị thông báo về đơn hàng thay đổi trạng thái
	orderID := ""
	orderNumber := ""
	status := ""
	prevStatus := ""

	if id, ok := orderData["order_id"].(string); ok {
		orderID = id
	}
	if num, ok := orderData["order_number"].(string); ok {
		orderNumber = num
	}
	if s, ok := orderData["status"].(string); ok {
		status = s
	}
	if ps, ok := orderData["previous_status"].(string); ok {
		prevStatus = ps
	}

	consumer.logger.Printf("Order status changed: #%s (ID: %s) from '%s' to '%s'",
		orderNumber, orderID, prevStatus, status)

	return nil
}

// handlePaymentSucceededEvent converts order.payment_succeeded event to log format and logs it
func (consumer *Consumer) handlePaymentSucceededEvent(event StandardEvent) error {
	// Lấy dữ liệu từ event
	var paymentData map[string]interface{}

	// Xử lý data dựa trên kiểu dữ liệu
	switch data := event.Data.(type) {
	case map[string]interface{}:
		paymentData = data
	case string:
		if err := json.Unmarshal([]byte(data), &paymentData); err != nil {
			return fmt.Errorf("failed to unmarshal payment data string: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type for payment event")
	}

	// Tạo event mới với định dạng log
	logEvent := StandardEvent{
		ID:         fmt.Sprintf("log_%s", event.ID),
		Name:       "log.INFO.order.payment_succeeded",
		Data:       paymentData,
		DataSchema: event.DataSchema,
		Source:     event.Source,
		CreatedAt:  time.Now(),
		Version:    event.Version,
	}

	// Ghi log vào logger-service
	err := consumer.logStandardEvent(logEvent)
	if err != nil {
		return fmt.Errorf("failed to log payment succeeded event: %w", err)
	}

	// Hiển thị thông báo về thanh toán thành công
	orderID := ""
	orderNumber := ""
	method := ""
	amount := 0.0

	if id, ok := paymentData["order_id"].(string); ok {
		orderID = id
	}
	if num, ok := paymentData["order_number"].(string); ok {
		orderNumber = num
	}
	if m, ok := paymentData["payment_method"].(string); ok {
		method = m
	}
	if a, ok := paymentData["amount"].(float64); ok {
		amount = a
	}

	consumer.logger.Printf("Payment succeeded: #%s (ID: %s) - Method: %s, Amount: $%.2f",
		orderNumber, orderID, method, amount)

	// Update user order data
	if userID, ok := paymentData["user_id"].(string); ok && userID != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := consumer.updateUserOrderData(ctx, userID); err != nil {
			consumer.logger.Printf("Failed to update user order data: %v", err)
			// Continue processing even if updating user data fails
		}
	}

	return nil
}

// handlePaymentFailedEvent converts order.payment_failed event to log format and logs it
func (consumer *Consumer) handlePaymentFailedEvent(event StandardEvent) error {
	// Lấy dữ liệu từ event
	var paymentData map[string]interface{}

	// Xử lý data dựa trên kiểu dữ liệu
	switch data := event.Data.(type) {
	case map[string]interface{}:
		paymentData = data
	case string:
		if err := json.Unmarshal([]byte(data), &paymentData); err != nil {
			return fmt.Errorf("failed to unmarshal payment data string: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type for payment event")
	}

	// Tạo event mới với định dạng log
	logEvent := StandardEvent{
		ID:         fmt.Sprintf("log_%s", event.ID),
		Name:       "log.INFO.order.payment_failed",
		Data:       paymentData,
		DataSchema: event.DataSchema,
		Source:     event.Source,
		CreatedAt:  time.Now(),
		Version:    event.Version,
	}

	// Ghi log vào logger-service
	err := consumer.logStandardEvent(logEvent)
	if err != nil {
		return fmt.Errorf("failed to log payment failed event: %w", err)
	}

	// Hiển thị thông báo về thanh toán thất bại
	orderID := ""
	orderNumber := ""
	method := ""

	if id, ok := paymentData["order_id"].(string); ok {
		orderID = id
	}
	if num, ok := paymentData["order_number"].(string); ok {
		orderNumber = num
	}
	if m, ok := paymentData["payment_method"].(string); ok {
		method = m
	}

	consumer.logger.Printf("Payment failed: #%s (ID: %s) - Method: %s",
		orderNumber, orderID, method)

	return nil
}

// handleLegacyEvent handles events in the legacy format for backward compatibility
func (consumer *Consumer) handleLegacyEvent(payload Payload) error {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			consumer.logger.Println(err)
			return err
		}

	case "user.registered", "user.profile_updated", "auth.password_reset_requested", "user.password_changed", "order.created":
		// Log sự kiện
		err := logEvent(payload)
		if err != nil {
			consumer.logger.Println("Error logging event:", err)
			return err
		}

		// Chuyển đổi và xử lý giống như StandardEvent để đảm bảo tính nhất quán
		// Tạo một StandardEvent từ payload legacy
		data := map[string]interface{}{}
		err = json.Unmarshal([]byte(payload.Data), &data)
		if err != nil {
			consumer.logger.Printf("Error parsing legacy event data: %v", err)
			return err
		}

		// Tạo một StandardEvent từ payload
		stdEvent := StandardEvent{
			ID:        fmt.Sprintf("converted_%s", uuid.New().String()),
			Name:      payload.Name,
			Data:      data,
			CreatedAt: time.Now(),
		}

		// Xử lý thông qua StandardEvent handler để đồng nhất cách xử lý
		return consumer.handleStandardEvent(stdEvent, payload.Name)

	default:
		err := logEvent(payload)
		if err != nil {
			consumer.logger.Println(err)
			return err
		}
	}

	return nil
}

// Log standardized event to logger service
func (consumer *Consumer) logStandardEvent(event StandardEvent) error {
	// Tạo LoggerClient với địa chỉ localhost mặc định
	consumer.logger.Printf("Creating logger client connection...")
	logHost := os.Getenv("LOGGER_SERVICE_HOST")
	if logHost == "" {
		logHost = "logger-service:50056"
	}
	loggerClient, err := NewLoggerClient(logHost)
	if err != nil {
		consumer.logger.Printf("Error creating logger client: %v", err)
		return err
	}
	defer loggerClient.Close()

	// Convert event to JSON
	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event: %w", err)
	}

	// Extract user ID if available for better logging
	userID := "unknown"
	if data, ok := event.Data.(map[string]interface{}); ok {
		if uid, ok := data["user_id"].(string); ok && uid != "" {
			userID = uid
		}
	}

	// Get action type from event name
	actionType := "event"
	if strings.HasPrefix(event.Name, "log.INFO.user.") {
		actionType = strings.TrimPrefix(event.Name, "log.INFO.user.")
	}

	consumer.logger.Printf("📤 Sending log to logger service: %s for user ID: %s", actionType, userID)

	// Ghi log qua gRPC
	ctx := context.Background()
	err = loggerClient.WriteLog(ctx, event.Name, string(jsonData))
	if err != nil {
		consumer.logger.Printf("Error logging event via gRPC: %v", err)
		return err
	}

	consumer.logger.Printf("User activity logged: %s for user ID: %s", actionType, userID)
	return nil
}

// Forward event to user-service for processing
func (consumer *Consumer) forwardToUserService(event StandardEvent) error {
	// Add detailed debug logs
	consumer.logger.Printf("DEBUG: Starting forwardToUserService for event %s (ID: %s)", event.Name, event.ID)

	// Serialize the event data to JSON string
	eventData, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	consumer.logger.Printf("DEBUG: Forwarding user.registered event with full data: %s", string(eventData))

	// Extract and log individual fields from user data
	var userData map[string]interface{}
	if err := json.Unmarshal(eventData, &userData); err == nil {
		consumer.logger.Printf("DEBUG: User data extracted - ID: %v, Email: %v, FirstName: %v, LastName: %v, Username: %v",
			userData["id"], userData["email"], userData["first_name"], userData["last_name"], userData["username"])
	}

	// Add connection timeout with longer duration for reliability
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create gRPC connection with enhanced retry
	var conn *grpc.ClientConn
	var dialErr error

	maxRetries := 5
	connected := false

	userServiceURL := os.Getenv("USER_SERVICE_HOST")
	// Set default URL if environment variable is not set or empty
	if userServiceURL == "" || userServiceURL == " " {
		userServiceURL = "user-service:50052"
	}

	// Try service discovery patterns if explicit URL fails
	for attempt := 0; attempt < maxRetries && !connected; attempt++ {
		// Use the configured URL for all attempts
		currentURL := userServiceURL

		consumer.logger.Printf("DEBUG: Attempting to connect to user service at %s (attempt %d/%d)",
			currentURL, attempt+1, maxRetries)

		conn, dialErr = grpc.DialContext(ctx, currentURL,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock())

		if dialErr == nil {
			connected = true
			consumer.logger.Printf("DEBUG: Successfully connected to user service at %s", currentURL)
			break
		}

		consumer.logger.Printf("WARNING: Failed to connect to user service at %s (attempt %d/%d): %v",
			currentURL, attempt+1, maxRetries, dialErr)

		if attempt < maxRetries-1 {
			// Exponential backoff
			backoffDuration := time.Duration(1<<uint(attempt+1)) * 100 * time.Millisecond
			if backoffDuration > 2*time.Second {
				backoffDuration = 2 * time.Second // Cap maximum backoff
			}
			consumer.logger.Printf("Retrying in %v...", backoffDuration)
			time.Sleep(backoffDuration)
		}
	}

	if !connected {
		return fmt.Errorf("failed to connect to user service after %d attempts: %w", maxRetries, dialErr)
	}

	defer conn.Close()

	// Create gRPC client
	client := userpb.NewUserServiceClient(conn)

	// Create context with timeout for the RPC call
	callCtx, callCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer callCancel()

	// Create event request
	eventRequest := &userpb.EventRequest{
		EventId:   event.ID,
		EventName: event.Name,
		EventData: string(eventData),
		Source:    event.Source,
		CreatedAt: event.CreatedAt.Format(time.RFC3339),
	}

	// consumer.logger.Printf("DEBUG: Calling ProcessEvent RPC with request: %+v", eventRequest)
	// consumer.logger.Printf("DEBUG: Full EventData being sent: %s", eventRequest.EventData)

	// Call the ProcessEvent RPC
	response, err := client.ProcessEvent(callCtx, eventRequest)
	if err != nil {
		// consumer.logger.Printf("ERROR: Failed calling user service ProcessEvent: %v", err)
		return fmt.Errorf("error calling user service ProcessEvent: %w", err)
	}

	if !response.Success {
		// consumer.logger.Printf("ERROR: User service reported failure: %s", response.Message)
		return fmt.Errorf("user service reported failure: %s", response.Message)
	}

	// consumer.logger.Printf("EVENT SUCCESS: Event %s processed by user-service: %s", event.Name, response.Message)
	return nil
}

// logEvent logs a legacy event to the logger service
func logEvent(entry Payload) error {
	// Tạo LoggerClient với địa chỉ localhost mặc định
	logHost := os.Getenv("LOGGER_SERVICE_HOST")
	if logHost == "" {
		logHost = "logger-service:50056"
	}
	loggerClient, err := NewLoggerClient(logHost)
	if err != nil {
		log.Printf("Error creating logger client: %v", err)
		return err
	}
	defer loggerClient.Close()

	// Convert payload to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	// Ghi log qua gRPC
	ctx := context.Background()
	err = loggerClient.WriteLog(ctx, entry.Name, string(jsonData))
	if err != nil {
		log.Printf("Error logging event via gRPC: %v", err)
		return err
	}

	log.Printf("Legacy event logged: %s", entry.Name)
	return nil
}

// Send welcome email for new user registration
func (consumer *Consumer) sendWelcomeEmail(event StandardEvent) error {
	// Extract user data
	var userData UserRegistered

	// Handle different data formats
	switch data := event.Data.(type) {
	case map[string]interface{}:
		// Convert map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal user data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &userData); err != nil {
			return fmt.Errorf("failed to unmarshal user data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &userData); err != nil {
			return fmt.Errorf("failed to unmarshal user data string: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type for user registration event")
	}

	// Create a mail client
	mailClient, err := NewMailClient("")
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}
	defer mailClient.Close()

	// Send welcome email
	err = mailClient.SendWelcomeEmail(userData.Email, userData.FirstName, userData.LastName)
	if err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	consumer.logger.Printf("Welcome email sent to %s", userData.Email)
	return nil
}

// sendPasswordResetEmail sends a password reset email based on an event
func (consumer *Consumer) sendPasswordResetEmail(event StandardEvent) error {
	// Kiểm tra ID sự kiện đã được xử lý chưa
	eventID := event.ID
	if eventID == "" {
		consumer.logger.Printf("Warning: Password reset event has no ID, will process anyway")
	} else {
		// Kiểm tra xem sự kiện này đã được xử lý chưa
		processor := consumer.getMessageProcessor()
		if processor != nil && processor.IsProcessed(eventID) {
			consumer.logger.Printf("Skipping duplicate password reset email event with ID: %s", eventID)
			return nil
		}
	}

	// Extract password reset data
	var resetData PasswordResetRequested

	// Handle different data formats
	switch data := event.Data.(type) {
	case map[string]interface{}:
		// Convert map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal password reset data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &resetData); err != nil {
			return fmt.Errorf("failed to unmarshal password reset data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &resetData); err != nil {
			return fmt.Errorf("failed to unmarshal password reset data string: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type for password reset event")
	}

	// Create a mail client
	mailClient, err := NewMailClient("")
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}
	defer mailClient.Close()

	// Send password reset email
	err = mailClient.SendPasswordResetEmail(resetData.Email, resetData.TokenHash, resetData.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	// Đánh dấu sự kiện đã được xử lý
	if eventID != "" && consumer.getMessageProcessor() != nil {
		consumer.getMessageProcessor().MarkProcessed(eventID)
	}

	consumer.logger.Printf("Password reset email sent to %s with event ID: %s", resetData.Email, eventID)
	return nil
}

// Send password changed notification email
func (consumer *Consumer) sendPasswordChangedEmail(event StandardEvent) error {
	// Extract password changed data
	var changedData PasswordChanged

	// Handle different data formats
	switch data := event.Data.(type) {
	case map[string]interface{}:
		// Convert map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal password changed data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &changedData); err != nil {
			return fmt.Errorf("failed to unmarshal password changed data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &changedData); err != nil {
			return fmt.Errorf("failed to unmarshal password changed data string: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type for password changed event")
	}

	// Create a mail client
	mailClient, err := NewMailClient("")
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}
	defer mailClient.Close()

	// Send password changed email
	err = mailClient.SendPasswordChangedEmail(changedData.Email, changedData.ChangedAt)
	if err != nil {
		return fmt.Errorf("failed to send password changed email: %w", err)
	}

	consumer.logger.Printf("Password changed email sent to %s", changedData.Email)
	return nil
}

// Send order confirmation email
func (consumer *Consumer) sendOrderConfirmationEmail(event StandardEvent) error {
	consumer.logger.Printf("Processing order confirmation event")

	// Extract order data
	var orderData OrderCreatedData

	// Handle different data formats
	switch data := event.Data.(type) {
	case map[string]interface{}:
		// Convert map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal order data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &orderData); err != nil {
			return fmt.Errorf("failed to unmarshal order data: %w", err)
		}

		// Get shipping information directly from data
		if name, ok := data["shipping_name"].(string); ok && name != "" {
			orderData.ShippingName = name
		}
		if address, ok := data["shipping_address"].(string); ok && address != "" {
			orderData.ShippingAddress = address
		}
		if phone, ok := data["shipping_phone"].(string); ok && phone != "" {
			orderData.ShippingPhone = phone
		}

		// Legacy code for backward compatibility
		if shippingData, ok := data["shipping"].(map[string]interface{}); ok {
			if name, ok := shippingData["name"].(string); ok && orderData.ShippingName == "" {
				orderData.ShippingName = name
			}
			if address, ok := shippingData["address"].(string); ok && orderData.ShippingAddress == "" {
				orderData.ShippingAddress = address
			}
			if phone, ok := shippingData["phone"].(string); ok && orderData.ShippingPhone == "" {
				orderData.ShippingPhone = phone
			}
		}
	case string:
		if err := json.Unmarshal([]byte(data), &orderData); err != nil {
			return fmt.Errorf("failed to unmarshal order data string: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type for order data")
	}

	// Validate required fields
	if orderData.UserEmail == "" {
		return fmt.Errorf("missing email in order data")
	}

	// Format items for display in email
	itemsStr := ""
	if len(orderData.Items) > 0 {
		for i, item := range orderData.Items {
			subtotal := float64(item.Quantity) * item.Price
			itemsStr += fmt.Sprintf("%s\n%d × $%.2f = $%.2f",
				item.Name,
				item.Quantity,
				item.Price,
				subtotal)

			// Add a separator between items except for the last one
			if i < len(orderData.Items)-1 {
				itemsStr += "\n\n"
			}
		}
	} else {
		itemsStr = "Your order items will be displayed here."
	}

	// Format total for display
	totalStr := fmt.Sprintf("%.2f", orderData.Total)

	// Format order date in a user-friendly way
	orderDate := orderData.CreatedAt.Format("Monday, January 2, 2006")

	// Prepare payment method display
	paymentMethod := orderData.PaymentMethod
	if paymentMethod == "cod" || paymentMethod == "COD" {
		paymentMethod = "Cash on Delivery (COD)"
	} else if paymentMethod == "card" {
		paymentMethod = "Credit/Debit Card"
	} else if paymentMethod == "bank_transfer" {
		paymentMethod = "Bank Transfer"
	}

	// Extract customer name from orderData if possible
	customerName := "Valued Customer"
	if orderData.ShippingName != "" {
		customerName = orderData.ShippingName
	} else {
		// If user ID is available, try to get user information from email
		if orderData.UserEmail != "" {
			parts := strings.Split(orderData.UserEmail, "@")
			if len(parts) > 0 && parts[0] != "" {
				customerName = strings.Title(parts[0])
			}
		}
	}

	// Create a mail client
	mailClient, err := NewMailClient("")
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}
	defer mailClient.Close()

	// Send order confirmation email
	err = mailClient.SendOrderConfirmationEmail(
		orderData.UserEmail,
		orderData.OrderID,
		orderData.OrderNumber,
		customerName,
		orderDate,
		orderData.Status,
		paymentMethod,
		totalStr,
		itemsStr,
		orderData.ShippingAddress,
		orderData.ShippingName,
		orderData.ShippingPhone,
	)
	if err != nil {
		return fmt.Errorf("failed to send order confirmation email: %w", err)
	}

	consumer.logger.Printf("Order confirmation email sent to %s", orderData.UserEmail)
	return nil
}

// sendOrderStatusEmail sends order status notification email
func (consumer *Consumer) sendOrderStatusEmail(event StandardEvent) error {
	consumer.logger.Printf("Processing order status change event")

	// Extract order status data
	var statusData OrderStatusChangedData

	// Handle different data formats
	switch data := event.Data.(type) {
	case map[string]interface{}:
		// Convert map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal order status data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &statusData); err != nil {
			return fmt.Errorf("failed to unmarshal order status data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &statusData); err != nil {
			return fmt.Errorf("failed to unmarshal order status data string: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type for order status data")
	}

	// Validate required fields
	if statusData.UserEmail == "" {
		return fmt.Errorf("missing email in order status data")
	}

	// Get appropriate subject based on status
	var subject string
	var template string

	switch statusData.Status {
	case "processing":
		subject = "Your Order is Being Processed"
		template = "order_processing.html.gohtml"
	case "shipped":
		subject = "Your Order Has Been Shipped"
		template = "order_shipped.html.gohtml"
	case "delivered":
		subject = "Your Order Has Been Delivered"
		template = "order_delivered.html.gohtml"
	case "cancelled":
		subject = "Your Order Has Been Cancelled"
		template = "order_cancelled.html.gohtml"
	default:
		subject = "Order Status Update: " + statusData.Status
		template = "order_status_change.html.gohtml"
	}

	// Format total for display
	totalStr := fmt.Sprintf("%.2f", statusData.Total)

	// Create variables map
	variables := map[string]string{
		"order_id":        statusData.OrderID,
		"order_number":    statusData.OrderNumber,
		"status":          statusData.Status,
		"previous_status": statusData.PreviousStatus,
		"total":           totalStr,
		"year":            fmt.Sprintf("%d", time.Now().Year()),
	}

	// Create a mail client
	mailClient, err := NewMailClient("")
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}
	defer mailClient.Close()

	// Send order status email
	err = mailClient.SendTemplateEmail(
		statusData.UserEmail,
		subject+" - Order #"+statusData.OrderNumber,
		template,
		variables,
	)
	if err != nil {
		return fmt.Errorf("failed to send order status email: %w", err)
	}

	consumer.logger.Printf("Order status email sent to %s", statusData.UserEmail)
	return nil
}

// sendPaymentSuccessEmail sends payment success notification email
func (consumer *Consumer) sendPaymentSuccessEmail(event StandardEvent) error {
	consumer.logger.Printf("Processing payment success event")

	// Extract payment data
	var paymentData PaymentSucceededData

	// Handle different data formats
	switch data := event.Data.(type) {
	case map[string]interface{}:
		// Convert map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal payment data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &paymentData); err != nil {
			return fmt.Errorf("failed to unmarshal payment data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &paymentData); err != nil {
			return fmt.Errorf("failed to unmarshal payment data string: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type for payment data")
	}

	// Validate required fields
	if paymentData.UserEmail == "" {
		return fmt.Errorf("missing email in payment data")
	}

	// Format currency amount for display
	amountStr := fmt.Sprintf("%.2f", paymentData.Amount)

	// Format payment date in a user-friendly way
	formattedPaymentDate := paymentData.PaymentDate
	// Try to parse the date if it's in a standard format and reformat it
	if parsedDate, err := time.Parse(time.RFC3339, paymentData.PaymentDate); err == nil {
		formattedPaymentDate = parsedDate.Format("Monday, January 2, 2006")
	} else if parsedDate, err := time.Parse("2006-01-02T15:04:05Z", paymentData.PaymentDate); err == nil {
		formattedPaymentDate = parsedDate.Format("Monday, January 2, 2006")
	}

	// Prepare payment method display
	paymentMethod := paymentData.PaymentMethod
	if paymentMethod == "card" {
		paymentMethod = "Credit/Debit Card"
	} else if paymentMethod == "paypal" {
		paymentMethod = "PayPal"
	} else if paymentMethod == "bank_transfer" {
		paymentMethod = "Bank Transfer"
	}

	// Prepare currency display
	currency := paymentData.Currency
	if currency == "USD" {
		currency = "USD ($)"
	} else if currency == "EUR" {
		currency = "EUR (€)"
	} else if currency == "GBP" {
		currency = "GBP (£)"
	} else if currency == "VND" {
		currency = "VND (₫)"
	}

	// Create variables map
	variables := map[string]string{
		"order_id":       paymentData.OrderID,
		"order_number":   paymentData.OrderNumber,
		"amount":         amountStr,
		"currency":       currency,
		"payment_method": paymentMethod,
		"transaction_id": paymentData.TransactionID,
		"payment_date":   formattedPaymentDate,
		"payment_status": paymentData.Status,
		"year":           fmt.Sprintf("%d", time.Now().Year()),
		"user_id":        paymentData.UserID,
		"user_email":     paymentData.UserEmail,
	}

	// Create a mail client
	mailClient, err := NewMailClient("")
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}
	defer mailClient.Close()

	// Send payment success email
	err = mailClient.SendTemplateEmail(
		paymentData.UserEmail,
		"Payment Confirmation for Order #"+paymentData.OrderNumber,
		"payment_success.html.gohtml",
		variables,
	)
	if err != nil {
		return fmt.Errorf("failed to send payment success email: %w", err)
	}

	consumer.logger.Printf("Payment success email sent to %s", paymentData.UserEmail)
	return nil
}

// Send OTP verification email
func (consumer *Consumer) sendOTPEmail(event StandardEvent) error {
	// Kiểm tra ID sự kiện đã được xử lý chưa
	eventID := event.ID
	if eventID == "" {
		consumer.logger.Printf("Warning: OTP event has no ID, will process anyway")
	} else {
		// Kiểm tra xem sự kiện này đã được xử lý chưa
		// Lưu ý: Đây là triển khai đơn giản, trong thực tế nên dùng Redis hoặc DB để lưu trữ ID đã xử lý
		processor := consumer.getMessageProcessor()
		if processor != nil && processor.IsProcessed(eventID) {
			consumer.logger.Printf("Skipping duplicate OTP email event with ID: %s", eventID)
			return nil
		}
	}

	// Extract OTP data
	var otpData OTPGeneratedData

	// Handle different data formats
	switch data := event.Data.(type) {
	case map[string]interface{}:
		// Convert map to JSON then to struct
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal OTP data: %w", err)
		}
		if err := json.Unmarshal(jsonData, &otpData); err != nil {
			return fmt.Errorf("failed to unmarshal OTP data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &otpData); err != nil {
			return fmt.Errorf("failed to unmarshal OTP data string: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type for OTP event")
	}

	// Create a mail client
	mailClient, err := NewMailClient("")
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}
	defer mailClient.Close()

	// Send OTP email
	err = mailClient.SendOTPEmail(
		otpData.Email,
		otpData.OTP,
		otpData.Purpose,
		int32(otpData.ExpiresIn),
		otpData.Message,
		otpData.ActionText,
	)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	// Đánh dấu sự kiện đã được xử lý
	if eventID != "" && consumer.getMessageProcessor() != nil {
		consumer.getMessageProcessor().MarkProcessed(eventID)
	}

	consumer.logger.Printf("OTP email sent to %s with event ID: %s", otpData.Email, eventID)
	return nil
}

// MessageProcessor xử lý theo dõi tin nhắn đã được xử lý
type MessageProcessor struct {
	processedMessages map[string]time.Time
	mutex             sync.RWMutex
	maxAge            time.Duration
}

// getMessageProcessor trả về message processor hoặc tạo mới nếu chưa có
func (consumer *Consumer) getMessageProcessor() *MessageProcessor {
	// Return the message processor stored in the consumer
	if consumer.messageProcessor == nil {
		consumer.messageProcessor = &MessageProcessor{
			processedMessages: make(map[string]time.Time),
			mutex:             sync.RWMutex{},
			maxAge:            15 * time.Minute,
		}
	}
	return consumer.messageProcessor
}

// IsProcessed kiểm tra xem một tin nhắn đã được xử lý chưa
func (mp *MessageProcessor) IsProcessed(messageID string) bool {
	mp.mutex.RLock()
	defer mp.mutex.RUnlock()

	processTime, exists := mp.processedMessages[messageID]
	if !exists {
		return false
	}

	// Nếu tin nhắn đã quá cũ, xem như chưa xử lý
	if time.Since(processTime) > mp.maxAge {
		delete(mp.processedMessages, messageID)
		return false
	}

	return true
}

// MarkProcessed đánh dấu tin nhắn đã được xử lý
func (mp *MessageProcessor) MarkProcessed(messageID string) {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()

	mp.processedMessages[messageID] = time.Now()

	// Dọn dẹp các tin nhắn cũ
	for id, t := range mp.processedMessages {
		if time.Since(t) > mp.maxAge {
			delete(mp.processedMessages, id)
		}
	}
}

// updateUserOrderData updates the user's order count and total spend in the user service
func (consumer *Consumer) updateUserOrderData(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required to update order data")
	}

	// Connect to user service if not already connected
	if err := consumer.connectToUserService(); err != nil {
		return fmt.Errorf("failed to connect to user service: %w", err)
	}

	// Get order count and total spend for the user from checkout service
	// This would typically involve a call to the checkout service
	// For now, we'll use a simple query to the order repository

	// Since we don't have direct access to checkout repository here,
	// we'll need to set up a dedicated endpoint or event to get this data

	// Placeholder values - in production, these would come from checkout service
	orderCount := 0
	totalSpend := 0.0

	// Call the user service to update the user order data
	_, err := consumer.userClient.SyncUserOrderData(ctx, &userpb.SyncUserOrderDataRequest{
		UserId:     userID,
		OrderCount: int32(orderCount),
		TotalSpend: totalSpend,
	})

	if err != nil {
		consumer.logger.Printf("Error updating user order data: %v", err)
		return fmt.Errorf("failed to sync user order data: %w", err)
	}

	consumer.logger.Printf("Updated order data for user %s: order count=%d, total spend=%.2f",
		userID, orderCount, totalSpend)
	return nil
}

// connectToUserService establishes a connection to the user service
func (consumer *Consumer) connectToUserService() error {
	consumer.userMutex.Lock()
	defer consumer.userMutex.Unlock()

	// If we already have a connection, return nil
	if consumer.userClient != nil {
		return nil
	}

	// Get user service host from environment or use default
	userHost := os.Getenv("USER_SERVICE_HOST")
	if userHost == "" {
		userHost = "user-service:50052"
	}

	// consumer.logger.Printf("DEBUG: Attempting to connect to user service at %s", userHost)

	// Set up connection to user service with retry logic
	maxRetries := 5
	var dialErr error
	var conn *grpc.ClientConn

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Set up connection with a timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		conn, dialErr = grpc.DialContext(
			ctx,
			userHost,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		cancel()

		// if dialErr == nil {
		// 	// consumer.logger.Printf("DEBUG: Successfully connected to user service at %s", userHost)
		// 	break
		// }

		// consumer.logger.Printf("WARNING: Failed to connect to user service at %s (attempt %d/%d): %v",
		// 	userHost, attempt, maxRetries, dialErr)

		// Wait before retrying
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	if dialErr != nil {
		return fmt.Errorf("failed to connect to user service after %d attempts: %w", maxRetries, dialErr)
	}

	// Create user service client
	consumer.userConn = conn
	consumer.userClient = userpb.NewUserServiceClient(conn)

	return nil
}

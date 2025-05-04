package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "listener-service/proto/user"
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
	conn      *amqp.Connection
	queueName string
	logger    *log.Logger
}

// Create a new consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	logger := log.New(os.Stdout, "[LISTENER] ", log.LstdFlags)

	consumer := Consumer{
		conn:   conn,
		logger: logger,
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

// declareRandomQueue declares a random queue
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name - empty for random
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
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
							consumer.logger.Printf("❌ Error processing standard event: %v", err)
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
	// Log dạng đơn giản khi nhận được event
	consumer.logger.Printf("📨 Received event: %s", event.Name)

	// Ghi log sự kiện vào logger-service
	err := consumer.logStandardEvent(event)
	if err != nil {
		consumer.logger.Printf("Warning: Failed to log event: %v", err)
		// Continue processing even if logging fails
	}

	// Process based on event name
	switch event.Name {
	case "log.INFO.user.login_success":
		// Hiển thị thông báo đăng nhập thành công
		if data, ok := event.Data.(map[string]interface{}); ok {
			userID, _ := data["user_id"].(string)
			email, _ := data["email"].(string)
			consumer.logger.Printf("👤 User login: %s (%s)", email, userID)
		}
		consumer.logger.Printf("✅ User activity logged: login_success")
		return nil

	case "log.INFO.user.login_failed":
		// Hiển thị thông báo đăng nhập thất bại
		if data, ok := event.Data.(map[string]interface{}); ok {
			email, _ := data["email"].(string)
			consumer.logger.Printf("🚫 Failed login attempt for user: %s", email)
		}
		consumer.logger.Printf("✅ User activity logged: login_failed")
		return nil

	case "log.INFO.user.registered", "log.INFO.user.profile_updated",
		"log.INFO.user.password_changed", "log.INFO.user.password_reset_requested":
		// Ghi log ngắn gọn
		if data, ok := event.Data.(map[string]interface{}); ok {
			action, _ := data["action"].(string)
			consumer.logger.Printf("✅ User activity logged: %s", action)
		} else {
			consumer.logger.Printf("✅ User activity logged: %s", event.Name)
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

			consumer.logger.Printf("🔒 User %s (%s) logged out from %s", userID, email, logoutType)
		} else {
			consumer.logger.Printf("🔒 User logged out (detailed info not available)")
		}
		consumer.logger.Printf("✅ User activity logged: logout")
		return nil

	case "user.registered":
		// Forward to user-service to store the user
		consumer.logger.Printf("👤 Forwarding new user registration to user-service")
		err := consumer.forwardToUserService(event)
		if err != nil {
			consumer.logger.Printf("Error forwarding to user-service: %v", err)
			return err
		}

		// Send welcome email
		consumer.logger.Printf("📧 Sending welcome email to new user")
		err = consumer.sendWelcomeEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending welcome email: %v", err)
			// Continue even if email fails
		}

	case "auth.password_reset_requested":
		// Send password reset email
		consumer.logger.Printf("📧 Sending password reset email")
		err := consumer.sendPasswordResetEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending password reset email: %v", err)
			return err
		}

	case "auth.password_changed":
		// Send password changed notification email
		consumer.logger.Printf("📧 Sending password changed notification email")
		err := consumer.sendPasswordChangedEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending password changed email: %v", err)
			return err
		}

	case "auth.otp_generated":
		// Send OTP verification email
		consumer.logger.Printf("📧 Sending OTP verification email")
		err := consumer.sendOTPEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending OTP email: %v", err)
			return err
		}

	case "order.created":
		// Send order confirmation email
		consumer.logger.Printf("📧 Sending order confirmation email")
		err := consumer.sendOrderConfirmationEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending order confirmation email: %v", err)
			// Continue even if email fails
		}

	case "order.payment_succeeded":
		// Send payment success email
		consumer.logger.Printf("📧 Sending payment success email")
		err := consumer.sendPaymentSuccessEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending payment success email: %v", err)
			// Continue even if email fails
		}

	case "order.status_changed":
		// Send order status notification email
		consumer.logger.Printf("📧 Sending order status update email")
		err := consumer.sendOrderStatusEmail(event)
		if err != nil {
			consumer.logger.Printf("Error sending order status email: %v", err)
			// Continue even if email fails
		}

	default:
		// Forward all other events to appropriate service
		consumer.logger.Printf("🔄 Forwarding event %s to appropriate service", event.Name)
	}

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

	case "user.registered":
		// Handle user registration event - send welcome email
		// First, log the event
		err := logEvent(payload)
		if err != nil {
			consumer.logger.Println("Error logging registration event:", err)
			return err
		}

		// Then publish email.send event
		err = publishEmailEvent("registration", payload.Data)
		if err != nil {
			consumer.logger.Println("Error publishing registration email event:", err)
			return err
		}

	case "auth.password_reset_requested":
		// Handle password reset request - send password reset email
		err := logEvent(payload)
		if err != nil {
			consumer.logger.Println("Error logging password reset request event:", err)
			return err
		}

		// Publish email.send event
		err = publishEmailEvent("reset_password", payload.Data)
		if err != nil {
			consumer.logger.Println("Error publishing password reset email event:", err)
			return err
		}

	case "user.password_changed":
		// Handle password change - send notification email
		err := logEvent(payload)
		if err != nil {
			consumer.logger.Println("Error logging password change event:", err)
			return err
		}

		// Publish email.send event
		err = publishEmailEvent("password_change", payload.Data)
		if err != nil {
			consumer.logger.Println("Error publishing password change email event:", err)
			return err
		}

	case "order.created":
		// Handle order creation - send order confirmation
		err := logEvent(payload)
		if err != nil {
			consumer.logger.Println("Error logging order creation event:", err)
			return err
		}

		// Publish email.send event
		err = publishEmailEvent("order_confirmation", payload.Data)
		if err != nil {
			consumer.logger.Println("Error publishing order confirmation email event:", err)
			return err
		}

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
	consumer.logger.Printf("🔄 Creating logger client connection...")
	loggerClient, err := NewLoggerClient("localhost:50001")
	if err != nil {
		consumer.logger.Printf("❌ Error creating logger client: %v", err)
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
		consumer.logger.Printf("❌ Error logging event via gRPC: %v", err)
		return err
	}

	consumer.logger.Printf("✅ User activity logged: %s for user ID: %s", actionType, userID)
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

	// Set up gRPC connection to user service
	userServiceURL := "localhost:50052" // Default for docker environment

	// Try multiple service discovery patterns
	userServiceOptions := []string{
		"user-service:50052",
		"localhost:50052",
		"user-service.default.svc.cluster.local:50052",
	}

	// Allow override from environment
	if envURL := os.Getenv("USER_SERVICE_URL"); envURL != "" {
		userServiceURL = envURL
		consumer.logger.Printf("DEBUG: Using USER_SERVICE_URL from environment: %s", userServiceURL)
	} else if os.Getenv("USE_LOCAL_SERVICES") == "true" || os.Getenv("LOCAL_DEVELOPMENT") == "true" {
		userServiceURL = "localhost:50052"
		consumer.logger.Printf("DEBUG: Using localhost for user-service due to local environment")
	}

	// Add connection timeout with longer duration for reliability
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create gRPC connection with enhanced retry
	var conn *grpc.ClientConn
	var dialErr error

	maxRetries := 5
	connected := false

	// Try service discovery patterns if explicit URL fails
	for attempt := 0; attempt < maxRetries && !connected; attempt++ {
		// On first attempt, use the configured URL
		// On subsequent attempts, try other service discovery patterns
		currentURL := userServiceURL
		if attempt > 0 && attempt-1 < len(userServiceOptions) {
			currentURL = userServiceOptions[attempt-1]
			consumer.logger.Printf("DEBUG: Retrying with alternative service URL: %s (attempt %d/%d)",
				currentURL, attempt+1, maxRetries)
		}

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
	client := pb.NewUserServiceClient(conn)

	// Create context with timeout for the RPC call
	callCtx, callCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer callCancel()

	// Create event request
	eventRequest := &pb.EventRequest{
		EventId:   event.ID,
		EventName: event.Name,
		EventData: string(eventData),
		Source:    event.Source,
		CreatedAt: event.CreatedAt.Format(time.RFC3339),
	}

	consumer.logger.Printf("DEBUG: Calling ProcessEvent RPC with request: %+v", eventRequest)
	consumer.logger.Printf("DEBUG: Full EventData being sent: %s", eventRequest.EventData)

	// Call the ProcessEvent RPC
	response, err := client.ProcessEvent(callCtx, eventRequest)
	if err != nil {
		consumer.logger.Printf("ERROR: Failed calling user service ProcessEvent: %v", err)
		return fmt.Errorf("error calling user service ProcessEvent: %w", err)
	}

	if !response.Success {
		consumer.logger.Printf("ERROR: User service reported failure: %s", response.Message)
		return fmt.Errorf("user service reported failure: %s", response.Message)
	}

	consumer.logger.Printf("EVENT SUCCESS: Event %s processed by user-service: %s", event.Name, response.Message)
	return nil
}

// logEvent logs a legacy event to the logger service
func logEvent(entry Payload) error {
	// Tạo LoggerClient với địa chỉ localhost mặc định
	loggerClient, err := NewLoggerClient("localhost:50001")
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

// publishEmailEvent publishes an email.send event to RabbitMQ
func publishEmailEvent(emailType, data string) error {
	// Parse the original data to extract email and other information
	var originalData map[string]interface{}
	err := json.Unmarshal([]byte(data), &originalData)
	if err != nil {
		return fmt.Errorf("error unmarshalling original data: %w", err)
	}

	// Ensure we have an email address
	email, ok := originalData["email"].(string)
	if !ok {
		return fmt.Errorf("email address not found in event data")
	}

	// Create email payload
	emailPayload := map[string]interface{}{
		"type":      emailType,
		"to":        email,
		"from":      "", // Will use default from mail-service
		"from_name": "", // Will use default from mail-service
	}

	// Set subject based on email type
	switch emailType {
	case "registration":
		emailPayload["subject"] = "Welcome to our platform!"
		// Add any registration-specific data
		if name, ok := originalData["name"].(string); ok {
			emailPayload["name"] = name
		}
		if username, ok := originalData["username"].(string); ok {
			emailPayload["username"] = username
		}
		if verificationToken, ok := originalData["verification_token"].(string); ok {
			emailPayload["token"] = verificationToken
		}

	case "reset_password":
		emailPayload["subject"] = "Password Reset Request"
		// Add reset token if available
		if resetToken, ok := originalData["reset_token"].(string); ok {
			emailPayload["token"] = resetToken
		}

	case "password_change":
		emailPayload["subject"] = "Your Password Has Been Changed"
		// No additional data needed

	case "order_confirmation":
		emailPayload["subject"] = "Order Confirmation"
		// Add order details if available
		if orderID, ok := originalData["order_id"].(string); ok {
			emailPayload["order_id"] = orderID
		}
		if total, ok := originalData["total"].(float64); ok {
			emailPayload["total"] = total
		}
		// Add items if available
		if items, ok := originalData["items"].([]interface{}); ok {
			emailPayload["items"] = items
		}

	default:
		emailPayload["subject"] = "Notification from Our Platform"
	}

	// Convert to JSON for RabbitMQ
	jsonData, err := json.Marshal(map[string]interface{}{
		"name": "email.send",
		"data": emailPayload,
	})
	if err != nil {
		return fmt.Errorf("error marshalling email payload: %w", err)
	}

	// Determine RabbitMQ connection URL
	rabbitURL := "amqp://guest:guest@localhost:5672"

	if os.Getenv("RABBITMQ_URL") != "" {
		rabbitURL = os.Getenv("RABBITMQ_URL")
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return fmt.Errorf("error connecting to RabbitMQ: %w", err)
	}
	defer conn.Close()

	// Create channel
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("error creating RabbitMQ channel: %w", err)
	}
	defer ch.Close()

	// Publish to RabbitMQ
	err = ch.Publish(
		"logs_topic", // exchange
		"email.send", // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		return fmt.Errorf("error publishing to RabbitMQ: %w", err)
	}

	return nil
}

// Implementation of additional methods for completeness

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

	// Create email data
	emailData := EmailData{
		Type:     "registration",
		To:       userData.Email,
		Subject:  "Welcome to our platform!",
		Template: "welcome.html.gohtml",
		Variables: map[string]string{
			"first_name": userData.FirstName,
			"last_name":  userData.LastName,
			"email":      userData.Email,
		},
	}

	// Create standard email event
	emailEvent := StandardEvent{
		ID:         fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		Name:       "email.send",
		Data:       emailData,
		DataSchema: "v1",
		Source:     "listener-service",
		CreatedAt:  time.Now(),
		Version:    "v1",
	}

	// Forward to mail service
	return consumer.publishEmailEvent(emailEvent)
}

// Send password reset email
func (consumer *Consumer) sendPasswordResetEmail(event StandardEvent) error {
	// Extract reset data
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

	// Create email data
	emailData := EmailData{
		Type:     "reset_password",
		To:       resetData.Email,
		Subject:  "Password Reset Request",
		Template: "password_reset.html.gohtml",
		Variables: map[string]string{
			"email":      resetData.Email,
			"token_hash": resetData.TokenHash,
			"expires_at": resetData.ExpiresAt,
			"year":       fmt.Sprintf("%d", time.Now().Year()),
		},
	}

	// Create standard email event
	emailEvent := StandardEvent{
		ID:         fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		Name:       "email.send",
		Data:       emailData,
		DataSchema: "v1",
		Source:     "listener-service",
		CreatedAt:  time.Now(),
		Version:    "v1",
	}

	// Forward to mail service
	return consumer.publishEmailEvent(emailEvent)
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

	// Create email data
	emailData := EmailData{
		Type:     "password_change",
		To:       changedData.Email,
		Subject:  "Your Password Has Been Changed",
		Template: "password_changed.html.gohtml",
		Variables: map[string]string{
			"email":      changedData.Email,
			"changed_at": changedData.ChangedAt,
			"year":       fmt.Sprintf("%d", time.Now().Year()),
		},
	}

	// Create standard email event
	emailEvent := StandardEvent{
		ID:         fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		Name:       "email.send",
		Data:       emailData,
		DataSchema: "v1",
		Source:     "listener-service",
		CreatedAt:  time.Now(),
		Version:    "v1",
	}

	// Forward to mail service
	return consumer.publishEmailEvent(emailEvent)
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

		// Lấy thông tin shipping trực tiếp từ data
		if name, ok := data["shipping_name"].(string); ok && name != "" {
			orderData.ShippingName = name
		}
		if address, ok := data["shipping_address"].(string); ok && address != "" {
			orderData.ShippingAddress = address
		}
		if phone, ok := data["shipping_phone"].(string); ok && phone != "" {
			orderData.ShippingPhone = phone
		}

		// Log thông tin shipping để debug
		consumer.logger.Printf("DEBUG: Shipping info from event - Name: %s, Address: %s, Phone: %s",
			orderData.ShippingName, orderData.ShippingAddress, orderData.ShippingPhone)

		// Khối code cũ, vẫn giữ lại để tương thích ngược
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

	// Format items for display in email - this creates a more user-friendly presentation
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
		// Ghi log để biết chúng ta cần lấy thông tin từ user service
		consumer.logger.Printf("DEBUG: ShippingName is empty, attempting to retrieve user info for user ID %s", orderData.UserID)

		// Nếu có user ID, thử lấy thông tin user từ user-service
		if orderData.UserID != "" {
			// TODO: Triển khai gọi user-service để lấy thông tin người dùng
			// Đây là nơi bạn sẽ thêm code để gọi user-service API
			// Ví dụ: userInfo, err := getUserInfoFromUserService(orderData.UserID)

			// Tạm thời sử dụng email làm customer name nếu không có shipping name
			if orderData.UserEmail != "" {
				parts := strings.Split(orderData.UserEmail, "@")
				if len(parts) > 0 && parts[0] != "" {
					customerName = strings.Title(parts[0])
					consumer.logger.Printf("DEBUG: Using email user part as customer name: %s", customerName)
				}
			}
		}
	}

	// Create the order tracking URL (replace with your actual domain)
	orderURL := fmt.Sprintf("https://vntech.com/orders/%s", orderData.OrderID)

	// Create email data with all variables needed by the template
	emailData := EmailData{
		Type:     "order_confirmation",
		To:       orderData.UserEmail,
		Subject:  "Order Confirmation - #" + orderData.OrderNumber,
		Template: "order_confirmation.html.gohtml",
		Variables: map[string]string{
			// Customer information
			"customer_name": customerName,

			// Order details
			"order_id":       orderData.OrderID,
			"order_number":   orderData.OrderNumber,
			"order_date":     orderDate,
			"status":         orderData.Status,
			"payment_method": paymentMethod,
			"total":          totalStr,
			"items":          itemsStr,

			// Shipping information
			"shipping_name":    orderData.ShippingName,
			"shipping_address": orderData.ShippingAddress,
			"shipping_phone":   orderData.ShippingPhone,
			"order_url":        orderURL,

			// Additional variables that might be useful
			"year":       fmt.Sprintf("%d", time.Now().Year()),
			"user_id":    orderData.UserID,
			"user_email": orderData.UserEmail,
		},
	}

	// Log the variables we're sending for debugging
	consumer.logger.Printf("Order confirmation email variables: %+v", emailData.Variables)

	// Create standard email event
	emailEvent := StandardEvent{
		ID:         fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		Name:       "email.send",
		Data:       emailData,
		DataSchema: "v1",
		Source:     "listener-service",
		CreatedAt:  time.Now(),
		Version:    "v1",
	}

	// Forward to mail service
	return consumer.publishEmailEvent(emailEvent)
}

// Send OTP verification email
func (consumer *Consumer) sendOTPEmail(event StandardEvent) error {
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

	// Set subject based on purpose
	subject := "Your Verification Code"
	if otpData.Purpose == "password_reset" {
		subject = "Password Reset Verification"
	} else if otpData.Purpose == "registration" {
		subject = "Complete Your Registration"
	}

	// Create email data
	emailData := EmailData{
		Type:     "otp",
		To:       otpData.Email,
		Subject:  subject,
		Template: "otp_verification.html.gohtml",
		Variables: map[string]string{
			"email":       otpData.Email,
			"otp":         otpData.OTP,
			"purpose":     otpData.Purpose,
			"expires_in":  fmt.Sprintf("%d", otpData.ExpiresIn),
			"message":     otpData.Message,
			"action_text": otpData.ActionText,
			"year":        fmt.Sprintf("%d", time.Now().Year()),
		},
	}

	// Create standard email event
	emailEvent := StandardEvent{
		ID:         fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		Name:       "email.send",
		Data:       emailData,
		DataSchema: "v1",
		Source:     "listener-service",
		CreatedAt:  time.Now(),
		Version:    "v1",
	}

	// Forward to mail service
	return consumer.publishEmailEvent(emailEvent)
}

// Publish email event to mail service
func (consumer *Consumer) publishEmailEvent(event StandardEvent) error {
	// Set the PublishedAt timestamp just before publishing
	event.PublishedAt = time.Now()

	// Convert event to JSON
	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal email event: %w", err)
	}

	// Create a channel
	ch, err := consumer.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}
	defer ch.Close()

	// Publish the message
	err = ch.Publish(
		"logs_topic", // exchange
		"email.send", // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish email event: %w", err)
	}

	consumer.logger.Printf("Email event published: %s for %s", event.Name, event.Data.(EmailData).To)
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

	// Create email data with improved variables for better user experience
	emailData := EmailData{
		Type:     "payment_success",
		To:       paymentData.UserEmail,
		Subject:  "Payment Confirmation for Order #" + paymentData.OrderNumber,
		Template: "payment_success.html.gohtml",
		Variables: map[string]string{
			// Order information
			"order_id":     paymentData.OrderID,
			"order_number": paymentData.OrderNumber,

			// Payment details
			"amount":         amountStr,
			"currency":       currency,
			"payment_method": paymentMethod,
			"transaction_id": paymentData.TransactionID,
			"payment_date":   formattedPaymentDate,
			"payment_status": paymentData.Status,

			// Additional variables
			"year":       fmt.Sprintf("%d", time.Now().Year()),
			"user_id":    paymentData.UserID,
			"user_email": paymentData.UserEmail,
		},
	}

	// Create email event - pass the emailData object directly, not as a string
	emailEvent := StandardEvent{
		ID:         fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		Name:       "email.send",
		Data:       emailData, // Pass emailData directly, not as a JSON string
		DataSchema: "v1",
		Source:     "listener-service",
		CreatedAt:  time.Now(),
		Version:    "v1",
	}

	// Log the email we're about to send
	consumer.logger.Printf("Sending payment success email to %s for order #%s",
		paymentData.UserEmail, paymentData.OrderNumber)

	// Publish email event
	return consumer.publishEmailEvent(emailEvent)
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

	// Get appropriate subject line and template based on status
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

	// Create email data
	emailData := EmailData{
		Type:     "order_status",
		To:       statusData.UserEmail,
		Subject:  subject + " - Order #" + statusData.OrderNumber,
		Template: template,
		Variables: map[string]string{
			"order_id":        statusData.OrderID,
			"order_number":    statusData.OrderNumber,
			"status":          statusData.Status,
			"previous_status": statusData.PreviousStatus,
			"total":           totalStr,
		},
	}

	// Create email event - pass the emailData object directly
	emailEvent := StandardEvent{
		ID:         fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		Name:       "email.send",
		Data:       emailData, // Pass emailData directly, not as a JSON string
		DataSchema: "v1",
		Source:     "listener-service",
		CreatedAt:  time.Now(),
		Version:    "v1",
	}

	// Publish email event
	return consumer.publishEmailEvent(emailEvent)
}

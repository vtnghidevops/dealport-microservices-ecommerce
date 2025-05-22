package helpers

import (
	"errors"
	"listener-service/event"
	userpb "listener-service/proto/user"
	"listener-service/tests/mocks"
	"sync"

	"github.com/stretchr/testify/mock"
)

// ErrUnknownEventType is an error for unknown event types
var ErrUnknownEventType = errors.New("unknown event type")

// TestConsumerDependencies holds mock dependencies for consumer testing
type TestConsumerDependencies struct {
	MailClient   *mocks.MockMailClient
	LoggerClient *mocks.MockLoggerClient
	UserClient   *mocks.MockUserServiceClient
	Processor    *MessageProcessor
}

// NewTestConsumerDependencies creates mock dependencies for consumer testing
func NewTestConsumerDependencies() *TestConsumerDependencies {
	return &TestConsumerDependencies{
		MailClient:   new(mocks.MockMailClient),
		LoggerClient: new(mocks.MockLoggerClient),
		UserClient:   new(mocks.MockUserServiceClient),
		Processor:    NewMessageProcessor(),
	}
}

// SetupDefaultMockBehaviors configures common mock behaviors
func (t *TestConsumerDependencies) SetupDefaultMockBehaviors() {
	// Setup default logger behavior
	t.LoggerClient.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Setup default mail behaviors
	t.MailClient.On("SendWelcomeEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	t.MailClient.On("SendOrderConfirmationEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	t.MailClient.On("SendPasswordResetEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	t.MailClient.On("SendOrderStatusEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	t.MailClient.On("SendPaymentConfirmationEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	t.MailClient.On("SendPasswordChangedEmail", mock.Anything, mock.Anything).Return(nil)
	t.MailClient.On("SendOTPEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Setup default user service behavior
	t.UserClient.On("ProcessEvent", mock.Anything, mock.Anything).Return(&userpb.EventResponse{Success: true}, nil)
	t.UserClient.On("SyncUserOrderData", mock.Anything, mock.Anything).Return(&userpb.SyncUserOrderDataResponse{Success: true}, nil)
}

// TestConsumer is a test implementation of the Consumer interface
type TestConsumer struct {
	sync.Mutex
	Topics          []string
	ProcessedEvents map[string]bool
	EventHandlers   map[string]EventHandlerFunc
	Dependencies    *TestConsumerDependencies
	ErrorOnUnknown  bool
}

// NewTestConsumer creates a new test consumer
func NewTestConsumer() *TestConsumer {
	deps := NewTestConsumerDependencies()
	deps.SetupDefaultMockBehaviors()

	return &TestConsumer{
		ProcessedEvents: make(map[string]bool),
		EventHandlers:   make(map[string]EventHandlerFunc),
		Dependencies:    deps,
	}
}

// Listen mocks listening for topics
func (c *TestConsumer) Listen(topics []string) error {
	c.Lock()
	defer c.Unlock()
	c.Topics = topics
	return nil
}

// RegisterHandler registers a handler for a specific event type
func (c *TestConsumer) RegisterHandler(eventType string, handler EventHandlerFunc) {
	c.Lock()
	defer c.Unlock()
	c.EventHandlers[eventType] = handler
}

// HandleStandardEvent handles a standard event in testing
func (c *TestConsumer) HandleStandardEvent(e event.StandardEvent, topic string) error {
	c.Lock()
	c.ProcessedEvents[e.ID] = true
	c.Unlock()

	if handler, ok := c.EventHandlers[e.Name]; ok {
		return handler(e)
	}

	if c.ErrorOnUnknown {
		return ErrUnknownEventType
	}

	return nil
}

// HandleLegacyEvent handles a legacy event in testing
func (c *TestConsumer) HandleLegacyEvent(payload event.Payload) error {
	if handler, ok := c.EventHandlers[payload.Name]; ok {
		// Create a pseudo-event to pass to the handler
		pseudoEvent := event.StandardEvent{
			ID:   "legacy-" + payload.Name,
			Name: payload.Name,
			Data: payload.Data,
		}
		return handler(pseudoEvent)
	}

	if c.ErrorOnUnknown {
		return ErrUnknownEventType
	}

	return nil
}

// IsProcessed checks if an event has been processed
func (c *TestConsumer) IsProcessed(eventID string) bool {
	c.Lock()
	defer c.Unlock()
	return c.ProcessedEvents[eventID]
}

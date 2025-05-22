package helpers

import (
	"errors"
	"listener-service/event"
)

// MessageProcessor is a test helper for processing messages
type MessageProcessor struct {
	Handlers         map[string]EventHandlerFunc
	ProcessedEvents  map[string]bool
	HandlerErrors    map[string]error
	DefaultHandler   EventHandlerFunc
	ErrorOnUnknown   bool
	RecordProcessing bool
}

// EventHandlerFunc is a function type that handles events
type EventHandlerFunc func(event event.StandardEvent) error

// NewMessageProcessor creates a new message processor for testing
func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{
		Handlers:         make(map[string]EventHandlerFunc),
		ProcessedEvents:  make(map[string]bool),
		HandlerErrors:    make(map[string]error),
		RecordProcessing: true,
	}
}

// RegisterHandler registers a handler function for a specific event type
func (p *MessageProcessor) RegisterHandler(eventType string, handlerFunc EventHandlerFunc) {
	p.Handlers[eventType] = handlerFunc
}

// RegisterErrorHandler registers a handler that returns an error for testing error scenarios
func (p *MessageProcessor) RegisterErrorHandler(eventType string, err error) {
	p.HandlerErrors[eventType] = err
	p.RegisterHandler(eventType, func(e event.StandardEvent) error {
		return err
	})
}

// SetDefaultHandler sets a default handler for events without specific handlers
func (p *MessageProcessor) SetDefaultHandler(handlerFunc EventHandlerFunc) {
	p.DefaultHandler = handlerFunc
}

// ProcessEvent processes an event based on registered handlers
func (p *MessageProcessor) ProcessEvent(e event.StandardEvent) error {
	if p.RecordProcessing {
		p.ProcessedEvents[e.ID] = true
	}

	// If error is configured for this event type, return it
	if err, ok := p.HandlerErrors[e.Name]; ok {
		return err
	}

	// If handler exists for this event type, use it
	if handler, ok := p.Handlers[e.Name]; ok {
		return handler(e)
	}

	// If default handler exists, use it
	if p.DefaultHandler != nil {
		return p.DefaultHandler(e)
	}

	// If configured to error on unknown event types
	if p.ErrorOnUnknown {
		return errors.New("unknown event type: " + e.Name)
	}

	// No handler found but not configured to error
	return nil
}

// IsProcessed checks if an event has been processed
func (p *MessageProcessor) IsProcessed(eventID string) bool {
	return p.ProcessedEvents[eventID]
}

// ProcessLegacyPayload processes a legacy payload
func (p *MessageProcessor) ProcessLegacyPayload(payload event.Payload) error {
	// If error is configured for this payload type, return it
	if err, ok := p.HandlerErrors[payload.Name]; ok {
		return err
	}

	// If handler exists for this payload type, create a pseudo-event to process
	if handler, ok := p.Handlers[payload.Name]; ok {
		pseudoEvent := event.StandardEvent{
			ID:   "legacy-" + payload.Name,
			Name: payload.Name,
			Data: payload.Data,
		}
		return handler(pseudoEvent)
	}

	// If default handler exists, use it with pseudo-event
	if p.DefaultHandler != nil {
		pseudoEvent := event.StandardEvent{
			ID:   "legacy-" + payload.Name,
			Name: payload.Name,
			Data: payload.Data,
		}
		return p.DefaultHandler(pseudoEvent)
	}

	// If configured to error on unknown payload types
	if p.ErrorOnUnknown {
		return errors.New("unknown payload type: " + payload.Name)
	}

	// No handler found but not configured to error
	return nil
}

package events

import (
	"time"
)

// EventType represents the type of event that occurred
type EventType string

const (
	// System events
	EventTypeSystemStartup    EventType = "system_startup"
	EventTypeSystemShutdown   EventType = "system_shutdown"
	EventTypeSystemError      EventType = "system_error"
	
	// User events
	EventTypeUserLogin        EventType = "user_login"
	EventTypeUserLogout       EventType = "user_logout"
	EventTypeUserRegistration EventType = "user_registration"
	
	// Business events
	EventTypeOrderCreated     EventType = "order_created"
	EventTypeOrderCompleted   EventType = "order_completed"
	EventTypeOrderCancelled   EventType = "order_cancelled"
	EventTypePaymentReceived  EventType = "payment_received"
	EventTypePaymentFailed    EventType = "payment_failed"
	
	// Alert events
	EventTypeHighCPUUsage     EventType = "high_cpu_usage"
	EventTypeHighMemoryUsage  EventType = "high_memory_usage"
	EventTypeDiskSpaceLow     EventType = "disk_space_low"
	EventTypeServiceDown      EventType = "service_down"
)

// Event represents a notification event
type Event struct {
	ID          string                 `json:"id"`
	Type        EventType              `json:"type"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Severity    Severity               `json:"severity"`
	Timestamp   time.Time              `json:"timestamp"`
	UserID      string                 `json:"user_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Channel     string                 `json:"channel,omitempty"`
}

// Severity represents the severity level of an event
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityError    Severity = "error"
	SeverityCritical Severity = "critical"
)

// EventBuilder provides a fluent interface for building events
type EventBuilder struct {
	event *Event
}

// NewEvent creates a new event builder
func NewEvent(eventType EventType) *EventBuilder {
	return &EventBuilder{
		event: &Event{
			ID:        generateEventID(),
			Type:      eventType,
			Timestamp: time.Now(),
			Severity:  SeverityInfo,
			Metadata:  make(map[string]interface{}),
		},
	}
}

// WithTitle sets the event title
func (eb *EventBuilder) WithTitle(title string) *EventBuilder {
	eb.event.Title = title
	return eb
}

// WithMessage sets the event message
func (eb *EventBuilder) WithMessage(message string) *EventBuilder {
	eb.event.Message = message
	return eb
}

// WithSeverity sets the event severity
func (eb *EventBuilder) WithSeverity(severity Severity) *EventBuilder {
	eb.event.Severity = severity
	return eb
}

// WithUserID sets the user ID associated with the event
func (eb *EventBuilder) WithUserID(userID string) *EventBuilder {
	eb.event.UserID = userID
	return eb
}

// WithMetadata adds metadata to the event
func (eb *EventBuilder) WithMetadata(key string, value interface{}) *EventBuilder {
	eb.event.Metadata[key] = value
	return eb
}

// WithChannel sets the target Slack channel
func (eb *EventBuilder) WithChannel(channel string) *EventBuilder {
	eb.event.Channel = channel
	return eb
}

// Build returns the constructed event
func (eb *EventBuilder) Build() *Event {
	return eb.event
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

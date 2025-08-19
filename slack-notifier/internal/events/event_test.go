package events

import (
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	event := NewEvent(EventTypeUserLogin).Build()

	if event.Type != EventTypeUserLogin {
		t.Errorf("Expected event type to be %s, got %s", EventTypeUserLogin, event.Type)
	}

	if event.Severity != SeverityInfo {
		t.Errorf("Expected default severity to be %s, got %s", SeverityInfo, event.Severity)
	}

	if event.ID == "" {
		t.Error("Expected event ID to be generated")
	}

	if event.Timestamp.IsZero() {
		t.Error("Expected event timestamp to be set")
	}

	if event.Metadata == nil {
		t.Error("Expected event metadata to be initialized")
	}
}

func TestEventBuilder(t *testing.T) {
	event := NewEvent(EventTypeOrderCreated).
		WithTitle("Test Order").
		WithMessage("This is a test order").
		WithSeverity(SeverityWarning).
		WithUserID("test@example.com").
		WithMetadata("amount", "$100").
		WithMetadata("items", 2).
		WithChannel("#test-channel").
		Build()

	if event.Title != "Test Order" {
		t.Errorf("Expected title 'Test Order', got '%s'", event.Title)
	}

	if event.Message != "This is a test order" {
		t.Errorf("Expected message 'This is a test order', got '%s'", event.Message)
	}

	if event.Severity != SeverityWarning {
		t.Errorf("Expected severity %s, got %s", SeverityWarning, event.Severity)
	}

	if event.UserID != "test@example.com" {
		t.Errorf("Expected user ID 'test@example.com', got '%s'", event.UserID)
	}

	if event.Channel != "#test-channel" {
		t.Errorf("Expected channel '#test-channel', got '%s'", event.Channel)
	}

	if len(event.Metadata) != 2 {
		t.Errorf("Expected 2 metadata items, got %d", len(event.Metadata))
	}

	if event.Metadata["amount"] != "$100" {
		t.Errorf("Expected metadata amount '$100', got '%v'", event.Metadata["amount"])
	}

	if event.Metadata["items"] != 2 {
		t.Errorf("Expected metadata items 2, got '%v'", event.Metadata["items"])
	}
}

func TestEventIDGeneration(t *testing.T) {
	event1 := NewEvent(EventTypeUserLogin).Build()
	event2 := NewEvent(EventTypeUserLogin).Build()

	if event1.ID == event2.ID {
		t.Error("Expected different event IDs for different events")
	}

	// Check if ID follows expected format (timestamp + random string)
	if len(event1.ID) < 20 {
		t.Errorf("Expected event ID to be at least 20 characters, got %d", len(event1.ID))
	}
}

func TestEventTypes(t *testing.T) {
	// Test all event types
	eventTypes := []EventType{
		EventTypeSystemStartup,
		EventTypeSystemShutdown,
		EventTypeSystemError,
		EventTypeUserLogin,
		EventTypeUserLogout,
		EventTypeUserRegistration,
		EventTypeOrderCreated,
		EventTypeOrderCompleted,
		EventTypeOrderCancelled,
		EventTypePaymentReceived,
		EventTypePaymentFailed,
		EventTypeHighCPUUsage,
		EventTypeHighMemoryUsage,
		EventTypeDiskSpaceLow,
		EventTypeServiceDown,
	}

	for _, eventType := range eventTypes {
		event := NewEvent(eventType).Build()
		if event.Type != eventType {
			t.Errorf("Expected event type %s, got %s", eventType, event.Type)
		}
	}
}

func TestSeverityLevels(t *testing.T) {
	severities := []Severity{
		SeverityInfo,
		SeverityWarning,
		SeverityError,
		SeverityCritical,
	}

	for _, severity := range severities {
		event := NewEvent(EventTypeUserLogin).WithSeverity(severity).Build()
		if event.Severity != severity {
			t.Errorf("Expected severity %s, got %s", severity, event.Severity)
		}
	}
}

func TestEventTimestamp(t *testing.T) {
	before := time.Now()
	event := NewEvent(EventTypeUserLogin).Build()
	after := time.Now()

	if event.Timestamp.Before(before) || event.Timestamp.After(after) {
		t.Errorf("Event timestamp %v should be between %v and %v",
			event.Timestamp, before, after)
	}
}

func TestMultipleMetadata(t *testing.T) {
	event := NewEvent(EventTypeUserLogin).
		WithMetadata("key1", "value1").
		WithMetadata("key2", "value2").
		WithMetadata("key3", "value3").
		Build()

	expectedMetadata := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for key, expectedValue := range expectedMetadata {
		if value, exists := event.Metadata[key]; !exists {
			t.Errorf("Expected metadata key '%s' to exist", key)
		} else if value != expectedValue {
			t.Errorf("Expected metadata value '%v' for key '%s', got '%v'",
				expectedValue, key, value)
		}
	}
}

func TestEventBuilderChaining(t *testing.T) {
	// Test that all builder methods can be chained
	event := NewEvent(EventTypePaymentReceived).
		WithTitle("Payment").
		WithMessage("Payment received").
		WithSeverity(SeverityInfo).
		WithUserID("user@example.com").
		WithMetadata("amount", 100).
		WithChannel("#payments").
		Build()

	if event.Title != "Payment" {
		t.Error("Chained WithTitle failed")
	}
	if event.Message != "Payment received" {
		t.Error("Chained WithMessage failed")
	}
	if event.Severity != SeverityInfo {
		t.Error("Chained WithSeverity failed")
	}
	if event.UserID != "user@example.com" {
		t.Error("Chained WithUserID failed")
	}
	if event.Channel != "#payments" {
		t.Error("Chained WithChannel failed")
	}
	if event.Metadata["amount"] != 100 {
		t.Error("Chained WithMetadata failed")
	}
}

func TestRandomStringGeneration(t *testing.T) {
	// Test that random strings are generated
	str1 := randomString(6)
	str2 := randomString(6)

	if len(str1) != 6 {
		t.Errorf("Expected random string length 6, got %d", len(str1))
	}

	if len(str2) != 6 {
		t.Errorf("Expected random string length 6, got %d", len(str2))
	}

	// They should be different (though there's a small chance they could be the same)
	if str1 == str2 {
		t.Log("Warning: Random strings are the same, this is possible but unlikely")
	}
}

func TestEventTypeConstants(t *testing.T) {
	// Verify all event type constants are defined
	expectedTypes := map[string]EventType{
		"system_startup":    EventTypeSystemStartup,
		"system_shutdown":   EventTypeSystemShutdown,
		"system_error":      EventTypeSystemError,
		"user_login":        EventTypeUserLogin,
		"user_logout":       EventTypeUserLogout,
		"user_registration": EventTypeUserRegistration,
		"order_created":     EventTypeOrderCreated,
		"order_completed":   EventTypeOrderCompleted,
		"order_cancelled":   EventTypeOrderCancelled,
		"payment_received":  EventTypePaymentReceived,
		"payment_failed":    EventTypePaymentFailed,
		"high_cpu_usage":    EventTypeHighCPUUsage,
		"high_memory_usage": EventTypeHighMemoryUsage,
		"disk_space_low":    EventTypeDiskSpaceLow,
		"service_down":      EventTypeServiceDown,
	}

	for name, eventType := range expectedTypes {
		if eventType == "" {
			t.Errorf("Event type constant for '%s' is empty", name)
		}
	}
}

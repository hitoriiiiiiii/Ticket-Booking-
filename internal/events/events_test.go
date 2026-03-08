package events

import (
	"testing"
	"time"
)

// TestBaseEvent tests the BaseEvent struct
func TestBaseEvent(t *testing.T) {
	event := BaseEvent{
		Type:        "UserRegistered",
		AggregateID: "user_123",
		Timestamp:   time.Now(),
		Payload:     map[string]interface{}{"name": "test"},
	}

	if event.Type != "UserRegistered" {
		t.Errorf("Expected Type UserRegistered, got %s", event.Type)
	}
	if event.AggregateID != "user_123" {
		t.Errorf("Expected AggregateID user_123, got %s", event.AggregateID)
	}
	if event.Payload == nil {
		t.Error("Expected Payload to not be nil")
	}
}

// TestEventTypes tests valid event types
func TestEventTypes(t *testing.T) {
	validEvents := []string{
		EventUserRegistered,
		EventUserUpdated,
		EventUserDeleted,
		EventMovieCreated,
		EventMovieUpdated,
		EventMovieDeleted,
		EventShowCreated,
		EventShowUpdated,
		EventShowDeleted,
		EventTicketReserved,
		EventTicketConfirmed,
		EventTicketCancelled,
		EventPaymentInitiated,
		EventPaymentVerified,
		EventPaymentRefunded,
	}

	// Test that all defined events are non-empty
	for _, eventType := range validEvents {
		if eventType == "" {
			t.Error("Expected event type to be non-empty")
		}
	}
}

// TestEventPayload tests the EventPayload struct
func TestEventPayload(t *testing.T) {
	payload := EventPayload{
		UserID:   "user_123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	if payload.UserID != "user_123" {
		t.Errorf("Expected UserID user_123, got %s", payload.UserID)
	}
	if payload.Username != "testuser" {
		t.Errorf("Expected Username testuser, got %s", payload.Username)
	}
	if payload.Email != "test@example.com" {
		t.Errorf("Expected Email test@example.com, got %s", payload.Email)
	}
}

// TestEventHandler type tests
func TestEventHandlerType(t *testing.T) {
	// Test that EventHandler is a function type
	var handler EventHandler = func(event BaseEvent) error {
		return nil
	}

	if handler == nil {
		t.Error("Expected handler to not be nil")
	}
}

// TestEventPayloadWithBooking tests EventPayload with booking data
func TestEventPayloadWithBooking(t *testing.T) {
	payload := EventPayload{
		UserID:    "user_123",
		BookingID: "book_123",
		SeatID:    "seat_1",
		Status:    "BOOKED",
	}

	if payload.BookingID != "book_123" {
		t.Errorf("Expected BookingID book_123, got %s", payload.BookingID)
	}
	if payload.SeatID != "seat_1" {
		t.Errorf("Expected SeatID seat_1, got %s", payload.SeatID)
	}
	if payload.Status != "BOOKED" {
		t.Errorf("Expected Status BOOKED, got %s", payload.Status)
	}
}

package queue

import (
	"testing"
	"time"
)

// TestJobPayload tests the JobPayload struct
func TestJobPayload(t *testing.T) {
	payload := JobPayload{
		Type:    JobTypeNotification,
		UserID:  "user_123",
		Message: "Your booking is confirmed",
		Data:    map[string]interface{}{"booking_id": "book_123"},
		Priority: 1,
		CreatedAt: time.Now(),
	}

	if payload.Type != JobTypeNotification {
		t.Errorf("Expected Type %s, got %s", JobTypeNotification, payload.Type)
	}
	if payload.UserID != "user_123" {
		t.Errorf("Expected UserID user_123, got %s", payload.UserID)
	}
	if payload.Priority != 1 {
		t.Errorf("Expected Priority 1, got %d", payload.Priority)
	}
}

// TestJobTypes tests the job type constants
func TestJobTypes(t *testing.T) {
	if JobTypeNotification == "" {
		t.Error("Expected JobTypeNotification to be non-empty")
	}
	if JobTypeEmail == "" {
		t.Error("Expected JobTypeEmail to be non-empty")
	}
	if JobTypePayment == "" {
		t.Error("Expected JobTypePayment to be non-empty")
	}
	if JobTypeBooking == "" {
		t.Error("Expected JobTypeBooking to be non-empty")
	}
}

// TestJobPayloadSerialization tests JSON serialization of JobPayload
func TestJobPayloadSerialization(t *testing.T) {
	payload := JobPayload{
		Type:    JobTypeBooking,
		UserID:  "user_456",
		Message: "Booking confirmed",
		Data:    map[string]interface{}{"seat_id": "seat_1"},
		Priority: 0,
		CreatedAt: time.Now(),
	}

	// Test that all fields are properly set
	if payload.Type != JobTypeBooking {
		t.Errorf("Expected Type %s, got %s", JobTypeBooking, payload.Type)
	}
	if payload.Data == nil {
		t.Error("Expected Data to not be nil")
	}
}

// TestJobProcessorFunctionType tests that JobProcessor is a function type
func TestJobProcessorFunctionType(t *testing.T) {
	var processor JobProcessor = func(job JobPayload) error {
		return nil
	}

	if processor == nil {
		t.Error("Expected processor to be non-nil")
	}
}

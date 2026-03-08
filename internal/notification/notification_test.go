package notification

import (
	"testing"
	"time"
)

// TestNotificationModel tests the Notification struct
func TestNotificationModel(t *testing.T) {
	notification := Notification{
		ID:        "notif_123",
		UserID:    "user_123",
		Message:   "Your booking is confirmed",
		Type:      "BOOKING",
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if notification.ID != "notif_123" {
		t.Errorf("Expected ID notif_123, got %s", notification.ID)
	}
	if notification.UserID != "user_123" {
		t.Errorf("Expected UserID user_123, got %s", notification.UserID)
	}
	if notification.Message != "Your booking is confirmed" {
		t.Errorf("Expected message, got %s", notification.Message)
	}
	if notification.Type != "BOOKING" {
		t.Errorf("Expected Type BOOKING, got %s", notification.Type)
	}
	if notification.IsRead != false {
		t.Errorf("Expected IsRead false, got %v", notification.IsRead)
	}
}

// TestNotificationTypes tests valid notification types
func TestNotificationTypes(t *testing.T) {
	validTypes := []string{"BOOKING", "PAYMENT", "SHOW", "GENERAL"}

	testCases := []struct {
		notifType string
		expected  bool
	}{
		{"BOOKING", true},
		{"PAYMENT", true},
		{"SHOW", true},
		{"GENERAL", true},
		{"INVALID", false},
		{"", false},
	}

	for _, tc := range testCases {
		isValid := false
		for _, ntype := range validTypes {
			if tc.notifType == ntype {
				isValid = true
				break
			}
		}
		if isValid != tc.expected {
			t.Errorf("Expected type %s to be %v", tc.notifType, tc.expected)
		}
	}
}

// TestNotificationMarkAsRead tests marking notification as read
func TestNotificationMarkAsRead(t *testing.T) {
	notification := Notification{
		ID:        "notif_456",
		UserID:    "user_123",
		Message:   "Test notification",
		Type:      "GENERAL",
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if notification.IsRead {
		t.Error("Expected IsRead to be false initially")
	}

	// Simulate marking as read
	notification.IsRead = true

	if !notification.IsRead {
		t.Error("Expected IsRead to be true after marking")
	}
}

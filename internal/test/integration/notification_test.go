// Package integration provides integration tests for the ticket booking system
package integration

import (
	"context"
	"testing"

	"github.com/hitorii/ticket-booking/internal/notification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNotificationIntegration_SendNotification tests notification sending
func TestNotificationIntegration_SendNotification(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Clean up before test
	svc.CommandDB.Pool.Exec(ctx, "DELETE FROM notifications")

	t.Run("SendSuccess", func(t *testing.T) {
		req := notification.SendNotificationRequest{
			UserID:  "user-1",
			Message: "Test notification message",
			Type:    "booking",
		}

		notif, err := svc.NotificationCmdSvc.SendNotification(ctx, req)
		require.NoError(t, err)
		assert.NotEmpty(t, notif.ID)
		assert.Equal(t, req.UserID, notif.UserID)
		assert.Equal(t, req.Message, notif.Message)
		assert.Equal(t, req.Type, notif.Type)
		assert.False(t, notif.IsRead)
	})

	t.Run("SendInvalidUserID", func(t *testing.T) {
		req := notification.SendNotificationRequest{
			UserID:  "",
			Message: "Test notification message",
			Type:    "booking",
		}

		_, err := svc.NotificationCmdSvc.SendNotification(ctx, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")
	})

	t.Run("SendEmptyMessage", func(t *testing.T) {
		req := notification.SendNotificationRequest{
			UserID:  "user-1",
			Message: "",
			Type:    "booking",
		}

		_, err := svc.NotificationCmdSvc.SendNotification(ctx, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "message is required")
	})
}

// TestNotificationIntegration_MarkAsRead tests marking notifications as read
func TestNotificationIntegration_MarkAsRead(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Clean up before test
	svc.CommandDB.Pool.Exec(ctx, "DELETE FROM notifications")

	// First create a notification
	req := notification.SendNotificationRequest{
		UserID:  "user-mark-read",
		Message: "Test notification for mark as read",
		Type:    "booking",
	}
	notif, err := svc.NotificationCmdSvc.SendNotification(ctx, req)
	require.NoError(t, err)

	t.Run("MarkAsRead", func(t *testing.T) {
		err := svc.NotificationCmdSvc.MarkAsRead(ctx, notif.ID)
		require.NoError(t, err)

		// Verify the notification was marked as read
		updatedNotif, err := svc.NotificationQuerySvc.GetNotificationByID(ctx, notif.ID)
		require.NoError(t, err)
		assert.True(t, updatedNotif.IsRead)
	})

	t.Run("MarkAsRead_InvalidID", func(t *testing.T) {
		err := svc.NotificationCmdSvc.MarkAsRead(ctx, "non-existent-id")
		require.Error(t, err)
	})
}

// TestNotificationIntegration_MarkAllAsRead tests marking all notifications as read for a user
func TestNotificationIntegration_MarkAllAsRead(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Clean up before test
	svc.CommandDB.Pool.Exec(ctx, "DELETE FROM notifications")

	// Create multiple notifications for the same user
	userID := "user-mark-all-read"
	for i := 0; i < 3; i++ {
		req := notification.SendNotificationRequest{
			UserID:  userID,
			Message: "Test notification message",
			Type:    "booking",
		}
		_, err := svc.NotificationCmdSvc.SendNotification(ctx, req)
		require.NoError(t, err)
	}

	t.Run("MarkAllAsRead", func(t *testing.T) {
		err := svc.NotificationCmdSvc.MarkAllAsRead(ctx, userID)
		require.NoError(t, err)

		// Verify all notifications were marked as read
		notifs, err := svc.NotificationQuerySvc.GetUserNotifications(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, notifs, 3)
		for _, n := range notifs {
			assert.True(t, n.IsRead)
		}
	})

	t.Run("MarkAllAsRead_InvalidUserID", func(t *testing.T) {
		err := svc.NotificationCmdSvc.MarkAllAsRead(ctx, "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")
	})
}

// TestNotificationIntegration_Delete tests deleting notifications
func TestNotificationIntegration_Delete(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Clean up before test
	svc.CommandDB.Pool.Exec(ctx, "DELETE FROM notifications")

	// First create a notification
	req := notification.SendNotificationRequest{
		UserID:  "user-delete",
		Message: "Test notification for delete",
		Type:    "booking",
	}
	notif, err := svc.NotificationCmdSvc.SendNotification(ctx, req)
	require.NoError(t, err)

	t.Run("Delete", func(t *testing.T) {
		err := svc.NotificationCmdSvc.DeleteNotification(ctx, notif.ID)
		require.NoError(t, err)

		// Verify the notification was deleted
		deletedNotif, err := svc.NotificationQuerySvc.GetNotificationByID(ctx, notif.ID)
		require.Error(t, err)
		assert.Nil(t, deletedNotif)
	})

	t.Run("Delete_InvalidID", func(t *testing.T) {
		err := svc.NotificationCmdSvc.DeleteNotification(ctx, "non-existent-id")
		require.Error(t, err)
	})
}

// Command handler for notification write operations (CQRS)

package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommandHandler struct {
	CommandService *CommandService
}

func NewCommandHandler(cs *CommandService) *CommandHandler {
	return &CommandHandler{CommandService: cs}
}

// SendNotification - Command handler for sending a notification
func (h *CommandHandler) SendNotification(c *gin.Context) {
	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	notification, err := h.CommandService.SendNotification(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send notification: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// MarkAsRead - Command handler for marking a notification as read
func (h *CommandHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	err := h.CommandService.MarkAsRead(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to mark notification as read: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Notification marked as read",
	})
}

// MarkAllAsRead - Command handler for marking all notifications as read for a user
func (h *CommandHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.Param("userID")

	err := h.CommandService.MarkAllAsRead(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to mark notifications as read: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "All notifications marked as read",
	})
}

// DeleteNotification - Command handler for deleting a notification
func (h *CommandHandler) DeleteNotification(c *gin.Context) {
	id := c.Param("id")

	err := h.CommandService.DeleteNotification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete notification: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Notification deleted",
	})
}

// SendNotificationRequest - Request model for sending notification
type SendNotificationRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Type    string `json:"type"` // email, sms, push, etc.
}

// Query handler for notification read operations (CQRS)

package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueryHandler struct {
	QueryService *QueryService
}

func NewQueryHandler(qs *QueryService) *QueryHandler {
	return &QueryHandler{QueryService: qs}
}

// GetUserNotifications - Query handler for getting all notifications for a user
func (h *QueryHandler) GetUserNotifications(c *gin.Context) {
	userID := c.Param("userID")

	notifications, err := h.QueryService.GetUserNotifications(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch notifications",
		})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GetUnreadNotifications - Query handler for getting unread notifications for a user
func (h *QueryHandler) GetUnreadNotifications(c *gin.Context) {
	userID := c.Param("userID")

	notifications, err := h.QueryService.GetUnreadNotifications(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch notifications",
		})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GetNotification - Query handler for getting a single notification by ID
func (h *QueryHandler) GetNotification(c *gin.Context) {
	id := c.Param("id")

	notification, err := h.QueryService.GetNotificationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}
	c.JSON(http.StatusOK, notification)
}

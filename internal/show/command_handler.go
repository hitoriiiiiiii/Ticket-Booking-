// Command handler for show write operations (CQRS)

package show

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CommandHandler struct {
	CommandService *CommandService
}

func NewCommandHandler(cs *CommandService) *CommandHandler {
	return &CommandHandler{CommandService: cs}
}

// CreateShow - Command handler for creating a new show
func (h *CommandHandler) CreateShow(c *gin.Context) {
	var req CreateShowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	show, err := h.CommandService.CreateShow(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create show: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, show)
}

// UpdateShow - Command handler for updating an existing show
func (h *CommandHandler) UpdateShow(c *gin.Context) {
	id := c.Param("id")
	
	var req UpdateShowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.CommandService.UpdateShow(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update show: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Show updated successfully"})
}

// DeleteShow - Command handler for deleting a show
func (h *CommandHandler) DeleteShow(c *gin.Context) {
	id := c.Param("id")

	err := h.CommandService.DeleteShow(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete show: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Show deleted successfully"})
}

// CreateShowRequest - Request model for creating a show
type CreateShowRequest struct {
	MovieID   int       `json:"movie_id"`
	Theater   string    `json:"theater"`
	StartTime time.Time `json:"start_time"`
}

// UpdateShowRequest - Request model for updating a show
type UpdateShowRequest struct {
	MovieID   int       `json:"movie_id"`
	Theater   string    `json:"theater"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

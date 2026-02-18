// Command handler for movie write operations (CQRS)

package movie

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

// CreateMovie - Command handler for creating a new movie
func (h *CommandHandler) CreateMovie(c *gin.Context) {
	var req CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	movie, err := h.CommandService.CreateMovie(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

// UpdateMovie - Command handler for updating an existing movie
func (h *CommandHandler) UpdateMovie(c *gin.Context) {
	id := c.Param("id")
	
	var req UpdateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.CommandService.UpdateMovie(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully"})
}

// DeleteMovie - Command handler for deleting a movie
func (h *CommandHandler) DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	err := h.CommandService.DeleteMovie(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}

// Command handler for user write operations (CQRS)

package user

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

// Register - Command handler for user registration
func (h *CommandHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.CommandService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser - Command handler for updating user information
func (h *CommandHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.CommandService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser - Command handler for deleting a user
func (h *CommandHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.CommandService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// RegisterRequest - Request model for user registration
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

// UpdateUserRequest - Request model for updating user
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

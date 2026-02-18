// Query handler for user read operations (CQRS)

package user

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

// Login - Query handler for user login (verifies credentials)
func (h *QueryHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.QueryService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ListUsers - Query handler for listing all users
func (h *QueryHandler) ListUsers(c *gin.Context) {
	users, err := h.QueryService.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser - Query handler for getting a single user by ID
func (h *QueryHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	
	user, err := h.QueryService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// LoginRequest - Request model for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

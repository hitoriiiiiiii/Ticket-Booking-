//handler 
package notifications

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{Repo: repo}
}

func (h *Handler) GetUserNotifications(c *gin.Context) {

	userID := c.Param("userID")

	list, err := h.Repo.GetUserNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch notifications",
		})
		return
	}

	c.JSON(http.StatusOK, list)
}

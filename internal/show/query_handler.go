// Query handler for show read operations (CQRS)

package show

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

// GetShows - Query handler for listing all shows
func (h *QueryHandler) GetShows(c *gin.Context) {
	shows, err := h.QueryService.GetAllShows(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shows"})
		return
	}
	c.JSON(http.StatusOK, shows)
}

// GetShow - Query handler for getting a single show by ID
func (h *QueryHandler) GetShow(c *gin.Context) {
	id := c.Param("id")
	
	show, err := h.QueryService.GetShowByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Show not found"})
		return
	}
	c.JSON(http.StatusOK, show)
}

// GetShowsByMovie - Query handler for getting shows by movie ID
func (h *QueryHandler) GetShowsByMovie(c *gin.Context) {
	movieID := c.Param("movieID")
	
	shows, err := h.QueryService.GetShowsByMovie(c.Request.Context(), movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shows"})
		return
	}
	c.JSON(http.StatusOK, shows)
}

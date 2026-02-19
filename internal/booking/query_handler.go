package booking

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

// CheckAvailability - Query handler for checking seat availability
func (h *QueryHandler) CheckAvailability(c *gin.Context) {
	seatID := c.Param("seat_id")

	status, err := h.QueryService.GetAvailability(c.Request.Context(), seatID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"seat_id": seatID,
			"status":  "AVAILABLE",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"seat_id": status.SeatID,
		"status":  status.Status,
	})
}

// GetEvents - Query handler for getting all events
func (h *QueryHandler) GetEvents(c *gin.Context) {
	eventList, err := h.QueryService.GetEvents(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch events"})
		return
	}

	c.JSON(200, eventList)
}

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

// GetUserReservations - Query handler for getting all reservations for a user
func (h *QueryHandler) GetUserReservations(c *gin.Context) {
	userID := c.Param("user_id")

	reservations, err := h.QueryService.GetUserReservations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}

	// Convert reservations to gin.H format
	result := make([]gin.H, len(reservations))
	for i, r := range reservations {
		result[i] = gin.H{
			"id":     r.ID,
			"user_id": r.UserID,
			"seat_id": r.SeatID,
			"status":  r.Status,
		}
	}

	if result == nil {
		result = []gin.H{}
	}

	c.JSON(http.StatusOK, result)
}

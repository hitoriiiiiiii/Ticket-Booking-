package booking

import (
	"github.com/gin-gonic/gin"
	"github.com/hitorii/ticket-booking/internal/events"
	"net/http"
)

// HealthCheck handles the health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

//Ticket Reservation handler (CQRS command handler)
type Handler struct {
	EventStore *events.Store
}

func (h * Handler) ReserveTicket(c *gin.Context){
	var req struct {
		UserID string `json:"user_id"`	
		SeatID string `json:"seat_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	//Append Events 
	err := h.EventStore.Append("TicketReserved", req.SeatID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "event store failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket Reserved ✅ Event Stored",
	})

}

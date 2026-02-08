package booking

import (
	"context"
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

	// Check if seat is already reserved
	reserved, err := h.EventStore.IsSeatReserved(req.SeatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "event store query failed"})
		return
	}

	if reserved {
		c.JSON(http.StatusConflict, gin.H{"error": "Seat already reserved"})
		return
	}

	//Append Events
	err = h.EventStore.Append("TicketReserved", req.SeatID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "event store failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket Reserved ✅ Event Stored",
	})
}

func (h *Handler) GetEvents(c *gin.Context) {

	rows, err := h.EventStore.DB.Query(
		context.Background(),
		"SELECT aggregate_id, event_type, payload, created_at FROM events",
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch events"})
		return
	}

	var events []map[string]interface{}

	for rows.Next() {
		var aggID, eventType string
		var payload []byte
		var createdAt string

		rows.Scan(&aggID, &eventType, &payload, &createdAt)

		events = append(events, gin.H{
			"seat_id": aggID,
			"type":   eventType,
			"payload": string(payload),
			"time":   createdAt,
		})
	}

	c.JSON(200, events)
}

//Cancel Ticket Handler
func (h *Handler) CancelTicket(c *gin.Context){
	var req struct {	
		UserID string `json:"user_id"`
		SeatID string `json:"seat_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	//Append Cancel Event
	err := h.EventStore.Append("TicketCancelled", req.SeatID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "event store failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket Cancelled ✅ Event Stored",
	})
}

//Check availability handler (CQRS query handler)
func (h *Handler) CheckAvailability(c *gin.Context) {

	seatID := c.Param("seat_id")

	available, err := h.EventStore.IsSeatAvailable(seatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "availability check failed"})
		return
	}

	if available {
		c.JSON(http.StatusOK, gin.H{
			"seat_id":   seatID,
			"available": true,
			"message":   "Seat is Available ✅",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"seat_id":   seatID,
			"available": false,
			"message":   "Seat is Reserved ❌",
		})
	}
}
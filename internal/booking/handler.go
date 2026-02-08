package booking

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hitorii/ticket-booking/internal/events"
)

type Handler struct {
	EventStore *events.Store
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func (h *Handler) ReserveTicket(c *gin.Context) {

	var req struct {
		UserID string `json:"user_id"`
		SeatID string `json:"seat_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := h.EventStore.DB.Exec(
		context.Background(),
		"INSERT INTO reservations (seat_id, user_id, status) VALUES ($1,$2,'HELD')",
		req.SeatID,
		req.UserID,
	)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Seat already reserved",
		})
		return
	}

	_ = h.EventStore.Append("TicketReserved", req.SeatID, req)

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket Reserved",
	})
}

func (h *Handler) ConfirmTicket(c *gin.Context) {

	var req struct {
		UserID string `json:"user_id"`
		SeatID string `json:"seat_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := h.EventStore.DB.Exec(
		context.Background(),
		`UPDATE reservations 
		 SET status='BOOKED'
		 WHERE seat_id=$1 AND user_id=$2 AND status='HELD'`,
		req.SeatID,
		req.UserID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to confirm booking",
		})
		return
	}

	if res.RowsAffected() == 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Seat not held or already booked",
		})
		return
	}

	_ = h.EventStore.Append("TicketConfirmed", req.SeatID, req)

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket Confirmed",
	})
}

func (h *Handler) CancelTicket(c *gin.Context) {

	var req struct {
		UserID string `json:"user_id"`
		SeatID string `json:"seat_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	res, err := h.EventStore.DB.Exec(
		context.Background(),
		"DELETE FROM reservations WHERE seat_id=$1 AND user_id=$2",
		req.SeatID,
		req.UserID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cancel failed",
		})
		return
	}

	if res.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No reservation found",
		})
		return
	}

	_ = h.EventStore.Append("TicketCancelled", req.SeatID, req)

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket Cancelled",
	})
}

func (h *Handler) CheckAvailability(c *gin.Context) {

	seatID := c.Param("seat_id")

	var status string

	err := h.EventStore.DB.QueryRow(
		context.Background(),
		"SELECT status FROM reservations WHERE seat_id=$1",
		seatID,
	).Scan(&status)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"seat_id": seatID,
			"status":  "AVAILABLE",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"seat_id": seatID,
		"status":  status,
	})
}

func (h *Handler) GetEvents(c *gin.Context) {

	rows, err := h.EventStore.DB.Query(
		context.Background(),
		"SELECT aggregate_id, event_type, payload, created_at FROM events ORDER BY created_at DESC",
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch events"})
		return
	}

	var eventsList []gin.H

	for rows.Next() {
		var aggID, eventType string
		var payload []byte
		var createdAt string

		rows.Scan(&aggID, &eventType, &payload, &createdAt)

		eventsList = append(eventsList, gin.H{
			"seat_id":  aggID,
			"type":     eventType,
			"payload":  string(payload),
			"time":     createdAt,
		})
	}

	c.JSON(200, eventsList)
}

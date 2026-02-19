package booking

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hitorii/ticket-booking/internal/events"
)

type CommandHandler struct {
	CommandService *CommandService
	EventStore     *events.Store
}

func NewCommandHandler(cs *CommandService, es *events.Store) *CommandHandler {
	return &CommandHandler{
		CommandService: cs,
		EventStore:     es,
	}
}

// ReserveTicket - Command handler for reserving a ticket
func (h *CommandHandler) ReserveTicket(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
		SeatID string `json:"seat_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.CommandService.ReserveTicket(c.Request.Context(), req.UserID, req.SeatID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	_ = h.EventStore.Append("TicketReserved", req.SeatID, req)

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket Reserved",
	})
}

// ConfirmTicket - Command handler for confirming a ticket reservation
func (h *CommandHandler) ConfirmTicket(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
		SeatID string `json:"seat_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.CommandService.ConfirmTicket(c.Request.Context(), req.UserID, req.SeatID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	_ = h.EventStore.Append("TicketConfirmed", req.SeatID, req)

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket Confirmed",
	})
}

// CancelTicket - Command handler for cancelling a ticket reservation
func (h *CommandHandler) CancelTicket(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
		SeatID string `json:"seat_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.CommandService.CancelTicket(c.Request.Context(), req.UserID, req.SeatID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	_ = h.EventStore.Append("TicketCancelled", req.SeatID, req)

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket Cancelled",
	})
}

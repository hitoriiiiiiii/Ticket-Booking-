// Command handler for payment write operations (CQRS)

package payments

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

// InitiatePayment - Command handler for initiating a new payment
func (h *CommandHandler) InitiatePayment(c *gin.Context) {
	var req InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	payment, err := h.CommandService.InitiatePayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to initiate payment: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// VerifyPayment - Command handler for verifying a payment
func (h *CommandHandler) VerifyPayment(c *gin.Context) {
	var req VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err := h.CommandService.VerifyPayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify payment: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Payment verification completed",
	})
}

// RefundPayment - Command handler for refunding a payment
func (h *CommandHandler) RefundPayment(c *gin.Context) {
	paymentID := c.Param("id")

	err := h.CommandService.RefundPayment(c.Request.Context(), paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to refund payment: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Payment refund completed",
	})
}

// InitiatePaymentRequest - Request model for initiating payment
type InitiatePaymentRequest struct {
	BookingID string `json:"booking_id"`
	UserID    string `json:"user_id"`
	Amount    int    `json:"amount"`
}

// VerifyPaymentRequest - Request model for verifying payment
type VerifyPaymentRequest struct {
	PaymentID string `json:"payment_id"`
	Mode      string `json:"mode"` // success/fail/random
}

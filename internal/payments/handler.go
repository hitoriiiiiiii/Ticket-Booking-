package payments

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

// POST /payments/initiate
func (h *Handler) InitiatePayment(c *gin.Context) {

	var req struct {
		BookingID string `json:"booking_id"`
		UserID    string `json:"user_id"`
		Amount    int    `json:"amount"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	payment, err := h.Service.InitiatePayment(req.BookingID, req.UserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to initiate payment",
		})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// POST /payments/verify
func (h *Handler) VerifyPayment(c *gin.Context) {

	var req struct {
		PaymentID string `json:"payment_id"`
		Mode      string `json:"mode"` // success/fail/random
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err := h.Service.Verify(req.PaymentID, req.Mode)
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

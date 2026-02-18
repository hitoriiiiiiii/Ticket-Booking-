// Query handler for payment read operations (CQRS)

package payments

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

// GetPayment - Query handler for getting a payment by ID
func (h *QueryHandler) GetPayment(c *gin.Context) {
	id := c.Param("id")
	
	payment, err := h.QueryService.GetPaymentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	c.JSON(http.StatusOK, payment)
}

// GetPaymentByBooking - Query handler for getting a payment by booking ID
func (h *QueryHandler) GetPaymentByBooking(c *gin.Context) {
	bookingID := c.Param("bookingID")
	
	payment, err := h.QueryService.GetPaymentByBookingID(c.Request.Context(), bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	c.JSON(http.StatusOK, payment)
}

// GetPaymentsByUser - Query handler for getting all payments for a user
func (h *QueryHandler) GetPaymentsByUser(c *gin.Context) {
	userID := c.Param("userID")
	
	payments, err := h.QueryService.GetPaymentsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}
	c.JSON(http.StatusOK, payments)
}

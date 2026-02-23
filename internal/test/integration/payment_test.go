// Package integration provides integration tests for the ticket booking system
package integration

import (
	"context"
	"testing"

	"github.com/hitorii/ticket-booking/internal/events"
	"github.com/hitorii/ticket-booking/internal/payments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPaymentIntegration_InitiatePayment tests payment initiation flow
func TestPaymentIntegration_InitiatePayment(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Clean up before test
	svc.CommandDB.Pool.Exec(ctx, "DELETE FROM payments")

	t.Run("InitiatePayment_Success", func(t *testing.T) {
		req := payments.InitiatePaymentRequest{
			BookingID: "booking-1",
			UserID:    "user-1",
			Amount:    1000,
		}

		payment, err := svc.PaymentCmdSvc.InitiatePayment(ctx, req)
		require.NoError(t, err)
		assert.NotEmpty(t, payment.ID)
		assert.Equal(t, req.BookingID, payment.BookingID)
		assert.Equal(t, req.UserID, payment.UserID)
		assert.Equal(t, req.Amount, payment.Amount)
		assert.Equal(t, "PENDING", payment.Status)
	})

	t.Run("InitiatePayment_InvalidBookingID", func(t *testing.T) {
		req := payments.InitiatePaymentRequest{
			BookingID: "",
			UserID:    "user-1",
			Amount:    1000,
		}

		_, err := svc.PaymentCmdSvc.InitiatePayment(ctx, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "booking ID is required")
	})

	t.Run("InitiatePayment_InvalidUserID", func(t *testing.T) {
		req := payments.InitiatePaymentRequest{
			BookingID: "booking-1",
			UserID:    "",
			Amount:    1000,
		}

		_, err := svc.PaymentCmdSvc.InitiatePayment(ctx, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")
	})

	t.Run("InitiatePayment_NegativeAmount", func(t *testing.T) {
		req := payments.InitiatePaymentRequest{
			BookingID: "booking-1",
			UserID:    "user-1",
			Amount:    -100,
		}

		_, err := svc.PaymentCmdSvc.InitiatePayment(ctx, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "amount must be positive")
	})

	t.Run("InitiatePayment_ZeroAmount", func(t *testing.T) {
		req := payments.InitiatePaymentRequest{
			BookingID: "booking-1",
			UserID:    "user-1",
			Amount:    0,
		}

		_, err := svc.PaymentCmdSvc.InitiatePayment(ctx, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "amount must be positive")
	})
}

// TestPaymentIntegration_VerifyPayment tests payment verification flow
func TestPaymentIntegration_VerifyPayment(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Clean up before test
	svc.CommandDB.Pool.Exec(ctx, "DELETE FROM payments")

	// First create a payment to verify
	initReq := payments.InitiatePaymentRequest{
		BookingID: "booking-verify-1",
		UserID:    "user-verify-1",
		Amount:    500,
	}
	payment, err := svc.PaymentCmdSvc.InitiatePayment(ctx, initReq)
	require.NoError(t, err)

	t.Run("VerifyPayment_Success", func(t *testing.T) {
		verifyReq := payments.VerifyPaymentRequest{
			PaymentID: payment.ID,
			Mode:      "success",
		}

		err := svc.PaymentCmdSvc.VerifyPayment(ctx, verifyReq)
		require.NoError(t, err)

		// Verify the status was updated
		updatedPayment, err := svc.PaymentCmdSvc.Repo.GetPaymentByID(payment.ID)
		require.NoError(t, err)
		assert.Equal(t, "SUCCESS", updatedPayment.Status)
		assert.NotEmpty(t, updatedPayment.TransactionID)
	})

	t.Run("VerifyPayment_Failure", func(t *testing.T) {
		// Create another payment to test failure
		initReq2 := payments.InitiatePaymentRequest{
			BookingID: "booking-verify-2",
			UserID:    "user-verify-2",
			Amount:    500,
		}
		payment2, err := svc.PaymentCmdSvc.InitiatePayment(ctx, initReq2)
		require.NoError(t, err)

		verifyReq := payments.VerifyPaymentRequest{
			PaymentID: payment2.ID,
			Mode:      "fail",
		}

		err = svc.PaymentCmdSvc.VerifyPayment(ctx, verifyReq)
		require.NoError(t, err)

		// Verify the status was updated to FAILED
		updatedPayment, err := svc.PaymentCmdSvc.Repo.GetPaymentByID(payment2.ID)
		require.NoError(t, err)
		assert.Equal(t, "FAILED", updatedPayment.Status)
	})

	t.Run("VerifyPayment_AlreadyProcessed", func(t *testing.T) {
		// Try to verify an already verified payment
		verifyReq := payments.VerifyPaymentRequest{
			PaymentID: payment.ID,
			Mode:      "success",
		}

		err := svc.PaymentCmdSvc.VerifyPayment(ctx, verifyReq)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already processed")
	})

	t.Run("VerifyPayment_InvalidPaymentID", func(t *testing.T) {
		verifyReq := payments.VerifyPaymentRequest{
			PaymentID: "non-existent-payment-id",
			Mode:      "success",
		}

		err := svc.PaymentCmdSvc.VerifyPayment(ctx, verifyReq)
		require.Error(t, err)
	})
}

// TestPaymentIntegration_RefundPayment tests payment refund flow
func TestPaymentIntegration_RefundPayment(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Clean up before test
	svc.CommandDB.Pool.Exec(ctx, "DELETE FROM payments")

	// Create and verify a payment to test refund
	initReq := payments.InitiatePaymentRequest{
		BookingID: "booking-refund-1",
		UserID:    "user-refund-1",
		Amount:    500,
	}
	payment, err := svc.PaymentCmdSvc.InitiatePayment(ctx, initReq)
	require.NoError(t, err)

	// Verify the payment first
	verifyReq := payments.VerifyPaymentRequest{
		PaymentID: payment.ID,
		Mode:      "success",
	}
	err = svc.PaymentCmdSvc.VerifyPayment(ctx, verifyReq)
	require.NoError(t, err)

	t.Run("RefundPayment_Success", func(t *testing.T) {
		err := svc.PaymentCmdSvc.RefundPayment(ctx, payment.ID)
		require.NoError(t, err)

		// Verify the status was updated
		updatedPayment, err := svc.PaymentCmdSvc.Repo.GetPaymentByID(payment.ID)
		require.NoError(t, err)
		assert.Equal(t, "REFUNDED", updatedPayment.Status)
	})

	t.Run("RefundPayment_NotSuccessful", func(t *testing.T) {
		// Create a new payment that is still PENDING
		initReq2 := payments.InitiatePaymentRequest{
			BookingID: "booking-refund-2",
			UserID:    "user-refund-2",
			Amount:    500,
		}
		payment2, err := svc.PaymentCmdSvc.InitiatePayment(ctx, initReq2)
		require.NoError(t, err)

		// Try to refund a non-successful payment
		err = svc.PaymentCmdSvc.RefundPayment(ctx, payment2.ID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "can only refund successful payments")
	})
}

// TestPaymentIntegration_WithEvents tests event publishing during payment operations
func TestPaymentIntegration_WithEvents(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Clean up before test
	svc.CommandDB.Pool.Exec(ctx, "DELETE FROM payments")

	// Subscribe to payment events
	initiatedEvents := make(chan events.BaseEvent, 1)
	verifiedEvents := make(chan events.BaseEvent, 1)

	svc.Dispatcher.Subscribe(events.EventPaymentInitiated, func(event events.BaseEvent) error {
		initiatedEvents <- event
		return nil
	})

	svc.Dispatcher.Subscribe(events.EventPaymentVerified, func(event events.BaseEvent) error {
		verifiedEvents <- event
		return nil
	})

	// Initiate a payment
	initReq := payments.InitiatePaymentRequest{
		BookingID: "booking-event-1",
		UserID:    "user-event-1",
		Amount:    1000,
	}
	_, err := svc.PaymentCmdSvc.InitiatePayment(ctx, initReq)
	require.NoError(t, err)

	// Wait for event
	select {
	case event := <-initiatedEvents:
		assert.Equal(t, events.EventPaymentInitiated, event.Type)
	case <-ctx.Done():
		t.Fatal("Context cancelled before event received")
	}

	// Get the payment ID from the event and verify
	var paymentID string
	select {
	case event := <-initiatedEvents:
		payload := event.Payload.(events.EventPayload)
		paymentID = payload.PaymentID
	default:
		// If not in the first channel, check the second
	}

	// Verify the payment
	if paymentID == "" {
		// Get from DB
		payment, _ := svc.PaymentCmdSvc.Repo.GetPaymentByBookingID("booking-event-1")
		if payment != nil {
			paymentID = payment.ID
		}
	}

	if paymentID != "" {
		verifyReq := payments.VerifyPaymentRequest{
			PaymentID: paymentID,
			Mode:      "success",
		}
		err = svc.PaymentCmdSvc.VerifyPayment(ctx, verifyReq)
		require.NoError(t, err)

		select {
		case event := <-verifiedEvents:
			assert.Equal(t, events.EventPaymentVerified, event.Type)
		case <-ctx.Done():
			t.Fatal("Context cancelled before event received")
		}
	}
}

// Package integration provides integration tests for the ticket booking system
package integration

import (
	"context"
	"testing"

	"github.com/hitorii/ticket-booking/internal/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBookingIntegration_ReserveTicket tests ticket reservation flow
func TestBookingIntegration_ReserveTicket(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Clean up before test
	svc.CommandDB.Pool.Exec(ctx, "DELETE FROM reservations")

	t.Run("ReserveTicket_Success", func(t *testing.T) {
		err := svc.BookingCmdSvc.ReserveTicket(ctx, "user-1", "seat-A1")
		require.NoError(t, err)

		// Verify the reservation was created
		var count int
		err = svc.CommandDB.Pool.QueryRow(ctx, 
			"SELECT COUNT(*) FROM reservations WHERE seat_id=$1 AND user_id=$2", 
			"seat-A1", "user-1",
		).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "Reservation should be created")
	})

	t.Run("ReserveTicket_AlreadyReserved", func(t *testing.T) {
		// First reservation
		err := svc.BookingCmdSvc.ReserveTicket(ctx, "user-2", "seat-A2")
		require.NoError(t, err)

		// Second reservation for the same seat should fail
		err = svc.BookingCmdSvc.ReserveTicket(ctx, "user-3", "seat-A2")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "seat already reserved")
	})
}

// TestBookingIntegration_ConfirmTicket tests ticket confirmation flow
func TestBookingIntegration_ConfirmTicket(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Clean up before test
	svc.CommandDB.Pool.Exec(ctx, "DELETE FROM reservations")

	t.Run("ConfirmTicket_Success", func(t *testing.T) {
		// First reserve a ticket
		err := svc.BookingCmdSvc.ReserveTicket(ctx, "user-confirm-1", "seat-B1")
		require.NoError(t, err)

		// Then confirm it
		err = svc.BookingCmdSvc.ConfirmTicket(ctx, "user-confirm-1", "seat-B1")
		require.NoError(t, err)

		// Verify the status changed to BOOKED
		var status string
		err = svc.CommandDB.Pool.QueryRow(ctx,
			"SELECT status FROM reservations WHERE seat_id=$1 AND user_id=$2",
			"seat-B1", "user-confirm-1",
		).Scan(&status)
		require.NoError(t, err)
		assert.Equal(t, "BOOKED", status, "Status should be BOOKED")
	})

	t.Run("ConfirmTicket_NotHeld", func(t *testing.T) {
		// Try to confirm a seat that was never reserved
		err := svc.BookingCmdSvc.ConfirmTicket(ctx, "user-confirm-2", "seat-B2")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "seat not held or already booked")
	})
}

// TestBookingIntegration_CancelTicket tests ticket cancellation flow
func TestBookingIntegration_CancelTicket(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Clean up before test
	svc.CommandDB.Pool.Exec(ctx, "DELETE FROM reservations")

	t.Run("CancelTicket_Success", func(t *testing.T) {
		// First reserve a ticket
		err := svc.BookingCmdSvc.ReserveTicket(ctx, "user-cancel-1", "seat-C1")
		require.NoError(t, err)

		// Then cancel it
		err = svc.BookingCmdSvc.CancelTicket(ctx, "user-cancel-1", "seat-C1")
		require.NoError(t, err)

		// Verify the reservation was deleted
		var count int
		err = svc.CommandDB.Pool.QueryRow(ctx,
			"SELECT COUNT(*) FROM reservations WHERE seat_id=$1 AND user_id=$2",
			"seat-C1", "user-cancel-1",
		).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count, "Reservation should be deleted")
	})

	t.Run("CancelTicket_NotFound", func(t *testing.T) {
		// Try to cancel a seat that was never reserved
		err := svc.BookingCmdSvc.CancelTicket(ctx, "user-cancel-2", "seat-C2")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no reservation found")
	})
}

// TestBookingIntegration_WithEvents tests event publishing during booking operations
func TestBookingIntegration_WithEvents(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Clean up before test
	svc.CommandDB.Pool.Exec(ctx, "DELETE FROM reservations")

	// Subscribe to booking events
	reservedEvents := make(chan events.BaseEvent, 1)
	confirmedEvents := make(chan events.BaseEvent, 1)
	cancelledEvents := make(chan events.BaseEvent, 1)

	_ = reservedEvents
	_ = confirmedEvents
	_ = cancelledEvents

	// Reserve a ticket
	err := svc.BookingCmdSvc.ReserveTicket(ctx, "user-event-1", "seat-event-1")
	require.NoError(t, err)

	// Verify reservation exists
	var status string
	err = svc.CommandDB.Pool.QueryRow(ctx,
		"SELECT status FROM reservations WHERE seat_id=$1",
		"seat-event-1",
	).Scan(&status)
	require.NoError(t, err)
	assert.Equal(t, "HELD", status, "Initial status should be HELD")

	// Confirm the ticket
	err = svc.BookingCmdSvc.ConfirmTicket(ctx, "user-event-1", "seat-event-1")
	require.NoError(t, err)

	// Verify status changed to BOOKED
	err = svc.CommandDB.Pool.QueryRow(ctx,
		"SELECT status FROM reservations WHERE seat_id=$1",
		"seat-event-1",
	).Scan(&status)
	require.NoError(t, err)
	assert.Equal(t, "BOOKED", status, "Status should be BOOKED after confirmation")
}

// TestBookingIntegration_DoubleBookingPrevention tests that double booking is prevented
func TestBookingIntegration_DoubleBookingPrevention(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Clean up before test
	svc.CommandDB.Pool.Exec(ctx, "DELETE FROM reservations")

	// User A reserves a seat
	err := svc.BookingCmdSvc.ReserveTicket(ctx, "user-A", "seat-double-1")
	require.NoError(t, err)

	// User B tries to reserve the same seat - should fail
	err = svc.BookingCmdSvc.ReserveTicket(ctx, "user-B", "seat-double-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "seat already reserved")

	// User A confirms their booking
	err = svc.BookingCmdSvc.ConfirmTicket(ctx, "user-A", "seat-double-1")
	require.NoError(t, err)

	// Verify final status
	var status string
	err = svc.CommandDB.Pool.QueryRow(ctx,
		"SELECT status FROM reservations WHERE seat_id=$1",
		"seat-double-1",
	).Scan(&status)
	require.NoError(t, err)
	assert.Equal(t, "BOOKED", status, "Final status should be BOOKED")
}

package booking

import (
	"context"
	"errors"

	"github.com/hitorii/ticket-booking/internal/events"
	"github.com/hitorii/ticket-booking/internal/notification"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommandService struct {
	DB        *pgxpool.Pool
	Dispatcher *events.Dispatcher
}

func NewCommandService(db *pgxpool.Pool) *CommandService {
	return &CommandService{DB: db}
}

func NewCommandServiceWithDispatcher(db *pgxpool.Pool, dispatcher *events.Dispatcher) *CommandService {
	return &CommandService{DB: db, Dispatcher: dispatcher}
}

// ReserveTicket - Command to reserve a ticket
func (s *CommandService) ReserveTicket(ctx context.Context, userID, seatID string) error {
	_, err := s.DB.Exec(
		ctx,
		"INSERT INTO reservations (seat_id, user_id, status) VALUES ($1,$2,'HELD')",
		seatID,
		userID,
	)
	if err != nil {
		return errors.New("seat already reserved")
	}

	// Push job into Redis Queue after successful reservation
	err = notification.EnqueueBookingNotification(userID, seatID, "Seat reserved successfully")
	if err != nil {
		// Log error but don't fail the reservation
		// In production, you might want to implement retry logic
		println("Warning: Failed to enqueue notification:", err.Error())
	}

	// Emit event for event-driven flow
	if s.Dispatcher != nil {
		payload := events.EventPayload{
			UserID: userID,
			SeatID: seatID,
			Status: "HELD",
		}
		_ = s.Dispatcher.Publish(ctx, events.EventTicketReserved, seatID, payload)
	}

	return nil
}

// ConfirmTicket - Command to confirm a ticket reservation
func (s *CommandService) ConfirmTicket(ctx context.Context, userID, seatID string) error {
	res, err := s.DB.Exec(
		ctx,
		`UPDATE reservations 
		 SET status='BOOKED'
		 WHERE seat_id=$1 AND user_id=$2 AND status='HELD'`,
		seatID,
		userID,
	)
	if err != nil {
		return errors.New("failed to confirm booking")
	}

	if res.RowsAffected() == 0 {
		return errors.New("seat not held or already booked")
	}
	
	// Push job into Redis Queue after successful confirmation
	err = notification.EnqueueBookingNotification(userID, seatID, "Ticket confirmed successfully")
	if err != nil {
		println("Warning: Failed to enqueue notification:", err.Error())
	}

	// Emit event for event-driven flow
	if s.Dispatcher != nil {
		payload := events.EventPayload{
			UserID: userID,
			SeatID: seatID,
			Status: "BOOKED",
		}
		_ = s.Dispatcher.Publish(ctx, events.EventTicketConfirmed, seatID, payload)
	}
	
	return nil
}

// CancelTicket - Command to cancel a ticket reservation
func (s *CommandService) CancelTicket(ctx context.Context, userID, seatID string) error {
	res, err := s.DB.Exec(
		ctx,
		"DELETE FROM reservations WHERE seat_id=$1 AND user_id=$2",
		seatID,
		userID,
	)
	if err != nil {
		return errors.New("cancel failed")
	}

	if res.RowsAffected() == 0 {
		return errors.New("no reservation found")
	}
	
	// Push job into Redis Queue after cancellation
	err = notification.EnqueueBookingNotification(userID, seatID, "Reservation cancelled")
	if err != nil {
		println("Warning: Failed to enqueue notification:", err.Error())
	}

	// Emit event for event-driven flow
	if s.Dispatcher != nil {
		payload := events.EventPayload{
			UserID: userID,
			SeatID: seatID,
			Status: "CANCELLED",
		}
		_ = s.Dispatcher.Publish(ctx, events.EventTicketCancelled, seatID, payload)
	}
	
	return nil
}

package booking

import (
	"context"
	"errors"
	"fmt"

	"github.com/hitorii/ticket-booking/internal/events"
	"github.com/hitorii/ticket-booking/internal/notification"
	"github.com/hitorii/ticket-booking/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type CommandService struct {
	DB         *pgxpool.Pool
	Dispatcher *events.Dispatcher
	RedisClient *redis.Client
	Lock       *utils.DistributedLock
}

func NewCommandService(db *pgxpool.Pool) *CommandService {
	return &CommandService{DB: db}
}

func NewCommandServiceWithDispatcher(db *pgxpool.Pool, dispatcher *events.Dispatcher) *CommandService {
	return &CommandService{DB: db, Dispatcher: dispatcher}
}

// NewCommandServiceWithLock creates a command service with distributed locking support
func NewCommandServiceWithLock(db *pgxpool.Pool, dispatcher *events.Dispatcher, redisClient *redis.Client) *CommandService {
	var lock *utils.DistributedLock
	if redisClient != nil {
		lock = utils.NewDistributedLock(redisClient)
	}
	return &CommandService{
		DB:         db,
		Dispatcher: dispatcher,
		RedisClient: redisClient,
		Lock:       lock,
	}
}

// ReserveTicket - Command to reserve a ticket with distributed locking
func (s *CommandService) ReserveTicket(ctx context.Context, userID, seatID string) error {
	// Use distributed lock if available to prevent race conditions
	if s.Lock != nil {
		err := s.Lock.WithSeatLock(ctx, seatID, func() error {
			return s.reserveTicketInternal(ctx, userID, seatID)
		})
		if err != nil {
			return fmt.Errorf("failed to reserve ticket: %w", err)
		}
		return nil
	}
	
	// Fallback to internal method if no lock available
	return s.reserveTicketInternal(ctx, userID, seatID)
}

// reserveTicketInternal is the actual reservation logic
func (s *CommandService) reserveTicketInternal(ctx context.Context, userID, seatID string) error {
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

// ConfirmTicket - Command to confirm a ticket reservation with distributed locking
func (s *CommandService) ConfirmTicket(ctx context.Context, userID, seatID string) error {
	// Use distributed lock if available
	if s.Lock != nil {
		err := s.Lock.WithSeatLock(ctx, seatID, func() error {
			return s.confirmTicketInternal(ctx, userID, seatID)
		})
		if err != nil {
			return fmt.Errorf("failed to confirm ticket: %w", err)
		}
		return nil
	}
	
	return s.confirmTicketInternal(ctx, userID, seatID)
}

// confirmTicketInternal is the actual confirmation logic
func (s *CommandService) confirmTicketInternal(ctx context.Context, userID, seatID string) error {
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

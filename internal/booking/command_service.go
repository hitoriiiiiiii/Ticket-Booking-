package booking

import (
	"context"
	"errors"

	"github.com/hitorii/ticket-booking/internal/notification"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommandService struct {
	DB *pgxpool.Pool
}

func NewCommandService(db *pgxpool.Pool) *CommandService {
	return &CommandService{DB: db}
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

	// Push job into NotificationQueue after successful reservation
	notification.NotificationQueue <- notification.Job{
		Type:    notification.JobTypeBooking,
		UserID:  userID,
		Message: "Seat reserved successfully",
		Data:    map[string]interface{}{"seatID": seatID},
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
	return nil
}

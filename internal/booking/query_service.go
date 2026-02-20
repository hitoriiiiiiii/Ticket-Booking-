package booking

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryService struct {
	DB *pgxpool.Pool
}

func NewQueryService(db *pgxpool.Pool) *QueryService {
	return &QueryService{DB: db}
}

// SeatStatus represents the status of a seat
type SeatStatus struct {
	SeatID string `json:"seat_id"`
	Status string `json:"status"`
}

// GetAvailability - Query to check seat availability
func (s *QueryService) GetAvailability(ctx context.Context, seatID string) (*SeatStatus, error) {
	var status string

	err := s.DB.QueryRow(
		ctx,
		"SELECT status FROM reservations WHERE seat_id=$1",
		seatID,
	).Scan(&status)

	if err != nil {
		// If no row found, seat is available
		return &SeatStatus{
			SeatID: seatID,
			Status: "AVAILABLE",
		}, nil
	}

	return &SeatStatus{
		SeatID: seatID,
		Status: status,
	}, nil
}

// Event represents an event from the events table
type Event struct {
	SeatID    string `json:"seat_id"`
	Type      string `json:"type"`
	Payload   string `json:"payload"`
	CreatedAt string `json:"time"`
}

// GetEvents - Query to get all events
func (s *QueryService) GetEvents(ctx context.Context) ([]Event, error) {
	rows, err := s.DB.Query(
		ctx,
		"SELECT aggregate_id, event_type, payload, created_at FROM events ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventsList []Event
	for rows.Next() {
		var aggID, eventType string
		var payload []byte
		var createdAt string

		rows.Scan(&aggID, &eventType, &payload, &createdAt)

		eventsList = append(eventsList, Event{
			SeatID:    aggID,
			Type:      eventType,
			Payload:   string(payload),
			CreatedAt: createdAt,
		})
	}

	return eventsList, nil
}

// Reservation represents a user's reservation
type Reservation struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	SeatID string `json:"seat_id"`
	Status string `json:"status"`
}

// GetUserReservations - Query to get all reservations for a user
func (s *QueryService) GetUserReservations(ctx context.Context, userID string) ([]Reservation, error) {
	rows, err := s.DB.Query(
		ctx,
		"SELECT id, user_id, seat_id, status FROM reservations WHERE user_id=$1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []Reservation
	for rows.Next() {
		var r Reservation
		err := rows.Scan(&r.ID, &r.UserID, &r.SeatID, &r.Status)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}

	return reservations, nil
}

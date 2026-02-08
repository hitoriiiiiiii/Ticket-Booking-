//Append + load events from the queue to the event store (Postgress) for durability and later processing by projections

package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	DB *pgxpool.Pool
}

func (s *Store) Append(eventType string, aggregateID string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(
		context.Background(),
		`INSERT INTO events (aggregate_id, event_type, payload)
		 VALUES ($1, $2, $3)`,
		aggregateID, eventType, data,
	)

	if err != nil {
		log.Printf("❌ Failed to append event: %v", err)
		return err
	}

	log.Printf("✅ Appended event: %s for aggregate: %s", eventType, aggregateID)
	return nil
}

func (s *Store) IsSeatReserved(seatID string) (bool, error){
	var count int

	err := s.DB.QueryRow(
		context.Background(),
		`SELECT COUNT(*) FROM events WHERE event_type = 'TicketReserved' AND aggregate_id = $1`,
		seatID,
	).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) IsSeatAvailable(seatID string) (bool, error) {

	var reserved bool
	var cancelled bool

	// Check reserved event exists
	err := s.DB.QueryRow(
		context.Background(),
		"SELECT EXISTS (SELECT 1 FROM events WHERE aggregate_id=$1 AND event_type='TicketReserved')",
		seatID,
	).Scan(&reserved)

	if err != nil {
		return false, err
	}

	// Check cancelled event exists
	err = s.DB.QueryRow(
		context.Background(),
		"SELECT EXISTS (SELECT 1 FROM events WHERE aggregate_id=$1 AND event_type='TicketCancelled')",
		seatID,
	).Scan(&cancelled)

	if err != nil {
		return false, err
	}

	// Seat is available if not reserved OR cancelled exists
	if !reserved || cancelled {
		return true, nil
	}

	return false, nil
}

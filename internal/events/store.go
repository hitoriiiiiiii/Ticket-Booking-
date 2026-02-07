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
		 VALUES (gen_random_uuid(), $1, $2)`,
		eventType, data,
	)

	if err != nil {
		log.Printf("❌ Failed to append event: %v", err)
		return err
	}

	log.Printf("✅ Appended event: %s for aggregate: %s", eventType, aggregateID)
	return nil
}

package booking

// Real model updater (CQRS projection)

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hitorii/ticket-booking/internal/events"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationProjection struct {
	DB         *pgxpool.Pool
	EventStore *events.Store
}

func NewReservationProjection(db *pgxpool.Pool, eventStore *events.Store) *ReservationProjection {
	return &ReservationProjection{DB: db, EventStore: eventStore}
}

func (p *ReservationProjection) Run() {
	// Wait a bit for database to be ready
	time.Sleep(5 * time.Second)

	// Get last processed event ID
	var lastEventID *string
	err := p.DB.QueryRow(context.Background(),
		"SELECT last_event_id FROM projection_state WHERE id = 1",
	).Scan(&lastEventID)
	if err != nil {
		log.Println("Failed to get projection state:", err)
		return
	}

	// Build query for incremental processing
	query := "SELECT id, aggregate_id, event_type, payload FROM events WHERE id > $1 ORDER BY id"
	var rows pgx.Rows
	if lastEventID == nil {
		rows, err = p.DB.Query(context.Background(), query, "")
	} else {
		rows, err = p.DB.Query(context.Background(), query, *lastEventID)
	}
	if err != nil {
		log.Fatal("Failed to query events:", err)
	}
	defer rows.Close()

	var lastProcessedID string
	for rows.Next() {
		var eventID, seatID, eventType string
		var payload []byte

		err := rows.Scan(&eventID, &seatID, &eventType, &payload)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}

		if eventType == "TicketReserved" {
			var data struct {
				UserID string `json:"user_id"`
				SeatID string `json:"seat_id"`
			}
			json.Unmarshal(payload, &data)

			// Insert or update reservation in projection
			_, err := p.DB.Exec(context.Background(),
				`INSERT INTO reservation_projection(seat_id, user_id, status)
				 VALUES ($1, $2, 'HELD')
				 ON CONFLICT (seat_id) DO UPDATE SET
				 user_id = EXCLUDED.user_id,
				 status = 'HELD',
				 updated_at = NOW()`,
				data.SeatID, data.UserID,
			)
			if err != nil {
				log.Println("Projection failed for reserve:", err)
			} else {
				log.Println("Seat reserved in projection:", data.SeatID)
			}

		} else if eventType == "TicketConfirmed" {
			var data struct {
				UserID string `json:"user_id"`
				SeatID string `json:"seat_id"`
			}
			json.Unmarshal(payload, &data)

			// Update reservation status to confirmed in projection
			_, err := p.DB.Exec(context.Background(),
				`UPDATE reservation_projection SET status = 'BOOKED', updated_at = NOW() WHERE seat_id = $1`,
				data.SeatID,
			)
			if err != nil {
				log.Println("Projection failed for confirm:", err)
			} else {
				log.Println("Seat confirmed in projection:", data.SeatID)
			}

		} else if eventType == "TicketCancelled" {
			var data struct {
				UserID string `json:"user_id"`
				SeatID string `json:"seat_id"`
			}
			json.Unmarshal(payload, &data)

			// Update reservation status to cancelled in projection
			_, err := p.DB.Exec(context.Background(),
				`UPDATE reservation_projection SET status = 'CANCELLED', updated_at = NOW() WHERE seat_id = $1`,
				data.SeatID,
			)
			if err != nil {
				log.Println("Projection failed for cancel:", err)
			} else {
				log.Println("Seat cancelled in projection:", data.SeatID)
			}
		}

		lastProcessedID = eventID
	}

	// Update projection state if we processed events
	if lastProcessedID != "" {
		_, err := p.DB.Exec(context.Background(),
			"UPDATE projection_state SET last_event_id = $1 WHERE id = 1",
			lastProcessedID,
		)
		if err != nil {
			log.Println("Failed to update projection state:", err)
		}
	}

	log.Println("Projection run completed")
}

// StartProjectionWorker starts the projection worker in a goroutine
func StartProjectionWorker(db *pgxpool.Pool, eventStore *events.Store) {
	projection := NewReservationProjection(db, eventStore)
	go func() {
		for {
			projection.Run()
			// Sleep before next run
			time.Sleep(10 * time.Second)
		}
	}()
}

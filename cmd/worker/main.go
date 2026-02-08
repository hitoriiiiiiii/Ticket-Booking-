package main

import (
	"log"
	"time"

	"github.com/hitorii/ticket-booking/internal/booking"
	"github.com/hitorii/ticket-booking/internal/config"
	"github.com/hitorii/ticket-booking/internal/db"
)

func main() {
	cfg := config.Load()

	pool := db.Connect(cfg.DatabaseURL)
	defer pool.Close()

	projection := &booking.ReservationProjection{DB: pool}

	log.Println("Starting projection worker...")

	for {
		projection.Run()
		time.Sleep(5 * time.Second) // Run projection every 5 seconds
	}
}

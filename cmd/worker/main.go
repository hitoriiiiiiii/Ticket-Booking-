package main

import (
	"log"
	"os"
	"time"

	"github.com/hitorii/ticket-booking/internal/booking"
	"github.com/hitorii/ticket-booking/internal/config"
	"github.com/hitorii/ticket-booking/internal/db"
	"github.com/hitorii/ticket-booking/internal/notification"
)

func main() {
	cfg := config.Load()

	// Load .env
	_ = os.Getenv("COMMAND_DATABASE_URL")

	// Connect to databases
	cmdDB := db.Connect(os.Getenv("COMMAND_DATABASE_URL"))
	queryDB := db.Connect(os.Getenv("QUERY_DATABASE_URL"))
	defer cmdDB.Close()
	defer queryDB.Close()

	// Start notification worker
	notificationRepo := notification.NewRepository(queryDB)
	notification.StartWorker(notificationRepo)
	log.Println("🔔 Notification worker started")

	// Start projection worker
	projection := &booking.ReservationProjection{
		DB:         queryDB,
		EventStore: &booking.EventStore{DB: cmdDB},
	}

	log.Println("📊 Starting projection worker...")

	for {
		projection.Run()
		time.Sleep(5 * time.Second) // Run projection every 5 seconds
	}
}

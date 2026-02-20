package main

import (
	"log"
	"os"
	"time"

	"github.com/hitorii/ticket-booking/internal/booking"
	"github.com/hitorii/ticket-booking/internal/config"
	"github.com/hitorii/ticket-booking/internal/db"
	"github.com/hitorii/ticket-booking/internal/notification"
	"github.com/hitorii/ticket-booking/internal/queue"
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

	// Initialize Redis for job queue
	log.Println("🔌 Connecting to Redis for job queue...")
	if err := queue.InitRedis(cfg.RedisURL); err != nil {
		log.Printf("❌ Failed to connect to Redis: %v", err)
		log.Println("⚠️  Worker will run in fallback mode with in-memory queue")
	} else {
		log.Println("✅ Redis connected successfully")
	}
	defer queue.Close()

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

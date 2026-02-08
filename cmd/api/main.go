package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/hitorii/ticket-booking/internal/booking"
	"github.com/hitorii/ticket-booking/internal/db"
	"github.com/hitorii/ticket-booking/internal/events"
)

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found, using system env variables")
	}

	r := gin.Default()

	// Read DB config from env
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Build Database URL
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	// Connect DB pool
	pool := db.Connect(databaseURL)

	// Event Store
	eventStore := &events.Store{DB: pool}

	// Booking Handler
	bookingHandler := &booking.Handler{EventStore: eventStore}

	//Projection
	projection := &booking.ReservationProjection{DB: pool}
	go projection.Run()

	// Routes
	r.POST("/reserve", bookingHandler.ReserveTicket)
	r.GET("/health", booking.HealthCheck)
    r.GET("/events", bookingHandler.GetEvents)
	r.POST("/cancel", bookingHandler.CancelTicket)
	r.GET("/availability/:seat_id", bookingHandler.CheckAvailability)
	r.POST("/confirm", bookingHandler.ConfirmTicket)

	// Server Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Server running on port", port)
	r.Run(":" + port)
}

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
	"github.com/hitorii/ticket-booking/internal/movie"
	"github.com/hitorii/ticket-booking/internal/payments"
	"github.com/hitorii/ticket-booking/internal/show"
	"github.com/hitorii/ticket-booking/internal/user"
	"github.com/hitorii/ticket-booking/internal/notifications"
)

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found, using system env variables")
	}

	r := gin.Default()

	// Read Database URL from env, or build from individual components
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")

		// Build Database URL
		databaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser,
			dbPassword,
			dbHost,
			dbPort,
			dbName,
		)
	}
	log.Println("Using DATABASE_URL:", databaseURL)

	// Connect DB pool
	pool := db.Connect(databaseURL)

	// Event Store
	eventStore := &events.Store{DB: pool}

	// Booking Handler
	bookingHandler := &booking.Handler{EventStore: eventStore}

	// User Handler
	userHandler := &user.Handler{DB: pool}

	// Movie Handler
	movieHandler := &movie.Handler{DB: pool}

	// Payment Handler
	paymentRepo := payments.NewRepository(pool)
	paymentService := payments.NewService(paymentRepo)
	paymentHandler := payments.NewHandler(paymentService)

	// Shows Handler
	showService := show.NewService(pool)
	showHandler := show.NewHandler(showService)

	//Projection
	projection := &booking.ReservationProjection{DB: pool}
	go projection.Run()

	// Routes
	r.POST("/reserve", bookingHandler.ReserveTicket)
	r.GET("/health", booking.HealthCheck)
	r.GET("/events", bookingHandler.GetEvents) // Admin endpoint to view events
	//booking
	r.POST("/cancel", bookingHandler.CancelTicket)
	r.GET("/availability/:seat_id", bookingHandler.CheckAvailability)
	r.POST("/confirm", bookingHandler.ConfirmTicket)
	//user
	r.POST("/users/register", userHandler.Register)
	r.POST("/users/login", userHandler.Login)
	r.GET("/users", userHandler.ListUsers)
	//movie
	r.GET("/movies", movieHandler.GetMovies)
	r.GET("/movies/:id", movieHandler.GetMovie)
	r.POST("/movies", movieHandler.CreateMovie)
	//shows
	r.GET("/shows", showHandler.GetShows)
	r.POST("/shows", showHandler.CreateShow)
	//payments
	r.POST("/payments/initiate", paymentHandler.InitiatePayment)
	r.POST("/payments/verify", paymentHandler.VerifyPayment)
	//notifications
	notificationRepo := notifications.NewRepository(pool)
	notificationHandler := notifications.NewHandler(notificationRepo)
	r.GET("/notifications/:user_id", notificationHandler.GetUserNotifications)

	// Server Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Server running on port", port)
	r.Run(":" + port)
}

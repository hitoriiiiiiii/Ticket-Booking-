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
	"github.com/hitorii/ticket-booking/internal/notification"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found, using system env variables")
	}

	r := gin.Default()

	// Connect to Command DB and Query DB
	cmdDB := db.Connect(os.Getenv("COMMAND_DATABASE_URL"))
	queryDB := db.Connect(os.Getenv("QUERY_DATABASE_URL"))

	// Event Store for commands
	eventStore := &events.Store{DB: cmdDB}

	
	// Booking
	bookingHandler := &booking.Handler{
		CmdDB:   cmdDB,
		QueryDB: queryDB,
		EventStore: eventStore,
	}

	// User
	userHandler := &user.Handler{
		CmdDB:   cmdDB,
		QueryDB: queryDB,
	}

	// Movie
	movieHandler := &movie.Handler{
		CmdDB:   cmdDB,
		QueryDB: queryDB,
	}

	// Show
	showHandler := &show.Handler{
		CmdDB:   cmdDB,
		QueryDB: queryDB,
	}

	// Payments
	paymentRepo := payments.NewRepository(cmdDB)
	paymentService := payments.NewService(paymentRepo)
	paymentHandler := payments.NewHandler(paymentService)

	// Notifications
	notificationRepo := notification.NewRepository(queryDB)
	notificationHandler := notification.NewHandler(notificationRepo)

	// COMMAND APIs (writes)
	r.POST("/cmd/reserve", bookingHandler.ReserveTicket)
	r.POST("/cmd/cancel", bookingHandler.CancelTicket)
	r.POST("/cmd/confirm", bookingHandler.ConfirmTicket)
	r.POST("/cmd/users/register", userHandler.Register)
	r.POST("/cmd/movies", movieHandler.CreateMovie)
	r.POST("/cmd/shows", showHandler.CreateShow)
	r.POST("/cmd/payments/initiate", paymentHandler.InitiatePayment)
	r.POST("/cmd/payments/verify", paymentHandler.VerifyPayment)

	// QUERY APIs (reads)
	r.GET("/query/reservations/:user_id", bookingHandler.GetUserReservations)
	r.GET("/query/availability/:seat_id", bookingHandler.CheckAvailability)
	r.GET("/query/events", bookingHandler.GetEvents)
	r.GET("/query/users", userHandler.ListUsers)
	r.GET("/query/movies", movieHandler.GetMovies)
	r.GET("/query/movies/:id", movieHandler.GetMovie)
	r.GET("/query/shows", showHandler.GetShows)
	r.GET("/query/notifications/:user_id", notificationHandler.GetUserNotifications)

	// Health
	r.GET("/health", booking.HealthCheck)

	//    Projection 
	projection := &booking.ReservationProjection{
		DB:      queryDB,
		EventStore: eventStore,
	}
	go projection.Run() // listens for events and updates query DB

	// Start Server 
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Server running on port", port)
	r.Run(":" + port)
}

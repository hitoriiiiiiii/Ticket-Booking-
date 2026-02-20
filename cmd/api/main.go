package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/hitorii/ticket-booking/internal/booking"
	"github.com/hitorii/ticket-booking/internal/config"
	"github.com/hitorii/ticket-booking/internal/db"
	"github.com/hitorii/ticket-booking/internal/events"
	"github.com/hitorii/ticket-booking/internal/movie"
	"github.com/hitorii/ticket-booking/internal/payments"
	"github.com/hitorii/ticket-booking/internal/queue"
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

	// Load config
	cfg := config.Load()

	// Connect to Command DB and Query DB
	cmdDB := db.Connect(os.Getenv("COMMAND_DATABASE_URL"))
	queryDB := db.Connect(os.Getenv("QUERY_DATABASE_URL"))

	// Initialize Redis for job queue
	log.Println("🔌 Connecting to Redis for job queue...")
	if err := queue.InitRedis(cfg.RedisURL); err != nil {
		log.Printf("❌ Failed to connect to Redis: %v", err)
		log.Println("⚠️  Job queue will not work - please ensure Redis is running")
	} else {
		log.Println("✅ Redis connected successfully")
	}
	defer queue.Close()

	// Event Store for commands
	eventStore := &events.Store{DB: cmdDB}

	// Booking - Command Handler for write operations (CQRS pattern)
	bookingCommandService := booking.NewCommandService(cmdDB)
	bookingCommandHandler := booking.NewCommandHandler(bookingCommandService, eventStore)
	// Booking - Query handler for read operations
	bookingQueryService := booking.NewQueryService(queryDB)
	bookingQueryHandler := booking.NewQueryHandler(bookingQueryService)

	// User - CQRS pattern
	userCmdService := user.NewCommandService(cmdDB)
	userQueryService := user.NewQueryService(queryDB)
	userCommandHandler := user.NewCommandHandler(userCmdService)
	userQueryHandler := user.NewQueryHandler(userQueryService)


	// Movie - CQRS pattern
	movieCmdService := movie.NewCommandService(cmdDB)
	movieQueryService := movie.NewQueryService(queryDB)
	movieCommandHandler := movie.NewCommandHandler(movieCmdService)
	movieQueryHandler := movie.NewQueryHandler(movieQueryService)


	// Show - CQRS pattern
	showCmdService := show.NewCommandService(cmdDB)
	showQueryService := show.NewQueryService(queryDB)
	showCommandHandler := show.NewCommandHandler(showCmdService)
	showQueryHandler := show.NewQueryHandler(showQueryService)

	// Payments - CQRS pattern
	paymentRepo := payments.NewRepository(cmdDB)
	paymentCmdService := payments.NewCommandService(paymentRepo)
	paymentCommandHandler := payments.NewCommandHandler(paymentCmdService)
	// For payments, we can also create query handler if needed
	paymentQueryService := payments.NewQueryService(queryDB)
	paymentQueryHandler := payments.NewQueryHandler(paymentQueryService)

	// Notifications - CQRS pattern
	notificationQueryService := notification.NewQueryService(queryDB)
	notificationQueryHandler := notification.NewQueryHandler(notificationQueryService)

	
	
	// COMMAND APIs (writes)
	r.POST("/cmd/reserve", bookingCommandHandler.ReserveTicket)
	r.POST("/cmd/cancel", bookingCommandHandler.CancelTicket)
	r.POST("/cmd/confirm", bookingCommandHandler.ConfirmTicket)
	r.POST("/cmd/users/register", userCommandHandler.Register)
	r.POST("/cmd/movies", movieCommandHandler.CreateMovie)
	r.POST("/cmd/shows", showCommandHandler.CreateShow)
	r.POST("/cmd/payments/initiate", paymentCommandHandler.InitiatePayment)
	r.POST("/cmd/payments/verify", paymentCommandHandler.VerifyPayment)

	// QUERY APIs (reads)
	r.GET("/query/reservations/:user_id", bookingQueryHandler.GetUserReservations)
	r.GET("/query/availability/:seat_id", bookingQueryHandler.CheckAvailability)
	r.GET("/query/events", bookingQueryHandler.GetEvents)
	r.GET("/query/users", userQueryHandler.ListUsers)
	r.GET("/query/movies", movieQueryHandler.GetMovies)
	r.GET("/query/movies/:id", movieQueryHandler.GetMovie)
	r.GET("/query/shows", showQueryHandler.GetShows)
	r.GET("/query/notifications/:user_id", notificationQueryHandler.GetUserNotifications)

	// Health
	r.GET("/health", booking.HealthCheck)

	// Projection 
	projection := &booking.ReservationProjection{
		DB:         queryDB,
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

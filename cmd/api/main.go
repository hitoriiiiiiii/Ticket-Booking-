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
	err := godotenv.Load()
	if err != nil {
		log.Println("âš ï¸ No .env file found, using system env variables")
	}

	r := gin.Default()
	cfg := config.Load()
	cmdDB := db.Connect(os.Getenv("COMMAND_DATABASE_URL"))
	queryDB := db.Connect(os.Getenv("QUERY_DATABASE_URL"))

	log.Println("ðŸ”Œ Connecting to Redis for job queue...")
	if err := queue.InitRedis(cfg.RedisURL); err != nil {
		log.Printf("âŒ Failed to connect to Redis: %v", err)
		log.Println("âš ï¸  Job queue will not work")
	} else {
		log.Println("âœ… Redis connected successfully")
	}
	defer queue.Close()

	eventStore := &events.Store{DB: cmdDB}

	var eventDispatcher *events.Dispatcher
	if queue.RedisClient != nil {
		eventDispatcher = events.NewDispatcher(queue.RedisClient, eventStore)
		log.Println("ðŸ“¢ Event dispatcher initialized")
	}

	bookingCommandService := booking.NewCommandServiceWithDispatcher(cmdDB, eventDispatcher)
	bookingCommandHandler := booking.NewCommandHandler(bookingCommandService, eventStore)
	bookingQueryService := booking.NewQueryService(queryDB)
	bookingQueryHandler := booking.NewQueryHandler(bookingQueryService)

	userCmdService := user.NewCommandServiceWithDispatcher(cmdDB, eventDispatcher)
	userQueryService := user.NewQueryService(queryDB)
	userCommandHandler := user.NewCommandHandler(userCmdService)
	userQueryHandler := user.NewQueryHandler(userQueryService)

	movieCmdService := movie.NewCommandServiceWithDispatcher(eventDispatcher)
	movieQueryService := movie.NewQueryService(queryDB)
	movieCommandHandler := movie.NewCommandHandler(movieCmdService)
	movieQueryHandler := movie.NewQueryHandler(movieQueryService)

	showCmdService := show.NewCommandServiceWithDispatcher(eventDispatcher)
	showQueryService := show.NewQueryService(queryDB)
	showCommandHandler := show.NewCommandHandler(showCmdService)
	showQueryHandler := show.NewQueryHandler(showQueryService)

	paymentRepo := payments.NewRepository(cmdDB)
	paymentCmdService := payments.NewCommandServiceWithDispatcher(paymentRepo, eventDispatcher)
	paymentCommandHandler := payments.NewCommandHandler(paymentCmdService)
	paymentQueryService := payments.NewQueryService(queryDB)
	paymentQueryHandler := payments.NewQueryHandler(paymentQueryService)

	notificationQueryService := notification.NewQueryService(queryDB)
	notificationQueryHandler := notification.NewQueryHandler(notificationQueryService)

	if eventDispatcher != nil {
		setupEventSubscribers(eventDispatcher)
	}

	r.POST("/cmd/reserve", bookingCommandHandler.ReserveTicket)
	r.POST("/cmd/cancel", bookingCommandHandler.CancelTicket)
	r.POST("/cmd/confirm", bookingCommandHandler.ConfirmTicket)
	r.POST("/cmd/users/register", userCommandHandler.Register)
	r.POST("/cmd/movies", movieCommandHandler.CreateMovie)
	r.POST("/cmd/shows", showCommandHandler.CreateShow)
	r.POST("/cmd/payments/initiate", paymentCommandHandler.InitiatePayment)
	r.POST("/cmd/payments/verify", paymentCommandHandler.VerifyPayment)

	r.GET("/query/reservations/:user_id", bookingQueryHandler.GetUserReservations)
	r.GET("/query/availability/:seat_id", bookingQueryHandler.CheckAvailability)
	r.GET("/query/events", bookingQueryHandler.GetEvents)
	r.GET("/query/users", userQueryHandler.ListUsers)
	r.GET("/query/movies", movieQueryHandler.GetMovies)
	r.GET("/query/movies/:id", movieQueryHandler.GetMovie)
	r.GET("/query/shows", showQueryHandler.GetShows)
	r.GET("/query/notifications/:user_id", notificationQueryHandler.GetUserNotifications)

	r.GET("/health", booking.HealthCheck)

	projection := &booking.ReservationProjection{
		DB:         queryDB,
		EventStore: eventStore,
	}
	go projection.Run()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("ðŸš€ Server running on port", port)
	r.Run(":" + port)
}

func setupEventSubscribers(dispatcher *events.Dispatcher) {
	dispatcher.Subscribe(events.EventTicketReserved, func(event events.BaseEvent) error {
		log.Printf("ðŸ“¥ Event received: %s for aggregate: %s", event.Type, event.AggregateID)
		return nil
	})
	dispatcher.Subscribe(events.EventTicketConfirmed, func(event events.BaseEvent) error {
		log.Printf("ðŸ“¥ Event received: %s for aggregate: %s", event.Type, event.AggregateID)
		return nil
	})
	dispatcher.Subscribe(events.EventTicketCancelled, func(event events.BaseEvent) error {
		log.Printf("ðŸ“¥ Event received: %s for aggregate: %s", event.Type, event.AggregateID)
		return nil
	})
	dispatcher.Subscribe(events.EventPaymentVerified, func(event events.BaseEvent) error {
		log.Printf("ðŸ“¥ Event received: %s for aggregate: %s", event.Type, event.AggregateID)
		return nil
	})
	dispatcher.Subscribe(events.EventUserRegistered, func(event events.BaseEvent) error {
		log.Printf("ðŸ“¥ Event received: %s for aggregate: %s", event.Type, event.AggregateID)
		return nil
	})
	log.Println("âœ… Event subscribers configured")
}

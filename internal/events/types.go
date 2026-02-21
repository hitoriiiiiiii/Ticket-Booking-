package events

import (
	"time"
)

// Event types for the event-driven system
const (
	// Booking events
	EventTicketReserved   = "TicketReserved"
	EventTicketConfirmed  = "TicketConfirmed"
	EventTicketCancelled  = "TicketCancelled"
	
	// Payment events
	EventPaymentInitiated = "PaymentInitiated"
	EventPaymentVerified  = "PaymentVerified"
	EventPaymentRefunded  = "PaymentRefunded"
	
	// User events
	EventUserRegistered   = "UserRegistered"
	EventUserUpdated      = "UserUpdated"
	EventUserDeleted      = "UserDeleted"
	
	// Show events
	EventShowCreated      = "ShowCreated"
	EventShowUpdated      = "ShowUpdated"
	EventShowDeleted      = "ShowDeleted"
	
	// Movie events
	EventMovieCreated     = "MovieCreated"
	EventMovieUpdated     = "MovieUpdated"
	EventMovieDeleted     = "MovieDeleted"
)

// BaseEvent represents a base event structure
type BaseEvent struct {
	Type        string      `json:"type"`
	AggregateID string      `json:"aggregate_id"`
	Timestamp   time.Time   `json:"timestamp"`
	Payload     interface{} `json:"payload"`
}

// EventPayload is the payload for events
type EventPayload struct {
	// Common fields
	UserID    string `json:"user_id,omitempty"`
	BookingID string `json:"booking_id,omitempty"`
	SeatID    string `json:"seat_id,omitempty"`
	
	// User specific
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	
	// Payment specific
	PaymentID string `json:"payment_id,omitempty"`
	Amount    int    `json:"amount,omitempty"`
	Status    string `json:"status,omitempty"`
	Mode      string `json:"mode,omitempty"`
	
	// Show/Movie specific
	ShowID    string `json:"show_id,omitempty"`
	MovieID   string `json:"movie_id,omitempty"`
	Name      string `json:"name,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	Price     int    `json:"price,omitempty"`
	
	// Admin specific
	IsAdmin bool `json:"is_admin,omitempty"`
}

// EventHandler is a function type for handling events
type EventHandler func(event BaseEvent) error

// EventSubscriber subscribes to specific event types
type EventSubscriber struct {
	EventType string
	Handler   EventHandler
}

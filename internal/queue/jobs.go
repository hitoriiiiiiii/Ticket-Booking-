// Job definitions and handlers for background processing
package queue

import "time"

// Job types
const (
	JobTypeNotification = "notification"
	JobTypeEmail        = "email"
	JobTypePayment      = "payment"
	JobTypeBooking      = "booking"
)

// JobPayload represents a job that can be serialized to Redis
type JobPayload struct {
	Type      string                 `json:"type"`
	UserID    string                 `json:"user_id"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
	Priority  int                    `json:"priority"`
	CreatedAt time.Time              `json:"created_at"`
}

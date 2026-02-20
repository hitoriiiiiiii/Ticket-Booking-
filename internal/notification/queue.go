// Queues
package notification

import (
	"context"
	"time"

	"github.com/hitorii/ticket-booking/internal/queue"
)

// Job types
const (
	JobTypeNotification = "notification"
	JobTypeEmail        = "email"
	JobTypePayment      = "payment"
	JobTypeBooking      = "booking"
)

// Job represents a background job
type Job struct {
	Type      string
	UserID    string
	Message   string
	Data      map[string]interface{}
	Priority  int
}

// Enqueue adds a job to the Redis queue
func Enqueue(job Job) error {
	return queue.ProduceJob(context.Background(), queue.JobPayload{
		Type:      job.Type,
		UserID:    job.UserID,
		Message:   job.Message,
		Data:      job.Data,
		Priority:  job.Priority,
		CreatedAt: time.Now(),
	})
}

// EnqueueNotification creates and enqueues a notification job
func EnqueueNotification(userID, message string) error {
	return Enqueue(Job{
		Type:    JobTypeNotification,
		UserID:  userID,
		Message: message,
	})
}

// EnqueueEmail creates and enqueues an email job
func EnqueueEmail(userID, subject, body string) error {
	return Enqueue(Job{
		Type: JobTypeEmail,
		UserID: userID,
		Message: subject,
		Data: map[string]interface{}{
			"body": body,
		},
	})
}

// EnqueueBookingNotification creates and enqueues a booking-related notification
func EnqueueBookingNotification(userID, seatID, eventType string) error {
	return Enqueue(Job{
		Type: JobTypeBooking,
		UserID: userID,
		Message: eventType,
		Data: map[string]interface{}{
			"seat_id": seatID,
			"event":   eventType,
		},
	})
}

// EnqueuePaymentNotification creates and enqueues a payment-related notification
func EnqueuePaymentNotification(userID, paymentID, status string) error {
	return Enqueue(Job{
		Type: JobTypePayment,
		UserID: userID,
		Message: status,
		Data: map[string]interface{}{
			"payment_id": paymentID,
			"status":     status,
		},
	})
}

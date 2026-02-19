//Queues
package notifications

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

// NotificationQueue is the main job queue channel
var NotificationQueue = make(chan Job, 100)

// Enqueue adds a job to the queue
func Enqueue(job Job) {
	NotificationQueue <- job
}

// EnqueueNotification creates and enqueues a notification job
func EnqueueNotification(userID, message string) {
	Enqueue(Job{
		Type:    JobTypeNotification,
		UserID:  userID,
		Message: message,
	})
}

// EnqueueEmail creates and enqueues an email job
func EnqueueEmail(userID, subject, body string) {
	Enqueue(Job{
		Type: JobTypeEmail,
		UserID: userID,
		Message: subject,
		Data: map[string]interface{}{
			"body": body,
		},
	})
}

// EnqueueBookingNotification creates and enqueues a booking-related notification
func EnqueueBookingNotification(userID, seatID, eventType string) {
	Enqueue(Job{
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
func EnqueuePaymentNotification(userID, paymentID, status string) {
	Enqueue(Job{
		Type: JobTypePayment,
		UserID: userID,
		Message: status,
		Data: map[string]interface{}{
			"payment_id": paymentID,
			"status":     status,
		},
	})
}

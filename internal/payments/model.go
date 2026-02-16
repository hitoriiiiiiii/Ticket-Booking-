// Models for payments database
package payments

import "time"

type Payment struct {
	ID        string `json:"id"`
	BookingID string `json:"booking_id"`
	UserID    string `json:"user_id"`
	Amount    int    `json:"amount"`

	Status        string `json:"status"` // PENDING, SUCCESS, FAILED
	TransactionID string `json:"transaction_id"`

	CreatedAt time.Time `json:"created_at"`
}

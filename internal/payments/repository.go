// Repository for payments
package payments

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreatePayment(payment *Payment) error {
	payment.ID = uuid.New().String()

	query := `
		INSERT INTO payments
		(id, booking_id, user_id, amount, status, transaction_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.DB.Exec(
		context.Background(),
		query,
		payment.ID,
		payment.BookingID,
		payment.UserID,
		payment.Amount,
		payment.Status,
		payment.TransactionID,
		payment.CreatedAt,
	)

	return err
}

func (r *Repository) GetPaymentByID(id string) (*Payment, error) {

	query := `
		SELECT id, booking_id, user_id, amount, status, transaction_id, created_at
		FROM payments
		WHERE id = $1
	`

	row := r.DB.QueryRow(
		context.Background(),
		query,
		id,
	)

	var payment Payment

	err := row.Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.UserID,
		&payment.Amount,
		&payment.Status,
		&payment.TransactionID,
		&payment.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// Extra UpdateStatus method (you can keep only one)
func (r *Repository) UpdateStatus(paymentID, status, txnID string) error {

	query := `
		UPDATE payments
		SET status=$1, transaction_id=$2
		WHERE id=$3
	`

	_, err := r.DB.Exec(
		context.Background(),
		query,
		status,
		txnID,
		paymentID,
	)

	return err
}

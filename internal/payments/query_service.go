// Query service for payment read operations (CQRS)

package payments

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryService struct {
	DB *pgxpool.Pool
}

func NewQueryService(db *pgxpool.Pool) *QueryService {
	return &QueryService{DB: db}
}

// GetPaymentByID - Query to get a payment by ID
func (s *QueryService) GetPaymentByID(ctx context.Context, id string) (*Payment, error) {
	var p Payment
	err := s.DB.QueryRow(ctx, `
		SELECT id, booking_id, user_id, amount, status, transaction_id, created_at
		FROM payments WHERE id = $1
	`, id).Scan(&p.ID, &p.BookingID, &p.UserID, &p.Amount, &p.Status, &p.TransactionID, &p.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetPaymentByBookingID - Query to get a payment by booking ID
func (s *QueryService) GetPaymentByBookingID(ctx context.Context, bookingID string) (*Payment, error) {
	var p Payment
	err := s.DB.QueryRow(ctx, `
		SELECT id, booking_id, user_id, amount, status, transaction_id, created_at
		FROM payments WHERE booking_id = $1
	`, bookingID).Scan(&p.ID, &p.BookingID, &p.UserID, &p.Amount, &p.Status, &p.TransactionID, &p.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetPaymentsByUserID - Query to get all payments for a user
func (s *QueryService) GetPaymentsByUserID(ctx context.Context, userID string) ([]Payment, error) {
	rows, err := s.DB.Query(ctx, `
		SELECT id, booking_id, user_id, amount, status, transaction_id, created_at
		FROM payments WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		err := rows.Scan(&p.ID, &p.BookingID, &p.UserID, &p.Amount, &p.Status, &p.TransactionID, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}

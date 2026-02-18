// Command service for payment write operations (CQRS)

package payments

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommandService struct {
	Repo *Repository
}

func NewCommandService(repo *Repository) *CommandService {
	return &CommandService{Repo: repo}
}

// InitiatePayment - Command to initiate a new payment
func (s *CommandService) InitiatePayment(ctx context.Context, req InitiatePaymentRequest) (*Payment, error) {
	// Validate input
	if req.BookingID == "" {
		return nil, errors.New("booking ID is required")
	}
	if req.UserID == "" {
		return nil, errors.New("user ID is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	payment := &Payment{
		BookingID: req.BookingID,
		UserID:    req.UserID,
		Amount:    req.Amount,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}

	err := s.Repo.CreatePayment(payment)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// VerifyPayment - Command to verify a payment
func (s *CommandService) VerifyPayment(ctx context.Context, req VerifyPaymentRequest) error {
	// Validate input
	if req.PaymentID == "" {
		return errors.New("payment ID is required")
	}

	payment, err := s.Repo.GetPaymentByID(req.PaymentID)
	if err != nil {
		return err
	}

	if payment.Status != "PENDING" {
		return errors.New("payment already processed")
	}

	var finalStatus string

	switch req.Mode {
	case "success":
		finalStatus = "SUCCESS"
	case "fail":
		finalStatus = "FAILED"
	case "random":
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(2) == 0 {
			finalStatus = "SUCCESS"
		} else {
			finalStatus = "FAILED"
		}
	default:
		return errors.New("invalid mode (use success/fail/random)")
	}

	txnID := fmt.Sprintf("mock_txn_%d", time.Now().Unix())

	return s.Repo.UpdateStatus(req.PaymentID, finalStatus, txnID)
}

// RefundPayment - Command to refund a payment
func (s *CommandService) RefundPayment(ctx context.Context, paymentID string) error {
	// Validate input
	if paymentID == "" {
		return errors.New("payment ID is required")
	}

	payment, err := s.Repo.GetPaymentByID(paymentID)
	if err != nil {
		return err
	}

	if payment.Status != "SUCCESS" {
		return errors.New("can only refund successful payments")
	}

	return s.Repo.UpdateStatus(paymentID, "REFUNDED", "")
}

// Initialize CommandService with database pool (for compatibility)
func NewCommandServiceWithDB(db *pgxpool.Pool) *CommandService {
	repo := NewRepository(db)
	return &CommandService{Repo: repo}
}

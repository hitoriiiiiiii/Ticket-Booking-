//Internal payment systems service layer

package payments

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

//Initiate payment process
//(Always creates PENDING payment)

func (s *Service) InitiatePayment(bookingID string, userID string, amount int) (*Payment, error) {
	payment := &Payment{
		BookingID: bookingID,
		UserID:    userID,
		Amount:    amount,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}

	err := s.Repo.CreatePayment(payment)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// Verify payment with external gateway (simulated here)
func (s *Service) Verify(paymentID string, mode string) error {

	payment, err := s.Repo.GetPaymentByID(paymentID)
	if err != nil {
		return err
	}

	if payment.Status != "PENDING" {
		return errors.New("payment already processed")
	}

	var finalStatus string

	switch mode {

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

	return s.Repo.UpdateStatus(paymentID, finalStatus, txnID)
}

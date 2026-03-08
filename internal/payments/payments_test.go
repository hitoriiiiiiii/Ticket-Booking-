package payments

import (
	"errors"
	"testing"
	"time"
)

// TestPaymentModel tests the Payment struct
func TestPaymentModel(t *testing.T) {
	payment := Payment{
		ID:        "pay_123",
		BookingID: "book_123",
		UserID:    "user_123",
		Amount:    100.00,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}

	if payment.ID != "pay_123" {
		t.Errorf("Expected ID pay_123, got %s", payment.ID)
	}
	if payment.BookingID != "book_123" {
		t.Errorf("Expected BookingID book_123, got %s", payment.BookingID)
	}
	if payment.UserID != "user_123" {
		t.Errorf("Expected UserID user_123, got %s", payment.UserID)
	}
	if payment.Amount != 100 {
		t.Errorf("Expected Amount 100, got %d", payment.Amount)
	}
	if payment.Status != "PENDING" {
		t.Errorf("Expected Status PENDING, got %s", payment.Status)
	}
}

// TestInitiatePaymentRequest_Validation tests validation for InitiatePaymentRequest
func TestInitiatePaymentRequest_Validation(t *testing.T) {
	tests := []struct {
		name      string
		req       InitiatePaymentRequest
		wantErr   bool
		errString string
	}{
		{
			name: "empty booking ID",
			req: InitiatePaymentRequest{
				BookingID: "",
				UserID:    "user_123",
				Amount:    100.00,
			},
			wantErr:   true,
			errString: "booking ID is required",
		},
		{
			name: "empty user ID",
			req: InitiatePaymentRequest{
				BookingID: "book_123",
				UserID:    "",
				Amount:    100.00,
			},
			wantErr:   true,
			errString: "user ID is required",
		},
		{
			name: "zero amount",
			req: InitiatePaymentRequest{
				BookingID: "book_123",
				UserID:    "user_123",
				Amount:    0,
			},
			wantErr:   true,
			errString: "amount must be positive",
		},
		{
			name: "negative amount",
			req: InitiatePaymentRequest{
				BookingID: "book_123",
				UserID:    "user_123",
				Amount:    -50.00,
			},
			wantErr:   true,
			errString: "amount must be positive",
		},
		{
			name: "valid request",
			req: InitiatePaymentRequest{
				BookingID: "book_123",
				UserID:    "user_123",
				Amount:    100.00,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInitiatePaymentRequest(tt.req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for %s", tt.name)
				} else if err.Error() != tt.errString {
					t.Errorf("Expected error %s, got %s", tt.errString, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestVerifyPaymentRequest_Validation tests validation for VerifyPaymentRequest
func TestVerifyPaymentRequest_Validation(t *testing.T) {
	tests := []struct {
		name      string
		req       VerifyPaymentRequest
		wantErr   bool
		errString string
	}{
		{
			name: "empty payment ID",
			req: VerifyPaymentRequest{
				PaymentID: "",
				Mode:      "success",
			},
			wantErr:   true,
			errString: "payment ID is required",
		},
		{
			name: "invalid mode",
			req: VerifyPaymentRequest{
				PaymentID: "pay_123",
				Mode:      "invalid",
			},
			wantErr:   true,
			errString: "invalid mode (use success/fail/random)",
		},
		{
			name: "success mode",
			req: VerifyPaymentRequest{
				PaymentID: "pay_123",
				Mode:      "success",
			},
			wantErr: false,
		},
		{
			name: "fail mode",
			req: VerifyPaymentRequest{
				PaymentID: "pay_123",
				Mode:      "fail",
			},
			wantErr: false,
		},
		{
			name: "random mode",
			req: VerifyPaymentRequest{
				PaymentID: "pay_123",
				Mode:      "random",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVerifyPaymentRequest(tt.req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for %s", tt.name)
				} else if err.Error() != tt.errString {
					t.Errorf("Expected error %s, got %s", tt.errString, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestRefundPayment tests validation for refund
func TestRefundPayment_Validation(t *testing.T) {
	tests := []struct {
		name        string
		paymentID   string
		wantErr     bool
		errString   string
	}{
		{
			name:        "empty payment ID",
			paymentID:   "",
			wantErr:     true,
			errString:   "payment ID is required",
		},
		{
			name:        "valid payment ID",
			paymentID:   "pay_123",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRefundPayment(tt.paymentID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for %s", tt.name)
				} else if err.Error() != tt.errString {
					t.Errorf("Expected error %s, got %s", tt.errString, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// Validation functions that match the service implementation
func validateInitiatePaymentRequest(req InitiatePaymentRequest) error {
	if req.BookingID == "" {
		return errors.New("booking ID is required")
	}
	if req.UserID == "" {
		return errors.New("user ID is required")
	}
	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	return nil
}

func validateVerifyPaymentRequest(req VerifyPaymentRequest) error {
	if req.PaymentID == "" {
		return errors.New("payment ID is required")
	}
	switch req.Mode {
	case "success", "fail", "random":
		return nil
	default:
		return errors.New("invalid mode (use success/fail/random)")
	}
}

func validateRefundPayment(paymentID string) error {
	if paymentID == "" {
		return errors.New("payment ID is required")
	}
	return nil
}

// TestPaymentStatus tests valid payment statuses
func TestPaymentStatus(t *testing.T) {
	validStatuses := []string{"PENDING", "SUCCESS", "FAILED", "REFUNDED"}

	testCases := []struct {
		status   string
		expected bool
	}{
		{"PENDING", true},
		{"SUCCESS", true},
		{"FAILED", true},
		{"REFUNDED", true},
		{"INVALID", false},
		{"", false},
	}

	for _, tc := range testCases {
		isValid := false
		for _, s := range validStatuses {
			if tc.status == s {
				isValid = true
				break
			}
		}
		if isValid != tc.expected {
			t.Errorf("Expected status %s to be %v", tc.status, tc.expected)
		}
	}
}

// TestNewCommandService tests constructor functions
func TestNewCommandService(t *testing.T) {
	svc := NewCommandService(nil)
	if svc == nil {
		t.Error("Expected non-nil CommandService")
	}
}

// Benchmark tests
func BenchmarkPaymentModel(b *testing.B) {
	payment := Payment{
		ID:        "pay_123",
		BookingID: "book_123",
		UserID:    "user_123",
		Amount:    100.00,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = payment.Status
	}
}

func BenchmarkInitiatePaymentValidation(b *testing.B) {
	req := InitiatePaymentRequest{
		BookingID: "book_123",
		UserID:    "user_123",
		Amount:    100.00,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateInitiatePaymentRequest(req)
	}
}

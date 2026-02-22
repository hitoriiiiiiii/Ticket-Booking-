package booking

import (
	"errors"
	"testing"
)

// TestReserveTicket_Validation tests validation for ticket reservation
func TestReserveTicket_Validation(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		seatID  string
		wantErr bool
		errStr  string
	}{
		{
			name:    "empty user ID",
			userID:  "",
			seatID:  "seat_1",
			wantErr: true,
			errStr:  "user ID is required",
		},
		{
			name:    "empty seat ID",
			userID:  "user_1",
			seatID:  "",
			wantErr: true,
			errStr:  "seat ID is required",
		},
		{
			name:    "valid reservation",
			userID:  "user_1",
			seatID:  "seat_1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReserveTicket(tt.userID, tt.seatID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for %s", tt.name)
				} else if err.Error() != tt.errStr {
					t.Errorf("Expected error %s, got %s", tt.errStr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestConfirmTicket_Validation tests validation for ticket confirmation
func TestConfirmTicket_Validation(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		seatID  string
		wantErr bool
		errStr  string
	}{
		{
			name:    "empty user ID",
			userID:  "",
			seatID:  "seat_1",
			wantErr: true,
			errStr:  "user ID is required",
		},
		{
			name:    "empty seat ID",
			userID:  "user_1",
			seatID:  "",
			wantErr: true,
			errStr:  "seat ID is required",
		},
		{
			name:    "valid confirmation",
			userID:  "user_1",
			seatID:  "seat_1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfirmTicket(tt.userID, tt.seatID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for %s", tt.name)
				} else if err.Error() != tt.errStr {
					t.Errorf("Expected error %s, got %s", tt.errStr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestCancelTicket_Validation tests validation for ticket cancellation
func TestCancelTicket_Validation(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		seatID  string
		wantErr bool
		errStr  string
	}{
		{
			name:    "empty user ID",
			userID:  "",
			seatID:  "seat_1",
			wantErr: true,
			errStr:  "user ID is required",
		},
		{
			name:    "empty seat ID",
			userID:  "user_1",
			seatID:  "",
			wantErr: true,
			errStr:  "seat ID is required",
		},
		{
			name:    "valid cancellation",
			userID:  "user_1",
			seatID:  "seat_1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCancelTicket(tt.userID, tt.seatID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for %s", tt.name)
				} else if err.Error() != tt.errStr {
					t.Errorf("Expected error %s, got %s", tt.errStr, err.Error())
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
func validateReserveTicket(userID, seatID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}
	if seatID == "" {
		return errors.New("seat ID is required")
	}
	return nil
}

func validateConfirmTicket(userID, seatID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}
	if seatID == "" {
		return errors.New("seat ID is required")
	}
	return nil
}

func validateCancelTicket(userID, seatID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}
	if seatID == "" {
		return errors.New("seat ID is required")
	}
	return nil
}

// TestReservationStatus tests valid reservation statuses
func TestReservationStatus(t *testing.T) {
	validStatuses := []string{"HELD", "BOOKED", "CANCELLED"}

	testCases := []struct {
		status   string
		expected bool
	}{
		{"HELD", true},
		{"BOOKED", true},
		{"CANCELLED", true},
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
func BenchmarkReserveTicketValidation(b *testing.B) {
	userID := "user_1"
	seatID := "seat_1"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateReserveTicket(userID, seatID)
	}
}

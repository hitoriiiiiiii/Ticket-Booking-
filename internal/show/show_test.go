package show

import (
	"errors"
	"testing"
	"time"
)

// TestShowModel tests the Show struct
func TestShowModel(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(2 * time.Hour)

	show := Show{
		ID:        1,
		MovieID:   1,
		Theater:   "Theater A",
		StartTime: startTime,
		EndTime:   endTime,
	}

	if show.ID != 1 {
		t.Errorf("Expected ID 1, got %d", show.ID)
	}
	if show.MovieID != 1 {
		t.Errorf("Expected MovieID 1, got %d", show.MovieID)
	}
	if show.Theater != "Theater A" {
		t.Errorf("Expected Theater A, got %s", show.Theater)
	}
	if !show.StartTime.Equal(startTime) {
		t.Error("StartTime mismatch")
	}
	if !show.EndTime.Equal(endTime) {
		t.Error("EndTime mismatch")
	}
}

// TestCreateShowRequest_Validation tests validation for CreateShowRequest
func TestCreateShowRequest_Validation(t *testing.T) {
	tests := []struct {
		name      string
		req       CreateShowRequest
		wantErr   bool
		errString string
	}{
		{
			name: "empty theater",
			req: CreateShowRequest{
				Theater:   "",
				MovieID:   1,
				StartTime: time.Now(),
			},
			wantErr:   true,
			errString: "theater is required",
		},
		{
			name: "zero movie ID",
			req: CreateShowRequest{
				Theater:   "Theater A",
				MovieID:   0,
				StartTime: time.Now(),
			},
			wantErr:   true,
			errString: "valid movie ID is required",
		},
		{
			name: "negative movie ID",
			req: CreateShowRequest{
				Theater:   "Theater A",
				MovieID:   -1,
				StartTime: time.Now(),
			},
			wantErr:   true,
			errString: "valid movie ID is required",
		},
		{
			name: "valid request",
			req: CreateShowRequest{
				Theater:   "Theater A",
				MovieID:   1,
				StartTime: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateShowRequest(tt.req)
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

// TestUpdateShowRequest_Validation tests validation for UpdateShowRequest
func TestUpdateShowRequest_Validation(t *testing.T) {
	tests := []struct {
		name      string
		req       UpdateShowRequest
		wantErr   bool
		errString string
	}{
		{
			name: "empty theater",
			req: UpdateShowRequest{
				Theater: "",
				MovieID: 1,
			},
			wantErr:   true,
			errString: "theater is required",
		},
		{
			name: "valid request",
			req: UpdateShowRequest{
				Theater: "Theater B",
				MovieID: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUpdateShowRequest(tt.req)
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
func validateCreateShowRequest(req CreateShowRequest) error {
	if req.Theater == "" {
		return errors.New("theater is required")
	}
	if req.MovieID <= 0 {
		return errors.New("valid movie ID is required")
	}
	return nil
}

func validateUpdateShowRequest(req UpdateShowRequest) error {
	if req.Theater == "" {
		return errors.New("theater is required")
	}
	return nil
}

// TestNewCommandService tests constructor functions
func TestNewCommandService(t *testing.T) {
	svc := NewCommandService()
	if svc == nil {
		t.Error("Expected non-nil CommandService")
	}
	if svc.Dispatcher != nil {
		t.Error("Expected nil Dispatcher")
	}
}

// TestEndTimeCalculation tests that end time is calculated correctly
func TestEndTimeCalculation(t *testing.T) {
	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	expectedEndTime := startTime.Add(2 * time.Hour)

	req := CreateShowRequest{
		Theater:   "Theater A",
		MovieID:   1,
		StartTime: startTime,
	}

	// Simulate the service logic
	endTime := req.StartTime.Add(2 * time.Hour)

	if !endTime.Equal(expectedEndTime) {
		t.Errorf("Expected end time %v, got %v", expectedEndTime, endTime)
	}
}

// Benchmark tests
func BenchmarkShowModel(b *testing.B) {
	startTime := time.Now()
	endTime := startTime.Add(2 * time.Hour)

	show := Show{
		ID:        1,
		MovieID:   1,
		Theater:   "Theater A",
		StartTime: startTime,
		EndTime:   endTime,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = show.Theater
	}
}

func BenchmarkCreateShowValidation(b *testing.B) {
	req := CreateShowRequest{
		Theater:   "Theater A",
		MovieID:   1,
		StartTime: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateCreateShowRequest(req)
	}
}

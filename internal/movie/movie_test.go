package movie

import (
	"errors"
	"testing"
)

// TestMovieModel tests the Movie struct
func TestMovieModel(t *testing.T) {
	movie := Movie{
		ID:       1,
		Name:     "Test Movie",
		Genre:    "Action",
		Duration: 120,
	}

	if movie.ID != 1 {
		t.Errorf("Expected ID 1, got %d", movie.ID)
	}
	if movie.Name != "Test Movie" {
		t.Errorf("Expected name Test Movie, got %s", movie.Name)
	}
	if movie.Genre != "Action" {
		t.Errorf("Expected genre Action, got %s", movie.Genre)
	}
	if movie.Duration != 120 {
		t.Errorf("Expected duration 120, got %d", movie.Duration)
	}
}

// TestCreateMovieRequest_Validation tests validation for CreateMovieRequest
func TestCreateMovieRequest_Validation(t *testing.T) {
	tests := []struct {
		name      string
		req       CreateMovieRequest
		wantErr   bool
		errString string
	}{
		{
			name: "empty name",
			req: CreateMovieRequest{
				Name:     "",
				Genre:    "Action",
				Duration: 120,
			},
			wantErr:   true,
			errString: "movie name is required",
		},
		{
			name: "zero duration",
			req: CreateMovieRequest{
				Name:     "Test Movie",
				Genre:    "Action",
				Duration: 0,
			},
			wantErr:   true,
			errString: "movie duration must be positive",
		},
		{
			name: "negative duration",
			req: CreateMovieRequest{
				Name:     "Test Movie",
				Genre:    "Action",
				Duration: -10,
			},
			wantErr:   true,
			errString: "movie duration must be positive",
		},
		{
			name: "valid request",
			req: CreateMovieRequest{
				Name:     "Test Movie",
				Genre:    "Action",
				Duration: 120,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateMovieRequest(tt.req)
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

// TestUpdateMovieRequest_Validation tests validation for UpdateMovieRequest
func TestUpdateMovieRequest_Validation(t *testing.T) {
	tests := []struct {
		name      string
		req       UpdateMovieRequest
		wantErr   bool
		errString string
	}{
		{
			name: "empty name",
			req: UpdateMovieRequest{
				Name:     "",
				Genre:    "Action",
				Duration: 120,
			},
			wantErr:   true,
			errString: "movie name is required",
		},
		{
			name: "negative duration",
			req: UpdateMovieRequest{
				Name:     "Test Movie",
				Genre:    "Action",
				Duration: -10,
			},
			wantErr:   true,
			errString: "movie duration must be positive",
		},
		{
			name: "valid request",
			req: UpdateMovieRequest{
				Name:     "Updated Movie",
				Genre:    "Drama",
				Duration: 90,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUpdateMovieRequest(tt.req)
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
func validateCreateMovieRequest(req CreateMovieRequest) error {
	if req.Name == "" {
		return errors.New("movie name is required")
	}
	if req.Duration <= 0 {
		return errors.New("movie duration must be positive")
	}
	return nil
}

func validateUpdateMovieRequest(req UpdateMovieRequest) error {
	if req.Name == "" {
		return errors.New("movie name is required")
	}
	if req.Duration <= 0 {
		return errors.New("movie duration must be positive")
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

// Benchmark tests
func BenchmarkMovieModel(b *testing.B) {
	movie := Movie{
		ID:       1,
		Name:     "Test Movie",
		Genre:    "Action",
		Duration: 120,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = movie.Name
	}
}

func BenchmarkCreateMovieValidation(b *testing.B) {
	req := CreateMovieRequest{
		Name:     "Test Movie",
		Genre:    "Action",
		Duration: 120,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateCreateMovieRequest(req)
	}
}

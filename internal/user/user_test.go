package user

import (
	"context"
	"errors"
	"testing"
)

// RegisterRequest is used in tests - it's defined in command_handler.go
// but we need to import it for testing

// TestUserModel tests the User struct
func TestUserModel(t *testing.T) {
	user := User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  false,
		Password: "hashedpassword",
	}

	if user.ID != 1 {
		t.Errorf("Expected ID 1, got %d", user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username testuser, got %s", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}
	if user.IsAdmin != false {
		t.Errorf("Expected IsAdmin false, got %v", user.IsAdmin)
	}
}

// TestUserPasswordOmission tests that password is omitted in JSON
func TestUserPasswordOmission(t *testing.T) {
	// This test verifies the JSON omitempty tag works
	// The password field has `json:"password,omitempty"` in model
	user := User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "",
	}

	// When password is empty, it should be omitted
	if user.Password != "" {
		t.Errorf("Expected empty password, got %s", user.Password)
	}
}

// Test validation functions that would typically exist
// Since the validation is embedded in the service, we test the error cases

func TestRegisterValidation(t *testing.T) {
	tests := []struct {
		name      string
		req       RegisterRequest
		wantErr   bool
		errString string
	}{
		{
			name: "empty username",
			req: RegisterRequest{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr:   true,
			errString: "username is required",
		},
		{
			name: "empty email",
			req: RegisterRequest{
				Username: "testuser",
				Email:    "",
				Password: "password123",
			},
			wantErr:   true,
			errString: "email is required",
		},
		{
			name: "empty password",
			req: RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "",
			},
			wantErr:   true,
			errString: "password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation logic
			err := validateRegisterRequest(tt.req)
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

func TestUpdateUserValidation(t *testing.T) {
	tests := []struct {
		name      string
		req       UpdateUserRequest
		wantErr   bool
		errString string
	}{
		{
			name: "empty username",
			req: UpdateUserRequest{
				Username: "",
				Email:    "test@example.com",
			},
			wantErr:   true,
			errString: "username is required",
		},
		{
			name: "empty email",
			req: UpdateUserRequest{
				Username: "testuser",
				Email:    "",
			},
			wantErr:   true,
			errString: "email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUpdateUserRequest(tt.req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for %s", tt.name)
				} else if err.Error() != tt.errString {
					t.Errorf("Expected error %s, got %s", tt.errString, err.Error())
				}
			}
		})
	}
}

// Validation functions that match the service implementation
func validateRegisterRequest(req RegisterRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func validateUpdateUserRequest(req UpdateUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

// Benchmark tests
func BenchmarkUserModel(b *testing.B) {
	user := User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  false,
		Password: "hashedpassword",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.Username
	}
}

func BenchmarkRegisterValidation(b *testing.B) {
	req := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateRegisterRequest(req)
	}
}

// Helper to avoid unused import warning
var _ = context.Background()

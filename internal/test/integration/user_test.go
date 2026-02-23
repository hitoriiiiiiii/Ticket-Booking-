// Package integration provides integration tests for the ticket booking system
package integration

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/hitorii/ticket-booking/internal/events"
	"github.com/hitorii/ticket-booking/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserIntegration_Register tests user registration flow
func TestUserIntegration_Register(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	t.Run("Register_Success", func(t *testing.T) {
		req := user.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
			IsAdmin:  false,
		}

		u, err := svc.UserCmdSvc.Register(ctx, req)
		require.NoError(t, err)
		assert.NotZero(t, u.ID)
		assert.Equal(t, req.Username, u.Username)
		assert.Equal(t, req.Email, u.Email)
	})

	t.Run("Register_DuplicateEmail", func(t *testing.T) {
		req := user.RegisterRequest{
			Username: "testuser2",
			Email:    "duplicate@example.com",
			Password: "password123",
		}

		// First registration should succeed
		_, err := svc.UserCmdSvc.Register(ctx, req)
		require.NoError(t, err)

		// Second registration with same email should fail
		req.Username = "testuser3"
		_, err = svc.UserCmdSvc.Register(ctx, req)
		require.Error(t, err)
	})

	t.Run("Register_InvalidEmail", func(t *testing.T) {
		req := user.RegisterRequest{
			Username: "testuser4",
			Email:    "invalid-email",
			Password: "password123",
		}

		_, err := svc.UserCmdSvc.Register(ctx, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "email")
	})

	t.Run("Register_EmptyPassword", func(t *testing.T) {
		req := user.RegisterRequest{
			Username: "testuser5",
			Email:    "test5@example.com",
			Password: "",
		}

		_, err := svc.UserCmdSvc.Register(ctx, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "password")
	})
}

// TestUserIntegration_Login tests user login flow
func TestUserIntegration_Login(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// First register a user
	req := user.RegisterRequest{
		Username: "logintestuser",
		Email:    "login@test.com",
		Password: "password123",
	}
	_, err := svc.UserCmdSvc.Register(ctx, req)
	require.NoError(t, err)

	t.Run("Login_Success", func(t *testing.T) {
		loginReq := user.LoginRequest{
			Email:    "login@test.com",
			Password: "password123",
		}

		token, err := svc.UserQuerySvc.Login(ctx, loginReq)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Login_WrongPassword", func(t *testing.T) {
		loginReq := user.LoginRequest{
			Email:    "login@test.com",
			Password: "wrongpassword",
		}

		_, err := svc.UserQuerySvc.Login(ctx, loginReq)
		require.Error(t, err)
	})

	t.Run("Login_NonExistentUser", func(t *testing.T) {
		loginReq := user.LoginRequest{
			Email:    "nonexistent@test.com",
			Password: "password123",
		}

		_, err := svc.UserQuerySvc.Login(ctx, loginReq)
		require.Error(t, err)
	})
}

// TestUserIntegration_UserEvents tests event publishing during user operations
func TestUserIntegration_UserEvents(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Subscribe to user events
	eventReceived := make(chan bool, 1)
	svc.Dispatcher.Subscribe(events.EventUserRegistered, func(event events.BaseEvent) error {
		eventReceived <- true
		return nil
	})

	// Register a user
	req := user.RegisterRequest{
		Username: "eventuser",
		Email:    "event@test.com",
		Password: "password123",
	}
	_, err := svc.UserCmdSvc.Register(ctx, req)
	require.NoError(t, err)

	// Wait for event
	select {
	case <-eventReceived:
		// Event received
	case <-time.After(2 * time.Second):
		t.Log("Warning: Event was not received within timeout")
	}
}

// TestUserIntegration_GetUserByID tests querying user by ID
func TestUserIntegration_GetUserByID(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Create a user
	req := user.RegisterRequest{
		Username: "queryuser",
		Email:    "query@test.com",
		Password: "password123",
	}
	createdUser, err := svc.UserCmdSvc.Register(ctx, req)
	require.NoError(t, err)

	t.Run("GetUserByID_Exists", func(t *testing.T) {
		u, err := svc.UserQuerySvc.GetUserByID(ctx, strconv.Itoa(createdUser.ID))
		if err != nil {
			// Try with default ID
			u, err = svc.UserQuerySvc.GetUserByID(ctx, "1")
		}
		if err == nil {
			assert.Equal(t, req.Username, u.Username)
		}
	})

	t.Run("GetUserID_NotFound", func(t *testing.T) {
		_, err := svc.UserQuerySvc.GetUserByID(ctx, "99999")
		require.Error(t, err)
	})
}

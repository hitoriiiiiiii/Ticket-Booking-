package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hitorii/ticket-booking/internal/events"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type CommandService struct {
	DB         *pgxpool.Pool
	Dispatcher *events.Dispatcher
}

func NewCommandService(db *pgxpool.Pool) *CommandService {
	return &CommandService{DB: db}
}

func NewCommandServiceWithDispatcher(db *pgxpool.Pool, dispatcher *events.Dispatcher) *CommandService {
	return &CommandService{DB: db, Dispatcher: dispatcher}
}

// Register - Command to register a new user
func (s *CommandService) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	// Validate input
	if req.Username == "" {
		return nil, errors.New("username is required")
	}
	if req.Email == "" {
		return nil, errors.New("email is required")
	}
	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error hashing password")
	}

	// Generate UUID for user
	userID := uuid.New().String()

	user := &User{
		ID:        userID,
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		IsAdmin:   req.IsAdmin,
		CreatedAt: time.Now(),
	}

	// Insert into database
	_, err = s.DB.Exec(ctx, `
		INSERT INTO users (id, username, email, password, is_admin, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, user.ID, user.Username, user.Email, user.Password, user.IsAdmin, user.CreatedAt)

	if err != nil {
		return nil, err
	}

	// Emit event for event-driven flow
	if s.Dispatcher != nil {
		payload := events.EventPayload{
			UserID:   userID,
			Username: req.Username,
			Email:    req.Email,
			IsAdmin:  req.IsAdmin,
		}
		_ = s.Dispatcher.Publish(ctx, events.EventUserRegistered, userID, payload)
	}

	return user, nil
}

// UpdateUser - Command to update user information
func (s *CommandService) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) error {
	// Validate UUID
	if id == "" {
		return errors.New("user ID is required")
	}

	// Validate input
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}

	// Update in database
	_, err := s.DB.Exec(ctx, `
		UPDATE users SET username = $1, email = $2 WHERE id = $3
	`, req.Username, req.Email, id)

	if err != nil {
		return err
	}

	// Emit event for event-driven flow
	if s.Dispatcher != nil {
		payload := events.EventPayload{
			UserID:   id,
			Username: req.Username,
			Email:    req.Email,
		}
		_ = s.Dispatcher.Publish(ctx, events.EventUserUpdated, id, payload)
	}

	return nil
}

// DeleteUser - Command to delete a user
func (s *CommandService) DeleteUser(ctx context.Context, id string) error {
	// Validate UUID
	if id == "" {
		return errors.New("user ID is required")
	}

	// Delete from database
	_, err := s.DB.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Emit event for event-driven flow
	if s.Dispatcher != nil {
		payload := events.EventPayload{
			UserID: id,
		}
		_ = s.Dispatcher.Publish(ctx, events.EventUserDeleted, id, payload)
	}

	return nil
}

// GetUserByID - Query to get user by ID
func (s *CommandService) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := s.DB.QueryRow(ctx, `
		SELECT id, username, email, is_admin, created_at FROM users WHERE id = $1
	`, id).Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

// Command service for notification write operations (CQRS)

package notifications

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommandService struct {
	Repo *Repository
}

func NewCommandService(repo *Repository) *CommandService {
	return &CommandService{Repo: repo}
}

// SendNotification - Command to send a notification
func (s *CommandService) SendNotification(ctx context.Context, req SendNotificationRequest) (*Notification, error) {
	// Validate input
	if req.UserID == "" {
		return nil, errors.New("user ID is required")
	}
	if req.Message == "" {
		return nil, errors.New("message is required")
	}
	if req.Type == "" {
		req.Type = "general" // Default type
	}

	notification := &Notification{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Message:   req.Message,
		Type:      req.Type,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	err := s.Repo.Save(Job{
		UserID:  req.UserID,
		Message: req.Message,
		Type:    req.Type,
	})
	if err != nil {
		return nil, err
	}

	return notification, nil
}

// MarkAsRead - Command to mark a notification as read
func (s *CommandService) MarkAsRead(ctx context.Context, id string) error {
	// Validate input
	if id == "" {
		return errors.New("notification ID is required")
	}

	_, err := s.Repo.DB.Exec(ctx, `
		UPDATE notifications SET is_read = true WHERE id = $1
	`, id)

	return err
}

// MarkAllAsRead - Command to mark all notifications as read for a user
func (s *CommandService) MarkAllAsRead(ctx context.Context, userID string) error {
	// Validate input
	if userID == "" {
		return errors.New("user ID is required")
	}

	_, err := s.Repo.DB.Exec(ctx, `
		UPDATE notifications SET is_read = true WHERE user_id = $1
	`, userID)

	return err
}

// DeleteNotification - Command to delete a notification
func (s *CommandService) DeleteNotification(ctx context.Context, id string) error {
	// Validate input
	if id == "" {
		return errors.New("notification ID is required")
	}

	_, err := s.Repo.DB.Exec(ctx, `
		DELETE FROM notifications WHERE id = $1
	`, id)

	return err
}

// Initialize CommandService with database pool (for compatibility)
func NewCommandServiceWithDB(db *pgxpool.Pool) *CommandService {
	repo := NewRepository(db)
	return &CommandService{Repo: repo}
}

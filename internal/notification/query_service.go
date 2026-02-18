// Query service for notification read operations (CQRS)

package notifications

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryService struct {
	DB *pgxpool.Pool
}

func NewQueryService(db *pgxpool.Pool) *QueryService {
	return &QueryService{DB: db}
}

// GetUserNotifications - Query to get all notifications for a user
func (s *QueryService) GetUserNotifications(ctx context.Context, userID string) ([]Notification, error) {
	query := `
		SELECT id, user_id, message, type, is_read, created_at
		FROM notifications
		WHERE user_id=$1
		ORDER BY created_at DESC
	`

	rows, err := s.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Notification

	for rows.Next() {
		var n Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Message,
			&n.Type,
			&n.IsRead,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, n)
	}

	return list, nil
}

// GetUnreadNotifications - Query to get unread notifications for a user
func (s *QueryService) GetUnreadNotifications(ctx context.Context, userID string) ([]Notification, error) {
	query := `
		SELECT id, user_id, message, type, is_read, created_at
		FROM notifications
		WHERE user_id=$1 AND is_read=false
		ORDER BY created_at DESC
	`

	rows, err := s.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Notification

	for rows.Next() {
		var n Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Message,
			&n.Type,
			&n.IsRead,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, n)
	}

	return list, nil
}

// GetNotificationByID - Query to get a notification (s *Query by ID
funcService) GetNotificationByID(ctx context.Context, id string) (*Notification, error) {
	var n Notification
	err := s.DB.QueryRow(ctx, `
		SELECT id, user_id, message, type, is_read, created_at
		FROM notifications WHERE id = $1
	`, id).Scan(&n.ID, &n.UserID, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &n, nil
}

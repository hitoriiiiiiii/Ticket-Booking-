package notification

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Save(job Job) error {

	query := `
	INSERT INTO notifications (id, user_id, message, type, is_read, created_at)
	VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := r.DB.Exec(
		context.Background(),
		query,
		uuid.New().String(),
		job.UserID,
		job.Message,
		job.Type,
		false,
		time.Now(),
	)

	return err
}

func (r *Repository) GetUserNotifications(userID string) ([]Notification, error) {

	query := `
	SELECT id, user_id, message, type, is_read, created_at
	FROM notifications
	WHERE user_id=$1
	ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Notification

	for rows.Next() {
		var n Notification
		rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Message,
			&n.Type,
			&n.IsRead,
			&n.CreatedAt,
		)
		list = append(list, n)
	}

	return list, nil
}

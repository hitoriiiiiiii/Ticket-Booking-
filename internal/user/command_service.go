// Command service for user write operations (CQRS)

package user

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type CommandService struct {
	DB *pgxpool.Pool
}

func NewCommandService(db *pgxpool.Pool) *CommandService {
	return &CommandService{DB: db}
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

	user := &User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		IsAdmin:   req.IsAdmin,
		CreatedAt: time.Now(),
	}

	// Insert into database
	_, err = s.DB.Exec(ctx, `
		INSERT INTO users (username, email, password, is_admin, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, user.Username, user.Email, user.Password, user.IsAdmin, user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser - Command to update user information
func (s *CommandService) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) error {
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

	return err
}

// DeleteUser - Command to delete a user
func (s *CommandService) DeleteUser(ctx context.Context, id string) error {
	// Delete from database
	_, err := s.DB.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

// Query service for user read operations (CQRS)

package user

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("supersecretkey")

type QueryService struct {
	DB *pgxpool.Pool
}

func NewQueryService(db *pgxpool.Pool) *QueryService {
	return &QueryService{DB: db}
}

// Login - Query to verify user credentials and generate JWT token
func (s *QueryService) Login(ctx context.Context, req LoginRequest) (string, error) {
	var user User
	err := s.DB.QueryRow(ctx, `
		SELECT id, username, email, password, is_admin, created_at
		FROM users WHERE email = $1
	`, req.Email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.CreatedAt)

	if err != nil {
		return "", err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", err
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     86400, // 24h expiry
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ListUsers - Query to get all users
func (s *QueryService) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := s.DB.Query(ctx, `
		SELECT id, username, email, is_admin, created_at
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserByID - Query to get a single user by ID
func (s *QueryService) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := s.DB.QueryRow(ctx, `
		SELECT id, username, email, is_admin, created_at
		FROM users WHERE id = $1
	`, id).Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

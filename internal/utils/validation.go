// Validation utilities for the ticket booking system
package utils

import (
	"context"
	"regexp"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UUIDRegex is the regex pattern for validating UUIDs
var UUIDRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(s string) bool {
	if s == "" {
		return false
	}
	// First try standard UUID parsing
	_, err := uuid.Parse(s)
	if err == nil {
		return true
	}
	// Fall back to regex check for lowercase hex format
	return UUIDRegex.MatchString(s)
}

// ValidateUUID returns an error if the string is not a valid UUID
func ValidateUUID(fieldName, value string) error {
	if !IsValidUUID(value) {
		return ErrInvalidUUID(fieldName)
	}
	return nil
}

// ErrInvalidUUID represents an invalid UUID error
type ErrInvalidUUID string

func (e ErrInvalidUUID) Error() string {
	return "invalid UUID format for " + string(e)
}

// UserExists checks if a user exists in the database
func UserExists(ctx context.Context, db *pgxpool.Pool, userID string) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SeatExists checks if a seat exists in the database
func SeatExists(ctx context.Context, db *pgxpool.Pool, seatID string) (bool, error) {
	var exists bool
	// Check in shows table for seat_id column
	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM shows WHERE id = $1)", seatID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// ShowExists checks if a show exists in the database
func ShowExists(ctx context.Context, db *pgxpool.Pool, showID string) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM shows WHERE id = $1)", showID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// ValidateUserExists returns an error if the user does not exist
func ValidateUserExists(ctx context.Context, db *pgxpool.Pool, userID string) error {
	exists, err := UserExists(ctx, db, userID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}
	return nil
}

// ErrUserNotFound represents a user not found error
var ErrUserNotFound = &UserNotFoundError{}

type UserNotFoundError struct{}

func (e *UserNotFoundError) Error() string {
	return "user not found"
}


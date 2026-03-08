// Models for user database

package user

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	Password  string    `json:"password,omitempty"` // omit password in JSON responses
	CreatedAt time.Time `json:"created_at"`
}

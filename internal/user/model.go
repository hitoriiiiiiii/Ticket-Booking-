//models for the user
package user

import "time"

type User struct {
	ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	Password  string    `json:"password,omitempty"` // omit password in JSON responses
	CreatedAt time.Time `json:"created_at"`
}

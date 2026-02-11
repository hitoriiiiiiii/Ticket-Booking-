//middleware for the authentication of the user and admin 

package middleware

import (
    "net/http"
    "strings"
    "ticket-booking/internal/user"
    "gorm.io/gorm"
    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("supersecretkey")

// AuthMiddleware verifies JWT token and adds user to context
func AuthMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Missing token", http.StatusUnauthorized)
                return
            }

            tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
            token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
                return jwtKey, nil
            })
            if err != nil || !token.Valid {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                http.Error(w, "Invalid token claims", http.StatusUnauthorized)
                return
            }

            userID := int(claims["user_id"].(float64))
            var u user.User
            if err := db.First(&u, userID).Error; err != nil {
                http.Error(w, "User not found", http.StatusUnauthorized)
                return
            }

            // Add user to context for downstream handlers
            ctx := r.Context()
            ctx = context.WithValue(ctx, "user", &u)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func AdminMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        u, ok := r.Context().Value("user").(*user.User)
        if !ok || u == nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if !u.IsAdmin {
            http.Error(w, "Admin access required", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}
// db/command.go
// Command DB Connection (Write Side)

package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectCommandDB creates Postgres connection for WRITE operations
func ConnectCommandDB(databaseURL string) *pgxpool.Pool {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal("❌ Unable to connect Command DB:", err)
	}

	// Ping check
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("❌ Command DB ping failed:", err)
	}

	log.Println("✅ Command DB Connected (Write Side)")
	return pool
}

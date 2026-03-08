// db/query.go
// Query DB Connection (Read Side)

package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectQueryDB creates Postgres connection for READ operations
func ConnectQueryDB(databaseURL string) *pgxpool.Pool {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal("❌ Unable to connect Query DB:", err)
	}

	// Ping check
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("❌ Query DB ping failed:", err)
	}

	log.Println("✅ Query DB Connected (Read Side)")
	return pool
}

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
// Configured for high concurrency (50K users support)
func ConnectCommandDB(databaseURL string) *pgxpool.Pool {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Parse config for connection pool settings
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatal("❌ Unable to parse Command DB config:", err)
	}

	// Configure connection pool for high concurrency
	poolConfig.MaxConns = 100          // Maximum number of connections
	poolConfig.MinConns = 20           // Minimum number of connections to maintain
	poolConfig.MaxConnLifetime = time.Hour    // Connection lifetime
	poolConfig.MaxConnIdleTime = 30 * time.Minute // Idle connection timeout
	poolConfig.HealthCheckPeriod = 5 * time.Minute // Health check interval

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatal("❌ Unable to connect Command DB:", err)
	}

	// Ping check
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("❌ Command DB ping failed:", err)
	}

	log.Println("✅ Command DB Connected (Write Side) - High concurrency pool configured")
	return pool
}

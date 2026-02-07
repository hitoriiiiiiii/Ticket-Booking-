//Postgress Connection

package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("Unable to connect DB:", err)
	}

	log.Println("✅ Connected to Postgres successfully")
	return pool
}

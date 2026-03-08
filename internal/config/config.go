package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CommandDatabaseURL string
	QueryDatabaseURL   string
	RedisURL           string
	Port               string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		CommandDatabaseURL: getEnv("COMMAND_DATABASE_URL", "postgres://ticket:ticket123@localhost:5433/ticket_cmd_db?sslmode=disable"),
		QueryDatabaseURL:   getEnv("QUERY_DATABASE_URL", "postgres://ticket:ticket123@localhost:5434/ticket_query_db?sslmode=disable"),
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		Port:               getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

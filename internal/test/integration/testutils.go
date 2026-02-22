// Package integration provides integration tests for the ticket booking system
package integration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hitorii/ticket-booking/internal/booking"
	"github.com/hitorii/ticket-booking/internal/events"
	"github.com/hitorii/ticket-booking/internal/queue"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

// TestDatabaseConfig holds database connection configuration for tests
type TestDatabaseConfig struct {
	Pool      *pgxpool.Pool
	Container *postgres.PostgresContainer
	DSN       string
	Host      string
	Port      int
}

// TestRedisConfig holds Redis connection configuration for tests
type TestRedisConfig struct {
	Client    *redis.Client
	Container *tcredis.RedisContainer
	URL       string
}

// TestServices holds all service instances for integration tests
type TestServices struct {
	CommandDB      *TestDatabaseConfig
	QueryDB        *TestDatabaseConfig
	Redis          *TestRedisConfig
	Dispatcher     *events.Dispatcher
	BookingCmdSvc  *booking.CommandService
	BookingQuerySvc *booking.QueryService
}

// isDockerError checks if the error is related to Docker not being available
func isDockerError(err string) bool {
	dockerErrors := []string{
		"docker",
		"container",
		"rootless",
		"elevated privileges",
		"pipe",
		"network",
	}
	errLower := strings.ToLower(err)
	for _, de := range dockerErrors {
		if strings.Contains(errLower, de) {
			return true
		}
	}
	return false
}

// SetupTestDatabase creates a PostgreSQL test container
func SetupTestDatabase(t *testing.T) (cfg *TestDatabaseConfig) {
	ctx := context.Background()

	// Skip if running in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("ticket_cmd_db"),
		postgres.WithUsername("ticket"),
		postgres.WithPassword("ticket123"),
	)
	if err != nil {
		// Check if it's a Docker-related error
		if isDockerError(err.Error()) {
			t.Skip("Docker is not available or not properly configured, skipping integration test")
		}
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	// Get host and port
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://ticket:ticket123@%s:%s/ticket_cmd_db?sslmode=disable", host, mappedPort.Port())

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Failed to create connection pool: %v", err)
	}

	// Run migrations
	if err := runMigrations(ctx, pool); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	cfg = &TestDatabaseConfig{
		Pool:      pool,
		Container: container,
		DSN:       dsn,
		Host:      host,
		Port:      int(mappedPort.Int()),
	}

	return cfg
}

// SetupTestRedis creates a Redis test container
func SetupTestRedis(t *testing.T) (cfg *TestRedisConfig) {
	ctx := context.Background()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	container, err := tcredis.RunContainer(ctx,
		testcontainers.WithImage("redis:7-alpine"),
	)
	if err != nil {
		if isDockerError(err.Error()) {
			t.Skip("Docker is not available or not properly configured, skipping integration test")
		}
		t.Fatalf("Failed to start Redis container: %v", err)
	}

	url, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get Redis connection string: %v", err)
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		t.Fatalf("Failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opts)

	// Wait for Redis to be ready
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	for {
		if err := client.Ping(ctx).Err(); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	cfg = &TestRedisConfig{
		Client:    client,
		Container: container,
		URL:       url,
	}

	return cfg
}

// SetupTestServices creates all necessary services for integration testing
func SetupTestServices(t *testing.T) *TestServices {
	ctx := context.Background()

	// Setup databases
	cmdDB := SetupTestDatabase(t)
	queryDB := SetupTestDatabase(t)

	// Update query DB name - create a new pool with different database
	queryDB.DSN = fmt.Sprintf("postgres://ticket:ticket123@%s:%d/ticket_query_db?sslmode=disable",
		queryDB.Host, queryDB.Port)

	// Create the query database
	_, err := cmdDB.Pool.Exec(ctx, "CREATE DATABASE ticket_query_db")
	if err != nil {
		// Database might already exist, ignore error
		t.Logf("Note: Query database might already exist: %v", err)
	}

	queryQ, err := pgxpool.New(ctx, queryDB.DSN)
	if err != nil {
		t.Fatalf("Failed to create query DB pool: %v", err)
	}
	queryDB.Pool = queryQ

	// Run migrations on query DB
	if err := runMigrations(ctx, queryDB.Pool); err != nil {
		t.Logf("Warning: Failed to run migrations on query DB: %v", err)
	}

	// Setup Redis
	redisCfg := SetupTestRedis(t)

	// Create event dispatcher
	dispatcher := events.NewDispatcher(redisCfg.Client, nil)

	// Create booking services
	bookingCmdSvc := booking.NewCommandServiceWithDispatcher(cmdDB.Pool, dispatcher)
	bookingQuerySvc := booking.NewQueryService(queryDB.Pool)

	return &TestServices{
		CommandDB:      cmdDB,
		QueryDB:        queryDB,
		Redis:          redisCfg,
		Dispatcher:     dispatcher,
		BookingCmdSvc:  bookingCmdSvc,
		BookingQuerySvc: bookingQuerySvc,
	}
}

// TeardownTestDatabase cleans up the PostgreSQL test container
func TeardownTestDatabase(t *testing.T, cfg *TestDatabaseConfig) {
	if cfg == nil || cfg.Container == nil {
		return
	}
	ctx := context.Background()
	if err := cfg.Container.Terminate(ctx); err != nil {
		t.Logf("Failed to terminate PostgreSQL container: %v", err)
	}
	if cfg.Pool != nil {
		cfg.Pool.Close()
	}
}

// TeardownTestRedis cleans up the Redis test container
func TeardownTestRedis(t *testing.T, cfg *TestRedisConfig) {
	if cfg == nil || cfg.Container == nil {
		return
	}
	ctx := context.Background()
	if err := cfg.Container.Terminate(ctx); err != nil {
		t.Logf("Failed to terminate Redis container: %v", err)
	}
	if cfg.Client != nil {
		cfg.Client.Close()
	}
}

// TeardownTestServices cleans up all test services
func TeardownTestServices(t *testing.T, svc *TestServices) {
	if svc == nil {
		return
	}
	TeardownTestDatabase(t, svc.CommandDB)
	TeardownTestDatabase(t, svc.QueryDB)
	TeardownTestRedis(t, svc.Redis)
	if svc.Dispatcher != nil {
		svc.Dispatcher.Close()
	}
}

// runMigrations executes database migrations for testing
func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	// Create reservations table
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS reservations (
			seat_id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			status TEXT NOT NULL,
			reserved_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create reservations table: %w", err)
	}

	// Create events table
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			aggregate_id TEXT NOT NULL,
			event_type TEXT NOT NULL,
			payload TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create events table: %w", err)
	}

	// Create users table
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			is_admin BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil
}

// SetupTestEnvironment sets up environment variables for testing
func SetupTestEnvironment() {
	os.Setenv("COMMAND_DATABASE_URL", "postgres://ticket:ticket123@localhost:5433/ticket_cmd_db?sslmode=disable")
	os.Setenv("QUERY_DATABASE_URL", "postgres://ticket:ticket123@localhost:5434/ticket_query_db?sslmode=disable")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("PORT", "8081")
}

// ClearRedisQueue clears all Redis queues for a clean test state
func ClearRedisQueue(ctx context.Context, rdb *redis.Client) error {
	// Clear job queue
	keys, err := rdb.Keys(ctx, "queue:*").Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		if err := rdb.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}

	// Clear event stream
	if err := rdb.Del(ctx, events.EventStreamName).Err(); err != nil {
		return err
	}

	return nil
}

// SetupIntegrationTest is a helper to run integration tests
// It handles setup and teardown automatically
func SetupIntegrationTest(t *testing.T) *TestServices {
	svc := SetupTestServices(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return nil
	}
	t.Cleanup(func() {
		TeardownTestServices(t, svc)
	})
	return svc
}

// TestFixture represents a test data fixture
type TestFixture struct {
	UserID string
	SeatID string
}

// CreateTestFixtures creates test data for integration tests
func CreateTestFixtures(svc *TestServices) (*TestFixture, error) {
	ctx := context.Background()

	fixture := &TestFixture{
		UserID: "test-user-1",
		SeatID: "seat-test-1",
	}

	// Clean up any existing test data
	_, _ = svc.CommandDB.Pool.Exec(ctx, "DELETE FROM reservations WHERE seat_id=$1 OR user_id=$1", fixture.SeatID)
	_, _ = svc.CommandDB.Pool.Exec(ctx, "DELETE FROM reservations WHERE user_id=$1", fixture.UserID)

	return fixture, nil
}

// WaitForEvent waits for an event to be published with timeout
func WaitForEvent(ctx context.Context, dispatcher *events.Dispatcher, eventType string, timeout time.Duration) error {
	eventCh := make(chan events.BaseEvent, 1)

	dispatcher.Subscribe(eventType, func(event events.BaseEvent) error {
		eventCh <- event
		return nil
	})

	select {
	case <-eventCh:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for event: %s", eventType)
	}
}

// GetJobFromQueue retrieves jobs from the Redis stream
func GetJobFromQueue(ctx context.Context, rdb *redis.Client) ([]redis.XMessage, error) {
	streams, err := rdb.XRead(ctx, &redis.XReadArgs{
		Streams: []string{queue.StreamName, "0"},
		Count:   1,
		Block:   2 * time.Second,
	}).Result()

	if err != nil {
		return nil, err
	}

	if len(streams) == 0 {
		return nil, nil
	}

	return streams[0].Messages, nil
}

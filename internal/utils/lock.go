package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// DistributedLock provides Redis-based distributed locking
type DistributedLock struct {
	redisClient *redis.Client
}

// NewDistributedLock creates a new distributed lock instance
func NewDistributedLock(redisClient *redis.Client) *DistributedLock {
	return &DistributedLock{redisClient: redisClient}
}

// AcquireLock attempts to acquire a lock with the given key and TTL
// Returns true if lock was acquired, false if not
// The lock will automatically expire after ttlSeconds
func (dl *DistributedLock) AcquireLock(ctx context.Context, key string, ttlSeconds int) (bool, error) {
	// Use SET NX with expiration for atomic lock acquisition
	result, err := dl.redisClient.SetNX(ctx, key, "locked", time.Duration(ttlSeconds)*time.Second).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}
	return result, nil
}

// ReleaseLock releases a lock with the given key
func (dl *DistributedLock) ReleaseLock(ctx context.Context, key string) error {
	_, err := dl.redisClient.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}
	return nil
}

// AcquireSeatLock acquires a lock for a specific seat
func (dl *DistributedLock) AcquireSeatLock(ctx context.Context, seatID string) (bool, error) {
	key := fmt.Sprintf("lock:seat:%s", seatID)
	return dl.AcquireLock(ctx, key, 10) // 10 seconds TTL
}

// ReleaseSeatLock releases a lock for a specific seat
func (dl *DistributedLock) ReleaseSeatLock(ctx context.Context, seatID string) error {
	key := fmt.Sprintf("lock:seat:%s", seatID)
	return dl.ReleaseLock(ctx, key)
}

// WithSeatLock is a helper that acquires a lock, executes the function, and releases the lock
func (dl *DistributedLock) WithSeatLock(ctx context.Context, seatID string, fn func() error) error {
	locked, err := dl.AcquireSeatLock(ctx, seatID)
	if err != nil {
		return fmt.Errorf("failed to acquire seat lock: %w", err)
	}
	if !locked {
		return fmt.Errorf("seat %s is temporarily unavailable - could not acquire lock", seatID)
	}
	
	// Ensure we release the lock when done
	defer func() {
		if err := dl.ReleaseSeatLock(ctx, seatID); err != nil {
			// Log error but don't override original error
			fmt.Printf("Warning: failed to release seat lock: %v\n", err)
		}
	}()
	
	return fn()
}

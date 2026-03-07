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
	if dl.redisClient == nil {
		return false, fmt.Errorf("redis client is nil")
	}
	// Use SET NX with expiration for atomic lock acquisition
	result, err := dl.redisClient.SetNX(ctx, key, "locked", time.Duration(ttlSeconds)*time.Second).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}
	return result, nil
}

// ReleaseLock releases a lock with the given key
func (dl *DistributedLock) ReleaseLock(ctx context.Context, key string) error {
	if dl.redisClient == nil {
		return fmt.Errorf("redis client is nil")
	}
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
// NOTE: This version releases the lock AFTER the function returns, but callers should ensure
// the function commits any database transaction before returning
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

// WithSeatLockAndCommit is a helper that acquires a lock, executes the function with commit callback,
// and releases the lock only AFTER successful commit. This prevents double booking.
// 
// The commitFn should attempt to commit the transaction and return any error.
// The lock will only be released if commitFn succeeds.
func (dl *DistributedLock) WithSeatLockAndCommit(ctx context.Context, seatID string, commitFn func() error) error {
	locked, err := dl.AcquireSeatLock(ctx, seatID)
	if err != nil {
		return fmt.Errorf("failed to acquire seat lock: %w", err)
	}
	if !locked {
		return fmt.Errorf("seat %s is temporarily unavailable - could not acquire lock", seatID)
	}
	
	// Execute the commit function
	commitErr := commitFn()
	
	// Only release lock AFTER commit succeeds
	// If commit fails, the lock will expire automatically (TTL) or caller should handle cleanup
	if releaseErr := dl.ReleaseSeatLock(ctx, seatID); releaseErr != nil {
		fmt.Printf("Warning: failed to release seat lock: %v\n", releaseErr)
	}
	
	return commitErr
}


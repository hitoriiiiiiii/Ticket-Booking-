package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache provides Redis-based caching functionality
type Cache struct {
	redisClient *redis.Client
	defaultTTL  time.Duration
}

// NewCache creates a new cache instance
func NewCache(redisClient *redis.Client, defaultTTL time.Duration) *Cache {
	return &Cache{
		redisClient: redisClient,
		defaultTTL:  defaultTTL,
	}
}

// Get retrieves a value from cache
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	result, err := c.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("cache miss for key: %s", key)
	}
	if err != nil {
		return fmt.Errorf("cache get error: %w", err)
	}

	return json.Unmarshal([]byte(result), dest)
}

// Set stores a value in cache with TTL
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}

	return c.redisClient.Set(ctx, key, data, ttl).Err()
}

// SetDefault uses the default TTL
func (c *Cache) SetDefault(ctx context.Context, key string, value interface{}) error {
	return c.Set(ctx, key, value, c.defaultTTL)
}

// Delete removes a key from cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.redisClient.Del(ctx, key).Err()
}

// Exists checks if a key exists
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// SeatAvailabilityCache provides caching for seat availability
type SeatAvailabilityCache struct {
	Cache *Cache
}

// NewSeatAvailabilityCache creates a seat availability cache
func NewSeatAvailabilityCache(redisClient *redis.Client) *SeatAvailabilityCache {
	return &SeatAvailabilityCache{
		Cache: NewCache(redisClient, 5*time.Second),
	}
}

// GetSeatAvailability retrieves cached seat availability
func (sac *SeatAvailabilityCache) GetSeatAvailability(ctx context.Context, showID string) (map[string]bool, error) {
	key := fmt.Sprintf("seat_availability:%s", showID)

	var availability map[string]bool
	err := sac.Cache.Get(ctx, key, &availability)
	if err != nil {
		return nil, err
	}
	return availability, nil
}

// SetSeatAvailability caches seat availability
func (sac *SeatAvailabilityCache) SetSeatAvailability(ctx context.Context, showID string, availability map[string]bool) error {
	key := fmt.Sprintf("seat_availability:%s", showID)
	return sac.Cache.Set(ctx, key, availability, 5*time.Second)
}

// InvalidateSeatAvailability removes cached seat availability
func (sac *SeatAvailabilityCache) InvalidateSeatAvailability(ctx context.Context, showID string) error {
	key := fmt.Sprintf("seat_availability:%s", showID)
	return sac.Cache.Delete(ctx, key)
}

// MovieCache provides caching for movie listings
type MovieCache struct {
	Cache *Cache
}

// NewMovieCache creates a movie cache
func NewMovieCache(redisClient *redis.Client) *MovieCache {
	return &MovieCache{
		Cache: NewCache(redisClient, 60*time.Second),
	}
}

// GetMovies retrieves cached movies
func (mc *MovieCache) GetMovies(ctx context.Context) ([]interface{}, error) {
	key := "movies:list"

	var movies []interface{}
	err := mc.Cache.Get(ctx, key, &movies)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

// SetMovies caches movies list
func (mc *MovieCache) SetMovies(ctx context.Context, movies []interface{}) error {
	key := "movies:list"
	return mc.Cache.Set(ctx, key, movies, 60*time.Second)
}

// InvalidateMovies removes cached movies
func (mc *MovieCache) InvalidateMovies(ctx context.Context) error {
	key := "movies:list"
	return mc.Cache.Delete(ctx, key)
}

// ShowCache provides caching for show listings
type ShowCache struct {
	Cache *Cache
}

// NewShowCache creates a show cache
func NewShowCache(redisClient *redis.Client) *ShowCache {
	return &ShowCache{
		Cache: NewCache(redisClient, 30*time.Second),
	}
}

// GetShows retrieves cached shows
func (sc *ShowCache) GetShows(ctx context.Context, movieID string) ([]interface{}, error) {
	key := fmt.Sprintf("shows:movie:%s", movieID)

	var shows []interface{}
	err := sc.Cache.Get(ctx, key, &shows)
	if err != nil {
		return nil, err
	}
	return shows, nil
}

// SetShows caches shows for a movie
func (sc *ShowCache) SetShows(ctx context.Context, movieID string, shows []interface{}) error {
	key := fmt.Sprintf("shows:movie:%s", movieID)
	return sc.Cache.Set(ctx, key, shows, 30*time.Second)
}

// InvalidateShows removes cached shows
func (sc *ShowCache) InvalidateShows(ctx context.Context, movieID string) error {
	key := fmt.Sprintf("shows:movie:%s", movieID)
	return sc.Cache.Delete(ctx, key)
}

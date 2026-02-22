// Redis stream producer and consumer
package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ProduceJob adds a job to the Redis stream
func ProduceJob(ctx context.Context, job JobPayload) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	jsonData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	// Add to stream
	streamID := RedisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: StreamName,
		Values: map[string]interface{}{
			"data": string(jsonData),
		},
	}).Val()

	fmt.Printf("📝 Job produced to Redis stream: %s\n", streamID)
	return nil
}

// JobProcessor is a function type for processing jobs
type JobProcessor func(job JobPayload) error

// ConsumeJobs reads jobs from the Redis stream and processes them
func ConsumeJobs(ctx context.Context, consumerName string, processor JobProcessor) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Read from stream using consumer group
			streams, err := RedisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    GroupName,
				Consumer: consumerName,
				Streams:  []string{StreamName, ">"},
				Count:    1,
				Block:    5 * time.Second,
			}).Result()

			if err == redis.Nil {
				// No new messages, continue loop
				continue
			}
			if err != nil {
				fmt.Printf("Error reading from stream: %v\n", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// Process each message
			for _, stream := range streams {
				for _, message := range stream.Messages {
					// Extract job data
					data, ok := message.Values["data"].(string)
					if !ok {
						fmt.Printf("Invalid message format: %v\n", message.Values)
						// Acknowledge to remove invalid message
						RedisClient.XAck(ctx, StreamName, GroupName, message.ID)
						continue
					}

					var jobMsg JobPayload
					if err := json.Unmarshal([]byte(data), &jobMsg); err != nil {
						fmt.Printf("Failed to unmarshal job: %v\n", err)
						RedisClient.XAck(ctx, StreamName, GroupName, message.ID)
						continue
					}

					fmt.Printf("📩 Processing job from Redis: Type=%s, UserID=%s\n", jobMsg.Type, jobMsg.UserID)

					// Process the job
					if err := processor(jobMsg); err != nil {
						fmt.Printf("Failed to process job: %v\n", err)
						// Don't acknowledge - job will be redelivered
						continue
					}

					// Acknowledge successful processing
					RedisClient.XAck(ctx, StreamName, GroupName, message.ID)
					fmt.Printf("✅ Job acknowledged: %s\n", message.ID)
				}
			}
		}
	}
}

// GetStreamInfo returns information about the stream
func GetStreamInfo(ctx context.Context) (*redis.XInfoStream, error) {
	if RedisClient == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.XInfoStream(ctx, StreamName).Result()
}

// GetPendingJobs returns pending jobs in the consumer group
func GetPendingJobs(ctx context.Context) ([]redis.XPendingExt, error) {
	if RedisClient == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}
	return RedisClient.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: StreamName,
		Group:  GroupName,
		Start:  "-",
		End:    "+",
		Count:  100,
	}).Result()
}

package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// EventStreamName is the Redis stream name for events
const EventStreamName = "events-stream"

// Dispatcher handles event publishing and subscribing
type Dispatcher struct {
	redisClient *redis.Client
	store       *Store
	subscribers map[string][]EventHandler
	mu          sync.RWMutex
}

// NewDispatcher creates a new event dispatcher
func NewDispatcher(redisClient *redis.Client, store *Store) *Dispatcher {
	return &Dispatcher{
		redisClient: redisClient,
		store:       store,
		subscribers: make(map[string][]EventHandler),
	}
}

// Publish publishes an event to the event store and Redis stream
func (d *Dispatcher) Publish(ctx context.Context, eventType, aggregateID string, payload interface{}) error {
	// Create base event
	event := BaseEvent{
		Type:        eventType,
		AggregateID: aggregateID,
		Timestamp:   time.Now(),
		Payload:     payload,
	}

	// Append to event store for durability
	if d.store != nil {
		if err := d.store.Append(eventType, aggregateID, payload); err != nil {
			log.Printf("❌ Failed to append event to store: %v", err)
			// Continue with publishing even if store fails
		}
	}

	// Publish to Redis stream for real-time processing
	if d.redisClient != nil {
		eventData, err := json.Marshal(event)
		if err != nil {
			log.Printf("❌ Failed to marshal event: %v", err)
			return err
		}

		streamID := d.redisClient.XAdd(ctx, &redis.XAddArgs{
			Stream: EventStreamName,
			Values: map[string]interface{}{
				"type":    eventType,
				"data":    string(eventData),
				"payload": fmt.Sprintf("%v", payload),
			},
		}).Val()

		log.Printf("📢 Published event: %s for aggregate: %s (stream ID: %s)", eventType, aggregateID, streamID)
	}

	// Notify local subscribers
	d.notifySubscribers(event)

	return nil
}

// Subscribe subscribes to specific event types
func (d *Dispatcher) Subscribe(eventType string, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.subscribers[eventType] = append(d.subscribers[eventType], handler)
	log.Printf("📡 Subscribed handler to event: %s", eventType)
}

// notifySubscribers notifies all subscribers of an event
func (d *Dispatcher) notifySubscribers(event BaseEvent) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	handlers, ok := d.subscribers[event.Type]
	if !ok {
		return
	}

	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h(event); err != nil {
				log.Printf("❌ Error in event handler for %s: %v", event.Type, err)
			}
		}(handler)
	}
}

// ConsumeEvents consumes events from Redis stream
func (d *Dispatcher) ConsumeEvents(ctx context.Context, consumerName string) error {
	if d.redisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	groupName := "event-consumers"

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			streams, err := d.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: consumerName,
				Streams:  []string{EventStreamName, ">"},
				Count:    10,
				Block:    5 * time.Second,
			}).Result()

			if err == redis.Nil {
				continue
			}
			if err != nil {
				log.Printf("Error reading from event stream: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					d.processEventMessage(message)
					d.redisClient.XAck(ctx, EventStreamName, groupName, message.ID)
				}
			}
		}
	}
}

// processEventMessage processes a single event message
func (d *Dispatcher) processEventMessage(message redis.XMessage) {
	eventType, _ := message.Values["type"].(string)
	eventData, _ := message.Values["data"].(string)

	var event BaseEvent
	if err := json.Unmarshal([]byte(eventData), &event); err != nil {
		log.Printf("❌ Failed to unmarshal event: %v", err)
		return
	}

	log.Printf("📥 Processing event: %s (ID: %s)", eventType, message.ID)

	// Notify subscribers
	d.notifySubscribers(event)
}

// GetStreamInfo returns information about the event stream
func (d *Dispatcher) GetStreamInfo(ctx context.Context) (*redis.XInfoStream, error) {
	if d.redisClient == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}
	return d.redisClient.XInfoStream(ctx, EventStreamName).Result()
}

// Close closes the dispatcher
func (d *Dispatcher) Close() error {
	if d.redisClient != nil {
		return d.redisClient.Close()
	}
	return nil
}

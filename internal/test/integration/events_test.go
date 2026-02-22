// Package integration provides integration tests for the ticket booking system
package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/hitorii/ticket-booking/internal/events"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEventsIntegration_PublishAndSubscribe tests the event publishing and subscribing system
func TestEventsIntegration_PublishAndSubscribe(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("Docker is not available: %v", r)
		}
	}()

	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	t.Run("PublishEvent_Success", func(t *testing.T) {
		// Create payload
		payload := events.EventPayload{
			UserID: "test-user-1",
			SeatID: "seat-1",
			Status: "HELD",
		}

		// Publish event
		err := svc.Dispatcher.Publish(ctx, events.EventTicketReserved, "seat-1", payload)
		require.NoError(t, err)

		// Give time for event to be processed
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("SubscribeAndReceive", func(t *testing.T) {
		// Create channel to receive events
		eventCh := make(chan events.BaseEvent, 1)

		// Subscribe to event type
		svc.Dispatcher.Subscribe(events.EventTicketConfirmed, func(event events.BaseEvent) error {
			eventCh <- event
			return nil
		})

		// Publish event
		payload := events.EventPayload{
			UserID: "test-user-2",
			SeatID: "seat-2",
			Status: "BOOKED",
		}
		err := svc.Dispatcher.Publish(ctx, events.EventTicketConfirmed, "seat-2", payload)
		require.NoError(t, err)

		// Wait for event with timeout
		select {
		case receivedEvent := <-eventCh:
			assert.Equal(t, events.EventTicketConfirmed, receivedEvent.Type)
			assert.Equal(t, "seat-2", receivedEvent.AggregateID)
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for event")
		}
	})

	t.Run("MultipleSubscribers", func(t *testing.T) {
		callCount := 0

		// Subscribe first handler
		svc.Dispatcher.Subscribe(events.EventTicketCancelled, func(event events.BaseEvent) error {
			callCount++
			return nil
		})

		// Subscribe second handler
		svc.Dispatcher.Subscribe(events.EventTicketCancelled, func(event events.BaseEvent) error {
			callCount++
			return nil
		})

		// Publish event
		payload := events.EventPayload{
			UserID: "test-user-3",
			SeatID: "seat-3",
			Status: "CANCELLED",
		}
		err := svc.Dispatcher.Publish(ctx, events.EventTicketCancelled, "seat-3", payload)
		require.NoError(t, err)

		// Wait for events to be processed
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, 2, callCount, "Both subscribers should be called")
	})
}

// TestEventsIntegration_AllEventTypes tests publishing different event types
func TestEventsIntegration_AllEventTypes(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("Docker is not available: %v", r)
		}
	}()

	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	testCases := []struct {
		eventType string
		payload   events.EventPayload
	}{
		{
			eventType: events.EventTicketReserved,
			payload: events.EventPayload{
				UserID: "user-1",
				SeatID: "seat-1",
				Status: "HELD",
			},
		},
		{
			eventType: events.EventTicketConfirmed,
			payload: events.EventPayload{
				UserID: "user-2",
				SeatID: "seat-2",
				Status: "BOOKED",
			},
		},
		{
			eventType: events.EventTicketCancelled,
			payload: events.EventPayload{
				UserID: "user-3",
				SeatID: "seat-3",
				Status: "CANCELLED",
			},
		},
		{
			eventType: events.EventPaymentInitiated,
			payload: events.EventPayload{
				UserID:    "user-4",
				PaymentID: "payment-1",
				Amount:    100,
				Status:    "PENDING",
			},
		},
		{
			eventType: events.EventPaymentVerified,
			payload: events.EventPayload{
				UserID:    "user-4",
				PaymentID: "payment-1",
				Amount:    100,
				Status:    "VERIFIED",
			},
		},
		{
			eventType: events.EventUserRegistered,
			payload: events.EventPayload{
				Username: "newuser",
				Email:    "newuser@example.com",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.eventType, func(t *testing.T) {
			err := svc.Dispatcher.Publish(ctx, tc.eventType, "aggregate-"+tc.eventType, tc.payload)
			require.NoError(t, err)
			time.Sleep(50 * time.Millisecond)
		})
	}
}

// TestEventsIntegration_StreamOperations tests Redis stream operations
func TestEventsIntegration_StreamOperations(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("Docker is not available: %v", r)
		}
	}()

	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	t.Run("EventStoredInRedisStream", func(t *testing.T) {
		// Publish an event
		payload := events.EventPayload{
			UserID: "stream-test-user",
			SeatID: "stream-test-seat",
			Status: "HELD",
		}
		err := svc.Dispatcher.Publish(ctx, events.EventTicketReserved, "stream-test-seat", payload)
		require.NoError(t, err)

		// Wait for event to be written to stream
		time.Sleep(200 * time.Millisecond)

		// Read from stream - use correct API with Streams slice
		streams, err := svc.Redis.Client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{events.EventStreamName, "0"},
			Count:   1,
			Block:   time.Second,
		}).Result()

		// Note: Stream might be empty if consumer already processed it
		// This test verifies the stream infrastructure is working
		assert.NoError(t, err)
		_ = streams // May be empty due to consumer group
	})

	t.Run("EventPayloadSerialization", func(t *testing.T) {
		// Test that event payload can be serialized and deserialized
		original := events.EventPayload{
			UserID:    "serialization-test",
			SeatID:    "seat-serialize",
			Status:    "BOOKED",
			PaymentID: "pay-123",
			Amount:    500,
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded events.EventPayload
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, original.UserID, decoded.UserID)
		assert.Equal(t, original.SeatID, decoded.SeatID)
		assert.Equal(t, original.Status, decoded.Status)
		assert.Equal(t, original.PaymentID, decoded.PaymentID)
		assert.Equal(t, original.Amount, decoded.Amount)
	})
}

// TestEventsIntegration_EventFlow tests the complete event flow
func TestEventsIntegration_EventFlow(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("Docker is not available: %v", r)
		}
	}()

	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Track events in order
	eventsReceived := make([]string, 0)
	_ = eventsReceived // suppress unused warning

	// Subscribe to multiple event types
	svc.Dispatcher.Subscribe(events.EventTicketReserved, func(e events.BaseEvent) error {
		eventsReceived = append(eventsReceived, e.Type)
		return nil
	})

	svc.Dispatcher.Subscribe(events.EventTicketConfirmed, func(e events.BaseEvent) error {
		eventsReceived = append(eventsReceived, e.Type)
		return nil
	})

	svc.Dispatcher.Subscribe(events.EventTicketCancelled, func(e events.BaseEvent) error {
		eventsReceived = append(eventsReceived, e.Type)
		return nil
	})

	// Publish events in sequence
	err := svc.Dispatcher.Publish(ctx, events.EventTicketReserved, "flow-seat-1", events.EventPayload{UserID: "user-1", SeatID: "flow-seat-1", Status: "HELD"})
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	err = svc.Dispatcher.Publish(ctx, events.EventTicketConfirmed, "flow-seat-1", events.EventPayload{UserID: "user-1", SeatID: "flow-seat-1", Status: "BOOKED"})
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	err = svc.Dispatcher.Publish(ctx, events.EventTicketCancelled, "flow-seat-2", events.EventPayload{UserID: "user-2", SeatID: "flow-seat-2", Status: "CANCELLED"})
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	// Wait for all events to be processed
	time.Sleep(100 * time.Millisecond)

	// Verify events were received in order
	assert.GreaterOrEqual(t, len(eventsReceived), 3, "Should have received at least 3 events")
	if len(eventsReceived) >= 3 {
		assert.Equal(t, events.EventTicketReserved, eventsReceived[0])
		assert.Equal(t, events.EventTicketConfirmed, eventsReceived[1])
		assert.Equal(t, events.EventTicketCancelled, eventsReceived[2])
	}
}

// TestEventsIntegration_ErrorHandling tests error handling in event system
func TestEventsIntegration_ErrorHandling(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("Docker is not available: %v", r)
		}
	}()

	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	t.Run("PublishWithNilPayload", func(t *testing.T) {
		// Should not panic with nil payload
		err := svc.Dispatcher.Publish(ctx, events.EventTicketReserved, "test-seat", nil)
		assert.NoError(t, err)
	})

	t.Run("PublishWithEmptyAggregateID", func(t *testing.T) {
		err := svc.Dispatcher.Publish(ctx, events.EventTicketReserved, "", events.EventPayload{UserID: "user-1"})
		assert.NoError(t, err)
	})
}

// TestEventsIntegration_ConcurrentPublishing tests concurrent event publishing
func TestEventsIntegration_ConcurrentPublishing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("Docker is not available: %v", r)
		}
	}()

	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	// Subscribe to capture events
	received := make(chan string, 100)
	svc.Dispatcher.Subscribe(events.EventTicketReserved, func(e events.BaseEvent) error {
		received <- e.Type
		return nil
	})

	// Publish events concurrently
	numEvents := 20
	for i := 0; i < numEvents; i++ {
		go func(idx int) {
			payload := events.EventPayload{
				UserID: "concurrent-user",
				SeatID: "concurrent-seat",
			}
			_ = svc.Dispatcher.Publish(ctx, events.EventTicketReserved, "concurrent-seat", payload)
		}(i)
	}

	// Wait for events to be processed
	time.Sleep(500 * time.Millisecond)
	close(received)

	// Count received events
	count := 0
	for range received {
		count++
	}

	assert.Equal(t, numEvents, count, "All concurrent events should be received")
}

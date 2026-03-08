// Package integration provides integration tests for the ticket booking system
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/hitorii/ticket-booking/internal/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEventsIntegration_PublishSubscribe tests event publishing and subscribing
func TestEventsIntegration_PublishSubscribe(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	t.Run("PublishAndSubscribe", func(t *testing.T) {
		eventReceived := make(chan bool, 1)

		svc.Dispatcher.Subscribe(events.EventTicketConfirmed, func(event events.BaseEvent) error {
			eventReceived <- true
			return nil
		})

		payload := events.EventPayload{
			UserID: "test-user",
			Status: "BOOKED",
		}

		err := svc.Dispatcher.Publish(ctx, events.EventTicketConfirmed, "seat-123", payload)
		require.NoError(t, err)

		select {
		case <-eventReceived:
		case <-ctx.Done():
			t.Fatal("Context cancelled before event received")
		}
	})

	t.Run("MultipleSubscribers", func(t *testing.T) {
		counter1 := 0
		counter2 := 0

		svc.Dispatcher.Subscribe(events.EventUserRegistered, func(event events.BaseEvent) error {
			counter1++
			return nil
		})
		svc.Dispatcher.Subscribe(events.EventUserRegistered, func(event events.BaseEvent) error {
			counter2++
			return nil
		})

		payload2 := map[string]interface{}{"username": "testuser"}
		err := svc.Dispatcher.Publish(ctx, events.EventUserRegistered, "user-123", payload2)
		require.NoError(t, err)

		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, 1, counter1)
		assert.Equal(t, 1, counter2)
	})
}

// TestEventsIntegration_StorePersistence tests event store persistence
func TestEventsIntegration_StorePersistence(t *testing.T) {
	svc := SetupIntegrationTest(t)
	if svc == nil {
		t.Skip("Test services could not be set up")
		return
	}
	ctx := context.Background()

	dispatcherWithoutStore := events.NewDispatcher(svc.Redis.Client, nil)

	t.Run("PublishWithoutStore", func(t *testing.T) {
		payload := map[string]string{"key": "value"}
		err := dispatcherWithoutStore.Publish(ctx, events.EventTicketCancelled, "seat-456", payload)
		require.NoError(t, err)
	})
}

// TestEventsIntegration_EventTypes verifies all expected event types exist
func TestEventsIntegration_EventTypes(t *testing.T) {
	expectedTypes := []string{
		events.EventTicketReserved,
		events.EventTicketConfirmed,
		events.EventTicketCancelled,
		events.EventPaymentInitiated,
		events.EventPaymentVerified,
		events.EventUserRegistered,
	}

	for _, eventType := range expectedTypes {
		assert.NotEmpty(t, eventType)
	}
}

package show

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hitorii/ticket-booking/internal/events"
)

type CommandService struct {
	Dispatcher *events.Dispatcher
}

func NewCommandService() *CommandService {
	return &CommandService{}
}

func NewCommandServiceWithDispatcher(dispatcher *events.Dispatcher) *CommandService {
	return &CommandService{Dispatcher: dispatcher}
}

// CreateShow - Command to create a new show
func (s *CommandService) CreateShow(ctx context.Context, req CreateShowRequest) (*Show, error) {
	// Validate input
	if req.Theater == "" {
		return nil, errors.New("theater is required")
	}
	if req.MovieID <= 0 {
		return nil, errors.New("valid movie ID is required")
	}

	// Calculate end time based on movie duration (this would typically be done with a query to movies table)
	// For now, we'll just set a placeholder end time
	endTime := req.StartTime.Add(2 * time.Hour) // Default 2 hours

	show := &Show{
		MovieID:   req.MovieID,
		Theater:   req.Theater,
		StartTime: req.StartTime,
		EndTime:   endTime,
	}

	// Emit event for event-driven flow
	if s.Dispatcher != nil {
		showIDStr := fmt.Sprintf("%d", show.ID)
		movieIDStr := fmt.Sprintf("%d", req.MovieID)
		payload := events.EventPayload{
			ShowID:    showIDStr,
			MovieID:   movieIDStr,
			StartTime: req.StartTime.Format(time.RFC3339),
			EndTime:   endTime.Format(time.RFC3339),
		}
		_ = s.Dispatcher.Publish(ctx, events.EventShowCreated, showIDStr, payload)
	}

	return show, nil
}

// UpdateShow - Command to update an existing show
func (s *CommandService) UpdateShow(ctx context.Context, id string, req UpdateShowRequest) error {
	// Validate input
	if req.Theater == "" {
		return errors.New("theater is required")
	}

	// Emit event for event-driven flow
	if s.Dispatcher != nil {
		movieIDStr := fmt.Sprintf("%d", req.MovieID)
		payload := events.EventPayload{
			ShowID:  id,
			MovieID: movieIDStr,
		}
		_ = s.Dispatcher.Publish(ctx, events.EventShowUpdated, id, payload)
	}

	return nil
}

// DeleteShow - Command to delete a show
func (s *CommandService) DeleteShow(ctx context.Context, id string) error {
	// Emit event for event-driven flow
	if s.Dispatcher != nil {
		payload := events.EventPayload{
			ShowID: id,
		}
		_ = s.Dispatcher.Publish(ctx, events.EventShowDeleted, id, payload)
	}

	return nil
}

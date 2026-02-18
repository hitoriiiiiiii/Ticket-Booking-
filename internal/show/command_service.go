// Command service for show write operations (CQRS)

package show

import (
	"context"
	"errors"
	"time"
)

type CommandService struct {
	// Could include event store or other command-related dependencies
}

func NewCommandService() *CommandService {
	return &CommandService{}
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

	return show, nil
}

// UpdateShow - Command to update an existing show
func (s *CommandService) UpdateShow(ctx context.Context, id string, req UpdateShowRequest) error {
	// Validate input
	if req.Theater == "" {
		return errors.New("theater is required")
	}

	// The actual database operation would be handled by a repository
	return nil
}

// DeleteShow - Command to delete a show
func (s *CommandService) DeleteShow(ctx context.Context, id string) error {
	// The actual database operation would be handled by a repository
	return nil
}

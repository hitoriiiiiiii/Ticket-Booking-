package movie

import (
	"context"
	"errors"
	"fmt"

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

// CreateMovieRequest - Request model for creating a movie
type CreateMovieRequest struct {
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	Duration int    `json:"duration"`
}

// UpdateMovieRequest - Request model for updating a movie
type UpdateMovieRequest struct {
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	Duration int    `json:"duration"`
}

// CreateMovie - Command to create a new movie
func (s *CommandService) CreateMovie(ctx context.Context, req CreateMovieRequest) (*Movie, error) {
	// Validate input
	if req.Name == "" {
		return nil, errors.New("movie name is required")
	}
	if req.Duration <= 0 {
		return nil, errors.New("movie duration must be positive")
	}

	// The actual database operation would be handled by a repository
	// Here we just return a movie object with the request data
	// The repository would handle the actual INSERT operation
	movie := &Movie{
		Name:     req.Name,
		Genre:    req.Genre,
		Duration: req.Duration,
	}

	// Emit event for event-driven flow
	if s.Dispatcher != nil {
		movieIDStr := fmt.Sprintf("%d", movie.ID)
		payload := events.EventPayload{
			MovieID: movieIDStr,
			Name:    req.Name,
		}
		_ = s.Dispatcher.Publish(ctx, events.EventMovieCreated, movieIDStr, payload)
	}

	return movie, nil
}

// UpdateMovie - Command to update an existing movie
func (s *CommandService) UpdateMovie(ctx context.Context, id string, req UpdateMovieRequest) error {
	// Validate input
	if req.Name == "" {
		return errors.New("movie name is required")
	}
	if req.Duration <= 0 {
		return errors.New("movie duration must be positive")
	}

	// Emit event for event-driven flow
	if s.Dispatcher != nil {
		payload := events.EventPayload{
			MovieID: id,
			Name:    req.Name,
		}
		_ = s.Dispatcher.Publish(ctx, events.EventMovieUpdated, id, payload)
	}

	return nil
}

// DeleteMovie - Command to delete a movie
func (s *CommandService) DeleteMovie(ctx context.Context, id string) error {
	// Emit event for event-driven flow
	if s.Dispatcher != nil {
		payload := events.EventPayload{
			MovieID: id,
		}
		_ = s.Dispatcher.Publish(ctx, events.EventMovieDeleted, id, payload)
	}

	return nil
}

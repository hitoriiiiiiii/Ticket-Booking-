// Query service for show read operations (CQRS)

package show

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryService struct {
	DB *pgxpool.Pool
}

func NewQueryService(db *pgxpool.Pool) *QueryService {
	return &QueryService{DB: db}
}

// GetAllShows - Query to get all shows
func (s *QueryService) GetAllShows(ctx context.Context) ([]Show, error) {
	rows, err := s.DB.Query(ctx, "SELECT id, movie_id, theater, start_time, end_time FROM shows")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []Show
	for rows.Next() {
		var sh Show
		if err := rows.Scan(&sh.ID, &sh.MovieID, &sh.Theater, &sh.StartTime, &sh.EndTime); err != nil {
			return nil, err
		}
		shows = append(shows, sh)
	}

	return shows, nil
}

// GetShowByID - Query to get a single show by ID
func (s *QueryService) GetShowByID(ctx context.Context, id string) (*Show, error) {
	var sh Show
	err := s.DB.QueryRow(ctx, "SELECT id, movie_id, theater, start_time, end_time FROM shows WHERE id=$1", id).Scan(&sh.ID, &sh.MovieID, &sh.Theater, &sh.StartTime, &sh.EndTime)
	if err != nil {
		return nil, err
	}
	return &sh, nil
}

// GetShowsByMovie - Query to get shows by movie ID
func (s *QueryService) GetShowsByMovie(ctx context.Context, movieID string) ([]Show, error) {
	rows, err := s.DB.Query(ctx, "SELECT id, movie_id, theater, start_time, end_time FROM shows WHERE movie_id=$1", movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []Show
	for rows.Next() {
		var sh Show
		if err := rows.Scan(&sh.ID, &sh.MovieID, &sh.Theater, &sh.StartTime, &sh.EndTime); err != nil {
			return nil, err
		}
		shows = append(shows, sh)
	}

	return shows, nil
}

// Query service for movie read operations (CQRS)

package movie

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

// GetAllMovies - Query to get all movies
func (s *QueryService) GetAllMovies(ctx context.Context) ([]Movie, error) {
	rows, err := s.DB.Query(ctx, "SELECT id, name, genre, duration FROM movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Name, &m.Genre, &m.Duration); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

// GetMovieByID - Query to get a single movie by ID
func (s *QueryService) GetMovieByID(ctx context.Context, id string) (*Movie, error) {
	var m Movie
	err := s.DB.QueryRow(ctx, "SELECT id, name, genre, duration FROM movies WHERE id=$1", id).Scan(&m.ID, &m.Name, &m.Genre, &m.Duration)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetMoviesByGenre - Query to get movies by genre
func (s *QueryService) GetMoviesByGenre(ctx context.Context, genre string) ([]Movie, error) {
	rows, err := s.DB.Query(ctx, "SELECT id, name, genre, duration FROM movies WHERE genre=$1", genre)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Name, &m.Genre, &m.Duration); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

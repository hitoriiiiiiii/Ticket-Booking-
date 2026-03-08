// Query handler for movie read operations (CQRS)

package movie

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryHandler struct {
	QueryService *QueryService
}

func NewQueryHandler(qs *QueryService) *QueryHandler {
	return &QueryHandler{QueryService: qs}
}

// GetMovies - Query handler for listing all movies
func (h *QueryHandler) GetMovies(c *gin.Context) {
	movies, err := h.QueryService.GetAllMovies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}
	c.JSON(http.StatusOK, movies)
}

// GetMovie - Query handler for getting a single movie by ID
func (h *QueryHandler) GetMovie(c *gin.Context) {
	id := c.Param("id")
	
	movie, err := h.QueryService.GetMovieByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}
	c.JSON(http.StatusOK, movie)
}

// Initialize QueryHandler with database pool (for compatibility with existing setup)
func NewQueryHandlerWithDB(db *pgxpool.Pool) *QueryHandler {
	qs := NewQueryService(db)
	return &QueryHandler{QueryService: qs}
}

//Handler for movie related endpoints

package movie

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	DB *pgxpool.Pool
}

func (h *Handler) GetMovies(c *gin.Context) {
	rows, err := h.DB.Query(c.Request.Context(), "SELECT id, name, genre, duration FROM movies")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Name, &m.Genre, &m.Duration); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse movie data"})
			return
		}
		movies = append(movies, m)
	}
	c.JSON(http.StatusOK, movies)
}

func (h *Handler) GetMovie(c *gin.Context) {
	id := c.Param("id")

	row := h.DB.QueryRow(c.Request.Context(), "SELECT id, name, genre, duration FROM movies WHERE id=$1", id)
	var m Movie
	if err := row.Scan(&m.ID, &m.Name, &m.Genre, &m.Duration); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *Handler) CreateMovie(c *gin.Context) {
	var m Movie
	if err := c.BindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.DB.QueryRow(c.Request.Context(),
		"INSERT INTO movies (name, genre, duration) VALUES ($1, $2, $3) RETURNING id",
		m.Name, m.Genre, m.Duration,
	).Scan(&m.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
		return
	}
	c.JSON(http.StatusCreated, m)
}

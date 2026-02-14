// Shows Handler for show related endpoints

package show

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) GetShows(c *gin.Context) {
	rows, err := h.Service.DB.Query(c.Request.Context(), "SELECT id, movie_id, theater, start_time, end_time FROM shows")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shows"})
		return
	}
	defer rows.Close()

	var shows []Show
	for rows.Next() {
		var s Show
		if err := rows.Scan(&s.ID, &s.MovieID, &s.Theater, &s.StartTime, &s.EndTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse show data"})
			return
		}
		shows = append(shows, s)
	}
	c.JSON(http.StatusOK, shows)
}

func (h *Handler) CreateShow(c *gin.Context) {
	var s Show
	if err := c.BindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Optional: calculate end_time based on duration
	var duration int
	err := h.Service.DB.QueryRow(c.Request.Context(), "SELECT duration FROM movies WHERE id=$1", s.MovieID).Scan(&duration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie not found"})
		return
	}
	s.EndTime = s.StartTime.Add(time.Duration(duration) * time.Minute)

	err = h.Service.DB.QueryRow(c.Request.Context(),
		"INSERT INTO shows (movie_id, theater, start_time, end_time) VALUES ($1, $2, $3, $4) RETURNING id",
		s.MovieID, s.Theater, s.StartTime, s.EndTime,
	).Scan(&s.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create show"})
		return
	}
	c.JSON(http.StatusCreated, s)
}

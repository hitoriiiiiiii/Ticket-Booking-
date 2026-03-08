package show

import "time"

type Show struct {
	ID        int       `json:"id"`
	MovieID   int       `json:"movie_id"`
	Theater   string    `json:"theater"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

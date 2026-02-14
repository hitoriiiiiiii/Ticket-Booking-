//Models for Movie database

package movie

type Movie struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Genre    string `json:"genre"`
	Duration int    `json:"duration"` // in minutes
}

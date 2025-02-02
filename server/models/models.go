package models

import "time"

// Repository represents repository data
type Repository struct {
	ID                 string    `json:"id"`
	Author             string    `json:"author"`
	Name               string    `json:"name"`
	URL                string    `json:"url"`
	Description        string    `json:"description"`
	Language           string    `json:"language"`
	Stars              int       `json:"stars"`
	Forks              int       `json:"forks"`
	CurrentPeriodStars int       `json:"current_period_stars"`
	UpdatedAt          time.Time `json:"updated_at"`
	// Additional interest metric (sum of stars and forks)
	InterestScore float64 `json:"interest_score"`
}

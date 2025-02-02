package store

import (
	"time"

	"github.com/estnafinema0/Github-Trends-Aggregator/server/models"
)

type Store struct {
	repos map[string]models.Repository
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) UpdateRepos(newRepos []models.Repository) {
	for _, repo := range newRepos {
		repo.UpdatedAt = time.Now()
		repo.InterestScore = float64(repo.Stars + repo.Forks)

		s.repos[repo.ID] = repo
	}
}

// GetRepos returns repositories by language
func (s *Store) GetRepos(language string) []models.Repository {
	result := []models.Repository{}

	for _, repo := range s.repos {
		if repo.Language == language {
			result = append(result, repo)
		}
	}

	return result
}

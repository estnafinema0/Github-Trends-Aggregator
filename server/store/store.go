package store

import (
	"sync"
	"time"

	"github.com/estnafinema0/Github-Trends-Aggregator/server/models"
)

// Store represents a structure for storing repositories
type Store struct {
	sync.RWMutex
	repos map[string]models.Repository
}

// NewStore creates a new Store instance
func NewStore() *Store {
	return &Store{
		repos: make(map[string]models.Repository),
	}
}

// UpdateRepos updates repositories
func (s *Store) UpdateRepos(newRepos []models.Repository) {
	s.Lock()
	defer s.Unlock()

	for _, repo := range newRepos {
		repo.UpdatedAt = time.Now()
		repo.InterestScore = float64(repo.Stars + repo.Forks)

		s.repos[repo.ID] = repo
	}
}

// GetRepos retrieves repositories by language
func (s *Store) GetRepos(language string) []models.Repository {
	s.RLock()
	defer s.RUnlock()

	var repos []models.Repository

	// Filter repositories by language
	for _, repo := range s.repos {
		if repo.Language == language {
			repos = append(repos, repo)
		}
	}

	return repos
}

package store

import (
	"sync"
	"time"

	"github.com/estnafinema0/Github-Trends-Aggregator/server/models"
)

type Store struct {
	sync.RWMutex
	repos map[string]models.Repository
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) UpdateRepos(newRepos []models.Repository) {
	s.Lock()
	defer s.Unlock()

	for _, repo := range newRepos {
		repo.UpdatedAt = time.Now()
		repo.InterestScore = float64(repo.Stars + repo.Forks)

		s.repos[repo.ID] = repo
	}
}

// GetRepos retrieves repositories based on the specified programming language
func (s *Store) GetRepos(language string) []models.Repository {
	s.RLock()
	defer s.RUnlock()

	var repos []models.Repository

	// Filter repos based on language
	for _, repo := range s.repos {
		if repo.Language == language {
			repos = append(repos, repo)
		}
	}

	return repos
}

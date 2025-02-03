package store

import (
	"sort"
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

// GetReposFiltered returns a list of repositories with filtering by language and sorting
func (s *Store) GetReposFiltered(language string, sortBy string) []models.Repository {
	s.RLock()
	defer s.RUnlock()

	result := []models.Repository{}

	for _, repo := range s.repos {
		if language == "" || repo.Language == language {
			result = append(result, repo)
		}
	}

	switch sortBy {
	case "stars":
		sort.Slice(result, func(i, j int) bool {
			return result[i].CurrentPeriodStars > result[j].CurrentPeriodStars
		})
	case "forks":
		sort.Slice(result, func(i, j int) bool {
			return result[i].Forks > result[j].Forks
		})
	case "current_period_stars":
		sort.Slice(result, func(i, j int) bool {
			return result[i].CurrentPeriodStars > result[j].CurrentPeriodStars
		})
	case "interest_score":
		sort.Slice(result, func(i, j int) bool {
			return result[i].InterestScore > result[j].InterestScore
		})
	}
	return result
}

// GetRepoByID returns a repository by its ID (e.g., "author/name")
func (s *Store) GetRepoByID(id string) (models.Repository, bool) {
	s.RLock()
	defer s.RLock()
	repo, ok := s.repos[id]
	return repo, ok

}

// GetAllRepos returns all repositories
func (s *Store) GetAllRepos() []models.Repository {
	s.RLock()
	defer s.RUnlock()

	var repos []models.Repository
	for _, repo := range s.repos {
		repos = append(repos, repo)
	}
	return repos
}

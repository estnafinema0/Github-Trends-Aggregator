package store

import (
	"sort"
	"sync"
	"time"

	"github.com/estnafinema0/Github-Trends-Aggregator/server/models"
	"github.com/estnafinema0/Github-Trends-Aggregator/server/config"
)

// Store represents a structure for storing repositories
type Store struct {
	sync.RWMutex
	repos       map[string]models.Repository
	NotifsList  map[int]string
	invRepos    map[int]models.Repository
	lastRepoID  int
	lastNotifID int
	reposHist   []models.StarsTimestamp
}

// NewStore creates a new Store instance
func NewStore() *Store {
	return &Store{
		repos:       make(map[string]models.Repository),
		NotifsList:  make(map[int]string),
		invRepos:    make(map[int]models.Repository),
		lastRepoID:  1,
		lastNotifID: 1,
	}
}

// UpdateRepos updates repositories
func (s *Store) UpdateRepos(newRepos []models.Repository) {
	s.Lock()
	defer s.Unlock()
	for k := range s.invRepos {
		delete(s.invRepos, k)
	}

	total := len(s.repos)
	var sumStars int

	for _, repo := range newRepos {
		repo.UpdatedAt = time.Now()
		repo.InterestScore = float64(repo.Stars + repo.Forks)
		sumStars += repo.Stars
		
		s.repos[repo.ID] = repo
		s.invRepos[repo.SecondaryID] = repo
	}
	avgStars := 0.0
	if total > 0 {
		avgStars = float64(sumStars) / float64(total)
	}
	if avgStars != 0.0 {
		s.reposHist = append(s.reposHist, models.StarsTimestamp { time.Now().Format("2006-01-02T15:04:05-0700"), avgStars })
		if len(s.reposHist) > config.HistoryLength {
			s.reposHist = s.reposHist[1:]
		}
	}
}

func (s *Store) UpdateNotifs(newNotif string) {
	s.Lock()
	defer s.Unlock()

	s.lastNotifID += 1
	s.NotifsList[s.lastNotifID] = newNotif
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

func (s *Store) GetRepoBySecondaryID(id int) (models.Repository, bool) {
	s.RLock()
	defer s.RLock()
	repo, ok := s.invRepos[id]
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

func (s *Store) GetTopRatedStatistics() (models.Repository, []models.StarsTimestamp) {
	var topRepo models.Repository
	for _, repo := range s.repos {
		if repo.InterestScore > topRepo.InterestScore {
			topRepo = repo
		}
	}
	return topRepo, s.reposHist
}

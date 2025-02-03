package api

import (
	"encoding/json"
	"net/http"

	"github.com/estnafinema0/Github-Trends-Aggregator/server/models"
	"github.com/estnafinema0/Github-Trends-Aggregator/server/store"
	"github.com/gorilla/mux"
)

// GetTrendsHandler returns handler for GET /trends
func GetTrendsHandler(s *store.Store) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		language := r.URL.Query().Get("language")
		sortBy := r.URL.Query().Get("sort_by")
		repos := s.GetReposFiltered(language, sortBy)
		rw.Header().Set("Content-Type", "application/json")

		json.NewEncoder(rw).Encode(repos)
	}
}

// GetRepoHandler returns detailed information about a repository, selected by its ID
func GetRepoHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		repo, ok := s.GetRepoByID(id)
		if !ok {
			http.Error(w, "Repository not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(repo)
	}
}

// Stats contains aggregated statistical data
type Stats struct {
	TotalRepos       int               `json:"total_repos"`
	AverageStars     float64           `json:"average_stars"`
	AverageForks     float64           `json:"average_forks"`
	ReposPerLanguage map[string]int    `json:"repos_per_language"`
	TopInterestRepo  models.Repository `json:"top_interest_repo"`
}

// GetStatsHandler returns aggregated statistics on repositories
func GetStatsHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.RLock()
		defer s.RUnlock()

		repos := s.GetAllRepos()
		total := len(repos)
		var sumStars, sumForks int
		reposPerLang := make(map[string]int)
		var topRepo models.Repository
		for _, repo := range repos {
			sumStars += repo.Stars
			sumForks += repo.Forks
			if repo.Language != "" {
				reposPerLang[repo.Language]++
			}
			if repo.InterestScore > topRepo.InterestScore {
				topRepo = repo
			}
		}
		avgStars := 0.0
		avgForks := 0.0
		if total > 0 {
			avgStars = float64(sumStars) / float64(total)
			avgForks = float64(sumForks) / float64(total)
		}
		stats := Stats{
			TotalRepos:       total,
			AverageStars:     avgStars,
			AverageForks:     avgForks,
			ReposPerLanguage: reposPerLang,
			TopInterestRepo:  topRepo,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}

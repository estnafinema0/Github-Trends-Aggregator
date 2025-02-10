package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
	"fmt"
	"strings"

	"github.com/estnafinema0/Github-Trends-Aggregator/server/models"
	"github.com/estnafinema0/Github-Trends-Aggregator/server/store"
	"github.com/gorilla/mux"
)

func EscapePlus(s string) string {
	return strings.ReplaceAll(s, "+", "%2b")
}

// GetTrendsHandler returns handler for GET /trends
func GetTrendsHandler(s *store.Store) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("static/trends.html"))

	return func(rw http.ResponseWriter, r *http.Request) {
		language := r.URL.Query().Get("language")
		sortBy := r.URL.Query().Get("sort_by")
		repos := s.GetReposFiltered(language, sortBy)
		type Temporary struct {
			Repos []models.Repository
			SelectedLanguage string
			Urlescape func(string) string
		}
		data := Temporary{repos, language, EscapePlus}
		fmt.Println(tmpl.Execute(rw, data))

		// rw.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(rw).Encode(repos)
	}
}

func GetIndexHandler(s *store.Store) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("static/index.html"))

	return func(rw http.ResponseWriter, r *http.Request) {
		repos := s.GetReposFiltered("", "stars")
		mlangs := make(map[string]bool)
		for _, v := range repos {
			mlangs[v.Language] = true
		}
		type TemporaryString struct { Language string }
		var langs []TemporaryString
		for k := range mlangs {
			langs = append(langs, TemporaryString{k})
		}
		type Temporary struct { Langs []TemporaryString }
		data := Temporary{langs}
		tmpl.Execute(rw, data)
		
		// rw.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(rw).Encode(repos)
	}
}

// GetRepoHandler returns detailed information about a repository, selected by its ID
func GetRepoHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		repo, ok := s.GetRepoBySecondaryID(id)
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

package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
	"strings"

	"github.com/estnafinema0/Github-Trends-Aggregator/server/models"
	"github.com/estnafinema0/Github-Trends-Aggregator/server/store"
	"github.com/gorilla/mux"
)

func Escape(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "+", "%2b"), "#", "%23")
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
		data := Temporary{repos, language, Escape}
		tmpl.Execute(rw, data)
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
	tmpl := template.Must(template.ParseFiles("static/stats.html"))

	return func(rw http.ResponseWriter, r *http.Request) {
		s.RLock()
		defer s.RUnlock()
		repo, list := s.GetTopRatedStatistics()
		type Temporary struct { Repo models.Repository; History []models.StarsTimestamp }
		data := Temporary{ repo, list }
		tmpl.Execute(rw, data)
	}
}

func GetSubscribeHandler(s *store.Store) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("static/subscribe.html"))

	return func(rw http.ResponseWriter, r *http.Request) {
		tmpl.Execute(rw, nil)
	}
}


func GetSubscribedHandler(s *store.Store) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("static/subscribed.html"))

	return func(rw http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		s.UpdateNotifs(email)
		tmpl.Execute(rw, nil)
	}
}
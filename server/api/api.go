package api

import (
	"encoding/json"
	"net/http"

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

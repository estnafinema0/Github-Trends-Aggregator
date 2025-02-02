package api

import (
	"encoding/json"
	"net/http"

	"github.com/estnafinema0/Github-Trends-Aggregator/server/store"
)

// GetTrendsHandler returns handler for GET /trends
func GetTrendsHandler(s *store.Store) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		language := r.URL.Query().Get("language")
		repos := s.GetRepos(language)
		rw.Header().Set("Content-Type", "application/json")

		json.NewEncoder(rw).Encode(repos)
	}
}

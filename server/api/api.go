package api

import (
	"net/http"
	// "github.com/estnafinema0/Github-Trends-Aggregator/server"
)

// GetTrendsHandler returns handler for GET /trends
func GetTrendsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//language := r.URL.Query().Get("language")

		w.Header().Set("Content-Type", "application/json")

	}
}

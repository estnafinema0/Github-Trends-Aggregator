package scheduler

import (
	"log"
	"time"

	"github.com/estnafinema0/Github-Trends-Aggregator/server/config"
	"github.com/estnafinema0/Github-Trends-Aggregator/server/fetcher"
	"github.com/estnafinema0/Github-Trends-Aggregator/server/store"
	"github.com/estnafinema0/Github-Trends-Aggregator/server/ws"
)

// StartScheduler starts the periodic fetching of trends
func StartScheduler(store *store.Store, hub *ws.Hub, l *log.Logger) {
	ticker := time.NewTicker(config.FetchInterval)
	defer ticker.Stop()

	// First run immediately
	for {
		l.Println("Starting trend fetching...")
		repos, err := fetcher.FetchTrendingRepos(l)
		if err != nil {
			l.Printf("Error fetching trends: %v\n", err)
		} else {
			store.UpdateRepos(repos)
			// Notify connected clients via WebSocket
			hub.Broadcast(l, repos)
		}
		<-ticker.C
	}
}

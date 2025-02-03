package fetcher

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/estnafinema0/Github-Trends-Aggregator/server/models"
)

const trendingURL = "https://github.com/trending"

// FetchTrendingRepos fetches trending repositories from GitHub Trending, parses HTML and returns list of repositories
func FetchTrendingRepos(l *log.Logger) ([]models.Repository, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(trendingURL)
	if err != nil {
		log.Printf("Error fetching data: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error parsing HTML: %v\n", err)
		return nil, err
	}

	var repos []models.Repository
	doc.Find("article.Box-row").Each(func(i int, s *goquery.Selection) {
		// Extract title containing author and repository name
		titleSel := s.Find("h2 > a")
		href, exists := titleSel.Attr("href")
		if !exists {
			l.Println("Failed to get href for repository")
			return
		}
		titleText := strings.TrimSpace(titleSel.Text())
		// Split string by "/" to get author and repository name
		parts := strings.Split(titleText, "/")
		if len(parts) < 2 {
			l.Println("Failed to get author and repository name")
			return
		}
		author := strings.TrimSpace(parts[0])
		name := strings.TrimSpace(parts[1])
		url := "https://github.com" + strings.TrimSpace(href)

		// l.Printf("Fetched repository: %s\n", url)

		// Repository description (if present)
		desc := strings.TrimSpace(s.Find("p.col-9.color-fg-muted.my-1.pr-4").Text())
		// l.Printf("Repository description: %s\n", desc)

		// Programming language
		language := strings.TrimSpace(s.Find("span[itemprop='programmingLanguage']").Text())
		// l.Printf("Programming language: %s\n", language)

		// Number of stars – select first link with "stargazers" in href
		starsStr := strings.TrimSpace(s.Find("a[href*='/stargazers']").First().Text())
		starsStr = strings.ReplaceAll(starsStr, ",", "")
		stars, _ := strconv.Atoi(starsStr)
		// l.Printf("Number of stars: %d\n", stars)

		// Number of forks – select first link with "network/members" in href
		forksStr := strings.TrimSpace(s.Find("a[href*='/network/members']").First().Text())
		forksStr = strings.ReplaceAll(forksStr, ",", "")
		forks, _ := strconv.Atoi(forksStr)
		// l.Printf("Number of forks: %d\n", forks)

		// Stars for current period – text from `span` with class "d-inline-block float-sm-right"
		periodStarsStr := strings.TrimSpace(s.Find("span.d-inline-block.float-sm-right").Text())
		// Usually string like "1,234 stars today" – take first part and remove commas
		periodStarsStr = strings.Split(periodStarsStr, " ")[0]
		periodStarsStr = strings.ReplaceAll(periodStarsStr, ",", "")
		currentPeriodStars, _ := strconv.Atoi(periodStarsStr)
		// l.Printf("Stars for current period: %d\n", currentPeriodStars)
		repo := models.Repository{
			ID:                 author + "/" + name,
			Author:             author,
			Name:               name,
			URL:                url,
			Description:        desc,
			Language:           language,
			Stars:              stars,
			Forks:              forks,
			CurrentPeriodStars: currentPeriodStars,
		}
		repos = append(repos, repo)
		// l.Printf("Added repository: %s\n", repo.ID)
	})

	l.Printf("Trending repos fetched: %d\n", len(repos))
	return repos, nil
}

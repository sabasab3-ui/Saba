package intelligence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WikipediaResearcher struct {
	Client     *http.Client
	MaxResults int
}

type wikiSearchResponse struct {
	Query struct {
		Search []struct {
			Title   string `json:"title"`
			Snippet string `json:"snippet"`
			PageID  int    `json:"pageid"`
		} `json:"search"`
	} `json:"query"`
}

func NewWikipediaResearcher() *WikipediaResearcher {
	return &WikipediaResearcher{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		MaxResults: 5,
	}
}

func (r *WikipediaResearcher) Research(
	ctx context.Context,
	question string,
) ([]Source, error) {

	question = strings.TrimSpace(question)

	if question == "" {
		return nil, fmt.Errorf("research question cannot be empty")
	}

	client := r.Client
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	maxResults := r.MaxResults
	if maxResults <= 0 {
		maxResults = 5
	}

	apiURL := "https://en.wikipedia.org/w/api.php"

	params := url.Values{}
	params.Set("action", "query")
	params.Set("list", "search")
	params.Set("srsearch", question)
	params.Set("format", "json")
	params.Set("utf8", "1")
	params.Set("srlimit", fmt.Sprintf("%d", maxResults))

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		apiURL+"?"+params.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"User-Agent",
		"SABA-Intelligence/1.0",
	)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("research request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"research service returned status %d",
			resp.StatusCode,
		)
	}

	var data wikiSearchResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf(
			"failed to decode research response: %w",
			err,
		)
	}

	sources := make([]Source, 0, len(data.Query.Search))

	for _, result := range data.Query.Search {
		title := strings.TrimSpace(result.Title)
		snippet := cleanHTML(result.Snippet)

		if title == "" {
			continue
		}

		pageURL := "https://en.wikipedia.org/?curid=" +
			fmt.Sprintf("%d", result.PageID)

		sources = append(sources, Source{
			Title:   title,
			URL:     pageURL,
			Content: snippet,
		})
	}

	if len(sources) == 0 {
		return nil, fmt.Errorf(
			"no research sources found for %q",
			question,
		)
	}

	return sources, nil
}

func cleanHTML(text string) string {
	text = strings.ReplaceAll(text, "<span class=\"searchmatch\">", "")
	text = strings.ReplaceAll(text, "</span>", "")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&#39;", "'")

	return strings.TrimSpace(text)
}

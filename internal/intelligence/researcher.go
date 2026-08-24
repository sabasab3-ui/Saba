package intelligence

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"sort"
	"strconv"
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

type wikiPage struct {
	PageID  int    `json:"pageid"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
	FullURL string `json:"fullurl"`
}

type wikiExtractResponse struct {
	Query struct {
		Pages map[string]wikiPage `json:"pages"`
	} `json:"query"`
}

type rankedWikiResult struct {
	Title   string
	Snippet string
	PageID  int
	Score   int
}

func NewWikipediaResearcher() *WikipediaResearcher {
	return &WikipediaResearcher{
		Client:     &http.Client{Timeout: 15 * time.Second},
		MaxResults: 8,
	}
}

func (r *WikipediaResearcher) Research(ctx context.Context, question string) ([]Source, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return nil, fmt.Errorf("research question cannot be empty")
	}

	client := r.Client
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	maxResults := r.MaxResults
	if maxResults <= 0 {
		maxResults = 8
	}

	results, err := searchWikipedia(ctx, client, question, maxResults)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no research sources found for %q", question)
	}

	// Rank by relevance and remove obvious false positives such as films,
	// songs and other pages that happen to share the query words.
	tokens := meaningfulTokens(question)
	ranked := make([]rankedWikiResult, 0, len(results))
	for _, result := range results {
		ranked = append(ranked, rankedWikiResult{
			Title: result.Title, Snippet: result.Snippet, PageID: result.PageID,
			Score: relevanceScore(result.Title, result.Snippet, tokens),
		})
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		return ranked[i].Score > ranked[j].Score
	})

	selected := make([]rankedWikiResult, 0, maxResults)
	for _, item := range ranked {
		if item.Score <= 0 && len(selected) >= 2 {
			continue
		}
		if isObviousNoise(item.Title, question) && len(selected) < 2 {
			continue
		}
		selected = append(selected, item)
		if len(selected) >= maxResults {
			break
		}
	}

	if len(selected) == 0 {
		// Keep the best result rather than failing a legitimate but unusual query.
		selected = append(selected, ranked[0])
	}

	pages, err := fetchWikipediaExtracts(ctx, client, selected)
	if err != nil {
		// Search snippets are still useful if the page-extract request fails.
		return snippetsToSources(selected), nil
	}

	sources := make([]Source, 0, len(selected))
	for _, item := range selected {
		page, ok := pages[strconv.Itoa(item.PageID)]
		if !ok {
			sources = append(sources, Source{
				Title:   item.Title,
				URL:     "https://en.wikipedia.org/?curid=" + strconv.Itoa(item.PageID),
				Content: cleanHTML(item.Snippet),
			})
			continue
		}

		content := strings.TrimSpace(page.Extract)
		if content == "" {
			content = cleanHTML(item.Snippet)
		}

		pageURL := page.FullURL
		if pageURL == "" {
			pageURL = "https://en.wikipedia.org/?curid=" + strconv.Itoa(item.PageID)
		}

		sources = append(sources, Source{
			Title:   page.Title,
			URL:     pageURL,
			Content: content,
		})
	}

	if len(sources) == 0 {
		return nil, fmt.Errorf("no usable research sources found for %q", question)
	}

	return deduplicateSources(sources), nil
}

func searchWikipedia(ctx context.Context, client *http.Client, question string, maxResults int) ([]struct {
	Title   string
	Snippet string
	PageID  int
}, error) {
	params := url.Values{}
	params.Set("action", "query")
	params.Set("list", "search")
	params.Set("srsearch", question)
	params.Set("format", "json")
	params.Set("utf8", "1")
	params.Set("srlimit", strconv.Itoa(maxResults))
	params.Set("srprop", "snippet|titlesnippet")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://en.wikipedia.org/w/api.php?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "SABA-Intelligence/2.0 (research service)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("research request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("research service returned status %d", resp.StatusCode)
	}

	var data wikiSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode research response: %w", err)
	}

	results := make([]struct {
		Title   string
		Snippet string
		PageID  int
	}, 0, len(data.Query.Search))
	for _, result := range data.Query.Search {
		if strings.TrimSpace(result.Title) == "" || result.PageID == 0 {
			continue
		}
		results = append(results, struct {
			Title   string
			Snippet string
			PageID  int
		}{result.Title, result.Snippet, result.PageID})
	}
	return results, nil
}

func fetchWikipediaExtracts(ctx context.Context, client *http.Client, selected []rankedWikiResult) (map[string]wikiPage, error) {
	ids := make([]string, 0, len(selected))
	for _, item := range selected {
		ids = append(ids, strconv.Itoa(item.PageID))
	}

	params := url.Values{}
	params.Set("action", "query")
	params.Set("pageids", strings.Join(ids, "|"))
	params.Set("prop", "extracts|info")
	params.Set("exintro", "1")
	params.Set("explaintext", "1")
	params.Set("inprop", "url")
	params.Set("format", "json")
	params.Set("utf8", "1")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://en.wikipedia.org/w/api.php?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "SABA-Intelligence/2.0 (research service)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("research extract service returned status %d", resp.StatusCode)
	}

	var data wikiExtractResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Query.Pages, nil
}

func meaningfulTokens(text string) []string {
	stop := map[string]bool{
		"what": true, "when": true, "where": true, "which": true, "who": true,
		"why": true, "how": true, "does": true, "did": true, "can": true,
		"could": true, "should": true, "would": true, "the": true, "and": true,
		"for": true, "with": true, "from": true, "into": true, "about": true,
		"this": true, "that": true, "are": true, "is": true, "was": true,
		"were": true, "today": true, "now": true, "best": true,
	}
	words := strings.Fields(strings.ToLower(text))
	seen := make(map[string]bool)
	out := make([]string, 0, len(words))
	for _, word := range words {
		word = strings.Trim(word, ".,!?;:()[]{}\"'")
		if len(word) < 3 || stop[word] || seen[word] {
			continue
		}
		seen[word] = true
		out = append(out, word)
	}
	return out
}

func relevanceScore(title, snippet string, tokens []string) int {
	text := strings.ToLower(title + " " + snippet)
	titleText := strings.ToLower(title)
	score := 0
	for _, token := range tokens {
		if strings.Contains(titleText, token) {
			score += 4
		} else if strings.Contains(text, token) {
			score++
		}
	}
	return score
}

func isObviousNoise(title, question string) bool {
	if strings.Contains(strings.ToLower(question), "movie") || strings.Contains(strings.ToLower(question), "film") {
		return false
	}
	lower := strings.ToLower(title)
	noise := []string{"film", "movie", "song", "album", "novel", "television series", "episode"}
	for _, word := range noise {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}

func snippetsToSources(results []rankedWikiResult) []Source {
	sources := make([]Source, 0, len(results))
	for _, result := range results {
		sources = append(sources, Source{
			Title:   result.Title,
			URL:     "https://en.wikipedia.org/?curid=" + strconv.Itoa(result.PageID),
			Content: cleanHTML(result.Snippet),
		})
	}
	return sources
}

func deduplicateSources(sources []Source) []Source {
	seen := make(map[string]bool)
	out := make([]Source, 0, len(sources))
	for _, source := range sources {
		key := strings.ToLower(strings.TrimSpace(source.Title))
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, source)
	}
	return out
}

func cleanHTML(text string) string {
	text = strings.ReplaceAll(text, "<span class=\"searchmatch\">", "")
	text = strings.ReplaceAll(text, "</span>", "")
	text = html.UnescapeString(text)
	return strings.Join(strings.Fields(strings.TrimSpace(text)), " ")
}

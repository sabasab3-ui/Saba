package intelligence

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

var errEmptyResearchQuestion = errors.New("empty research question")

type WebResearcher struct {
	Client     *http.Client
	MaxResults int
}

func NewWebResearcher() *WebResearcher {
	return &WebResearcher{
		Client:     &http.Client{Timeout: 10 * time.Second},
		MaxResults: 10,
	}
}

type ddgResponse struct {
	AbstractText string `json:"AbstractText"`
	AbstractURL  string `json:"AbstractURL"`
	Heading      string `json:"Heading"`
	Answer       string `json:"Answer"`
	Results      []struct {
		FirstURL string `json:"FirstURL"`
		Text     string `json:"Text"`
		Result   string `json:"Result"`
	} `json:"Results"`
	RelatedTopics []struct {
		FirstURL string `json:"FirstURL"`
		Text     string `json:"Text"`
		Topics   []struct {
			FirstURL string `json:"FirstURL"`
			Text     string `json:"Text"`
		} `json:"Topics"`
	} `json:"RelatedTopics"`
}

type googleNewsRSS struct {
	Channel struct {
		Items []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (r *WebResearcher) Research(ctx context.Context, question string) ([]Source, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return nil, errEmptyResearchQuestion
	}

	client := r.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	maxResults := r.MaxResults
	if maxResults <= 0 {
		maxResults = 10
	}

	queries := buildResearchQueries(question)

	type result struct {
		sources []Source
	}
	results := make(chan result, len(queries))
	var wg sync.WaitGroup

	for _, query := range queries {
		query := query
		wg.Add(1)
		go func() {
			defer wg.Done()

			var combined []Source

			if items, err := searchDuck(ctx, client, query); err == nil {
				combined = append(combined, items...)
			}
			if items, err := searchGoogleNews(ctx, client, query); err == nil {
				combined = append(combined, items...)
			}

			results <- result{sources: combined}
		}()
	}

	wg.Wait()
	close(results)

	var all []Source
	for item := range results {
		all = append(all, item.sources...)
	}

	all = dedupeSources(all)
	all = rankForQuestion(question, all)

	if len(all) > maxResults {
		all = all[:maxResults]
	}

	if len(all) == 0 {
		return nil, fmt.Errorf("no relevant web evidence found")
	}

	return all, nil
}

func buildResearchQueries(question string) []string {
	q := strings.TrimSpace(question)
	lower := strings.ToLower(q)
	queries := []string{q}

	if strings.Contains(lower, "africa") {
		queries = append(queries,
			"artificial intelligence Africa opportunities",
			"AI Africa agriculture healthcare education finance",
			"AI Africa jobs startups digital transformation",
		)
	}
	if strings.Contains(lower, "uganda") {
		queries = append(queries,
			"AI Uganda opportunities",
			"artificial intelligence Uganda",
		)
	}
	if strings.Contains(lower, "kenya") {
		queries = append(queries, "AI Kenya", "artificial intelligence Kenya")
	}
	if strings.Contains(lower, "nigeria") {
		queries = append(queries, "AI Nigeria", "artificial intelligence Nigeria")
	}

	seen := map[string]bool{}
	out := make([]string, 0, len(queries))
	for _, item := range queries {
		key := strings.ToLower(strings.TrimSpace(item))
		if key != "" && !seen[key] {
			seen[key] = true
			out = append(out, item)
		}
	}
	return out
}

func searchDuck(ctx context.Context, client *http.Client, query string) ([]Source, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("format", "json")
	params.Set("no_html", "1")
	params.Set("skip_disambig", "1")
	params.Set("no_redirect", "1")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.duckduckgo.com/?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "SABA-Intelligence/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DuckDuckGo HTTP %d", resp.StatusCode)
	}

	var data ddgResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var out []Source
	add := func(title, link, content string) {
		title = strings.TrimSpace(stripMarkup(title))
		link = strings.TrimSpace(link)
		content = strings.TrimSpace(stripMarkup(content))
		if title == "" || link == "" || content == "" {
			return
		}
		out = append(out, Source{Title: title, URL: link, Content: content})
	}

	if data.AbstractText != "" && data.AbstractURL != "" {
		add(data.Heading, data.AbstractURL, data.AbstractText)
	}
	if data.Answer != "" {
		add("DuckDuckGo web answer", "https://duckduckgo.com/", data.Answer)
	}
	for _, item := range data.Results {
		add(item.Result, item.FirstURL, item.Text)
	}
	for _, topic := range data.RelatedTopics {
		add(topic.Text, topic.FirstURL, topic.Text)
		for _, nested := range topic.Topics {
			add(nested.Text, nested.FirstURL, nested.Text)
		}
	}

	return dedupeSources(out), nil
}

func searchGoogleNews(ctx context.Context, client *http.Client, query string) ([]Source, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("hl", "en-US")
	params.Set("gl", "US")
	params.Set("ceid", "US:en")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://news.google.com/rss/search?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "SABA-Intelligence/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google News HTTP %d", resp.StatusCode)
	}

	var feed googleNewsRSS
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}

	out := make([]Source, 0, len(feed.Channel.Items))
	for _, item := range feed.Channel.Items {
		title := strings.TrimSpace(stripMarkup(item.Title))
		link := strings.TrimSpace(item.Link)
		description := strings.TrimSpace(stripMarkup(item.Description))

		if title == "" || link == "" {
			continue
		}
		if description == "" {
			description = title
		}

		content := description
		if item.PubDate != "" {
			content += " Published: " + item.PubDate
		}

		out = append(out, Source{Title: title, URL: link, Content: content})
		if len(out) >= 15 {
			break
		}
	}

	return dedupeSources(out), nil
}

func rankForQuestion(question string, sources []Source) []Source {
	tokens := meaningfulResearchTokens(question)
	locationTokens := researchLocationTokens(question)

	type scored struct {
		source Source
		score  int
	}

	ranked := make([]scored, 0, len(sources))

	for _, source := range sources {
		title := strings.ToLower(source.Title)
		content := strings.ToLower(source.Content)
		full := title + " " + content
		score := 0

		for _, token := range tokens {
			if strings.Contains(title, token) {
				score += 8
			} else if strings.Contains(content, token) {
				score += 2
			}
		}

		for _, token := range locationTokens {
			if strings.Contains(full, token) {
				score += 10
			}
		}

		if len(tokens) > 0 {
			matches := 0
			for _, token := range tokens {
				if strings.Contains(full, token) {
					matches++
				}
			}
			if matches == 0 {
				score -= 20
			}
		}

		ranked = append(ranked, scored{source: source, score: score})
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].score == ranked[j].score {
			return len(ranked[i].source.Content) > len(ranked[j].source.Content)
		}
		return ranked[i].score > ranked[j].score
	})

	out := make([]Source, 0, len(ranked))
	for _, item := range ranked {
		if item.score >= 2 {
			out = append(out, item.source)
		}
	}

	return out
}

func meaningfulResearchTokens(question string) []string {
	stop := map[string]bool{
		"what": true, "are": true, "the": true, "most": true, "important": true,
		"things": true, "should": true, "do": true, "to": true, "be": true,
		"become": true, "a": true, "an": true, "and": true, "of": true,
		"for": true, "in": true, "on": true, "with": true, "how": true,
		"is": true, "it": true, "used": true, "today": true, "can": true,
		"could": true, "would": true, "into": true, "from": true,
	}

	seen := map[string]bool{}
	var out []string

	for _, word := range strings.Fields(strings.ToLower(question)) {
		word = strings.Trim(word, ".,!?;:\"'()[]{}")
		if len(word) < 3 || stop[word] || seen[word] {
			continue
		}
		seen[word] = true
		out = append(out, word)
	}

	return out
}

func researchLocationTokens(question string) []string {
	lower := strings.ToLower(question)

	switch {
	case strings.Contains(lower, "africa"):
		return []string{"africa", "african", "sub-saharan"}
	case strings.Contains(lower, "uganda"):
		return []string{"uganda", "ugandan"}
	case strings.Contains(lower, "kenya"):
		return []string{"kenya", "kenyan"}
	case strings.Contains(lower, "nigeria"):
		return []string{"nigeria", "nigerian"}
	default:
		return nil
	}
}

func dedupeSources(sources []Source) []Source {
	seenURL := map[string]bool{}
	seenText := map[string]bool{}
	out := make([]Source, 0, len(sources))

	for _, source := range sources {
		u := strings.ToLower(strings.TrimSpace(source.URL))
		t := strings.ToLower(strings.Join(strings.Fields(source.Title+" "+source.Content), " "))

		if u == "" || seenURL[u] {
			continue
		}
		if t != "" && seenText[t] {
			continue
		}

		seenURL[u] = true
		if t != "" {
			seenText[t] = true
		}
		out = append(out, source)
	}

	return out
}

func stripMarkup(text string) string {
	replacements := []string{
		"<span class=\"searchmatch\">", "",
		"</span>", "",
		"<b>", "", "</b>", "",
		"<em>", "", "</em>", "",
	}

	for i := 0; i+1 < len(replacements); i += 2 {
		text = strings.ReplaceAll(text, replacements[i], replacements[i+1])
	}

	return strings.Join(strings.Fields(html.UnescapeString(text)), " ")
}

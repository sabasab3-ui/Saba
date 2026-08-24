package intelligence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

// WebResearcher gathers evidence from more than one public web source.
type WebResearcher struct {
	Client     *http.Client
	MaxResults int
}

func NewWebResearcher() *WebResearcher {
	return &WebResearcher{
		Client:     &http.Client{Timeout: 15 * time.Second},
		MaxResults: 10,
	}
}

type duckResponse struct {
	AbstractText string `json:"AbstractText"`
	AbstractURL  string `json:"AbstractURL"`
	Heading      string `json:"Heading"`
	Answer       string `json:"Answer"`
	Results      []struct {
		FirstURL string `json:"FirstURL"`
		Text     string `json:"Text"`
		Result   string `json:"Result"`
	} `json:"Results"`
	RelatedTopics []duckTopic `json:"RelatedTopics"`
}

type duckTopic struct {
	FirstURL string      `json:"FirstURL"`
	Text     string      `json:"Text"`
	Result   string      `json:"Result"`
	Topics   []duckTopic `json:"Topics"`
}

// Research runs the public web providers in parallel, tolerates one provider
// being unavailable, ranks evidence, and removes duplicate URLs.
func (r *WebResearcher) Research(ctx context.Context, question string) ([]Source, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return nil, ErrEmptyQuestion
	}

	client := r.Client
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	maxResults := r.MaxResults
	if maxResults <= 0 {
		maxResults = 10
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	var sources []Source
	var failures []string

	wg.Add(2)

	go func() {
		defer wg.Done()
		items, err := searchDuckDuckGo(ctx, client, question)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			failures = append(failures, "web search: "+err.Error())
			return
		}
		sources = append(sources, items...)
	}()

	go func() {
		defer wg.Done()
		items, err := searchWikipedia(ctx, client, question, maxResults)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			failures = append(failures, "Wikipedia: "+err.Error())
			return
		}
		for _, item := range items {
			sources = append(sources, Source{
				Title:   item.Title,
				URL:     "https://en.wikipedia.org/?curid=" + fmt.Sprintf("%d", item.PageID),
				Content: cleanHTML(item.Snippet),
			})
		}
	}()

	wg.Wait()

	sources = rankSources(question, deduplicateByURL(sources))
	if len(sources) > maxResults {
		sources = sources[:maxResults]
	}

	if len(sources) == 0 {
		if len(failures) > 0 {
			return nil, fmt.Errorf("all research providers failed: %s", strings.Join(failures, "; "))
		}
		return nil, fmt.Errorf("no research sources found for %q", question)
	}

	return sources, nil
}

func searchDuckDuckGo(ctx context.Context, client *http.Client, question string) ([]Source, error) {
	params := url.Values{}
	params.Set("q", question)
	params.Set("format", "json")
	params.Set("no_html", "1")
	params.Set("skip_disambig", "1")
	params.Set("no_redirect", "1")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.duckduckgo.com/?"+params.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "SABA-Intelligence/3.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var data duckResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	out := make([]Source, 0, 16)

	add := func(title, link, content string) {
		title = strings.TrimSpace(title)
		link = strings.TrimSpace(link)
		content = strings.TrimSpace(strings.Join(strings.Fields(content), " "))

		if title == "" || link == "" || content == "" {
			return
		}

		out = append(out, Source{
			Title:   title,
			URL:     link,
			Content: content,
		})
	}

	if data.AbstractText != "" && data.AbstractURL != "" {
		add(data.Heading, data.AbstractURL, data.AbstractText)
	}

	if data.Answer != "" {
		add("DuckDuckGo answer", "https://duckduckgo.com/", data.Answer)
	}

	for _, item := range data.Results {
		add(stripHTML(item.Result), item.FirstURL, stripHTML(item.Text))
	}

	for _, topic := range data.RelatedTopics {
		collectDuckTopic(&out, topic)
	}

	return deduplicateByURL(out), nil
}

func collectDuckTopic(out *[]Source, topic duckTopic) {
	if topic.FirstURL != "" && topic.Text != "" {
		*out = append(*out, Source{
			Title:   firstSentence(stripHTML(topic.Text)),
			URL:     topic.FirstURL,
			Content: stripHTML(topic.Text),
		})
	}

	for _, nested := range topic.Topics {
		collectDuckTopic(out, nested)
	}
}

func rankSources(question string, sources []Source) []Source {
	tokens := meaningfulTokens(question)

	type scored struct {
		source Source
		score  int
	}

	ranked := make([]scored, 0, len(sources))

	for _, source := range sources {
		text := strings.ToLower(source.Title + " " + source.Content)
		title := strings.ToLower(source.Title)
		score := 0

		for _, token := range tokens {
			if strings.Contains(title, token) {
				score += 5
			} else if strings.Contains(text, token) {
				score++
			}
		}

		if strings.Contains(source.URL, "wikipedia.org") {
			score++
		}

		ranked = append(ranked, scored{
			source: source,
			score:  score,
		})
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].score == ranked[j].score {
			return len(ranked[i].source.Content) > len(ranked[j].source.Content)
		}
		return ranked[i].score > ranked[j].score
	})

	out := make([]Source, 0, len(ranked))
	for _, item := range ranked {
		out = append(out, item.source)
	}

	return out
}

func deduplicateByURL(sources []Source) []Source {
	seen := make(map[string]bool)
	out := make([]Source, 0, len(sources))

	for _, source := range sources {
		key := strings.ToLower(strings.TrimSpace(source.URL))

		if key == "" || seen[key] {
			continue
		}

		seen[key] = true
		out = append(out, source)
	}

	return out
}

func stripHTML(text string) string {
	text = strings.ReplaceAll(text, "<b>", "")
	text = strings.ReplaceAll(text, "</b>", "")
	text = strings.ReplaceAll(text, "<span>", "")
	text = strings.ReplaceAll(text, "</span>", "")
	return cleanHTML(text)
}

func firstSentence(text string) string {
	text = strings.TrimSpace(text)

	for i, r := range text {
		if (r == '.' || r == '!' || r == '?') && i > 0 {
			return strings.TrimSpace(text[:i])
		}
	}

	return text
}

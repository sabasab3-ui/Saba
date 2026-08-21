package agent

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type WebScraperTool struct{}

func (t *WebScraperTool) Name() string {
	return "scraper"
}

func (t *WebScraperTool) Description() string {
	return "Fetch and extract readable text from a public webpage"
}

func (t *WebScraperTool) Execute(ctx context.Context, params map[string]interface{}) error {
	url, ok := params["url"].(string)
	if !ok || url == "" {
		return fmt.Errorf("missing url parameter")
	}

	if !strings.HasPrefix(url, "http://") &&
		!strings.HasPrefix(url, "https://") {
		return fmt.Errorf("url must start with http:// or https://")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "SABA-Agent/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("website returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return err
	}

	text := string(body)

	reScript := regexp.MustCompile(`(?is)<script[^>]*>.*?</script>|<style[^>]*>.*?</style>`)
	text = reScript.ReplaceAllString(text, "")

	reTags := regexp.MustCompile(`(?s)<[^>]*>`)
	text = reTags.ReplaceAllString(text, " ")

	text = html.UnescapeString(text)
	text = strings.Join(strings.Fields(text), " ")

	if len(text) > 10000 {
		text = text[:10000]
	}

	fmt.Println(text)

	return nil
}

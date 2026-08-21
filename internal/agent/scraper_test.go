package agent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebScraperTool(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
<html>
<head><title>SABA Test</title></head>
<body>
<h1>Hello from SABA</h1>
<p>Web scraping works.</p>
</body>
</html>
`))
	}))
	defer server.Close()

	tool := &WebScraperTool{}

	err := tool.Execute(context.Background(), map[string]interface{}{
		"url": server.URL,
	})

	if err != nil {
		t.Fatalf("scraper failed: %v", err)
	}
}

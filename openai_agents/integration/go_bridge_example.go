// Example integration boundary for SABA.
// This file is intentionally not wired into the repository automatically.
// It shows the small HTTP client boundary the Go gateway can use.

package integration

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
)

type AgentsRequest struct {
    Task  string `json:"task"`
    Mode  string `json:"mode"`
    Input string `json:"input,omitempty"`
}

type AgentsResponse struct {
    Status string `json:"status"`
    Mode   string `json:"mode"`
    Agent  string `json:"agent"`
    Output string `json:"output"`
}

func RunAgents(ctx context.Context, baseURL string, req AgentsRequest) (AgentsResponse, error) {
    payload, err := json.Marshal(req)
    if err != nil {
        return AgentsResponse{}, err
    }

    url := fmt.Sprintf("%s/run", baseURL)
    request, err := http.NewRequestWithContext(
        ctx, http.MethodPost, url, bytes.NewReader(payload),
    )
    if err != nil {
        return AgentsResponse{}, err
    }

    request.Header.Set("Content-Type", "application/json")

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        return AgentsResponse{}, err
    }
    defer response.Body.Close()

    if response.StatusCode < 200 || response.StatusCode >= 300 {
        return AgentsResponse{}, fmt.Errorf("agents service returned %s", response.Status)
    }

    var result AgentsResponse
    if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
        return AgentsResponse{}, err
    }

    return result, nil
}

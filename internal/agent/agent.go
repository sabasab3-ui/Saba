package agent

import (
	"context"
	"fmt"
	"strings"
)

type Agent struct {
	Name        string
	Model       string
	Memory      *Memory
	Tools       *ToolKit
	Language    string
	Description string
}

func NewAgent(name, model, language string) *Agent {
	tools := NewToolKit()

	// Register built-in SABA tools
	tools.Register(&InventoryTool{})
	tools.Register(&OrderTool{})
	tools.Register(&PaymentTool{})
	tools.Register(&WebScraperTool{})

	return &Agent{
		Name:        name,
		Model:       model,
		Memory:      NewMemory(),
		Tools:       tools,
		Language:    language,
		Description: "SABA - Smart Automation for Business in Africa",
	}
}

type Message struct {
	UserID    string
	Content   string
	Timestamp int64
	Language  string
}

type Response struct {
	Text       string
	Action     string
	Params     map[string]interface{}
	Confidence float64
}

func (a *Agent) Process(ctx context.Context, msg *Message) (*Response, error) {
	a.Memory.AddTurn(msg.UserID, msg.Content)

	text := strings.ToLower(strings.TrimSpace(msg.Content))

	// Simple task detection
	task := ""

	switch {
	case strings.Contains(text, "inventory"):
		task = "inventory"

	case strings.Contains(text, "order"):
		task = "orders"

	case strings.Contains(text, "payment"):
		task = "payments"

	case strings.Contains(text, "scrape"):
		task = "scraper"
	}

	// No task detected: normal response
	if task == "" {
		return &Response{
			Text:       fmt.Sprintf("SABA received: %s", msg.Content),
			Action:     "echo",
			Params:     map[string]interface{}{},
			Confidence: 0.5,
		}, nil
	}

	// Execute detected task
	result, err := a.ExecuteTask(ctx, task, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:       fmt.Sprintf("SABA executed %s: %v", task, result),
		Action:     task,
		Params:     map[string]interface{}{},
		Confidence: 0.9,
	}, nil
}

func (a *Agent) ExecuteTask(
	ctx context.Context,
	taskName string,
	params map[string]interface{},
) (map[string]interface{}, error) {

	tool, exists := a.Tools.Get(taskName)

	if !exists {
		return nil, fmt.Errorf("tool not found: %s", taskName)
	}

	err := tool.Execute(ctx, params)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status": "ok",
		"task":   taskName,
	}, nil
}

func (a *Agent) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"name":        a.Name,
		"model":       a.Model,
		"language":    a.Language,
		"status":      "active",
		"description": a.Description,
	}
}

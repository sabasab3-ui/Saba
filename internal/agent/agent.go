package agent

import (
	"context"
	"fmt"
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

	response := &Response{
		Text:       fmt.Sprintf("SABA received: %s", msg.Content),
		Action:     "echo",
		Confidence: 0.5,
	}

	return response, nil
}

func (a *Agent) ExecuteTask(ctx context.Context, taskName string, params map[string]interface{}) error {
	tool, exists := a.Tools.Get(taskName)
	if !exists {
		return fmt.Errorf("tool not found: %s", taskName)
	}

	return tool.Execute(ctx, params)
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

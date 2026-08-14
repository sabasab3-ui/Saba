package agent

import (
	"context"
	"fmt"
	"sync"
)

type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, params map[string]interface{}) error
}

type ToolKit struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

func NewToolKit() *ToolKit {
	return &ToolKit{
		tools: make(map[string]Tool),
	}
}

func (tk *ToolKit) Register(tool Tool) {
	tk.mu.Lock()
	defer tk.mu.Unlock()
	tk.tools[tool.Name()] = tool
}

func (tk *ToolKit) Get(name string) (Tool, bool) {
	tk.mu.RLock()
	defer tk.mu.RUnlock()

	tool, exists := tk.tools[name]
	return tool, exists
}

func (tk *ToolKit) List() []Tool {
	tk.mu.RLock()
	defer tk.mu.RUnlock()

	tools := make([]Tool, 0, len(tk.tools))
	for _, tool := range tk.tools {
		tools = append(tools, tool)
	}

	return tools
}

type InventoryTool struct{}

func (t *InventoryTool) Name() string {
	return "inventory"
}

func (t *InventoryTool) Description() string {
	return "Manage and query inventory levels"
}

func (t *InventoryTool) Execute(ctx context.Context, params map[string]interface{}) error {
	action, ok := params["action"].(string)
	if !ok {
		return fmt.Errorf("missing action parameter")
	}

	switch action {
	case "check", "update":
		return nil
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

type OrderTool struct{}

func (t *OrderTool) Name() string {
	return "orders"
}

func (t *OrderTool) Description() string {
	return "Create, view, and manage orders"
}

func (t *OrderTool) Execute(ctx context.Context, params map[string]interface{}) error {
	action, ok := params["action"].(string)
	if !ok {
		return fmt.Errorf("missing action parameter")
	}

	switch action {
	case "create", "list", "update":
		return nil
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

type PaymentTool struct{}

func (t *PaymentTool) Name() string {
	return "payments"
}

func (t *PaymentTool) Description() string {
	return "Process payments (MTN Mobile Money, Airtel Money)"
}

func (t *PaymentTool) Execute(ctx context.Context, params map[string]interface{}) error {
	provider, ok := params["provider"].(string)
	if !ok {
		return fmt.Errorf("missing provider parameter")
	}

	switch provider {
	case "mtn", "airtel":
		return nil
	default:
		return fmt.Errorf("unknown provider: %s", provider)
	}
}

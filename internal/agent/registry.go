package agent

import (
	"context"
	"fmt"
	"sort"
)

type ToolInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// NewDefaultToolKit creates the standard SABA business-tool package.
func NewDefaultToolKit() *ToolKit {
	tk := NewToolKit()

	tk.Register(&InventoryTool{})
	tk.Register(&OrderTool{})
	tk.Register(&PaymentTool{})

	return tk
}

// Catalog returns stable metadata for all registered tools.
func (tk *ToolKit) Catalog() []ToolInfo {
	tk.mu.RLock()
	defer tk.mu.RUnlock()

	items := make([]ToolInfo, 0, len(tk.tools))
	for name, tool := range tk.tools {
		items = append(items, ToolInfo{
			Name:        name,
			Description: tool.Description(),
			Category:    categoryFor(name),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return items
}

func categoryFor(name string) string {
	switch name {
	case "inventory", "orders":
		return "business"
	case "payments":
		return "finance"
	default:
		return "general"
	}
}

// Execute runs a registered tool by name.
func (tk *ToolKit) Execute(ctx context.Context, name string, params map[string]interface{}) error {
	tool, ok := tk.Get(name)
	if !ok {
		return fmt.Errorf("tool not found: %s", name)
	}

	return tool.Execute(ctx, params)
}

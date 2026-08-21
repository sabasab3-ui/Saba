package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/sabasab3-ui/saba/internal/database"
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

	db, err := database.OpenInventoryDB("inventory.db")
	if err != nil {
		return fmt.Errorf("open inventory database: %w", err)
	}
	defer db.Close()

	switch action {
	case "check":
		name, ok := params["name"].(string)
		if !ok || name == "" {
			return fmt.Errorf("missing product name")
		}

		product, err := db.CheckProduct(name)
		if err != nil {
			return err
		}

		if product == nil {
			return fmt.Errorf("product not found: %s", name)
		}

		return nil

	case "update":
		name, ok := params["name"].(string)
		if !ok || name == "" {
			return fmt.Errorf("missing product name")
		}

		v, ok := params["quantity"]
		if !ok {
			return fmt.Errorf("missing quantity")
		}

		var quantity int
		switch n := v.(type) {
		case int:
			quantity = n
		case float64:
			quantity = int(n)
		default:
			return fmt.Errorf("invalid quantity")
		}

		return db.UpdateQuantity(name, quantity)

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

	db, err := database.OpenOrderDB("orders.db")
	if err != nil {
		return fmt.Errorf("open orders database: %w", err)
	}
	defer db.Close()

	switch action {

	case "create":
		customer, ok := params["customer"].(string)
		if !ok || customer == "" {
			return fmt.Errorf("missing customer")
		}

		product, ok := params["product"].(string)
		if !ok || product == "" {
			return fmt.Errorf("missing product")
		}

		v, ok := params["quantity"]
		if !ok {
			return fmt.Errorf("missing quantity")
		}

		var quantity int

		switch n := v.(type) {
		case int:
			quantity = n
		case float64:
			quantity = int(n)
		default:
			return fmt.Errorf("invalid quantity")
		}

		_, err := db.CreateOrder(customer, product, quantity)
		if err != nil {
			return err
		}

		return nil

	case "list":
		_, err := db.ListOrders()
		if err != nil {
			return err
		}

		return nil

	case "update":
		v, ok := params["id"]
		if !ok {
			return fmt.Errorf("missing order id")
		}

		var id int64

		switch n := v.(type) {
		case int:
			id = int64(n)
		case int64:
			id = n
		case float64:
			id = int64(n)
		default:
			return fmt.Errorf("invalid order id")
		}

		status, ok := params["status"].(string)
		if !ok || status == "" {
			return fmt.Errorf("missing status")
		}

		return db.UpdateStatus(id, status)

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

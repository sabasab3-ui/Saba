package agent

import (
	"context"
	"github.com/sabasab3-ui/saba/internal/database"
	"testing"
)

func TestAgentToolsAndMemory(t *testing.T) {
	agent := NewAgent("SABA", "local", "en")
	db, err := database.OpenInventoryDB("inventory.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.AddProduct("test-product", 10); err != nil {
		t.Fatal(err)
	}

	// Test memory
	agent.Memory.AddTurn("user1", "Hello SABA")

	turns := agent.Memory.GetTurns("user1", 10)

	if len(turns) != 1 {
		t.Fatalf("expected 1 memory turn, got %d", len(turns))
	}

	// Test inventory tool
	_, err = agent.ExecuteTask(
		context.Background(),
		"inventory",
		map[string]interface{}{
			"action": "check",
			"name":   "test-product",
		},
	)

	if err != nil {
		t.Fatalf("inventory tool failed: %v", err)
	}

	// Test payment tool
	_, err = agent.ExecuteTask(
		context.Background(),
		"payments",
		map[string]interface{}{
			"provider": "mtn",
		},
	)

	if err != nil {
		t.Fatalf("payment tool failed: %v", err)
	}
}

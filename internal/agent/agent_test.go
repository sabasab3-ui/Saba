package agent

import (
	"context"
	"testing"
)

func TestAgentToolsAndMemory(t *testing.T) {
	agent := NewAgent("SABA", "local", "en")

	// Test memory
	agent.Memory.AddTurn("user1", "Hello SABA")

	turns := agent.Memory.GetTurns("user1", 10)

	if len(turns) != 1 {
		t.Fatalf("expected 1 memory turn, got %d", len(turns))
	}

	// Test inventory tool
	err := agent.ExecuteTask(
		context.Background(),
		"inventory",
		map[string]interface{}{
			"action": "check",
		},
	)

	if err != nil {
		t.Fatalf("inventory tool failed: %v", err)
	}

	// Test payment tool
	err = agent.ExecuteTask(
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

package agent

import (
	"sync"
	"time"
)

type Memory struct {
	mu    sync.RWMutex
	turns map[string][]Turn
}

type Turn struct {
	UserID    string
	UserMsg   string
	AgentMsg  string
	Timestamp int64
}

func NewMemory() *Memory {
	return &Memory{
		turns: make(map[string][]Turn),
	}
}

func (m *Memory) AddTurn(userID, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.turns[userID] = append(m.turns[userID], Turn{
		UserID:    userID,
		UserMsg:   message,
		Timestamp: time.Now().Unix(),
	})
}

func (m *Memory) GetTurns(userID string, limit int) []Turn {
	m.mu.RLock()
	defer m.mu.RUnlock()

	turns := m.turns[userID]

	if limit <= 0 || len(turns) <= limit {
		return append([]Turn(nil), turns...)
	}

	return append([]Turn(nil), turns[len(turns)-limit:]...)
}

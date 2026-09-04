package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Agent represents an AI agent instance
type Agent struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Type        string    `db:"type"`        // research, analysis, business, coding, reasoning
	Description string    `db:"description"`
	Status      string    `db:"status"`      // active, inactive, training
	Version     string    `db:"version"`
	Capabilities string   `db:"capabilities"` // JSON array
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// AgentSession represents a conversation session with an agent
type AgentSession struct {
	ID        string    `db:"id"`
	AgentID   string    `db:"agent_id"`
	UserID    string    `db:"user_id"`
	Status    string    `db:"status"`       // active, completed, failed
	Context   string    `db:"context"`      // JSON session context
	StartedAt time.Time `db:"started_at"`
	EndedAt   *time.Time `db:"ended_at"`
	MessagesCount int   `db:"messages_count"`
}

// CreateAgentsTable creates the agents table
func (db *DB) CreateAgentsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS agents (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			description TEXT,
			status TEXT DEFAULT 'active',
			version TEXT,
			capabilities TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.conn.ExecContext(ctx, query)
	return err
}

// CreateAgentSessionsTable creates the agent_sessions table
func (db *DB) CreateAgentSessionsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS agent_sessions (
			id TEXT PRIMARY KEY,
			agent_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			status TEXT DEFAULT 'active',
			context TEXT,
			started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			ended_at DATETIME,
			messages_count INTEGER DEFAULT 0,
			FOREIGN KEY (agent_id) REFERENCES agents(id)
		)
	`
	_, err := db.conn.ExecContext(ctx, query)
	return err
}

// AddAgent adds a new agent
func (db *DB) AddAgent(ctx context.Context, agent *Agent) error {
	agent.CreatedAt = time.Now()
	agent.UpdatedAt = time.Now()

	query := `
		INSERT INTO agents (id, name, type, description, status, version, capabilities, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.conn.ExecContext(ctx, query,
		agent.ID, agent.Name, agent.Type, agent.Description, agent.Status,
		agent.Version, agent.Capabilities, agent.CreatedAt, agent.UpdatedAt,
	)
	return err
}

// GetAgent retrieves an agent by ID
func (db *DB) GetAgent(ctx context.Context, id string) (*Agent, error) {
	agent := &Agent{}
	query := `SELECT id, name, type, description, status, version, capabilities, created_at, updated_at FROM agents WHERE id = ?`
	err := db.conn.QueryRowContext(ctx, query, id).Scan(
		&agent.ID, &agent.Name, &agent.Type, &agent.Description, &agent.Status,
		&agent.Version, &agent.Capabilities, &agent.CreatedAt, &agent.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("agent not found")
		}
		return nil, err
	}
	return agent, nil
}

// ListAgentsByType retrieves agents by type
func (db *DB) ListAgentsByType(ctx context.Context, agentType string) ([]*Agent, error) {
	query := `SELECT id, name, type, description, status, version, capabilities, created_at, updated_at FROM agents WHERE type = ? AND status = 'active' ORDER BY created_at DESC`
	rows, err := db.conn.QueryContext(ctx, query, agentType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*Agent
	for rows.Next() {
		agent := &Agent{}
		err := rows.Scan(
			&agent.ID, &agent.Name, &agent.Type, &agent.Description, &agent.Status,
			&agent.Version, &agent.Capabilities, &agent.CreatedAt, &agent.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		agents = append(agents, agent)
	}
	return agents, rows.Err()
}

// CreateAgentSession creates a new agent session
func (db *DB) CreateAgentSession(ctx context.Context, session *AgentSession) error {
	session.StartedAt = time.Now()
	query := `
		INSERT INTO agent_sessions (id, agent_id, user_id, status, context, started_at, messages_count)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.conn.ExecContext(ctx, query,
		session.ID, session.AgentID, session.UserID, session.Status,
		session.Context, session.StartedAt, session.MessagesCount,
	)
	return err
}

// GetAgentSession retrieves a session by ID
func (db *DB) GetAgentSession(ctx context.Context, id string) (*AgentSession, error) {
	session := &AgentSession{}
	query := `SELECT id, agent_id, user_id, status, context, started_at, ended_at, messages_count FROM agent_sessions WHERE id = ?`
	err := db.conn.QueryRowContext(ctx, query, id).Scan(
		&session.ID, &session.AgentID, &session.UserID, &session.Status,
		&session.Context, &session.StartedAt, &session.EndedAt, &session.MessagesCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}
	return session, nil
}

package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Automation represents an automation workflow
type Automation struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Type        string    `db:"type"`         // workflow, schedule, trigger
	Status      string    `db:"status"`       // active, inactive, paused
	Config      string    `db:"config"`       // JSON configuration
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	LastRunAt   *time.Time `db:"last_run_at"`
}

// AutomationRun represents a single execution of an automation
type AutomationRun struct {
	ID             string    `db:"id"`
	AutomationID   string    `db:"automation_id"`
	Status         string    `db:"status"`        // success, failed, running
	StartedAt      time.Time `db:"started_at"`
	CompletedAt    *time.Time `db:"completed_at"`
	Duration       int64     `db:"duration"`      // milliseconds
	ErrorMessage   string    `db:"error_message"`
	Output         string    `db:"output"`        // JSON result
}

// CreateAutomationsTable creates the automations table
func (db *DB) CreateAutomationsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS automations (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			type TEXT NOT NULL,
			status TEXT DEFAULT 'active',
			config TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_run_at DATETIME
		)
	`
	_, err := db.conn.ExecContext(ctx, query)
	return err
}

// CreateAutomationRunsTable creates the automation_runs table
func (db *DB) CreateAutomationRunsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS automation_runs (
			id TEXT PRIMARY KEY,
			automation_id TEXT NOT NULL,
			status TEXT NOT NULL,
			started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			completed_at DATETIME,
			duration INTEGER,
			error_message TEXT,
			output TEXT,
			FOREIGN KEY (automation_id) REFERENCES automations(id)
		)
	`
	_, err := db.conn.ExecContext(ctx, query)
	return err
}

// AddAutomation adds a new automation
func (db *DB) AddAutomation(ctx context.Context, automation *Automation) error {
	automation.CreatedAt = time.Now()
	automation.UpdatedAt = time.Now()

	query := `
		INSERT INTO automations (id, name, description, type, status, config, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.conn.ExecContext(ctx, query,
		automation.ID, automation.Name, automation.Description, automation.Type,
		automation.Status, automation.Config, automation.CreatedAt, automation.UpdatedAt,
	)
	return err
}

// GetAutomation retrieves an automation by ID
func (db *DB) GetAutomation(ctx context.Context, id string) (*Automation, error) {
	auto := &Automation{}
	query := `SELECT id, name, description, type, status, config, created_at, updated_at, last_run_at FROM automations WHERE id = ?`
	err := db.conn.QueryRowContext(ctx, query, id).Scan(
		&auto.ID, &auto.Name, &auto.Description, &auto.Type, &auto.Status,
		&auto.Config, &auto.CreatedAt, &auto.UpdatedAt, &auto.LastRunAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("automation not found")
		}
		return nil, err
	}
	return auto, nil
}

// ListAutomations retrieves all automations
func (db *DB) ListAutomations(ctx context.Context) ([]*Automation, error) {
	query := `SELECT id, name, description, type, status, config, created_at, updated_at, last_run_at FROM automations ORDER BY created_at DESC`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var automations []*Automation
	for rows.Next() {
		auto := &Automation{}
		err := rows.Scan(
			&auto.ID, &auto.Name, &auto.Description, &auto.Type, &auto.Status,
			&auto.Config, &auto.CreatedAt, &auto.UpdatedAt, &auto.LastRunAt,
		)
		if err != nil {
			return nil, err
		}
		automations = append(automations, auto)
	}
	return automations, rows.Err()
}

// RecordAutomationRun records an automation execution
func (db *DB) RecordAutomationRun(ctx context.Context, run *AutomationRun) error {
	query := `
		INSERT INTO automation_runs (id, automation_id, status, started_at, completed_at, duration, error_message, output)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.conn.ExecContext(ctx, query,
		run.ID, run.AutomationID, run.Status, run.StartedAt, run.CompletedAt,
		run.Duration, run.ErrorMessage, run.Output,
	)
	return err
}

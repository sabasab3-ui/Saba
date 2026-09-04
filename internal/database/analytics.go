package database

import (
	"context"
	"database/sql"
	"time"
)

// Analytics represents system analytics
type Analytics struct {
	ID              string    `db:"id"`
	MetricType      string    `db:"metric_type"` // agents_used, automations_run, tasks_completed, errors
	MetricValue     int64     `db:"metric_value"`
	Dimensions      string    `db:"dimensions"`  // JSON: {agent_type, country, industry}
	RecordedAt      time.Time `db:"recorded_at"`
}

// CreateAnalyticsTable creates the analytics table
func (db *DB) CreateAnalyticsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS analytics (
			id TEXT PRIMARY KEY,
			metric_type TEXT NOT NULL,
			metric_value INTEGER NOT NULL,
			dimensions TEXT,
			recorded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_metric_time (metric_type, recorded_at)
		)
	`
	_, err := db.conn.ExecContext(ctx, query)
	return err
}

// RecordMetric records an analytics metric
func (db *DB) RecordMetric(ctx context.Context, analytics *Analytics) error {
	analytics.RecordedAt = time.Now()
	query := `
		INSERT INTO analytics (id, metric_type, metric_value, dimensions, recorded_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := db.conn.ExecContext(ctx, query,
		analytics.ID, analytics.MetricType, analytics.MetricValue,
		analytics.Dimensions, analytics.RecordedAt,
	)
	return err
}

// GetMetricsByType retrieves metrics by type over a time range
func (db *DB) GetMetricsByType(ctx context.Context, metricType string, from, to time.Time) ([]*Analytics, error) {
	query := `
		SELECT id, metric_type, metric_value, dimensions, recorded_at
		FROM analytics
		WHERE metric_type = ? AND recorded_at BETWEEN ? AND ?
		ORDER BY recorded_at DESC
	`
	rows, err := db.conn.QueryContext(ctx, query, metricType, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*Analytics
	for rows.Next() {
		m := &Analytics{}
		err := rows.Scan(
			&m.ID, &m.MetricType, &m.MetricValue, &m.Dimensions, &m.RecordedAt,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, rows.Err()
}

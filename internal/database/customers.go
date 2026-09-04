package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Customer represents a business customer
type Customer struct {
	ID            string    `db:"id"`
	Name          string    `db:"name"`
	Email         string    `db:"email"`
	Phone         string    `db:"phone"`
	Country       string    `db:"country"`
	Industry      string    `db:"industry"`
	CompanySize   string    `db:"company_size"`
	Status        string    `db:"status"` // active, inactive, pending
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	LastContactAt *time.Time `db:"last_contact_at"`
}

// CreateCustomersTable creates the customers table
func (db *DB) CreateCustomersTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS customers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			phone TEXT,
			country TEXT,
			industry TEXT,
			company_size TEXT,
			status TEXT DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_contact_at DATETIME
		)
	`
	_, err := db.conn.ExecContext(ctx, query)
	return err
}

// AddCustomer adds a new customer
func (db *DB) AddCustomer(ctx context.Context, customer *Customer) error {
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	query := `
		INSERT INTO customers (id, name, email, phone, country, industry, company_size, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.conn.ExecContext(ctx, query,
		customer.ID, customer.Name, customer.Email, customer.Phone,
		customer.Country, customer.Industry, customer.CompanySize,
		customer.Status, customer.CreatedAt, customer.UpdatedAt,
	)
	return err
}

// GetCustomer retrieves a customer by ID
func (db *DB) GetCustomer(ctx context.Context, id string) (*Customer, error) {
	customer := &Customer{}
	query := `SELECT id, name, email, phone, country, industry, company_size, status, created_at, updated_at, last_contact_at FROM customers WHERE id = ?`
	err := db.conn.QueryRowContext(ctx, query, id).Scan(
		&customer.ID, &customer.Name, &customer.Email, &customer.Phone,
		&customer.Country, &customer.Industry, &customer.CompanySize,
		&customer.Status, &customer.CreatedAt, &customer.UpdatedAt, &customer.LastContactAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, err
	}
	return customer, nil
}

// ListCustomers retrieves all customers
func (db *DB) ListCustomers(ctx context.Context) ([]*Customer, error) {
	query := `SELECT id, name, email, phone, country, industry, company_size, status, created_at, updated_at, last_contact_at FROM customers ORDER BY created_at DESC`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		customer := &Customer{}
		err := rows.Scan(
			&customer.ID, &customer.Name, &customer.Email, &customer.Phone,
			&customer.Country, &customer.Industry, &customer.CompanySize,
			&customer.Status, &customer.CreatedAt, &customer.UpdatedAt, &customer.LastContactAt,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, rows.Err()
}

// UpdateCustomer updates an existing customer
func (db *DB) UpdateCustomer(ctx context.Context, customer *Customer) error {
	customer.UpdatedAt = time.Now()
	query := `
		UPDATE customers
		SET name = ?, email = ?, phone = ?, country = ?, industry = ?, company_size = ?, status = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := db.conn.ExecContext(ctx, query,
		customer.Name, customer.Email, customer.Phone,
		customer.Country, customer.Industry, customer.CompanySize,
		customer.Status, customer.UpdatedAt, customer.ID,
	)
	return err
}

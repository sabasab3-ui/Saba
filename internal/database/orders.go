package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type OrderDB struct {
	db *sql.DB
}

type Order struct {
	ID        int64  `json:"id"`
	Customer  string `json:"customer"`
	Product   string `json:"product"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func OpenOrderDB(path string) (*OrderDB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS orders (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	customer TEXT NOT NULL,
	product TEXT NOT NULL,
	quantity INTEGER NOT NULL DEFAULT 1,
	status TEXT NOT NULL DEFAULT 'pending',
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL
)
`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("create orders table: %w", err)
	}

	return &OrderDB{db: db}, nil
}

func (o *OrderDB) CreateOrder(customer, product string, quantity int) (*Order, error) {
	if customer == "" {
		return nil, fmt.Errorf("customer is required")
	}

	if product == "" {
		return nil, fmt.Errorf("product is required")
	}

	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than zero")
	}

	now := time.Now().Unix()

	result, err := o.db.Exec(`
INSERT INTO orders
(customer, product, quantity, status, created_at, updated_at)
VALUES (?, ?, ?, 'pending', ?, ?)
`, customer, product, quantity, now, now)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Order{
		ID:        id,
		Customer:  customer,
		Product:   product,
		Quantity:  quantity,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (o *OrderDB) GetOrder(id int64) (*Order, error) {
	order := &Order{}

	err := o.db.QueryRow(`
SELECT id, customer, product, quantity, status, created_at, updated_at
FROM orders
WHERE id = ?
`, id).Scan(
		&order.ID,
		&order.Customer,
		&order.Product,
		&order.Quantity,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *OrderDB) ListOrders() ([]Order, error) {
	rows, err := o.db.Query(`
SELECT id, customer, product, quantity, status, created_at, updated_at
FROM orders
ORDER BY id DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order

	for rows.Next() {
		var order Order

		err := rows.Scan(
			&order.ID,
			&order.Customer,
			&order.Product,
			&order.Quantity,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *OrderDB) UpdateStatus(id int64, status string) error {
	if status == "" {
		return fmt.Errorf("status is required")
	}

	result, err := o.db.Exec(`
UPDATE orders
SET status = ?, updated_at = ?
WHERE id = ?
`, status, time.Now().Unix(), id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("order not found: %d", id)
	}

	return nil
}

func (o *OrderDB) Close() error {
	return o.db.Close()
}

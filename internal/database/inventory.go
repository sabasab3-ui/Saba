package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type InventoryDB struct {
	db *sql.DB
}

type Product struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func OpenInventoryDB(path string) (*InventoryDB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS products (
id INTEGER PRIMARY KEY AUTOINCREMENT,
name TEXT NOT NULL UNIQUE,
quantity INTEGER NOT NULL DEFAULT 0,
created_at INTEGER NOT NULL,
updated_at INTEGER NOT NULL
)
`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("create products table: %w", err)
	}

	return &InventoryDB{db: db}, nil
}

func (i *InventoryDB) AddProduct(name string, quantity int) error {
	now := time.Now().Unix()

	_, err := i.db.Exec(`
INSERT INTO products (name, quantity, created_at, updated_at)
VALUES (?, ?, ?, ?)
ON CONFLICT(name) DO UPDATE SET
quantity = excluded.quantity,
updated_at = excluded.updated_at
`, name, quantity, now, now)

	return err
}

func (i *InventoryDB) CheckProduct(name string) (*Product, error) {
	product := &Product{}

	err := i.db.QueryRow(`
SELECT id, name, quantity, created_at, updated_at
FROM products
WHERE name = ?
`, name).Scan(
		&product.ID,
		&product.Name,
		&product.Quantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (i *InventoryDB) UpdateQuantity(name string, quantity int) error {
	result, err := i.db.Exec(`
UPDATE products
SET quantity = ?, updated_at = ?
WHERE name = ?
`, quantity, time.Now().Unix(), name)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("product not found: %s", name)
	}

	return nil
}

func (i *InventoryDB) Close() error {
	return i.db.Close()
}

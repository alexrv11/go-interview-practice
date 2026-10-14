package main

import (
    "fmt"
	"database/sql"
	"errors"
	"context"

	_ "github.com/mattn/go-sqlite3"
)

var ErrNotFound = errors.New("product not found")

// Product represents a product in the inventory system
type Product struct {
	ID       int64
	Name     string
	Price    float64
	Quantity int
	Category string
}

// ProductStore manages product operations
type ProductStore struct {
	db *sql.DB
}

// NewProductStore creates a new ProductStore with the given database connection
func NewProductStore(db *sql.DB) *ProductStore {
	return &ProductStore{db: db}
}

// InitDB sets up a new SQLite database and creates the products table
func InitDB(dbPath string) (*sql.DB, error) {
	// TODO: Open a SQLite database connection
	// TODO: Create the products table if it doesn't exist
	// The table should have columns: id, name, price, quantity, category
	conn, err := sql.Open("sqlite3", dbPath)
    if err != nil {
    	return nil, err
    }
    
    sqlStmt := `
	CREATE TABLE IF NOT EXISTS products (
        id          INTEGER PRIMARY KEY,
        name        TEXT    NOT NULL CHECK (length(trim(name)) > 0),
        price       REAL NOT NULL CHECK (price >= 0),
        quantity    INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
        category    TEXT
    ) STRICT;
    
    CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
	`
	_, err = conn.Exec(sqlStmt)
	if err != nil {
		
		return nil, err
	}
	
	return conn, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// TODO: Insert the product into the database
	// TODO: Update the product.ID with the database-generated ID
	res, err := ps.db.ExecContext(context.Background(),
		`INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)`,
		product.Name, product.Price, product.Quantity, product.Category,
	)
	if err != nil {
		return err
	}
	
	lastID, err := res.LastInsertId() 
	
	if err != nil {
	    return err
	}
	
	product.ID = lastID 
	
	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// TODO: Query the database for a product with the given ID
	// TODO: Return a Product struct populated with the data or an error if not found
	var p Product
	ctx := context.Background()
	err := ps.db.QueryRowContext(ctx, `
		SELECT p.id, p.name, p.price, p.quantity, p.category 
		FROM products p
		WHERE p.id = ?`, id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
 
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	
	if err != nil {
		return nil, fmt.Errorf("get product %d: %w", id, err)
	}
	
	return &p, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// TODO: Update the product in the database
	// TODO: Return an error if the product doesn't exist
	ctx := context.Background()
	_, err := ps.db.ExecContext(ctx, `
		UPDATE products
		SET name = ?, price = ?, quantity = ?, category = ?
		WHERE id = ?`,
		product.Name, product.Price, product.Quantity, product.Category, product.ID,
	)
	
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound // no row matched → product doesn't exist
	}
	
	if err != nil {
		return err
	}
	
	
	
	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	// TODO: Return an error if the product doesn't exist
	ctx := context.Background()
	_, err := ps.db.ExecContext(ctx, `DELETE FROM products WHERE id = ?`, id)
	
	if err != nil {
		return err
	}
	
	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// TODO: Query the database for products
	// TODO: If category is not empty, filter by category
	// TODO: Return a slice of Product pointers
	ctx := context.Background()
	query := `
		SELECT p.id, p.name, p.price, p.quantity, p.category
		FROM products p
		`
	if len(category) > 0 {
	    query = `
		    SELECT p.id, p.name, p.price, p.quantity, p.category
		    FROM products p
		    WHERE p.category = ? COLLATE NOCASE
		`     
	}
	
	rows, err := ps.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
 
	var out []*Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category); err != nil {
			return nil, err
		}
		out = append(out, &p)
	}
	return out, rows.Err()
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// TODO: Start a transaction
	// TODO: For each product ID in the updates map, update its quantity
	// TODO: If any update fails, roll back the transaction
	// TODO: Otherwise, commit the transaction
	ctx := context.Background()
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // undoes everything unless Commit succeeds
 
	// Prepare once, execute many times inside the transaction.
	stmt, err := tx.PrepareContext(ctx, `UPDATE products SET quantity = ? WHERE id = ?`)
	
	if err != nil {
		return fmt.Errorf("error prepare context: %w", err)
	}
	
	defer stmt.Close()
 
	for id, q := range updates {
		res, err := stmt.ExecContext(ctx, q, id)
		if err != nil {
			return fmt.Errorf("error execute query %d: %w", id, err) // e.g. CHECK quantity >= 0
		}
		
		n, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("error rows affected %d: %w", id, err)
		}
		
		if n == 0 {
			// Returning here triggers the deferred Rollback,
			// undoing every update already done in this loop.
			return fmt.Errorf("product %d: %w", id, ErrNotFound)
		}
	}
 
	return tx.Commit()
}

func main() {
	// Optional: you can write code here to test your implementation
}

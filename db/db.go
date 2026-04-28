// Database connection and configuration module
package db

import (
	"database/sql"

	// Import postgres driver
	_ "github.com/lib/pq"
)

// Database wraps the SQL connection pool
type Database struct {
	db *sql.DB
}

// NewDatabase establishes a connection to the PostgreSQL database
func NewDatabase() (*Database, error) {
	// Open connection using provided credentials and database name
	db, err := sql.Open("postgres", "postgresql://root:password@localhost:5433/go-chat?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

// Close terminates the database connection
func (d Database) Close() {
	d.db.Close()
}

// GetDB returns the underlying SQL database instance
func (d *Database) GetDB() *sql.DB {
	return d.db
}

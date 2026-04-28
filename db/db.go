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
	db, err := sql.Open("postgres", "postgresql://postgres:Nikhras%40122@localhost:5432/go_realtime_chat_app?sslmode=disable")
	// Open connection using provided credentials and database name
	db, err = sql.Open("postgres", "postgresql://root:password@localhost:5433/go-chat?sslmode=disable")
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

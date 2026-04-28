// Database repository for user operations
package user

import (
	"context"
	"database/sql"
)

// DBTX defines a generic interface for database operations (supports DB and Tx)
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// repository implements the Repository interface for user data
type repository struct {
	db DBTX
}

// NewRepository initializes a new user repository
func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}  

// CreateUser inserts a new user record into the database
func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastInsertId int
	query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id"
	
	// Execute insertion and retrieve the newly generated ID
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&lastInsertId)
	if err != nil {
		return &User{}, err
	}

	user.ID = int64(lastInsertId)
	return user, nil
}

// GetUserByEmail retrieves a user record by their email address
func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}
	query := "SELECT id, email, username, password FROM users WHERE email = $1"
	
	// Execute query and map results to the user struct
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return &User{}, nil
	}

	return &u, nil
}

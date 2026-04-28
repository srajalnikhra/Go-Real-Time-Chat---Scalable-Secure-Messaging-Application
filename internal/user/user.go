// Core user domain models and interfaces
package user

import "context"

// User represents a user record in the database
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUserReq represents the request payload for creating a new user
type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUserRes represents the response payload after a user is successfully created
type CreateUserRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// LoginUserReq represents the request payload for user authentication
type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUserRes represents the response payload after successful login
type LoginUserRes struct {
	accessToken string
	ID          string `json:"id"`
	Username    string `json:"username"`
}

// Repository defines the interface for database operations related to users
type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

// Service defines the interface for user-related business logic
type Service interface {
	CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error)
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error)
}

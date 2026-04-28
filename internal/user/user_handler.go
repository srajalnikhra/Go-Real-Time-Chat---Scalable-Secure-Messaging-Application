// HTTP handlers for user-related endpoints
package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler wraps the user service to handle HTTP requests
type Handler struct {
	Service  
}

// NewHandler initializes a new user handler
func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

// CreateUser handles the registration of a new user
func (h *Handler) CreateUser(c *gin.Context) {
	var u CreateUserReq
	
	// Parse and validate request body
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delegate user creation to the service layer
	res, err := h.Service.CreateUser(c.Request.Context(), &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, res)
}

// Login handles user authentication and session creation
func (h *Handler) Login(c *gin.Context) {
	var user LoginUserReq

	// Parse login credentials from request body
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate user via service layer
	u, err := h.Service.Login(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set JWT token as an HTTP-only cookie
	c.SetCookie("jwt", u.accessToken, 3600, "/", "localhost", false, true)

	// Prepare and return user details (excluding token)
	res := &LoginUserRes{
		Username: u.Username,
		ID:       u.ID,
	}
	c.JSON(http.StatusOK, res)
}

// Logout clears the user session by removing the JWT cookie
func (h *Handler) Logout(c *gin.Context) {
	// Invalidate the JWT cookie
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

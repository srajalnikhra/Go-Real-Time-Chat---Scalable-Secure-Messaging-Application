// HTTP handlers for WebSocket room and client operations
package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Handler wraps the WebSocket hub for HTTP endpoints
type Handler struct {
	hub *Hub
}

// NewHandler initializes a new WebSocket HTTP handler
func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

// CreateRoomReq represents the request payload for creating a chat room
type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateRoom handles the creation of a new chat room
func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	
	// Parse room details from request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Register the new room in the hub
	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, req)
}

// upgrader is used to upgrade HTTP connections to WebSocket protocol
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// JoinRoom upgrades the HTTP connection and adds the client to the specified room
func (h *Handler) JoinRoom(c *gin.Context) {
	// Upgrade the HTTP connection to a WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract connection parameters
	roomID := c.Param("roomId")
	clientID := c.Query("userId")
	username := c.Query("username")

	// Initialize the new client connection
	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	// Prepare join announcement message
	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	// Register client and announce their arrival
	h.hub.Register <- cl
	h.hub.Broadcast <- m

	// Start message pumps in background goroutines
	go cl.writeMessage()
	cl.readMessage(h.hub)
}

// RoomRes represents the response payload for a single chat room
type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetRooms returns a list of all active chat rooms
func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]RoomRes, 0)

	// Collect all rooms from the hub
	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	c.JSON(http.StatusOK, rooms)
}

// ClientRes represents the response payload for a single connected client
type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// GetClients returns a list of all users connected to a specific room
func (h *Handler) GetClients(c *gin.Context) {
	var clients []ClientRes
	roomId := c.Param("roomId")

	// Return empty list if room does not exist
	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	// Collect all clients in the specified room
	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}

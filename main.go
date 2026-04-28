// Application entry point
package main

import (
	"log"
	"server/db"
	"server/internal/user"
	"server/internal/ws"
	"server/router"
)

func main() {
	// Initialize database connection
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	// Initialize user domain layers
	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)
 
	// Initialize WebSocket hub and handler
	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)
	
	// Start WebSocket hub in a background goroutine
	go hub.Run()

	// Register routes and start HTTP server
	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:8080")
}

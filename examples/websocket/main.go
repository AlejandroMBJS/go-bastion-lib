package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alejandrombjs/go-bastion-lib/pkg/bastion"
	"github.com/alejandrombjs/go-bastion-lib/pkg/middleware"
	"github.com/alejandrombjs/go-bastion-lib/pkg/response"
	"github.com/alejandrombjs/go-bastion-lib/pkg/router"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for this example. In production, restrict this.
		return true
	},
}

func main() {
	cfg := bastion.DefaultConfig()
	cfg.Port = 8086 // Use a different port for this example

	app := bastion.NewApp(cfg)
	r := app.Router()

	// Global middlewares
	app.Use(
		middleware.DefaultLogging(),
		middleware.DefaultRecovery(),
	)

	// --- WebSocket Endpoint ---
	r.GET("/ws", func(ctx *router.Context) {
		conn, err := upgrader.Upgrade(ctx.ResponseWriter(), ctx.Request().Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v", err)
			response.Error(ctx, http.StatusInternalServerError, "websocket_upgrade_failed", "Could not open websocket connection")
			return
		}
		defer conn.Close()

		log.Printf("WebSocket client connected from %s", conn.RemoteAddr())

		// Simple echo loop
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error from %s: %v", conn.RemoteAddr(), err)
				break
			}
			log.Printf("Received from %s: %s", conn.RemoteAddr(), p)

			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Printf("WebSocket write error to %s: %v", conn.RemoteAddr(), err)
				break
			}
		}
		log.Printf("WebSocket client disconnected from %s", conn.RemoteAddr())
	})

	// --- Health Check (regular HTTP) ---
	r.GET("/health", func(ctx *router.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"status": "ok", "type": "http"})
	})

	log.Printf("WebSocket example server starting on :%d", cfg.Port)
	if err := app.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

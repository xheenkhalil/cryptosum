package api

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/cryptosum/backend/internal/config"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Hub maintains the set of active clients
type Hub struct {
	Clients    map[*websocket.Conn]bool
	Broadcast  chan []byte
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	mu         sync.Mutex
}

var WsHub = &Hub{
	Broadcast:  make(chan []byte),
	Register:   make(chan *websocket.Conn),
	Unregister: make(chan *websocket.Conn),
	Clients:    make(map[*websocket.Conn]bool),
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			log.Println("New WebSocket client connected")
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Close()
			}
			h.mu.Unlock()
			log.Println("WebSocket client disconnected")
		case message := <-h.Broadcast:
			h.mu.Lock()
			for client := range h.Clients {
				if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
					client.Close()
					delete(h.Clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func SetupWebSockets(app *fiber.App) {
	go WsHub.Run()

	// MVP: Listen to all trade results and broadcast to all connected clients
	go func() {
		rdb := redis.NewClient(&redis.Options{
			Addr: config.AppConfig.RedisURL,
		})
		
		pubsub := rdb.PSubscribe(context.Background(), "trade_results:*")
		ch := pubsub.Channel()
		for msg := range ch {
			out := fmt.Sprintf(`{"type":"TRADE_RESULT", "payload": %s}`, msg.Payload)
			WsHub.Broadcast <- []byte(out)
		}
	}()

	// Middleware to check if it's a websocket request
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/stream", websocket.New(func(c *websocket.Conn) {
		WsHub.Register <- c
		defer func() {
			WsHub.Unregister <- c
		}()

		// Keep connection alive and listen for messages
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			// We can process incoming client ws messages here if needed
			log.Printf("Received WS message: %s", msg)
		}
	}))
}

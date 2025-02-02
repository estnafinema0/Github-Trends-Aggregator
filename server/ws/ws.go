package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/estnafinema0/Github-Trends-Aggregator/server/models"
	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Hub manages connections and message broadcasting
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Broadcast sends updates to all clients
func (h *Hub) Broadcast(l *log.Logger, repos []models.Repository) {
	message, err := json.Marshal(repos)
	if err != nil {
		l.Printf("Error marshalling data: %v\n", err)
		return
	}
	h.broadcast <- message
}

func (h *Hub) Run(l *log.Logger) {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			l.Println("New client connected")
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				l.Println("Client disconnected")
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow connections from any source (for MVP)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeWs upgrades HTTP connection to WebSocket
func ServeWs(hub *Hub, l *log.Logger, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l.Printf("Error upgrading to WebSocket: %v\n", err)
		return
	}

	client := &Client{conn: conn, send: make(chan []byte, 256)}
	hub.register <- client

	go client.writePump()
	go client.readPump(hub)
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		// Processing incoming messages is not required (in this MVP)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for {
		message, ok := <-c.send
		if !ok {
			// Channel closed
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

// Ensure ServeWs returns an http.HandlerFunc
func ServeWsHandler(hub *Hub, l *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, l, w, r)
	}
}

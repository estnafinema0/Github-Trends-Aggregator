package ws

import (
	"encoding/json"
	"log"
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
	l          *log.Logger
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
func (h *Hub) Broadcast(repos []models.Repository) {
	message, err := json.Marshal(repos)
	if err != nil {
		h.l.Printf("Error marshalling data: %v\n", err)
		return
	}
	h.broadcast <- message
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.l.Println("New client connected")
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.l.Println("Client disconnected")
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

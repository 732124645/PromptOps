package ws

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Hub keeps track of connected SDK/UI clients and broadcasts hot-reload events.
type Hub struct {
	mu        sync.RWMutex
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte, 64),
	}
}

func (h *Hub) Run() {
	for msg := range h.broadcast {
		h.mu.RLock()
		for c := range h.clients {
			if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
				_ = c.Close()
			}
		}
		h.mu.RUnlock()
	}
}

// Broadcast queues a message for all connected clients; drops it if the queue is full.
func (h *Hub) Broadcast(msg []byte) {
	select {
	case h.broadcast <- msg:
	default:
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWS upgrades an HTTP request to a WebSocket connection for hot-reload events.
func (h *Hub) ServeWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	go func() {
		defer func() {
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
			_ = conn.Close()
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
}

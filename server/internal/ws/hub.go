package ws

import (
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client describes a single connected hot-reload session — an SDK instance or
// a browser tab. The metadata is collected from the WebSocket handshake.
type Client struct {
	ID          string    `json:"id"`
	ClientType  string    `json:"client_type"`
	App         string    `json:"app"`
	Namespace   string    `json:"namespace"`
	RemoteAddr  string    `json:"remote_addr"`
	UserAgent   string    `json:"user_agent"`
	ConnectedAt time.Time `json:"connected_at"`
}

// Hub keeps track of connected SDK/UI clients and broadcasts hot-reload events.
type Hub struct {
	mu        sync.RWMutex
	clients   map[*websocket.Conn]*Client
	broadcast chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]*Client),
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

// Clients returns a snapshot of the currently connected clients, newest first.
func (h *Hub) Clients() []Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]Client, 0, len(h.clients))
	for _, cl := range h.clients {
		out = append(out, *cl)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ConnectedAt.After(out[j].ConnectedAt)
	})
	return out
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWS upgrades an HTTP request to a WebSocket connection for hot-reload
// events, recording who connected from the handshake query parameters.
func (h *Hub) ServeWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	namespace := c.Query("namespace")
	if namespace == "" {
		namespace = c.Query("env")
	}
	meta := &Client{
		ID:          uuid.NewString(),
		ClientType:  c.DefaultQuery("client", "unknown"),
		App:         c.Query("app"),
		Namespace:   namespace,
		RemoteAddr:  c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		ConnectedAt: time.Now(),
	}

	h.mu.Lock()
	h.clients[conn] = meta
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

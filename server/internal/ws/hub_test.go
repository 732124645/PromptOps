package ws

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// newTestHub starts a hub behind an httptest server exposing /ws.
func newTestHub(t *testing.T) (*Hub, *httptest.Server) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	hub := NewHub()
	go hub.Run()

	r := gin.New()
	r.GET("/ws", hub.ServeWS)
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return hub, srv
}

// waitClients polls hub.Clients() until it reaches want or times out.
func waitClients(hub *Hub, want int) []Client {
	for i := 0; i < 100; i++ {
		clients := hub.Clients()
		if len(clients) == want {
			return clients
		}
		time.Sleep(5 * time.Millisecond)
	}
	return hub.Clients()
}

func TestHubTracksClientMetadata(t *testing.T) {
	hub, srv := newTestHub(t)
	wsURL := strings.Replace(srv.URL, "http", "ws", 1) +
		"/ws?client=node-sdk&namespace=prod&app=checkout"

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	clients := waitClients(hub, 1)
	if len(clients) != 1 {
		t.Fatalf("expected 1 connected client, got %d", len(clients))
	}
	got := clients[0]
	if got.ClientType != "node-sdk" {
		t.Errorf("client_type = %q, want node-sdk", got.ClientType)
	}
	if got.Namespace != "prod" {
		t.Errorf("namespace = %q, want prod", got.Namespace)
	}
	if got.App != "checkout" {
		t.Errorf("app = %q, want checkout", got.App)
	}
	if got.ID == "" || got.ConnectedAt.IsZero() {
		t.Errorf("expected ID and ConnectedAt to be set, got %+v", got)
	}

	// Disconnecting removes the client from the registry.
	conn.Close()
	if clients := waitClients(hub, 0); len(clients) != 0 {
		t.Fatalf("expected client removed after disconnect, got %d", len(clients))
	}
}

func TestHubDefaultsClientType(t *testing.T) {
	hub, srv := newTestHub(t)
	wsURL := strings.Replace(srv.URL, "http", "ws", 1) + "/ws"

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	clients := waitClients(hub, 1)
	if len(clients) != 1 || clients[0].ClientType != "unknown" {
		t.Fatalf("expected one client with client_type 'unknown', got %+v", clients)
	}
}

func TestHubBroadcastReachesClient(t *testing.T) {
	hub, srv := newTestHub(t)
	wsURL := strings.Replace(srv.URL, "http", "ws", 1) + "/ws?client=node-sdk"

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	waitClients(hub, 1)

	hub.Broadcast([]byte(`{"type":"update","key":"demo"}`))

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read broadcast: %v", err)
	}
	if !strings.Contains(string(msg), `"key":"demo"`) {
		t.Fatalf("unexpected broadcast payload: %s", msg)
	}
}

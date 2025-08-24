package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/application"
	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	clients     map[string]*Client
	upgrader    websocket.Upgrader
	locationSvc application.LocationService
	mu          sync.Mutex
}

func NewWebSocketHandler(locationSvc application.LocationService) *WebsocketHandler {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	var clients = make(map[string]*Client)
	return &WebsocketHandler{
		clients:     clients,
		upgrader:    upgrader,
		locationSvc: locationSvc,
	}
}

func (ws *WebsocketHandler) WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection: ", err)
		return
	}

	// TODO: securely handle this with zero-trust principle
	driverID := r.Header.Get("X-User-ID")
	if driverID == "" {
		return
	}

	client := &Client{
		ID:          driverID,
		Conn:        conn,
		Send:        make(chan any),
		locationSvc: ws.locationSvc,
	}

	ws.mu.Lock()
	ws.clients[driverID] = client
	ws.mu.Unlock()

	log.Printf("Driver %s connected", r.RemoteAddr)

	client.Start()
}

func (ws *WebsocketHandler) NotifyClient(driverID string, payload any) error {
	ws.mu.Lock()
	client, ok := ws.clients[driverID]
	ws.mu.Unlock()

	if !ok {
		return fmt.Errorf("driver %s not connected", driverID)
	}

	select {
	case client.Send <- payload:
		return nil
	default:
		return fmt.Errorf("send buffer full for driver %s", driverID)
	}
}

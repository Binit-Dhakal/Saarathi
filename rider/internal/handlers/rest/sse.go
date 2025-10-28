package rest

import (
	"fmt"
	"net/http"
	"sync"
)

type TripUpdateHandler struct {
	connections map[string]*Client
	mu          sync.RWMutex
}

func NewTripUpdateHandler() *TripUpdateHandler {
	return &TripUpdateHandler{
		connections: map[string]*Client{},
	}
}

func (t *TripUpdateHandler) TripUpdate(w http.ResponseWriter, r *http.Request) {
	tripID := r.URL.Query().Get("tripId")
	if tripID == "" {
		http.Error(w, "missing tripId", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	client := &Client{
		id:      tripID,
		send:    make(chan any, 10),
		done:    make(chan struct{}),
		flusher: flusher,
		w:       w,
	}
	fmt.Printf("User connected for tripID: %v\n", tripID)

	go client.writePump()

	t.mu.Lock()
	t.connections[tripID] = client
	t.mu.Unlock()

	<-r.Context().Done()

	t.mu.Lock()
	delete(t.connections, tripID)
	t.mu.Unlock()

	close(client.done)
	close(client.send)
}

func (t *TripUpdateHandler) NotifyRider(tripID string, payload any) {
	t.mu.RLock()
	client, ok := t.connections[tripID]
	t.mu.RUnlock()

	fmt.Println(tripID, payload)

	if ok {
		select {
		case client.send <- payload:
			return
		default:
			// Client channel is backed up; drop the message
		}
	}
}
